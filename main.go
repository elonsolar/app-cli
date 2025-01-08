package main

import (
	"app-cli/conf"
	"app-cli/data"
	"app-cli/handler"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config") // 配置文件名称（不带扩展名）
	viper.AddConfigPath(".")      // 配置文件所在路径
	err := viper.ReadInConfig()   // 读取配置文件
	if err != nil {
		slog.Error("read config failed", "err", err)
		return
	}

	var cfg conf.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("unmarshal config failed", "err", err)
		return
	}

	db := data.NewGormDB(&cfg)
	data := data.NewData(&cfg, db)

	handler := handler.NewHandler(&cfg, data)

	router := gin.Default()
	router.StaticFS("/static/download", http.Dir(cfg.DownloadDir))
	router.POST("/packTask", handler.CreatePackTask)
	router.GET("/packTask", handler.ListPackTask)
	router.POST("/packTask/download", handler.Download)

	router.Run(":8080")
}
