FROM golang:1.23.3

WORKDIR /app

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app/app.go

EXPOSE 8888

CMD ["./main"]