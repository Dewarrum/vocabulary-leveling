<script lang="ts">
	function onSignIn() {
		window.location.href = '/auth/sign-in';
	}
</script>

<div class="flex flex-col items-center gap-4 rounded-md border px-16 py-8 shadow-md">
	<h1 class="text-3xl">Sign in via Logto</h1>
	<button class="rounded-md border px-4 py-2" on:click={onSignIn}>Sign in</button>
</div>
