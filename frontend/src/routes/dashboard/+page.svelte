<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { todos, type Todo } from '$lib/stores/todos';
	import {
		getTodos,
		createTodo as apiCreateTodo,
		updateTodo as apiUpdateTodo,
	} from '$lib/api';
	import TodoCard from '$lib/components/TodoCard.svelte';
	import TodoModal from '$lib/components/TodoModal.svelte';

	let loading = $state(true);
	let submitting = $state(false);
	let todoList: Todo[] = $state([]);
	let showModal = $state(false);
	let editingTodo: Todo | null = $state(null);
	let filter = $state<'all' | 'pending' | 'completed'>('all');
	let searchQuery = $state('');
	let page = $state(1);
	let hasMore = $state(true);
	let loadingMore = $state(false);
	let username = $state('');

	todos.subscribe((v) => (todoList = v));
	auth.user.subscribe((u) => {
		if (u) username = u.username;
	});

	const filteredTodos = $derived.by(() => {
		let list = todoList;
		if (filter !== 'all') {
			list = list.filter((t) => t.status === filter);
		}
		if (searchQuery.trim()) {
			const q = searchQuery.toLowerCase();
			list = list.filter(
				(t) =>
					t.title.toLowerCase().includes(q) ||
					(t.description && t.description.toLowerCase().includes(q))
			);
		}
		return list;
	});

	const stats = $derived.by(() => {
		const total = todoList.length;
		const completed = todoList.filter((t) => t.status === 'completed').length;
		const pending = total - completed;
		const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
		return { total, completed, pending, percentage };
	});

	onMount(async () => {
		let isAuth = false;
		const unsub = auth.isAuthenticated.subscribe((v) => (isAuth = v));
		unsub();

		if (!isAuth) {
			goto('/login');
			return;
		}

		await fetchTodos();
	});

	async function fetchTodos() {
		loading = true;
		try {
			const data = await getTodos(1, 50);
			todos.set(data || []);
		} catch {
			// handled by api
		} finally {
			loading = false;
		}
	}

	async function loadMore() {
		if (loadingMore) return;
		loadingMore = true;
		page += 1;
		try {
			const data = await getTodos(page, 50);
			if (!data || data.length === 0) {
				hasMore = false;
			} else {
				todos.update((current) => [...current, ...data]);
			}
		} catch {
			//
		} finally {
			loadingMore = false;
		}
	}

	function openCreate() {
		editingTodo = null;
		showModal = true;
	}

	function openEdit(todo: Todo) {
		editingTodo = todo;
		showModal = true;
	}

	async function handleSubmit(title: string, description: string) {
		submitting = true;
		try {
			if (editingTodo) {
				const updated = await apiUpdateTodo(editingTodo.id, title, description);
				todos.updateTodo(updated);
			} else {
				const res = await apiCreateTodo(title, description);
				todos.add(res.todo);
			}
			showModal = false;
			editingTodo = null;
		} catch (e) {
			console.error(e);
		} finally {
			submitting = false;
		}
	}

	function getGreeting() {
		const hour = new Date().getHours();
		if (hour < 12) return 'Good morning';
		if (hour < 18) return 'Good afternoon';
		return 'Good evening';
	}
</script>

