<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s ease;
  overflow: hidden;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: white;
  font-size: 20px;
  font-weight: bold;
  background-color: #2b2f3a;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  border-bottom: 1px solid #eee;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.toggle-button {
  cursor: pointer;
  font-size: 24px;
  padding: 0 15px;
  display: flex;
}

.header-right {
  margin-right: 20px;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
}

.content-box {
  background-color: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>

<template>
  <el-container class="layout-container">
    <!-- 左侧导航栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo" v-if="!isCollapse">管理系统</div>
      <el-scrollbar>
        <el-menu
          :default-openeds="['1']"
          :collapse="isCollapse"
          :collapse-transition="false"
          router
        >
          <el-sub-menu index="1">
            <template #title>
              <el-icon><Odometer /></el-icon>
              <span v-show="!isCollapse">节点管理</span>
            </template>
            <el-menu-item index="/dashboard/temperature">
              <el-icon><Histogram /></el-icon>
              <span v-show="!isCollapse">温度</span>
            </el-menu-item>
            <el-menu-item index="/dashboard/moisture">
              <el-icon><TrendCharts /></el-icon>
              <span v-show="!isCollapse">湿度</span>
            </el-menu-item>
            <el-menu-item index="/dashboard/ppm">
              <el-icon><PictureRounded /></el-icon>
              <span v-show="!isCollapse">PPM</span>
            </el-menu-item>
            <el-menu-item index="/dashboard/cron">
              <el-icon><Setting /></el-icon>
              <span v-show="!isCollapse">任务</span>
            </el-menu-item>
              <el-menu-item index="/dashboard/door">
              <el-icon><Unlock /></el-icon>
              <span v-show="!isCollapse">远程开关</span>
            </el-menu-item>
            <!-- <el-menu-item-group>
              <template #title>
                <span v-show="!isCollapse">节点注册</span>
              </template>
              <el-menu-item index="/dashboard/page3">
                <el-icon><SetUp /></el-icon>
                <span v-show="!isCollapse">注册节点</span>
              </el-menu-item>
            </el-menu-item-group> -->
          </el-sub-menu>

          <el-sub-menu index="2">
            <template #title>
              <el-icon><Headset /></el-icon>
              <span v-show="!isCollapse">用户管理</span>
            </template>
            <el-menu-item index="/register">
              <el-icon><Refrigerator /></el-icon>
              <span v-show="!isCollapse">创建用户</span>
            </el-menu-item>
            <el-menu-item index="/dashboard/userStatus">
              <el-icon><Menu /></el-icon>
              <span v-show="!isCollapse">管理用户</span>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <!-- 右侧主内容区域 -->
    <el-container>
      <!-- 头部 -->
      <el-header class="header">
        <div class="toggle-button" @click="toggleCollapse">
          <el-icon>
            <component :is="iconComponent" />
          </el-icon>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ username }}
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <el-button @click="logout">退出登录</el-button>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主体内容 -->
      <el-main class="main-content">
        <div class="content-box">
            <router-view />
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import useCookies from 'universal-cookie'
import { useRouter } from 'vue-router'
import { Expand, Fold } from '@element-plus/icons-vue'

import {
  ArrowDown
} from '@element-plus/icons-vue'

// 控制导航栏折叠状态
const isCollapse = ref(false)
const cookies = new useCookies()
const username = cookies.get('username')
const router = useRouter()

function logout() {
  cookies.remove('username')
  cookies.remove('userId')
  cookies.remove('admin-token')
  router.push('/login')
}

// 切换导航栏折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const iconComponent = computed(() => {
  return isCollapse.value ? Expand : Fold
})

</script>

