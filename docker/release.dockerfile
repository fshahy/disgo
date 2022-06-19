FROM golang:1.18

COPY . /app

WORKDIR /app/release-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "opd_data", "-dbUser", "root", "-dbPassword", "Root@1985"]