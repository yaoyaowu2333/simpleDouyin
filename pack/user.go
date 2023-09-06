package pack

import (
	"simpleDouyin/dao"
	"simpleDouyin/entity"
	"strconv"
)

// User if param is nil then return nil
// 将dao.user数据通过User进行处理，在拷贝的过程中对数据进行组装
func User(userModel *dao.User) *entity.User {
	if userModel != nil {
		return &entity.User{
			Id:             userModel.Id,
			Name:           userModel.Name,
			FollowCount:    userModel.FollowCount,
			FollowerCount:  userModel.FollowerCount,
			TotalFavorited: strconv.FormatInt(userModel.TotalFavorited, 10),
			WorkCount:      userModel.WorkCount,
			FavoriteCount:  userModel.FavoriteCount,
		}
	}
	return nil
}

func Users(userModels []*dao.User) []*entity.User {
	if userModels != nil {
		var users = make([]*entity.User, 0, len(userModels))
		for _, model := range userModels {
			users = append(users, User(model))
		}
		return users
	}
	return nil
}

// MUser if param is nil then return empty map
// 将userModels数据通过MUser进行处理，在拷贝的过程中对数据进行组装
func MUser(userModels map[int64]dao.User) map[int64]*entity.User {
	if userModels != nil {
		var users = make(map[int64]*entity.User, len(userModels))
		for id, userModel := range userModels {
			users[id] = User(&userModel)
		}
		return users
	}
	return nil
}

func MUserByName(userModels map[string]dao.User) map[string]entity.User {
	if userModels != nil {
		var users = make(map[string]entity.User, len(userModels))
		for name, userModel := range userModels {
			users[name] = *User(&userModel)
		}
		return users
	}
	return nil
}
