/**
 * plugins/index.js
 *
 * Automatically included in `./src/main.js`
 */

// Plugins
import vuetify from './vuetify'
import router from '@/router'
import voteSDKPlugin from "@/plugins/voteSDKPlugin";
import snackbarPlugin from "@/plugins/snackbar";

export function registerPlugins (app) {
  app
    .use(vuetify)
    .use(router)
    .use(voteSDKPlugin, { baseURL: 'http://localhost:8080' })
    .use(snackbarPlugin)
}
