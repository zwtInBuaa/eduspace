package utils

import (
	"EDU_TH_2_backend/gin/config"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
)

func CasbinInit() *casbin.Enforcer {
	var enforcer *casbin.Enforcer

	switch config.GetString("casbin.adapter_type") {
	case "mysql":
		// 初始化 CasbinEnforcer 和 XormAdapter
		// 使用 MySQL 数据库初始化一个 Xorm 适配器
		a, err := xormadapter.NewAdapter("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			config.GetString("casbin.mysql.user"),
			config.GetString("casbin.mysql.password"),
			config.GetString("casbin.mysql.host"),
			config.GetString("casbin.mysql.port"),
		))
		if err != nil {
			panic(fmt.Sprintf("error: adapter: %s", err))
		}

		// 从 model.conf 文件中加载 Casbin 模型
		m, err := model.NewModelFromFile(config.GetString("casbin.model_path"))
		if err != nil {
			panic(fmt.Sprintf("error: model: %s", err))
		}

		enforcer, err = casbin.NewEnforcer(m, a)
		if err != nil {
			panic(fmt.Sprintf("error: enforcer: %s", err))
		}
	default:
		// 处理不支持的适配器类型错误
		panic(fmt.Sprintf("error: adapter_type: %s", config.GetString("casbin.adapter_type")))
	}

	// 先清除所有策略
	enforcer.ClearPolicy()

	// 角色权限
	// 管理员可以访问所有路由
	enforcer.AddPolicy("管理员", "/*", ".*")

	// 老师相比于助教额外可以访问的路由
	enforcer.AddPolicy("老师", "/courses/*", ".*")
	enforcer.AddPolicy("老师", "/users/*", ".*")
	enforcer.AddPolicy("老师", "/questions/*", ".*")
	enforcer.AddPolicy("老师", "/exams/*", ".*")

	// 助教相比于学生额外可以访问的路由
	enforcer.AddPolicy("助教", "/users/:user_id", ".*")

	// 学生可以访问的路由
	enforcer.AddPolicy("学生", "/users/:user_id", "PUT") // 仅限自己的
	enforcer.AddPolicy("学生", "/users/get-avatar/:user_id", ".*")
	enforcer.AddPolicy("学生", "/users/logout", ".*")                  // 仅限自己的
	enforcer.AddPolicy("学生", "/users/:user_id/changePassword", ".*") // 仅限自己的
	enforcer.AddPolicy("学生", "/users/upload-avatar/:user_id", ".*")
	enforcer.AddPolicy("学生", "/users/:user_id/weakness", ".*")
	enforcer.AddPolicy("学生", "/users/:user_id/recQuestion", ".*")
	enforcer.AddPolicy("学生", "/users/:user_id/questionOverview", ".*")
	enforcer.AddPolicy("学生", "/posts/*", ".*")    // 删除和修改仅限自己的
	enforcer.AddPolicy("学生", "/comments/*", ".*") // 删除和修改仅限自己的
	enforcer.AddPolicy("学生", "/courses/getall", ".*")
	enforcer.AddPolicy("学生", "/courses/:course_id", "GET")
	enforcer.AddPolicy("学生", "/courses/:course_id/exams", "GET")
	enforcer.AddPolicy("学生", "/courses/:course_id/questions", "GET")
	enforcer.AddPolicy("学生", "/courses/:course_id/teachers", "GET")
	enforcer.AddPolicy("学生", "/questions/getall", ".*")
	enforcer.AddPolicy("学生", "/questions/:question_id", "GET")
	enforcer.AddPolicy("学生", "/questions/:question_id/submit", "POST")
	enforcer.AddPolicy("学生", "/exams/:exam_id", "GET")
	enforcer.AddPolicy("学生", "/exams/:exam_id/questions", "GET")
	enforcer.AddPolicy("学生", "/utils/imgdbs", "POST")
	enforcer.AddPolicy("学生", "/visualizations/*", ".*")

	enforcer.AddGroupingPolicy("管理员", "老师")
	enforcer.AddGroupingPolicy("老师", "助教")
	enforcer.AddGroupingPolicy("助教", "学生")

	// 用户角色
	//enforcer.AddGroupingPolicy("1", "管理员")
	//enforcer.AddGroupingPolicy("user2", "老师")
	//enforcer.AddGroupingPolicy("user3", "学生")
	//enforcer.AddGroupingPolicy("user4", "助教")

	// 更改user1的管理员角色为老师
	//enforcer.RemoveGroupingPolicy("user1", "管理员")
	//enforcer.AddGroupingPolicy("user1", "老师")

	// 给user1增加/login的POST权限
	//enforcer.RemovePolicy("1", "/users/upload-avatar/1", "POST")

	enforcer.SavePolicy()

	return enforcer
}
