export enum env {
  BACKEND_BASE_URL = import.meta.env.VITE_BASE_API_URL,
}
export { default as ROUTES } from './routes.const';
export { default as MAIN_ROUTES } from './dashbordRoutes.const';
export { NAVBAR_AFTER_LOGIN ,NAVBAR_BEFORE_LOGIN } from './navbar.const';
