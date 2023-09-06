package dao

import (
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Video struct {
	Id            int64
	AuthorId      int64
	PlayUrl       string
	CoverUrl      string
	Title         string
	CreateAt      time.Time
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

// NewVideoDaoInstance Singleton
// 初始化一个*VideoDao类型的对象
func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

// QueryVideoById will return nil if no user is found
// 根据视频id查询视频
func (*VideoDao) QueryVideoById(id int64) (*Video, error) {
	var video Video
	err := db.Where("id = ?", id).First(&video).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Fatal("find video by id err:" + err.Error())
		return nil, err
	}
	return &video, nil
}

// QueryVideoBeforeTime will return empty array if no user is found
// 依据一个时间，来获取这个时间之前的一些视频
func (*VideoDao) QueryVideoBeforeTime(time time.Time, limit int) ([]*Video, error) {
	var videos []*Video
	err := db.Where("create_at < ?", time).Order("create_at DESC").Limit(limit).Find(&videos).Error

	if err != nil {
		log.Fatal("batch find video before time err:" + err.Error())
		return nil, err
	}
	return videos, nil
}

// CreateVideo
// 将视频信息存入视频表中
func (*VideoDao) CreateVideo(video *Video) error {
	return db.Create(&video).Error
}

// QueryVideoByAuthorId
// 返回数据库中获取的给定作者ID相关的视频列表
// 输入：*VideoDao类型
// 输出：视频列表
func (*VideoDao) QueryVideoByAuthorId(authorId int64) ([]*Video, error) {
	var videos []*Video
	// 查询数据库中具有给定作者ID的视频。db 是一个数据库连接对象
	err := db.Where("author_id = ?", authorId).Find(&videos).Error
	if err != nil {
		log.Printf("查询数据库中给定作者ID的视频失败！")
		log.Fatal("batch find video by author_id err:" + err.Error())
		return nil, err
	}
	return videos, nil
}

// UpdateCommentByID
// 更新评论数
func (*VideoDao) UpdateCommentByID(id int64, count int64) error {
	err := db.Model(&Video{}).Where("id = ?", id).UpdateColumn("comment_count", count).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFavoriteByID
// 更新点赞数
func (*VideoDao) UpdateFavoriteByID(id int64, actionType int64) error {
	var video Video
	err := db.Where("id = ?", id).First(&video).Error
	if err == gorm.ErrRecordNotFound {
		log.Printf("数据库查询指定id的用户，查询失败！")
		return err
	}
	if actionType == 1 {
		video.FavoriteCount += 1
	} else {
		video.FavoriteCount -= 1
	}
	return db.Save(&video).Error
}
