package sm

import (
	"context"

	"github.com/go-redis/redis/v8"
	"k8s.io/klog"
)

func NewRedisClient(ctx context.Context, db int) (*redis.Client, error) {
	klog.Info("establishing connection with redis server")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis.redis.svc.cluster.local:6379",
		Password: "redisWipro",
		DB:       db,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		klog.Error("Unable to establish connection")
		return nil, err
	} else {
		klog.Info("Established connection with redis server")
	}
	//verification
	val, err := client.Get(ctx, "66218A9E").Result()
	if err == redis.Nil {
		klog.Error("Key does not exist")
	} else if err != nil {
		klog.Error("Error retrieving key: %v", err)
	} else {
		klog.Info("Key value: %s", val)
	}
	return client, nil
}
