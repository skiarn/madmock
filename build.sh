go get golang.org/x/net/websocket 
GOOS=windows GOARCH=amd64 go build -o madmock.exe main.go;
GOOS=linux GOARCH=amd64 go build -o linux_madmock main.go;
GOOS=darwin GOARCH=386 go build -o mac_madmock main.go;
