<script lang="ts">
	import './layout.css';
	import Navbar from '$lib/components/Navbar.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<title>TodoFlow — Organize your life</title>
	<meta name="description" content="A beautiful, modern todo app to stay organized." />
</svelte:head>

<div class="flex min-h-screen flex-col">
	<Navbar />
	<main class="flex-1">
		{@render children()}
	</main>
	<footer class="border-t py-5 sm:py-6" style="border-color: var(--border-color);">
		<div class="mx-auto max-w-6xl px-4 sm:px-6">
			<p class="text-center text-sm" style="color: var(--text-muted);">
				&copy; {new Date().getFullYear()} TodoFlow. Built with Svelte & Go.
			</p>
		</div>
	</footer>
</div>
