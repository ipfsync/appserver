package appserver

import (
	"log"

	"github.com/robfig/cron"
)

// appCron is for scheduled jobs for UI data pushing only.
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
	log.Println("Cron job stopped")
}

func (c *appCron) sendBroadcast(event string, data map[string]interface{}) {
	msg := &MessageBroadcast{
		Event: event,
		Data:  data,
	}

	c.srv.Broadcast(msg)
}

func (c *appCron) peers() {

	peers, changed, err := c.srv.api.Peers()
	if err != nil {
		log.Println("Unable to fetch peers. Error %v", err)
		return
	}

	if !changed {
		return
	}

	c.sendBroadcast("peers", map[string]interface{}{
		"peers": peers,
	})
}
