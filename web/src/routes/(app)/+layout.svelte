<script lang="ts">
	import { page } from '$app/stores';
	import Header from '$lib/components/Header.svelte';

	const query = $page.url.searchParams.get('query') ?? '';
</script>

<Header showQuerySearch={true} queryText={query} />
<slot />
