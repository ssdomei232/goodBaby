// internal/timer_manager.go
package internal

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TimerManager struct {
	duration  time.Duration
	resetTime time.Time
	mutex     sync.RWMutex
}

var GlobalTimerManager *TimerManager

func InitTimerManager(duration time.Duration) {
	GlobalTimerManager = &TimerManager{
		duration:  duration,
		resetTime: time.Now(),
	}
}

func (tm *TimerManager) Reset() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.resetTime = time.Now()
}

func (tm *TimerManager) GetRemainingTime() time.Duration {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	elapsed := time.Since(tm.resetTime)
	remaining := tm.duration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

func HandleTimerStatus(c *gin.Context) {
	remainingTime := GlobalTimerManager.GetRemainingTime()

	c.JSON(http.StatusOK, gin.H{
		"code":           200,
		"remaining_time": remainingTime.Seconds(), // 返回剩余秒数
		"message":        "ok",
	})
}
