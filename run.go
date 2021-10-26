package main

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	s "github.com/webtor-io/content-enhancer/services"
)

func makeRunCMD() cli.Command {
	cmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "runs content enhancement",
		Action:  run,
	}
	configureRun(&cmd)
	return cmd
}

func configureRun(c *cli.Command) {
	c.Flags = s.RegisterClickHouseDBFlags([]cli.Flag{})
	c.Flags = s.RegisterUrlBuilderFlags(c.Flags)
}

func run(c *cli.Context) error {
	db := s.NewClickHouseDB(c)
	ch := s.NewClickHouse(c, db)
	ub := s.NewUrlBuilder(c)
	recs, err := ch.GetTopContent()
	cl := &http.Client{}
	if err != nil {
		return errors.Wrapf(err, "failed to get stat records")
	}
	for _, r := range recs {
		log.Infof("got record=%+v", r)
		u := ub.Build(&r)
		log.Infof("invoking url=\"%v\"", u)
		resp, err := cl.Get(u)
		if err != nil {
			log.WithError(err).Warnf("failed to invoke job url=\"%v\"", u)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Warnf("failed to invoke job url=\"%v\" status=\"%v\" code=%v", u, resp.Status, resp.StatusCode)
		} else {
			log.Infof("job invoked url=\"%v\" status=\"%v\" code=%v", u, resp.Status, resp.StatusCode)
		}
	}
	return nil
}
