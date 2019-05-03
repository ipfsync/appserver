module github.com/ipfsync/appserver

go 1.12

require (
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/ipfsync/ipfsmanager v0.0.0
	github.com/ipfsync/ipfsync v0.0.0
	github.com/ugorji/go v1.1.4 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager

replace github.com/ipfsync/ipfsync => ../ipfsync
