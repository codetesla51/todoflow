import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createThemeStore() {
	const stored = browser ? localStorage.getItem('theme') : null;
	const prefersDark = browser ? window.matchMedia('(prefers-color-scheme: dark)').matches : false;
	const initial = stored || (prefersDark ? 'dark' : 'light');

	const theme = writable<'light' | 'dark'>(initial as 'light' | 'dark');

	if (browser) {
		document.documentElement.classList.toggle('dark', initial === 'dark');
	}

	function toggle() {
		theme.update((current) => {
			const next = current === 'light' ? 'dark' : 'light';
			if (browser) {
				localStorage.setItem('theme', next);
				document.documentElement.classList.toggle('dark', next === 'dark');
			}
			return next;
		});
	}

	function set(value: 'light' | 'dark') {
		theme.set(value);
		if (browser) {
			localStorage.setItem('theme', value);
			document.documentElement.classList.toggle('dark', value === 'dark');
		}
	}

	return {
		subscribe: theme.subscribe,
		toggle,
		set
	};
}

export const theme = createThemeStore();
