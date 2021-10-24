package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/mongo"
	sv "github.com/core-go/service"
	"github.com/core-go/sql"
)

type Root struct {
	Server     sv.ServerConf     `mapstructure:"server"`
	Provider   string            `mapstructure:"provider"`
	Sql        sql.Config        `mapstructure:"sql"`
	Mongo      mongo.MongoConfig `mapstructure:"mongo"`
	Log        log.Config        `mapstructure:"log"`
	MiddleWare mid.LogConfig     `mapstructure:"middleware"`
	Status     *sv.StatusConfig  `mapstructure:"status"`
	Action     *sv.ActionConfig  `mapstructure:"action"`
}
