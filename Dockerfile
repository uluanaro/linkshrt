FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o linkshrt .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/linkshrt .
EXPOSE 8080
CMD ["./linkshrt"]