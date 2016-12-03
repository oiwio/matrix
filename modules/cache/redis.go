package cache

import (
	"github.com/garyburd/redigo/redis"
    // "fmt"
)

type Cache struct{
    RedisConn redis.Conn
}

func (c Cache) Set(key string,value string) error {
    var(
        err error
    )
    _,err=c.RedisConn.Do("SET",key,value)
    return err
}

func (c Cache) Get (key string) (string,error){
    var(
        err error
        value string
    )
    value,err=redis.String(c.RedisConn.Do("GET",key))
    return value,err
}
