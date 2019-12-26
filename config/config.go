package config

import (
	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	MainServer mainServer `toml:"main_server"`
	Log        logSection `toml:"log"`
	Db         DbConf     `toml:"db"`
}

func (c *Config) SetLogger(logger *log.Logger) {
	if len(c.Log.Filename) > 0 {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   c.Log.Filename,
			MaxSize:    c.Log.MaxSize, // megabytes
			MaxBackups: c.Log.MaxBackups,
			MaxAge:     c.Log.MaxAge,   // days
			Compress:   c.Log.Compress, // disabled by default
		})
	}

	logger.SetLevel(c.GetLogLever())
}

func (c *Config) Load(confPath string) error {
	_, err := toml.DecodeFile(confPath, c)
	return err
}

func (c *Config) GetLogLever() log.Lvl {
	return c.Log.GetLevel()
}

type mainServer struct {
	CountWorkers int    `toml:"count_workers"`
	ServiceName  string `toml:"service_name"`
	Host         string `toml:"host"`
	Port         string `toml:"port"`
}

type logSection struct {
	Level      string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	Filename   string `toml:"filename,omitempty"`
}

type DbConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	User string `toml:"user"`
	Name string `toml:"name"`
	Pass string `toml:"pass"`
}

func (l *logSection) GetLevel() log.Lvl {
	var lvl log.Lvl

	switch l.Level {
	case "DEBUG":
		lvl = log.DEBUG
	case "INFO":
		lvl = log.INFO
	case "WARN":
		lvl = log.WARN
	case "ERROR":
		lvl = log.ERROR
	default:
		lvl = log.INFO
	}
	return lvl
}
