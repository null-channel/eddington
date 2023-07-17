export enum env {
  BACKEND_BASE_URL = import.meta.env.VITE_BASE_API_URL,
  ORY_URL = import.meta.env.VITE_ORY_URL || "https://auth.nullcloud.io",
}
export { default as ROUTES } from "./routes.const";
export { default as MAIN_ROUTES } from "./dashboardRoutes";
export { NAVBAR_AFTER_LOGIN, NAVBAR_BEFORE_LOGIN } from "./navbar.const";
