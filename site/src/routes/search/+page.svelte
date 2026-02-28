<script lang="ts">
    import Footer from '$lib/components/footer.svelte';
    import Pagination from '$lib/components/pagination.svelte';
    import ResultCard from '$lib/components/result-card.svelte';
    import HeaderSearchForm from '$lib/components/header-search-form.svelte';
    import type {PageProps} from './$types';

    let {data}: PageProps = $props();
</script>

<header class="header">
    <a href="#main-content" class="sr-only skip-link">skip to search results</a>
    <div class="header-content">
        <span class="logo">Ploogle</span>
        <HeaderSearchForm query={data.message.Query}/>
        <a href="https://humandomestication.guide">back to wiki</a>
    </div>
</header>

<main class="main-content" id="#main-content">
    <div class="container">
        <ol>
            {#each data.message.Hits as hit}
                <li>
                    <ResultCard result={hit}/>
                </li>
            {/each}
        </ol>
        <Pagination
                baseUrl={`/search?q=${data.message.Query}`}
                pageParamName={'p'}
                totalPages={data.message.Page.Pages}
                pageNr={data.message.Page.Page}
        />
    </div>
</main>
<span>total results: {data.message.Page.Results}</span>
<span>database took: {(data.message.Performance.DeltaTime / 1e9).toFixed(3)}s</span>


<Footer/>

<style>
    ol {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        list-style: none;
        padding-bottom: 1rem;
    }

    .logo {
        float: left;
    }

    .header {
        background: url("/img/banner.gif");
        height: 5rem;
        border-bottom: 1px var(--affini-pink-main) solid;
        padding: 0 1rem;
    }

    .header-content {
        height: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem 1rem;
        width: 100%;
        background: rgba(4 11 32 / .9);
    }
</style>