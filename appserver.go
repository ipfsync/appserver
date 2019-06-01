package appserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/websocket"

	"github.com/ipfsync/ipfsync/core/api"

	"github.com/gin-gonic/gin"
)

type MessageError struct {
	Code    int
	Message string
}

type MessageCmd struct {
	Id   string
	Cmd  string
	Data map[string]interface{}
}

type MessageReply struct {
	Id    string
	Ok    bool
	Data  map[string]interface{}
	Error MessageError
}

type MessageBroadcast struct {
	Event string
	Data  map[string]interface{}
}

type AppServer struct {
	router    *gin.Engine
	httpsrv   *http.Server
	api       *api.Api
	cfg       *viper.Viper
	cron      *appCron
	wsClients map[*wsClient]bool
}

func NewAppServer(api *api.Api, cfg *viper.Viper) *AppServer {
	srv := &AppServer{
		router:    gin.Default(),
		api:       api,
		cfg:       cfg,
		wsClients: make(map[*wsClient]bool),
	}
	cron := newCron(srv)
	srv.cron = cron
	srv.buildRoutes()
	return srv
}

func (srv *AppServer) buildRoutes() {
	srv.router.GET("/ws", func(c *gin.Context) {
		srv.wsServe(c.Writer, c.Request)
	})

	srv.router.GET("/test", func(c *gin.Context) {
		_, _ = srv.api.NewCollection("", "")
	})
}

func (srv *AppServer) Start() {

	// TODO: automatically choose an unused port
	srv.httpsrv = &http.Server{
		Addr:    ":8080",
		Handler: srv.router,
	}

	go func() {
		if err := srv.httpsrv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()
}

func (srv *AppServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.httpsrv.Shutdown(ctx)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (srv *AppServer) wsServe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := newWsClient(srv, conn)
	srv.registerWsClient(client)

	go client.readPump()
	go client.writePump()
}

func (srv *AppServer) registerWsClient(c *wsClient) {
	srv.wsClients[c] = true

	// Start cron jobs
	srv.cron.start()
}

func (srv *AppServer) unregisterWsClient(c *wsClient) {
	delete(srv.wsClients, c)

	if len(srv.wsClients) == 0 {
		srv.cron.stop()
	}
}

func (srv *AppServer) Broadcast(msg interface{}) {
	for client := range srv.wsClients {
		client.send <- msg
	}
}
