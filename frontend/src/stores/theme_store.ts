import { defineStore } from 'pinia';
import { useTheme } from 'vuetify';
import { useStorage } from '@vueuse/core';
import { TThemes } from '@/pages/room_page/settings_menu/TThemes';

export const useThemeStore = defineStore('theme', () => {
  const storedTheme = useStorage('theme', TThemes.LightMode);
  const theme = useTheme();
  theme.change(storedTheme.value);

  const toggleTheme = () => {
    storedTheme.value =
      storedTheme.value === TThemes.LightMode ? TThemes.DarkMode : TThemes.LightMode;
    theme.change(storedTheme.value);
  };
  return { storedTheme, toggleTheme };
});
