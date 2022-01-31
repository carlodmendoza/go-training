FROM golang:1.16-alpine
WORKDIR /contacts
COPY main.go /contacts
COPY go.mod /contacts
RUN go build -o server 
EXPOSE 8080
CMD ["./server"]
