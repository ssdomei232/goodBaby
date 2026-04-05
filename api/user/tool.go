package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	"golang.org/x/crypto/bcrypt"
)

// 加密密码
func encryptPassword(password string) (string, error) {
	hashedID, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedID), nil
}

// 验证密码
func verifyPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func createUser(u *model.User) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	hashedPassword, err := encryptPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	result := gormDB.Create(u)
	return result.Error
}

// 校验密码
func verifyUser(u *model.User) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	var hashedPassword string
	result := gormDB.Where("username = ?", u.Username).First(&model.User{}).Select("password").Scan(&hashedPassword)
	if result.Error != nil {
		return result.Error
	}

	return verifyPassword(hashedPassword, u.Password)
}

// 通过用户名获取用户信息
func GetUserInfo(username string) (*model.User, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var user model.User
	result := gormDB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 通过 gin.Context 获取用户信息
func GetUserInfoByGinCtx(c *gin.Context) (userInfo *model.User, err error) {
	session := sessions.Default(c)
	username := session.Get("username")
	userInfo, err = GetUserInfo(username.(string))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
