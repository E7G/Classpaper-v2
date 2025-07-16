@ECHO OFF
go install github.com/cratonica/2goarray@latest
2goarray IconData main < favicon.ico > resources.go
