FROM golang:1.16-alpine

WORKDIR /final-project/server

COPY go.mod go.sum ../
COPY server/*.go ./
COPY server/internal/ ./internal/
COPY server/pkg/ ./pkg/
COPY server/storage/ ./storage/

RUN go build -o /server

EXPOSE 8080

ENTRYPOINT ["/server"]