package service

import (
	"log"
	"simpleDouyin/dao"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
	"strconv"
	"sync"
	"time"
)

type VideoService struct {
}

var videoService *VideoService
var serviceOnce sync.Once

func NewVideoServiceInstance() *VideoService {
	serviceOnce.Do(
		func() {
			videoService = &VideoService{}
		})
	return videoService
}

func (s *VideoService) FindVideoById(id int64) (*entity.Video, error) {
	videoModel, err := dao.NewVideoDaoInstance().QueryVideoById(id)
	if err != nil {
		return nil, err
	}

	if videoModel == nil {
		return nil, nil
	}

	userModel, err := dao.NewUserDaoInstance().QueryUserById(videoModel.AuthorId)
	if err != nil {
		return nil, err
	}

	user := pack.User(userModel)
	video := pack.Video(videoModel)

	video.Author = *user
	return video, nil
}

// Feed
// 通过传入时间戳，用户token，返回对应的视频数组，以及视频数组中最早的发布时间
// 获取视频数组大小是可以控制的，在utils中的DefaultLimit变量
func (s *VideoService) Feed(latestTime int64, token string, limit int) (*int64, []*entity.Video, error) {
	var lastTime time.Time
	if latestTime == 0 {
		lastTime = time.Now()
	} else {
		lastTime = time.UnixMilli(latestTime)
	}
	log.Printf("获取到时间戳%v\n", lastTime)
	videoModels, err := dao.NewVideoDaoInstance().QueryVideoBeforeTime(lastTime, limit)
	if err != nil {
		log.Printf("方法dao.QueryVideoBeforeTime(lastTime,limit) 失败：%v\n", err)
		return nil, nil, err
	}
	log.Printf("方法dao.QueryVideoBeforeTime(lastTime,limit) 成功：%v\n", videoModels)
	// 获取返回视频的作者id
	authorIds := pack.AuthorIds(videoModels)
	// 依据id获取用户信息
	userModelMap, err := dao.NewUserDaoInstance().MQueryUserById(authorIds)
	if err != nil {
		log.Printf("方法dao.MQueryUserById(authorIds) 失败：%v\n", err)
		return nil, nil, err
	}
	log.Printf("方法dao.MQueryUserById(authorIds) 成功：%v\n", userModelMap)
	// 将userModelMap数据通过MUser进行处理，在拷贝的过程中对数据进行组装
	userMap := pack.MUser(userModelMap)

	// 获取当前用户
	curUserId, err := dao.NewLoginStatusDaoInstance().QueryUserIdByToken(token)
	if err != nil {
		log.Printf("方法dao.QueryUserIdByToken(token) 失败：%v\n", err)
		return nil, nil, err
	}
	log.Printf("方法dao.QueryUserIdByToken(token) 成功：%v\n", curUserId)
	if curUserId != -1 {
		for uid := range userMap {
			//用户是否关注视频作者,这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
			userMap[uid].IsFollow = dao.NewRelationDaoInstance().IsFollow(curUserId, uid)
		}
	}
	// 将videoModels数据通过Videos进行处理，在拷贝的过程中对数据进行组装
	videos := pack.Videos(videoModels)

	for i, video := range videos {
		//插入Author
		video.Author = *userMap[authorIds[i]]

		//获取该视频的评论数字
		commentCount, _, err := dao.NewCommentDaoInstance().QueryCommentByVideoId(video.Id)
		if err != nil {
			log.Printf("方法dao.QueryCommentByVideoId(video.Id) 失败：%v\n", err)
			return nil, nil, err
		}
		log.Printf("方法dao.QueryCommentByVideoId(video.Id) 成功：%v\n", commentCount)
		video.CommentCount = commentCount

		//获取该视频的点赞数
		favoriteCount, err := dao.NewFavoriteDaoInstance().QueryFavoriteByVideoId(video.Id)
		if err != nil {
			log.Printf("方法dao.QueryFavoriteByVideoId(video.Id) 失败：%v\n", err)
			return nil, nil, err
		}
		log.Printf("方法dao.QueryFavoriteByVideoId(video.Id) 成功：%v\n", favoriteCount)
		video.FavoriteCount = favoriteCount

		//获取当前用户是否点赞了该视频
		video.IsFavorite = dao.NewFavoriteDaoInstance().QueryFavoriteByUserToken(video.Id, token)
	}

	var nextTime int64
	if len(videoModels) > 0 {
		//获得视频中最早的时间返回
		nextTime = videoModels[len(videoModels)-1].CreateAt.UnixMilli()
	} else {
		nextTime = time.Now().UnixMilli()
	}
	println("latest time: " + strconv.FormatInt(latestTime, 10))
	println("next time: " + strconv.FormatInt(nextTime, 10))

	return &nextTime, videos, nil
}

func (s *VideoService) PublishList(authorId int64) ([]*entity.Video, error) {
	// invalid authorId
	if authorId <= 0 {
		return nil, nil
	}

	videoModels, err := dao.NewVideoDaoInstance().QueryVideoByAuthorId(authorId)
	if err != nil {
		return nil, err
	}
	authorIds := pack.AuthorIds(videoModels)

	userModelMap, err := dao.NewUserDaoInstance().MQueryUserById(authorIds)
	if err != nil {
		return nil, err
	}

	userMap := pack.MUser(userModelMap)
	videos := pack.Videos(videoModels)

	for i, video := range videos {
		video.Author = *userMap[authorIds[i]]
	}

	return videos, nil
}

// Publish
// 将视频信息存入视频数据库,视频数加一
func (s VideoService) Publish(token, playUrl, coverUrl, title string) error {
	if playUrl == "" || coverUrl == "" || title == "" {
		return utils.Error{Msg: "参数不能为空"}
	}
	// 查询用户
	userId, err := dao.NewLoginStatusDaoInstance().QueryUserIdByToken(token)
	if err != nil {
		return err
	}
	if userId == -1 {
		return utils.Error{Msg: "user not exist"}
	}

	// 保存 video
	videoModel := dao.Video{
		AuthorId:      userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		Title:         title,
		CreateAt:      time.Now(),
		FavoriteCount: 0,
		CommentCount:  0,
	}
	//将视频信息存入视频数据库
	err = dao.NewVideoDaoInstance().CreateVideo(&videoModel)
	if err != nil {
		return err
	}
	// 用户的视频数增加
	err = dao.NewUserDaoInstance().IncreaseVideoCountByOne(userId)
	if err != nil {
		return err
	}
	return nil
}
