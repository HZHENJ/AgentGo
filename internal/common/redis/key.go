package redis

import (
	"fmt"
	"agentgo/pkg/conf"
)
// github.com/go-redis/redis/v8
// GenerateCaptchaKey generates key for captcha based on email
func GenerateCaptchaKey(email string) string {
	return fmt.Sprintf(conf.DefaultRedisKeyConfig.CaptchaPrefix, email)
}