import {JsSDK} from './jsSdk';

const VoteSDKPlugin = {
  install(app, options) {
    app.config.globalProperties.$sdk = new JsSDK(options.baseURL);
  }
};

export default VoteSDKPlugin;
