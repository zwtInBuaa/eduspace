package services

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
	"EDU_TH_2_backend/gin/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type UserService interface {
	Signup(c *gin.Context, username, buaaId, password string, role int64) (uint, error)
	Login(buaaId, password string) (string, uint, string, string, int64, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id uint, username, buaaId string, userId uint) error
	DeleteUser(c *gin.Context, id uint, userId uint) error
	SaveAvatar(c *gin.Context, id string, s string) (string, error)
	GetAvatar(c *gin.Context, id string) (string, error)
	GetTeacherCourses(teacherID uint) ([]*models.Course, error)
	GetStudentCourses(studentID uint) ([]map[string]interface{}, error)
	UpdatePassword(userID uint, oldPassword, newPassword string) error
	ResetPasswordByBuaaId(resetPasswordFrom *models.ResetPasswordForm, userId uint) error
	GetWeakness(userID uint) (string, error)
	RecQuestion(userID uint) ([]*models.Question, error)
	QuestionOverview(userID uint) (map[string]int, error)
	GetWeaknessLabel(userID uint) (map[string]interface{}, error)
}

type userService struct {
	userRepo  repositories.UserRepository
	tokenUtil utils.TokenUtil
}

func NewUserService(userRepo repositories.UserRepository, tokenUtil utils.TokenUtil) UserService {
	return &userService{userRepo: userRepo, tokenUtil: tokenUtil}
}

func (s *userService) Signup(c *gin.Context, username, buaaId, password string, role int64) (uint, error) {
	minioClient := utils.GetMinioClient()

	user := &models.User{Username: username, BuaaId: buaaId, Password: password, Role: role}
	if err := s.userRepo.Create(user); err != nil {
		return 0, err
	}

	bucketName := fmt.Sprintf("user-%d", int(user.ID))

	exists, err := minioClient.BucketExists(c, bucketName)
	if err != nil {
		s.userRepo.Delete(user.ID)
		return 0, err
	}

	if !exists {
		err = minioClient.MakeBucket(c, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			s.userRepo.Delete(user.ID)
			return 0, err
		}
	}

	return user.ID, nil
}

