package dao

// TODO

import (
	"gorm.io/gorm"
	"log"
	"regexp"
	"sync"
)

type User struct {
	Id            int64
	Name          string `gorm:"unique;not null"`
	Password      string
	FollowCount   int64
	FollowerCount int64
	VideoCount    int64
	LikeCount     int64
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

// NewUserDaoInstance Singleton
func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

// CreateUser
// 向用户表插入用户数据
func (*UserDao) CreateUser(user *User) error {
	return db.Create(&user).Error
}

// QueryUserById
// 根据用户id在用户表中查询用户
func (*UserDao) QueryUserById(id int64) (*User, error) {
	user := new(User) //实例化对象
	log.Printf("开始查询数据库中，id为%v的用户信息", id)
	result := db.Where("id = ?", id).First(&user)
	err := result.Error
	if err == gorm.ErrRecordNotFound {
		log.Printf("查询失败！")
		return nil, nil
	}
	if err != nil {
		log.Fatal("find user by id err:" + err.Error())
		return nil, err
	}
	return user, nil
}


// 作用：根据一组给定的用户ID查询用户信息
// 输入：用户ID切片，返回的是用户信息哈希表
// MQueryUserById will return empty array if no user is found
// 如果找不到用户，MQueryUserById将返回空数组
// MQueryUserById will return empty array if no user is found
// 依据id获取用户信息
func (*UserDao) MQueryUserById(ids []int64) (map[int64]User, error) {
	var users []*User
	err := db.Where("id in (?)", ids).Find(&users).Error
	if err != nil {
		log.Printf("查找给定用户id的用户信息，数据库查询失败！")
		return nil, err
	}
	var userMap = make(map[int64]User, len(users))
	for _, user := range users {
		id := user.Id
		userMap[id] = *user
	}
	return userMap, nil
}

func (d *UserDao) MQueryUserByName(names []string) (map[string]User, error) {
	var users []*User
	err := db.Where("name in (?)", names).Find(&users).Error
	if err != nil {
		return nil, err
	}
	var userMap = make(map[string]User, len(users))
	for _, user := range users {
		userMap[user.Name] = *user
	}
	return userMap, nil
}

// QueryUserByName
// 通过用户名从用户表中查询
func (*UserDao) QueryUserByName(name string) (*User, error) {
	var user *User
	err := db.Where("name = ?", name).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 根据用户令牌（token）查询用户信息
// 输入：token
// 输出：返回值是一个 *User 指针，表示查询到的用户对象，以及一个 error，表示可能的错误

// QueryUserByToken
// 从token中获取用户名与密码
// 效验用户名与密码
func (*UserDao) QueryUserByToken(token string) (*User, error) {
	var users *User //实例化对象
	// 创建了一个正则表达式对象 re， 并从re中提取用户名和密码
	re, err := regexp.Compile("[A-Za-z0-9_@.\\-\u4e00-\u9fa5]+")
	if err != nil {
		log.Printf("用户名、密码提取失败！")
		return nil, err
	}
	name := re.FindAllString(token, 2)[0]
	password := re.FindAllString(token, 2)[1]
	// 根据用户名和密码在数据库中查询匹配的用户
	err = db.Debug().Where("name = ? and password = ?", name, password).First(&users).Error
	// 没查到的情况
	if err == gorm.ErrRecordNotFound {
		log.Printf("数据库查询失败！没有查到与提供的用户名及密码匹配的用户")
		return nil, err
	}
	if err != nil {
		//fmt.Println("record not found!")
		return nil, err
	}
	// 查到了， 返回
	return users, nil
}


func (*UserDao) Save(user *User) error {
	result := db.Create(&user)
	err := result.Error
	if err != nil {
		return err
	}
	return nil
}

func (*UserDao) Total() (int64, error) {
	// 获取全部记录
	var count int64
	result := db.Table("users").Count(&count)
	err := result.Error
	if err != nil {
		log.Fatal("total user err:" + err.Error())
		return -1, err
	}
	return count, nil
}

func (*UserDao) MaxId() (int64, error) {
	// 获取全部记录
	var lastRec *User
	result := db.Table("users").Last(&lastRec)
	err := result.Error
	if err != nil {
		return 0, err
	}
	return lastRec.Id, nil
}

// IncreaseVideoCountByOne
// 使得发布视频数加一
func (*UserDao) IncreaseVideoCountByOne(id int64) error {
	var user *User
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Printf("数据库查询指定id的用户，查询失败！")
		return err
	}
	user.VideoCount = user.VideoCount + 1
	return db.Save(&user).Error
}
