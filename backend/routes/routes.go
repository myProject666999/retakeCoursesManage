package routes

import (
	"retakeCoursesManage/controllers"
	"retakeCoursesManage/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/user/me", controllers.GetCurrentUser)
			auth.PUT("/user/password", controllers.UpdatePassword)
			auth.PUT("/user/profile", controllers.UpdateProfile)

			auth.GET("/announcements", controllers.GetAnnouncements)
			auth.GET("/announcements/:id", controllers.GetAnnouncementDetail)

			auth.GET("/retake-courses", controllers.GetRetakeCourses)

			student := auth.Group("")
			student.Use(middleware.RoleMiddleware("student"))
			{
				student.GET("/student/grades", controllers.GetMyGrades)
				student.GET("/student/reports", controllers.GetMyReports)
				student.POST("/student/reports", controllers.CreateReport)
				student.GET("/student/applications", controllers.GetMyApplications)
				student.POST("/student/applications", controllers.CreateApplication)
				student.PUT("/student/applications/:id/cancel", controllers.CancelApplication)
			}

			teacher := auth.Group("")
			teacher.Use(middleware.RoleMiddleware("teacher"))
			{
				teacher.GET("/teacher/students", controllers.GetStudents)
				teacher.GET("/teacher/teaching-courses", controllers.GetMyTeachingCourses)
				teacher.GET("/teacher/courses/:course_id/students", controllers.GetCourseStudents)
				teacher.PUT("/teacher/grades/:id", controllers.EnterGrade)
				teacher.GET("/teacher/grades", controllers.GetMyGradesAsTeacher)
				teacher.GET("/teacher/reports", controllers.GetReportsForTeacher)
				teacher.PUT("/teacher/reports/:id/reply", controllers.ReplyReport)
				teacher.PUT("/teacher/retake-courses/:id", controllers.UpdateRetakeCourse)
				teacher.GET("/teacher/applications", controllers.GetApplicationsForTeacher)
				teacher.PUT("/teacher/applications/:id/review", controllers.ReviewApplication)
			}

			admin := auth.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/admin/majors", controllers.GetMajors)
				admin.POST("/admin/majors", controllers.CreateMajor)
				admin.PUT("/admin/majors/:id", controllers.UpdateMajor)
				admin.DELETE("/admin/majors/:id", controllers.DeleteMajor)

				admin.GET("/admin/classes", controllers.GetClasses)
				admin.POST("/admin/classes", controllers.CreateClass)
				admin.PUT("/admin/classes/:id", controllers.UpdateClass)
				admin.DELETE("/admin/classes/:id", controllers.DeleteClass)

				admin.GET("/admin/users", controllers.GetAllUsers)
				admin.POST("/admin/users", controllers.CreateUser)
				admin.PUT("/admin/users/:id", controllers.UpdateUser)
				admin.DELETE("/admin/users/:id", controllers.DeleteUser)
				admin.GET("/admin/users/export", controllers.ExportUsers)
				admin.POST("/admin/users/import", controllers.ImportUsers)
				admin.GET("/admin/teachers", controllers.GetAllTeachers)

				admin.GET("/admin/announcements", controllers.GetAllAnnouncements)
				admin.POST("/admin/announcements", controllers.CreateAnnouncement)
				admin.PUT("/admin/announcements/:id", controllers.UpdateAnnouncement)
				admin.DELETE("/admin/announcements/:id", controllers.DeleteAnnouncement)

				admin.GET("/admin/courses", controllers.GetAllCourses)
				admin.POST("/admin/courses", controllers.CreateCourse)
				admin.PUT("/admin/courses/:id", controllers.UpdateCourse)
				admin.DELETE("/admin/courses/:id", controllers.DeleteCourse)

				admin.GET("/admin/retake-courses", controllers.GetAllRetakeCourses)
				admin.POST("/admin/retake-courses", controllers.CreateRetakeCourse)
				admin.PUT("/admin/retake-courses/:id", controllers.AdminUpdateRetakeCourse)
				admin.DELETE("/admin/retake-courses/:id", controllers.DeleteRetakeCourse)

				admin.GET("/admin/grades", controllers.GetAllGrades)
				admin.PUT("/admin/grades/:id", controllers.UpdateGrade)

				admin.GET("/admin/reports", controllers.GetAllReports)
				admin.PUT("/admin/reports/:id/reply", controllers.AdminReplyReport)
				admin.DELETE("/admin/reports/:id", controllers.DeleteReport)

				admin.GET("/admin/applications", controllers.GetAllApplications)
				admin.PUT("/admin/applications/:id/review", controllers.AdminReviewApplication)
			}
		}
	}

	return r
}
