<script lang="ts">
	import { createProfileQuery } from '$lib/api/auth';
	import Header from '$lib/components/Header.svelte';

	const profileQuery = createProfileQuery();
</script>

<Header />
{#if $profileQuery.isSuccess && $profileQuery.data && $profileQuery.data.roles.includes('Admin')}
	<slot />
{:else}
	<div>You do not have access to this page</div>
{/if}
