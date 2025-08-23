import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { md } from 'vuetify/iconsets/md'
import 'vuetify/styles'

export default createVuetify({
  components,
  directives,
  icons:{
    defaultSet: 'md',
    sets:{
      md,
    }
  }
})
