import { MAIN_ROUTES } from '@constants';
import { lazyLoad } from '@helpers';

export default [
  {
    ...MAIN_ROUTES.CONTAINER,
    component: () =>
      lazyLoad(() => import('@pages/container/container.vue')),
  },

];
