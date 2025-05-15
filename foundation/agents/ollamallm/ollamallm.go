// Package ollamallm provides a client to act as the user in the
// client system.
package ollamallm

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Agent struct {
	llm    *ollama.LLM
	prompt string
}

func New(profilePath string) (*Agent, error) {
	f, err := os.Open(profilePath)
	if err != nil {
		return nil, fmt.Errorf("profile: %w", err)
	}
	defer f.Close()

	prompt, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading profile: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("llama3.2:latest"))
	if err != nil {
		return nil, fmt.Errorf("ollama: %w", err)
	}

	a := Agent{
		llm:    llm,
		prompt: string(prompt),
	}

	return &a, nil
}

func (a *Agent) Chat(ctx context.Context, input string, history []string) (string, error) {
	var b strings.Builder
	for _, h := range history {
		b.WriteString(fmt.Sprintf("%s\n\n", h))
	}

	prompt := fmt.Sprintf(a.prompt, b.String(), input)

	result, err := a.llm.Call(ctx, prompt, llms.WithMaxTokens(500), llms.WithTemperature(1.0))
	if err != nil {
		return "", fmt.Errorf("call: %w", err)
	}

	return result, nil
}
