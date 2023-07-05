import { defineStore } from "pinia";
import router from "@router";

const useUserStore = defineStore("user", {
  state: () => ({
  }),
  getters: {
    user() {
      const user = localStorage.getItem("session");
      return user ? JSON.parse(user) : null;
    },
  },
  actions: {
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
    logout(logoutUrl: string) {
      window.$axios
        .get(logoutUrl, {
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
          },
        })
        .then(() => {
          localStorage.removeItem("session");
          router.push("/login");
        });
    },
  },
});
export default useUserStore;
