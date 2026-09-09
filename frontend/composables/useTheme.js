import { ref, watchEffect } from 'vue';

const STORAGE_KEY = 'sr-theme';
const theme = ref(localStorage.getItem(STORAGE_KEY) || 'dark');

watchEffect(() => {
    document.documentElement.dataset.theme = theme.value;
    localStorage.setItem(STORAGE_KEY, theme.value);
});

export function useTheme() {
    const toggle = () => { theme.value = theme.value === 'dark' ? 'light' : 'dark'; };
    return { theme, toggle };
}
