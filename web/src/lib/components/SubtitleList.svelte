<script lang="ts">
	import { page } from '$app/stores';
	import type { Subtitle } from '$lib/api/subtitles';

	export let subtitles: Subtitle[];

	function renderPreviewVideoLink(subtitleId: string) {
		const query = new URLSearchParams($page.url.search);
		query.set('subtitleId', subtitleId);
		return `/search?${query.toString()}`;
	}
</script>

<div class="flex flex-col gap-4">
	{#each subtitles as subtitle}
		<div class="flex gap-2 items-center">
			<div class="flex flex-col">
				<div class="text-sm font-medium">
					<h2 class="text-lg">
						<a href={renderPreviewVideoLink(subtitle.id)}>{subtitle.text}</a>
					</h2>
				</div>
			</div>
		</div>
	{/each}
</div>
