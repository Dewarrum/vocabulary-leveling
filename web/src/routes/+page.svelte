<script lang="ts">
	import { goto } from '$app/navigation';
	import SimpleHeader from '$lib/components/SimpleHeader.svelte';

	let query = '';
	function onSubmit(event: SubmitEvent) {
		if (query === '') {
			return;
		}

		goto(`/search?query=${query}`);
	}
</script>

<SimpleHeader />
<main class="mt-96 flex flex-col items-center">
	<h1 class="text-3xl">Vocabulary Leveling</h1>
	<form class="w-[586px] p-4" on:submit|preventDefault={onSubmit}>
		<input class="w-full rounded-md border-2 border-gray-300 p-4" type="text" bind:value={query} />
		<button class="hidden">Search</button>
	</form>
</main>
