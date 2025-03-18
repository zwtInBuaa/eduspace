import { deleteRequest, getRequest, postRequest, putRequest } from './request';

export function getAllExercise() {
  return getRequest('/questions/getall');
}

export function getAllExerciseSet() {
  return getRequest('/exams/getall');
}

export function createNewExercise(data) {
  return postRequest('/questions/addquestion', data);
}

export function editExercise(id, data) {
  return putRequest('/questions/' + id, data);
}

export function submitexerciseAnswer(id, data) {
  return postRequest('/questions/' + id + '/submit', data);
}

export function getExerciseDetail(id) {
  return getRequest('/questions/' + id);
}

export function getExerciseSetDetail(id) {
  return getRequest('/exams/' + id);
}

// export function getExerciseSetAllQuestions(id) {
//   return getRequest('/exams/' + id + '/questions');
// }

export function createNewExerciseSet(data) {
  return postRequest('/exams/addexam', data);
}

export function editExerciseSet(id, data) {
  return putRequest('/exams/' + id, data);
}

export function deleteExercise(id) {
  return deleteRequest('/questions/' + id);
}

export function deleteExerciseSet(id) {
  return deleteRequest('/exams/' + id);
}
