// Package pcsweb web前端包
package pcsweb

import (
	"fmt"
	"net/http"

	"github.com/GeertJohan/go.rice"
)

var (
	staticBox    *rice.Box // go.rice 文件盒子
	templatesBox *rice.Box // go.rice 文件盒子
)

func boxInit() (err error) {
	staticBox, err = rice.FindBox("static")
	if err != nil {
		return
	}

	templatesBox, err = rice.FindBox("template")
	if err != nil {
		return
	}

	return nil
}

// StartServer 开启web服务
func StartServer(listen string, port uint) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}

	err := boxInit()
	if err != nil {
		return err
	}

	http.HandleFunc("/", rootMiddleware)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox())))
	http.HandleFunc("/about.html", middleware(aboutPage))
	http.HandleFunc("/login.html", middleware(loginPage))
	http.HandleFunc("/index.html", activeAuthMiddleware(indexPage))
	http.HandleFunc("/tasks.html", activeAuthMiddleware(tasksPage))

	http.HandleFunc("/cgi-bin/baidu/pcs/login", middleware(loginAPI))
	http.HandleFunc("/cgi-bin/baidu/pcs/file/list", activeAuthMiddleware(fileList))
	http.HandleFunc("/cgi-bin/baidu/pcs/file/download", activeAuthMiddleware(downloadProcess))
	http.HandleFunc("/cgi-bin/baidu/pcs/user/info", activeAuthMiddleware(userInfo))
	http.HandleFunc("/cgi-bin/baidu/pcs/user/list", activeAuthMiddleware(userList))
	http.HandleFunc("/cgi-bin/baidu/pcs/user/change", activeAuthMiddleware(userChange))
	http.HandleFunc("/cgi-bin/baidu/pcs/user/login", activeAuthMiddleware(userLogin))
	http.HandleFunc("/cgi-bin/baidu/pcs/user/logout", activeAuthMiddleware(userLogout))
	http.HandleFunc("/cgi-bin/baidu/pcs/config/saveto", activeAuthMiddleware(downloadDir))

	return http.ListenAndServe(fmt.Sprintf("%s:%d", listen, port), nil)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl := boxTmplParse("index", "index.html", "about.html")
	checkErr(tmpl.Execute(w, nil))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	tmpl := boxTmplParse("index", "index.html", "baidu/userinfo.html")
	checkErr(tmpl.Execute(w, r.Form.Get("path")))
}

func tasksPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	tmpl := boxTmplParse("tasks", "tasks.html", "baidu/taskinfo.html")
	checkErr(tmpl.Execute(w, r.Form.Get("path")))
}
