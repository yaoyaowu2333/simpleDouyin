package main

import (
	//"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"simpleDouyin/dao"
	"simpleDouyin/service"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)
	err := initDB()
	if err != nil {
		println(err.Error())
		return
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initDB() error {
	return dao.Init(true)
}
