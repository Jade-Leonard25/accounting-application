import { createRouter, createWebHashHistory } from 'vue-router'
import { authClient } from '../lib/auth-client'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    guestOnly?: boolean
  }
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    {
      path: '/login',
      name: 'login',
      component: () => import('../routers/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../routers/RegisterView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../routers/Dashboard.vue'),
      meta: { requiresAuth: true },
    },
    // 👇 Added Lazy Loaded Payroll Routes
    {
      path: '/payroll',
      name: 'payroll',
      component: () => import('../routers/PayrollView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/payroll/templates',
      name: 'payroll-templates',
      component: () => import('../routers/PayrollView.vue'),
      meta: { requiresAuth: true },
    },
    // Dev only, remove before shipping
    ...(import.meta.env.DEV
      ? [{ path: '/tester', name: 'tester', component: () => import('../routers/Tewster.vue') }]
      : []),
    { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
  ],
})

router.beforeEach(async (to) => {
  // Skip the network call for routes that don't care about auth
  if (!to.meta.requiresAuth && !to.meta.guestOnly) return

  let session = null
  try {
    const res = await authClient.getSession()
    session = res.data
  } catch {
    session = null
  }

  if (to.meta.requiresAuth && !session) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && session) {
    return { name: 'dashboard' }
  }
})

export default router