package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
	"strconv"
)

// 评论列表获取响应结构体
type CommentListResponse struct {
	entity.Response
	CommentList []entity.Comment `json:"comment_list,omitempty"`
}

// CommentList ..
// 评论列表获取接口
// 输入：查询到的视频id与用户token
// 输出：CommentListFunc函数的返回值  的json形式
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListFunc(
		c.Query("video_id"),
		c.Query("token"),
	))
}

// 输入：videoID, token
// 输出：CommentListResponse
func CommentListFunc(videoID, token string) CommentListResponse {
	// TODO 使用token进行鉴权
	log.Printf("开始对token  %v鉴权", token)
	if videoID == "" {
		log.Printf("token为空！")
		return ErrorCommentListResponse(utils.Error{Msg: "empty token or user_id"})
	}
	log.Printf("将视频id转换为整型")
	vid, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Printf("视频id转换整型失败！")
		return ErrorCommentListResponse(err)
	}
	log.Printf("开始根据视频id查询所有评论")
	comments, err := commentService.LoadComments(vid)
	if err != nil {
		log.Printf("commentService.LoadComments(vid)方法执行有误，获取指定id视频的评论列表失败！")
		ErrorCommentListResponse(err)
	}
	return CommentListResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "Load comments success!",
		},
		CommentList: pack.CommentsPtrs(comments),
	}
}

// 获取评论列表错误响应
func ErrorCommentListResponse(err error) CommentListResponse {
	return CommentListResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		CommentList: nil,
	}
}
