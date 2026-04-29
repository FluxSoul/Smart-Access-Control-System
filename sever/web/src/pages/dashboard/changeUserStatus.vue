<template>
  <div class="user-status-container">
    <el-table 
      :data="users" 
      style="width: 100%" 
      v-loading="loading"
      element-loading-text="加载中..."
    >
      <el-table-column prop="id" label="用户ID" width="180" />
      <el-table-column prop="username" label="用户名" width="180" />
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-switch
            v-model="scope.row.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    
    <div class="toolbar">
      <el-button 
        type="primary" 
        @click="submitAllChanges" 
        :disabled="changedUsers.length === 0"
      >
        提交所有更改
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/axios'

interface User {
  id: number
  username: string
  status: number // 0: 禁用, 1: 启用
}

const users = ref<User[]>([])
const loading = ref(false)
const changedUsers = ref<User[]>([])

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await axios.get('/admin/getAllUser')
    users.value = response.data.users || []
  } catch (error) {
    ElMessage.error('获取用户列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 处理状态变更
const handleStatusChange = (user: User) => {
  // 检查是否已在变更列表中
  const index = changedUsers.value.findIndex(u => u.id === user.id)
  if (index > -1) {
    changedUsers.value[index] = user
  } else {
    changedUsers.value.push(user)
  }
}

// 提交所有更改
// 提交所有更改
const submitAllChanges = async () => {
  if (changedUsers.value.length === 0) {
    ElMessage.info('没有需要提交的更改')
    return
  }

  try {
    // 逐个发送用户状态更新请求
    const updatePromises = changedUsers.value.map(user => 
      axios.post('/admin/changeUserStatus', {
        id: user.id,
        status: user.status
      })
    )
    
    await Promise.all(updatePromises)
    
    ElMessage.success('所有更改已提交成功')
    changedUsers.value = []
    
    // 重新获取数据以确保一致性
    await fetchUsers()
  } catch (error) {
    ElMessage.error('提交更改失败')
    console.error(error)
  }
}

// 页面加载时自动获取数据
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-status-container {
  padding: 20px;
}

.toolbar {
  margin-top: 20px;
  text-align: right;
}
</style>