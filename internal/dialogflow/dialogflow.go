package dialogflow

import (
	"context"
	"fmt"

	cx "cloud.google.com/go/dialogflow/cx/apiv3"
	"cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"github.com/morikuni/failure"
	"github.com/toshi0607/dfcx/internal/logger"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Config struct {
	Location        string
	BaseAgentID     string
	BaseProjectID   string
	BaseEnvID       string
	TargetAgentID   string
	TargetProjectID string
	TargetEnvID     string
}

const (
	asiaNorthEast1Endpoint = "asia-northeast1-dialogflow.googleapis.com"
	defaultFlowID          = "00000000-0000-0000-0000-000000000000"
)

func Deploy(ctx context.Context, cfg Config, version string) error {
	logger.Logger.Info(fmt.Sprintf("deploy version %s to %s started", version, cfg.TargetProjectID))
	if cfg.BaseAgentID != cfg.TargetAgentID {
		ac, err := cx.NewAgentsClient(ctx, option.WithEndpoint(fmt.Sprintf("%s:443", asiaNorthEast1Endpoint)))
		if err != nil {
			return failure.Wrap(err)
		}
		defer func() {
			if err := ac.Close(); err != nil {
				logger.Logger.Error("failed to close agent client", err)
			}
		}()

		exportReq := &cxpb.ExportAgentRequest{
			Name:        baseAgent(cfg),
			Environment: baseEnvironment(cfg),
		}
		exportOp, err := ac.ExportAgent(ctx, exportReq)
		if err != nil {
			return failure.Wrap(fmt.Errorf("failed to export agent, error: %v", err))
		}

		exportedAgent, err := exportOp.Wait(ctx)
		if err != nil {
			return failure.Wrap(fmt.Errorf("failed to wait an exportOp, error: %v", err))
		}

		content := exportedAgent.GetAgentContent()

		restoreReq := &cxpb.RestoreAgentRequest{
			Name:  targetAgent(cfg),
			Agent: &cxpb.RestoreAgentRequest_AgentContent{AgentContent: content},
		}
		restoreOp, err := ac.RestoreAgent(ctx, restoreReq)
		if err != nil {
			return failure.Wrap(fmt.Errorf("failed to restore the agent, error: %v", err))
		}

		err = restoreOp.Wait(ctx)
		if err != nil {
			return failure.Wrap(fmt.Errorf("failed to wait an exportOp, error: %v", err))
		}
	}

	vc, err := cx.NewVersionsClient(ctx, option.WithEndpoint(fmt.Sprintf("%s:443", asiaNorthEast1Endpoint)))
	if err != nil {
		return failure.Wrap(err)
	}
	defer func() {
		if err := vc.Close(); err != nil {
			logger.Logger.Error("failed to close version client", err)
		}
	}()

	vReq := &cxpb.CreateVersionRequest{
		Parent: targetFlow(cfg),
		Version: &cxpb.Version{
			DisplayName: version,
		},
	}
	vOp, err := vc.CreateVersion(ctx, vReq)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to create the version, error: %v", err))
	}

	_, err = vOp.Wait(ctx)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wait a create op, error: %v", err))
	}
	logger.Logger.Info("version created", "version", vReq.Version.DisplayName)

	m, err := vOp.Metadata()
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to get version metadata, error: %v", err))
	}

	ec, err := cx.NewEnvironmentsClient(ctx, option.WithEndpoint(fmt.Sprintf("%s:443", asiaNorthEast1Endpoint)))
	if err != nil {
		return failure.Wrap(err)
	}
	defer func() {
		if err := ec.Close(); err != nil {
			logger.Logger.Error("failed to close environments client", err)
		}
	}()

	var messageType *cxpb.Environment
	mask, err := fieldmaskpb.New(messageType, "version_configs")
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to create mask, error: %v", err))
	}

	eReq := &cxpb.UpdateEnvironmentRequest{
		Environment: &cxpb.Environment{
			Name: targetEnvironment(cfg),
			// Only default flow is assumed.
			// All reachable flows must be set.
			VersionConfigs: []*cxpb.Environment_VersionConfig{{Version: m.Version}},
		},
		UpdateMask: mask,
	}

	uOp, err := ec.UpdateEnvironment(ctx, eReq)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to update the environment, error: %v", err))
	}

	_, err = uOp.Wait(ctx)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wait a update op, error: %v", err))
	}
	logger.Logger.Info(fmt.Sprintf("updated the version to %s", vReq.Version.DisplayName))

	logger.Logger.Info("deploy finished")

	return nil
}

func targetAgent(cfg Config) string {
	return fmt.Sprintf("projects/%s/locations/%s/agents/%s", cfg.TargetProjectID, cfg.Location, cfg.TargetAgentID)
}

func targetEnvironment(cfg Config) string {
	return fmt.Sprintf("%s/environments/%s", targetAgent(cfg), cfg.TargetEnvID)
}

func targetFlow(cfg Config) string {
	return fmt.Sprintf("%s/flows/%s", targetAgent(cfg), defaultFlowID)
}

func baseAgent(cfg Config) string {
	return fmt.Sprintf("projects/%s/locations/%s/agents/%s", cfg.BaseProjectID, cfg.Location, cfg.BaseAgentID)
}

func baseEnvironment(cfg Config) string {
	return fmt.Sprintf("%s/environments/%s", baseAgent(cfg), cfg.TargetEnvID)
}
