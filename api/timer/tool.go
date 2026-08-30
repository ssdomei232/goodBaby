package timer

import (
	"strconv"
	"time"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func boolOr(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

// findTimer 查找属于该用户的 Timer
func findTimer(timerID, uid uint) (*model.Timer, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var timer model.Timer
	if err := gormDB.Where("id = ? AND uid = ?", timerID, uid).First(&timer).Error; err != nil {
		return nil, err
	}
	return &timer, nil
}

func countRules(timerID, uid uint) (int64, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var count int64
	err = gormDB.Model(&model.Rule{}).Where("timer_id = ? AND uid = ?", timerID, uid).Count(&count).Error
	return count, err
}

// signTimer 完成一次签到：重置计时、清除提醒与触发标记
func signTimer(timer *model.Timer) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if err := gormDB.Model(&model.Timer{}).Where("id = ?", timer.ID).
		Updates(map[string]any{
			"last_sign":   now,
			"last_remind": 0,
			"triggered":   false,
		}).Error; err != nil {
		return err
	}

	timer.LastSign = now
	timer.LastRemind = 0
	timer.Triggered = false
	return nil
}

// getTimerOwnerUID 获取 Timer 的所属用户 ID
func getTimerOwnerUID(timerID uint) (uint, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var timer model.Timer
	if err := gormDB.Select("uid").Where("id = ?", timerID).First(&timer).Error; err != nil {
		return 0, err
	}
	return timer.UID, nil
}
