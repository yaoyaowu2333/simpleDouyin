package service

// TODO
import (
	"errors"
	"log"
	"simpleDouyin/dao"
	"simpleDouyin/entity"
	"simpleDouyin/pack"
	"simpleDouyin/utils"
	"strings"
	"sync"
)

type UserService struct {
}

var userService *UserService
var userOnce sync.Once
var usersLoginInfo = dao.CopyULI()

func NewUserServiceInstance() *UserService {
	userOnce.Do(
		func() {
			userService = &UserService{}
		})
	return userService
}

// UserInfo
// 查询用户信息,重组用户信息并返回
func (s *UserService) UserInfo(uid int64) (*entity.User, error) {
	// 查询用户信息
	userModel, err := dao.NewUserDaoInstance().QueryUserById(uid)
	if err != nil {
		log.Printf("方法userDao.QueryUserById(uid)失败,%v\n", err)
		return nil, err
	}
	log.Printf("方法userDao.QueryUserById(uid)成功,%v\n", userModel)
	// 包装用户信息
	user := pack.User(userModel)
	user.IsFollow = true
	return user, nil
}

// FindUserByName
// 通过用户名查询用户,并通过pack.User(user)对查询到的数据进行组装
func (s *UserService) FindUserByName(name string) (*entity.User, error) {
	user, err := dao.NewUserDaoInstance().QueryUserByName(name)
	if err != nil || user == nil {
		log.Printf("方法userDao.QueryUserByName(name)失败,%v\n", err)
		return nil, err
	}
	log.Printf("方法userDao.QueryUserByName(name)成功\n")
	return pack.User(user), nil
}

// AddUser 创建用户和token
// 对密码进行加密存储
// 记录登录日志
func (s *UserService) AddUser(username, password string) error {
	// 用户注册
	password = utils.Md5(password)
	token := "<" + username + "><" + password + ">"
	newUser := &dao.User{
		Name:     username,
		Password: password,
	}
	//将用户信息插入用户表
	err := dao.NewUserDaoInstance().CreateUser(newUser)
	if err != nil {
		return err
	}

	// 创建登录日志
	loginStatus := &dao.LoginStatus{
		UserId: newUser.Id,
		Token:  token,
	}
	//将新用户存入登录缓冲中
	usersLoginInfo[loginStatus.Token] = *pack.User(newUser)
	//记录登录信息到登录日志表中
	err = dao.NewLoginStatusDaoInstance().CreateLoginStatus(loginStatus)
	if err != nil {
		log.Printf("方法LoginStatusDao.CreateLoginStatus(loginStatus)失败,%v\n", err)
		return err
	}
	log.Printf("方法LoginStatusDao.CreateLoginStatus(loginStatus)成功\n")
	return nil
}

// Register
// 对新注册的用户名与密码进行效验
// 生成token,并检查是否该用户名是否存在,若没有则成功注册
func (s *UserService) Register(username, password string) error {
	// 用户输入验证
	err := InfoVerify(username, password)
	if err != nil {
		log.Printf("效验失败：%v\n", err)
		return err
	}
	log.Printf("方法InfoVerify(username, password)效验成功\n")
	token := "<" + username + "><" + password + ">"
	// 先查缓存 ..
	if _, exist := usersLoginInfo[token]; !exist {
		//通过用户名查询是否曾经注册过
		if user, _ := userService.FindUserByName(username); user == nil {
			//将用户添加到user表
			err = s.AddUser(username, password)
			if err != nil {
				log.Printf("用户添加失败\n")
				return utils.Error{Msg: "User register failed, Please retry for a minute!\n"}
			}
			log.Printf("用户注册成功\n")
			return err
		}
		log.Printf("用户已经存在,不需要注册\n")
	}
	log.Printf("用户已经登录,不需要注册\n")
	return utils.Error{Msg: "User already exist, don't register again!"}
}

func (s *UserService) Login(username, password string) (*int64, *string, error) {
	// 用户校验
	password = utils.Md5(password)
	token := "<" + username + "><" + password + ">"

	user, _ := s.FindUserByName(username)
	if user == nil {
		return nil, nil, utils.Error{Msg: "User doesn't exist, Please Register! "}
	}
	usersLoginInfo[token] = *user
	// 密码校验
	result, _ := dao.NewUserDaoInstance().QueryUserByToken(token)
	if result == nil {
		return nil, nil, utils.Error{Msg: "Password Wrong!"}
	}
	// 创建token
	loginStatus := &dao.LoginStatus{
		UserId: user.Id,
		Token:  token,
		//Token:  utils.GenerateUUID(),
	}
	err := dao.NewLoginStatusDaoInstance().CreateLoginStatus(loginStatus)
	if err != nil {
		return nil, nil, err
	}
	return &user.Id, &token, nil
}

// InfoVerify
// 对用户名与密码分别进行验证
func InfoVerify(username string, password string) error {
	if Check(username) {
		return errors.New("Please Check Username!\nThe length is controlled within 4-32 characters, and <, >, \\is not allowed")
	}
	if Check(password) {
		return errors.New("Please Check Password!\nThe length is controlled within 4-32 characters, and <, >, \\is not allowed")
	}
	return nil
}

// Check
// 验证信息,要求如下
// 要求1:长度大于等于4且小于等于32的
// 要求2:不包含<>, /与\\
func Check(str string) bool {
	length := len(str)
	if length < 4 || length > 32 {
		return true
	}
	if strings.Contains(str, "<") || strings.Contains(str, ">") ||
		strings.Contains(str, "/") || strings.Contains(str, "\\") {
		return true
	}
	return false
}
