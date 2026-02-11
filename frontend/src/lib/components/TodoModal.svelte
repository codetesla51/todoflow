<script lang="ts">
	import type { Todo } from '$lib/stores/todos';

	interface Props {
		show?: boolean;
		editingTodo?: Todo | null;
		submitting?: boolean;
		onClose: () => void;
		onSubmit: (title: string, description: string) => void;
	}
	let {
		show = false,
		editingTodo = null,
		submitting = false,
		onClose,
		onSubmit
	}: Props = $props();

	let title = $state('');
	let description = $state('');

	$effect(() => {
		if (editingTodo) {
			title = editingTodo.title;
			description = editingTodo.description || '';
		} else {
			title = '';
			description = '';
		}
	});

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (!title.trim()) return;
		onSubmit(title.trim(), description.trim());
		title = '';
		description = '';
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onClose();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 z-50 flex items-end justify-center p-0 sm:items-center sm:p-4"
		role="dialog"
		aria-modal="true"
	>
		<div
			class="absolute inset-0 bg-black/40 backdrop-blur-sm"
			onclick={onClose}
			role="button"
			tabindex="-1"
			aria-label="Close modal"
			onkeydown={(e) => { if (e.key === 'Enter') onClose(); }}
		></div>

		<!-- Modal -->
		<div
			class="relative w-full max-w-md self-end rounded-t-xl border p-5 shadow-lg sm:self-center sm:rounded-lg"
			style="border-color: var(--border-color); background: var(--bg-card);"
		>
			<!-- Header -->
			<div class="mb-5 flex items-center justify-between">
				<h2 class="text-base font-semibold" style="color: var(--text-primary);">
					{editingTodo ? 'Edit task' : 'New task'}
				</h2>
				<button
					onclick={onClose}
					class="flex h-8 w-8 items-center justify-center rounded-lg transition-colors cursor-pointer"
					style="color: var(--text-muted);"
					onmouseenter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--bg-hover)'; }}
					onmouseleave={(e) => { (e.currentTarget as HTMLElement).style.background = 'transparent'; }}
					aria-label="Close"
				>
					<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<!-- Form -->
			<form onsubmit={handleSubmit}>
				<div class="space-y-4">
					<div>
						<label for="title" class="mb-1.5 block text-sm font-medium" style="color: var(--text-secondary);">
							Title <span class="text-red-400">*</span>
						</label>
						<input
							id="title"
							type="text"
							bind:value={title}
							placeholder="What needs to be done?"
							disabled={submitting}
							class="w-full rounded-lg border px-3.5 py-2.5 text-sm outline-none transition-colors duration-150 focus:border-accent-500 disabled:opacity-50"
							style="border-color: var(--border-color); background: var(--bg-input); color: var(--text-primary);"
							required
						/>
					</div>
					<div>
						<label for="description" class="mb-1.5 block text-sm font-medium" style="color: var(--text-secondary);">
							Description
						</label>
						<textarea
							id="description"
							bind:value={description}
							placeholder="Add some details..."
							rows="3"
							disabled={submitting}
							class="w-full resize-none rounded-lg border px-3.5 py-2.5 text-sm outline-none transition-colors duration-150 focus:border-accent-500 disabled:opacity-50"
							style="border-color: var(--border-color); background: var(--bg-input); color: var(--text-primary);"
						></textarea>
					</div>
				</div>

				<div class="mt-5 flex items-center justify-end gap-2">
					<button
						type="button"
						onclick={onClose}
						disabled={submitting}
						class="rounded-lg border px-4 py-2 text-sm font-medium transition-colors duration-150 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
						style="border-color: var(--border-color); color: var(--text-secondary); background: var(--bg-secondary);"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={!title.trim() || submitting}
						class="inline-flex items-center gap-2 rounded-lg bg-accent-500 px-4 py-2 text-sm font-medium text-white transition-colors duration-150 hover:bg-accent-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
					>
						{#if submitting}
							<svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
							</svg>
							{editingTodo ? 'Saving...' : 'Creating...'}
						{:else}
							{editingTodo ? 'Save' : 'Create'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
