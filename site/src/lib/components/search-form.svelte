<script lang="ts">
	let { query = '' } = $props();
</script>

<form method="POST" action="?/search">
	<label>
		Query
		<input name="query" type="text" value={query} />
	</label>

	<input type="submit" id="submit" value="Search" />
	<button formaction="?/plucky">I'm feeling plucky!</button>
</form>
