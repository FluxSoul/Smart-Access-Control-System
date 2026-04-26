import axios, { type InternalAxiosRequestConfig } from "axios";
import { ElMessage } from "element-plus";
import useCookies from "universal-cookie";

const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 5000,
    headers: {
        "Content-Type": "application/json",
    },
});

service.interceptors.request.use(function (config: InternalAxiosRequestConfig) {
    const cookies = new useCookies()
    const token = cookies.get("admin-token")
    if (token) {
        config.headers["Authorization"] = token
    }
    return config
}, function (error) {
    return Promise.reject(error)
})

service.interceptors.response.use(function (response) {
    return response
}, function (error) {
    ElMessage.error(error.message)
    return Promise.reject(error)
})


export default service;
