package gsession

type Session interface {
	Set(key string, value interface{}) error // set session value
	SetMulti(map[string]interface{}) error
	Get(key string) (interface{}, error)  // get session value
	GetMulti(value interface{}) error
	Delete(key string) error     // delete session value
	SessionID() string                // get current sessionID
	Clear() error
}

type SessionProvider interface {
	SessionInit(sid string, maxLife int64) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC()
}


type CommonConf struct {
	Host        string
	Port        uint16
	UserName    string
	Password    string
	DB          string
}