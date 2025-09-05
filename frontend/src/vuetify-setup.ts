import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md } from 'vuetify/iconsets/md'
import 'vuetify/styles'
import GoogleIcon from '@/login_page/GoogleIcon.vue';
import { mount } from '@vue/test-utils';
import type { Component } from 'vue';
import vuetify from '@/vuetify-setup.ts'


export const customIcons = {
  googleIcon: GoogleIcon,
}

export default createVuetify({
  components,
  directives,
  icons:{
    defaultSet: 'md',
    sets:{
      md,
    },
    aliases: {
      ...customIcons
    },
  }
})

export function mountVuetify(component: Component) {
  return mount(component, {
    props:{},
    global: {
      components: { component },
      plugins: [vuetify],
    },
  });
}
