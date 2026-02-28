<script lang="ts">
	const { result } = $props();
</script>

<div class="card">
	<a href={result.Chapter ? result.Chapter_Url : result.Book_Url}>
		<h2>{result.Title}{result.Chapter ? `, Chapter ${result.Chapter}` : ''}</h2></a
	>
	<span class="author font-boring">{result.Author}</span>
	<p class="summary">{result.Summary}</p>

	<blockquote class="excerpt">{@html result.Excerpt}</blockquote>
	<span class="font-boring smol url-hint">{result.Chapter ? result.Chapter_Url : result.Book_Url}</span>
</div>

<style>
	a {
		text-decoration: underline 0.15em rgba(0, 0, 0, 0);
		transition: text-decoration-color 100ms;
		cursor: pointer;
		color: var(--a-light-blue);
		&:visited {
			color: var(--a-mid-teal);
		}
		&:hover:not(:visited) {
			text-decoration: underline 0.15em var(--a-light-blue);
		}
		&:hover:visited {
			text-decoration: underline 0.15em var(--a-mid-teal);
		}
	}

	.card {
		display: flex;
		flex-direction: column;
		gap: .5rem;
	}



	.author {
	}
	.author::before {
		content: 'by ';
		color: oklch(from var(--a-light-blue) l a b / 1);
	}

	.excerpt :global{
		& > em {
			color: var(--a-yellow);
			font-weight: bold;
		}
	}

	.excerpt {
		color: lightgray;
		position: relative;
		padding-inline-start: .5rem;
		margin-inline-start: .5rem;
		border-inline-start: 2px solid var(--a-light-blue);
		z-index: 2;
		padding-block: .5rem;
		border-radius: 0 1rem 0 0;
		background: oklch( from var(--a-light-blue) l c h / 0.1);
	}

	.excerpt::before {
		content: '';
		position: absolute;
		left: -5px;
		top: 0;
		bottom: 0;
		background: var(--a-light-blue);
		width: 5px;
		border-radius: 10px 0 0 10px;
	}

	/*.excerpt::after {*/
	/*	content: '\"';*/
	/*	position: absolute;*/
	/*	top: -5rem;*/
	/*	right: .5rem;*/
	/*	font-size: 18rem;*/
	/*	font-family: "Fugaz One";*/
	/*	color: oklab( from var(--a-blue-low) l a b / 1);*/
	/*	z-index: -1;*/
	/*}*/
	.url-hint {
		display: none;
	}
	@media screen and (min-width: 768px) {
		.url-hint {
			display: block;
		}
	}

</style>