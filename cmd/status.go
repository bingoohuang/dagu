package main

import (
	"log"

	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/controller"
	"github.com/yohamta/dagu/internal/models"
	"github.com/yohamta/dagu/internal/utils"

	"github.com/urfave/cli/v2"
)

func newStatusCommand() *cli.Command {
	cl := &config.Loader{
		HomeDir: utils.MustGetUserHomeDir(),
	}
	return &cli.Command{
		Name:  "status",
		Usage: "dagu status <config>",
		Action: func(c *cli.Context) error {
			configFilePath := c.Args().Get(0)
			cfg, err := cl.Load(configFilePath, "")
			if err != nil {
				return err
			}
			return queryStatus(cfg)
		},
	}
}

func queryStatus(cfg *config.Config) error {
	status, err := controller.New(cfg).GetStatus()
	if err != nil {
		return err
	}
	res := &models.StatusResponse{
		Status: status,
	}
	log.Printf("Pid=%d Status=%s", res.Status.Pid, res.Status.Status)
	return nil
}
