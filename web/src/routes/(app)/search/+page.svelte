<script lang="ts">
	import { page } from '$app/stores';
	import { searchSubtitles } from '$lib/api/subtitles';
	import SubtitleList from '$lib/components/SubtitleList.svelte';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	$: searchQuery = $page.url.searchParams.get('query') ?? '';
	$: subtitles = createQuery({
		queryKey: ['subtitles', searchQuery],
		queryFn: () => searchSubtitles(searchQuery)
	});

	$: subtitleId = $page.url.searchParams.get('subtitleId') ?? '';
</script>

<main class="flex flex-row gap-4 pt-8 w-3/4 mx-auto">
	<div class="flex flex-col basis-3/4">
		{#if $subtitles.isLoading}
			<div class="flex justify-center items-center">
				<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-gray-900" />
			</div>
		{:else if $subtitles.isError}
			<div class="flex justify-center items-center">
				<div class="rounded-full h-8 w-8 border-t-2 border-b-2 border-red-500" />
			</div>
		{:else if $subtitles.isSuccess}
			<SubtitleList subtitles={$subtitles.data} />
		{/if}
	</div>

	{#if subtitleId}
		<div class="flex-grow basis-1/4">
			<VideoPlayer {subtitleId} width={360} />
		</div>
	{/if}
</main>
