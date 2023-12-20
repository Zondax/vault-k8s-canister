package icp

import (
	"fmt"
	"net/url"
	"os"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"go.uber.org/zap"
)

var client *Agent

func GetBackendClient() (*Agent, error) {
	if client == nil {
		canisterID := os.Getenv("CANISTER_ID")
		icpNodeUrl := os.Getenv("ICP_NODE_URL")

		if canisterID == "" {
			return nil, fmt.Errorf("CANISTER_ID env var must be defined")
		}

		if icpNodeUrl == "" {
			return nil, fmt.Errorf("ICP_NODE_URL env var must be defined")
		}

		ledgerID, err := principal.Decode(canisterID)
		if err != nil {
			return nil, err
		}

		localReplica, err := url.Parse(icpNodeUrl)
		if err != nil {
			return nil, err
		}

		ccfg := &agent.ClientConfig{Host: localReplica}

		randomIdentity, err := identity.NewRandomSecp256k1Identity()
		if err != nil {
			zap.S().Errorf("[ICP] Failed to create client: %+v", err)
		}

		a, err := NewAgent(ledgerID, agent.Config{ClientConfig: ccfg, FetchRootKey: true, Identity: randomIdentity})
		if err != nil {
			return nil, err
		}

		client = a
	}

	return client, nil
}
