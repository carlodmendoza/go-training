FROM golang:1.16-alpine
WORKDIR /final-project/client
COPY go.mod ./
COPY main.go ./
COPY models ./models
COPY utils ./utils
RUN go build -o /client
CMD ["/client"]