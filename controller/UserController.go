package controller

import (
	"cgo/cgo"
	"cgo/constant"
	"cgo/service"
	"cgo/utils"
	"log"
	"net/http"
)

type UserController struct {
	cgo.ApiController
}

var userService = new(service.UserService)

func (p *UserController) Router(router *cgo.RouterHandler) {
	router.Router("/register", p.register)
	router.Router("/login", p.login)
	router.Router("/findAll", p.findAll)
	router.Router("/findUser", p.findUser)
}

func (p *UserController) register(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if utils.Empty(username) || utils.Empty(password) {
		cgo.ResultFail(w, "username or password cannot be empty.")
		return
	}
	exist := userService.GetUserByName(username)
	if exist != nil {
		cgo.ResultFail(w, "User already exists.")
		return
	}
	id := userService.Insert(username, password)
	if id <= 0 {
		cgo.ResultFail(w, "Register failed")
		return
	}
	cgo.ResultOk(w, "Register success")
}

func (p *UserController) login(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if utils.Empty(username) || utils.Empty(password) {
		cgo.ResultFail(w, "username or password cannot be empty.")
		return
	}
	users := userService.SelectUserByName(username)
	if len(users) == 0 {
		cgo.ResultFail(w, "user does not exist")
		return
	}
	user := users[0]
	if user.Password != password {
		cgo.ResultFail(w, "password error")
		return
	}
	session := cgo.GlobalSession().SessionStart(w, r)
	session.Set(constant.KEY_USER, user)
	cgo.ResultOk(w, "login success")
}

func (p *UserController) findAll(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("GSESSION")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(cookie.Value)
	}
	users := userService.SelectAllUser()
	cgo.ResultJsonOk(w, users)
}

func (p *UserController) findUser(w http.ResponseWriter, r *http.Request) {
	user := p.GetUser(w, r)
	cgo.ResultJsonOk(w, user)
}
