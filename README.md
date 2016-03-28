# mad-mock
**The application is to be used as a http mock.**

## Manually
### How to get started ###

* Install go - https://golang.org/dl/
* git clone http://github.com/skiarn/mad-mock

Configuration
* linux/mac: ```export GOPATH=/project/path/madmock```
* windows: ```set GOPATH=/project/path/madmock```

### Build application ###
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
* git clone http://github.com/skiarn/mad-mock
* cd mad-mock
* Make sure you edit: ```ENTRYPOINT /go/bin/madmock -u=http://github.com``` change ```-u=http://github.com``` to the base url of what system or page you wish to mock.
* ```docker build -t madmock . ```
* ```docker run -d -p 8080:8080 madmock```
