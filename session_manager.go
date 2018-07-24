package gsession

import (
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	maxlife   int64 = 3600 * 24 * 15
	randChars       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	shard           = "session:"
)

// 注入cookie
type Manager struct {
	cookieName string
	lock       sync.RWMutex
	provider   SessionProvider
	maxlife    int64
}

func (m *Manager) CreateSessionId() string {
	return randStr(32)
}

func (m *Manager) SesstionStart(req *http.Request, writer http.ResponseWriter) (sess Session) {
	// 减少db压力
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := req.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid := shard + m.CreateSessionId()
		sess, err = m.provider.SessionInit(sid, m.maxlife)
		if err != nil {
			panic(err)
		}
		newCookie := http.Cookie{Name: m.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(m.maxlife)}
		http.SetCookie(writer, &newCookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		sess, err = m.provider.SessionRead(sid)
		if err != nil {
			m.SesstionDestroy(req, writer)
			panic(err)
		}
	}
	return
}

func (m *Manager) SesstionDestroy(req *http.Request, writer http.ResponseWriter) {
	cookie, err := req.Cookie(m.cookieName)
	if err != nil {
		panic(err)
	}
	if cookie.Value == "" {
		return
	} else {
		m.provider.SessionDestroy(cookie.Value)
		newCookie := http.Cookie{Name: m.cookieName, Path: "/", HttpOnly: true, MaxAge: -1}
		http.SetCookie(writer, &newCookie)
	}
}

func (m *Manager) GC() {
	time.AfterFunc(time.Duration(m.maxlife), func() {
		m.provider.SessionGC()
	})
}

// 获取manager
func NewManager(cookieName string, provider SessionProvider, life ...int64) *Manager {
	if len(life) > 0 {
		maxlife = life[0]
	}

	managerInstance := &Manager{
		cookieName: cookieName,
		provider:   provider,
		maxlife:    maxlife,
	}
	managerInstance.GC()

	return managerInstance
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
