import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.cqupt.nodedesign.app',
  appName: 'emqx-admin',
  webDir: 'dist',
    server: {
    androidScheme: 'http',
        allowNavigation: [
            '172.20.10.5'
        ]
  }
};

export default config;
