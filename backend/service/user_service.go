package service

import (
	"errors"
	"ipicture-backend/models"
	"ipicture-backend/utils"
	"ipicture-backend/config"
)

func RegisterUser(req models.UserRegisterRequest) error {
	db := config.DB
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err == nil {
		return errors.New("用户名已存在")
	}
	user = models.User{
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
	}
	return db.Create(&user).Error
}

func LoginUser(req models.UserLoginRequest) (string, error) {
	db := config.DB
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", errors.New("用户不存在")
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", errors.New("密码错误")
	}
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "", errors.New("生成token失败")
	}
	return token, nil
} 