<script lang="ts">
    import Footer from '$lib/components/footer.svelte';
    import Pagination from '$lib/components/pagination.svelte';
    import ResultCard from '$lib/components/result-card.svelte';
    import HeaderSearchForm from '$lib/components/header-search-form.svelte';
    import type {PageProps} from './$types';
    import {page} from "$app/state";

    let {data}: PageProps = $props();

    function hasHiddenItems() {
        const hasHidden = data.message.Hits.some(v => v.Blacklisted)
        console.log("has hidden items:", hasHidden)
        return hasHidden
    }
</script>

<div class="height-container">
    <header class="header ">
        <a href="#main-content" class="sr-only skip-link">skip to search results</a>
        <div class="header-content container">
            <a href="/" class="logo logo-hover interactive">Ploogle</a>
            <HeaderSearchForm query={data.message.Query}/>
            <a href="https://humandomestication.guide" class="wiki-link font-boring phone-hidden">back to hd.g</a>
        </div>
    </header>

    <main class="main-content" id="#main-content">
        <div class="container content">
        <span class="results-count">found {data.message.Page.Results}
            results in {(data.message.Performance.DeltaTime / 1e9).toFixed(3)} s </span>

            <ol class="results-list">
                {#each data.message.Hits as hit}
                    {#if !hit.Blacklisted || data.showHidden}
                        <li class="results-item">
                            <ResultCard result={hit}/>
                        </li>
                    {/if}
                {/each}
            </ol>
            {#if hasHiddenItems() && !data.showHidden}
                <form method="post" action="?/showHidden" class="show-hidden-form">
                    <span> Some results were skipped</span>
                    <input type="hidden" name="target" value={page.url}/>
                    <input type="submit" id="showHidden" value="Show skipped results"/>
                </form>
            {/if}

            {#if data.message.Hits.length == 0}
                <section>
                    <h2>
                        No results found
                    </h2>

                    <span>
                        maybe you should write something about <em>{data.message.Query}</em> :3
                    </span>
                </section>
            {/if}

            <Pagination
                    baseUrl={`/search?q=${data.message.Query}`}
                    pageParamName={'p'}
                    totalPages={data.message.Page.Pages}
                    pageNr={data.message.Page.Page}
            />
        </div>


    </main>
    <Footer/>
</div>
<style>
    ol {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        list-style: none;
        padding-bottom: 1rem;
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
            padding: 1rem 0;
        }
    }


    h2 {
        padding-bottom: 1rem;
    }

    .main-content {
        position: relative;
        z-index: 0;
    }

    .container {
        max-width: var(--content-width);
        margin: 0 auto;
        padding: 1rem 1rem 3rem 1rem;
        width: 100%;
        height: 100%;
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

    .show-hidden-form {
        display: flex;
        flex-direction: column;
        max-width: 30rem;
        margin: 0 auto;
        padding: 1rem 1rem 1rem 1rem;
        margin-bottom: 2rem;
        border: 1px solid var(--a-blue-low);
        border-radius: 1rem 1rem 0 1rem;
        background-color: var(--space-blue-low);
        gap: .75rem;
        position: relative;

        &::after {
            content: " ";
            position: absolute;
            right: 1rem;
            height: 40px;
            width: 40px;
            background: url("/img/info.svg") no-repeat center center  / 32px 32px;
        }

        & > input {
            border: none;
            background: none;
            text-align: left;
            cursor: pointer;
            text-decoration: underline;
            text-decoration-color: rgba(from white r g b / 0.4);
        }

    }


    @media screen and (min-width: 768px) {
        .height-container {
            display: grid;
            grid-template-rows: auto 1fr auto;
            justify-items: stretch;
            width: 100%;
            min-height: 100svh;
        }

        .header-content {
            flex-direction: row;

            & > * {
                padding: unset;
            }
        }

        .container {
            padding: 1rem 1rem;
        }


    }

</style>