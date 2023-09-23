package service

import (
	"Gous/internal/cache"
	"Gous/internal/dao"
	"Gous/internal/model"
	"Gous/internal/utils"
	"Gous/pkg/constant"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Register 用户注册
// 真正操作数据库
func Register(req *RegisterRequest) error {
	//校验参数合法性
	if req.UserName == "" || req.PassWord == "" || req.Age <= 0 ||
		!utils.Contains([]string{constant.GenderMale, constant.GenderFeMale}, req.Gender) {
		log.Errorf("Gous：register param invalid")
		return fmt.Errorf("register param invalid")
	}

	// 数据库操作
	existedUser, err := dao.GetUserByName(req.UserName)
	// 查询出错
	if err != nil {
		log.Errorf("Gous: Register | error: %v", err)
		return fmt.Errorf("gous: Register | error: %v", err)
	}

	// 数据库已存在该用户
	if existedUser != nil {
		log.Errorf("Gous: 用户已经注册，user_name=%s", req.UserName)
		return fmt.Errorf("gous: 用户已经注册，user_name=%s", req.UserName)
	}

	// 往数据库中插入该记录
	user := &model.User{
		CreateModel: model.CreateModel{Creator: req.UserName},
		ModifyModel: model.ModifyModel{Modifier: req.UserName},
		Name:        req.UserName,
		Gender:      req.Gender,
		NickName:    req.NickName,
		Age:         req.Age,
		PassWord:    req.PassWord,
	}
	log.Infof("Gous：user ====== %+v", user)
	if err := dao.CreateUser(user); err != nil {
		log.Errorf("Gous：Register failed | error: %v", err)
		return fmt.Errorf("gous：register failed | error: %v", err)
	}
	return nil
}

// Login 查询是否存在该用户，并创建一个会话 session
func Login(ctx context.Context, req *LoginRequest) (string, error) {
	uuid := ctx.Value(constant.ReqUuid)
	log.Debugf("%s | Login access from:%s,@,%s", uuid, req.UserName, req.PassWord)

	// 获取数据库中该用户信息
	user, err := getUserInfo(req.UserName)

	// 查询失败或没有查到，就返回
	if err != nil {
		log.Errorf("Login | %v", err)
		return "", fmt.Errorf("login|%v", err)
	}

	// 密码不正确
	if req.PassWord != req.PassWord {
		log.Errorf("Login|password err : req.password=%s|user.password=%s", req.PassWord, user.PassWord)
		return "", fmt.Errorf("password is not correct")
	}

	// 创建会话 ID session
	session := utils.GenerateSession(user.Name)
	// 缓存 session
	err = cache.SetSessionInfo(user, session)

	if err != nil {
		log.Errorf(" Login|Failed to SetSessionInfo, uuid=%s|user_name=%s|session=%s|err=%v", uuid, user.Name, session, err)
		return "", fmt.Errorf("login|SetSessionInfo fail:%v, err")
	}

	log.Infof("Login successfully, %s@%s with redis_session session_%s", req.UserName, req.PassWord, session)
	return session, nil
}

func getUserInfo(userName string) (*model.User, error) {
	// 查询 redis 缓存
	user, err := cache.GetUserInfoFromCache(userName)
	if err == nil && user.Name == userName {
		log.Infof("cache_user ===== %v", user)
		return user, nil
	}
	// 查询数据库
	user, err = dao.GetUserByName(userName)
	if err != nil {
		return user, err
	}

	if user == nil {
		return nil, fmt.Errorf("用户尚未注册")
	}
	log.Infof("user === %+v", user)

	// 缓存 用户信息
	err = cache.SetUserCacheInfo(user)
	if err != nil {
		log.Error("cache userinfo failed for user:", user.Name, " with err:", err.Error())
	}
	log.Infof("setUserInfo successfully, with key userinfo_%s", user.Name)
	return user, nil
}

func Logout(ctx context.Context, req *LogoutRequest) error {
	// 获取 uuid
	uuid := ctx.Value(constant.ReqUuid)
	// 获取上下文中传递的 session
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|Logout access from,user_name=%s|session=%s", uuid, req.UserName, session)
	// 查询 redis 中是否存在该 session，判断是否处于登录状态
	_, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("Logout|GetSessionInfo err:%v", err)
	}

	// 从 redis 中删除，并返回错误信息
	err = cache.DelSessionInfo(session)
	if err != nil { // 删除失败
		log.Errorf("%s|Failed to delSessionInfo :%s", uuid, session)
		return fmt.Errorf("del session err:%v", err)
	}
	// 删除成功
	log.Infof("%s|Success to delSessionInfo :%s", uuid, session)
	return nil
}

