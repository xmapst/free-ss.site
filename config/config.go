package config

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github/xmapst/free-ss.site/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Config  Configure
	version = "0.0.1"
	envfile = kingpin.Flag("envfile", "Load the environment variable file").Default(".envfile").String()
)

// Configure stores configuration.
type Configure struct {
	RunMode      string   `envconfig:"RUN_MODE" default:"release"`
	HTTPPort     int      `envconfig:"HTTP_PORT" default:"80"`
	HTTPAddr     string   `envconfig:"HTTP_ADDR" default:"0.0.0.0"`
	LogLevel     string   `envconfig:"LOG_LEVEL" default:"debug"`
	PrefixUrl    string   `envconfig:"PREFIX_URL" default:"https://sub.yfdou.com"`
	ReadTimeout  int64    `envconfig:"READ_TIMEOUT" default:"60"`
	WriteTimeout int64    `envconfig:"WRITE_TIMEOUT" default:"60"`
	IPWhiteList  []string `envconfig:"IP_WHITE_LIST" default:"*"`
}

func init() {
	logrus.SetFormatter(&utils.ConsoleFormatter{})
	logrus.SetReportCaller(true)
	logrus.Debug("Load variable")
	kingpin.Version(version)
	kingpin.Parse()
	// load environment variables from file.
	_ = godotenv.Load(*envfile)

	// load the configuration from the environment.
	err := envconfig.Process("", &Config)
	if err != nil {
		logrus.Fatalln(err)
	}
	level, err := logrus.ParseLevel(Config.LogLevel)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.SetLevel(level)
}
