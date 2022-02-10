FROM golang:1.17-buster AS builder

WORKDIR /app

COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -o go-mfa .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/go-mfa .

CMD ["./go-mfa"]
