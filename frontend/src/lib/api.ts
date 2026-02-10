import { auth } from '$lib/stores/auth';

const API_BASE = 'http://localhost:8080';

async function request(endpoint: string, options: RequestInit = {}) {
	const token = auth.getToken();
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string>)
	};

	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const res = await fetch(`${API_BASE}${endpoint}`, {
		...options,
		headers
	});

	const data = await res.json();

	if (res.status === 401 && !endpoint.startsWith('/auth/')) {
		auth.logout();
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}
		throw new Error('Unauthorized');
	}

	if (!res.ok) {
		throw new Error(data.error || 'Something went wrong');
	}

	return data;
}

// Auth
export async function register(username: string, email: string, password: string) {
	return request('/auth/register', {
		method: 'POST',
		body: JSON.stringify({ username, email, password })
	});
}

export async function login(email: string, password: string) {
	return request('/auth/login', {
		method: 'POST',
		body: JSON.stringify({ email, password })
	});
}

// Profile
export async function getProfile() {
	return request('/api/profile');
}

// Todos
export async function getTodos(page = 1, limit = 20) {
	return request(`/api/todos?page=${page}&limit=${limit}`);
}

export async function getTodo(id: number) {
	return request(`/api/todos/${id}`);
}

export async function createTodo(title: string, description: string) {
	return request('/api/todos', {
		method: 'POST',
		body: JSON.stringify({ title, description })
	});
}

export async function updateTodo(id: number, title: string, description: string) {
	return request(`/api/todos/${id}`, {
		method: 'PUT',
		body: JSON.stringify({ title, description })
	});
}

export async function updateTodoStatus(id: number, status: 'pending' | 'completed') {
	return request(`/api/todos/${id}/status`, {
		method: 'PATCH',
		body: JSON.stringify({ status })
	});
}

export async function deleteTodo(id: number) {
	return request(`/api/todos/${id}`, {
		method: 'DELETE'
	});
}
