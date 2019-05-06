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

type AppServer struct {
	router    *gin.Engine
	httpsrv   *http.Server
	api       *core.Api
	wsClients map[*WsClient]bool
}

func NewAppServer(api *core.Api) *AppServer {
	srv := &AppServer{router: gin.Default(), api: api}
	srv.buildRoutes()
	return srv
}

func (srv *AppServer) buildRoutes() {
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
}

func (srv *AppServer) wsServe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	srv.registerWsClient(&wsClient{srv: srv, conn: conn})
}

func (srv *AppServer) registerWsClient(c *wsClient) {
	srv.wsClients[c] = true
}

func (srv *AppServer) unregisterWsClient(c *wsClient) {
	delete(srv.wsClients, c)
}
