import { useCookies } from 'vue3-cookies';
import { ROUTES } from '@constants';

const AUTH_ROUTES = [ROUTES.LOGIN.name, ROUTES.RESET_PASSWORD.name, ROUTES.UPDATE_PASSWORD.name];

export default async (to: any) => {
  const { cookies } = useCookies();

  if (!AUTH_ROUTES.includes(to.name) && !cookies.get('user-token')) {
    return { name: ROUTES.LOGIN.name };
  }

  if (AUTH_ROUTES.includes(to.name) && cookies.get('user-token')) {
    return { name: ROUTES.HOME.name };
  }
  return true;
};
