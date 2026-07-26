package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/model"
	"gorm.io/driver/postgres"
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
	if err := ensureAdmin(gormDB); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}
	return gormDB
}

// ensureAdmin 保证系统里至少有一个管理员。
//
// IsAdmin 是后加的字段，老数据库升级上来时所有人都不是管理员，
// 这里把最早注册的那个用户提升为管理员。
func ensureAdmin(gormDB *gorm.DB) error {
	var adminCount int64
	if err := gormDB.Model(&model.User{}).Where("is_admin = ?", true).Count(&adminCount).Error; err != nil {
		return err
	}
	if adminCount > 0 {
		return nil
	}

	var first model.User
	if err := gormDB.Order("id ASC").First(&first).Error; err != nil {
		// 一个用户都没有，等第一个注册的人来当管理员
		return nil
	}

	if err := gormDB.Model(&model.User{}).Where("id = ?", first.ID).
		Update("is_admin", true).Error; err != nil {
		return err
	}
	log.Printf("已将用户 %s (ID: %d) 提升为管理员", first.Username, first.ID)
	return nil
}

func open() (*gorm.DB, error) {
	config, err := configs.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	var (
		gormDB *gorm.DB
		driver = config.DatabaseDriver
	)

	switch driver {
	case configs.DriverPostgres:
		if config.DatabaseDSN == "" {
			return nil, fmt.Errorf("使用 postgres 时必须配置 database_dsn")
		}
		gormDB, err = gorm.Open(postgres.Open(config.DatabaseDSN), gormConfig)
	case configs.DriverSQLite, "":
		driver = configs.DriverSQLite
		dsn := fmt.Sprintf(
			"file:%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)",
			config.DatabasePath,
		)
		gormDB, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s（可选 sqlite / postgres）", config.DatabaseDriver)
	}
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	if driver == configs.DriverSQLite {
		// sqlite 单写多读，限制连接数避免写冲突
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	} else {
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// 提前握手，配置写错时在启动阶段就报错，而不是等到第一次查询
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	log.Printf("数据库已连接 (driver=%s)", driver)
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
