import axios from "@/axios";
import type { AxiosPromise } from "axios";

// 管理员登陆
function login(username: string, password: string): AxiosPromise<any> {
    return axios.post('/admin/login', {
        username,
        password
    })
}

// 管理员注册
function register(username: string, password: string): AxiosPromise<any> {
    return axios.post('/admin/register', {
        username,
        password
    })
}

// 获取管理员信息
function getUserInfo(): AxiosPromise<any> {
    return axios.get('/admin/getinfo')
}

export {
    login,
    register,
    getUserInfo
}