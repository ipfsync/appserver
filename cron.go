package appserver

type appCron struct {
	srv *AppServer
}

func newCron(srv *AppServer) *appCron {
	c := &appCron{srv: srv}
	return c
}

func (c *appCron) Peers() {
	peers, _ := c.srv.api.Peers()
	var str string
	for _, p := range peers {
		str += p.Address().String() + "\n"
	}
}
