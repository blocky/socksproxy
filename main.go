package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/armon/go-socks5"
)

const (
	allowed            = true
	denied             = false
	letsEncryptProd    = "acme-v02.api.letsencrypt.org"
	letsEncryptStaging = "acme-staging-v02.api.letsencrypt.org"
)

type allowLetsEncrypt struct{}

func (m allowLetsEncrypt) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	for _, fqdn := range []string{letsEncryptProd, letsEncryptStaging} {
		if req.DestAddr.FQDN == fqdn {
			return ctx, allowed
		}
	}
	return ctx, denied
}

type logAll struct {
	rules socks5.RuleSet
}

func (m logAll) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	ctx, allowed := m.rules.Allow(ctx, req)

	prefix := "Denying"
	if allowed {
		prefix = "Allowing"
	}
	from := req.RemoteAddr
	to := req.DestAddr
	log.Printf("%s connection request from %s:%d (%s) to %s:%d (%s).",
		prefix,
		from.IP, from.Port, from.FQDN,
		to.IP, to.Port, to.FQDN,
	)

	return ctx, allowed
}

type config struct {
	addr         string
	fqdnAllowAll bool
	verbose      bool
}

func parseFlags() config {
	c := config{}

	flag.StringVar(&c.addr, "addr", ":1080", "Address to listen on.")
	flag.BoolVar(&c.fqdnAllowAll, "fqdn-allow-all", false, "Allow all FQDNs")
	flag.BoolVar(&c.verbose, "verbose", false, "log verbosely")
	flag.Parse()

	return c
}

func newServer(cfg config) (*socks5.Server, error) {
	var rules socks5.RuleSet = allowLetsEncrypt{}
	if cfg.fqdnAllowAll {
		rules = socks5.PermitAll()
	}

	if cfg.verbose {
		rules = logAll{rules}
	}

	conf := &socks5.Config{Rules: rules}
	server, err := socks5.New(conf)
	if err != nil {
		return nil, fmt.Errorf("creating proxy server: %w", err)
	}

	return server, nil
}

func main() {
	cfg := parseFlags()

	server, err := newServer(cfg)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting SOCKSv5 server on %s.", cfg.addr)
	if err := server.ListenAndServe("tcp", cfg.addr); err != nil {
		panic(err)
	}
}
