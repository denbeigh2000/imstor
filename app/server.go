package imstor

import "log"

// Imstor is the main type used by application servers. It handles incoming
// user requests, makes calls to data stores, and queues asynchronous actions
// like thumbnailing.
type Imstor struct {
	UserAPI
	Servers []Server
}

func (i *Imstor) Serve() {
	for _, server := range i.Servers {
		go func(s Server) {
			for {
				log.Printf("Error serving: %v", s.Serve(i))
			}
		}(server)
	}

	// TODO: Have a more graceful way of serving until closed
	// also grumble github.com/golang/go/issues/4674
	select {}
}
