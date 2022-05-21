package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Config 连接配置
type Config struct {
	Host     string
	Port     int64
	Password string
	DB       int
	PoolSize int
}

var clich = make(chan map[string]*redis.Client)

// InitWithName 初始化redis连接
func InitWithName(name string) error {
	var conf Config
	sub := viper.Sub("redis." + name)
	if sub == nil {
		return fmt.Errorf("invalid redis config under %s", name)
	}
	if err := sub.Unmarshal(&conf); err != nil {
		return err
	}

	if conf.Port == 0 {
		conf.Port = 6379
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	}

	if conf.PoolSize > 0 {
		opts.PoolSize = conf.PoolSize
	}

	cli := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return err
	}

	// set cli
	clich <- map[string]*redis.Client{name: cli}

	return nil
}

// GetClient 获取redis连接客户端
func GetClient(name string) *redis.Client {
	climap := <-clich
	if cli, ok := climap[name]; ok {
		return cli
	}

	return nil
}

func cliPool() {
	var clis = make(map[string]*redis.Client)
	for {
		select {
		case climap := <-clich:
			for name, cli := range climap {
				clis[name] = cli
			}
		case clich <- clis:
		}
	}
}

func init() {
	go cliPool()
}
