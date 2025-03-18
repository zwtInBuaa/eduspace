package main_test

import (
	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/database"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	fmt.Println("Debug started")

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

	//// 初始化路由
	//r := routers.InitRouter(db)
	//
	//user0 := models.User{
	//	Role: 0,
	//}
	//
	//utils.

	//// 创建一个新的 HTTP 请求
	//req, err := http.NewRequest("GET", "/", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// 创建一个新的 HTTP 记录器
	//recorder := httptest.NewRecorder()
	//
	//// 调用处理程序函数并传入记录器和请求
	//r.ServeHTTP(recorder, req)
	//
	//// 检查响应代码
	//assert.Equal(t, http.StatusOK, recorder.Code)

	// 检查响应正文
	//assert.Contains(t, recorder.Body.String(), "Hello World!")
}
