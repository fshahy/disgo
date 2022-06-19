module release-service

go 1.18

replace shared => ../shared

require (
	shared v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30
	github.com/nats-io/nuid v1.0.1
)
