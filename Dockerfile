FROM golang:1.18-alpine as builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -ldflags "-s -w" -o drop

FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=builder /src/drop /app/applicaiton

ENV dsn ""
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/applicaiton"]

