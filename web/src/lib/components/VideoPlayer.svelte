<script lang="ts">
	import videojs from 'video.js';
	import 'video.js/dist/video-js.css';

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
			player.src(
				'http://localhost:3000/api/videos/manifest.mpd?videoId=82a926c9-1b5e-4d6a-8e09-1e92bfe38bab'
				// '/video.mpd'
			);
			player.controls(true);
		}
	}
</script>

<div data-vjs-player>
	<div bind:this={videoRef} />
</div>