<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 sm:py-8">
	<!-- Header Section -->
	<div class="mb-6 flex flex-col gap-3 sm:mb-8 sm:flex-row sm:items-end sm:justify-between">
		<div>
			<p class="text-sm" style="color: var(--text-muted);">
				{getGreeting()},
			</p>
			<h1 class="mt-1 text-xl font-semibold tracking-tight sm:text-2xl" style="color: var(--text-primary);">
				{username}
			</h1>
		</div>
		<button
			onclick={openCreate}
			class="inline-flex w-full items-center justify-center gap-2 rounded-lg bg-accent-500 px-4 py-2.5 text-sm font-medium text-white transition-colors duration-150 hover:bg-accent-600 active:scale-[0.98] cursor-pointer sm:w-auto sm:py-2"
		>
			<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
			</svg>
			New task
		</button>
	</div>

	<!-- Stats -->
	<div class="mb-6 grid grid-cols-2 gap-2.5 sm:grid-cols-4 sm:gap-3">
		<div class="rounded-lg border p-3 sm:p-4" style="border-color: var(--border-color); background: var(--bg-card);">
			<p class="text-xl font-semibold tabular-nums sm:text-2xl" style="color: var(--text-primary);">{stats.total}</p>
			<p class="text-xs" style="color: var(--text-muted);">Total</p>
		</div>
		<div class="rounded-lg border p-3 sm:p-4" style="border-color: var(--border-color); background: var(--bg-card);">
			<p class="text-xl font-semibold tabular-nums sm:text-2xl" style="color: var(--text-primary);">{stats.pending}</p>
			<p class="text-xs" style="color: var(--text-muted);">Pending</p>
		</div>
		<div class="rounded-lg border p-3 sm:p-4" style="border-color: var(--border-color); background: var(--bg-card);">
			<p class="text-xl font-semibold tabular-nums sm:text-2xl" style="color: var(--text-primary);">{stats.completed}</p>
			<p class="text-xs" style="color: var(--text-muted);">Done</p>
		</div>
		<div class="rounded-lg border p-3 sm:p-4" style="border-color: var(--border-color); background: var(--bg-card);">
			<p class="text-xl font-semibold tabular-nums sm:text-2xl" style="color: var(--text-primary);">{stats.percentage}%</p>
			<p class="text-xs" style="color: var(--text-muted);">Progress</p>
		</div>
	</div>

	{#if stats.total > 0}
		<div class="mb-6">
			<div class="h-1.5 w-full overflow-hidden rounded-full" style="background: var(--bg-tertiary);">
				<div
					class="h-full rounded-full bg-accent-500 transition-all duration-500 ease-out"
					style="width: {stats.percentage}%"
				></div>
			</div>
		</div>
	{/if}

	<!-- Filters & Search -->
	<div class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex w-full items-center gap-1 overflow-x-auto rounded-lg border p-0.5 sm:w-auto" style="border-color: var(--border-color); background: var(--bg-secondary); -webkit-overflow-scrolling: touch;">
			{#each [['all', 'All'], ['pending', 'Pending'], ['completed', 'Done']] as [value, label]}
				<button
					onclick={() => (filter = value as typeof filter)}
					class="flex-1 whitespace-nowrap rounded-md px-3.5 py-1.5 text-sm font-medium transition-colors duration-150 cursor-pointer sm:flex-none"
					style={filter === value
						? 'background: var(--bg-card); color: var(--text-primary);'
						: 'color: var(--text-muted);'}
				>
					{label}
				</button>
			{/each}
		</div>

		<div class="relative w-full sm:w-auto">
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
				<svg class="h-4 w-4" style="color: var(--text-muted);" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
				</svg>
			</div>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search tasks..."
				class="w-full rounded-lg border py-2.5 pl-9 pr-4 text-sm outline-none transition-colors duration-150 focus:border-accent-500 sm:w-56 sm:py-2"
				style="border-color: var(--border-color); background: var(--bg-input); color: var(--text-primary);"
			/>
		</div>
	</div>

	<!-- Todo List -->
	{#if loading}
		<div class="flex flex-col items-center justify-center py-20">
			<svg class="h-10 w-10 animate-spin text-accent-500" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
			</svg>
			<p class="mt-4 text-sm font-medium" style="color: var(--text-muted);">Loading your tasks...</p>
		</div>
	{:else if filteredTodos.length === 0}
		<div class="flex flex-col items-center justify-center py-16">
			<div class="mb-3 flex h-14 w-14 items-center justify-center rounded-xl" style="background: var(--bg-tertiary);">
				<svg class="h-7 w-7" style="color: var(--text-muted);" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
				</svg>
			</div>
			{#if searchQuery}
				<h3 class="text-lg font-semibold" style="color: var(--text-primary);">No results found</h3>
				<p class="mt-1 text-sm" style="color: var(--text-muted);">Try a different search term</p>
			{:else if filter !== 'all'}
				<h3 class="text-lg font-semibold" style="color: var(--text-primary);">No {filter} tasks</h3>
				<p class="mt-1 text-sm" style="color: var(--text-muted);">
					{filter === 'pending' ? 'All caught up! Great job.' : 'Complete some tasks to see them here.'}
				</p>
			{:else}
				<h3 class="text-lg font-semibold" style="color: var(--text-primary);">No tasks yet</h3>
				<p class="mt-1 text-sm" style="color: var(--text-muted);">Create your first task to get started</p>
				<button
					onclick={openCreate}
					class="mt-4 inline-flex items-center gap-2 rounded-lg bg-accent-500 px-4 py-2 text-sm font-medium text-white transition-colors duration-150 hover:bg-accent-600 cursor-pointer"
				>
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
					</svg>
					Create Task
				</button>
			{/if}
		</div>
	{:else}
		<div class="space-y-3">
			{#each filteredTodos as todo (todo.id)}
				<TodoCard {todo} onEdit={openEdit} />
			{/each}
		</div>

		{#if hasMore && filteredTodos.length >= 50}
			<div class="mt-6 flex justify-center">
				<button
					onclick={loadMore}
					disabled={loadingMore}
					class="inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-medium transition-colors duration-150 cursor-pointer"
					style="border-color: var(--border-color); color: var(--text-secondary); background: var(--bg-secondary);"
				>
					{#if loadingMore}
						<svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
						</svg>
						Loading...
					{:else}
						Load more
					{/if}
				</button>
			</div>
		{/if}
	{/if}
</div>

<TodoModal
	show={showModal}
	{editingTodo}
	{submitting}
	onClose={() => {
		showModal = false;
		editingTodo = null;
	}}
	onSubmit={handleSubmit}
/>
