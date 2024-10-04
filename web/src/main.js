/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import withUUID from "vue-uuid";

const app = withUUID(
    createApp(App)
);

registerPlugins(app)

app.mount('#app')
