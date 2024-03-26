package main

import (
	"fmt"
	"net"
	"os/exec"
	"testing"
)

func TestStartServer(t *testing.T) {
	other := "https://example.com"
	letsEncrypt := fmt.Sprintf("https://%s", letsEncryptStaging)

	for _, tt := range []struct {
		allowAll bool
		dst      string
		success  bool
	}{
		{allowAll: false, dst: letsEncrypt, success: true},
		{allowAll: false, dst: other, success: false},
		{allowAll: true, dst: letsEncrypt, success: true},
		{allowAll: true, dst: other, success: true},
	} {
		cfg := config{
			addr:         ":1081",
			verbose:      false,
			fqdnAllowAll: tt.allowAll,
		}

		server, err := newServer(cfg)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		listener, err := net.Listen("tcp", cfg.addr)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		wait := make(chan struct{})
		go func() {
			// we don't care about the error here
			_ = server.Serve(listener)
			close(wait)
		}()

		_, err = exec.Command(
			"curl",
			"--socks5-hostname",
			"localhost"+cfg.addr,
			tt.dst,
		).Output()

		if tt.success && err != nil {
			t.Errorf("expected no error, got %v", err)
		} else if !tt.success && err == nil {
			t.Errorf("expected error, got no error")
		}

		listener.Close()
		<-wait
	}
}
