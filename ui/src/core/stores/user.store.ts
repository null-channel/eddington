import { defineStore } from "pinia";
import { env } from "@constants";
import { useCookies } from "vue3-cookies";
import router from "@router";

const { cookies } = useCookies();
export const useUserStore = defineStore("user", {
  state: () => ({
    user: {  },
  }),
  getters: {},
  actions: {
    async login(login: { email: string; password: string }) {
      return window.$axios
        .post(`${env.BACKEND_BASE_URL}/auth/signin`, login)
        .then(async (data: any) => {
          cookies.set("user-token", data.access_token);
          this.$patch({ user: data.user });
          await router.push("/");
        });
    },
    logout() {
      cookies.remove("user-token");
      router.push("/signin");
    },
  },
});
