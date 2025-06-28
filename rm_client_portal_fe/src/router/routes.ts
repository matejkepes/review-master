import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  // redirect to /dashboard if user is authenticated
  {
    path: '/',
    redirect: () => {
      if (localStorage.getItem('USER') === null) {
        return '/login'
      } else {
        return '/dashboard'
      }
    },
  },
  {
    path: '/dashboard',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/DashboardPage.vue') }],
    meta: { requiresAuth: true },
  },
  {
    path: '/reports',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/ReportsListPage.vue') }],
    meta: { requiresAuth: true },
  },
  {
    path: '/reports/:id',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/ReportViewPage.vue') }],
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    component: () => import('layouts/LandingLayout.vue'),
    children: [{ path: '', component: () => import('pages/LoginPage.vue') }],
    meta: { guestOnly: true },
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes
