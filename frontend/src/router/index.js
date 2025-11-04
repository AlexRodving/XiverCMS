import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/Register.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/',
      component: () => import('../layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('../views/Dashboard.vue'),
        },
        {
          path: 'content-types',
          name: 'ContentTypes',
          component: () => import('../views/ContentTypes.vue'),
        },
        {
          path: 'content-types/:uid',
          name: 'ContentTypeDetail',
          component: () => import('../views/ContentTypeDetail.vue'),
        },
        {
          path: 'content-types/:uid/entries',
          name: 'ContentEntries',
          component: () => import('../views/ContentEntries.vue'),
        },
        {
          path: 'content-types/:uid/entries/:id',
          name: 'ContentEntryDetail',
          component: () => import('../views/ContentEntryDetail.vue'),
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('../views/Users.vue'),
          meta: { requiresAdmin: true },
        },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/')
  } else if (to.meta.requiresAdmin && !authStore.isSuperAdmin()) {
    next('/')
  } else {
    next()
  }
})

export default router

