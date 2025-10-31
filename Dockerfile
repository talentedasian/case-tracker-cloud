FROM golang:1.25.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /inmate main.go

CMD ["/inmate"]