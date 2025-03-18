import axios from 'axios';
import store from '@/store';
import router from '@/router';

var HttpRequestUrl = 'http://114.116.211.180:8081/'; // 发布的url

axios.interceptors.request.use(
  (request) => {
    request.baseURL = HttpRequestUrl;
    let token = store.getters['user/token'];
    if (token !== null && token !== '') {
      request.headers['Authorization'] = 'Bearer ' + token;
    }
    return request;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          store.commit('user/setToken', '');
          if (!router.currentRoute.fullPath.includes('/login')) {
            await router.push({
              path: '/login',
              query: {
                redirect: router.currentRoute.fullPath
              }
            });
          }
          return Promise.reject('登录状态失效');
        case 403:
          return Promise.reject('拒绝访问');
        case 404:
          return Promise.reject('请求地址出错');
        case 429:
          return Promise.reject('请求冷却中，请稍等 5 秒后重试');
        case 500:
          return error.response.data === null
            ? Promise.reject('服务器错误')
            : Promise.reject(error.response.status + '：' + error.response.data.error);
        default:
          return error.response.data === null
            ? Promise.reject('请求错误：' + error.response.status)
            : Promise.reject(error.response.status + '：' + error.response.data.error);
      }
    }
    return Promise.reject(error);
  }
);

export default axios;
