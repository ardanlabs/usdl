// Package ollamallm provides a client to act as the user in the
// client system.
package ollamallm

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

//go:embed prompts/prompt.txt
var prompt string

type Agent struct {
	llm *ollama.LLM
}

func New() (*Agent, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2:latest"))
	if err != nil {
		return nil, fmt.Errorf("ollama: %w", err)
	}

	a := Agent{
		llm: llm,
	}

	return &a, nil
}

func (a *Agent) Chat(ctx context.Context, input string, history []string) (string, error) {
	var b strings.Builder
	for _, h := range history {
		b.WriteString(fmt.Sprintf("%s\n\n", h))
	}

	prompt := fmt.Sprintf(prompt, b.String(), input)

	result, err := a.llm.Call(ctx, prompt, llms.WithMaxTokens(500), llms.WithTemperature(1.0))
	if err != nil {
		return "", fmt.Errorf("call: %w", err)
	}

	return result, nil
}
