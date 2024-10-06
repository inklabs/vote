import {VoteSDK} from './voteSdk';

const VoteSDKPlugin = {
  install(app, options) {
    app.config.globalProperties.$sdk = new VoteSDK(options.baseURL);
  }
};

export default VoteSDKPlugin;
