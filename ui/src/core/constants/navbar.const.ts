import { MAIN_ROUTES, ROUTES } from ".";

export const NAVBAR_BEFORE_LOGIN = [
  { ...MAIN_ROUTES.CONTAINER },
  { ...ROUTES.CONTACT_US },
  { ...ROUTES.LOGIN, hiddenOnDesktop: true },
  { ...ROUTES.SIGNUP, hiddenOnDesktop: true },
];
export const NAVBAR_AFTER_LOGIN = [
  { ...MAIN_ROUTES.CONTAINER },
  { ...ROUTES.CONTACT_US },
];
