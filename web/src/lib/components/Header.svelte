<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { createProfileQuery } from '$lib/api/auth';

	export let showQuerySearch = false;
	export let queryText = '';

	const profileQuery = createProfileQuery();

	function onSubmit() {
		const query = new URLSearchParams($page.url.search);
		query.set('query', queryText);
		goto(`/search?${query.toString()}`);
	}

	function onSignOut() {
		window.location.href = '/auth/sign-out';
	}
</script>

<div class="flex h-16 w-full flex-row items-center gap-4 border-b p-4">
	<h1 class="text-2xl">
		<a href="/">Vocabulary Leveling</a>
	</h1>

	{#if showQuerySearch}
		<form class="w-[586px] p-4" on:submit|preventDefault={onSubmit}>
			<input
				class="w-full rounded-md border-2 border-gray-300 px-2 py-1"
				type="text"
				bind:value={queryText}
			/>
			<button class="hidden">Search</button>
		</form>
	{/if}

	<span class="flex-grow"></span>

	{#if $profileQuery.isSuccess}
		{#if $profileQuery.data}
			<div class="flex flex-row">
				<span class="mr-1">{$profileQuery.data.name}</span>
				<button on:click={onSignOut}>(Sign out)</button>
			</div>
		{:else}
			<span>Unknown</span>
		{/if}
	{/if}
</div>
