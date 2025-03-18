package services

import (
	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type VisualizationService interface {
	BinarySearchSubmitCode(codeBlankString []string) (map[string]interface{}, error)
	SortGenCode(sortType string) (map[string]interface{}, error)
	SortSubmitCode(sortType string, codeBlankString []string) (map[string]interface{}, error)
	GraphGenCode(sortType string) (map[string]interface{}, error)
	GraphSubmitCode(sortType string, codeBlankString []string, graphInfo string) (map[string]interface{}, error)
	SortSubmitCodeNotInDocker(sortType string, codeBlankString []string) (map[string]interface{}, error)
	CustomSubmitCode(form *models.CustomSubmitCodeForm) (map[string]interface{}, error)
}

type visualizationService struct {
}

func NewVisualizationService() VisualizationService {
	return &visualizationService{}
}

func (s *visualizationService) BinarySearchSubmitCode(codeBlankString []string) (map[string]interface{}, error) {
	var run_py string
	var real_py = "binary_search.py"

	run_py = "run/" + strings.Replace(real_py, ".py", "_run_"+fmt.Sprintf("%d", time.Now().Unix())+".py", 1)
	real_py = real_py

	path := config.GetString("scripts.path")

	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	// 输入有$$CodeBlank$$，报错
	if len(codeBlankString) > 0 {
		// codeBlankString[0] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[0], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}
	for i := 1; i < len(codeBlankString); i++ {
		// codeBlankString[i] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[i], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}

	// 2. 构建 sed 命令
	var sedCmdBuilder strings.Builder
	if len(codeBlankString) > 0 {
		sedCmdBuilder.WriteString("s/\\$\\$CodeBlank\\[0\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[0])
		sedCmdBuilder.WriteString(" /g")
	}
	for i := 1; i < len(codeBlankString); i++ {
		sedCmdBuilder.WriteString("; s/\\$\\$CodeBlank\\[")
		sedCmdBuilder.WriteString(fmt.Sprintf("%d", i))
		sedCmdBuilder.WriteString("\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[i])
		sedCmdBuilder.WriteString(" /g")
	}
	sedCmd := sedCmdBuilder.String()

	// 3. 修改 run.py
	cmd := exec.Command("sed", "-i", sedCmd, path+run_py)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 4. 执行 run.py
	output, err := RunCodeInDocker("python:latest", "python", "/scripts/"+run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}

func (s *visualizationService) SortGenCode(sortType string) (map[string]interface{}, error) {
	var run_py string
	var real_py string

	if sortType == "bubbleSort" {
		real_py = "gen_bubble_sort_code.py"
	} else if sortType == "selectionSort" {
		real_py = "gen_selection_sort_code.py"
	} else if sortType == "insertionSort" {
		real_py = "gen_insertion_sort_code.py"
	} else if sortType == "quickSort" {
		real_py = "gen_quick_sort_code.py"
	}

	path := config.GetString("scripts.path")

	run_py = GetRunPy(real_py)
	real_py = real_py

	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 4. 执行 run.py
	output, err := RunCodeNotInDocker(run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}

func (s *visualizationService) SortSubmitCode(sortType string, codeBlankString []string) (map[string]interface{}, error) {
	var run_py string
	var real_py string

	if sortType == "bubbleSort" {
		real_py = "bubble_sort.py"
	} else if sortType == "selectionSort" {
		real_py = "selection_sort.py"
	} else if sortType == "insertionSort" {
		real_py = "insertion_sort.py"
	} else if sortType == "quickSort" {
		real_py = "quick_sort.py"
	}

	path := config.GetString("scripts.path")

	run_py = GetRunPy(real_py)
	real_py = real_py

	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 输入有$$CodeBlank$$，报错
	if len(codeBlankString) > 0 {
		// codeBlankString[0] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[0], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}
	for i := 1; i < len(codeBlankString); i++ {
		// codeBlankString[i] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[i], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}

	// 2. 构建 sed 命令
	var sedCmdBuilder strings.Builder
	if len(codeBlankString) > 0 {
		sedCmdBuilder.WriteString("s/\\$\\$CodeBlank\\[0\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[0])
		sedCmdBuilder.WriteString(" /g")
	}
	for i := 1; i < len(codeBlankString); i++ {
		sedCmdBuilder.WriteString("; s/\\$\\$CodeBlank\\[")
		sedCmdBuilder.WriteString(fmt.Sprintf("%d", i))
		sedCmdBuilder.WriteString("\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[i])
		sedCmdBuilder.WriteString(" /g")
	}
	sedCmd := sedCmdBuilder.String()

	// 3. 修改 run.py
	cmd := exec.Command("sed", "-i", sedCmd, path+run_py)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 4. 执行 run.py
	output, err := RunCodeInDocker("python:latest", "python", "/scripts/"+run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}

func (s *visualizationService) GraphGenCode(graphType string) (map[string]interface{}, error) {
	var run_py string
	var real_py string

	if graphType == "bfs" {
		real_py = "bfs.py"
	} else if graphType == "dfs" {
		real_py = "dfs.py"
	} else if graphType == "dijkstra" {
		real_py = "dijkstra.py"
	} else if graphType == "kruskal" {
		real_py = "kruskal.py"
	} else if graphType == "prim" {
		real_py = "prim.py"
	}

	path := config.GetString("scripts.path")

	run_py = GetRunPy(real_py)
	real_py = "graph/loadCode/" + real_py
	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 4. 执行 run.py
	output, err := RunCodeNotInDocker(run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}
func (s *visualizationService) GraphSubmitCode(graphType string, codeBlankString []string, graphInfo string) (map[string]interface{}, error) {
	var run_py string
	var real_py string

	if graphType == "bfs" {
		real_py = "bfs.py"
	} else if graphType == "dfs" {
		real_py = "dfs.py"
	} else if graphType == "dijkstra" {
		real_py = "dijkstra.py"
	} else if graphType == "kruskal" {
		real_py = "kruskal.py"
	} else if graphType == "prim" {
		real_py = "prim.py"
	}

	// run_py 加上当前时间戳防止重名，并且放入run文件夹中

	run_py = GetRunPy(real_py)
	real_py = "graph/submitCode/" + real_py

	path := config.GetString("scripts.path")

	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	//defer func() {
	//	// 删除运行临时文件
	//	err := os.Remove(path + run_py)
	//	if err != nil && !os.IsNotExist(err) {
	//		// 如果错误不是文件不存在的错误，则输出错误信息
	//		fmt.Println("delete run_py error: ", err)
	//	}
	//}()

	// 输入有$$CodeBlank$$，报错
	if len(codeBlankString) > 0 {
		// codeBlankString[0] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[0], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
		if strings.Contains(codeBlankString[0], "$$GraphInfo$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}
	for i := 1; i < len(codeBlankString); i++ {
		// codeBlankString[i] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[i], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
		if strings.Contains(codeBlankString[0], "$$GraphInfo$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}

	// 2. 构建 sed 命令
	var sedCmdBuilder strings.Builder
	if len(codeBlankString) > 0 {
		sedCmdBuilder.WriteString("s/\\$\\$CodeBlank\\[0\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[0])
		sedCmdBuilder.WriteString(" /g")
	}
	for i := 1; i < len(codeBlankString); i++ {
		sedCmdBuilder.WriteString("; s/\\$\\$CodeBlank\\[")
		sedCmdBuilder.WriteString(fmt.Sprintf("%d", i))
		sedCmdBuilder.WriteString("\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[i])
		sedCmdBuilder.WriteString(" /g")
	}
	sedCmdBuilder.WriteString("; s/\\$\\$GraphInfo\\$\\$/")
	sedCmdBuilder.WriteString(graphInfo)
	sedCmdBuilder.WriteString(" /g")
	sedCmd := sedCmdBuilder.String()

	// 3. 修改 run.py
	cmd := exec.Command("sed", "-i", sedCmd, path+run_py)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 4. 执行 run.py
	output, err := RunCodeInDocker("python:latest", "python", "/scripts/"+run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}

func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func getOutputMap(output string) (map[string]interface{}, error) {
	var outputMap map[string]interface{}
	err := json.Unmarshal([]byte(output), &outputMap)
	if err != nil {
		return nil, err
	}
	return outputMap, nil
}

func RunCodeInDocker(imageName string, args ...string) (string, error) {

	ctx := context.Background()

	// 创建 Docker 客户端
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", errors.New("创建 Docker 客户端失败: " + err.Error())
	}

	// 拉取 Python 镜像
	reader, err := dockerClient.ImagePull(ctx, "python:latest", types.ImagePullOptions{})
	if err != nil {
		return "", errors.New("拉取 Python 镜像失败: " + err.Error())
	}
	_, _ = io.ReadAll(reader)

	// 获取当前工作目录所在的驱动器盘符
	//drive := filepath.VolumeName(os.Getenv("PWD"))

	// 创建容器，并将上传的 Python 脚本作为卷挂载到容器中
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   append([]string{"timeout", "-s", "KILL", "-v", strconv.FormatInt(60, 10)}, args...),
	}, &container.HostConfig{
		Binds:       []string{fmt.Sprintf("%s:/scripts/run", "/root/backend/scripts/run")},
		NetworkMode: "none",
		Resources: container.Resources{
			Memory:     32 * 1024 * 1024,
			MemorySwap: 32 * 1024 * 1024,
		},
	}, nil, nil, "python_runner_"+uuid.New().String())
	if err != nil {
		return "", errors.New("创建容器失败: " + err.Error())
	}

	// 启动容器
	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", errors.New("启动容器失败: " + err.Error())
	}

	defer func() {
		// 检查容器是否存在
		_, err = dockerClient.ContainerInspect(ctx, resp.ID)
		if err != nil {
			if client.IsErrNotFound(err) {
				// 容器不存在，不需要进行停止和删除操作
				return
			}
		}
		// 停止并删除容器
		err = dockerClient.ContainerStop(ctx, resp.ID, container.StopOptions{})
		if err != nil {
			logger.Error("停止容器失败: " + err.Error())
		}

		err = dockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
		if err != nil {
			logger.Error("删除容器失败: " + err.Error())
		}
	}()

	//time.Sleep(1 * time.Second)

	// 等待容器运行完成
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	var status container.WaitResponse

	select {
	case err := <-errCh:
		if err != nil {
			// 处理错误
			return "", errors.New("容器运行错误: " + err.Error())
		}
	case status = <-statusCh:
	}

	// 获取容器的输出日志
	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: false})
	if err != nil {
		return "", errors.New("获取容器输出日志失败: " + err.Error())
	}

	defer func() {
		if err := out.Close(); err != nil {
			return
		}
	}()

	stdout, stderr := bytes.Buffer{}, bytes.Buffer{}
	if _, err := stdcopy.StdCopy(&stdout, &stderr, out); err != nil {
		return "", errors.New("读取容器输出日志失败: " + err.Error())
	}

	// 如果标准错误输出不为空，则返回错误信息
	if stderr.Len() != 0 {
		return "", errors.New("程序运行错误: " + stderr.String())
		//return "", errors.New("程序运行错误！")
	}

	if status.StatusCode != 0 {
		// 容器退出时的处理
		return "", errors.New("容器运行错误: " + "错误码: " + strconv.Itoa(int(status.StatusCode)))
	}

	return stdout.String() + stderr.String(), nil
}

func RunCodeNotInDocker(filename string) (string, error) {
	path := config.GetString("scripts.path")

	// 4. 执行 run.py
	cmd := exec.Command("python3", path+filename)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	// 编译信息不输出
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// 获取编译错误信息
		//errMsg := stderr.String()
		//return "", fmt.Errorf("程序运行错误 : %v\n%s", err, errMsg)
		return "", fmt.Errorf("程序运行错误!")
	}
	// 5. 获取 run.py 中的输出
	output := stdout.String()

	return output, nil
}

func (s *visualizationService) SortSubmitCodeNotInDocker(sortType string, codeBlankString []string) (map[string]interface{}, error) {
	var run_py string
	var real_py string

	if sortType == "bubbleSort" {
		real_py = "bubble_sort.py"
	} else if sortType == "selectionSort" {
		real_py = "selection_sort.py"
	} else if sortType == "insertionSort" {
		real_py = "insertion_sort.py"
	} else if sortType == "quickSort" {
		real_py = "quick_sort.py"
	}

	path := config.GetString("scripts.path")

	run_py = GetRunPy(real_py)
	real_py = real_py

	// 1. 拷贝文件
	err := copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 输入有$$CodeBlank$$，报错
	if len(codeBlankString) > 0 {
		// codeBlankString[0] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[0], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}
	for i := 1; i < len(codeBlankString); i++ {
		// codeBlankString[i] 含有 $$CodeBlank$$
		if strings.Contains(codeBlankString[i], "$$CodeBlank$$") {
			return nil, fmt.Errorf("内容不合法")
		}
	}

	// 2. 构建 sed 命令
	var sedCmdBuilder strings.Builder
	if len(codeBlankString) > 0 {
		sedCmdBuilder.WriteString("s/\\$\\$CodeBlank\\[0\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[0])
		sedCmdBuilder.WriteString(" /g")
	}
	for i := 1; i < len(codeBlankString); i++ {
		sedCmdBuilder.WriteString("; s/\\$\\$CodeBlank\\[")
		sedCmdBuilder.WriteString(fmt.Sprintf("%d", i))
		sedCmdBuilder.WriteString("\\]\\$\\$/")
		sedCmdBuilder.WriteString(codeBlankString[i])
		sedCmdBuilder.WriteString(" /g")
	}
	sedCmd := sedCmdBuilder.String()

	// 3. 修改 run.py
	cmd := exec.Command("sed", "-i", sedCmd, path+run_py)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 4. 执行 run.py
	output, err := RunCodeNotInDocker(run_py)
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}

func GetRunPy(real_py string) string {
	run_py := "run/" + strings.Replace(real_py, ".py", "_run_"+uuid.New().String()+".py", 1)
	return run_py
}

func GetCodePy(real_py string) string {
	run_py := "run/" + strings.Replace(real_py, ".py", "_code_"+uuid.New().String()+".py", 1)
	return run_py
}

func (s *visualizationService) CustomSubmitCode(customSubmitCodeForm *models.CustomSubmitCodeForm) (map[string]interface{}, error) {
	code := customSubmitCodeForm.Code

	real_py := "pythonExecParser.py"
	run_py := GetRunPy(real_py)
	code_py := GetCodePy(real_py)
	real_py = "custom/" + real_py

	// 保存 code 到 /run 目录下的code_py文件中
	path := config.GetString("scripts.path")
	file, err := os.OpenFile(path+code_py, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	// defer 删除code_py文件
	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + code_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete code_py error: ", err)
		}
	}()

	if _, err := fmt.Fprintf(file, "%s\n", code); err != nil {
		return nil, err
	}

	// 1. 拷贝文件
	err = copyFile(path+real_py, path+run_py)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 删除运行临时文件
		err := os.Remove(path + run_py)
		if err != nil && !os.IsNotExist(err) {
			// 如果错误不是文件不存在的错误，则输出错误信息
			fmt.Println("delete run_py error: ", err)
		}
	}()

	// 4. 执行 run.py
	output, err := RunCodeInDocker("python-with-pytutor", "python", "/scripts/"+run_py, "/scripts/"+code_py, "True")
	if err != nil {
		return nil, err
	}

	outputMap, err := getOutputMap(output)
	if err != nil {
		return nil, err
	}
	// 输出解析后的数据
	//fmt.Println(outputMap)

	return outputMap, nil
}
