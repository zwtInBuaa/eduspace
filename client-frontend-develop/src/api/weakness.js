import { getRequest } from './request';

export async function loadFailureRate(user_id) {
  return getRequest(`/users/${user_id}/weakness`);
}

export async function loadRecommendQuestions(user_id) {
  return getRequest(`/users/${user_id}/recQuestion`);
}

export async function loadQuestionOverview(user_id) {
  return getRequest(`/users/${user_id}/questionOverview`);
}
