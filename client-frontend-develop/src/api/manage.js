import { deleteRequest, getRequest, postRequest, putRequest } from './request';

export function getOneUser(id) {
  return getRequest('/users/' + id);
}

export function getOneCourse(course_id) {
  return getRequest('/courses/' + course_id);
}

/*ManageUsersView*/

export function getAllUsers() {
  return getRequest('/users/getall');
}

export function postUsers(user) {
  return postRequest('/users/signup', user);
}

export function deleteUser(id) {
  return deleteRequest('/users/' + id);
}

export function putUser(id, user) {
  return putRequest('/users/' + id, user);
}

export function resetPassword(data) {
  return postRequest('/users/resetPassword', data);
}

/*ManageCourseView*/
export function getAllCourses() {
  return getRequest('/courses/getall');
}

export function postCourse(data) {
  return postRequest('/courses/signup', data);
}

export function deleteCourse(id) {
  return deleteRequest('/courses/' + id);
}

export function putCourse(id, data) {
  return putRequest('/courses/' + id, data);
}

/*ManageCourseStudent 课程的全部学生*/
export function getCourseAllStudents(id) {
  return getRequest('/courses/' + id + '/students');
}

export function getCourseAllTeachers(course_id) {
  return getRequest('/courses/' + course_id + '/teachers');
}

export function postStudentsToCourse(courseId, data) {
  return postRequest('/courses/' + courseId + '/add_students', data);
}

export function deleteStudentFromCourse(course_id, student_id) {
  return deleteRequest('/courses/' + course_id + '/students/' + student_id);
}

export function postTeachersToCourse(course_id, data) {
  return postRequest('/courses/' + course_id + '/add_teachers', data);
}

export function deleteTeacherFromCourse(course_id, teacher_id) {
  return deleteRequest('/courses/' + course_id + '/teachers/' + teacher_id);
}

/*获取教师全部课程*/
export function getStudentsAllCourse(user_id) {
  return getRequest('/users/' + user_id + '/student_courses');
}

export function getUserAllCourses(user_id) {
  return getRequest('/users/' + user_id + '/student_courses');
}

export function userAllCourse(user_id) {
  let courses = [];
  getUserAllCourses(user_id)
    .then(function (response) {
      for (let i = 0; i < response.data.length; i++) {
        courses.push({
          id: response.data[i].id,
          name: response.data[i].name
        });
      }
    })
    .catch(function () {
      return courses;
    });
  return courses;
}
