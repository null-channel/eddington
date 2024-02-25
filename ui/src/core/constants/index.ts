export enum env {
  BACKEND_BASE_URL = import.meta.env.VITE_BASE_API_URL,
  CLERK_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY,
  CLERK_URL = import.meta.env.VITE_CLERK_URL || "https://nullcloud.clerk.dev",
}
export { default as ROUTES } from "./routes.const";
export { default as MAIN_ROUTES } from "./dashboardRoutes";
export { NAVBAR_AFTER_LOGIN, NAVBAR_BEFORE_LOGIN } from "./navbar.const";
