FROM golang:1.17-alpine
WORKDIR /app
COPY . /app
RUN go mod download

RUN go build cmd

CMD ["./main"]