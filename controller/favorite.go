package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/service"
	"strconv"
)

var favoriteService = service.NewFavoriteServiceInstance()

// 点赞操作响应结构体
type FavoriteActionResponse struct {
	entity.Response
}

// FavoriteAction
// 点赞接口
// 接受三个查询参数：视频id， token， 点赞操作类型
// 返回包含 FavoriteActionFunc 函数结果的JSON作为响应
// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	c.JSON(http.StatusOK, FavoriteActionFunc(
		c.Query("video_id"),
		c.Query("token"),
		c.Query("action_type"),
	))
}

// FavoriteActionFunc
// 输入：videoId，token 和 actionType
// 输出：FavoriteActionResponse类型
func FavoriteActionFunc(videoId, token, actionType string) FavoriteActionResponse {
	// 将 videoId 解析为整数类型
	vid, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		log.Printf("视频Id解析整数失败！")
		return ErrorFavoriteResponse(err)
	}

	// 如果点赞类型为1，进行点赞
	if actionType == "1" {
		err = favoriteService.Add(vid, token)
		if err != nil {
			log.Printf("favoriteService.Add(vid, token)方法出错，点赞操作失败！")
			return ErrorFavoriteResponse(err)
		}
		return FavoriteActionResponse{
			Response: entity.Response{
				StatusCode: 0,
				StatusMsg:  "Thanks for your favorite! ",
			},
		}
	} else if actionType == "2" { // 如果点赞类型为2，进行取消点赞操作
		err := favoriteService.Withdraw(vid, token)
		if err != nil {
			log.Printf("favoriteService.Withdraw(vid, token)方法出错，取消赞操作失败！")
			return ErrorFavoriteResponse(err)
		}
		return FavoriteActionResponse{
			Response: entity.Response{
				StatusCode: 0,
				StatusMsg:  "Please Favorite Next Time! ",
			},
		}
	} else { // 如果点赞类型不为1，不为2，进行错误响应
		return FavoriteActionResponse{
			Response: entity.Response{
				StatusCode: 1,
				StatusMsg:  "Service Wrong!",
			},
		}
	}
}

// 错误类型响应
func ErrorFavoriteResponse(err error) FavoriteActionResponse {
	return FavoriteActionResponse{
		Response: entity.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		},
	}
}
