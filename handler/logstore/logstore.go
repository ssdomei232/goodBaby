// Package logstore 负责写入与清理规则执行日志
package logstore

import (
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// Record 写入一条执行日志，失败只记录到标准输出，不影响主流程
func Record(entry *model.ExecutionLog) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("写入执行日志失败: %v", err)
		return
	}

	if entry.CreateAt == 0 {
		entry.CreateAt = time.Now().Unix()
	}

	if err := gormDB.Create(entry).Error; err != nil {
		log.Printf("写入执行日志失败: %v", err)
		return
	}

	trim(entry.UID)
}

// trim 按配置保留每个用户最近的 N 条日志
func trim(uid uint) {
	config, err := configs.GetConfig()
	if err != nil || config.LogRetainCount <= 0 {
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		return
	}

	// 用 Find 而不是 First：日志没超出保留条数是常态，
	// First 找不到记录会返回 ErrRecordNotFound，每写一条日志都刷一行 gorm 警告。
	var cutoff []model.ExecutionLog
	if err := gormDB.Where("uid = ?", uid).
		Order("id DESC").
		Offset(config.LogRetainCount).
		Limit(1).
		Find(&cutoff).Error; err != nil {
		log.Printf("清理执行日志失败: %v", err)
		return
	}
	if len(cutoff) == 0 {
		// 没有超出保留条数
		return
	}

	if err := gormDB.Where("uid = ? AND id <= ?", uid, cutoff[0].ID).
		Delete(&model.ExecutionLog{}).Error; err != nil {
		log.Printf("清理执行日志失败: %v", err)
	}
}
