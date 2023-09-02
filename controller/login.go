package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
)

type UserLoginResponse struct {
	entity.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	c.JSON(http.StatusOK, LoginFunc(username, password))
}

// Register POST /douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	c.JSON(http.StatusOK, RegisterFunc(username, password))
}

func LoginFunc(username, password string) UserLoginResponse {
	userId, token, err := userService.Login(username, password)
	if err != nil {
		log.Printf("方法userService.Login(username, password) 失败：%v", err)
		return ErrorUserLoginResponse(err)
	}
	log.Printf("方法userService.Login(username, password)成功\n")
	return UserLoginResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: *userId,
		Token:  *token,
	}
}

func RegisterFunc(username, password string) UserLoginResponse {
	if err := userService.Register(username, password); err != nil {
		log.Printf("方法userService.Register(username, password) 失败：%v", err)
		return ErrorUserLoginResponse(err)
	}
	log.Printf("方法userService.Register(username, password)成功\n")
	return LoginFunc(username, password)
}

func ErrorUserLoginResponse(err error) UserLoginResponse {
	return UserLoginResponse{
		Response: entity.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		},
	}
}
