################
# BUILD BINARY #
################
# golang:1.18.2-alpine3.16
# FROM golang@sha256:5dbae0a41feefa0d9e0e455ddf82672bf439baac577650e6c6b81e98b6efe97a as builder
# FROM alpine:3.18 AS build
FROM golang:1.22.0-alpine3.18 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/go-charai
COPY . .

RUN echo $PWD && ls -la

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

#CMD go build -v
# go build command with the -ldflags="-w -s" option to produce a smaller binary file by stripping debug information and symbol tables. 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/go-charai .

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.18

RUN apk update

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
COPY --from=builder /go/bin/go-charai /go/bin/go-charai

ENTRYPOINT ["/go/bin/go-charai"]