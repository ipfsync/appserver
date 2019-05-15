package appserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ipfsync/ipfsync/core"

	"github.com/gin-gonic/gin"
)

type Message interface {
}

type AppServer struct {
	router    *gin.Engine
	httpsrv   *http.Server
	api       *core.Api
	cron      *appCron
	wsClients map[*wsClient]bool
}

func NewAppServer(api *core.Api) *AppServer {
	srv := &AppServer{
		router:    gin.Default(),
		api:       api,
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

	srv.router.GET("/peers", func(c *gin.Context) {
		peers, _ := srv.api.Peers()
		var str string
		for _, p := range peers {
			str += p.Address().String() + "\n"
		}
		c.String(http.StatusOK, str)
	})
}

func (srv *AppServer) Start() {

	srv.httpsrv = &http.Server{
		Addr:    ":8080",
		Handler: srv.router,
	}

	// Start cron jobs
	srv.cron.start()

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
}

func (srv *AppServer) unregisterWsClient(c *wsClient) {
	delete(srv.wsClients, c)
}

func (srv *AppServer) Broadcast(msg interface{}) {
	for client := range srv.wsClients {
		client.send <- msg
	}
}
