<script lang="ts"></script>

<main class="mx-auto flex w-3/4 flex-col items-center pt-8">
	<form
		class="flex flex-col gap-4"
		action="/api/admin/videos/upload"
		method="POST"
		enctype="multipart/form-data"
	>
		<label for="videoFile">Select video</label>
		<input id="videoFile" type="file" name="video" />

		<label for="videoName">Video name</label>
		<input
			id="videoName"
			type="text"
			name="videoName"
			class="rounded-md border-2 border-gray-300 px-4 py-2"
		/>

		<label for="subtitlesFile">Select subtitles</label>
		<input id="subtitlesFile" type="file" name="subtitles" />

		<button class="rounded-md border px-4 py-2">Upload</button>
	</form>
</main>
