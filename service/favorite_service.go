package service

// TODO
import (
	"simpleDouyin/dao"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
	"sync"
	"sync/atomic"
	"time"
	"log"
)

type FavoriteService struct {
}

var favoriteService *FavoriteService
var favoriteOnce sync.Once

func NewFavoriteServiceInstance() *FavoriteService {
	favoriteOnce.Do(
		func() {
			favoriteService = &FavoriteService{}
		})
	return favoriteService
}

func (s *FavoriteService) FindUserByToken(token string) (*entity.User, error) {
	user, err := dao.NewUserDaoInstance().QueryUserByToken(token)
	if err != nil {
		return nil, err
	}
	return pack.User(user), err
}

// 输入：token
// 输出：用户点赞的视频切片和执行报错
func (s *FavoriteService) FindVideosByToken(token string) ([]*entity.Video, error) {
	// invalid token
	if token == "" {
		log.Printf("用户token为空！失败！")
		return nil, nil
	}
	// 查询与用户点赞的视频ID列表
	videoIds, err := dao.NewFavoriteDaoInstance().QueryVideoIdByToken(token)
	if err != nil {
		log.Printf("dao.NewFavoriteDaoInstance().QueryVideoIdByToken(token)方法有误，查询与用户点赞的视频ID列表失败！")
		return nil, err
	}
	var videos []*entity.Video
	for _, id := range videoIds {  // 将查询到的视频对象 video 添加到 videos 列表中，最终形成了用户收藏的视频列表
		// 查询视频的详细信息
		video, _ := NewVideoServiceInstance().FindVideoById(id)
		//video.IsFavorite = true
		videos = append(videos, video)
	}
	return videos, nil
}

func (s *FavoriteService) TotalComment() (int64, error) {
	count, err := dao.NewFavoriteDaoInstance().Total()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (s *FavoriteService) LastId() (int64, error) {
	count, err := dao.NewFavoriteDaoInstance().MaxId()
	if err != nil {
		return count, err
	}
	return count, nil
}


// 输入：视频id以及token
// 输出：点赞操作执行是否成功的报错，如若执行成功，则返回nil
func (s *FavoriteService) Add(videoId int64, token string) error {
	// 先查缓存 ..
	if _, exist := usersLoginInfo[token]; !exist {
		user, _ := dao.NewUserDaoInstance().QueryUserByToken(token)
		if user == nil {
			log.Printf("dao.NewUserDaoInstance().QueryUserByToken(token)方法失败，用户为空！")
			return utils.Error{Msg: "User doesn't exist, Please Register! "}
		}
		usersLoginInfo[token] = *pack.User(user)
	}
	// 点赞
	// 获取当前点赞的最后一个ID
	favoriteIdSequence, _ := favoriteService.LastId()
	// 上一步递增的点赞ID递增
	atomic.AddInt64(&favoriteIdSequence, 1)
	newFavorite := &dao.Favorite{
		Id:        favoriteIdSequence,
		UserToken: token,
		VideoId:   videoId,
		CreateAt:  time.Now(),
	}
	// 存新的点赞记录
	err := dao.NewFavoriteDaoInstance().Save(newFavorite)
	if err != nil {
		log.Printf("dao.NewFavoriteDaoInstance().Save(newFavorite)方法失败，保存点赞记录操作失误！")
		return err
	}
	return nil
}

// 输入：视频id和token
// 输出：取消点赞的操作的执行报错情况，当执行成功时，返回为nil
func (s *FavoriteService) Withdraw(videoId int64, token string) error {
	// 删除评论
	err := dao.NewFavoriteDaoInstance().Delete(videoId, token)
	if err != nil {
		log.Printf("dao.NewFavoriteDaoInstance().Delete(videoId, token)方法失败，取消点赞记录操作失误！")
		return err
	}
	return nil
}

// 服务端获取用户点赞列表操作的响应
func (s *FavoriteService) FavoriteList(token string) ([]*entity.Video, error) {
	return s.FindVideosByToken(token)
}
