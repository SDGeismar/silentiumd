package grpcservice

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/louisinger/silentiumd/internal/application"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

type Config struct {
	Port       uint32
	AppService application.SilentiumService
	NoTLS      bool
	HostName   string
}

func (c Config) Validate() error {
	lis, err := net.Listen("tcp", c.address())
	if err != nil {
		return fmt.Errorf("invalid port: %s", err)
	}
	defer lis.Close()

	return nil
}

func (c Config) insecure() bool {
	return c.NoTLS
}

func (c Config) address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c Config) gatewayAddress() string {
	return fmt.Sprintf("localhost:%d", c.Port)
}

func (c Config) tlsConfig() (*tls.Config, error) {
	hostPolicy := autocert.HostWhitelist(c.HostName)

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
	}

	return &tls.Config{
		Rand:           rand.Reader,
		Time:           time.Now,
		NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
		MinVersion:     tls.VersionTLS12,
		GetCertificate: m.GetCertificate,
	}, nil
}
