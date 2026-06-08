import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue'), meta: { public: true } },
  { path: '/', name: 'Dashboard', component: () => import('../views/DashboardView.vue') },
  { path: '/sales', name: 'Sales', component: () => import('../views/SalesView.vue') },
  { path: '/reports', name: 'Reports', component: () => import('../views/ReportsView.vue') },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
