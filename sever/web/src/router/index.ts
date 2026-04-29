import {
    createRouter,
    createWebHashHistory,
    type RouteRecordRaw
} from 'vue-router'

import Index from '@/pages/index.vue'
import About from '@/pages/about.vue'
import Login from '@/pages/login.vue'
import NotFound from '@/pages/404.vue'
import EmqxDashboard from '@/pages/emqxDashboard.vue'
import register from '@/pages/register.vue'
import ChangeUserStatus from '@/pages/dashboard/changeUserStatus.vue'
import DashboardTemperature from '@/pages/dashboard/dashboardTemperature.vue'
import DashboardMoisture from '@/pages/dashboard/dashboardMoisture.vue'
import DashboardPPM from '@/pages/dashboard/dashboardPPM.vue'
import DashboardCron from '@/pages/dashboard/dashboardCron.vue'
import DashboardDoor from '@/pages/dashboard/dashboardDoor.vue'

const routes: RouteRecordRaw[] = [{
    path: '/',
    component: Index,
    meta: {
        title: '首页'
    }
},
{
    path: '/about',
    component: About,
    meta: {
        title: '关于'
    }
},
{
    path: '/:pathMatch(.*)*',
    component: NotFound,
    meta: {
        title: '404'
    }
},
{
    path: '/login',
    component: Login,
    meta: {
        title: '登录'
    }
},
{
    path: '/register',
    component: register,
    meta: {
        title: '注册'
    }

},
{
    path: '/dashboard',
    component: EmqxDashboard,
    children: [
        {
            path: 'temperature',
            component: DashboardTemperature,
            meta: {
                title: '温度视图'
            }
        },
        {
            path: 'moisture',
            component: DashboardMoisture,
            meta: {
                title: '湿度视图'
            }
        },
        {
            path: 'ppm',
            component: DashboardPPM,
            meta: {
                title: 'PPM'
            }
        },
        {
            path: "cron",
            component: DashboardCron,
            meta: {
                title: "Cron"
            }
        },
        {
            path: "door",
            component: DashboardDoor,
            meta: {
                title: "门禁"
            }
        },
        {
            path: 'userStatus',
            component: ChangeUserStatus,
            meta: {
                title: '用户管理'
            }
        }
    ]
}
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router


