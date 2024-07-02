<script lang="ts">
	import { page } from '$app/stores';
	import { searchSubtitles, type Subtitle } from '$lib/api/subtitles';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	$: searchQuery = $page.url.searchParams.get('query') ?? '';
	$: subtitles = createQuery({
		queryKey: ['subtitles', searchQuery],
		queryFn: () => searchSubtitles(searchQuery)
	});

	let selectedSubtitleIndex = 0;
	let selectedSubtitle: Subtitle | undefined;
	$: {
		if ($subtitles.isSuccess && $subtitles.data.length > 0) {
			selectedSubtitle = $subtitles.data[selectedSubtitleIndex];
		}
	}

	function onNextSubtitleSelected() {
		if (!$subtitles.isSuccess) return;

		if (selectedSubtitleIndex < $subtitles.data.length - 1) {
			selectedSubtitleIndex++;
		}
	}

	function onPreviousSubtitleSelected() {
		if (!$subtitles.isSuccess) return;

		if (selectedSubtitleIndex > 0) {
			selectedSubtitleIndex--;
		}
	}
</script>

<main class="mx-auto flex w-3/4 flex-col gap-4 pt-8">
	{#if $subtitles.isSuccess && selectedSubtitle}
		<div class="flex flex-row">
			<button
				class="rounded-md border border-gray-500 px-4 py-2 disabled:opacity-50"
				disabled={selectedSubtitleIndex === 0}
				on:click={onPreviousSubtitleSelected}>Previous</button
			>
			<span class="flex-grow"></span>
			<button
				class="rounded-md border border-gray-500 px-4 py-2 disabled:opacity-50"
				disabled={selectedSubtitleIndex === $subtitles.data.length - 1}
				on:click={onNextSubtitleSelected}>Next</button
			>
		</div>
		<div class="w-full">
			<VideoPlayer subtitle={selectedSubtitle} />
		</div>
		<div>
			<span class="text-2xl">{selectedSubtitle.text}</span>
		</div>
	{/if}
</main>
