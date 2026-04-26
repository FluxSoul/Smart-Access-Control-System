<template>
  <div class="curtain-control-container">
    <el-card class="control-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><House /></el-icon>
          <span>远程控制</span>
        </div>
      </template>

      <!-- 节点选择 -->
      <div class="node-selector">
        <el-form :model="controlForm" label-width="80px">
          <el-form-item label="节点ID">
            <el-select
              v-model="controlForm.nodeId"
              placeholder="请选择节点"
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="node in nodeList"
                :key="node.id"
                :value="String(node.user_id)"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- 控制按钮 -->
      <div class="control-buttons">
        <el-button
          type="success"
          :icon="Open"
          @click="openCurtain"
          :loading="loading.open"
          :disabled="!controlForm.nodeId || loading.close"
          size="large"
          class="control-btn"
        >
          打开
        </el-button>

        <el-button
          type="danger"
          :icon="Close"
          @click="closeCurtain"
          :loading="loading.close"
          :disabled="!controlForm.nodeId || loading.open"
          size="large"
          class="control-btn"
        >
          关闭
        </el-button>
      </div>

      <!-- 状态显示 -->
      <div class="status-display" v-if="currentStatus">
        <el-tag
          :type="currentStatus === 'opened' ? 'success' : 'info'"
          size="large"
          class="status-tag"
        >
          <el-icon>
            <component :is="currentStatus === 'opened' ? Open : Close" />
          </el-icon>
          {{ currentStatus === 'opened' ? '已打开' : '已关闭' }}
        </el-tag>
        <p class="status-time" v-if="lastOperationTime">
          最后操作: {{ lastOperationTime }}
        </p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { House, Open, Close } from '@element-plus/icons-vue'
import axios from '@/axios'


// 接口定义
interface NodeOption {
  id: number,
  user_id: number,
}

interface ControlForm {
  nodeId: string
}


// 响应数据接口
interface ApiResponse {
  nodes: any
  message: string
  data?: any
}

// 组件状态
const controlForm = reactive<ControlForm>({
  nodeId: ''
})

const loading = reactive({
  open: false,
  close: false
})

const currentStatus = ref<string>('')
const lastOperationTime = ref<string>('')

// 节点列表（示例数据，实际应从API获取）
const nodeList = ref<NodeOption[]>()

const fetchNodeList = async () => {
  try {
    const response = await axios.get<ApiResponse>('/admin/getAllNode')
    console.log(response)
    if (response.status === 200 && Array.isArray(response.data.nodes)) {
      nodeList.value = response.data.nodes.map((node: NodeOption) => ({
        id: node.id + 1,
        user_id: node.user_id  // 添加label，如果没有则使用value
      }))

    } else {
        console.log(nodeList.value)
      ElMessage.error(response.data.message)
    }
  } catch (error) {
    ElMessage.error('获取节点列表失败')
  }
}



// 打开窗帘
const openCurtain = async () => {
  if (!controlForm.nodeId) {
    ElMessage.warning('请先选择节点')
    return
  }

  loading.open = true
  try {
    const response = await axios.get<ApiResponse>(
      `/empx/openTheDoor/${controlForm.nodeId}`
    )

    if (response.status === 200) {
      ElMessage.success('窗帘打开成功')
      currentStatus.value = 'opened'
      lastOperationTime.value = new Date().toLocaleString()
    } else {
      ElMessage.error(`操作失败: ${response.data.message}`)
    }
  } catch (error: any) {
    console.error('打开窗帘失败:', error)
    ElMessage.error(`打开窗帘失败: ${error.message || '网络错误'}`)
  } finally {
    loading.open = false
  }
}

// 关闭窗帘
const closeCurtain = async () => {
  if (!controlForm.nodeId) {
    ElMessage.warning('请先选择节点')
    return
  }

  loading.close = true
  try {
    const response = await axios.get<ApiResponse>(
      `/empx/closeTheDoor/${controlForm.nodeId}`
    )

    if (response.status === 200) {
      ElMessage.success('关闭成功')
      currentStatus.value = 'closed'
      lastOperationTime.value = new Date().toLocaleString()
    } else {
      ElMessage.error(`操作失败: ${response.data.message}`)
    }
  } catch (error: any) {
    console.error('关闭失败:', error)
    ElMessage.error(`关闭失败: ${error.message || '网络错误'}`)
  } finally {
    loading.close = false
  }
}

onMounted(() => {
    fetchNodeList()
})

</script>

<style scoped lang="scss">
.curtain-control-container {
  width: 100%;
  max-width: 500px;
  margin: 0 auto;

  .control-card {
    border-radius: 12px;
    border: 1px solid #ebeef5;

    .card-header {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 18px;
      font-weight: 600;
      color: #303133;

      .header-icon {
        font-size: 20px;
        color: #409eff;
      }
    }
  }

  .node-selector {
    margin-bottom: 30px;

    :deep(.el-form-item__label) {
      font-weight: 500;
      color: #606266;
    }
  }

  .control-buttons {
    display: flex;
    gap: 20px;
    justify-content: center;
    margin-bottom: 30px;

    .control-btn {
      min-width: 140px;
      height: 50px;
      font-size: 16px;
      font-weight: 500;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        transition: all 0.3s ease;
      }
    }
  }

  .status-display {
    text-align: center;
    padding-top: 20px;
    border-top: 1px dashed #e4e7ed;

    .status-tag {
      padding: 10px 20px;
      font-size: 16px;
      margin-bottom: 10px;

      .el-icon {
        margin-right: 6px;
      }
    }

    .status-time {
      margin: 0;
      font-size: 14px;
      color: #909399;
    }
  }
}

@media (max-width: 768px) {
  .curtain-control-container {
    max-width: 100%;

    .control-buttons {
      flex-direction: column;

      .control-btn {
        width: 100%;
      }
    }
  }
}
</style>
