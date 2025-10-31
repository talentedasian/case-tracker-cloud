FROM golang:1.25.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o inmate

FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/inmate /inmate

# Set entrypoint
CMD ["/inmate"]
