package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once     sync.Once
	instance *gorm.DB
	initErr  error
)

// GetGormDB 返回全局唯一的 gorm 连接。
//
// 之前每次调用都会新开一个 sqlite 连接，在并发执行规则时很容易触发
// "database is locked"，这里改为单例并开启 WAL。
func GetGormDB() (*gorm.DB, error) {
	once.Do(func() {
		instance, initErr = open()
	})
	return instance, initErr
}

// MustInit 在启动阶段初始化数据库并执行迁移，失败直接退出
func MustInit() *gorm.DB {
	gormDB, err := GetGormDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	if err := AutoMigrate(gormDB); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	return gormDB
}

func open() (*gorm.DB, error) {
	config, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", config.DatabasePath)
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	// sqlite 单写多读，限制连接数避免写冲突
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	return gormDB, nil
}

// AutoMigrate 建表 / 补齐新增字段
func AutoMigrate(gormDB *gorm.DB) error {
	return gormDB.AutoMigrate(
		&model.User{},
		&model.Timer{},
		&model.Rule{},
		&model.Account{},
		&model.ExecutionLog{},
	)
}
