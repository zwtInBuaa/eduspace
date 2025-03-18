package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var v *viper.Viper

// InitConfig 读取配置文件
func InitConfig() error {
	v = viper.New()
	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath("./gin/config/")

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	return nil
}

func GetV() *viper.Viper {
	return v
}

func GetString(key string) string {
	return v.GetString(key)
}

func GetBool(key string) bool {
	return v.GetBool(key)
}

func GetInt(key string) int {
	return v.GetInt(key)
}

func GetInt64(key string) int64 {
	return v.GetInt64(key)
}

func GetFloat64(key string) float64 {
	return v.GetFloat64(key)
}
