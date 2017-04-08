# mad-mock
**The application is to be used as a http mock.**
Purpose of tool/server is to ensure communication to another resource is read only and fakeable without entering the external system. 

It works like this, request to other systems passes through this tool and this tool will cache the results. If the request is a modifiable call such as POST/PUT/DELETE the request is never executed against the target system but instead a faked response is sent back.

Its possible to modify and create fake responses using a web gui or modify raw files on disk.

Try it out using dockerhub!
`docker run -e TARGET=github.com -e PORT=7070 -p 8080:7070 --rm -it skiarn/madmock bash`

## Manually
### How to get started ###

* Install go - https://golang.org/dl/
* git clone http://github.com/skiarn/mad-mock

Configuration
* linux/mac: ```export GOPATH=/project/path/madmock```
* windows: ```set GOPATH=/project/path/madmock```

### Build application ###
go get golang.org/x/net/websocket
go install madmock

## Run application
Application is located in $GOPATH/bin.
* For help ```./madmock --help ```

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

## Docker
Try it out using dockerhub!
`docker run -e TARGET=github.com -e PORT=7070 -p 8080:7070 --rm -it skiarn/madmock bash`

Or build yourself:
* git clone http://github.com/skiarn/mad-mock
* cd mad-mock
* Make sure you edit: ```ENTRYPOINT /go/bin/madmock -u=http://github.com``` change ```-u=http://github.com``` to the base url of what system or page you wish to mock.
* ```docker build -t madmock . ```
* ```docker run -d -p 8080:8080 madmock```
