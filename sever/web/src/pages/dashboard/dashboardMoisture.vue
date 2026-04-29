<template>
  <div class="flex flex-col items-center p-4">
    <!-- 日期选择器 -->
    <div class="flex justify-center items-center mb-4 flex-wrap gap-2">
      <el-date-picker
          v-model="startTime"
          type="datetime"
          placeholder="开始时间"
          value-format="YYYY-MM-DD HH:mm:ss"
          class="w-48"
      />
      <span class="mx-2 text-gray-500">至</span>
      <el-date-picker
          v-model="endTime"
          type="datetime"
          placeholder="结束时间"
          value-format="YYYY-MM-DD HH:mm:ss"
          class="w-48"
      />
      <el-button
          @click="handleSearch"
          type="primary"
          class="ml-2"
          :loading="loading"
      >
        查询
      </el-button>
    </div>

    <!-- 图表容器 -->
    <div id="chart" class="w-full md:w-4/5 lg:w-3/4 xl:w-2/3 h-[50vh] min-h-[300px]"></div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import VChart from '@visactor/vchart';
import axios from '@/axios';

interface NodeData {
  nodeID: number;
  ts: string;
  value: number;
}

interface ChartData {
  id: number;
  time: string;
  value: number;
}

export default defineComponent({
  name: 'DashboardTemperature',
  setup() {
    const startTime = ref<string>('');
    const endTime = ref<string>('');
    const chartData = ref<ChartData[]>([]);
    const loading = ref<boolean>(false);

    let vchart: VChart | null = null;
    let resizeObserver: ResizeObserver | null = null;

    // 格式化时间戳（优化版）
    const formatDateTime = (timestamp: string): string => {
      try {
        // 处理 ISO 8601 格式，兼容时区
        const date = new Date(timestamp.endsWith('Z') ? timestamp : timestamp + 'Z');
        if (isNaN(date.getTime())) {
          throw new Error('Invalid date');
        }

        const year = date.getFullYear();
        const month = String(date.getUTCMonth() + 1).padStart(2, '0');
        const day = String(date.getUTCDate()).padStart(2, '0');
        const hour = String(date.getUTCHours()).padStart(2, '0');
        const minute = String(date.getUTCMinutes()).padStart(2, '0');

        return `${year}-${month}-${day} ${hour}:${minute}`;
      } catch (error) {
        console.error('时间格式化错误:', timestamp, error);
        return timestamp; // 返回原始值作为降级方案
      }
    };

    // 处理原始数据
    const processData = (data: NodeData[]): ChartData[] => {
      return data.map((item) => ({
        id: item.nodeID,
        time: formatDateTime(item.ts),
        value: Number(item.value) || 0, // 确保值为数字
      }));
    };

    // 初始化或更新图表
    const updateChart = (data: ChartData[]) => {
      if (!data || data.length === 0) {
        ElMessage.warning('暂无数据');
        return;
      }

      const spec = {
        type: 'line',
        data: {
          values: data,
        },
        xField: 'time',
        yField: 'value',
        seriesField: 'id',
        invalidType: 'link',
        line: {
          style: {
            curveType: 'monotone', // 平滑曲线
          },
        },
        axes: [
          {
            orient: 'bottom',
            type: 'band',
            label: {
              autoRotate: true,
              autoHide: true,
            },
          },
          {
            orient: 'left',
            type: 'linear',
            title: {
              visible: true,
              text: '湿度值',
            },
          },
        ],
        tooltip: {
          mark: {
            content: [
              {
                key: '时间',
                value: (datum: any) => datum.time,
              },
              {
                key: '节点',
                value: (datum: any) => `节点${datum.id}`,
              },
              {
                key: '湿度',
                value: (datum: any) => `${datum.value}°C`,
              },
            ],
          },
        },
        legend: {
          visible: true,
          title: {
            visible: true,
            text: '节点ID',
          },
        },
      };

      // 如果图表已存在，只更新数据；否则创建新实例
      if (vchart) {
        vchart.updateSpec(spec, true);
      } else {
        vchart = new VChart(spec, { dom: 'chart' });
        vchart.renderSync();
      }
    };

    // 加载数据（带日期过滤）
    const loadData = async () => {
      if (!startTime.value || !endTime.value) {
        ElMessage.warning('请选择完整的时间范围');
        return;
      }

      loading.value = true;
      try {
        const params: Record<string, string> = {
          startTime: startTime.value,
          endTime: endTime.value,
        };

        const res = await axios.get('/empx/getMessage/7', { params });

        if (!res.data || !Array.isArray(res.data)) {
          throw new Error('该时间段没有数据');
        }

        const processedData = processData(res.data);
        chartData.value = processedData;

        updateChart(processedData);
        ElMessage.success('数据加载成功');
      } catch (error: any) {
        console.error('数据加载失败:', error);
        ElMessage.error(error.message || '数据加载失败，请检查网络连接');
      } finally {
        loading.value = false;
      }
    };

    // 手动触发查询
    const handleSearch = () => {
      loadData();
    };

    // 监听容器大小变化
    const setupResizeObserver = () => {
      const chartContainer = document.getElementById('chart');
      if (chartContainer && 'ResizeObserver' in window) {
        resizeObserver = new ResizeObserver(() => {
          if (vchart) {
            vchart.renderSync();
          }
        });
        resizeObserver.observe(chartContainer);
      } else {
        // 降级方案：监听 window resize
        window.addEventListener('resize', () => {
          if (vchart) {
            vchart.renderSync();
          }
        });
      }
    };

    // 组件挂载时加载初始数据
    onMounted(() => {
      // 设置默认时间为最近24小时
      const end_for = new Date();
      const end = new Date(end_for.getTime() + 8 * 60 * 60 * 1000);
      const start = new Date(end.getTime() - 24 * 60 * 60 * 1000);

      startTime.value = formatDateTime(start.toISOString());
      endTime.value = formatDateTime(end.toISOString());

      loadData();
      setupResizeObserver();
    });

    // 组件卸载时清理资源
    onUnmounted(() => {
      if (vchart) {
        vchart.release();
        vchart = null;
      }
      if (resizeObserver) {
        resizeObserver.disconnect();
      }
    });

    return {
      startTime,
      endTime,
      loading,
      handleSearch,
    };
  },
});
</script>

<style scoped>
#chart {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>
