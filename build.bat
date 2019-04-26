set GOPATH=D:\gowork
set CGO_ENABLED=0

go install

set GOOS=linux
set GOARCH=amd64
go install -ldflags "-s -w"

set GOOS=darwin
set GOARCH=amd64
go install -ldflags "-s -w"