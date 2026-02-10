import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
	id: number;
	username: string;
	email: string;
}

function createAuthStore() {
	const storedToken = browser ? localStorage.getItem('token') : null;
	const storedUser = browser ? localStorage.getItem('user') : null;

	const token = writable<string | null>(storedToken);
	const user = writable<User | null>(storedUser ? JSON.parse(storedUser) : null);

	const isAuthenticated = derived(token, ($token) => !!$token);

	function login(t: string, u: User) {
		token.set(t);
		user.set(u);
		if (browser) {
			localStorage.setItem('token', t);
			localStorage.setItem('user', JSON.stringify(u));
		}
	}

	function logout() {
		token.set(null);
		user.set(null);
		if (browser) {
			localStorage.removeItem('token');
			localStorage.removeItem('user');
		}
	}

	function getToken(): string | null {
		if (browser) {
			return localStorage.getItem('token');
		}
		return null;
	}

	return {
		token,
		user,
		isAuthenticated,
		login,
		logout,
		getToken
	};
}

export const auth = createAuthStore();
