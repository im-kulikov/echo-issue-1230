package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("toml")
	v.SetConfigFile("config.cfg")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Defaults:
	v.SetDefault("logs.log_level", "info")
	v.SetDefault("logs.log_path", "app.log")
	v.SetDefault("server.listen_ip", "0.0.0.0")
	v.SetDefault("server.listen_port", 8888)

	// Echo instance
	e := echo.New()

	// Setup application
	fd, err := os.OpenFile(
		v.GetString("logs.log_path"),
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		panic(err)
	}

	e.Logger.SetOutput(fd)

	switch lvl := v.GetString("logs.log_level"); lvl {
	case "DEBUG", "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "INFO", "info":
		e.Logger.SetLevel(log.INFO)
	case "WARN", "warn":
		e.Logger.SetLevel(log.WARN)
	case "ERROR", "error":
		e.Logger.SetLevel(log.ERROR)
	case "OFF", "off":
		e.Logger.SetLevel(log.OFF)
	default:
		panic(fmt.Sprintf("unknown log level: %s",
			lvl))
	}

	host := v.GetString("server.listen_ip")
	port := v.GetInt64("server.listen_port")

	// Middleware
	logConfig := middleware.DefaultLoggerConfig
	logConfig.Output = fd

	e.Use(middleware.LoggerWithConfig(logConfig))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(host + ":" + strconv.FormatInt(port, 10)))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
