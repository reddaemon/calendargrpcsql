package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/reddaemon/calendargrpcsql/config"
)

func GetDb(c *config.Config) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Db.Host, c.Db.Port, c.Db.User, c.Db.Pass, c.Db.Name)

	return sqlx.Connect("postgres", psqlInfo)
}
