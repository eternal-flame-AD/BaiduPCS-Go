package pcsweb

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/eternal-flame-AD/BaiduPCS-Go/internal/pcsconfig"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var (
	clientNonce = make([]string, 0)
	clientToken = make([]string, 0)
	stripSign   = regexp.MustCompile("&sign=[^&]+")
)

func issueToken() string {
	token := randStringBytes(15)
	clientToken = append(clientToken, token)
	return token
}

func inStringArray(a string, b []string) bool {
	for _, val := range b {
		if val == a {
			return true
		}
	}
	return false
}

func checkAuth(r *http.Request) bool {
	fmt.Println(r.Header)
	token, err := r.Cookie("token")
	if token == nil || !inStringArray(token.Value, clientToken) || err != nil {
		return false
	}
	return true
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := boxTmplParse("login", "login.html")
	checkErr(tmpl.Execute(w, r.Form.Get("path")))
}
func loginAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nonce := r.Form.Get("nonce")
	sign := r.Form.Get("sign")
	if nonce == "" || inStringArray(nonce, clientNonce) {
		w.WriteHeader(400)
		w.Write((&ErrInfo{ErrroCode: 31023, ErrorMsg: "incorrect parameter"}).JSON())
		return
	}
	m5 := md5.New()
	m5.Write([]byte(pcsconfig.Config.GetWebUsername() + ":" + pcsconfig.Config.GetWebPassword() + ":" + nonce))
	if sign == hex.EncodeToString(m5.Sum(nil)) {
		w.Write([]byte("{\"error_code\":0,\"token\":\"" + issueToken() + "\"}"))
		return
	}
	w.WriteHeader(403)
	w.Write((&ErrInfo{ErrroCode: 31023, ErrorMsg: "incorrect parameter"}).JSON())
	return
}
