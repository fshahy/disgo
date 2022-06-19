FROM golang:1.18

COPY . /app

WORKDIR /app/registration-service/cmd
RUN go build ./...

ENTRYPOINT [ "./cmd", "-dbName", "disgo_data", "-dbUser", "root", "-dbPassword", "root"]