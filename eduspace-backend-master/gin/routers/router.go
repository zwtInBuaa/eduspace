package routers

import (
	"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/utils"
	"bytes"
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc

	Middlewares []gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

var tokenUtil utils.TokenUtil
var enforcer *casbin.Enforcer

func InitRouter(db *gorm.DB) *gin.Engine {
	// 初始化路由
	router := gin.Default()

	// 将本地文件夹/gin/uploads映射到URL路径/uploads
	//router.Static("/uploads", "/gin/uploads")
	router.Static("/static", "gin/static")

	tokenUtil = utils.NewTokenUtil("my-secret-key", 3600)

	// 初始化Casbin
	enforcer = utils.CasbinInit()

	// 初始化下游Controller
	NewUserRouter(db)
	NewPostRouter(db)
	NewCommentRouter(db)
	NewUtilRouter(db)
	NewVisualizationRouter(db)
	NewCourseRouter(db)
	NewExamRouter(db)
	NewQuestionRouter(db)
	NewPermissionRouter(db)

	// 设置中间件

	// 添加CORS中间件
	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:8080", "http://114.116.211.180:8080"} // 允许来自http://localhost:8080的请求
	config.AllowOrigins = []string{"*"} // 允许来自任何域的请求
	//config.AllowHeaders = []string{"Authorization"} // 要写所有的，而不只是authorization
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// 注册路由组，注册控制器
	SetupRouter(router.Group("/"), routes)
	SetupRouter(router.Group("/users"), userRoutes)
	SetupRouter(router.Group("/posts"), postRoutes)
	SetupRouter(router.Group("/comments"), commentRoutes)
	SetupRouter(router.Group("/utils"), utilRoutes)
	SetupRouter(router.Group("/visualizations"), visualizationRoutes)
	SetupRouter(router.Group("/courses"), courseRoutes)
	SetupRouter(router.Group("/exams"), examRoutes)
	SetupRouter(router.Group("/questions"), questionRoutes)

	return router
}

func SetupRouter(r *gin.RouterGroup, routes Routes) {

	for _, route := range routes {

		midwares := route.Middlewares
		midwares = append(midwares, middlewares.CasbinMiddleware(enforcer), route.HandlerFunc)

		switch route.Method {
		case http.MethodGet:
			r.GET(route.Pattern, midwares...)
		case http.MethodPost:
			r.POST(route.Pattern, midwares...)
		case http.MethodPut:
			r.PUT(route.Pattern, midwares...)
		case http.MethodDelete:
			r.DELETE(route.Pattern, midwares...)
		case http.MethodPatch:
			r.PATCH(route.Pattern, midwares...)
		}
	}
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/",
		Index,
		[]gin.HandlerFunc{},
	},
	{
		"TestDocker",
		http.MethodGet,
		"/testdocker",
		TestDocker,
		[]gin.HandlerFunc{},
	},
}

func TestDocker(c *gin.Context) {
	ctx := context.Background()

	// 创建 Docker 客户端
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "创建 Docker 客户端失败: " + err.Error(),
		})
		return
	}

	// 拉取 Python 镜像
	reader, err := dockerClient.ImagePull(ctx, "python:latest", types.ImagePullOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "拉取 Python 镜像失败: " + err.Error(),
		})
		return
	}
	_, _ = ioutil.ReadAll(reader)

	// 获取当前工作目录所在的驱动器盘符
	//drive := filepath.VolumeName(os.Getenv("PWD"))

	path, err := os.Getwd()
	if err != nil {
		// 处理错误
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "获取当前工作目录失败: " + err.Error(),
		})
	}

	// 创建容器，并将上传的 Python 脚本作为卷挂载到容器中
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "python:latest",
		Cmd:   []string{"python", "/scripts/run/" + "docker_test.py"},
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/scripts", path+"/gin/scripts")},
	}, nil, nil, "")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "创建容器失败: " + err.Error(),
		})
		return
	}

	// 启动容器
	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "启动容器失败: " + err.Error(),
		})
		return
	}

	time.Sleep(5 * time.Second)

	// 等待容器运行完成
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	// 获取容器的输出日志
	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "获取容器输出日志失败: " + err.Error(),
		})
		return
	}

	defer func() {
		if err := out.Close(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "关闭容器输出日志失败: " + err.Error(),
			})
		}
	}()
	stdout, stderr := bytes.Buffer{}, bytes.Buffer{}
	if _, err := stdcopy.StdCopy(&stdout, &stderr, out); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "读取容器输出日志失败: " + err.Error(),
		})
		return
	}

	// 读取输出日志并删除容器
	err = dockerClient.ContainerStop(ctx, resp.ID, container.StopOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "停止容器失败: " + err.Error(),
		})
		return
	}

	err = dockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "删除容器失败: " + err.Error(),
		})
		return
	}

	select {
	case err := <-errCh:
		if err != nil {
			// 处理错误
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "容器运行错误: " + err.Error(),
			})
			return
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			// 容器退出时的处理
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "容器运行错误: " + "错误码: " + strconv.Itoa(int(status.StatusCode)),
			})
			return
		}
	}

	// 如果标准错误输出不为空，则返回错误信息
	if stderr.Len() != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "容器运行错误: " + stderr.String(),
		})
		return
	}

	c.String(http.StatusOK, stdout.String()+stderr.String())
}
