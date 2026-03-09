package cache

import (
	"time"
	"context"
	"github.com/go-redis/redis/v8"
	"strings"

	rediskeys "agentgo/internal/common/redis"
)

type UserCacheDao interface {
	SetCaptchaForEmail(ctx context.Context, email, captcha string) error
	CheckCaptchaForEmail(ctx context.Context, email, captcha string) (bool, error)
}

type userCacheDao struct {
	rdb *redis.Client
}

func NewUserCacheDao(rdb *redis.Client) UserCacheDao {
	return &userCacheDao{
		rdb: rdb,
	}
}

// SetCaptchaForEmail stores the captcha code for a given email in Redis with an expiration time.
func (dao *userCacheDao) SetCaptchaForEmail(ctx context.Context, email, captcha string) error {
	key := rediskeys.GenerateCaptchaKey(email)
	expire := 2 * time.Minute
	return dao.rdb.Set(ctx, key, captcha, expire).Err()
}

// CheckCaptchaForEmail checks if the provided captcha matches the one stored in Redis for the given email.
func (dao *userCacheDao) CheckCaptchaForEmail(ctx context.Context, email, captcha string) (bool, error) {
	key := rediskeys.GenerateCaptchaKey(email)

	value, err := dao.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Key does not exist
		}
		return false, err // Other Redis error
	}

	if strings.EqualFold(value, captcha) {
		// if the captcha matches, delete the key to prevent reuse
		if err := dao.rdb.Del(ctx, key).Err(); err != nil {
		}
		return true, nil
	}
	return false, nil
}