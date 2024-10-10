import emitter from 'tiny-emitter/instance'

export const SNACKBAR = 'snackbar';

const SnackbarPlugin = {
  install(app, options) {
    app.config.globalProperties.$onSnackbar = callback => {
      emitter.on(SNACKBAR, message => callback(message));
    };
    app.config.globalProperties.$showSnackbar = message => {
      emitter.emit(SNACKBAR, message);
    }
  }
};

export default SnackbarPlugin;
