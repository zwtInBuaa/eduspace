import { getRequest, postRequest, putRequest, deleteRequest } from './request';

export function getAllBlogs() {
  return getRequest('/posts/getall');
}

export function getBlogByID(id) {
  return getRequest('/posts/' + id);
}

export function postBlog(data) {
  return postRequest('/posts/postblog', data);
}

export function getBlogRemark(id) {
  return getRequest('/comments/getall/' + id);
}

export function postBlogRemark(data) {
  return postRequest('/comments/addcomment', data);
}

export async function postBlogEdit(id, title, content) {
  return await putRequest('/posts/' + id, { title, content });
}

export async function postBlogRemarkEdit(id, content) {
  return await putRequest('/comments/' + id, { content });
}

export async function postBlogDelete(id) {
  return deleteRequest('/posts/' + id);
}

export async function postBlogRemarkDelete(id) {
  return deleteRequest('/comments/' + id);
}
