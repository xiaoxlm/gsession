package tests

import (
	"gsession"
	"gsession/instance"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"fmt"
)

var (
	cookieName     = "cookieset"
	sessionManager *gsession.Manager

	username = "xxx"
	tel      = "123"

	getInfoValue string
)

func TestCookie(t *testing.T) {
	// sessionInit
	http.HandleFunc("/login", login)
	req, _ := http.NewRequest("GET", "/login", nil)
	respWrite := httptest.NewRecorder()
	login(respWrite, req)

	cookie := respWrite.HeaderMap.Get("Set-Cookie")
	strSlice := strings.Split(cookie, ";")
	cookie = strSlice[0]
	cookieValue := strings.Split(cookie, "=")[1]

	if !strings.Contains(cookieValue, "session") {
		t.Fatal("cookieValue incorrect")
	}

	//get session
	http.HandleFunc("/getUsername", getUsername)
	req, _ = http.NewRequest("GET", "/getUsername", nil)
	req.Header.Set("Cookie", cookie)
	respWrite = httptest.NewRecorder()
	getUsername(respWrite, req)
	if getInfoValue != username {
		t.Fatal("username incorrect")
	}

	//multi
	http.HandleFunc("/getMulti", getMulti)
	req, _ = http.NewRequest("GET", "/getMulti", nil)
	req.Header.Set("Cookie", cookie)
	respWrite = httptest.NewRecorder()
	getMulti(respWrite, req)
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessionManager.SesstionStart(r, w)
	session.Set("username", username)
	session.Set("tel", tel)
}

func getUsername(w http.ResponseWriter, r *http.Request) {
	//getCookie := r.Header.Get("Cookie")
	//fmt.Println(getCookie)
	//co, _ := r.Cookie(cookieName)
	//fmt.Println(co.Name, co.Value)
	session := sessionManager.SesstionStart(r, w)
	val, _ := session.Get("username")
	getInfoValue = val.(string)
}

func getMulti(w http.ResponseWriter, r *http.Request) {
	type multi struct {
		Username string `redis:"username"`
		Tel      string `redis:"tel"`
	}
	m := multi{}
	session := sessionManager.SesstionStart(r, w)
	session.GetMulti(&m)
	fmt.Println(m)
}

func init() {
	redisProvider := instance.RedisProvider{}
	redisProvider.SetConn(gsession.CommonConf{
		Host: "127.0.0.1",
		Port: 6379,
	})
	sessionManager = gsession.NewManager(cookieName, &redisProvider)
}
