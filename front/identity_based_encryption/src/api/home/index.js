// 统一请求路径前缀在libs/axios.js中修改
import { getRequest, postRequest } from '../../utils/request';

// 创建账户
export const submitData = (params) => {
    return postRequest('/api/submitData', params)
}
// 创建账户
export const createAccount = (params) => {
    return getRequest('/api/createAccount', params)
}