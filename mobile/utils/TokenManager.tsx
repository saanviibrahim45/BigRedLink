import * as Keychain from 'react-native-keychain';

const ACCESS_KEY  = 'ACCESS_TOKEN';
const REFRESH_KEY = 'REFRESH_TOKEN';

export default {

  save: async (accessToken: string, refreshToken: string): Promise<void> => {
    try {
      await Keychain.setGenericPassword(ACCESS_KEY, accessToken);
      await Keychain.setInternetCredentials(REFRESH_KEY, REFRESH_KEY, refreshToken);
    } catch (e) {
      console.error('TokenManager.save failed', e);
      throw e;
    }
  },


  getAccess: async (): Promise<string | null> => {
    try {
      const result = await Keychain.getGenericPassword();
      if (result && result.username === ACCESS_KEY) {
        return result.password;
      }
      return null;
    } catch (e) {
      console.error('TokenManager.getAccess failed', e);
      return null;
    }
  },


  getRefresh: async (): Promise<string | null> => {
    try {
      const result = await Keychain.getInternetCredentials(REFRESH_KEY);
      if (result && result.username === REFRESH_KEY) {
        return result.password;
      }
      return null;
    } catch (e) {
      console.error('TokenManager.getRefresh failed', e);
      return null;
    }
  },


  clear: async (): Promise<void> => {
    try {
      await Keychain.resetGenericPassword();
      await Keychain.resetInternetCredentials({ server: REFRESH_KEY });
    } catch (e) {
      console.error('TokenManager.clear failed', e);
      throw e;
    }
  },
};
