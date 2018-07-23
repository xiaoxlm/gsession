package instance

import (
	"github.com/gomodule/redigo/redis"
	"gsession"
	"bytes"
	"strconv"
)

type pool struct {
	MaxIdle int
	MaxActive int
}

type RedisSession struct {

}

func (*RedisSession) Set(key string, value interface{}) error {}

func (*RedisSession) Get(key string) (interface{}, error) {}

func (*RedisSession) Delete(key string) error {}

func (*RedisSession) SessionID() string {}


type RedisProvider struct {
	conf gsession.CommonConf
}

func (*RedisProvider) SessionInit(sid string, maxLife int64) (Session, error) {

}

func (*RedisProvider) SessionRead(sid string) (Session, error) {}

func (*RedisProvider) SessionDestroy(sid string) error {}

func (*RedisProvider) SessionGC() {}
