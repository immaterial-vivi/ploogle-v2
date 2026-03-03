<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css'
	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Ploogle - Search HDG stories</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</svelte:head>

{@render children()}

<style>
	body :global{

	}

</style>