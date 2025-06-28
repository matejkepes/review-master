
const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/Index.vue') },
      { path: '/login', component: () => import('pages/Login.vue') },
      { path: '/clients', component: () => import('pages/Clients.vue'), meta: { requiresAuth: true } },
      { path: '/simpleclient/:clientIdProp/edit', props: true, component: () => import('pages/SimpleClient.vue'), meta: { requiresAuth: true } },
      { path: '/simpleclient/', component: () => import('pages/SimpleClient.vue'), meta: { requiresAuth: true } },
      { path: '/client/:clientIdProp/edit', props: true, component: () => import('pages/Client.vue'), meta: { requiresAuth: true } },
      { path: '/client/', component: () => import('pages/Client.vue'), meta: { requiresAuth: true } },
      { path: '/grconfig/:clientIdProp/add', props: true, component: () => import('pages/GRConfigNew.vue'), meta: { requiresAuth: true } },
      { path: '/grctime/:clientIdProp/:grctIdProp/add', props: true, component: () => import('pages/GRConfigTimeNew.vue'), meta: { requiresAuth: true } },
      { path: '/sendTest/', component: () => import('pages/SendTest.vue'), meta: { requiresAuth: true } },
      { path: '/stats/', component: () => import('pages/Stats.vue'), meta: { requiresAuth: true } },
      { path: '/statsnew/', component: () => import('pages/StatsNew.vue'), meta: { requiresAuth: true } },
      { path: '/checklog/', component: () => import('pages/CheckLog.vue'), meta: { requiresAuth: true } },
      { path: '/checknothingsent/', component: () => import('pages/CheckNothingSent.vue'), meta: { requiresAuth: true } },
      { path: '/gmbreport/', component: () => import('pages/GMBReport.vue'), meta: { requiresAuth: true } },
      { path: '/users', component: () => import('pages/Users.vue'), meta: { requiresAuth: true } },
      { path: '/user/', component: () => import('pages/User.vue'), meta: { requiresAuth: true } },
      { path: '/user/:userIdProp/edit', props: true, component: () => import('pages/User.vue'), meta: { requiresAuth: true } }
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/Error404.vue')
  }
]

export default routes
