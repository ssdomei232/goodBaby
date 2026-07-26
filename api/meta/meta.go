// Package meta 向 WebUI 暴露驱动能力与站点信息
package meta

import (
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal/accountConfigChecker"
	"github.com/ssdomei232/goodBaby/internal/ruleConfigChecker"
)

// HandleGetProviders 返回所有账号类型与规则类型的元数据
//
// 前端据此动态渲染配置表单，新增驱动时前端无需改动。
func HandleGetProviders(c *gin.Context) {
	response.OK(c, gin.H{
		"accounts": accountConfigChecker.InitValidatorRegistry().Metas(),
		"rules":    ruleConfigChecker.InitValidatorRegistry().Metas(),
	})
}

// HandleGetSiteInfo 返回无需登录即可获取的站点信息
func HandleGetSiteInfo(c *gin.Context) {
	config, err := configs.GetConfig()
	if err != nil {
		response.ServerError(c, "获取配置失败")
		return
	}

	userCount, err := user.CountUsers()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	response.OK(c, gin.H{
		"name": "goodBaby",
		// 系统内还没有用户时始终允许注册，方便全新部署创建第一个账号
		"enable_registry":        config.EnableRegistry || userCount == 0,
		"need_initial_user":      userCount == 0,
		"check_interval_minutes": config.CheckIntervalMinutes,
	})
}
