package gateway

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

type webhookRequest struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

func token() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "gw_" + hex.EncodeToString(b), nil
}

func HandleList(c *gin.Context) {
	u, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}
	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "获取网关失败")
		return
	}
	var items []model.MessageGateway
	if err := dbConn.Where("uid = ?", u.ID).Order("id DESC").Find(&items).Error; err != nil {
		response.ServerError(c, "获取网关失败")
		return
	}
	response.OK(c, items)
}

func HandleCreate(c *gin.Context) {
	u, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}
	var req model.MessageGatewayRequest
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" {
		response.BadRequest(c, "网关名称和规则不能为空")
		return
	}
	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "创建网关失败")
		return
	}
	key, err := token()
	if err != nil {
		response.ServerError(c, "生成网关 Token 失败")
		return
	}
	item := model.MessageGateway{UID: u.ID, Name: strings.TrimSpace(req.Name), Token: key, CreateAt: time.Now().Unix()}
	if err := dbConn.Create(&item).Error; err != nil {
		response.ServerError(c, "创建网关失败")
		return
	}
	response.OK(c, item)
}

func HandleDelete(c *gin.Context) {
	u, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}
	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "删除网关失败")
		return
	}
	if err := dbConn.Where("id = ? AND uid = ?", c.Param("gatewayID"), u.ID).Delete(&model.MessageGateway{}).Error; err != nil {
		response.ServerError(c, "删除网关失败")
		return
	}
	response.OK(c, "网关已删除")
}

func HandleWebhook(c *gin.Context) {
	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "网关不可用")
		return
	}
	var gw model.MessageGateway
	if err := dbConn.Where("token = ?", c.Param("token")).First(&gw).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "网关不存在")
		return
	}
	var req webhookRequest
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Message) == "" {
		response.BadRequest(c, "message 不能为空")
		return
	}
	var rules []model.Rule
	if err := dbConn.Where("uid = ? AND gateway_id = ? AND enabled = ?", gw.UID, gw.ID, true).Find(&rules).Error; err != nil {
		response.ServerError(c, "读取网关规则失败")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()
	fails := make([]string, 0)
	for i := range rules {
		rule := rules[i]
		cfg, err := overrideMessage(rule.ConfigJson, req.Title, req.Message)
		if err != nil {
			fails = append(fails, rule.Name+": 规则不支持消息网关")
			continue
		}
		rule.ConfigJson = cfg
		if err := runner.ExecuteRuleWithContext(ctx, &rule, "webhook"); err != nil {
			fails = append(fails, rule.Name+": "+err.Error())
		}
	}
	response.OK(c, gin.H{"total": len(rules), "failed": fails})
}

func overrideMessage(raw, title, message string) (string, error) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return raw, err
	}
	matched := false
	if _, ok := obj["msg"]; ok {
		obj["msg"] = message
		matched = true
	}
	if _, ok := obj["message"]; ok {
		obj["message"] = message
		matched = true
	}
	if _, ok := obj["body"]; ok {
		obj["body"] = message
		matched = true
	}
	if !matched {
		return raw, fmt.Errorf("message field not found")
	}
	if title != "" {
		if _, ok := obj["title"]; ok {
			obj["title"] = title
		}
	}
	b, err := json.Marshal(obj)
	return string(b), err
}
