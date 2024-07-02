<script lang="ts">
	import type { Subtitle } from '$lib/api/subtitles';
	import videojs from 'video.js';
	import 'video.js/dist/video-js.css';

	export let subtitle: Subtitle;

	let videoRef: HTMLDivElement;
	const videoElement = document.createElement('video-js');
	videoElement.classList.add('w-full');
	videoElement.classList.add('h-auto');
	videoElement.classList.add('aspect-video');

	let player = videojs(videoElement, {
		type: 'application/x-mpegURL',
		html5: {
			nativeTextTracks: false
		}
	});

	let textTrack = player.addTextTrack('subtitles', 'korean', 'kr');

	$: {
		if (videoRef) {
			videoRef.appendChild(videoElement);
		}
	}

	$: {
		if (player) {
			player.src(`/api/videos/manifest.mpd?subtitleId=${subtitle.id}`);
			player.controls(true);
			player.height('auto');
			player.volume(0.1);
		}
	}

	$: {
		if (textTrack) {
			textTrack.mode = 'showing';

			let cue = textTrack.cues?.getCueById('cue');
			if (cue) {
				textTrack.removeCue(cue);
			}

			cue = new VTTCue(0, 1000, subtitle.text);
			cue.id = 'cue';
			textTrack.addCue(cue);
		}
	}
</script>

<div data-vjs-player bind:this={videoRef}></div>
