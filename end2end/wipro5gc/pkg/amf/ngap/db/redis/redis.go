package redis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-redis/redis/v8"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
)

type AmfUeRan struct {
	Supi        string
	AmfUeNgapId uint64
	RanUeNgapId uint64
}

/*type AmfUeNgapIdToSupi struct {
	Supi string
}

type SupiToAmfUeNgapId struct {
	AmfUeNgapId string
}*/

type RedisClient struct {
	c *redis.Client
}

func NewRedisClient(config config.RedisConfig) *RedisClient {
	opts := redis.Options{
		Addr:     config.IP + ":" + config.Port,
		Password: config.Pass,
		DB:       config.Db,
	}
	client := redis.NewClient(&opts)
	return &RedisClient{c: client}
}

func (r *RedisClient) Start() (string, error) {
	ctx := context.Background()
	res, err := r.c.Ping(ctx).Result()
	return res, err
}

func (r *RedisClient) Read(key string) (*AmfUeRan, error) {
	ctx := context.Background()
	klog.Info("key : ", key)
	res, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	klog.Info("value read from db : ", res)
	amfUeRan := &AmfUeRan{}
	if err = json.Unmarshal([]byte(res), amfUeRan); err != nil {
		return nil, err
	}
	return amfUeRan, nil
}

func (r *RedisClient) Write(key string, val string) (string, error) {
	ctx := context.Background()
	res, err := r.c.Set(ctx, key, val, 0).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RedisClient) Update(key string, updates map[string]interface{}) (string, error) {
	var amfUeRan AmfUeRan
	ctx := context.Background()
	res, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return "", errors.New("failed to get value from redis : " + err.Error())
	}
	err = json.Unmarshal([]byte(res), &amfUeRan)
	if err != nil {
		return "", errors.New("failed to unmarshal json value : " + err.Error())
	}

	for field, newValue := range updates {
		fieldValue := reflect.ValueOf(&amfUeRan).Elem().FieldByName(field)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			newValueReflect := reflect.ValueOf(newValue)
			if newValueReflect.Type().AssignableTo(fieldValue.Type()) {
				fieldValue.Set(newValueReflect)
			}
		}
	}

	updatedJson, err := json.Marshal(amfUeRan)
	if err != nil {
		return "", errors.New("failed to marshal to struct : " + err.Error())
	}
	res, err = r.c.Set(ctx, key, updatedJson, 0).Result()
	if err != nil {
		return "", errors.New("failed to set value in redis : " + err.Error())
	}
	return res, nil
}
