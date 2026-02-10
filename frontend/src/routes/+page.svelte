<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let isAuth = $state(false);
	auth.isAuthenticated.subscribe((v) => (isAuth = v));

	onMount(() => {
		if (isAuth) goto('/dashboard');
	});
</script>

<section class="mx-auto max-w-5xl px-6">
	<!-- Hero -->
	<div class="flex flex-col items-center pt-20 pb-16 sm:pt-28 sm:pb-20 lg:pt-36">
		<h1 class="text-center text-4xl font-bold tracking-tight sm:text-5xl lg:text-[3.5rem] lg:leading-[1.15]" style="color: var(--text-primary);">
			Get things done,<br />without the noise.
		</h1>

		<p class="mx-auto mt-5 max-w-lg text-center text-base leading-relaxed" style="color: var(--text-secondary);">
			A minimal task manager that stays out of your way. Fast backend, clean interface, zero clutter.
		</p>

		<div class="mt-9 flex items-center gap-3">
			<a
				href="/register"
				class="rounded-lg bg-accent-500 px-6 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:bg-accent-600 active:scale-[0.97]"
			>
				Get started
			</a>
			<a
				href="/login"
				class="rounded-lg border px-6 py-2.5 text-sm font-semibold transition-all duration-200 active:scale-[0.97]"
				style="border-color: var(--border-color); color: var(--text-primary);"
			>
				Sign in
			</a>
		</div>
	</div>

	<!-- Divider -->
	<div class="mx-auto h-px w-full max-w-3xl" style="background: var(--border-color);"></div>

	<!-- Features -->
	<div class="grid gap-0 sm:grid-cols-3 py-16 sm:py-20">
		<div class="flex flex-col items-center px-6 py-6 text-center sm:border-r" style="border-color: var(--border-color);">
			<div class="mb-3 flex h-10 w-10 items-center justify-center rounded-lg" style="background: var(--bg-tertiary);">
				<svg class="h-5 w-5 text-accent-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
					<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
				</svg>
			</div>
			<h3 class="text-sm font-semibold" style="color: var(--text-primary);">Fast</h3>
			<p class="mt-1.5 text-xs leading-relaxed" style="color: var(--text-muted);">
				Go + Redis. Sub-millisecond responses, no spinners.
			</p>
		</div>

		<div class="flex flex-col items-center px-6 py-6 text-center sm:border-r" style="border-color: var(--border-color);">
			<div class="mb-3 flex h-10 w-10 items-center justify-center rounded-lg" style="background: var(--bg-tertiary);">
				<svg class="h-5 w-5 text-accent-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
					<path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
				</svg>
			</div>
			<h3 class="text-sm font-semibold" style="color: var(--text-primary);">Secure</h3>
			<p class="mt-1.5 text-xs leading-relaxed" style="color: var(--text-muted);">
				JWT auth, rate limiting, encrypted passwords.
			</p>
		</div>

		<div class="flex flex-col items-center px-6 py-6 text-center">
			<div class="mb-3 flex h-10 w-10 items-center justify-center rounded-lg" style="background: var(--bg-tertiary);">
				<svg class="h-5 w-5 text-accent-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
					<path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
				</svg>
			</div>
			<h3 class="text-sm font-semibold" style="color: var(--text-primary);">Dark mode</h3>
			<p class="mt-1.5 text-xs leading-relaxed" style="color: var(--text-muted);">
				Adapts to your system or toggle manually.
			</p>
		</div>
	</div>
</section>
