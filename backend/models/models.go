package models

import (
	"retakeCoursesManage/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Password     string    `gorm:"not null" json:"-"`
	RealName     string    `json:"real_name"`
	Role         Role      `gorm:"type:enum('student','teacher','admin')" json:"role"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	ClassID      *uint     `json:"class_id"`
	Class        *Class    `gorm:"foreignkey:ClassID" json:"class,omitempty"`
	MajorID      *uint     `json:"major_id"`
	Major        *Major    `gorm:"foreignkey:MajorID" json:"major,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Class struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	MajorID   uint      `json:"major_id"`
	Major     Major     `gorm:"foreignkey:MajorID" json:"major,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Major struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Code      string    `gorm:"unique;not null" json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Course struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Code      string    `gorm:"unique;not null" json:"code"`
	Credits   float64   `json:"credits"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RetakeCourse struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CourseID    uint      `json:"course_id"`
	Course      Course    `gorm:"foreignkey:CourseID" json:"course,omitempty"`
	TeacherID   uint      `json:"teacher_id"`
	Teacher     User      `gorm:"foreignkey:TeacherID" json:"teacher,omitempty"`
	Semester    string    `gorm:"not null" json:"semester"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Classroom   string    `json:"classroom"`
	MaxStudents int       `json:"max_students"`
	Status      string    `gorm:"default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Grade struct {
	ID             uint        `gorm:"primary_key" json:"id"`
	StudentID      uint        `json:"student_id"`
	Student        User        `gorm:"foreignkey:StudentID" json:"student,omitempty"`
	RetakeCourseID uint        `json:"retake_course_id"`
	RetakeCourse   RetakeCourse `gorm:"foreignkey:RetakeCourseID" json:"retake_course,omitempty"`
	Score          *float64    `json:"score"`
	Status         string      `gorm:"default:'pending'" json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type Announcement struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	IsTop     bool      `gorm:"default:false" json:"is_top"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Report struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Type      string    `gorm:"type:enum('suggestion','report')" json:"type"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignkey:UserID" json:"user,omitempty"`
	Reply     string    `gorm:"type:text" json:"reply"`
	ReplyUserID *uint    `json:"reply_user_id"`
	ReplyUser *User     `gorm:"foreignkey:ReplyUserID" json:"reply_user,omitempty"`
	Status    string    `gorm:"default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RetakeApplication struct {
	ID             uint         `gorm:"primary_key" json:"id"`
	StudentID      uint         `json:"student_id"`
	Student        User         `gorm:"foreignkey:StudentID" json:"student,omitempty"`
	RetakeCourseID uint         `json:"retake_course_id"`
	RetakeCourse   RetakeCourse `gorm:"foreignkey:RetakeCourseID" json:"retake_course,omitempty"`
	Reason         string       `gorm:"type:text" json:"reason"`
	Status         string       `gorm:"default:'pending'" json:"status"`
	ReviewerID     *uint        `json:"reviewer_id"`
	Reviewer       *User        `gorm:"foreignkey:ReviewerID" json:"reviewer,omitempty"`
	ReviewComment  string       `gorm:"type:text" json:"review_comment"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func InitDB() error {
	dsn := config.AppConfig.DBUser + ":" + config.AppConfig.DBPassword + "@tcp(" +
		config.AppConfig.DBHost + ":" + config.AppConfig.DBPort + ")/" +
		config.AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}

	DB.AutoMigrate(&User{}, &Class{}, &Major{}, &Course{}, &RetakeCourse{},
		&Grade{}, &Announcement{}, &Report{}, &RetakeApplication{})

	return nil
}
