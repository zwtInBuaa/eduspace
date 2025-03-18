import axios from '@/api/intercept';
import store from '@/store';

export async function getRequest(url, params) {
  try {
    return await axios(url, {
      method: 'get',
      params: params
    });
  } catch (error) {
    store.dispatch('snackbar/error', error);
    return Promise.reject(error);
  }
}

export async function postRequest(url, data) {
  try {
    return await axios(url, {
      method: 'post',
      data: data
    });
  } catch (error) {
    store.dispatch('snackbar/error', error);
    return Promise.reject(error);
  }
}

export async function putRequest(url, data) {
  try {
    return await axios(url, {
      method: 'put',
      data: data
    });
  } catch (error) {
    store.dispatch('snackbar/error', error);
    return Promise.reject(error);
  }
}

export async function deleteRequest(url, data) {
  try {
    return await axios(url, {
      method: 'delete',
      data: data
    });
  } catch (error) {
    store.dispatch('snackbar/error', error);
    return Promise.reject(error);
  }
}
