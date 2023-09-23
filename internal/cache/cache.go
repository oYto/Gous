package cache

import (
	"Gous/config"
	"Gous/internal/model"
	"Gous/internal/utils"
	"Gous/pkg/constant"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetUserInfoFromCache 查询用户缓存信息
func GetUserInfoFromCache(username string) (*model.User, error) {
	// 根据 userinfo_ + username 设置 key 值
	redisKey := constant.UserInfoPrefix + username
	// 查询信息
	val, err := utils.GetRedisCLi().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), user)
	return user, err
}

// SetUserCacheInfo 缓存用户信息
func SetUserCacheInfo(user *model.User) error {
	rediskey := constant.UserInfoPrefix + user.Name
	val, err := json.Marshal(user)
	// 解析失败
	if err != nil {
		return err
	}
	// 缓存过期时间
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.UserExpired)
	// 键 值 过期时间
	_, err = utils.GetRedisCLi().Set(context.Background(), rediskey, val, expired*time.Second).Result()
	return err
}

// SetSessionInfo 缓存会话 session
func SetSessionInfo(user *model.User, session string) error {
	//生成一个redis key，该key的格式是 "session:{session}"，其中的"{session}"是一个变量，代表了用户的session ID
	redisKey := constant.SessionKeyPrefix + session
	//将用户对象序列化为JSON格式的字符串。json.Marshal函数接受一个任意类型的值作为参数，并返回一个[]byte类型的字节数组和一个错误对象。
	val, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	//根据全局配置文件中的 Cache.SessionExpired 字段获取缓存的过期时间，然后将其转换为 time.Duration 类型，乘以 time.Second，最后得到一个 time.Duration 类型的过期时间。
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.SessionExpired)
	//使用了Go Redis客户端库（github.com/go-redis/redis）提供的Set方法，将用户信息序列化为JSON字符串后存储在Redis中，并设置了过期时间为配置文件中设定的Session过期时间。
	//该方法返回一个Redis的Reply对象，其中包含了当前键值的状态信息和执行结果。
	//在这里使用了一个匿名变量来忽略掉状态信息，只关心执行结果的error类型，以便上层业务可以判断是否存储成功。
	_, err = utils.GetRedisCLi().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

// GetSessionInfo 查询缓存中是否存在该 session，用于判断是否处于登录状态
func GetSessionInfo(session string) (*model.User, error) {
	redisKey := constant.SessionKeyPrefix + session
	//调用 Redis 客户端对象的 Get 方法，从 Redis 中获取 key 对应的值。
	val, err := utils.GetRedisCLi().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func DelSessionInfo(session string) error {
	log.Infof("session is:%s", session)
	redisKey := constant.SessionKeyPrefix + session
	_, err := utils.GetRedisCLi().Del(context.Background(), redisKey).Result()
	return err
}

// DelUserCacheInfo 删除缓存中的用户信息，注销时用
func DelUserCacheInfo(user *model.User) error {
	redisKey := constant.UserInfoPrefix + user.Name
	log.Infof("rediskey=============%v", redisKey)
	_, err := utils.GetRedisCLi().Del(context.Background(), redisKey).Result()
	return err
}

// UpdateCachedUserInfo 更新用户在redis中的信息
func UpdateCachedUserInfo(user *model.User) error {
	err := SetUserCacheInfo(user) //直接添加同名key，原本的数据会被覆盖
	if err != nil {
		redisKey := constant.UserInfoPrefix + user.Name
		//如果存储失败，则说明 Redis 存储出现问题，这时候需要将原先的缓存删除以避免缓存数据和数据库中数据不一致。
		utils.GetRedisCLi().Del(context.Background(), redisKey).Result()
	}
	return err
}