// Logoff 注销
func Logoff(req *LogoffRequest, ctx context.Context) error {
	existedUser, err := dao.GetUserByName(req.UserName)
	// 查询出错
	if err != nil {
		log.Errorf("Logoff|%v", err)
		return fmt.Errorf("logoff|%v", err)
	}

	// 数据库中查不到
	if existedUser == nil {
		log.Errorf("UserNotExistedUser%v", err)
		return fmt.Errorf("userNotexisteduser%v", err)
	}

	// 删除 session 会话
	session := ctx.Value(constant.SessionKey).(string)
	err = cache.DelSessionInfo(session)
	if err != nil {
		log.Errorf("|Failed to delSessionInfo :%s", session)
		return fmt.Errorf("del delsessioninfo err:%v", err)
	}

	// 清空 redis 缓存中的用户信息
	if err := cache.DelUserCacheInfo(existedUser); err != nil {
		log.Errorf("DeleteCacheInfo|%v", err)
		return fmt.Errorf("deletecacheinfo|%v", err)
	}

	// 删除数据库信息
	log.Infof("user ====== %+v", existedUser)
	if err := dao.DeleteUser(existedUser); err != nil {
		log.Errorf("DeleteDB|%v", err)
		return fmt.Errorf("deletedb|%v", err)
	}
	return nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	fmt.Println(ctx)
	uuid := ctx.Value(constant.ReqUuid)
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|GetUserInfo access from,user_name=%s|session=%s", uuid, req.UserName, session)

	if session == "" || req.UserName == "" {
		return nil, fmt.Errorf("GetUserInfo|request params invalid")
	}

	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return nil, fmt.Errorf("getUserInfo|GetSessionInfo err:%v", err)
	}

	if user.Name != req.UserName {
		log.Errorf("%s|session info not match with username=%s", uuid, req.UserName)
	}
	log.Infof("%s|Succ to GetUserInfo|user_name=%s|session=%s", uuid, req.UserName, session)
	return &GetUserInfoResponse{
		UserName: user.Name,
		Age:      user.Age,
		Gender:   user.Gender,
		PassWord: user.PassWord,
		NickName: user.NickName,
	}, nil
}

// UpdateUserNickName 更新用户昵称
func UpdateUserNickName(ctx context.Context, req *UpdateNickNameRequest) error {
	uuid := ctx.Value(constant.ReqUuid)
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|UpdateUserNickName access from,user_name=%s|session=%s", uuid, req.UserName, session)
	log.Infof("UpdateUserNickName|req==%v", req)

	if session == "" || req.UserName == "" {
		return fmt.Errorf("UpdateUserNickName|request params invalid")
	}

	//从缓存中获取用户信息
	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("UpdateUserNickName|GetSessionInfo err:%v", err)
	}

	if user.Name != req.UserName {
		log.Errorf("UpdateUserNickName|%s|session info not match with username=%s", uuid, req.UserName)
	}

	updateUser := &model.User{
		NickName: req.NewNickName,
	}

	return updateUserInfo(updateUser, req.UserName, session)
}

// 更新数据库中用户昵称
func updateUserInfo(user *model.User, userName, session string) error {
	//更新数据库中的昵称
	affectedRows := dao.UpdateUserInfo(userName, user)

	// db更新成功
	if affectedRows == 1 {
		user, err := dao.GetUserByName(userName)
		if err == nil {
			cache.UpdateCachedUserInfo(user)
			if session != "" {
				err = cache.SetSessionInfo(user, session)
				if err != nil {
					log.Error("update session failed:", err.Error())
					cache.DelSessionInfo(session)
				}
			}
		} else {
			log.Error("Failed to get dbUserInfo for cache, username=%s with err:", userName, err.Error())
		}
	}
	return nil
}
