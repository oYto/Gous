package config

import (
	rlog "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"

	"sync"
)

var (
	config GlobalConfig // 全局配置文件
	once   sync.Once    // 只执行一次的代码
)

// LogConf 日志配置
type LogConf struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"` // 日志模式
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`       // 日志路径
	SaveDays   uint   `yaml:"save_days" mapstructure:"save_days"`     // 日志保存天数
	Level      string `yaml:"level" mapstructure:"level"`             // 日志级别
}

// DbConf 数据库配置
type DbConf struct {
	Host        string `yaml:"host" mapstructure:"host"`                   // 主机地址
	Port        string `yaml:"port" mapstructure:"port"`                   // 端口号
	User        string `yaml:"user" mapstructure:"user"`                   // 用户名
	Password    string `yaml:"password" mapstructure:"password"`           // 密码
	Dbname      string `yaml:"dbname" mapstructure:"dbname"`               // 数据库名
	MaxIdleConn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"` // 最大空闲连接数
	MaxOpenConn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"` // 最大打开连接数
	MaxIdleTime int64  `yaml:"max_idle_time" mapstructure:"max_idle_time"` // 连接最大空闲时间
}

// AppConf 服务配置
type AppConf struct {
	AppName string `yaml:"app_name" mapstructure:"app_name"` //	业务名
	Version string `yaml:"version" mapstructure:"version"`   // 版本
	Port    int    `yaml:"port" mapstructure:"port"`         // 端口
	RunMode string `yaml:"run_mode" mapstructure:"run_mode"` // 运行模式
}

// RedisConf 配置
type RedisConf struct {
	Host     string `yaml:"rhost" mapstructure:"rhost"`       // db主机地址
	Port     int    `yaml:"rport" mapstructure:"rport"`       // db端口
	DB       int    `yaml:"rdb" mapstructure:"rdb"`           // 数据库
	PassWord string `yaml:"passwd" mapstructure:"passwd"`     // 密码
	PoolSile int    `yaml:"poolsize" mapstructure:"poolsize"` // 连接池大小，即最大连接数
}

// Cache 配置
type Cache struct {
	SessionExpired int `yaml:"session_expired" mapstructure:"session_expired"` // 会话过期时间
	UserExpired    int `yaml:"user_expired" mapstructure:"user_expired"`       // 用户信息过期时间
}

// GlobalConfig 业务配置结构体
type GlobalConfig struct {
	AppConfig   AppConf   `yaml:"app" mapstructure:"app"`                 // 服务配置
	CorsOrigin  []string  `yaml:"cors_origin" mapstructure:"cors_origin"` // 跨域源列表
	DbConfig    DbConf    `yaml:"db" mapstructure:"db"`                   // 数据库配置
	LogConfig   LogConf   `yaml:"log" mapstructure:"log"`                 // 日志配置
	RedisConfig RedisConf `yaml:"redis" mapstructure:"redis"`             // redis 配置
	Cache       Cache     `yaml:"cache" mapstructure:"cache"`             // 缓存配置
}

// GetGlobalConf 获取全局配置文件
func GetGlobalConf() *GlobalConfig {
	once.Do(readConf)
	return &config
}

// 将配置文件中的信息全部加载到 全局配置文件中
func readConf() {
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		panic("read config file err:" + err.Error())
	}
	err = viper.Unmarshal(&config) // 将配置信息反序列化填充到全局配置文件中
	if err != nil {
		panic("config file unmarshal err:" + err.Error())
	}
	log.Infof("config === %+v", config)
}

// InitConfig 初始化日志
func InitConfig() {
	globalConf := GetGlobalConf() // 获取全局配置文件
	// 根据我们自己设置的日志级别来设置日志级别
	level, err := log.ParseLevel(globalConf.LogConfig.Level)
	if err != nil {
		panic("log level parse err:" + err.Error())
	}
	log.SetFormatter(&logFormatter{
		log.TextFormatter{
			DisableColors:   true,                  // 禁止日志输出中的颜色
			TimestampFormat: "2006-01-02 15:04:05", // 指定日志时间戳
			FullTimestamp:   true,                  // 启用完整的时间戳格式
		}})
	log.SetReportCaller(true) // 打印文件位置，行号
	log.SetLevel(level)       // 设置日志的级别
	switch globalConf.LogConfig.LogPattern {
	case "stdout": // 将日志输出到标准输出
		log.SetOutput(os.Stdout)
	case "stderr": // 将日志输出到标准错误输出
		log.SetOutput(os.Stderr)
	case "file": // 使用第三方日志库 rlog 创建一个新的日志记录器，并将其设置为日志的输出目标
		logger, err := rlog.New( // 基于文件的日志记录器，将日志输出到指定的文件中。
			globalConf.LogConfig.LogPath+".%Y%m%d",                // 文件名：路径+日期
			rlog.WithRotationCount(globalConf.LogConfig.SaveDays), // 设置日志文件保留天数
			rlog.WithRotationTime(time.Hour*24),                   // 设置日志文件的轮转时间间隔
		)
		if err != nil {
			panic("log conf err:" + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log conf err, check log_pattern in logsvr.yaml")
	}
}
