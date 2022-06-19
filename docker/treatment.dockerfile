FROM golang:1.18

COPY . /app

WORKDIR /app/treatment-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_data", "-dbUser", "root", "-dbPassword", "root"]