package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrUserNotFound 用户不存在
var ErrUserNotFound = errors.New("用户不存在")

// ErrWrongPassword 密码错误
var ErrWrongPassword = errors.New("用户名或密码错误")

// contextKey AuthMiddleware 把当前用户放进 gin.Context 时使用的 key
const contextKey = "currentUser"

// 加密密码
func encryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
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
	if u.APIKey == "" {
		u.APIKey, err = generateAPIKey()
		if err != nil {
			return err
		}
	}
	u.CreateAt = time.Now().Unix()

	return gormDB.Create(u).Error
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "gb_" + hex.EncodeToString(b), nil
}

func ensureAPIKey(user *model.User) error {
	if user.APIKey != "" {
		return nil
	}
	key, err := generateAPIKey()
	if err != nil {
		return err
	}
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}
	if err := gormDB.Model(&model.User{}).Where("id = ? AND (api_key IS NULL OR api_key = '')", user.ID).Update("api_key", key).Error; err != nil {
		return err
	}
	user.APIKey = key
	return nil
}

// IsUsernameTaken 用户名是否已被占用
func IsUsernameTaken(username string) (bool, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return false, err
	}

	var count int64
	if err := gormDB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountUsers 返回已注册用户数，用于首个用户的引导注册
func CountUsers() (int64, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var count int64
	err = gormDB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// verifyUser 校验用户名与密码，成功时返回用户记录
//
// 原实现把 First 与 Select().Scan() 串在同一条语句上，实际拿不到密码哈希，
// 且用 sql.ErrNoRows 去比较 gorm 的错误永远不成立。
func verifyUser(username, password string) (*model.User, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := gormDB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := verifyPassword(user.Password, password); err != nil {
		return nil, ErrWrongPassword
	}
	return &user, nil
}

// GetUserInfo 通过用户名获取用户信息
func GetUserInfo(username string) (*model.User, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := gormDB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID 通过 ID 获取用户信息
func GetUserByID(id uint) (*model.User, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := gormDB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if err := ensureAPIKey(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserInfoByGinCtx 获取当前登录用户
//
// AuthMiddleware 已经查过一次库并写入了 context，这里优先读缓存。
func GetUserInfoByGinCtx(c *gin.Context) (*model.User, error) {
	if cached, ok := c.Get(contextKey); ok {
		if user, ok := cached.(*model.User); ok {
			return user, nil
		}
	}

	session := sessions.Default(c)
	uid, ok := session.Get("uid").(uint)
	if !ok {
		return nil, ErrUserNotFound
	}

	user, err := GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	c.Set(contextKey, user)
	return user, nil
}

// setSession 登录成功后写入会话
func setSession(c *gin.Context, user *model.User) error {
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Set("username", user.Username)
	return session.Save()
}

// validateNotifyConfig 校验钉钉提醒配置
func validateNotifyConfig(raw *string) error {
	if raw == nil || *raw == "" {
		return nil
	}

	var config struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal([]byte(*raw), &config); err != nil {
		return fmt.Errorf("钉钉配置 JSON 格式错误: %v", err)
	}
	if config.AccessToken == "" {
		return fmt.Errorf("钉钉配置中 access_token 不能为空")
	}
	return nil
}
