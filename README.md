# mad-mock

### How to get started ###

* Install go - https://golang.org/dl/

* Configuration
linux/mac: export GOPATH=/pjoject/path/
windows: set GOPATH=/pjoject/path/

### Build application ###
go install mad-mock

## Run application
* For help ```./mad-mock --help ```

```
Usage of ./mad-mock:
  -d="mad-mock-store": Directory path to mock data and config files.
  -p=8080: What port the mock should run on.
  -u="": Base url to system to be mocked (request will be fetched once and stored locally).
```

Example, logs data to server.log and mocks 127.0.0.1:9090.
Open your web browser and visit http://localhost:8080/mock to view mocked urls and their request data.
```
nohup ./mad-mock -u=127.0.0.1:9090 > server.log 2>&1 &
```

**The application is to be used as a http mock.**
