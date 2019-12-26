package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendargrpcsql/config"
)

type App struct {
	Config *config.Config
	Logger *log.Logger
	Db     *sqlx.DB
}
