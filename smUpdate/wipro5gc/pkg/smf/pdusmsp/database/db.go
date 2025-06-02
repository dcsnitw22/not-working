package db

import (
	"k8s.io/klog"
	redisClient "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/database/redis"
)

// TODO crud functions for db-parameter client db name, database name
type DBManager interface {
	Start()
}

func (db *DBInfo) Start() {
	klog.Info("DataBase initialize")

}

// add all database client in struct
type DBInfo struct {
	Redis redisClient.RedisInfo
}

func NewDBManager() *DBInfo {
	rClient := redisClient.NewRedisDbManager()

	// //Redis Client initialization
	// client, err := redisClient.NewRedisClient(ctx, 0)
	// if err != nil {
	// 	klog.Errorf("unable to connect to database:ERROR:%s", err.Error())
	// }
	// klog.Info("Connected to redis Database")
	return &DBInfo{
		Redis: *rClient,
	}
}

// func (db *DBInfo) CreateDbContext(key string, data string) error {

// 	// Use the Redis client to create a new entry in the database
// 	_, err := redisClient.Create(key, data, redisClient.SessionDb, db.RedisClient)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *DBInfo) ReadDbContext(key string) (string, error) {
// 	// Use the Redis client to retrieve data
// 	data, err := redisClient.Read(key, redisClient.SessionDb, db.RedisClient)
// 	if err != nil {
// 		return "", err
// 	}
// 	return data, nil
// }

// func (db *DBInfo) UpdateDbContext(key string, data string) error {

// 	// Use the Redis client to update the entry in the database
// 	_, err := redisClient.Update(key, data, redisClient.SessionDb, db.RedisClient)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *DBInfo) DeleteDbContext(key string) error {
// 	// Use the Redis client to delete the entry from the database
// 	_, err := redisClient.Delete(key, redisClient.SessionDb, db.RedisClient)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
