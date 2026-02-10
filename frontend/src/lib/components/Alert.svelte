<script lang="ts">
	interface Props {
		message?: string;
		type?: 'error' | 'success';
	}
	let { message = '', type = 'error' }: Props = $props();
</script>

{#if message}
	<div
		class="flex items-center gap-2.5 rounded-xl border px-4 py-3 text-sm font-medium animate-in slide-in-from-top-2 fade-in duration-300 {type === 'error' ? 'bg-red-500/10 border-red-500/20 text-red-500' : 'bg-emerald-500/10 border-emerald-500/20 text-emerald-500'}"
	>
		{#if type === 'error'}
			<svg class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
		{:else}
			<svg class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
		{/if}
		{message}
	</div>
{/if}
