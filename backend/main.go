package main

import (
	"fmt"
	"log"
	"retakeCoursesManage/config"
	"retakeCoursesManage/models"
	"retakeCoursesManage/routes"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer models.DB.Close()

	initDefaultData()

	r := routes.SetupRouter()

	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDefaultData() {
	var adminCount int
	models.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			RealName: "系统管理员",
			Role:     models.RoleAdmin,
			Email:    "admin@example.com",
		}
		models.DB.Create(&admin)
		fmt.Println("Default admin created: username=admin, password=admin123")
	}

	var majorCount int
	models.DB.Model(&models.Major{}).Count(&majorCount)
	if majorCount == 0 {
		majors := []models.Major{
			{Name: "计算机科学与技术", Code: "CS001"},
			{Name: "软件工程", Code: "SE001"},
			{Name: "电子信息工程", Code: "EE001"},
		}
		for _, m := range majors {
			models.DB.Create(&m)
		}
		fmt.Println("Default majors created")
	}

	var classCount int
	models.DB.Model(&models.Class{}).Count(&classCount)
	if classCount == 0 {
		var major models.Major
		models.DB.First(&major)
		if major.ID > 0 {
			classes := []models.Class{
				{Name: "计算机2101班", MajorID: major.ID},
				{Name: "计算机2102班", MajorID: major.ID},
			}
			for _, c := range classes {
				models.DB.Create(&c)
			}
			fmt.Println("Default classes created")
		}
	}

	var teacherCount int
	models.DB.Model(&models.User{}).Where("role = ?", "teacher").Count(&teacherCount)
	if teacherCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		teachers := []models.User{
			{
				Username: "teacher001",
				Password: string(hashedPassword),
				RealName: "张教授",
				Role:     models.RoleTeacher,
				Email:    "zhang@example.com",
			},
			{
				Username: "teacher002",
				Password: string(hashedPassword),
				RealName: "李教授",
				Role:     models.RoleTeacher,
				Email:    "li@example.com",
			},
		}
		for _, t := range teachers {
			models.DB.Create(&t)
		}
		fmt.Println("Default teachers created: username=teacher001/teacher002, password=123456")
	}

	var studentCount int
	models.DB.Model(&models.User{}).Where("role = ?", "student").Count(&studentCount)
	if studentCount == 0 {
		var class models.Class
		models.DB.First(&class)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		var classID *uint
		if class.ID > 0 {
			classID = &class.ID
		}

		students := []models.User{
			{
				Username: "student001",
				Password: string(hashedPassword),
				RealName: "王小明",
				Role:     models.RoleStudent,
				Email:    "wang@example.com",
				ClassID:  classID,
			},
			{
				Username: "student002",
				Password: string(hashedPassword),
				RealName: "李小红",
				Role:     models.RoleStudent,
				Email:    "li2@example.com",
				ClassID:  classID,
			},
		}
		for _, s := range students {
			models.DB.Create(&s)
		}
		fmt.Println("Default students created: username=student001/student002, password=123456")
	}

	var courseCount int
	models.DB.Model(&models.Course{}).Count(&courseCount)
	if courseCount == 0 {
		courses := []models.Course{
			{Name: "高等数学", Code: "MATH101", Credits: 5.0},
			{Name: "线性代数", Code: "MATH102", Credits: 3.0},
			{Name: "数据结构", Code: "CS101", Credits: 4.0},
			{Name: "操作系统", Code: "CS102", Credits: 4.0},
			{Name: "计算机网络", Code: "CS103", Credits: 3.0},
		}
		for _, c := range courses {
			models.DB.Create(&c)
		}
		fmt.Println("Default courses created")
	}
}