func (s *userService) Login(buaaId, password string) (string, uint, string, string, int64, error) {
	user, err := s.userRepo.FindByBuaaId(buaaId)
	if err != nil {
		return "", 0, "", "", 0, errors.New("buaaId错误！")
	}
	if !user.CheckPassword(password) {
		return "", 0, "", "", 0, errors.New("密码错误！")
	}
	token, err := s.tokenUtil.GenerateToken(user)
	if err != nil {
		return "", 0, "", "", 0, err
	}

	return token, user.ID, user.Username, user.BuaaId, user.Role, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) UpdateUser(id uint, username, buaaId string, userId uint) error {
	// role只能向下修改
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	changeUser, err := s.userRepo.FindByID(userId)
	if err != nil {
		return err
	}
	if changeUser.ID != user.ID && changeUser.Role != 0 && changeUser.Role >= user.Role {
		return errors.New("权限不足！")
	}

	user.Username = username
	user.BuaaId = buaaId
	if err := s.userRepo.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUser(c *gin.Context, id uint, userId uint) error {
	// role只能向下修改
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	changeUser, err := s.userRepo.FindByID(userId)
	if err != nil {
		return err
	}

	if changeUser.ID != user.ID && changeUser.Role != 0 && changeUser.Role >= user.Role {
		return errors.New("权限不足！")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return err
	}
	//// 如果用户有头像则删除
	//avatarPath := filepath.Join(".", "gin", "uploads", "avatars", strconv.FormatUint(uint64(id), 10)+"-"+"avatar.png") // 假设头像图片的路径为 /uploads/{userID}/avatar.jpg
	//// 检查文件是否存在
	//_, err := os.Stat(avatarPath)
	//fileExists := !os.IsNotExist(err)
	//
	//// 如果文件存在，则删除头像文件
	//if fileExists {
	//	err = os.Remove(avatarPath)
	//	if err != nil {
	//		// 处理文件删除错误
	//		return fmt.Errorf("Failed to delete avatar file for user %d %v", id, err.Error())
	//	}
	//}

	// 删除用户的桶
	minioClient := utils.GetMinioClient()

	bucketName := fmt.Sprintf("user-%d", int(id))

	err = minioClient.RemoveBucket(c, bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) SaveAvatar(c *gin.Context, userID, image string) (string, error) {
	filename := "avatar.png"

	// 将Base64编码数据保存到文件中
	//savePath := filepath.Join(".", "gin", "uploads", "avatars", userID+"-"+filename)
	//saveFile, err := os.Create(savePath)
	//if err != nil {
	//	return "", fmt.Errorf("unable to create save file: %v", err)
	//}
	//defer saveFile.Close()

	// 去除前缀信息
	image = strings.ReplaceAll(image, "data:image/png;base64,", "")
	saveData, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return "", fmt.Errorf("unable to decode Base64 encoded data: %v", err)
	}

	// 将Base64编码数据保存到MinIO中
	minioClient := utils.GetMinioClient()

	bucketName := fmt.Sprintf("user-%s", userID)

	_, err = minioClient.PutObject(c, bucketName, filename, bytes.NewReader(saveData), int64(len(saveData)), minio.PutObjectOptions{ContentType: "image/png"})

	// 生成一个长期有效的签名 URL
	url, err := utils.GetAvatarURL(c, userID)
	if err != nil {
		return "", err
	}

	//_, err = saveFile.Write(saveData)
	//if err != nil {
	//	return "", fmt.Errorf("unable to write file: %v", err)
	//}

	return url, nil
}

func (s *userService) GetAvatar(c *gin.Context, userID string) (string, error) {
	// 功能迁移到minio_util.go中
	return "", nil
}

func (s *userService) GetTeacherCourses(teacherID uint) ([]*models.Course, error) {
	return s.userRepo.GetTeacherCourses(teacherID)
}

func (s *userService) GetStudentCourses(studentID uint) ([]map[string]interface{}, error) {
	courses, err := s.userRepo.GetStudentCourses(studentID)
	if err != nil {
		return nil, err
	}

	var courseInfos []map[string]interface{}
	for _, course := range courses {
		courseInfo := map[string]interface{}{
			"id":          course.ID,
			"name":        course.Name,
			"description": course.Description,
			"created_at":  course.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at":  course.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		courseInfos = append(courseInfos, courseInfo)
	}
	return courseInfos, nil
}

func (s *userService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	// 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 检查旧密码是否正确
	if !user.CheckPassword(oldPassword) {
		return errors.New("Invalid old password")
	}

	// 更新密码
	user.Password = newPassword

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// log:哪个用户干了啥事
	logger.Info("UpdatePassword: " + user.Username + " (" + user.BuaaId + ")")

	return nil
}

func (s *userService) ResetPasswordByBuaaId(resetPasswordFrom *models.ResetPasswordForm, userId uint) error {
	// role只能向下修改
	user, err := s.userRepo.FindByBuaaId(resetPasswordFrom.BuaaID)
	if err != nil {
		return err
	}
	changeUser, err := s.userRepo.FindByID(userId)
	if err != nil {
		return err
	}

	if changeUser.ID != user.ID && changeUser.Role != 0 && changeUser.Role >= user.Role {
		return errors.New("权限不足！")
	}

	user.Password = resetPasswordFrom.NewPassword

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetWeakness(userID uint) (string, error) {
	// Get all of the user's submission history
	history, err := s.userRepo.GetUserSubmitHistory(userID)
	if err != nil {
		return "", err
	}

	// Count the number of correct and incorrect submissions for each tag
	tagCounts := make(map[string][2]int)

	for _, submission := range history {
		if submission.IsCorrect {
			for _, Tag := range submission.Tags {
				count := tagCounts[Tag]
				count[0]++
				tagCounts[Tag] = count
			}
		} else {
			for _, Tag := range submission.Tags {
				count := tagCounts[Tag]
				count[1]++
				tagCounts[Tag] = count
			}
		}
	}

	//logger.Info(fmt.Sprintf("%v", tagCounts))

	// Calculate the ratio of correct to incorrect submissions for each tag
	tagRatios := make(map[string]float64)
	for tag, counts := range tagCounts {
		if counts[0]+counts[1] > 0 {
			tagRatios[tag] = float64(counts[0]) / float64(counts[0]+counts[1])
		} else {
			tagRatios[tag] = 0.0
		}
	}

	//调用 scipts/weakness_predict/main.py,分别传入tagRatios的每一项"排序，图，树，分治，基础知识，循环和分支，模拟", 返回预测结果
	//err = os.Chdir("./gin/scripts/weakness_predict")
	//if err != nil {
	//	return "", err
	//}

	//cmd := exec.Command("D:\\PycharmProjects\\weakness\\venv\\Scripts\\python",
	cmd := exec.Command("python3",
		"-c",
		"import os; os.chdir('./gin/scripts/weakness_predict'); import main;",
		strconv.FormatFloat(tagRatios["排序"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["图"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["树"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["分治"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["基础知识"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["循环和分支"], 'f', -1, 64),
		strconv.FormatFloat(tagRatios["模拟"], 'f', -1, 64),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	// TODO 编译信息不输出
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// 获取编译错误信息
		errMsg := stderr.String()
		return "", fmt.Errorf("running py failed: %v\n%s", err, errMsg)
		//return nil, err
	}
	output := stdout.String()
	//// 5. 获取 run.py 中的输出, Convert GBK to UTF-8
	//outputReader := transform.NewReader(&stdout, simplifiedchinese.GB18030.NewDecoder())
	//outputBytes, err := ioutil.ReadAll(outputReader)
	//if err != nil {
	//	return "", err
	//}
	//output := string(outputBytes)
	output = strings.TrimRight(output, "\r\n")

	logger.Info(fmt.Sprintf("GetWeakness: %v", output))

	return output, nil
}

func (s *userService) RecQuestion(userID uint) ([]*models.Question, error) {
	// 先调用Weakness，获取用户的弱项
	weakness, err := s.GetWeakness(userID)
	fmt.Println(weakness)
	if err != nil {
		return nil, err
	}

	weakness = strings.ReplaceAll(weakness, "'", "\"")

	// 把string的weakness（json形式）转换成[]string
	var weaknessList []string
	err = json.Unmarshal([]byte(weakness), &weaknessList)
	if err != nil {
		return nil, err
	}

	// 使用weaknessList中的每一项，调用questionRepo的GetQuestionByTag，获取推荐的题目
	var questionList []*models.Question
	for _, tag := range weaknessList {
		questions, err := s.userRepo.GetQuestionByTag(tag)
		if err != nil {
			return nil, err
		}
		for i := range questions {
			questionList = append(questionList, questions[i])
		}
	}

	//logger.Info(fmt.Sprintf("questionList: %v", questionList))

	// questionList中可能有重复的题目，去重
	questionMap := make(map[uint]*models.Question)
	for _, question := range questionList {
		questionMap[question.ID] = question
	}
	questionList = make([]*models.Question, 0)
	for _, question := range questionMap {
		questionList = append(questionList, question)
	}

	// 如果推荐的题目数量不足10题，再从所有题目中随机选取
	// 如果推荐的题目数量超过10题，随机选取10题
	if len(questionList) < 10 {
		questions, err := s.userRepo.GetAllQuestion()
		if err != nil {
			return nil, err
		}
		for i := range questions {
			questionList = append(questionList, questions[i])
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questionList), func(i, j int) {
			questionList[i], questionList[j] = questionList[j], questionList[i]
		})
		questionList = questionList[:10]
	}

	return questionList, nil
}

func (s *userService) QuestionOverview(userID uint) (map[string]int, error) {
	// 返回正确题目数量，错误题目数量，未做题目数量
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	// Get all of the user's submission history
	history, err := s.userRepo.GetUserSubmitHistory(userID)
	if err != nil {
		return nil, err
	}
	// 获取上述history中该userID的所有题目id，统计对了，错了，没做的数量（占所有question的比例）
	var doneQuestion, correctQuestion, wrongQuestion []uint
	for i := range history {
		if history[i].UserID != userID {
			continue
		}
		if history[i].IsCorrect == true {
			correctQuestion = append(correctQuestion, history[i].QuestionID)
		} else {
			wrongQuestion = append(wrongQuestion, history[i].QuestionID)
		}
		doneQuestion = append(doneQuestion, history[i].QuestionID)
	}
	// 获取所有题目的id
	questionIDs, err := s.userRepo.GetAllQuestionID()
	if err != nil {
		return nil, err
	}
	// 计算正确题目数量，错误题目数量，未做题目数量
	var correctNum, wrongNum, notDoneNum int
	for i := range questionIDs {
		if contains(doneQuestion, questionIDs[i]) {
			if contains(correctQuestion, questionIDs[i]) {
				correctNum++
			} else {
				wrongNum++
			}
		} else {
			notDoneNum++
		}
	}
	return map[string]int{
		"correctNum": correctNum,
		"wrongNum":   wrongNum,
		"notDoneNum": notDoneNum,
	}, nil
}

func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *userService) GetWeaknessLabel(userID uint) (map[string]interface{}, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	// Get all of the user's submission history
	history, err := s.userRepo.GetUserSubmitHistory(userID)
	if err != nil {
		return nil, err
	}
	// 获取history中该userID的所有题目tag, 统计每个tag的正确数量和错误数量
	var tagCorrectNum, tagWrongNum map[string]int
	tagCorrectNum = make(map[string]int)
	tagWrongNum = make(map[string]int)
	for i := range history {
		if history[i].UserID != userID {
			continue
		}
		for _, tag := range history[i].Tags {
			if history[i].IsCorrect == true {
				tagCorrectNum[tag]++
			} else {
				tagWrongNum[tag]++
			}
		}
	}
	// 计算每个tag的正确率
	tagRatios := make(map[string]float64)
	for tag, correctNum := range tagCorrectNum {
		tagRatios[tag] = float64(correctNum) / float64(tagCorrectNum[tag]+tagWrongNum[tag])
	}

	var labels []map[string]interface{}
	var data []map[string]interface{}

	needLabel := []string{"排序", "图", "树", "分治", "基础知识", "循环和分支", "模拟"}
	for i := range needLabel {
		if _, ok := tagRatios[needLabel[i]]; !ok {
			tagRatios[needLabel[i]] = 0
		}
	}
	for i := range needLabel {
		labels = append(labels, map[string]interface{}{
			"name": needLabel[i],
			"max":  1,
		})
		data = append(data, map[string]interface{}{
			"label":  needLabel[i],
			"weight": tagRatios[needLabel[i]],
		})
	}

	return map[string]interface{}{
		"labels": labels,
		"data":   data,
	}, nil
}
