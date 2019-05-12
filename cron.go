package appserver

import "github.com/robfig/cron"

type appCron struct {
	srv *AppServer
	rc  *cron.Cron
}

func newCron(srv *AppServer) *appCron {
	c := &appCron{srv: srv}

	// Cron jobs
	rc := cron.New()

	// Peers data
	_ = rc.AddFunc("@every 1s", c.Peers)

	rc.Start()

	c.rc = rc

	return c
}

func (c *appCron) Peers() {
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
