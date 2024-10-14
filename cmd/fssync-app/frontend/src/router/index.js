import { createMemoryHistory, createRouter } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'server',
    component: () => import('@/views/Server.vue'),
    // component: () => import('@/views/Client.vue'),
    meta: {
      keepAlive: false
    }
  },
  {
    path: '/client',
    name: 'client',
    component: () => import('@/views/Client.vue'),
    meta: {
      keepAlive: false
    }
  }
];

const router = createRouter({
  history: createMemoryHistory(),
  // base: process.env.BASE_URL,
  routes,
})

export default router;
