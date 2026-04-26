import router from '@/router'
import useCookies from 'universal-cookie'
import { ElMessage } from 'element-plus'

// 全局前置守卫
router.beforeEach((to, from, next) => { 
    const cookies = new useCookies()
    const token = cookies.get("admin-token")
    if (!token && to.path != '/login') {
        ElMessage.error('请先登录')
        return next('/login')
    }
    if (token && to.path == '/login') {
        ElMessage.error('请勿重复登录')
        return next( (from.path ? from.path : '/') )
    }

    let title = (to.meta.title ? to.meta.title : '后台管理') + ' - 节点设计'
    document.title = title

    next()
})
