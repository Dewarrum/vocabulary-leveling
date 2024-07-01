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

	let subtitleId = '';
	$: {
		if ($subtitles.isSuccess) {
			subtitleId = $subtitles.data[0].id;
		}
	}
</script>

<main class="flex flex-row gap-4 pt-8 w-3/4 mx-auto">
	{#if subtitleId}
		<div class="w-full mx-auto">
			<VideoPlayer {subtitleId} />
		</div>
	{/if}
</main>
