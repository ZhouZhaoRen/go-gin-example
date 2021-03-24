package gredis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"go-gin-example/pkg/setting"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// 往Redis中添加元素数据
func Set(key string,data interface{},time int) error{
	conn:=RedisConn.Get()
	defer conn.Close()

	value,err:=json.Marshal(data)
	if err != nil {
		return err
	}

	_,err=conn.Do("SET",key,value)
	if err != nil {
		return err
	}
	//
	_,err=conn.Do("EXPIRE",key,time)
	if err!=nil{
		return err
	}
	//
	return nil
}

// 查找Redis中是否存在某个key
func Exists(key string) bool{
	conn:=RedisConn.Get()
	defer conn.Close()

	exists,err:=redis.Bool(conn.Do("EXISTS",key))
	if err != nil {
		return false
	}

	return exists
}

// 从Redis中获取数据,返回的是字节数组
func Get(key string)([]byte,error){
	conn:=RedisConn.Get()
	defer conn.Close()

	reply,err:=redis.Bytes(conn.Do("GET",key))
	if err != nil {
		return nil, err
	}
	return reply,nil
}

// 删除Redis中的数据
func Delete(key string)(bool,error){
	conn:=RedisConn.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DELETE",key))
}

// 删除相似的keys
func LikeDeletes(key string)error{
	conn:=RedisConn.Get()
	defer conn.Close()

	// 使用keys命令从Redis中查找出全部的key，然后逐个删除
	keys,err:=redis.Strings(conn.Do("KEYS","*"+key+"*"))
	if err != nil {
		return err
	}
	//
	for _,key:=range keys{
		_,err=Delete(key)
		if err != nil {
			return err
		}
	}
	//
	return nil
}