module github.com/ipfsync/appserver

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/websocket v1.4.0
	github.com/ipfsync/ipfsync v0.0.0
	github.com/robfig/cron v1.1.0
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager

replace github.com/ipfsync/ipfsync => ../ipfsync
