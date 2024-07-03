<script lang="ts">
	import { goto } from '$app/navigation';
	import { createProfileQuery } from '$lib/api/auth';
	import Header from '$lib/components/Header.svelte';
	import SignInForm from '$lib/components/SignInForm.svelte';

	const profileQuery = createProfileQuery();

	let query = '';
	function onSubmit() {
		if (query === '') {
			return;
		}

		goto(`/search?query=${query}`);
	}
</script>

<Header />
<main class="mt-96 flex flex-col items-center">
	{#if $profileQuery.isSuccess && $profileQuery.data}
		<h1 class="text-3xl">Vocabulary Leveling</h1>
		<form class="w-[586px] p-4" on:submit|preventDefault={onSubmit}>
			<input
				class="w-full rounded-md border-2 border-gray-300 p-4"
				type="text"
				bind:value={query}
			/>
			<button class="hidden">Search</button>
		</form>
	{:else}
		<SignInForm />
	{/if}
</main>
