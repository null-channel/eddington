import { MAIN_ROUTES } from '../constants';

export default [
  {
    ...MAIN_ROUTES.APPS,
    component: () =>import('@pages/container/container.vue'),
  },

];
