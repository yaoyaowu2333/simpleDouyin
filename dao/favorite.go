package dao

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Favorite struct {
	Id        int64
	UserToken string
	VideoId   int64
	CreateAt  time.Time
}

type FavoriteDao struct {
}

var favoriteDao *FavoriteDao
var favoriteOnce sync.Once

// NewFavoriteDaoInstance Singleton
func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDao{}
		})
	return favoriteDao
}

// QueryFavoriteByVideoId
// 获取点赞数
func (d *FavoriteDao) QueryFavoriteByVideoId(videoID int64) (int64, error) {
	var favoriteCount int64
	result := db.Table("videos").Select("favorite_count").Where("id = ?", videoID).Find(&favoriteCount)
	err := result.Error
	if err != nil {
		return 0, err
	}
	return favoriteCount, nil
}


// 根据用户token查询用户点赞的视频的ID列表
func (d *FavoriteDao) QueryVideoIdByToken(token string) ([]int64, error) {
	var ids []int64
	err := db.Select("video_id").Table("favorites").Where("user_token = ?", token).Find(&ids).Error
	if err != nil {
		log.Printf("数据库查询指定用户点赞视频的列表失败！！")
		return nil, err
	}
	return ids, nil
}

// QueryFavoriteByUserToken
// 登录用户是否点赞该视频
func (d *FavoriteDao) QueryFavoriteByUserToken(videoId int64, token string) bool {
	err := db.Where("video_id = ? AND user_token = ?", videoId, token).First(&Favorite{}).Error
	if err != nil {
		return false
	}
	return true
}


// 用于保存点赞（收藏）记录到数据库中
func (d *FavoriteDao) Save(favorite *Favorite) error {
	result := db.Create(&favorite)
	err := result.Error
	if err != nil {
		log.Printf("数据库创建点赞数据操作失败！")
		return err
	}

	// 更新与点赞相关的视频记录
	err = db.Debug().Model(&Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if err != nil {
		log.Printf("数据库更新点赞信息操作失败！")
		fmt.Println(err)
		return err
	}
	return nil
}

// 用户删除某条点赞记录
func (d *FavoriteDao) Delete(videoId int64, token string) error {
	err := db.Where("user_token = ? AND video_id = ?", token, videoId).Delete(&Favorite{}).Error
	if err != nil {
		log.Printf("数据库查询指定用户token以及视频id对应点赞的操作失败！")
		return err
	}

	// 更新数据库
	err = db.Debug().Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	if err != nil {
		log.Printf("数据库更新失败！")
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *FavoriteDao) Total() (int64, error) {
	// 获取全部记录
	var count int64
	result := db.Table("comments").Count(&count)
	err := result.Error
	if err != nil {
		log.Fatal("total user err:" + err.Error())
		return -1, err
	}
	return count, nil
}

func (d *FavoriteDao) MaxId() (int64, error) {
	// 获取全部记录
	var lastRec *Comment
	result := db.Table("favorites").Last(&lastRec)
	err := result.Error
	if err != nil {
		//log.Fatal("max id err:" + err.Error())
		return 0, err
	}
	return lastRec.Id, nil
}
