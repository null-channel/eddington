import { ROUTES } from "@constants";
import { createRouter, createWebHistory } from "vue-router";

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
      // beforeEnter: authGuard,
    },
    // {
    //   ...ROUTES.CONTACT_US,
    //   component: () => lazyLoad(() => import("@pages/contactUs/contactUs.vue")),
    //   beforeEnter: authGuard,
    // },
    // {
    //   ...ROUTES.RESET_PASSWORD,
    //   component: () =>
    //     lazyLoad(() => import("@pages/resetPassword/resetPassword.vue")),
    //   beforeEnter: authGuard,
    // },
    // {
    //   ...ROUTES.UPDATE_PASSWORD,
    //   component: () =>
    //     lazyLoad(() => import("@pages/updatePassword/updatePassword.vue")),
    //   beforeEnter: authGuard,
    // },
    // {
    //   ...ROUTES.MAIN,
    //   component: () => lazyLoad(() => import("@pages/main/main.vue")),
    //   //beforeEnter: authGuard,
    //   children: mainRouter as any,
    // },
    // {
    //   ...ROUTES.NOT_FOUND,
    //   redirect: ROUTES.SIGN_IN.name,
    // },
  ],
});
export default router;
