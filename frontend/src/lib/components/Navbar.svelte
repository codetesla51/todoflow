<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import ThemeToggle from './ThemeToggle.svelte';

	let isAuth = $state(false);
	let username = $state('');
	let dropdownOpen = $state(false);

	auth.isAuthenticated.subscribe((v) => (isAuth = v));
	auth.user.subscribe((u) => {
		if (u) username = u.username;
	});

	function handleLogout() {
		auth.logout();
		dropdownOpen = false;
		goto('/login');
	}

	function closeDropdown() {
		dropdownOpen = false;
	}
</script>

<svelte:window onclick={() => { if (dropdownOpen) dropdownOpen = false; }} />

<header
	class="sticky top-0 z-50 border-b transition-all duration-200"
	style="border-color: var(--border-color); background: var(--bg-primary);"
>
	<div class="mx-auto flex h-14 max-w-5xl items-center justify-between px-4 sm:px-6">
		<a href="/" class="flex items-center gap-2 group">
			<span class="flex h-7 w-7 items-center justify-center rounded-md bg-accent-500 transition-transform duration-200 group-hover:scale-105">
				<svg class="h-4 w-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
				</svg>
			</span>
			<span class="text-base font-semibold tracking-tight" style="color: var(--text-primary);">
				TodoFlow
			</span>
		</a>

		<div class="flex items-center gap-2.5">
			<ThemeToggle />

			{#if isAuth}
				<div class="relative">
					<button
						onclick={(e) => { e.stopPropagation(); dropdownOpen = !dropdownOpen; }}
						class="flex items-center gap-2 rounded-lg border px-3 py-1.5 text-sm transition-colors duration-150 cursor-pointer"
						style="border-color: var(--border-color); background: var(--bg-secondary);"
					>
						<span class="flex h-6 w-6 items-center justify-center rounded-md bg-accent-500/15 text-accent-500 text-xs font-bold">
							{username.charAt(0).toUpperCase()}
						</span>
						<span class="hidden sm:block font-medium" style="color: var(--text-primary);">{username}</span>
						<svg class="h-3.5 w-3.5 transition-transform duration-150" class:rotate-180={dropdownOpen} style="color: var(--text-muted);" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
						</svg>
					</button>

					{#if dropdownOpen}
						<!-- svelte-ignore a11y_no_static_element_interactions -->
						<div
							class="absolute right-0 mt-1.5 w-48 rounded-lg border py-1 shadow-lg"
							style="border-color: var(--border-color); background: var(--bg-card);"
							onclick={(e) => e.stopPropagation()}
							onkeydown={(e) => { if (e.key === 'Escape') dropdownOpen = false; }}
						>
							<div class="px-3 py-2 border-b" style="border-color: var(--border-color);">
								<p class="text-sm font-medium" style="color: var(--text-primary);">{username}</p>
							</div>
							<button
								onclick={handleLogout}
								class="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer"
							>
								<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
									<path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
								</svg>
								Sign out
							</button>
						</div>
					{/if}
				</div>
			{:else}
				<a
					href="/login"
					class="rounded-lg px-3.5 py-1.5 text-sm font-medium transition-colors duration-150"
					style="color: var(--text-secondary);"
				>
					Sign in
				</a>
				<a
					href="/register"
					class="rounded-lg bg-accent-500 px-3.5 py-1.5 text-sm font-medium text-white transition-colors duration-150 hover:bg-accent-600"
				>
					Get started
				</a>
			{/if}
		</div>
	</div>
</header>
