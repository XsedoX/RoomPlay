import { describe, it } from 'vitest';
import { mountVuetify } from '@/vuetify-setup.ts';
import RoomData from '@/login_page/RoomData.vue';


describe('Component with room name and password is rendered correcty', () => {
  it("checks if room name text field is rendered", () =>{
    const component = mountVuetify(RoomData);


  });
});
