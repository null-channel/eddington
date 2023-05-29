export enum env {
  BACKEND_BASE_URL = import.meta.env.VITE_BASE_API_URL,
}
export { default as ROUTES } from './routes.const';
export { default as MAIN_ROUTES } from './dashbordRoutes.const';
export { default as NAVBAR } from './navbar.const';
