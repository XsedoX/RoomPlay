import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import { aliases as mdAliases, md } from 'vuetify/iconsets/md';
import 'vuetify/styles';
import GoogleIcon from '@/pages/login_page/GoogleIcon.vue';
import { darkTheme, lightTheme } from '@/assets/themes.ts';
import { TThemes } from './pages/room_page/settings_menu/TThemes';

export const customIcons = {
  googleIcon: GoogleIcon,
};

export default createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'md',
    sets: {
      md,
    },
    aliases: {
      ...mdAliases,
      ...customIcons,
    },
  },
  theme: {
    defaultTheme: TThemes.LightMode,
    themes: {
      [TThemes.LightMode]: lightTheme,
      [TThemes.DarkMode]: darkTheme,
    },
  },
});
