FROM golang:1.18

COPY . /app

WORKDIR /app/release-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_release_db", "-dbUser", "root", "-dbPassword", "root"]