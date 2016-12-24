package cache

import (
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = newPool()
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err)
			}
			return c, err
		},
	}
}

func set(args... string) error {
	var (
		err error
	)
	c := pool.Get()
	defer c.Close()
	_, err = c.Do("SET", args)
	return err
}

func get(key string) (string, error) {
	var (
		err   error
		value string
	)
	c := pool.Get()
	defer c.Close()
	value, err = redis.String(c.Do("GET", key))
	return value, err
}

func del(key string)error{
    var(err error)
    c:=pool.Get()
    defer c.Close()
    _,err=c.Do("DEL",key)
    return err
}