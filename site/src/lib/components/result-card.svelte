<script lang="ts">

    import romLogo from '$lib/assets/rom-icon.png';
    import ao3Logo from '$lib/assets/ao3-icon.png';

    const {result} = $props();

    const chapterUrl = () => getChapterUrl(result)

    function selectIconForUrl() {
        if (chapterUrl().startsWith('https://readonlymind.com')) {
            return {n: "rom", l: romLogo};
        }
        return {n: "ao3", l: ao3Logo};
    }

    function getChapterUrl(result) {
        return result.Chapter ? result.Chapter_Url : result.Book_Url
    }

    function shortChapterUrl() {
        const url = new URL(chapterUrl())

        return url.pathname

    }

</script>
    <section class="card {result.Blacklisted ? 'blacklisted': ''} {result.Direct_Title_Match ? "direct-match" :""}">
        <span class="sr-only">⮦</span>
        <a aria-hidden="true" class="result-header" href={chapterUrl()} target="_blank"
           rel="noopener noreferrer">
            <img class="source-icon" alt="{selectIconForUrl(result).n}-icon" src={selectIconForUrl(result).l}/>
            <span class="font-boring smol url-hint">{shortChapterUrl()}</span>
        </a>
        <a href={chapterUrl()} target="_blank" rel="noopener noreferrer" class="title">
            <h2>{result.Title}{result.Chapter ? `, Chapter ${result.Chapter}` : ''}</h2></a
        >
        {#if result.Blacklisted}
        <span class="blacklist-reason">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                 stroke-width="1.5" stroke="currentColor" class="size-6">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>
            </svg>
            {result.Blacklist_Reason}</span>
        {:else}

            <span class="author font-boring">
                <span class="sr-only">by</span>
                {result.Author}
            </span>
            <p class="summary"><span class="sr-only">Summary: &nbsp; </span>{result.Summary}</p>
            {#if !result.Direct_Title_Match}
            <span class="sr-only">Most relevant section:</span>
            <blockquote class="excerpt">[...]{@html result.Excerpt}[...]</blockquote>
            {/if}
        {/if}
    </section>
<style>

    h2 {
        font-size: 1.75rem;
        font-family: "Fugaz One";
    }

    .result-header {
        display: flex;
        align-items: center;
        gap: .15rem;
    }

    .source-icon {
        width: 1.4rem;
        object-fit: contain;
        object-position: center;
    }

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
        /*border: 1px solid var(--a-blue-low);*/
        border-radius: 0.5rem;
    }

    .blacklisted {
        filter: grayscale(1);
        font-size: smaller;

        > a > h2 {
            font-size: 1rem !important;
        }
    }

    .blacklist-reason {
        display: flex;
        flex-direction: row;
        gap: .5rem;
        align-items: center;

        > * {
            height: 1.5rem;
        }
    }

    .author {
    }

    .author::before {
        content: 'by ';
        color: grey;
    }

    .excerpt :global {
        & > em {
            color: var(--a-yellow);
            font-weight: bold;
        }
    }

    .excerpt {
        color: lightgray;
        position: relative;
        margin-inline-start: .5rem;
        border-inline-start: 2px solid rgba( 0 0 0/ 0 );
        z-index: 2;
        padding-block: .5rem;
        padding-inline: 1rem 1rem;
        border-radius: 0 0 1rem 0;
        corner-shape: bevel;
        background: oklch(from var(--a-light-blue) l c h / 0.1);
        clip: unset;
        /*clip-path: polygon(120% -10%, 120% calc(100% - 16px), calc(100% - 16px) 100%, -10% 100%, -10% -10%);*/
    }

    .excerpt::before {
        content: '';
        position: absolute;
        left: -5px;
        top: 0;
        bottom: 0;
        background: var(--a-light-blue);
        clip-path: polygon(5px 0%, 100% 0%, 100% 100%, 0% 100%, 0% 7px);
        width: 5px;
        /*border-radius: 10px 0 0 10px;*/
    }

    .excerpt::after {
        content: '';
        position: absolute;
        right: -5px;
        top: -10px;
        background: url("/img/quotation-marks.svg") no-repeat center center / contain;
        width: 30px;
        height: 30px;
        /*border-radius: 10px 0 0 10px;*/
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
    /*.url-hint {*/
    /*    display: none;*/
    /*}*/

    .direct-match {
        /*border: var(--a-yellow) solid 1px;*/
        /*box-shadow: var(--a-yellow) 0px 0px 1px;*/
        & .title {
            display: flex;
            align-items: center;
            gap: .5rem;
        }
        & .title::before {
            content: "★";
            color: var(--a-yellow);
            font-size: larger;
        }
    }

    @media screen and (min-width: 768px) {
        .card {
            padding: 0 1rem;
        }

        .url-hint {
            display: block;
        }
    }


</style>