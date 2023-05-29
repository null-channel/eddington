import { defineAsyncComponent } from 'vue';

export const lazyLoad = (loader: any) => {
  return defineAsyncComponent({
    loader,
  });
};
