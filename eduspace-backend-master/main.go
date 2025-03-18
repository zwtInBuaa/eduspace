package main

import (
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"log"
	"net/http"

	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/database"
	//"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/routers"
)

func main() {
	fmt.Println("Server started")

	// 读取配置文件
	err := config.InitConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to init config: %v", err))
	}

	// 初始化 log
	logger.InitLogger() // 服务重新启动时，日志会追加，不会删除

	// 初始化数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	// Close
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	defer sqlDB.Close()

	// 初始化 MinIO
	err = utils.InitMinio()
	if err != nil {
		panic(fmt.Sprintf("Failed to init MinIO: %v", err))
	}

	// 初始化路由
	r := routers.InitRouter(db)

	// 启动服务
	port := config.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Listening on %s...\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
