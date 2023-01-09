package dialogflow

import (
	"context"
	"fmt"

	cx "cloud.google.com/go/dialogflow/cx/apiv3"
	"cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"github.com/morikuni/failure"
	"github.com/toshi0607/dfcx/internal/logger"
	"google.golang.org/api/option"
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
			Name: version,
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

	// update env for target

	return nil
}

func targetAgent(cfg Config) string {
	return fmt.Sprintf("projects/%s/locations/%s/agents/%s", cfg.TargetProjectID, cfg.Location, cfg.TargetAgentID)
}

func targetEnvironment(cfg Config) string {
	return fmt.Sprintf("%s/locations/%s/agents/%s", targetAgent(cfg), cfg.TargetEnvID)
}

func targetFlow(cfg Config) string {
	return fmt.Sprintf("%s/locations/%s/agents/%s", targetAgent(cfg), defaultFlowID)
}

func baseAgent(cfg Config) string {
	return fmt.Sprintf("projects/%s/locations/%s/agents/%s", cfg.BaseProjectID, cfg.Location, cfg.BaseAgentID)
}

func baseEnvironment(cfg Config) string {
	return fmt.Sprintf("%s/locations/%s/agents/%s", baseAgent(cfg), cfg.TargetEnvID)
}