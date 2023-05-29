import { createApp } from "vue";
import "./index.css";
import App from "./App.vue";
import { axiosHelper } from "@helpers";
import { createPinia } from "pinia";
import { plugin, defaultConfig } from "@formkit/vue";
import genesis from "./formkit.config";
import Notifications from "@kyvg/vue3-notification";

const app = createApp(App);
const pinia = createPinia();
app.use(plugin, defaultConfig(genesis));
app.use(axiosHelper);
app.use(pinia);
app.use(Notifications);
app.mount("#app");
