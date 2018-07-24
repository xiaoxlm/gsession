package instance

import (
	"bytes"
	"github.com/gomodule/redigo/redis"
	"gsession"
	"strconv"
	"sync"
	"errors"

)

// 后面加入connect pool
type pool struct {
	MaxIdle   int
	MaxActive int
}

type RedisSession struct {
	sid      string
	provider *RedisProvider
}

func (rs *RedisSession) Set(key string, value interface{}) error {
	return rs.provider.hSet(rs.sid, key, value)
}

func (rs *RedisSession) SetMulti(kv map[string]interface{}) error {
	return rs.provider.hmset(rs.sid, kv)
}

func (rs *RedisSession) Get(key string) (interface{}, error) {
	return rs.provider.hGet(rs.sid, key)
}

func (rs *RedisSession) GetMulti(value interface{}) error {
	return rs.provider.hGetAll(rs.sid, value)
}

func (rs *RedisSession) Delete(key string) error {
	return rs.provider.hDel(rs.sid, key)
}

func (rs *RedisSession) SessionID() string {
	return rs.sid
}

func (rs *RedisSession) Clear() error {
	return rs.provider.SessionDestroy(rs.sid)
}

type RedisProvider struct {
	lock sync.RWMutex //用来锁
	conn redis.Conn
}

func (rp *RedisProvider) SessionInit(sid string, maxLife int64) (gsession.Session, error) {
	kv := map[string]interface{}{
		"field": "",
	}

	if err := rp.hmset(sid, kv); err != nil {
		panic(err)
	}
	if err := rp.expire(sid, maxLife); err != nil {
		panic(err)
	}

	return &RedisSession{sid, rp}, nil
}

func (rp *RedisProvider) SessionRead(sid string) (gsession.Session, error) {
	if !rp.exist(sid) {
		return &RedisSession{}, errors.New("sid read fail")
	}
	return &RedisSession{sid, rp}, nil
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
	err := rp.del(sid)
	return err
}

func (*RedisProvider) SessionGC() {
	return
}

func (rp *RedisProvider) SetConn(conf gsession.CommonConf) {
	var addr bytes.Buffer
	addr.WriteString(conf.Host)
	addr.WriteString(":")
	addr.WriteString(strconv.Itoa(int(conf.Port)))

	conn, err := redis.Dial("tcp", addr.String())
	if err != nil {
		panic(err)
	}

	if conf.Password != "" {
		_, err := conn.Do("AUTH", conf.Password)
		if err != nil {
			panic(err)
		}
	}
	if conf.DB != "" {
		conn.Do("SELECT", conf.DB)
	}
	rp.conn = conn

	return
}

func (rp *RedisProvider) hmset(key string, kv map[string]interface{}) error {
	var args []interface{}
	args = append(args, key)
	for k, v := range kv {
		args = append(args, k, v)
	}

	_, err := rp.conn.Do("HMSET", args...)
	return err
}

func (rp *RedisProvider) hSet(key string, field string, value interface{}) error {
	_, err := rp.conn.Do("HSET", key, field, value)
	return err
}

func (rp *RedisProvider) hGet(key string, field string) (string, error) {
	r, err := redis.String(rp.conn.Do("HGET", key, field))
	return r,err
}

func (rp *RedisProvider) hGetAll(key string, value interface{}) error {
	reply, err := rp.conn.Do("HGETALL", key)
	v, err := redis.Values(reply, err)
	err = redis.ScanStruct(v, value)
	return err
}

func (rp *RedisProvider) hDel(key string, fields ...interface{}) error {
	var args []interface{}
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}

	_, err := rp.conn.Do("HDEL", args...)
	return err
}

func (rp *RedisProvider) expire(key string, expireTime int64) error {
	_, err := rp.conn.Do("EXPIRE", key, expireTime)
	return err
}

func (rp *RedisProvider) del(key ...interface{}) error {
	_, err := rp.conn.Do("DEL", key...)
	return err
}

func (rp *RedisProvider) exist(key string) bool {
	r, _ := redis.Bool(rp.conn.Do("EXISTS", key))
	return r
}
