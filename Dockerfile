# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Build the mad-mock command inside the container.
RUN go get github.com/skiarn/mad-mock/src/mad-mock
RUN go install github.com/skiarn/mad-mock/src/mad-mock

# Run the outyet command by default when the container starts.
#                            Enter host i  wish to mock. Change this to your system ip. ex: 127.0.0.1:9090
ENTRYPOINT /go/bin/mad-mock -u=http://github.com

# Document that the service listens on port 8080.
EXPOSE 8080
