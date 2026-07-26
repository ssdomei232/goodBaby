package configs

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"sync"
)

// 默认配置文件路径，可通过环境变量 GOODBABY_CONFIG 覆盖
const defaultConfigPath = "config.json"

type Config struct {
	// HTTP 监听地址，如 ":8088"
	ListenAddr string `json:"listen_addr"`
	// 是否开放注册
	EnableRegistry bool `json:"enable_registry"`
	// 规则执行的最长重试时间(小时)，用于指数退避的整体超时
	TimeoutDurationHours int `json:"timeout_duration_hours"`
	// 检查 timer 的间隔(分钟)
	CheckIntervalMinutes int `json:"check_interval_minutes"`
	// 数据库驱动，"sqlite"(默认) 或 "postgres"
	DatabaseDriver string `json:"database_driver"`
	// sqlite 数据库文件路径，DatabaseDriver 为 sqlite 时使用
	DatabasePath string `json:"database_path"`
	// postgres 连接串，DatabaseDriver 为 postgres 时使用
	// 例如 postgres://user:pass@localhost:5432/goodbaby?sslmode=disable
	DatabaseDSN string `json:"database_dsn"`
	// session 加密密钥，为空时自动生成并写回配置文件
	SessionSecret string `json:"session_secret"`
	// session 有效期(小时)
	SessionMaxAgeHours int `json:"session_max_age_hours"`
	// 允许跨域访问的来源，开发前端时使用，如 ["http://localhost:5173"]
	AllowedOrigins []string `json:"allowed_origins"`
	// 每个用户保留的执行日志条数，<=0 表示不限制
	LogRetainCount int `json:"log_retain_count"`
}

// 支持的数据库驱动
const (
	DriverSQLite   = "sqlite"
	DriverPostgres = "postgres"
)

func defaultConfig() Config {
	return Config{
		ListenAddr:           ":8088",
		EnableRegistry:       true,
		TimeoutDurationHours: 6,
		CheckIntervalMinutes: 10,
		DatabaseDriver:       DriverSQLite,
		DatabasePath:         "data.db",
		SessionMaxAgeHours:   24 * 7,
		AllowedOrigins:       []string{},
		LogRetainCount:       500,
	}
}

var (
	once     sync.Once
	cached   Config
	loadErr  error
	cacheMux sync.RWMutex
)

// GetConfig 读取配置，只在首次调用时读盘，之后返回缓存。
//
// 配置文件不存在时会以默认值创建一份，缺失的字段会被补齐并写回。
func GetConfig() (Config, error) {
	once.Do(func() {
		cfg, err := load()
		cacheMux.Lock()
		cached, loadErr = cfg, err
		cacheMux.Unlock()
	})

	cacheMux.RLock()
	defer cacheMux.RUnlock()
	return cached, loadErr
}

// MustGetConfig 与 GetConfig 相同，但在出错时直接退出，用于启动阶段
func MustGetConfig() Config {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	return cfg
}

func configPath() string {
	if p := os.Getenv("GOODBABY_CONFIG"); p != "" {
		return p
	}
	return defaultConfigPath
}

func load() (Config, error) {
	path := configPath()
	config := defaultConfig()

	content, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := json.Unmarshal(content, &config); err != nil {
			return config, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}
	case errors.Is(err, fs.ErrNotExist):
		log.Printf("配置文件 %s 不存在，使用默认配置创建", path)
	default:
		return config, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}

	changed := normalize(&config)
	applyEnvOverrides(&config)

	if changed {
		if err := save(path, config); err != nil {
			// 只读挂载等场景下写回失败不应阻断启动
			log.Printf("警告: 写回配置文件失败: %v", err)
		}
	}

	return config, nil
}

// normalize 补齐缺失/非法的配置项，返回是否发生了修改
func normalize(c *Config) bool {
	def := defaultConfig()
	changed := false

	if c.ListenAddr == "" {
		c.ListenAddr, changed = def.ListenAddr, true
	}
	if c.TimeoutDurationHours <= 0 {
		c.TimeoutDurationHours, changed = def.TimeoutDurationHours, true
	}
	if c.CheckIntervalMinutes <= 0 {
		c.CheckIntervalMinutes, changed = def.CheckIntervalMinutes, true
	}
	if c.DatabasePath == "" {
		c.DatabasePath, changed = def.DatabasePath, true
	}
	if c.DatabaseDriver == "" {
		c.DatabaseDriver, changed = def.DatabaseDriver, true
	}
	if c.SessionMaxAgeHours <= 0 {
		c.SessionMaxAgeHours, changed = def.SessionMaxAgeHours, true
	}
	if c.SessionSecret == "" {
		// 随机生成一次并持久化，避免每次重启都让所有用户掉线
		c.SessionSecret, changed = randomHex(32), true
	}
	if c.AllowedOrigins == nil {
		c.AllowedOrigins, changed = def.AllowedOrigins, true
	}
	return changed
}

func applyEnvOverrides(c *Config) {
	if v := os.Getenv("GOODBABY_LISTEN_ADDR"); v != "" {
		c.ListenAddr = v
	}
	if v := os.Getenv("GOODBABY_DB_PATH"); v != "" {
		c.DatabasePath = v
	}
	if v := os.Getenv("GOODBABY_DB_DRIVER"); v != "" {
		c.DatabaseDriver = v
	}
	if v := os.Getenv("GOODBABY_DB_DSN"); v != "" {
		c.DatabaseDSN = v
		// 给了连接串却没显式指定驱动时，按 postgres 处理
		if os.Getenv("GOODBABY_DB_DRIVER") == "" {
			c.DatabaseDriver = DriverPostgres
		}
	}
	if v := os.Getenv("GOODBABY_SESSION_SECRET"); v != "" {
		c.SessionSecret = v
	}
	if v := os.Getenv("GOODBABY_ENABLE_REGISTRY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			c.EnableRegistry = b
		}
	}
}

func save(path string, c Config) error {
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o600)
}

// Save 持久化配置并刷新缓存，供管理接口修改运行时配置使用
func Save(c Config) error {
	normalize(&c)
	if err := save(configPath(), c); err != nil {
		return err
	}

	cacheMux.Lock()
	cached, loadErr = c, nil
	cacheMux.Unlock()
	return nil
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("无法生成随机密钥: %v", err))
	}
	return hex.EncodeToString(buf)
}
