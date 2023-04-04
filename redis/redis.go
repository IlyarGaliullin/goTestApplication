package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"testApplication/models"
	"testApplication/utils"
	"time"
)

var tokenLifespan = time.Hour * 6

type Connection struct {
	client *redis.Client
}

var ErrUnauthorized = errors.New("unauthorized")

func NewConn() (*Connection, error) {

	host := utils.Conf.GetString("redis.host")
	port := utils.Conf.GetString("redis.port")
	password := utils.Conf.GetString("redis.password")
	db := utils.Conf.GetInt("redis.db")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	return &Connection{client: client}, nil
}

func (redisConn *Connection) AddToken(ctx context.Context, user models.User, token string) error {

	userKey := fmt.Sprintf("user:%d", user.Id)

	err := redisConn.client.Set(ctx, token, userKey, tokenLifespan).Err()
	if err != nil {
		return err
	}
	return nil
}

func (redisConn *Connection) CheckToken(ctx context.Context, token string) (userId int, err error) {

	result, err := redisConn.client.Get(ctx, token).Result()

	if err != nil || result == "" {
		return userId, ErrUnauthorized
	}
	userId, err = strconv.Atoi(strings.Split(result, ":")[1])

	if err != nil {
		return userId, ErrUnauthorized
	}
	return
}

func (redisConn *Connection) RemoveToken(ctx context.Context, token string) error {

	val, err := redisConn.client.Del(ctx, token).Result()

	if err != nil {
		return err
	}
	if val == 0 {
		return errors.New("no token exists")
	}
	return nil
}
