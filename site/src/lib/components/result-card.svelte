<script lang="ts">
	const { result } = $props();
</script>

<div>
	<a href={result.Chapter ? result.Chapter_Url : result.Book_Url}>
		<h2>{result.Title}{result.Chapter ? `, Chapter ${result.Chapter}` : ''}</h2></a
	>
	<span>{result.Author}</span>
	<p>{result.Excerpt}</p>
	<span>{result.Chapter ? result.Chapter_Url : result.Book_Url}</span>
</div>
