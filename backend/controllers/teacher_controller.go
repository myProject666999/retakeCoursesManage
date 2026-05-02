package controllers

import (
	"net/http"
	"retakeCoursesManage/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStudents(c *gin.Context) {
	var students []models.User
	query := models.DB.Preload("Class").Preload("Major").Where("role = ?", "student")

	if classID := c.Query("class_id"); classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.User{}).Where("role = ?", "student").Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学生列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      students,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMyTeachingCourses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var courses []models.RetakeCourse

	query := models.DB.Preload("Course").Where("teacher_id = ?", userID)

	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取授课课程失败"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func GetCourseStudents(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var course models.RetakeCourse
	if err := models.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看该课程学生"})
		return
	}

	var grades []models.Grade
	if err := models.DB.Preload("Student.Class").Preload("Student.Major").
		Where("retake_course_id = ?", courseID).Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学生列表失败"})
		return
	}

	c.JSON(http.StatusOK, grades)
}

func EnterGrade(c *gin.Context) {
	userID, _ := c.Get("user_id")
	gradeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Score float64 `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var grade models.Grade
	if err := models.DB.Preload("RetakeCourse").First(&grade, gradeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "成绩记录不存在"})
		return
	}

	if grade.RetakeCourse.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限录入成绩"})
		return
	}

	if req.Score < 0 || req.Score > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "成绩应在0-100之间"})
		return
	}

	grade.Score = &req.Score
	grade.Status = "submitted"

	if err := models.DB.Save(&grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "录入成绩失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成绩录入成功"})
}

func GetMyGradesAsTeacher(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var grades []models.Grade

	query := models.DB.Preload("Student").Preload("RetakeCourse.Course")

	var courses []models.RetakeCourse
	models.DB.Where("teacher_id = ?", userID).Find(&courses)

	var courseIDs []uint
	for _, c := range courses {
		courseIDs = append(courseIDs, c.ID)
	}

	if len(courseIDs) > 0 {
		query = query.Where("retake_course_id IN ?", courseIDs)
	} else {
		query = query.Where("1 = 0")
	}

	if courseID := c.Query("course_id"); courseID != "" {
		query = query.Where("retake_course_id = ?", courseID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

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

func GetReportsForTeacher(c *gin.Context) {
	var reports []models.Report
	query := models.DB.Preload("User").Preload("ReplyUser")

	if reportType := c.Query("type"); reportType != "" {
		query = query.Where("type = ?", reportType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

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

func ReplyReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Reply string `json:"reply" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var report models.Report
	if err := models.DB.First(&report, reportID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报告不存在"})
		return
	}

	report.Reply = req.Reply
	report.ReplyUserID = &[]uint{userID.(uint)}[0]
	report.Status = "replied"

	if err := models.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "回复成功"})
}

func UpdateRetakeCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Classroom   string `json:"classroom"`
		MaxStudents int    `json:"max_students"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var course models.RetakeCourse
	if err := models.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改该课程"})
		return
	}

	if req.Classroom != "" {
		course.Classroom = req.Classroom
	}
	if req.MaxStudents > 0 {
		course.MaxStudents = req.MaxStudents
	}
	if req.Status != "" {
		course.Status = req.Status
	}

	if err := models.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程更新成功"})
}

func ReviewApplication(c *gin.Context) {
	userID, _ := c.Get("user_id")
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Status        string `json:"status" binding:"required"`
		ReviewComment string `json:"review_comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var application models.RetakeApplication
	if err := models.DB.Preload("RetakeCourse").First(&application, appID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
		return
	}

	if application.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核待审核的申请"})
		return
	}

	if application.RetakeCourse.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限审核该申请"})
		return
	}

	application.Status = req.Status
	application.ReviewerID = &[]uint{userID.(uint)}[0]
	application.ReviewComment = req.ReviewComment

	if err := models.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	if req.Status == "approved" {
		grade := models.Grade{
			StudentID:      application.StudentID,
			RetakeCourseID: application.RetakeCourseID,
			Status:         "pending",
		}
		models.DB.Create(&grade)
	}

	c.JSON(http.StatusOK, gin.H{"message": "审核完成"})
}

func GetApplicationsForTeacher(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var applications []models.RetakeApplication

	var courses []models.RetakeCourse
	models.DB.Where("teacher_id = ?", userID).Find(&courses)

	var courseIDs []uint
	for _, c := range courses {
		courseIDs = append(courseIDs, c.ID)
	}

	query := models.DB.Preload("Student").Preload("RetakeCourse.Course").Preload("Reviewer")

	if len(courseIDs) > 0 {
		query = query.Where("retake_course_id IN ?", courseIDs)
	} else {
		query = query.Where("1 = 0")
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

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
