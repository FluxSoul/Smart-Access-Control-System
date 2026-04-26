<template>
  <div class="homepage-container">
    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="nav-content">
        <div class="logo-section">
            <a href="https://www.cqupt.edu.cn" target="_blank">
                <img src="https://www.cqupt.edu.cn/images/logo1.png" class="h-7" alt="重庆邮电大学">
            </a>
            <el-icon class="logo-icon"><Connection /></el-icon>
            <h1 class="title">IoT Node Management</h1>
        </div>
        <el-button
          type="primary"
          class="login-btn"
          @click="goToLogin"
        >
          <el-icon class="btn-icon"><User /></el-icon>
          登录
        </el-button>
      </div>
    </header>

    <!-- 面包屑导航 -->
    <div class="breadcrumb-container">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/dashboard'}">管理平台</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 主要内容区域 -->
    <main class="main-content">
      <div class="content-wrapper">
        <!-- 项目介绍卡片 -->
        <el-card class="hero-card">
          <h2 class="hero-title">重庆邮电大学物联网工程·节点设计</h2>
          <p class="hero-subtitle">多节点智能数据收发管理系统</p>
          <div class="tech-stack-tags">
            <el-tag type="success">MQTT</el-tag>
            <el-tag type="info">EMQX</el-tag>
            <el-tag type="warning">Gin</el-tag>
            <el-tag type="primary">Vue3</el-tag>
            <el-tag type="info">ElementPlus</el-tag>
            <el-tag type="success">TypeScript</el-tag>
            <el-tag type="danger">STM32</el-tag>
            <el-tag type="info">ESP32+8266</el-tag>
          </div>
        </el-card>

        <!-- 功能特性网格 -->
        <el-row :gutter="24" class="features-grid">
          <el-col :span="8">
            <el-card class="feature-card" shadow="hover">
              <div class="feature-icon">
                <el-icon class="icon"><Odometer /></el-icon>
              </div>
              <h3>实时监控</h3>
              <p>支持多节点传感器数据实时采集与可视化展示</p>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="feature-card" shadow="hover">
              <div class="feature-icon">
                <el-icon class="icon"><Upload /></el-icon>
              </div>
              <h3>智能收发</h3>
              <p>基于MQTT协议实现高效、稳定的数据通信</p>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="feature-card" shadow="hover">
              <div class="feature-icon">
                <el-icon class="icon"><Operation /></el-icon>
              </div>
              <h3>节点管理</h3>
              <p>STM32节点远程配置、状态监控</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <p>© 2025 重庆邮电大学 | 物联网工程节点设计系统</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import {
  Connection,
  User,
  Odometer,
  Upload,
  Operation
} from '@element-plus/icons-vue';

const router = useRouter();

const goToLogin = () => {
  router.push('/login');
  ElMessage.success('即将跳转到登录页面');
};
</script>

<style scoped lang="scss">
.homepage-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

// 顶部导航栏
.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);

  .nav-content {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
  }

  .logo-section {
    display: flex;
    align-items: center;
    gap: 12px;

    .logo-icon {
      font-size: 28px;
      color: #409eff;
    }

    .title {
      font-size: 20px;
      font-weight: 600;
      margin: 0;
      color: #303133;
    }
  }

  .login-btn {
    font-weight: 500;

    .btn-icon {
      margin-right: 6px;
    }
  }
}

// 面包屑容器
.breadcrumb-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.08);

  :deep(.el-breadcrumb) {
    .el-breadcrumb__item {
      .el-breadcrumb__inner {
        color: rgba(255, 255, 255, 0.8);

        &:hover {
          color: #fff;
        }
      }

      &:last-child .el-breadcrumb__inner {
        color: #fff;
        font-weight: normal;
      }
    }
  }
}

// 主要内容区域
.main-content {
  padding: 40px 24px;
  min-height: 600px;

  .content-wrapper {
    max-width: 1200px;
    margin: 0 auto;
  }
}

// 主卡片样式
.hero-card {
  border: none;
  background: rgba(255, 255, 255, 0.95);
  margin-bottom: 40px;
  padding: 40px 32px;
  text-align: center;
  border-radius: 12px;

  .hero-title {
    font-size: 32px;
    font-weight: 700;
    margin: 0 0 12px 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text; /* 先写厂商前缀版本 */
    background-clip: text;         /* 再写标准版本 */
    -webkit-text-fill-color: transparent;
  }

  .hero-subtitle {
    font-size: 18px;
    color: #606266;
    margin: 0 0 24px 0;
  }

  .tech-stack-tags {
    display: flex;
    justify-content: center;
    gap: 12px;
    flex-wrap: wrap;

    .el-tag {
      font-weight: 500;
      border-radius: 20px;
      padding: 0 16px;
    }
  }
}

// 功能特性网格
.features-grid {
  margin-bottom: 60px;

  .feature-card {
    border: none;
    background: rgba(255, 255, 255, 0.95);
    padding: 32px 24px;
    text-align: center;
    border-radius: 12px;
    height: 100%;
    transition: transform 0.3s ease;

    &:hover {
      transform: translateY(-4px);
    }

    .feature-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 60px;
      height: 60px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 50%;
      margin-bottom: 20px;

      .icon {
        font-size: 28px;
        color: #fff;
      }
    }

    h3 {
      font-size: 20px;
      color: #303133;
      margin: 0 0 12px 0;
    }

    p {
      color: #606266;
      margin: 0;
      line-height: 1.8;
    }
  }

  @media (max-width: 768px) {
    .el-col {
      margin-bottom: 24px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}

// 页脚
.footer {
  background: rgba(0, 0, 0, 0.2);
  padding: 24px;
  text-align: center;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}
</style>
