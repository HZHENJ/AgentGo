package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"agentgo/internal/model"
	"agentgo/pkg/conf" // 你的配置结构体
	"github.com/cloudwego/eino-ext/components/model/openai"
	einoModel "github.com/cloudwego/eino/components/model"
)

func init() {
	Register(TypeOpenAI, func(ctx context.Context, cfg *conf.LLMConfig) (Model, error) {
		return NewOpenAIAdapter(ctx, cfg)
	})
}

type OpenAIAdapter struct {
	llm einoModel.ToolCallingChatModel
}

// NewOpenAIAdapter
func NewOpenAIAdapter(ctx context.Context, cfg *conf.LLMConfig) (*OpenAIAdapter, error) {
	llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: cfg.BaseURL,
		Model:   cfg.ModelName,
		APIKey:  cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("create openai model failed: %v", err)
	}
	return &OpenAIAdapter{llm: llm}, nil
}

func (o *OpenAIAdapter) GenerateResponse(ctx context.Context, messages []*model.Message) (string, error) {
	einoMsgs := toEinoMessages(messages)
	
	resp, err := o.llm.Generate(ctx, einoMsgs)
	if err != nil {
		return "", fmt.Errorf("openai generate failed: %v", err)
	}
	return resp.Content, nil 
}

func (o *OpenAIAdapter) StreamResponse(ctx context.Context, messages []*model.Message, cb func(string)) (string, error) {
	einoMsgs := toEinoMessages(messages)
	stream, err := o.llm.Stream(ctx, einoMsgs)
	if err != nil {
		return "", fmt.Errorf("openai stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("openai stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content)
			cb(msg.Content)
		}
	}
	return fullResp.String(), nil
}

func (o *OpenAIAdapter) GetModelType() string { return TypeOpenAI }