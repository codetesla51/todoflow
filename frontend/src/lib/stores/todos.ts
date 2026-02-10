import { writable } from 'svelte/store';

export interface Todo {
	id: number;
	user_id: number;
	title: string;
	description: string;
	status: 'pending' | 'completed';
	created_at: string;
	updated_at: string;
}

function createTodoStore() {
	const { subscribe, set, update } = writable<Todo[]>([]);

	return {
		subscribe,
		set,
		update,
		add(todo: Todo) {
			update((todos) => [todo, ...todos]);
		},
		remove(id: number) {
			update((todos) => todos.filter((t) => t.id !== id));
		},
		updateTodo(updated: Todo) {
			update((todos) => todos.map((t) => (t.id === updated.id ? updated : t)));
		},
		clear() {
			set([]);
		}
	};
}

export const todos = createTodoStore();
