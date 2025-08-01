// src/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Dialog from 'primevue/dialog';
import Toast from 'primevue/toast';
import Aura from '@primeuix/themes/aura';


const app = createApp(App)
app.use(ToastService);

app.component('Button', Button);
app.component('InputText', InputText);
app.component('Dialog', Dialog);
app.component('Toast', Toast);

const pinia = createPinia()

app.use(pinia)

// Используйте PrimeVue
app.use(PrimeVue, {
    theme: {
        preset: Aura,
    },
});

// Другие плагины
// app.use(router)

app.mount('#app')