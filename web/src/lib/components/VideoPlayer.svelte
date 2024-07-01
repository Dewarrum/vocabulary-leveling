<script lang="ts">
	import { PUBLIC_API } from '$env/static/public';
	import videojs from 'video.js';
	import 'video.js/dist/video-js.css';

	export let subtitleId;

	let videoRef: HTMLDivElement;
	const videoElement = document.createElement('video-js');
	videoElement.classList.add('w-full');
	videoElement.classList.add('h-auto');
	videoElement.classList.add('aspect-video');

	$: player = videojs(videoElement, {
		type: 'application/x-mpegURL'
	});

	$: {
		if (videoRef) {
			videoRef.appendChild(videoElement);
		}

		if (player) {
			player.src(`${PUBLIC_API}/api/videos/manifest.mpd?subtitleId=${subtitleId}`);
			// player.src(`/manifest.mpd`);
			player.controls(true);
			player.height('auto');
			player.volume(0.1);
		}
	}
</script>

<div data-vjs-player bind:this={videoRef}></div>
