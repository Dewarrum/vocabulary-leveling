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
				'http://localhost:3000/api/videos/manifest.mpd?subtitleCueId=7f8d9ef9-91b0-4fb9-95f7-9a39c3bc51dc'
				// '/video.mpd'
			);
			player.controls(true);
		}
	}
</script>

<div data-vjs-player>
	<div bind:this={videoRef} />
</div>
