package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/service"
	"strconv"
)

var commentService = service.NewCommentServiceInstance()

// 视频评论响应结构体
type CommentActionResponse struct {
	entity.Response
	Comment entity.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
// 视频评论接口
// 接受两个查询参数：视频id、token、评论操作类型、评论id、评论内容
// 返回包含 CommentActionFunc 函数结果的JSON作为响应
func CommentAction(c *gin.Context) {
	c.JSON(http.StatusOK, CommentActionFunc(
		c.Query("video_id"),
		c.Query("token"),
		c.Query("action_type"),
		c.Query("comment_id"),
		c.Query("comment_text"),
	))
}

// CommentActionFunc
// 输入：videoId, token, actionType, commentId, text
// CommentActionResponse类型
func CommentActionFunc(videoId, token, actionType, commentId, text string) CommentActionResponse {
	//
	log.Printf("将视频id转换为整型")
	vid, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		log.Printf("视频id转换为整型失败！")
		return ErrorCommentResponse(err)
	}
	if actionType == "1" {
		log.Printf("评论操作类型为1，下面开始执行评论的操作")
		comment, err := commentService.Add(vid, token, text)
		if err != nil {
			log.Printf("commentService.Add(vid, token, text)方法执行有误，评论添加失败！")
			return ErrorCommentResponse(err)
		}
		if comment == nil {
			log.Printf("评论内容为空，请检查评论内容！")
			return FailCommentResponse("Comments are not allowed to be empty! ")
		}
		return CommentActionResponse{
			Response: entity.Response{
				StatusCode: 0,
				StatusMsg:  "Add comment success! ",
			},
			Comment: *comment,
		}
	} else if actionType == "2" {
		log.Printf("评论操作类型为2，下面开始执行取消评论的操作")
		cid, err := strconv.ParseInt(commentId, 10, 64)
		if err != nil {
			log.Printf("评论id转换为整型操作失败！")
			return ErrorCommentResponse(err)
		}
		comment, err := commentService.Withdraw(cid, token, vid)
		if err != nil {
			log.Printf("strconv.ParseInt(commentId, 10, 64)方法执行有误，评论取消失败！")
			return ErrorCommentResponse(err)
		}
		if comment == nil {
			return FailCommentResponse("Withdraw failed, Please try again later! ")
		}
		return CommentActionResponse{
			Response: entity.Response{
				StatusCode: 0,
				StatusMsg:  "Withdraw comment success! ",
			},
			Comment: *comment,
		}
	} else {
		return CommentActionResponse{
			Response: entity.Response{
				StatusCode: 1,
				StatusMsg:  "Service Wrong!",
			},
		}
	}
}

// ErrorCommentResponse 评论操作错误
func ErrorCommentResponse(err error) CommentActionResponse {
	return CommentActionResponse{
		Response: entity.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		},
	}
}

// FailCommentResponse 评论操作失败
func FailCommentResponse(msg string) CommentActionResponse {
	return CommentActionResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  msg,
		},
	}
}
