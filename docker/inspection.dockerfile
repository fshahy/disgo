FROM golang:1.18

COPY . /app

WORKDIR /app/inspection-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_inspection_db", "-dbUser", "root", "-dbPassword", "root"]