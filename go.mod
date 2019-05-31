module github.com/ipfsync/appserver

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/websocket v1.4.0
	github.com/ipfsync/ipfsync v0.0.0
	github.com/robfig/cron v1.1.0
	github.com/spf13/viper v1.4.0
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
)

replace github.com/ipfsync/appserver => ../appserver

replace github.com/ipfsync/ipfsmanager => ../ipfsmanager

replace github.com/ipfsync/ipfsync => ../ipfsync

replace github.com/ipfsync/common => ../common

replace github.com/ipfsync/resource => ../resource
