package appserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AppServer struct {
	router  *gin.Engine
	httpsrv *http.Server
}

func NewAppServer() *AppServer {
	srv := &AppServer{router: gin.Default()}
	srv.buildRoutes()
	return srv
}

func (srv *AppServer) buildRoutes() {

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
