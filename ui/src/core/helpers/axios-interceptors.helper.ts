import { useCookies } from "vue3-cookies";
import { notify } from "@kyvg/vue3-notification";
import router from "@router";

const baseURL = import.meta.env.VITE_BASE_API_URL;

function axiosInterceptor(axios: any) {
  const { cookies } = useCookies();

  axios.interceptors.request.use((request: any) => {
    const token = cookies.get("user-token");
    const isApiUrl = request.url.startsWith(baseURL);
    if (token && isApiUrl) {
      request.headers.Authorization = `Bearer ${token}`;
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
        cookies.remove("user-token");
        router.push("/signin");
      }
      throw error;
    }
  );
  return axios;
}
export default axiosInterceptor