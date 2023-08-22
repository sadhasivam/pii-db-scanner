FROM golang:1.21 as build 

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o scannercli ./cmd/scannercli

# Use a lightweight alpine image as a second stage.
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from the previous stage.
COPY --from=build /app/scannercli /usr/local/bin/

# Set the binary as the entrypoint of the container.
ENTRYPOINT ["scannercli"]

# Command line arguments can be provided after the image name in the `docker run` command.
CMD ["--help"]