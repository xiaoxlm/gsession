package gsession

import (
	"bytes"
	"github.com/gomodule/redigo/redis"
	"gsession/instance"
	"math/rand"
	"strconv"
	"time"
	"net/http"
	"net/url"
)

const (
	REDIS = iota
	MYSQL
)

var (
	maxlife int64 = 3600 * 24 * 15
	randChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// 注入cookie
type Manager struct {
	cookieName  string
	provider    SessionProvider
	maxlife     int64
}

func (m *Manager) CreateSessionId() string {
	return randStr(32)
}

func (m *Manager) SesstionStart(req *http.Request, writer http.ResponseWriter) {
	cookie, err := req.Cookie(m.cookieName)
	if err != nil && cookie.Value == "" {
		sid := m.CreateSessionId()

		cookie := http.Cookie{Name: m.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(m.maxlife)}
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)

	}
}

func (*Manager) SesstionDestroy() {

}

// 获取manager
func NewManager (cookieName string, provider SessionProvider, life ...int64) *Manager {

	if len(life) > 0 {
		maxlife = life[0]
	}

	return &Manager{
		cookieName: cookieName,
		provider: provider,
		maxlife: maxlife,
	}
}


func randStr(l int) string {
	le := len(randChars)
	data := make([]byte, l, l)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		data[i] = byte(randChars[rand.Intn(le)])
	}
	return string(data)
}


func redisClient(conf instance.RedisConf) redis.Conn {
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
	return conn
}
