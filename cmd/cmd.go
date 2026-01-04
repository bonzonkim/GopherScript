package cmd

import (
	"github.com/bonzonkim/gopher-script/config"
	"github.com/bonzonkim/gopher-script/internal/logger"
)

type Cmd struct {
	Config *config.Config
	Logger *logger.Logger
}

func NewCmd() *Cmd {
	c := &Cmd{
		Config: config.NewConfig(),
	}
	c.Logger = logger.NewLogger(c.Config.Env)

	return c
}
