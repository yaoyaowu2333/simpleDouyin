package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"simpleDouyin/entity"
	"simpleDouyin/service"
	"simpleDouyin/utils"
)

// 是否自动生成封面，需要配置环境，默认为否
// var useGeneratedCover = utils.UseGeneratedCover
var useGeneratedCover = true

//var useGeneratedCover = false

// Publish POST /publish/action/
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, entity.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, PublishFunc(token, title, data, c))
}

func PublishFunc(token, title string, data *multipart.FileHeader, c *gin.Context) entity.Response {
	//检查token是否为空
	if token == "" {
		return ErrorResponse(utils.Error{Msg: "token不能为空"})
	}
	//检查title是否为空
	if token == "" {
		return ErrorResponse(utils.Error{Msg: "title不能为空"})
	}
	//检查文件是否为空
	if data == nil {
		return ErrorResponse(utils.Error{Msg: "empty data file"})
	}
	//检查后缀名
	ext := filepath.Ext(data.Filename)
	if ext != ".mp4" {
		return ErrorResponse(utils.Error{Msg: "unsupported file extension"})
	}
	//获取上传视频的文件名,不包含路径。
	filepath.Base(data.Filename)
	//生成一个随机UUID作为文件名的一部分。
	fileName := utils.GenerateUUID()
	//生成视频文件名与封面文件名
	videoFileName := fmt.Sprintf("%s%s", fileName, ext)
	coverName := fmt.Sprintf("%s%s", fileName, ".jpg")

	var dir = "./public/"
	_, err := os.Stat(dir)
	// 判断文件夹是否存在,
	if os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	//保存文件
	saveFile := filepath.Join(dir, videoFileName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		log.Printf("视频文件保存失败\n")
		return ErrorResponse(err)
	}
	log.Printf("视频文件保存成功\n")
	//生成视频url信息
	// TODO 目前是数据库硬编码ip:port，后续可改成动态
	playUrl := utils.VideoUrlPrefix + videoFileName
	//封面
	var coverUrl string
	if useGeneratedCover {
		//生成封面url信息
		coverUrl = utils.VideoUrlPrefix + coverName
		//提取视频第一帧并保存封面
		coverFilePath := filepath.Join(dir, coverName)
		utils.ReadFrameAsJpeg(saveFile, 1, coverFilePath)
	} else {
		coverUrl = utils.DefaultCoverUrl
	}
	log.Printf("视频封面保存成功")
	err = service.NewVideoServiceInstance().Publish(token, playUrl, coverUrl, title)
	if err != nil {
		return ErrorResponse(err)
	}
	return entity.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}
}

func ErrorResponse(err error) entity.Response {
	return entity.Response{
		StatusCode: 1,
		StatusMsg:  err.Error(),
	}
}
