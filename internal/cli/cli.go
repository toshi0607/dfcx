package cli

import (
	"os"

	"github.com/toshi0607/dfcx/internal/command"
	"github.com/toshi0607/dfcx/internal/logger"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func Run(argv []string) error {
	var logLevel string

	app := &cli.App{
		Name:  "dfcx",
		Usage: "operate dialogflow cx",
		Commands: []*cli.Command{
			command.Agent(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "Log level [debug|info|warn|error]",
				EnvVars:     []string{"DF_LOG_LEVEL"},
				Destination: &logLevel,
				Value:       "info",
			},
		},
	}

	logger.Logger = slog.New(slog.HandlerOptions{Level: level(logLevel)}.NewJSONHandler(os.Stdout))

	if err := app.Run(argv); err != nil {
		logger.Logger.Error("exit with error", err)
		return err
	}

	return nil
}

func level(logLevel string) slog.Level {
	switch logLevel {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
