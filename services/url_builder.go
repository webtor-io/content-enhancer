package services

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/urfave/cli"
)

const (
	proxyHost = "proxy-host"
	proxyPort = "proxy-port"
	apiKey    = "api-key"
)

func RegisterUrlBuilderFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.IntFlag{
			Name:   proxyPort,
			Usage:  "proxy port",
			EnvVar: "TORRENT_HTTP_PROXY_SERVICE_PORT",
		},
		cli.StringFlag{
			Name:   proxyHost,
			Usage:  "proxy host",
			EnvVar: "TORRENT_HTTP_PROXY_SERVICE_HOST",
		},
		cli.StringFlag{
			Name:   apiKey,
			Usage:  "api key",
			EnvVar: "API_KEY",
		},
	)
}

type UrlBuilder struct {
	proxyHost string
	proxyPort int
	apiKey    string
}

func NewUrlBuilder(c *cli.Context) *UrlBuilder {
	return &UrlBuilder{
		proxyHost: c.String(proxyHost),
		proxyPort: c.Int(proxyPort),
		apiKey:    c.String(apiKey),
	}
}

func (s *UrlBuilder) Build(sr *StatRecord) string {
	return fmt.Sprintf("http://%v:%v/%v/%v~tc~mhls~mhlsp?api-key=%v",
		s.proxyHost,
		s.proxyPort,
		sr.InfoHash,
		url.PathEscape(strings.TrimLeft(sr.OriginalPath, "/")),
		s.apiKey,
	)
}
