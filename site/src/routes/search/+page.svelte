<script lang="ts">
	import Footer from '$lib/components/footer.svelte';
	import Pagination from '$lib/components/pagination.svelte';
	import ResultCard from '$lib/components/result-card.svelte';
	import SearchForm from '$lib/components/search-form.svelte';
	import type { PageProps } from './$types';
	let { data }: PageProps = $props();
</script>

<SearchForm query={data.message.Query} />

<ol>
	{#each data.message.Hits as hit}
		<li><ResultCard result={hit} /></li>
	{/each}
</ol>
<span>total results: {data.message.Page.Results}</span> |
<span>database took: {(data.message.Performance.DeltaTime / 1e9).toFixed(3)}s</span>

<Pagination
	baseUrl={`/search?q=${data.message.Query}`}
	pageParamName={'p'}
	totalPages={data.message.Page.Pages}
	pageNr={data.message.Page.Page}
/>

<Footer />
