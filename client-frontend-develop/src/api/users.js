import { postRequest, putRequest } from './request';

export async function login(data) {
  return postRequest('/users/login', data);
}

export async function logout() {
  return postRequest('/users/logout');
}

/*用户中心*/
export async function changeAvatar(id, data) {
  return postRequest('/users/upload-avatar/' + id, data);
}

export async function changePassword(id, data) {
  return putRequest('/users/' + id + '/changePassword', data);
}
