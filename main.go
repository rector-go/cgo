package main

import (
	"cgo/cgo"
	"cgo/controller"
	"log"
	"net/http"
	"time"
)

func main() {
	cgo.InitDB()
	cgo.CreateTable()

	server := &http.Server{
		Addr:        ":8080",
		Handler:     cgo.Router,
		ReadTimeout: 5 * time.Second,
	}
	RegisterRouter(cgo.Router)
	err := server.ListenAndServe()
	if err !=nil{
		log.Panic(err)
	}
}

func RegisterRouter(handler *cgo.RouterHandler) {
	new(controller.UserController).Router(handler)
}
