import { MAIN_ROUTES, ROUTES } from ".";

export const NAVBAR_BEFORE_LOGIN = [
  { ...ROUTES.HOME },
  { ...ROUTES.LOGIN, hiddenOnDesktop: true },
  { ...ROUTES.SIGNUP, hiddenOnDesktop: true },
];
export const NAVBAR_AFTER_LOGIN = [
  { ...ROUTES.HOME },
  {
    ...MAIN_ROUTES.CONTAINER,
    path: `${ROUTES.MAIN.path}/${MAIN_ROUTES.CONTAINER.path}`,
  },
];
