package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"retakeCoursesManage/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetMajors(c *gin.Context) {
	var majors []models.Major
	query := models.DB

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Major{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&majors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取专业列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      majors,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateMajor(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var existing models.Major
	if err := models.DB.Where("name = ? OR code = ?", req.Name, req.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "专业名称或代码已存在"})
		return
	}

	major := models.Major{
		Name: req.Name,
		Code: req.Code,
	}

	if err := models.DB.Create(&major).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建专业失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "专业创建成功", "id": major.ID})
}

func UpdateMajor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var major models.Major
	if err := models.DB.First(&major, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "专业不存在"})
		return
	}

	if req.Name != "" {
		var existing models.Major
		if err := models.DB.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "专业名称已存在"})
			return
		}
		major.Name = req.Name
	}

	if req.Code != "" {
		var existing models.Major
		if err := models.DB.Where("code = ? AND id != ?", req.Code, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "专业代码已存在"})
			return
		}
		major.Code = req.Code
	}

	if err := models.DB.Save(&major).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新专业失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "专业更新成功"})
}

func DeleteMajor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	models.DB.Model(&models.Class{}).Where("major_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该专业下存在班级，无法删除"})
		return
	}

	if err := models.DB.Delete(&models.Major{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除专业失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "专业删除成功"})
}

func GetClasses(c *gin.Context) {
	var classes []models.Class
	query := models.DB.Preload("Major")

	if majorID := c.Query("major_id"); majorID != "" {
		query = query.Where("major_id = ?", majorID)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Class{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取班级列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      classes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateClass(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		MajorID uint   `json:"major_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var existing models.Class
	if err := models.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "班级名称已存在"})
		return
	}

	var major models.Major
	if err := models.DB.First(&major, req.MajorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "专业不存在"})
		return
	}

	class := models.Class{
		Name:    req.Name,
		MajorID: req.MajorID,
	}

	if err := models.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建班级失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "班级创建成功", "id": class.ID})
}

func UpdateClass(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Name    string `json:"name"`
		MajorID uint   `json:"major_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var class models.Class
	if err := models.DB.First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "班级不存在"})
		return
	}

	if req.Name != "" {
		var existing models.Class
		if err := models.DB.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "班级名称已存在"})
			return
		}
		class.Name = req.Name
	}

	if req.MajorID > 0 {
		var major models.Major
		if err := models.DB.First(&major, req.MajorID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "专业不存在"})
			return
		}
		class.MajorID = req.MajorID
	}

	if err := models.DB.Save(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新班级失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "班级更新成功"})
}

