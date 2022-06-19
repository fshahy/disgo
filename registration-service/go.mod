module registration-service

go 1.18

replace shared => ../shared

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/nats-io/nats.go v1.16.0
	github.com/nats-io/nuid v1.0.1
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.11.12 // indirect
	github.com/minio/highwayhash v1.0.1 // indirect
	github.com/nats-io/jwt/v2 v2.0.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
)
