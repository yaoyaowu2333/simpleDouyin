package dao

import (
	"gorm.io/gorm"
	"log"
	"simpleDouyin/entity"
	"sync"
)

type LoginStatus struct {
	Id     int64
	UserId int64
	Token  string
}

type LoginStatusDao struct {
}

var loginStatusDao *LoginStatusDao
var loginStatusOnce sync.Once
var usersLoginInfo = map[string]entity.User{}

func CopyULI() map[string]entity.User {
	return usersLoginInfo
}

// NewLoginStatusDaoInstance Singleton
func NewLoginStatusDaoInstance() *LoginStatusDao {
	loginStatusOnce.Do(
		func() {
			loginStatusDao = &LoginStatusDao{}
		})
	return loginStatusDao
}

func (*LoginStatusDao) QueryTokenByUserId(userId int64) (*LoginStatus, error) {
	var loginStatus *LoginStatus
	err := db.Where("user_id = ?", userId).First(&loginStatus).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Fatal("find loginStatus by user_id err:" + err.Error())
		return nil, err
	}
	return loginStatus, nil
}

// CreateLoginStatus
// 向登录日志表插入登录数据
func (*LoginStatusDao) CreateLoginStatus(loginStatus *LoginStatus) error {
	return db.Create(&loginStatus).Error
}

// QueryUserIdByToken
// 根据token获取登录用户的id
func (*LoginStatusDao) QueryUserIdByToken(token string) (int64, error) {
	var loginStatus *LoginStatus
	err := db.Where("token = ?", token).First(&loginStatus).Error
	if err == gorm.ErrRecordNotFound {
		return -1, nil
	}
	if err != nil {
		log.Fatal("find loginStatus by token err:" + err.Error())
		return -1, err
	}
	return loginStatus.UserId, nil
}
