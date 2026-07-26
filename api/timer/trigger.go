package timer

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleTriggerTimer 手动触发一个 Timer，用于调试。
//
// 会真实执行该 Timer 下所有启用的规则，但不会改变 Timer 的签到/触发状态；
// 使用较短的测试超时，避免在页面上等待数小时的重试。
func HandleTriggerTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	if _, err := findTimer(timerID, userInfo.ID); err != nil {
		response.NotFound(c, "Timer 不存在")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	var rules []model.Rule
	if err := gormDB.Where("timer_id = ? AND uid = ? AND enabled = ?", timerID, userInfo.ID, true).
		Find(&rules).Error; err != nil {
		response.ServerError(c, "获取关联规则失败")
		return
	}

	if len(rules) == 0 {
		response.BadRequest(c, "该定时器下没有启用的规则")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	// 并发执行所有规则，收集失败信息
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		fails []string
	)
	for i := range rules {
		rule := rules[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runner.ExecuteRuleWithContext(ctx, &rule, model.TriggerManual); err != nil {
				mu.Lock()
				fails = append(fails, rule.Name+": "+err.Error())
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	response.OK(c, gin.H{
		"total":  len(rules),
		"failed": fails,
	})
}
