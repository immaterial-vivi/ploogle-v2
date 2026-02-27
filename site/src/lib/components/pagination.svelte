<script lang="ts">
	let { pageNr, totalPages, baseUrl, pageParamName } = $props();

	let pages = [];
	for (
		var i = Math.max(1, pageNr - 5);
		i <= Math.min(Math.max(1, pageNr - 5) + 9, totalPages);
		i++
	) {
		pages.push({
			num: i,
			active: pageNr == i
		});
	}

	const hasNextPage = pageNr < totalPages;
	const hasPrevPage = pageNr > 1;

	let nextPage = $state(1);
	let prevPage = $state(1);
	if (hasNextPage) {
		nextPage = pageNr + 1;
	}
	if (hasPrevPage) {
		prevPage = pageNr - 1;
	}
</script>

{#if totalPages > 1}
	<ol class="pages">
		{#if hasPrevPage}
			<a title="go to previous page" class="pageButton" href="{baseUrl}&{pageParamName}={prevPage}">
				<li>&lt;</li>
			</a>
		{/if}
		{#each pages as page}
			<li class="page">
				{#if page.active}
					{page.num}
				{:else}
					<a href="{baseUrl}&{pageParamName}={page.num}">
						{page.num}
					</a>
				{/if}
			</li>
		{/each}
		{#if hasNextPage}
			<a title="go to next page" class="pageButton" href="{baseUrl}&{pageParamName}={nextPage}">
				<li>&gt;</li>
			</a>
		{/if}
	</ol>
{/if}

<style>
	.pages {
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
		list-style: none;
		gap: 1rem;
	}
</style>