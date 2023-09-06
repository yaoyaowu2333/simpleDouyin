package service

// TODO
import (
	"fmt"
	"log"
	"simpleDouyin/dao"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
	"sync"
	"sync/atomic"
	"time"
)

type CommentService struct {
}

var commentService *CommentService
var commentOnce sync.Once

func NewCommentServiceInstance() *CommentService {
	commentOnce.Do(
		func() {
			commentService = &CommentService{}
		})
	return commentService
}

// 加载视频所有评论
func (s *CommentService) LoadComments(videoId int64) ([]*entity.Comment, error) {
	return s.FindCommentByVideoId(videoId)
}

func (s *CommentService) FindCommentByName(name string) (*entity.Comment, error) {
	// 查询用户信息
	commentModel, err := dao.NewCommentDaoInstance().QueryCommentByName(name)
	if err != nil {
		log.Printf("dao.NewCommentDaoInstance().QueryCommentByName(name)方法执行有误！")
		return nil, err
	}

	// 包装用户信息
	return pack.Comment(commentModel), nil
}

// 获取视频所有评论
func (s *CommentService) FindCommentByVideoId(videoID int64) ([]*entity.Comment, error) {
	// invalid authorId
	if videoID <= 0 {
		log.Printf("视频id有问题！")
		return nil, nil
	}

	_, commentModels, err := dao.NewCommentDaoInstance().QueryCommentByVideoId(videoID)
	userNames := pack.UserNames(commentModels)
	fmt.Println(userNames)

	userModelMap, err := dao.NewUserDaoInstance().MQueryUserByName(userNames)
	if err != nil {
		log.Printf("dao.NewUserDaoInstance().MQueryUserByName(userNames)方法执行有误！")
		return nil, err
	}
	userMap := pack.MUserByName(userModelMap)
	comments := pack.Comments(commentModels)

	for i, comment := range comments {
		comment.User = userMap[userNames[i]]
	}

	return comments, nil
}

func (s *CommentService) TotalComment() (int64, error) {
	count, err := dao.NewCommentDaoInstance().Total()
	if err != nil {
		log.Printf("dao.NewCommentDaoInstance().Total()方法执行有误！")
		return -1, err
	}
	return count, nil
}

func (s *CommentService) LastId() (int64, error) {
	count, err := dao.NewCommentDaoInstance().MaxId()
	if err != nil {
		log.Printf("dao.NewCommentDaoInstance().MaxId()方法执行有误！")
		return count, err
	}
	return count, nil
}

// 添加评论的操作
func (s *CommentService) Add(videoId int64, token, text string) (*entity.Comment, error) {
	// 先查缓存,查看用户是否登录 ..
	var user *dao.User
	if token == "" {
		log.Printf("token为空")
		return nil, utils.Error{Msg: "User doesn't login"}
	}
	if _, exist := usersLoginInfo[token]; !exist {
		return nil, utils.Error{Msg: "User doesn't login"}
	}
	user, _ = dao.NewUserDaoInstance().QueryUserByToken(token)
	if user == nil {
		log.Printf("dao.NewUserDaoInstance().QueryUserByToken(token)方法失败，用户为空！")
		return nil, utils.Error{Msg: "User doesn't exist, Please Register! "}
	}
	// 评论
	log.Printf("下面获取当前最大的评论id，自加一，并生成新的评论，添加进数据库中...")
	commentIdSequence, _ := commentService.LastId()
	atomic.AddInt64(&commentIdSequence, 1)
	newComment := &dao.Comment{
		Id:       commentIdSequence,
		VideoId:  videoId,
		UserName: usersLoginInfo[token].Name,
		Content:  text,
		CreateAt: time.Now().Format("01-02"),
	}
	fmt.Println(newComment)
	comment, err := dao.NewCommentDaoInstance().Save(newComment)
	if err != nil {
		log.Printf("dao.NewCommentDaoInstance().Save(newComment)方法执行有误，添加评论操作失败！")
		return nil, err
	}
	// 修改该视频的评论数
	err = dao.NewVideoDaoInstance().UpdateCommentByID(videoId, 1)
	if err != nil {
		log.Printf("dao.NewVideoDaoInstance().UpdateFavoriteByID(videoId)方法失败，修改视频点赞数操作失误！")
		return nil, err
	}
	return pack.Comment(comment), nil
}

// 删除评论的操作
func (s *CommentService) Withdraw(id int64, token string, videoId int64) (*entity.Comment, error) {
	//
	log.Printf("开始删除评论！")
	// 先查缓存,查看用户是否登录 ..
	var user *dao.User
	if token == "" {
		log.Printf("token为空")
		return nil, utils.Error{Msg: "User doesn't login"}
	}
	if _, exist := usersLoginInfo[token]; !exist {
		return nil, utils.Error{Msg: "User doesn't login"}
	}
	user, _ = dao.NewUserDaoInstance().QueryUserByToken(token)
	if user == nil {
		log.Printf("dao.NewUserDaoInstance().QueryUserByToken(token)方法失败，用户为空！")
		return nil, utils.Error{Msg: "User doesn't exist, Please Register! "}
	}
	oldComment, err := dao.NewCommentDaoInstance().DeleteById(id)
	if err != nil {
		log.Printf("dao.NewCommentDaoInstance().DeleteById(id)方法执行有误，评论删除失败！")
		return nil, err
	}
	// 修改该视频的评论数
	err = dao.NewVideoDaoInstance().UpdateCommentByID(videoId, 2)
	if err != nil {
		log.Printf("dao.NewVideoDaoInstance().UpdateFavoriteByID(videoId)方法失败，修改视频点赞数操作失误！")
		return nil, err
	}
	return pack.Comment(oldComment), nil
}
