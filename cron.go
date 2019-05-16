package appserver

import (
	"log"
	"sort"
	"time"

	net "github.com/libp2p/go-libp2p-net"

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

type peerinfo struct {
	Address   string
	Direction net.Direction
	Latency   time.Duration
}

var peersinfo []peerinfo

func (c *appCron) peers() {
	peers, err := c.srv.api.Peers()
	if err != nil {
		return
	}
	peersinfo = nil
	for _, p := range peers {
		l, _ := p.Latency()
		peersinfo = append(peersinfo, peerinfo{
			Address:   p.Address().String(),
			Direction: p.Direction(),
			Latency:   l,
		})
	}

	// Sort
	sort.Slice(peersinfo, func(i, j int) bool {
		return peersinfo[i].Address < peersinfo[j].Address
	})

	msg := &MessageBroadcast{
		Event: "peers",
		Data: map[string]interface{}{
			"peers": peersinfo,
		},
	}

	c.srv.Broadcast(msg)
}
