import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases as mdAliases, md } from 'vuetify/iconsets/md'
import 'vuetify/styles'
import GoogleIcon from '@/pages/login_page/GoogleIcon.vue';
import { mount } from '@vue/test-utils';
import type { Component } from 'vue';
import vuetify from '@/vuetify-setup.ts'
import { darkTheme, lightTheme } from '@/assets/themes.ts';


export const customIcons = {
  googleIcon: GoogleIcon,
}

export default createVuetify({
  components,
  directives,
  icons:{
    defaultSet: 'md',
    sets:{
      md
    },
    aliases: {
      ...mdAliases,
      ...customIcons
    },
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: lightTheme,
      dark: darkTheme
    }
  }
});

export function mountVuetify(component: Component, customProps?: Record<string, unknown>) {
  return mount(component, {
    ...(customProps && { props: customProps }),
    global: {
      components: { component },
      plugins: [vuetify],
    },
  });
}
