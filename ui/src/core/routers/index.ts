import { ROUTES } from "@constants";
import { createRouter, createWebHistory } from "vue-router";
import mainRouter from "./main.router";
import authGuard from "../guards/auth.guard";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      ...ROUTES.HOME,
      component: () => import("@pages/Home/home.vue"),
      //beforeEnter: authGuard,
    },
    {
      ...ROUTES.LOGIN,
      component: () => import("@pages/Log-in/log-in.vue"),
      // beforeEnter: authGuard,
    },
    {
      ...ROUTES.SIGNUP,
      component: () => import("@pages/Sign-up/sign-up.vue"),
      //beforeEnter: authGuard,
    },
    {
      ...ROUTES.MAIN,
      component: () => import("@pages/Dashboard/dashboard.vue"),
      beforeEnter: authGuard,
      children: mainRouter as any,
    },
  ],
});
export default router;
