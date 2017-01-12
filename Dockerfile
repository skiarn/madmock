# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ENV TARGET github.com

# Download and Build mad-mock application inside the container.
RUN go get github.com/skiarn/madmock
RUN go install github.com/skiarn/madmock

# Run the application when the container starts.
# Enter host to mock. Change this to your system ip. ex: -u=127.0.0.1:9090
ENTRYPOINT /go/bin/madmock -u=$TARGET

# Service listens on port 9988.
EXPOSE 9988
