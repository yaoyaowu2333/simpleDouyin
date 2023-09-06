package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
)

// 点赞列表响应结构体
type FavoriteListResponse struct {
	entity.Response
	VideoList []entity.Video `json:"video_list"`
}

// FavoriteList ..
// 点赞列表接口
// 接受一个查询参数：token
// 返回包含 FavoriteListFunc 函数结果的JSON作为响应
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, FavoriteListFunc(
		c.Query("token"),
	))
}

// 输入：token
// 输出：FavoriteListResponse类型
func FavoriteListFunc(token string) FavoriteListResponse {
	// TODO 使用token进行鉴权
	if token == "" {
		log.Printf("用户token为空！")
		return ErrorFavoriteListResponse(utils.Error{Msg: "empty token"})
	}
	// 获取与该用户关联的点赞视频列表
	videos, err := favoriteService.FavoriteList(token)
	if err != nil {
		log.Printf("favoriteService.FavoriteList(token)方法有误，获取用户点赞的视频列表失败！")
		ErrorFavoriteListResponse(err)
	}
	return FavoriteListResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "Load Favorites success!",
		},
		VideoList: pack.VideoPtrs(videos),
	}
}

// 获取点赞视频列表错误响应
func ErrorFavoriteListResponse(err error) FavoriteListResponse {
	return FavoriteListResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		VideoList: nil,
	}
}
