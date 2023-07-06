import { createRouter, createWebHistory } from 'vue-router';
import AdminGameView from '../views/AdminGameView.vue';
import AdminView from '../views/AdminView.vue';
import GameView from '../views/GameView.vue';
import HomeView from '../views/HomeView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/game',
      name: 'game',
      component: GameView,
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminView,
    },
    {
      path: '/admin/games/:gameId',
      name: 'admin-game',
      component: AdminGameView,
    },
  ],
});

export default router;
