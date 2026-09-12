// Package gateway 提供消息网关的 HTTP 接口。
//
// 数据库读写都发生在这一层，真正的投递逻辑在 internal/gateway 里，
// 这里只负责鉴权、取数、组装投递任务。
package gateway

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	gatewaycore "github.com/ssdomei232/goodBaby/internal/gateway"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
	"gorm.io/gorm"
)

// webhookRequest 外部系统投递消息的请求体
type webhookRequest struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

// HandleList 获取当前用户的所有消息网关
func HandleList(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}

	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "获取网关失败")
		return
	}

	items := []model.MessageGateway{}
	if err := dbConn.Where("uid = ?", userInfo.ID).Order("id DESC").Find(&items).Error; err != nil {
		response.ServerError(c, "获取网关失败")
		return
	}
	response.OK(c, items)
}

// HandleCreate 创建消息网关
func HandleCreate(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.MessageGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if err := req.Validate(); err != nil {
		response.FromError(c, err, "创建网关失败")
		return
	}

	// 不填类型时用默认网关，填了就必须是已注册的类型
	if req.Type == "" {
		req.Type = model.GatewayTypeWebhook
	}
	if _, ok := gatewaycore.InitGatewayRegistry().Resolve(req.Type); !ok {
		response.BadRequest(c, "不支持的消息网关类型: "+req.Type)
		return
	}

	token, err := gatewaycore.NewToken()
	if err != nil {
		response.ServerError(c, "生成网关 Token 失败")
		return
	}

	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "创建网关失败")
		return
	}

	item := model.MessageGateway{
		UID:      userInfo.ID,
		Name:     req.Name,
		Type:     req.Type,
		Token:    token,
		CreateAt: time.Now().Unix(),
	}
	if err := dbConn.Create(&item).Error; err != nil {
		response.ServerError(c, "创建网关失败")
		return
	}

	response.OK(c, item)
}

// HandleDelete 删除消息网关，绑定在它下面的规则一并删除
func HandleDelete(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "未登录")
		return
	}

	gatewayID, err := parseID(c.Param("gatewayID"))
	if err != nil {
		response.BadRequest(c, "网关 ID 格式错误")
		return
	}

	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "删除网关失败")
		return
	}

	err = dbConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("gateway_id = ? AND uid = ?", gatewayID, userInfo.ID).
			Delete(&model.GatewayRule{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND uid = ?", gatewayID, userInfo.ID).
			Delete(&model.MessageGateway{}).Error
	})
	if err != nil {
		response.ServerError(c, "删除网关失败")
		return
	}

	response.OK(c, "网关已删除")
}

// HandleWebhook 接收外部系统投递的消息，触发绑定在该网关上的规则
func HandleWebhook(c *gin.Context) {
	var req webhookRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		response.BadRequest(c, "message 不能为空")
		return
	}

	dbConn, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "网关不可用")
		return
	}

	var target model.MessageGateway
	if err := dbConn.Where("token = ?", c.Param("token")).First(&target).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "网关不存在")
		return
	}

	rules := []model.GatewayRule{}
	if err := dbConn.Where("uid = ? AND gateway_id = ? AND enabled = ?", target.UID, target.ID, true).
		Find(&rules).Error; err != nil {
		response.ServerError(c, "读取网关规则失败")
		return
	}

	// 投递需要在 HTTP 请求内返回结果，因此用较短的超时
	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	result, err := gatewaycore.InitGatewayRegistry().Deliver(ctx, target.Type, &gatewaycore.Task{
		Gateway: &target,
		Rules:   rules,
		Message: gatewaycore.Message{
			Title:   strings.TrimSpace(req.Title),
			Content: strings.TrimSpace(req.Message),
		},
	})
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.OK(c, result)
}
