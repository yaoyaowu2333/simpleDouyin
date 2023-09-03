package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/service"
	"simpleDouyin/utils"
	"strconv"
	"log"
)

// 视频列表相应结构体
type VideoListResponse struct {
	entity.Response
	VideoList []entity.Video `json:"video_list"`
}


// 发布列表接口
// 接受两个查询参数：token 和 user_id，可能用于认证和标识用户      
// 这个token是当前账号使用者的标识，user_id是其他用户的标识 ？？？
// 返回包含 publishListFunc 函数结果的JSON作为响应
// PublishList /douyin/publish
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, publishListFunc(c.Query("token"), c.Query("user_id")))
}


// 输入：token 和 user_id
// 输出：VideoListResponse类型
func publishListFunc(token, userId string) VideoListResponse {
	// TODO 使用token进行鉴权
	if token == "" || userId == "" {		// 如果token 和 userId 为空，返回一个错误响应
		log.Printf("错误，token或者用户id为空！")
		return ErrorVideoListResponse(utils.Error{Msg: "empty token or user_id"})
	}
	uid, err := strconv.ParseInt(userId, 10, 64)		// 将 userId 从字符串转换为整数uid
	if err != nil {
		log.Printf("用户id转化为整型出错！")
		return ErrorVideoListResponse(err)
	}


	// 调用一个服务（service.NewVideoServiceInstance().PublishList(uid)）
	// 来获取与给定用户uid相关联的视频列表
	// NewVideoServiceInstance见service/video.go文件
	// PublishList见service/video.go文件中的函数PublishList
	videos, err := service.NewVideoServiceInstance().PublishList(uid)
	if err != nil {
		log.Printf("方法service.NewVideoServiceInstance().PublishList(uid) 失败")
		return ErrorVideoListResponse(err)
	}

	return VideoListResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: pack.VideoPtrs(videos),
	}

}


// 视频列表获取错误的响应
func ErrorVideoListResponse(err error) VideoListResponse {
	return VideoListResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		VideoList: nil,
	}
}
