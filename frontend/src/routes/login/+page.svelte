<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { login as apiLogin } from '$lib/api';
	import Alert from '$lib/components/Alert.svelte';
	import { onMount } from 'svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	let showPassword = $state(false);

	onMount(() => {
		let isAuth = false;
		const unsub = auth.isAuthenticated.subscribe((v) => (isAuth = v));
		if (isAuth) goto('/dashboard');
		unsub();
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (loading) return;
		error = '';
		loading = true;
		try {
			const res = await apiLogin(email, password);
			auth.login(res.token, res.user);
			await goto('/dashboard');
		} catch (err: any) {
			error = err.message || 'Invalid email or password';
			loading = false;
		}
	}
</script>

<div class="flex min-h-[calc(100vh-8rem)] items-center justify-center px-4 py-12">
	<div class="w-full max-w-sm">
		<div class="mb-7">
			<h1 class="text-2xl font-bold" style="color: var(--text-primary);">Sign in</h1>
			<p class="mt-1 text-sm" style="color: var(--text-muted);">Enter your credentials to continue</p>
		</div>

		{#if error}
			<div class="mb-5">
				<Alert message={error} type="error" />
			</div>
		{/if}

		<form onsubmit={handleSubmit} class="space-y-4">
			<div>
				<label for="email" class="mb-1.5 block text-sm font-medium" style="color: var(--text-secondary);">
					Email
				</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					placeholder="you@example.com"
					required
					class="w-full rounded-lg border px-3.5 py-2.5 text-sm outline-none transition-colors duration-150 focus:border-accent-500 focus:ring-1 focus:ring-accent-500"
					style="border-color: var(--border-color); background: var(--bg-input); color: var(--text-primary);"
				/>
			</div>

			<div>
				<label for="password" class="mb-1.5 block text-sm font-medium" style="color: var(--text-secondary);">
					Password
				</label>
				<div class="relative">
					<input
						id="password"
						type={showPassword ? 'text' : 'password'}
						bind:value={password}
						placeholder="••••••••"
						required
						class="w-full rounded-lg border px-3.5 py-2.5 pr-10 text-sm outline-none transition-colors duration-150 focus:border-accent-500 focus:ring-1 focus:ring-accent-500"
						style="border-color: var(--border-color); background: var(--bg-input); color: var(--text-primary);"
					/>
					<button
						type="button"
						onclick={() => (showPassword = !showPassword)}
						class="absolute inset-y-0 right-0 flex items-center pr-3 cursor-pointer"
						style="color: var(--text-muted);"
						aria-label="Toggle password"
					>
						{#if showPassword}
							<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
							</svg>
						{:else}
							<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
								<path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
							</svg>
						{/if}
					</button>
				</div>
			</div>

			<button
				type="submit"
				disabled={loading}
				class="mt-2 flex w-full items-center justify-center gap-2 rounded-lg bg-accent-500 py-2.5 text-sm font-semibold text-white transition-colors duration-150 hover:bg-accent-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
			>
				{#if loading}
					<svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
					</svg>
				{:else}
					Sign in
				{/if}
			</button>
		</form>

		<p class="mt-6 text-center text-sm" style="color: var(--text-muted);">
			No account?
			<a href="/register" class="font-medium text-accent-500 hover:text-accent-600 transition-colors">
				Create one
			</a>
		</p>
	</div>
</div>
