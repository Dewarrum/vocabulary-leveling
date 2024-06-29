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

<div class="w-full flex flex-row gap-4 h-16 items-center p-4 border-b">
	<h1 class="text-2xl">
		<a href="/">Vocabulary Leveling</a>
	</h1>

	<form class="w-[586px] p-4" on:submit|preventDefault={onSubmit}>
		<input
			class="w-full border-2 border-gray-300 rounded-md py-1 px-2"
			type="text"
			bind:value={queryText}
		/>
		<button class="hidden">Search</button>
	</form>

	<span class="flex-grow"></span>

	<span>User Name</span>
</div>
