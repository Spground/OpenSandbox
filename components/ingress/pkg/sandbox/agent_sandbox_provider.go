package sandbox

import "context"

type AgentSandboxProvider struct {
}

func NewAgentSandboxProvider() *AgentSandboxProvider {
	return &AgentSandboxProvider{}
}

func (a *AgentSandboxProvider) GetEndpoint(name string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AgentSandboxProvider) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

var _ Provider = (*AgentSandboxProvider)(nil)
