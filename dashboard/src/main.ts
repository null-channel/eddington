import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { clerkPlugin } from 'vue-clerk/plugin'
import myhttp from '../src/axios/axios'

import App from './App.vue'
import router from './router'
import type { AxiosStatic } from 'axios'
declare module 'vue' {
  interface ComponentCustomProperties {
    $http: AxiosStatic
    $translate: (key: string) => string
  }
}
const app = createApp(App)

app.use(createPinia())

// TODO: change hard coded things to env variables.

// Add the Clerk plugin
app.use(clerkPlugin, {
  publishableKey: "pk_test_aW50aW1hdGUta3JpbGwtNzguY2xlcmsuYWNjb3VudHMuZGV2JA",
})

const http = myhttp;
// Set the base URL for axios
http.defaults.baseURL = 'http://localhost:8080';

console.log('Marek is unhappy');
// Add a response interceptor
http.interceptors.response.use(
  response => {
    console.log('Marek is something for sure');
    console.log(response);
    // If the response is successful, just return it
    if (response.status >= 200 && response.status < 233) {
      return response;
    }
    if (response.status === 401) {
      // If the response is 401, redirect to the login page
      router.push('/login');
    }
    if (response.status === 404) {
      // If the response is 404, redirect to the error page
      router.push('/error');
    }
    if (response.status == 234) {
      console.log('Marek is angry... this will work?');
      console.log(response.headers['location']);
      router.push(response.headers['Location']);
    }
    return response;
  },
  error => {
    console.log('Marek is happy?');
    if (error.response.status === 234) {

      // If the response is 500, redirect to the error page
      router.push(error.response.headers.location);

    }
    // Check if it's a redirect response (e.g., 302)
    if (error.response && [301, 302, 307, 308].includes(error.response.status)) {
      // Use the Vue Router to navigate to the new location
      // This assumes you have access to the router instance, possibly via the Composition API's useRouter hook
      router.push(error.response.headers.location);
    }
    return Promise.reject(error);
  }
);
app.provide('$http', http);

app.use(router)
// Mount the app

app.mount('#app')

