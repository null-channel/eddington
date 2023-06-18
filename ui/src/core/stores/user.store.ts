import { defineStore } from "pinia";
import { env } from "@constants";
import { useCookies } from "vue3-cookies";
import router from "@router";
import { $ory, injectStrict } from "@helpers";

const { cookies } = useCookies();
export const useUserStore = defineStore("user", {
  state: () => ({
    user: {
      session: {},
      authenticated: false,
      logoutUrl: "",
    },
  }),
  getters: {
    
  },
  actions: {
    getUser() {},
    login(url: string, headers: any, formData: any) {
      return window.$axios
        .post(url, formData, {
          headers,
        })
    },
    logout() {
      cookies.remove("user-token");
      router.push("/signin");
    },
 
  },
});