func DeleteClass(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	models.DB.Model(&models.User{}).Where("class_id = ? AND role = ?", id, "student").Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该班级下存在学生，无法删除"})
		return
	}

	if err := models.DB.Delete(&models.Class{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除班级失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "班级删除成功"})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	query := models.DB.Preload("Class").Preload("Major")

	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if classID := c.Query("class_id"); classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RealName string `json:"real_name"`
		Role     string `json:"role" binding:"required"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		ClassID  *uint  `json:"class_id"`
		MajorID  *uint  `json:"major_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var existing models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		RealName: req.RealName,
		Role:     models.Role(req.Role),
		Email:    req.Email,
		Phone:    req.Phone,
		ClassID:  req.ClassID,
		MajorID:  req.MajorID,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户创建成功", "id": user.ID})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		RealName string `json:"real_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		ClassID  *uint  `json:"class_id"`
		MajorID  *uint  `json:"major_id"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.ClassID != nil {
		user.ClassID = req.ClassID
	}
	if req.MajorID != nil {
		user.MajorID = req.MajorID
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := models.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

func ExportUsers(c *gin.Context) {
	var users []models.User
	if err := models.DB.Preload("Class").Preload("Major").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户失败"})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=users_%s.csv", time.Now().Format("20060102150405")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"ID", "用户名", "姓名", "角色", "邮箱", "电话", "班级", "专业", "创建时间"})

	for _, user := range users {
		roleName := map[string]string{
			"student": "学生",
			"teacher": "教师",
			"admin":   "管理员",
		}[string(user.Role)]

		className := ""
		if user.Class != nil {
			className = user.Class.Name
		}

		majorName := ""
		if user.Major != nil {
			majorName = user.Major.Name
		}

		writer.Write([]string{
			strconv.Itoa(int(user.ID)),
			user.Username,
			user.RealName,
			roleName,
			user.Email,
			user.Phone,
			className,
			majorName,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

func ImportUsers(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传文件"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件读取失败"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件格式错误"})
		return
	}

	defaultPassword := "123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)

	var successCount int
	var errorMessages []string

	for i, record := range records[1:] {
		if len(record) < 3 {
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行：数据不完整", i+2))
			continue
		}

		username := strings.TrimSpace(record[0])
		realName := strings.TrimSpace(record[1])
		roleStr := strings.TrimSpace(record[2])

		if username == "" {
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行：用户名为空", i+2))
			continue
		}

		var existing models.User
		if err := models.DB.Where("username = ?", username).First(&existing).Error; err == nil {
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行：用户名%s已存在", i+2, username))
			continue
		}

		role := models.RoleStudent
		switch strings.ToLower(roleStr) {
		case "教师", "teacher":
			role = models.RoleTeacher
		case "管理员", "admin":
			role = models.RoleAdmin
		}

		email := ""
		if len(record) > 3 {
			email = strings.TrimSpace(record[3])
		}

		phone := ""
		if len(record) > 4 {
			phone = strings.TrimSpace(record[4])
		}

		user := models.User{
			Username: username,
			Password: string(hashedPassword),
			RealName: realName,
			Role:     role,
			Email:    email,
			Phone:    phone,
		}

		if err := models.DB.Create(&user).Error; err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行：创建用户失败", i+2))
			continue
		}

		successCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("导入完成，成功%d条", successCount),
		"success_count": successCount,
		"errors":        errorMessages,
	})
}

func GetAllAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	query := models.DB.Preload("Author")

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Announcement{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      announcements,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateAnnouncement(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		IsTop   bool   `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	announcement := models.Announcement{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID.(uint),
		IsTop:    req.IsTop,
	}

	if err := models.DB.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "公告创建成功", "id": announcement.ID})
}

func UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		IsTop   *bool  `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var announcement models.Announcement
	if err := models.DB.First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	if req.Title != "" {
		announcement.Title = req.Title
	}
	if req.Content != "" {
		announcement.Content = req.Content
	}
	if req.IsTop != nil {
		announcement.IsTop = *req.IsTop
	}

	if err := models.DB.Save(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "公告更新成功"})
}

func DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := models.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "公告删除成功"})
}

func GetAllCourses(c *gin.Context) {
	var courses []models.Course
	query := models.DB

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.Course{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateCourse(c *gin.Context) {
	var req struct {
		Name    string  `json:"name" binding:"required"`
		Code    string  `json:"code" binding:"required"`
		Credits float64 `json:"credits"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var existing models.Course
	if err := models.DB.Where("name = ? OR code = ?", req.Name, req.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程名称或代码已存在"})
		return
	}

	course := models.Course{
		Name:    req.Name,
		Code:    req.Code,
		Credits: req.Credits,
	}

	if err := models.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程创建成功", "id": course.ID})
}

func UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Name    string  `json:"name"`
		Code    string  `json:"code"`
		Credits float64 `json:"credits"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var course models.Course
	if err := models.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	if req.Name != "" {
		var existing models.Course
		if err := models.DB.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "课程名称已存在"})
			return
		}
		course.Name = req.Name
	}

	if req.Code != "" {
		var existing models.Course
		if err := models.DB.Where("code = ? AND id != ?", req.Code, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "课程代码已存在"})
			return
		}
		course.Code = req.Code
	}

	if req.Credits != 0 {
		course.Credits = req.Credits
	}

	if err := models.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程更新成功"})
}

func DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	models.DB.Model(&models.RetakeCourse{}).Where("course_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该课程下存在重修课程，无法删除"})
		return
	}

	if err := models.DB.Delete(&models.Course{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}

func GetAllRetakeCourses(c *gin.Context) {
	var courses []models.RetakeCourse
	query := models.DB.Preload("Course").Preload("Teacher")

	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if teacherID := c.Query("teacher_id"); teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	models.DB.Model(&models.RetakeCourse{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取重修课程列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreateRetakeCourse(c *gin.Context) {
	var req struct {
		CourseID    uint   `json:"course_id" binding:"required"`
		TeacherID   uint   `json:"teacher_id" binding:"required"`
		Semester    string `json:"semester" binding:"required"`
		Classroom   string `json:"classroom"`
		MaxStudents int    `json:"max_students"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var course models.Course
	if err := models.DB.First(&course, req.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程不存在"})
		return
	}

	var teacher models.User
	if err := models.DB.Where("id = ? AND role = ?", req.TeacherID, "teacher").First(&teacher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教师不存在"})
		return
	}

	retakeCourse := models.RetakeCourse{
		CourseID:    req.CourseID,
		TeacherID:   req.TeacherID,
		Semester:    req.Semester,
		Classroom:   req.Classroom,
		MaxStudents: req.MaxStudents,
		Status:      "active",
	}

	if req.Status != "" {
		retakeCourse.Status = req.Status
	}

	if err := models.DB.Create(&retakeCourse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建重修课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重修课程创建成功", "id": retakeCourse.ID})
}

func AdminUpdateRetakeCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		CourseID    uint   `json:"course_id"`
		TeacherID   uint   `json:"teacher_id"`
		Semester    string `json:"semester"`
		Classroom   string `json:"classroom"`
		MaxStudents int    `json:"max_students"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var course models.RetakeCourse
	if err := models.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "重修课程不存在"})
		return
	}

	if req.CourseID > 0 {
		var foundCourse models.Course
		if err := models.DB.First(&foundCourse, req.CourseID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "课程不存在"})
			return
		}
		course.CourseID = req.CourseID
	}

	if req.TeacherID > 0 {
		var t models.User
		if err := models.DB.Where("id = ? AND role = ?", req.TeacherID, "teacher").First(&t).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "教师不存在"})
			return
		}
		course.TeacherID = req.TeacherID
	}

	if req.Semester != "" {
		course.Semester = req.Semester
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新重修课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重修课程更新成功"})
}

func DeleteRetakeCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	models.DB.Model(&models.Grade{}).Where("retake_course_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该重修课程下存在成绩记录，无法删除"})
		return
	}

	models.DB.Model(&models.RetakeApplication{}).Where("retake_course_id = ?", id).Delete(&models.RetakeApplication{})

	if err := models.DB.Delete(&models.RetakeCourse{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除重修课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重修课程删除成功"})
}

func GetAllGrades(c *gin.Context) {
	var grades []models.Grade
	query := models.DB.Preload("Student").Preload("RetakeCourse.Course")

	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if courseID := c.Query("retake_course_id"); courseID != "" {
		query = query.Where("retake_course_id = ?", courseID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成绩列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      grades,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateGrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Score  float64 `json:"score"`
		Status string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var grade models.Grade
	if err := models.DB.First(&grade, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "成绩记录不存在"})
		return
	}

	if req.Score > 0 || req.Score == 0 {
		if req.Score < 0 || req.Score > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "成绩应在0-100之间"})
			return
		}
		grade.Score = &req.Score
	}

	if req.Status != "" {
		grade.Status = req.Status
	}

	if err := models.DB.Save(&grade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新成绩失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成绩更新成功"})
}

func GetAllReports(c *gin.Context) {
	var reports []models.Report
	query := models.DB.Preload("User").Preload("ReplyUser")

	if reportType := c.Query("type"); reportType != "" {
		query = query.Where("type = ?", reportType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报告列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminReplyReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var req struct {
		Reply  string `json:"reply" binding:"required"`
		Status string `json:"status"`
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

	if req.Status != "" {
		report.Status = req.Status
	}

	if err := models.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "回复成功"})
}

func DeleteReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := models.DB.Delete(&models.Report{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "报告删除成功"})
}

func GetAllApplications(c *gin.Context) {
	var applications []models.RetakeApplication
	query := models.DB.Preload("Student").Preload("RetakeCourse.Course").Preload("Reviewer")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if courseID := c.Query("retake_course_id"); courseID != "" {
		query = query.Where("retake_course_id = ?", courseID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取申请列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      applications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminReviewApplication(c *gin.Context) {
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
	if err := models.DB.First(&application, appID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
		return
	}

	if application.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核待审核的申请"})
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

func GetAllTeachers(c *gin.Context) {
	var teachers []models.User
	query := models.DB.Where("role = ?", "teacher")

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教师列表失败"})
		return
	}

	c.JSON(http.StatusOK, teachers)
}
