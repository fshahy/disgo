package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"release-service"
	"shared"
)

func main() {
	var (
		showHelp     bool
		showVersion  bool
		serverListen string
		natsServers  string
		dbUser       string
		dbPass       string
		dbName       string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: release-service [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.StringVar(&serverListen, "listen", "0.0.0.0:9093", "Network host:port to listen on")
	flag.StringVar(&natsServers, "nats", "nats://nats-server:4222", "List of NATS Servers to connect")
	flag.StringVar(&dbUser, "dbUser", "", "Database username")
	flag.StringVar(&dbPass, "dbPassword", "", "Database password")
	flag.StringVar(&dbName, "dbName", "", "Database name")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATS Microservices OPD Sample - Release Service v%s\n", release.Version)
		os.Exit(0)
	}
	log.Printf("Starting NATS Microservices OPD Sample - Release Service version %s", release.Version)

	// Register new component within the system.
	comp := shared.NewComponent("release-service")

	// Connect to NATS and setup discovery subscriptions.
	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Database.
	err = comp.SetupConnectionToDB("mysql", dbUser+":"+dbPass+"@tcp(release-db:3306)/"+dbName)
	if err != nil {
		log.Fatal(err)
	}

	s := release.Server{
		Component: comp,
	}

	err = s.ListenAndServe(serverListen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for HTTP requests on %v", serverListen)
	runtime.Goexit()
}
