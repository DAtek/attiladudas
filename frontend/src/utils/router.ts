import { apiClient } from '@/utils/api_client'
import { createRouter, createWebHistory } from 'vue-router'

export const routeNames = {
  ABOUT: 'about',
  ADMIN: 'admin',
  LOGIN: 'login',
  GALLERY_LIST: 'galleries',
  GALLERY: 'gallery',
  FIVE_IN_A_ROW: 'five_in_a_row'
}

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: routeNames.ABOUT,
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/admin/',
      name: routeNames.ADMIN,
      component: () => import('../views/AdminView.vue')
    },
    {
      path: '/login/',
      name: routeNames.LOGIN,
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/galleries/:slug/',
      name: routeNames.GALLERY,
      component: () => import('../views/gallery_view/GalleryView.vue')
    },
    {
      path: '/galleries/',
      name: routeNames.GALLERY_LIST,
      component: () => import('../views/GalleryList.vue')
    },
    {
      path: '/five-in-a-row/',
      name: routeNames.FIVE_IN_A_ROW,
      component: () => import('../views/five_in_a_row/FiveInARowView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.name === routeNames.ADMIN && !apiClient.token) {
    next({ name: routeNames.LOGIN })
    return
  }

  next()
})
