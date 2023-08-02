# Latest golang image on apline linux
FROM golang:latest as builder

# Env variables
ENV GOOS linux
ENV CGO_ENABLED 0

# Work directory
WORKDIR /lingo-docker

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Building the application
RUN go build -o lingo-docker

# Fetching the latest nginx image
FROM alpine:latest as production

# Certificates
RUN apk add --no-cache ca-certificates

# Copying built assets from builder
COPY --from=builder lingo-docker .

# Starting our application
CMD ./lingo-docker

# Exposing server port
EXPOSE 5000