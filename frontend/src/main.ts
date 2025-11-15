import { createApp } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import createVuetify from '@/vuetify-setup.ts';
import App from '@/App.vue';
import router from '@/router';
const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

app.use(createVuetify);
app.use(pinia);
app.use(router);

app.mount('#app');
