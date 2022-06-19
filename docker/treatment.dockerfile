FROM golang:1.18

COPY . /app

WORKDIR /app/treatment-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_treatment_db", "-dbUser", "root", "-dbPassword", "root"]