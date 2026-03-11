package llm

import (
	"context"
	"agentgo/internal/model"
	"github.com/cloudwego/eino/schema"
)

type Model interface {
	GenerateResponse(ctx context.Context, messages []*model.Message) (string, error)
	StreamResponse(ctx context.Context, messages []*model.Message, cb func(string)) (string, error)
	GetModelType() string
}

// toEinoMessages converts a slice of model.Message to a slice of schema.Message for use with the Eino library.
func toEinoMessages(msgs []*model.Message) []*schema.Message {
	if msgs == nil {
		return nil
	}
	res := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant // default to assistant role
		if m.IsUser {
			role = schema.User   // set to user role if IsUser is true
		}
		res = append(res, &schema.Message{
			Role:    role,
			Content: m.Content,
		})
	}
	return res
}
