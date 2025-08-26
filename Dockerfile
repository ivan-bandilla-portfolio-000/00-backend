FROM golang:1.24.5 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# ensure ca-certificates are available in builder
RUN apt-get update && apt-get install -y ca-certificates
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /app/main ./...

# use minimal final image
FROM scratch
# copy CA certs so HTTPS calls work
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/main /app/main
EXPOSE 8080
ENTRYPOINT ["/app/main"]