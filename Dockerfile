# Build stage
FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o kwait ./cmd/kwait

# Final stage
FROM gcr.io/distroless/base-debian11

COPY --from=build /app/kwait /usr/local/bin/kwait

ENTRYPOINT ["/usr/local/bin/kwait"]
