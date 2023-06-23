import { defineStore } from "pinia";
import { useCookies } from "vue3-cookies";
import router from "@router";

const { cookies } = useCookies();
const useUserStore = defineStore("user", {
  state: () => ({
    user: {
      session: {},
      authenticated: false,
      logoutUrl: "",
    },
  }),
  getters: {},
  actions: {
    getUser() {},
    login(url: string, headers: any, formData: any) {
      return window.$axios.post(url, formData, {
        headers,
      });
    },
    signUp(url: string, headers: any, formData: any) {
      return window.$axios.post(url, formData, {
        headers,
      });
    },
    logout() {
      cookies.remove("user-token");
      router.push("/signin");
    },
  },
});
export default useUserStore