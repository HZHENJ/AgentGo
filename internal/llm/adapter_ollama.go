package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"agentgo/internal/model"
	"agentgo/pkg/conf"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	einoModel "github.com/cloudwego/eino/components/model"
)


func init() {
	Register(TypeOllama, func(ctx context.Context, cfg *conf.LLMConfig) (Model, error) {
		return NewOllamaAdapter(ctx, cfg)
	})
}

// OllamaAdapter 
type OllamaAdapter struct {
	llm einoModel.ToolCallingChatModel
}

// NewOllamaAdapter 构造函数
func NewOllamaAdapter(ctx context.Context, cfg *conf.LLMConfig) (*OllamaAdapter, error) {
	llm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.BaseURL,
		Model:   cfg.ModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create ollama model failed: %v", err)
	}
	return &OllamaAdapter{llm: llm}, nil
}

// GenerateResponse 
func (o *OllamaAdapter) GenerateResponse(ctx context.Context, messages []*model.Message) (string, error) {
	einoMsgs := toEinoMessages(messages)
	
	resp, err := o.llm.Generate(ctx, einoMsgs)
	if err != nil {
		return "", fmt.Errorf("ollama generate failed: %v", err)
	}
	return resp.Content, nil
}

// StreamResponse 
func (o *OllamaAdapter) StreamResponse(ctx context.Context, messages []*model.Message, cb func(string)) (string, error) {
	einoMsgs := toEinoMessages(messages)
	
	stream, err := o.llm.Stream(ctx, einoMsgs)
	if err != nil {
		return "", fmt.Errorf("ollama stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			
			return "", fmt.Errorf("ollama stream recv failed: %v", err)
		}
		
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content)
			cb(msg.Content)
		}
	}

	return fullResp.String(), nil 
}

func (o *OllamaAdapter) GetModelType() string { return TypeOllama }