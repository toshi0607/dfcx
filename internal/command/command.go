package command

import (
	"github.com/toshi0607/dfcx/internal/dialogflow"
	"github.com/urfave/cli/v2"
)

var (
	config  dialogflow.Config
	version string
)

func Agent() *cli.Command {
	return &cli.Command{
		Name:        "agent",
		Description: "dialogflow cx agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "location",
				EnvVars:     []string{"DF_LOCATION"},
				Usage:       "agent location",
				DefaultText: "asia-northeast1",
				Destination: &config.Location,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "base-project",
				EnvVars:     []string{"DF_BASE_PROJECT"},
				Usage:       "base project name",
				Destination: &config.BaseProjectID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "base-agent",
				EnvVars:     []string{"DF_BASE_AGENT"},
				Usage:       "base agent ID",
				Destination: &config.BaseAgentID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "base-env",
				EnvVars:     []string{"DF_BASE_ENV"},
				Usage:       "base environment ID",
				Destination: &config.BaseEnvID,
				Required:    true,
			},
		},
		Subcommands: []*cli.Command{
			deploy(),
		},
	}
}

func deploy() *cli.Command {
	return &cli.Command{
		Name:        "deploy",
		Description: "Operate database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "version",
				Aliases:     []string{"v"},
				Usage:       "version",
				Destination: &version,
				Required:    true,
			},
		},
		Subcommands: []*cli.Command{
			base(),
			stg(),
			prd(),
		},
	}
}

func base() *cli.Command {
	return &cli.Command{
		Name:        "base",
		Description: "deploy base agent",
		Flags:       nil,
		Action: func(context *cli.Context) error {
			config.TargetEnvID = config.BaseEnvID
			config.TargetAgentID = config.BaseAgentID
			config.TargetProjectID = config.BaseProjectID
			if err := dialogflow.Deploy(context.Context, config, version); err != nil {
				return err
			}
			return nil
		},
	}
}

func stg() *cli.Command {
	return &cli.Command{
		Name:        "stg",
		Description: "deploy stg agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "stg-project",
				EnvVars:     []string{"DF_STG_PROJECT"},
				Usage:       "base project name",
				Destination: &config.TargetProjectID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "stg-agent",
				EnvVars:     []string{"DF_STG_AGENT"},
				Usage:       "base agent ID",
				Destination: &config.TargetAgentID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "stg-env",
				EnvVars:     []string{"DF_STG_ENV"},
				Usage:       "staging environment ID",
				Destination: &config.TargetEnvID,
				Required:    true,
			},
		},
		Action: func(context *cli.Context) error {
			if err := dialogflow.Deploy(context.Context, config, version); err != nil {
				return err
			}
			return nil
		},
	}
}

func prd() *cli.Command {
	return &cli.Command{
		Name:        "prd",
		Description: "deploy prd agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "prd-project",
				EnvVars:     []string{"DF_PRD_PROJECT"},
				Usage:       "production project name",
				Destination: &config.TargetProjectID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "prd-agent",
				EnvVars:     []string{"DF_PRD_AGENT"},
				Usage:       "production agent ID",
				Destination: &config.TargetAgentID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "prd-env",
				EnvVars:     []string{"DF_PRD_ENV"},
				Usage:       "production environment ID",
				Destination: &config.TargetEnvID,
				Required:    true,
			},
		},
		Action: func(context *cli.Context) error {
			if err := dialogflow.Deploy(context.Context, config, version); err != nil {
				return err
			}
			return nil
		},
	}
}
