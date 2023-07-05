import { ROUTES } from '@constants';
import { Ory } from '@helpers';

const AUTH_ROUTES = [ROUTES.LOGIN.name, ROUTES.RESET_PASSWORD.name, ROUTES.UPDATE_PASSWORD.name];

export default async (to: any) => {
  let authenticated = false;
  Ory.toSession()
    .then(() => {
      authenticated=true
    })
  if (!AUTH_ROUTES.includes(to.name) && !authenticated) {
    return { name: ROUTES.LOGIN.name };
  }

  if (AUTH_ROUTES.includes(to.name) && authenticated) {
    return { name: ROUTES.HOME.name };
  }
  return true;
};
