import axios from 'axios';
import { Message } from 'element-ui';
// 超时设定
axios.defaults.timeout = 60000 * 5;

axios.interceptors.request.use(config => {
    return config;
}, err => {
    Message.error({content: '请求超时',duration: 5});
    return Promise.resolve(err);
});

// http response 拦截器
axios.interceptors.response.use(response => {
    return response.data;
}, (err) => {
});

export const getRequest = (url, params) => {
    return axios({
        method: 'get',
        url: `${url}`,
        params: params,
    });
};

export const postRequest = (url, params) => {
    return axios({
        method: 'post',
        contentType: "application/json;charset=UTF-8",
        url: `${url}`,
        data: params,
    });
};