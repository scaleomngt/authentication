// 统一请求路径前缀在libs/axios.js中修改
import { getRequest, postRequest } from '../../utils/request';

// 创建账户
export const calcData = (params) => {
    return postRequest('/api/calcData', params)
}