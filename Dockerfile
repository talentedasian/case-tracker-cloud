FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=$TARGETOS go build -o inmate

FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/inmate /inmate

ENTRYPOINT ["/inmate/go-gcp"]
