<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	export let queryText = '';
	function onSubmit() {
		const query = new URLSearchParams($page.url.search);
		query.set('query', queryText);
		goto(`/search?${query.toString()}`);
	}
</script>

<div class="flex h-16 w-full flex-row items-center gap-4 border-b p-4">
	<h1 class="text-2xl">
		<a href="/">Vocabulary Leveling</a>
	</h1>

	<form class="w-[586px] p-4" on:submit|preventDefault={onSubmit}>
		<input
			class="w-full rounded-md border-2 border-gray-300 px-2 py-1"
			type="text"
			bind:value={queryText}
		/>
		<button class="hidden">Search</button>
	</form>

	<span class="flex-grow"></span>

	<span>User Name</span>
</div>
