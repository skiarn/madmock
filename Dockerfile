# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Download and Build mad-mock application inside the container.
RUN go get github.com/skiarn/mad-mock/src/madmock
RUN go install github.com/skiarn/mad-mock/src/madmock

# Run the application when the container starts.
#                            Enter host i  wish to mock. Change this to your system ip. ex: -u=127.0.0.1:9090
ENTRYPOINT /go/bin/madmock -u=http://github.com

# Document that the service listens on port 8080.
EXPOSE 8080
