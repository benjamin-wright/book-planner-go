package tokens

import (
	"context"
	"fmt"
	"time"

	r "github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/book-planner-go/src/pkg/redis"
)

type Tokens struct {
	redis *r.Client
}

func New() (*Tokens, error) {
	config, err := redis.ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection config: %+v", err)
	}

	r, err := redis.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %+v", err)
	}

	tokens := Tokens{
		redis: r,
	}

	return &tokens, nil
}

func (t *Tokens) DeleteToken(id string, kind string) error {
	key := fmt.Sprintf("%s.%s", id, kind)
	err := t.redis.Del(context.Background(), key).Err()

	if err != nil {
		return fmt.Errorf("error deleting token %s: %+v", key, err)
	}

	return nil
}

func (t *Tokens) GetToken(id string, kind string) (string, error) {
	key := fmt.Sprintf("%s.%s", id, kind)
	value, err := t.redis.Get(context.Background(), key).Result()
	if err == r.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch token: %+v", err)
	}

	return value, nil
}

func (t *Tokens) AddToken(token string, id string, kind string, expiration time.Duration) (string, error) {
	key := fmt.Sprintf("%s.%s", id, kind)
	err := t.redis.Set(context.Background(), key, token, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save token: %+v", err)
	}

	return token, nil
}

func (t *Tokens) AddPasswordHash(id string, password string) error {
	key := fmt.Sprintf("%s.password", id)

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %+v", err)
	}
	hash := string(bytes)

	err = t.redis.Set(context.Background(), key, hash, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to send hashed password to redis: %+v", err)
	}

	return nil
}

func (t *Tokens) CheckPassword(id string, password string) (bool, error) {
	key := fmt.Sprintf("%s.password", id)

	hash, err := t.redis.Get(context.Background(), key).Result()
	if err == r.Nil {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to fetch password: %+v", err)
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}
