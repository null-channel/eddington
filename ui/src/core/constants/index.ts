export enum env {
  BACKEND_BASE_URL = import.meta.env.VITE_BASE_API_URL,
<<<<<<< HEAD
  ORY_URL = import.meta.env.VITE_ORY_URL || "http://localhost:4000",
=======
  ORY_URL = import.meta.env.VUE_APP_ORY_URL || "http://localhost:4000",
>>>>>>> 7470203 (adding ory plugin)
}
export { default as ROUTES } from "./routes.const";
export { default as MAIN_ROUTES } from "./dashbordRoutes.const";
export { NAVBAR_AFTER_LOGIN, NAVBAR_BEFORE_LOGIN } from "./navbar.const";
