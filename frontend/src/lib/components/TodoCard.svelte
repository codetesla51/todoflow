<script lang="ts">
	import type { Todo } from '$lib/stores/todos';
	import { updateTodoStatus, deleteTodo as apiDeleteTodo } from '$lib/api';
	import { todos } from '$lib/stores/todos';

	interface Props {
		todo: Todo;
		onEdit: (todo: Todo) => void;
	}
	let { todo, onEdit }: Props = $props();
	let loading = $state(false);
	let deleting = $state(false);

	function formatDate(dateStr: string) {
		const d = new Date(dateStr);
		const now = new Date();
		const diff = now.getTime() - d.getTime();
		const minutes = Math.floor(diff / 60000);
		const hours = Math.floor(minutes / 60);
		const days = Math.floor(hours / 24);

		if (minutes < 1) return 'Just now';
		if (minutes < 60) return `${minutes}m ago`;
		if (hours < 24) return `${hours}h ago`;
		if (days < 7) return `${days}d ago`;
		return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	}

	async function toggleStatus() {
		if (loading) return;
		loading = true;
		try {
			const newStatus = todo.status === 'completed' ? 'pending' : 'completed';
			const updated = await updateTodoStatus(todo.id, newStatus);
			todos.updateTodo(updated);
		} catch (e) {
			console.error(e);
		} finally {
			loading = false;
		}
	}

	async function handleDelete() {
		if (deleting) return;
		deleting = true;
		try {
			await apiDeleteTodo(todo.id);
			todos.remove(todo.id);
		} catch (e) {
			console.error(e);
		} finally {
			deleting = false;
		}
	}
</script>

<div
	class="group relative rounded-lg border p-3.5 transition-colors duration-150 sm:p-4"
	style="border-color: var(--border-color); background: var(--bg-card);"
	class:opacity-60={todo.status === 'completed'}
>
	<div class="flex items-start gap-3">
		<!-- Checkbox -->
		<button
			onclick={toggleStatus}
			disabled={loading}
			class="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded border-2 transition-colors duration-150 cursor-pointer sm:h-5 sm:w-5"
			class:bg-accent-500={todo.status === 'completed'}
			class:border-accent-500={todo.status === 'completed'}
			style={todo.status !== 'completed' ? 'border-color: var(--border-hover);' : ''}
		>
			{#if todo.status === 'completed'}
				<svg class="h-3.5 w-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
					<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
				</svg>
			{/if}
		</button>

		<!-- Content -->
		<div class="min-w-0 flex-1">
			<h3
				class="text-[0.95rem] font-medium leading-snug transition-all duration-200"
				class:line-through={todo.status === 'completed'}
				style="color: {todo.status === 'completed' ? 'var(--text-muted)' : 'var(--text-primary)'};"
			>
				{todo.title}
			</h3>
			{#if todo.description}
				<p
					class="mt-1.5 text-sm leading-relaxed line-clamp-2"
					style="color: var(--text-secondary);"
				>
					{todo.description}
				</p>
			{/if}
			<div class="mt-2.5 flex items-center gap-3">
				<span
					class="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium {todo.status === 'pending' ? 'bg-amber-500/10 text-amber-600' : 'bg-emerald-500/10 text-emerald-600'}"
				>
					<span
						class="h-1.5 w-1.5 rounded-full"
						class:bg-amber-500={todo.status === 'pending'}
						class:bg-emerald-500={todo.status === 'completed'}
					></span>
					{todo.status === 'completed' ? 'Done' : 'Pending'}
				</span>
				<span class="text-xs" style="color: var(--text-muted);">{formatDate(todo.created_at)}</span>
			</div>
		</div>

		<!-- Actions -->
		<div class="flex shrink-0 items-center gap-1 opacity-100 sm:opacity-0 sm:transition-opacity sm:duration-200 sm:group-hover:opacity-100">
			<button
				onclick={() => onEdit(todo)}
				class="flex h-9 w-9 items-center justify-center rounded-lg transition-colors duration-200 cursor-pointer active:scale-95 sm:h-8 sm:w-8"
				style="color: var(--text-muted);"
				onmouseenter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--bg-hover)'; (e.currentTarget as HTMLElement).style.color = 'var(--text-primary)'; }}
				onmouseleave={(e) => { (e.currentTarget as HTMLElement).style.background = 'transparent'; (e.currentTarget as HTMLElement).style.color = 'var(--text-muted)'; }}
				aria-label="Edit todo"
			>
				<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
				</svg>
			</button>
			<button
				onclick={handleDelete}
				disabled={deleting}
				class="flex h-9 w-9 items-center justify-center rounded-lg transition-colors duration-200 cursor-pointer text-red-400 active:scale-95 hover:bg-red-500/10 hover:text-red-500 sm:h-8 sm:w-8"
				aria-label="Delete todo"
			>
				{#if deleting}
					<svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
					</svg>
				{:else}
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
					</svg>
				{/if}
			</button>
		</div>
	</div>
</div>
