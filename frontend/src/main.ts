import { createApp } from 'vue';
import { createPinia } from 'pinia';

import createVuetify from '@/vuetify-setup.ts';
import App from '@/App.vue';
import router from '@/router';
const app = createApp(App);


app.use(createVuetify);
app.use(createPinia());
app.use(router);

app.mount('#app');
