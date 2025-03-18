package controllers

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/services"
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (ctrl *CommentController) CreateComment(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content, ok := reqBody["content"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
		return
	}

	postIdFloat64, ok := reqBody["post_id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	postId, _ := utils.ParseFloat64ToUint(postIdFloat64)

	userId := c.GetInt("userId")

	if err := ctrl.commentService.CreateComment(content, uint(userId), postId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户在哪个帖子下发表了评论
	logger.Info(fmt.Sprintf("user %d create comment in post %d", userId, postId))

	c.JSON(http.StatusOK, gin.H{"message": "create comment success"})
}

func (ctrl *CommentController) GetComment(c *gin.Context) {
	commentId, _ := utils.ParseStringToUint(c.Param("comment_id"))

	comment, err := ctrl.commentService.GetCommentByID(commentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (ctrl *CommentController) GetCommentsByPost(c *gin.Context) {
	postId, _ := utils.ParseStringToUint(c.Param("post_id"))

	comments, err := ctrl.commentService.GetCommentsByPostID(postId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range comments {
		comments[i].Avatar, err = utils.GetAvatarURL(c, strconv.Itoa(int(comments[i].UserID)))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (ctrl *CommentController) UpdateComment(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("comment_id"))

	userId := c.GetInt("userId")
	userRole := c.GetInt("role")

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content, ok := reqBody["content"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
		return
	}

	if err := ctrl.commentService.UpdateComment(id, content, userId, userRole); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// logs: 哪个用户修改了哪条评论
	logger.Info(fmt.Sprintf("user %d update comment %d", userId, id))

	c.JSON(http.StatusOK, gin.H{"message": "update comment success"})
}

func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	commentId, _ := utils.ParseStringToUint(c.Param("comment_id"))

	userId := c.GetInt("userId")
	userRole := c.GetInt("role")

	if err := ctrl.commentService.DeleteComment(commentId, userId, userRole); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// logs: 哪个用户删除了哪条评论
	logger.Info(fmt.Sprintf("user %d delete comment %d", userId, commentId))

	c.JSON(http.StatusOK, gin.H{"message": "delete comment success"})
}
