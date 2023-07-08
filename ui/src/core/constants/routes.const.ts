const ROUTES = {
  MAIN: { path: "/main", name: "dashboard_home" },
  NOT_FOUND: { path: "/:catchAll(.*)", name: "NotFound" },
  HOME: { name: "home", path: "/" },
  LOGIN: { name: "login", path: "/login" },
  SIGNUP: { name: "sign-up", path: "/sign-up" },
  RESET_PASSWORD: { name: "reset-password", path: "/resetpassword" },
  UPDATE_PASSWORD: { name: "update-password", path: "/updatepassword" },
};
export default ROUTES;
