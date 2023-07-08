import { MAIN_ROUTES } from '../constants';

export default [
  {
    ...MAIN_ROUTES.CONTAINER,
    component: () =>import('@pages/container/container.vue'),
  },

];
