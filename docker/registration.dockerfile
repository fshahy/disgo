FROM golang:1.18

COPY . /app

WORKDIR /app/registration-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_registration_db", "-dbUser", "root", "-dbPassword", "root"]