FROM golang:1.16-alpine
WORKDIR /final-project/server
COPY go.mod ./
COPY main.go ./
COPY data ./data
COPY models ./models
COPY utils ./utils
RUN go build -o /server
EXPOSE 8080
CMD ["/server"]