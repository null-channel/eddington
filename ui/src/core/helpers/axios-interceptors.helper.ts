import { useCookies } from "vue3-cookies";
import { notify } from "@kyvg/vue3-notification";
import router from "@router";
import { AxiosInstance } from "axios";
import { useUserStore } from "../stores/user.store";

const baseURL = import.meta.env.VITE_BASE_API_URL;

function axiosInterceptor(axios: AxiosInstance) {
  const userStore =useUserStore()

  axios.interceptors.request.use((request: any) => {
    const isApiUrl = request.url.startsWith(baseURL);
    if (userStore.user.authenticated && isApiUrl) {
      request.withCredentials = true;
    }

    return request;
  });

  axios.interceptors.response.use(
    (res: any) => {
      return res.data;
    },
    (error: any) => {
      notify({ type: "error", text: error.response?.data.message });
      if (error.response?.status == 401) {
        userStore.$reset()
        router.push("/login");
      }
      throw error;
    }
  );
  return axios;
}
export default axiosInterceptor