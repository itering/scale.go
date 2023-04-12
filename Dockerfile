FROM golang:1.18

WORKDIR app

COPY . .

RUN go mod download

CMD go test -v ./...