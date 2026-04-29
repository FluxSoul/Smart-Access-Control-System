<template>
  <div class="p-6 max-w-6xl mx-auto">
    <h2 class="text-2xl font-bold mb-6">任务管理</h2>

    <!-- 任务列表 -->
    <el-table
      :data="tasks"
      v-loading="loading"
      border
      highlight-current-row
      @current-change="handleTaskSelect"
      class="mb-6"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="task_name" label="任务名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="cron_expr" label="Cron表达式" width="180">
        <template #default="{ row }">
          <el-tag type="info" size="small" class="font-mono">{{ row.cron_expr }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === true ? 'success' : 'info'">
            {{ row.status === true ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="params" label="参数" width="150">
        <template #default="{ row }">
          <el-popover placement="top" :width="250" trigger="hover" v-if="row.params">
            <pre class="text-xs m-0 p-2">{{ JSON.stringify(row.params, null, 2) }}</pre>
            <template #reference>
              <el-button link type="primary" size="small">查看</el-button>
            </template>
          </el-popover>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column prop="updated_at" label="更新时间" width="180" />
    </el-table>

    <!-- 编辑面板 -->
    <el-card v-if="selectedTask" class="mt-6">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-medium">编辑任务: {{ selectedTask.task_name }}</span>
          <el-tag type="info">ID: {{ selectedTask.id }}</el-tag>
        </div>
      </template>

      <el-form :model="editForm" label-width="120px" class="mt-4">
        <!-- Cron 表达式类型选择 -->
        <el-form-item label="表达式类型">
          <el-radio-group v-model="cronType">
            <el-radio label="standard">标准 Cron</el-radio>
            <el-radio label="every">@every 语法</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- Cron 表达式输入 -->
        <el-form-item label="Cron表达式" required>
          <el-input
            v-model="editForm.cron_expr"
            :placeholder="cronType === 'every' ? '例如: @every 5m' : '例如: 0 */5 * * * ?'"
            style="width: 300px;"
          />
          <el-popover placement="right" :width="300" trigger="hover" class="ml-2">
            <template #reference>
              <el-icon class="cursor-pointer text-gray-500"><InfoFilled /></el-icon>
            </template>
            <div class="text-sm">
              <div v-if="cronType === 'every'">
                <p>格式: @every &lt;数字&gt;&lt;单位&gt;</p>
                <p>单位: s(秒), m(分), h(时), d(天)</p>
                <p>示例: @every 30s, @every 2h</p>
              </div>
              <div v-else>
                <p>标准 Quartz Cron 格式</p>
                <p>示例: 0 0 12 * * ?</p>
              </div>
            </div>
          </el-popover>
        </el-form-item>

        <!-- 状态切换 -->
        <el-form-item label="任务状态">
          <el-switch
            v-model="editForm.status"
            :active-value="true"
            :inactive-value="false"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>

        <!-- 当前值预览 -->
        <el-descriptions :column="2" border class="mt-4">
          <el-descriptions-item label="原Cron">{{ selectedTask.cron_expr }}</el-descriptions-item>
          <el-descriptions-item label="原状态">
            <el-tag :type="selectedTask.status === true ? 'success' : 'info'">
              {{ selectedTask.status === true ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="新Cron">{{ editForm.cron_expr || '-' }}</el-descriptions-item>
          <el-descriptions-item label="新状态">
            <el-tag :type="editForm.status === true ? 'success' : 'info'">
              {{ editForm.status === true ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 操作按钮 -->
        <el-form-item class="mt-6">
          <el-button type="primary" @click="handleSave" :loading="saving" :icon="Check">
            保存修改
          </el-button>
          <el-button @click="handleCancel" :icon="Close">取消选择</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 无选择状态提示 -->
    <el-empty v-else description="请选择要编辑的任务" class="mt-6" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue';
import axios from '@/axios';
import { ElMessage } from 'element-plus';
import { InfoFilled, Check, Close } from '@element-plus/icons-vue';

// 类型定义
interface Task {
  id: number;
  task_name: string;
  cron_expr: string;
  description: string;
  status: boolean;
  params: Record<string, any>;
  created_at: string;
  updated_at: string;
}

// 响应式数据
const tasks = ref<Task[]>([]);
const loading = ref(false);
const selectedTask = ref<Task | null>(null);
const saving = ref(false);
const cronType = ref<'standard' | 'every'>('standard');

// 编辑表单
const editForm = reactive({
  cron_expr: '',
  status: true,
});

// 监听选中任务变化
watch(selectedTask, (newTask) => {
  if (newTask) {
    editForm.cron_expr = newTask.cron_expr;
    editForm.status = newTask.status;
  }
});

// 加载任务列表
const fetchTasks = async () => {
  loading.value = true;
  try {
    const response = await axios.get<{ data: Task[] }>('/task');
    tasks.value = response.data.data || response.data;
  } catch (error: any) {
    ElMessage.error(`加载失败: ${error.message || '未知错误'}`);
    tasks.value = [];
  } finally {
    loading.value = false;
  }
};

// 选择任务
const handleTaskSelect = (row: Task | null) => {
  selectedTask.value = row;
};

// 取消选择
const handleCancel = () => {
  selectedTask.value = null;
  editForm.cron_expr = '';
  editForm.status = true;
};

// 保存修改
const handleSave = async () => {
  if (!selectedTask.value) {
    ElMessage.warning('请先选择任务');
    return;
  }

  if (!editForm.cron_expr) {
    ElMessage.warning('请输入Cron表达式');
    return;
  }

  // 验证 Cron 格式
  if (cronType.value === 'every' && !validateEvery(editForm.cron_expr)) {
    ElMessage.error('@every 格式不正确，应为: @every <数字><s|m|h|d>');
    return;
  }

  saving.value = true;

  try {
    // 1. 更新 Cron 表达式（如果已修改）
    if (editForm.cron_expr !== selectedTask.value.cron_expr) {
      await axios.put(`/task/${selectedTask.value.task_name}/cron`, {
        cronExpr: editForm.cron_expr,
      });
      ElMessage.success('Cron表达式更新成功');
    }

    // 2. 更新状态（如果已修改）
    if (editForm.status !== selectedTask.value.status) {
      await axios.put(`/task/${selectedTask.value.task_name}/status`, {
        status: editForm.status === true,
      });
      ElMessage.success('状态更新成功');
    }

    // 3. 刷新列表
    await fetchTasks();

    // 4. 重新选中当前任务
    const updatedTask = tasks.value.find(t => t.id === selectedTask.value?.id);
    if (updatedTask) {
      selectedTask.value = updatedTask;
    }
  } catch (error: any) {
    ElMessage.error(`保存失败: ${error.response?.data?.error || error.message}`);
  } finally {
    saving.value = false;
  }
};

// 验证 @every 格式
const validateEvery = (value: string): boolean => {
  return /^@every \d+(s|m|h|d)$/.test(value);
};

// 初始化
onMounted(() => {
  fetchTasks();
});
</script>

<style scoped>
:deep(.el-table .current-row) {
  background-color: #f0f7ff;
}
</style>
