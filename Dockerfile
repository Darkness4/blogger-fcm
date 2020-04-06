# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.14 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/Darkness4/blogger-fcm
COPY . .

# Build the command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -v -o blogger-fcm

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/Darkness4/blogger-fcm/blogger-fcm /blogger-fcm

# Run the web service on container startup.
CMD ["/blogger-fcm"]