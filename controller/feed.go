package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/service"
	"simpleDouyin/utils"
	"strconv"
	"time"
)

type FeedResponse struct {
	entity.Response
	VideoList []entity.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

// Feed GET /feed/
func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedFunc(c.Query("latest_time"), c.Query("token")))
}

func FeedFunc(latestTime string, token string) FeedResponse {
	log.Printf("传入的时间" + latestTime + "\n")
	timeInt, _ := strconv.ParseInt(latestTime, 10, 64)
	log.Printf("获取到用户token:%v\n", token)
	nextTime, videos, err := service.NewVideoServiceInstance().Feed(timeInt, token, utils.DefaultLimit)
	// service层出错
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, token, utils.DefaultLimit) 失败：%v\n", err)
		return ErrorFeedResponse(err)
	}
	log.Printf("方法videoService.Feed(lastTime, token, utils.DefaultLimit) 成功\n")
	return FeedResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: pack.VideoPtrs(videos),
		NextTime:  *nextTime,
	}
}

// ErrorFeedResponse
// 如果service层出错,返回状态码-1
func ErrorFeedResponse(err error) FeedResponse {
	return FeedResponse{
		Response: entity.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		},
		VideoList: nil,
		NextTime:  time.Now().UnixMilli(),
	}
}
