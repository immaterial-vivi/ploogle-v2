<script lang="ts">
    import Footer from '$lib/components/footer.svelte';
    import Pagination from '$lib/components/pagination.svelte';
    import ResultCard from '$lib/components/result-card.svelte';
    import HeaderSearchForm from '$lib/components/header-search-form.svelte';
    import type {PageProps} from './$types';

    let {data}: PageProps = $props();
</script>

<header class="header ">
    <a href="#main-content" class="sr-only skip-link">skip to search results</a>
    <div class="header-content container">
        <a href="/" class="logo">Ploogle</a>
        <HeaderSearchForm query={data.message.Query}/>
        <a href="https://humandomestication.guide" class="wiki-link font-boring">back to hd.g</a>
    </div>
</header>

<main class="main-content" id="#main-content">
    <div class="container content">
        <span class="results-count">found {data.message.Page.Results}
            results in {(data.message.Performance.DeltaTime / 1e9).toFixed(3)} s </span>

        <ol class="results-list">
            {#each data.message.Hits as hit}
                <li class="results-item">
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
        text-decoration: none;
        cursor: pointer;
        font-family: Fugaz One, sans, sans-serif;
        font-weight: bold;
        font-size: 2rem;
        color: var(--affini-pink-main);
        transition: all 0.2s;

        &:hover {
            text-decoration: underline;
            color: white;

        }
    }

    .header {
        background: url("/img/banner.gif");
        border-bottom: 1px var(--affini-pink-main) solid;
        padding: 0 1rem;
        z-index: 9;
    }

    .header-content {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        & > * {
            padding:  1rem 0;
        }
    }

    .main-content {
        position: relative;
        z-index: 0;
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 1rem 1rem 3rem 1rem;
        width: 100%;
    }

    .content {
        background: rgba(4 11 32 / .95);
    }

    .results-count {
        position: relative;
        display: inline-block;
        border: 1px solid var(--a-light-blue);
        background: oklab(from var(--a-light-blue) l a b / 0.1);
        padding: 0.5rem 1rem;
        margin: 1rem auto;
        text-align: center;
        color: var(--a-light-blue);
        border-radius: 1rem 1rem 0 1rem;
    }


    .results-list {
        padding-bottom: 2rem;
    }

    .results-list > *:not(:last-child)::after {
        content: "";
        position: absolute;
        bottom: -2px;
        width: 8rem;
        margin: auto;
        left: 0;
        right: 0;
        height: 24px;
        background: /*url("/img/favicon.svg") no-repeat center center / 24px 24px,*/ linear-gradient(var(--a-blue-low)) no-repeat 0 center / 2rem 2px,
        linear-gradient(var(--a-blue-low)) no-repeat 6rem center / 2rem 2px;
        /*            radial-gradient(*/
        /*                    var(--a-blue-low) 0 1px,*/
        /*                rgba( 0 0 0 / 0) 1px 4px*/
        /*            )  center center / 4px 4px;*/
    }

    .results-list > *:not(:last-child)::before {
        content: "";
        position: absolute;
        bottom: -0px;
        width: 48px;
        margin: auto;
        left: 0;
        right: 0;
        height: 18px;
        clip-path: polygon(25% 0%, 49% 0%, 23% 100%, 34% 100%, 61% 0%, 76% 0%, 47% 100%, 61% 100%, 90% 0%, 100% 0%, 73% 100%, 50% 100%, 13% 100%, 0% 100%);
        background-color: var(--a-blue-low);

    }


    .results-list > *:not(:last-child) {
        position: relative;
        padding-bottom: 2rem;
    }

    .wiki-link::before {
        content: "";
        position: absolute;
        top:   0;
        left: -2rem;
        width: 2rem;
        bottom: 0;
        background: url("/img/favicon.svg") no-repeat center center / contain;
    }
    .wiki-link {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        margin: 0 ;
        max-width: 20rem;
        color: rgba(0 0 0 / 0 );
        text-decoration: underline rgba(255 255 255 / 0);
        transition: text-decoration-color 100ms;
        &:hover {
            text-decoration: underline rgba(255 255 255 / 1);
        }
    }

    @media screen and (min-width: 768px) {

        .header-content {
            flex-direction: row;
            & > * {
                padding:  unset;
            }
        }

        .container {
            padding: 1rem 1rem;
        }
        .wiki-link {
            position: relative;
            color: var(--a-yellow);
            width: unset;

        }
    }

</style>