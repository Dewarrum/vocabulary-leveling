<script lang="ts">
	import { PUBLIC_API } from '$env/static/public';
	import videojs from 'video.js';
	import 'video.js/dist/video-js.css';

	export let subtitleId;
	export let width: number | string | undefined = undefined;
	export let height: number | string | undefined = undefined;

	let videoRef: HTMLDivElement;
	const videoElement = document.createElement('video-js');

	$: player = videojs(
		videoElement,
		{
			type: 'application/x-mpegURL',
			html5: {
				dash: {
					useTTML: true
				}
			}
		},
		() => {
			videojs.log('player is ready');
		}
	);

	$: {
		if (videoRef) {
			videoRef.appendChild(videoElement);
		}

		if (player) {
			player.src(`${PUBLIC_API}/api/videos/manifest.mpd?subtitleId=${subtitleId}`);
			player.controls(true);
			player.width(width);
			player.height(height);
		}
	}
</script>

<div data-vjs-player>
	<div bind:this={videoRef} />
</div>
