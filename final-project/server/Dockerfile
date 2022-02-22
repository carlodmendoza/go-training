FROM golang:1.16-alpine
WORKDIR /final-project/server
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY data ./data
COPY models ./models
COPY redis ./redis
COPY utils ./utils
RUN go build -o /server
EXPOSE 8080
CMD ["/server"]