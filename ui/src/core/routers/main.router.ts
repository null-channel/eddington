import { MAIN_ROUTES } from '@constants';
import { lazyLoad } from '@helpers';

export default [
  {
    ...MAIN_ROUTES.DASHBORD,
    component: () =>
      lazyLoad(() => import('@pages/main/dashbord/dashbord.vue')),
  },

];
