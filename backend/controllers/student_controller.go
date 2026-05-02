package controllers

import (
	"net/http"
	"retakeCoursesManage/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	query := models.DB.Preload("Author").Order("is_top DESC, created_at DESC")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Announcement{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       announcements,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

func GetAnnouncementDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var announcement models.Announcement
	if err := models.DB.Preload("Author").First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

func GetRetakeCourses(c *gin.Context) {
	var courses []models.RetakeCourse
	query := models.DB.Preload("Course").Preload("Teacher")

	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.RetakeCourse{}).Where(query.Where).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取重修课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMyGrades(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var grades []models.Grade

	query := models.DB.Preload("RetakeCourse.Course").Preload("Student").
		Where("student_id = ?", userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Grade{}).Where("student_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成绩失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      grades,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMyReports(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var reports []models.Report

	query := models.DB.Preload("User").Preload("ReplyUser").
		Where("user_id = ?", userID)

	if reportType := c.Query("type"); reportType != "" {
		query = query.Where("type = ?", reportType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Report{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	report := models.Report{
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		UserID:  userID.(uint),
		Status:  "pending",
	}

	if err := models.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "报告创建成功", "id": report.ID})
}

func GetMyApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var applications []models.RetakeApplication

	query := models.DB.Preload("Student").Preload("RetakeCourse.Course").
		Preload("Reviewer").Where("student_id = ?", userID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.RetakeApplication{}).Where("student_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取申请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      applications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		RetakeCourseID uint   `json:"retake_course_id" binding:"required"`
		Reason         string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var existingApplication models.RetakeApplication
	if err := models.DB.Where("student_id = ? AND retake_course_id = ?", userID, req.RetakeCourseID).
		First(&existingApplication).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已申请过该重修课程"})
		return
	}

	application := models.RetakeApplication{
		StudentID:      userID.(uint),
		RetakeCourseID: req.RetakeCourseID,
		Reason:         req.Reason,
		Status:         "pending",
	}

	if err := models.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建申请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "申请创建成功", "id": application.ID})
}

func CancelApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var application models.RetakeApplication
	if err := models.DB.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
		return
	}

	if application.StudentID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	if application.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能取消待审核的申请"})
		return
	}

	application.Status = "cancelled"
	if err := models.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消申请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "申请已取消"})
}
