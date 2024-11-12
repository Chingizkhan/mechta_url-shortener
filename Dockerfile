## Build the application from source
#FROM golang:1.22.1 AS build-stage
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#RUN go mod download
#
#COPY . .
#
#RUN CGO_ENABLED=0 go build -o /main
#RUN chmod +x /main  # Устанавливаем права на выполнение
#
## Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...
#
## Deploy the application binary into a lean image
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
#
#WORKDIR /
#
#COPY --from=build-stage /main /main
#
#EXPOSE 8888
#
#USER nonroot:nonroot
#
#ENTRYPOINT ["/main"]


## Stage 1: Build stage
#FROM golang:1.22.3-alpine AS build
#
## Set the working directory
#WORKDIR /app
#
## Copy and download dependencies
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copy the source code
#COPY . .
#
## Build the Go application
#RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .
#RUN chmod +x /myapp  # Устанавливаем права на выполнение
#
## Stage 2: Final stage
#FROM alpine:edge
#
## Set the working directory
#WORKDIR /app
#
## Copy the binary from the build stage
#COPY --from=build /app/myapp .
#
## Set the timezone and install CA certificates
#RUN apk --no-cache add ca-certificates tzdata
#
## Set the entrypoint command
#ENTRYPOINT ["/app/myapp"]

FROM golang:1.23.3

WORKDIR /app

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app/app.go

EXPOSE 8888

CMD ["./main"]