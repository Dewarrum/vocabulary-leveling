<script lang="ts">
	import { page } from '$app/stores';
	import QueryHeader from '$lib/components/QueryHeader.svelte';

	const query = $page.url.searchParams.get('query') ?? '';
</script>

<QueryHeader queryText={query} />
<slot />
