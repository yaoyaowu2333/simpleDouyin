package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/service"
	"strconv"
)

var userService = service.NewUserServiceInstance()

type UserResponse struct {
	entity.Response
	User entity.User `json:"user"`
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, UserInfoFunc(
		c.Query("user_id"),
		c.Query("token"),
	))
}

func UserInfoFunc(userId, token string) UserResponse {
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return ErrorUserResponse(err)
	}
	log.Printf("查询用户的id: %v\n", uid)
	user, err := service.NewUserServiceInstance().UserInfo(uid)
	if err != nil {
		log.Printf("方法userService.UserInfo(uid)失败: %v\n", err)
		return ErrorUserResponse(err)
	}
	log.Printf("方法userService.UserInfo(uid)成功: %v\n", user)
	if user == nil {
		return FailUserResponse("user not exist: uid " + strconv.FormatInt(uid, 10))
	}
	return UserResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		User: *user,
	}
}

func ErrorUserResponse(err error) UserResponse {
	return UserResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
	}
}

func FailUserResponse(msg string) UserResponse {
	return UserResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  msg,
		},
	}
}
