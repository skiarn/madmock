# mad-mock
**The application is to be used as a http mock.**
Purpose of tool/server is to ensure communication to another resource is read only and fakeable without entering the external system.

It works like this, request to other systems passes through this tool and this tool will cache the results. If the request is a modifiable call such as POST/PUT/DELETE the request is never executed against the target system but instead a faked response is sent back.

Its possible to modify and create fake responses using a web gui or modify raw files on disk.

Visit http://localhost:<port>/mock to view a web interface to view and modify payloads. 
The browser wil try open mock monitor view automatically.

Build: `go run main.go`

**Try it out using dockerhub!**
`docker run -p 8080:10 -it --rm skiarn/madmock -u=apple.com -p=10`

Need to setup or modify messages?
curl -H "Content-Type: application/json" -X POST -d '{"uri":"/examplePost","method":"GET","contenttype":"text/html; charset=utf-8","status":200, "body": "body payload..."}' localhost:8080/mock/api/mock/

```
Usage of ./madmock:
  -d="mad-mock-store": Directory path to mock data and config files.
  -p=8080: What port the mock should run on.
  -u="": Base url to system to be mocked (request will be fetched once and stored locally).
```

Example, logs data to server.log and mocks 127.0.0.1:9090.
Open your web browser and visit http://localhost:8080/mock to view mocked urls and their request data.
```
nohup ./madmock -u=127.0.0.1:9090 > server.log 2>&1 &
```

Or want to build from source using docker,
`docker run -it -p 9988:9988 --rm golang:1.9 sh -c "mkdir -p /go/src/github.com/skiarn; git clone http://github.com/skiarn/madmock /go/src/github.com/skiarn/madmock; go get golang.org/x/net/websocket; go run /go/src/github.com/skiarn/madmock/main.go -u=google.se; bash"`
