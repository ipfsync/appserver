package appserver

import (
	"log"

	"github.com/robfig/cron"
)

type appCron struct {
	srv *AppServer
	rc  *cron.Cron
}

func newCron(srv *AppServer) *appCron {
	c := &appCron{srv: srv}

	// Cron jobs
	rc := cron.New()
	c.rc = rc

	c.buildJobs()

	return c
}

func (c *appCron) buildJobs() {
	// Peers data
	err := c.rc.AddFunc("@every 1s", c.peers)
	if err != nil {
		log.Printf("Unable to add job: %v", err)
	}
}

func (c *appCron) start() {
	c.rc.Start()
	log.Println("Cron job started")
}

func (c *appCron) stop() {
	c.rc.Stop()
}

func (c *appCron) peers() {
	peers, err := c.srv.api.Peers()
	if err != nil {
		return
	}
	var addr []string
	for _, p := range peers {
		addr = append(addr, p.Address().String())
	}

	msg := &MessageBroadcast{
		Event: "peers",
		Data: map[string]interface{}{
			"addresses": addr,
		},
	}

	c.srv.Broadcast(msg)
}
