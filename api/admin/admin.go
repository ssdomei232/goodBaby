// Package admin 提供管理员可用的系统配置接口
package admin

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/handler/scheduler"
)

// HandleGetConfig 返回当前可编辑的系统配置，以及一些只读的运行时信息
func HandleGetConfig(c *gin.Context) {
	config, err := configs.GetConfig()
	if err != nil {
		response.ServerError(c, "读取配置失败")
		return
	}

	response.OK(c, gin.H{
		"config": config.Editable(),
		// 只读信息，让管理员知道当前跑在什么环境上，但不能在这里改
		"readonly": gin.H{
			"listen_addr":     config.ListenAddr,
			"database_driver": config.DatabaseDriver,
		},
	})
}

// HandleUpdateConfig 更新可编辑的系统配置
func HandleUpdateConfig(c *gin.Context) {
	var req configs.EditableConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	before, err := configs.GetConfig()
	if err != nil {
		response.ServerError(c, "读取配置失败")
		return
	}

	updated, err := configs.ApplyEditable(req)
	if err != nil {
		// Validate 的错误可以直接展示给用户
		response.BadRequest(c, err.Error())
		return
	}

	// 检查间隔变了就重排定时任务，免得改完还要重启
	if before.CheckIntervalMinutes != updated.CheckIntervalMinutes {
		if err := scheduler.Reschedule(updated.CheckIntervalMinutes); err != nil {
			log.Printf("重排定时任务失败: %v", err)
			response.ServerError(c, "配置已保存，但定时任务重排失败，请重启服务")
			return
		}
		log.Printf("检查间隔已更新为 %d 分钟", updated.CheckIntervalMinutes)
	}

	response.OK(c, updated.Editable())
}
