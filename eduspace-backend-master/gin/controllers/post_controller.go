package controllers

import (
	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/services"
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type PostController struct {
	postService services.PostService
}

func NewPostController(postService services.PostService) *PostController {
	return &PostController{postService: postService}
}

func (ctrl *PostController) CreatePost(c *gin.Context) {

	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "你每 5 秒只能发表一次帖子",
			})
			return
		}
	}

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title, ok := reqBody["title"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
		return
	}

	content, ok := reqBody["content"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
		return
	}

	userId := c.GetInt("userId")

	if err := ctrl.postService.CreatePost(title, content, uint(userId)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	c.JSON(http.StatusOK, gin.H{"message": "create post success"})
}

func (ctrl *PostController) GetPost(c *gin.Context) {
	postId, _ := utils.ParseStringToUint(c.Param("post_id"))

	post, err := ctrl.postService.GetPostByID(postId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	avatar, err := utils.GetAvatarURL(c, strconv.Itoa(int(post["poster_id"].(uint))))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post["poster_avatar"] = avatar

	c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) GetAllPosts(c *gin.Context) {
	posts, err := ctrl.postService.GetAllPosts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range posts {
		avatar, err := utils.GetAvatarURL(c, strconv.Itoa(int(posts[i]["poster_id"].(uint))))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		posts[i]["poster_avatar"] = avatar
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

//func (ctrl *PostController) GetPostsByUser(c *gin.Context) {
//	userID := utils.GetUserIDFromContext(c)
//
//	posts, err := ctrl.postService.GetPostsByUserID(userID)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"posts": posts})
//}

func (ctrl *PostController) UpdatePost(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("post_id"))
	userId := c.GetInt("userId")
	userRole := c.GetInt("role")

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title, ok := reqBody["title"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
		return
	}

	content, ok := reqBody["content"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
		return
	}

	if err := ctrl.postService.UpdatePost(id, title, content, userId, userRole); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update post success"})
}

func (ctrl *PostController) DeletePost(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("post_id"))

	userId := c.GetInt("userId")
	userRole := c.GetInt("role")

	if err := ctrl.postService.DeletePost(id, userId, userRole); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户删除了哪个帖子
	logger.Info(fmt.Sprintf("user %d delete post %d", userId, id))

	c.JSON(http.StatusOK, gin.H{"message": "delete post success"})
}
