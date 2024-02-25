import { createApp } from "vue";
import "./index.css";
import App from "./App.vue";
import { axiosHelper } from "@helpers";
import { createPinia } from "pinia";
import { plugin, defaultConfig } from "@formkit/vue";
import genesis from "./formkit.config";
import Notifications from "@kyvg/vue3-notification";
import router from "@router";
import { clerkPlugin } from 'vue-clerk/plugin'

async function run(){
    const pinia = createPinia();
    const app = createApp(App);
    app.use(pinia);
    app.use(clerkPlugin, {
        publishableKey: import.meta.env.VITE_CLERK_PUBLISHABLE_KEY,});
    app.use(plugin, defaultConfig(genesis));
    app.use(axiosHelper);
    app.use(router);
    app.use(Notifications);
    app.mount("#app");
}
run()
