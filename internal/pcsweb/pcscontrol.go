package pcsweb

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eternal-flame-AD/BaiduPCS-Go/internal/pcsconfig"
)

func toJSON(s interface{}) []byte {
	res, err := json.Marshal(s)
	if err != nil {
		return (&ErrInfo{
			ErrroCode: 1,
			ErrorMsg:  err.Error(),
		}).JSON()
	}

	return res
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	userData := pcsconfig.Config.ActiveUser()

	w.Write(toJSON(userData))
}

func userList(w http.ResponseWriter, r *http.Request) {
	users := pcsconfig.Config.BaiduUserList()

	w.Write(toJSON(users))
}

func userChange(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uid, err := strconv.ParseUint(r.Form.Get("uid"), 10, 64)

	if err != nil {
		w.WriteHeader(403)
		w.Write((&ErrInfo{ErrroCode: 31045, ErrorMsg: "user not exists"}).JSON())
		return
	}

	_, err = pcsconfig.Config.SwitchUser(&pcsconfig.BaiduBase{
		UID: uid,
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write((&ErrInfo{ErrroCode: 31209, ErrorMsg: err.Error()}).JSON())
		return
	}

	w.Write((&ErrInfo{ErrroCode: 0, ErrorMsg: "no error"}).JSON())
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	bduss := r.Form.Get("bduss")

	_, err := pcsconfig.Config.SetupUserByBDUSS(bduss, "", "")

	if err != nil {
		w.WriteHeader(500)
		w.Write((&ErrInfo{ErrroCode: 31209, ErrorMsg: err.Error()}).JSON())
		return
	}

	w.WriteHeader(501)
	w.Write((&ErrInfo{ErrroCode: 31210, ErrorMsg: "not implemented"}).JSON())
}

func userLogout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uid, err := strconv.ParseUint(r.Form.Get("uid"), 10, 64)

	if err != nil {
		w.WriteHeader(403)
		w.Write((&ErrInfo{ErrroCode: 31045, ErrorMsg: "user not exists"}).JSON())
		return
	}

	_, err = pcsconfig.Config.DeleteUser(&pcsconfig.BaiduBase{
		UID: uid,
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write((&ErrInfo{ErrroCode: 31209, ErrorMsg: err.Error()}).JSON())
		return
	}

	w.Write((&ErrInfo{ErrroCode: 0, ErrorMsg: "no error"}).JSON())
}
