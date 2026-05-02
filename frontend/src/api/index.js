import request from '@/utils/request'

export function login(data) {
  return request.post('/login', data)
}

export function getCurrentUser() {
  return request.get('/user/me')
}

export function updatePassword(data) {
  return request.put('/user/password', data)
}

export function updateProfile(data) {
  return request.put('/user/profile', data)
}

export function getAnnouncements(params) {
  return request.get('/announcements', { params })
}

export function getAnnouncementDetail(id) {
  return request.get(`/announcements/${id}`)
}

export function getRetakeCourses(params) {
  return request.get('/retake-courses', { params })
}

export function getMyGrades(params) {
  return request.get('/student/grades', { params })
}

export function getMyReports(params) {
  return request.get('/student/reports', { params })
}

export function createReport(data) {
  return request.post('/student/reports', data)
}

export function getMyApplications(params) {
  return request.get('/student/applications', { params })
}

export function createApplication(data) {
  return request.post('/student/applications', data)
}

export function cancelApplication(id) {
  return request.put(`/student/applications/${id}/cancel`)
}

export function getStudents(params) {
  return request.get('/teacher/students', { params })
}

export function getMyTeachingCourses(params) {
  return request.get('/teacher/teaching-courses', { params })
}

export function getCourseStudents(courseId) {
  return request.get(`/teacher/courses/${courseId}/students`)
}

export function enterGrade(id, data) {
  return request.put(`/teacher/grades/${id}`, data)
}

export function getTeacherGrades(params) {
  return request.get('/teacher/grades', { params })
}

export function getTeacherReports(params) {
  return request.get('/teacher/reports', { params })
}

export function replyReport(id, data) {
  return request.put(`/teacher/reports/${id}/reply`, data)
}

export function updateTeacherRetakeCourse(id, data) {
  return request.put(`/teacher/retake-courses/${id}`, data)
}

export function getTeacherApplications(params) {
  return request.get('/teacher/applications', { params })
}

export function reviewApplication(id, data) {
  return request.put(`/teacher/applications/${id}/review`, data)
}

export function getMajors(params) {
  return request.get('/admin/majors', { params })
}

export function createMajor(data) {
  return request.post('/admin/majors', data)
}

export function updateMajor(id, data) {
  return request.put(`/admin/majors/${id}`, data)
}

export function deleteMajor(id) {
  return request.delete(`/admin/majors/${id}`)
}

export function getClasses(params) {
  return request.get('/admin/classes', { params })
}

export function createClass(data) {
  return request.post('/admin/classes', data)
}

export function updateClass(id, data) {
  return request.put(`/admin/classes/${id}`, data)
}

export function deleteClass(id) {
  return request.delete(`/admin/classes/${id}`)
}

export function getAllUsers(params) {
  return request.get('/admin/users', { params })
}

export function createUser(data) {
  return request.post('/admin/users', data)
}

export function updateUser(id, data) {
  return request.put(`/admin/users/${id}`, data)
}

export function deleteUser(id) {
  return request.delete(`/admin/users/${id}`)
}

export function exportUsers() {
  return request.get('/admin/users/export', { responseType: 'blob' })
}

export function importUsers(formData) {
  return request.post('/admin/users/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getAllTeachers() {
  return request.get('/admin/teachers')
}

export function getAllAnnouncements(params) {
  return request.get('/admin/announcements', { params })
}

export function createAnnouncement(data) {
  return request.post('/admin/announcements', data)
}

export function updateAnnouncement(id, data) {
  return request.put(`/admin/announcements/${id}`, data)
}

export function deleteAnnouncement(id) {
  return request.delete(`/admin/announcements/${id}`)
}

export function getAllCourses(params) {
  return request.get('/admin/courses', { params })
}

export function createCourse(data) {
  return request.post('/admin/courses', data)
}

export function updateCourse(id, data) {
  return request.put(`/admin/courses/${id}`, data)
}

export function deleteCourse(id) {
  return request.delete(`/admin/courses/${id}`)
}

export function getAllRetakeCourses(params) {
  return request.get('/admin/retake-courses', { params })
}

export function createRetakeCourse(data) {
  return request.post('/admin/retake-courses', data)
}

export function updateRetakeCourse(id, data) {
  return request.put(`/admin/retake-courses/${id}`, data)
}

export function deleteRetakeCourse(id) {
  return request.delete(`/admin/retake-courses/${id}`)
}

export function getAllGrades(params) {
  return request.get('/admin/grades', { params })
}

export function updateGrade(id, data) {
  return request.put(`/admin/grades/${id}`, data)
}

export function getAllReports(params) {
  return request.get('/admin/reports', { params })
}

export function adminReplyReport(id, data) {
  return request.put(`/admin/reports/${id}/reply`, data)
}

export function deleteReport(id) {
  return request.delete(`/admin/reports/${id}`)
}

export function getAllApplications(params) {
  return request.get('/admin/applications', { params })
}

export function adminReviewApplication(id, data) {
  return request.put(`/admin/applications/${id}/review`, data)
}
