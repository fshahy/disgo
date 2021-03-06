package inspection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"shared"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

const (
	Version = "0.1.0"
)

// Server is a component.
type Server struct {
	*shared.Component
}

// Listen to patient registration events
func (s *Server) ListenRegisterEvents() error {
	nc := s.NATS()
	nc.Subscribe("patient.register", func(msg *nats.Msg) {
		var req *shared.RegistrationEvent
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}

		log.Printf("New Patient Registration Event received for PatientID %d with Token  %d\n",
			req.ID, req.Token)

		// Insert data to the database
		db := s.DB()

		insForm, err := db.Prepare("INSERT INTO patient_registrations(id, token) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(req.ID, req.Token)

	})

	return nil

}

// HandlePending processes requests to view pending inspections.
func (s *Server) HandlePending(w http.ResponseWriter, r *http.Request) {
	// Retrieve pending inspections from the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM patient_registrations")
	if err != nil {
		panic(err.Error())
	}

	type allRegistrations []shared.RegistrationEvent
	var registrations = allRegistrations{}

	for selDB.Next() {
		var newRegistration shared.RegistrationEvent
		var id int
		var token uint64
		err = selDB.Scan(&id, &token)
		if err != nil {
			panic(err.Error())
		}
		newRegistration.ID = id
		newRegistration.Token = token
		registrations = append(registrations, newRegistration)
	}

	fmt.Println(registrations)
	json.NewEncoder(w).Encode(registrations)
}

// HandleRecord processes patient inspection record requests.
func (s *Server) HandleRecord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var inspection *shared.InspectionRequest
	err = json.Unmarshal(body, &inspection)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert data to the database
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO inspection_details(id, time, observations, medication, tests, notes) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(inspection.ID, inspection.Time, inspection.Observations, inspection.Medication, inspection.Tests, inspection.Notes)

	// Remove the entry from pending inspections table if it exists
	removeData, err := db.Prepare("DELETE FROM patient_registrations WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	removeData.Exec(inspection.ID)

	// Tag the request with an ID for tracing in the logs.
	inspection.RequestID = nuid.Next()
	fmt.Println(inspection)

	// Publish event to the NATS server
	nc := s.NATS()

	inspection_event := shared.InspectionEvent{inspection.ID, inspection.Medication, inspection.Tests, inspection.Notes}
	reg_event, err := json.Marshal(inspection_event)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("requestID:%s - Publishing inspection event with patientID %d\n", inspection.RequestID, inspection.ID)
	// Publishing the message to NATS Server
	nc.Publish("patient.treatment", reg_event)

	json.NewEncoder(w).Encode(inspection_event)
}

// HandleHistory processes requests to view inspection history.
func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["id"]
	// Read data from the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM inspection_details WHERE ID=?", patientID)
	if err != nil {
		panic(err.Error())
	}

	type allInspections []shared.InspectionRequest
	var inspections = allInspections{}

	for selDB.Next() {
		var newInspection shared.InspectionRequest
		var id int
		var time, observations, medication, tests, notes string
		err = selDB.Scan(&id, &time, &observations, &medication, &tests, &notes)
		if err != nil {
			panic(err.Error())
		}
		newInspection.ID = id
		newInspection.Time = time
		newInspection.Observations = observations
		newInspection.Medication = medication
		newInspection.Tests = tests
		newInspection.Notes = notes
		inspections = append(inspections, newInspection)
	}

	fmt.Println(inspections)
	json.NewEncoder(w).Encode(inspections)
}

// HandleHomeLink respond with the version of the service
func (s *Server) HandleHomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Inspection Service v%s\n", Version))
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func (s *Server) ListenAndServe(addr string) error {

	// Start listening to patient registration events
	s.ListenRegisterEvents()

	r := mux.NewRouter()
	router := r.PathPrefix("/opd/inspection/").Subrouter()

	// Handle base path requests
	// GET /opd/inspection
	router.HandleFunc("/", s.HandleHomeLink)

	// Handle inspection record requests
	// POST /opd/inspection/record/{id}
	router.HandleFunc("/record", s.HandleRecord).Methods("POST")

	// Handle history view requests
	// GET /opd/inspection/history/{id}
	router.HandleFunc("/history/{id}", s.HandleHistory).Methods("GET")

	// Handle pending inspections view requests
	// GET /opd/inspection/pending
	router.HandleFunc("/pending", s.HandlePending).Methods("GET")

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go srv.Serve(l)

	return nil
}
