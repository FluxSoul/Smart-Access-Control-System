<style scoped>

</style>
<template>
    <el-row class=" min-h-screen bg-light-blue-500">
        <el-col :lg="16" :md="12" class=" flex items-center justify-center">
            <div>
                <div class="font-bold text-5xl text-light-50 mb-4">
                    欢迎👏
                </div>
                <div class="text-gray-200 text-sm">
                    管理员注册
                </div>
            </div>
        </el-col>
        <el-col :lg="8" :md="12" class="bg-light-50 flex items-center justify-center flex-col">
            <h2 class="font-bold text-3xl text-gray-800">
                请注册
            </h2> 
            <div class="flex items-center justify-center my-5 text-gray-300 space-x-2">
                <span class="h-[1px] w-16 bg-gray-200"></span>
                <span>账号密码</span>
                <span class="h-[1px] w-16 bg-gray-200"></span>
            </div>
            <el-form ref="formRef" :model="form" class="w-[250px]" :rules="rules">
                <el-form-item prop="username">
                    <el-input class="my-1" v-model="form.username" type="username" placeholder="请输入用户名">
                        <template #prefix>
                            <el-icon>
                                <User/>
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input class="mb-1" v-model="form.password" type="password" placeholder="请输入密码" show-password>
                        <template #prefix>
                            <el-icon>
                                <Lock/> 
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password_confirm">
                    <el-input class="mb-1" v-model="form.password_confirm" type="password" placeholder="请输入确认密码" show-password>
                        <template #prefix>
                            <el-icon>
                                <Lock/> 
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item>
                    <el-button class="w-full" type="primary" @click="handleLogin" :loading="loading">登录</el-button>
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>

<script setup lang="ts">

import { ref, reactive } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import { useRouter } from 'vue-router'
import { register } from '@/api/manager'
import Cookie from 'universal-cookie'

const form = reactive({
    username: '',
    password: '',
    password_confirm: ''
})

const cookie = new Cookie()
const loading = ref(false)
const router = useRouter()

const rules = {
    username: [
        {required: true, message: '请输入用户名', trigger: 'blur'}
    ],
    password: [
        {required: true, message: '请输入密码', trigger: 'blur'},
        {min: 4, message: '密码长度不能小于4个字符', trigger: 'blur'}
    ]
}

const formRef = ref<InstanceType<typeof ElForm> | null>(null)

const handleLogin = () => {
    formRef.value?.validate((vailed: boolean) => {
        if (!vailed) {
            ElMessage.error('请输入正确的用户名和密码')
            return
        }
        if (form.password !== form.password_confirm) {
            ElMessage.error('密码不一致')
            return
        }
        // console.log(vailed)
        loading.value = true
        register(form.username, form.password)
        .then(res => {
            // console.log(res.data.user.token)
            // 提示成功
            cookie.remove('admin-token')
            ElMessage.success('注册成功, 1秒后自动跳转' + res.data.user.created_time)
            // 计时跳转
            setTimeout(() => {
                router.push('/login')
            }, 1000)
        })
        .catch(() => {
            ElMessage.error('注册失败')
        })
        .finally(() => {
            setTimeout(() => {
                loading.value = false
            }, 1000)
        })
    })
}

</script>