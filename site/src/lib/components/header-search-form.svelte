<script lang="ts">
    import SearchIcon from '$lib/components/search-icon.svelte';

    let {query = ''} = $props();

</script>

<form method="POST" action="?/search">
    <label>
        <span class="sr-only">search query</span>
        <input name="query" type="text" value={query}/>
    </label>

    <button type="submit" id="submit" class="search-button">
        <SearchIcon class="button-icon"/>
        <span class="sr-only">Search</span>
    </button>
    <!--	<button formaction="?/plucky">I'm feeling plucky!</button>-->

    <div class="hint-button" tabindex="0">
        <svg xmlns="http://www.w3.org/2000/svg" fill="#040B20" viewBox="0 0 24 24" stroke-width="1.5"
             stroke="currentColor"
             class="icon" width="24" height="24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"/>
        </svg>
    </div>

    <ul class="hint">
        <span class="sr-only">Tips for more accurate results:</span>
        <li>
            <span>"quoted text"</span> to search for phrases
        </li>
        <li>
            <span>unquoted text</span> to find stories containing both 'unquoted' <em>and</em> 'text'
        </li>
        <li>
            put <span>OR</span> between terms to find stories containing <em>either</em> term
        </li>
        <li>
            <span>-word</span> exclude stories containing 'word'
        </li>
    </ul>

</form>

<style>
    form {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        /*background: #1aa179;*/
        background: rgba(250 250 250 /  0.10);
        border-radius: 0.25rem;
        backdrop-filter: blur(20px);
        box-shadow: inset 0px 4px 4px 0px rgba(0 0 0 /.3);
        box-shadow: 0px 4px 16px 4px rgba(0 0 0 /.3);

    }

    input[type=text] {
        position: relative;
        border: none;
        background: none;
        font-size: 1.25rem;
        color: var(--a-yellow);
        margin: 0.5rem 0.5rem 0.5rem 1.5rem;
    }

    input[type=text]:focus {
        outline: none;
    }

    label::after {
        content: '';
        display: block;
        position: absolute;
        bottom: .35rem;
        left: .75rem;
        right: 4.25rem;
        margin: auto;
        background: var(--a-yellow);
        height: 1px;
        transition: transform 0.1s ease-out;
        transform: scaleX(0) scaleY(0);
        transform-origin: center bottom;

    }

    label:focus-within::after {
        background: var(--a-yellow);
        height: 1px;
        transform: scaleX(1.0) scaleY(1.0);
    }


    #submit {
        border: none;
        background: none;
        color: var(--a-yellow);
        display: flex;
        margin: 0.5rem 1rem 0.5rem 0;
        cursor: pointer;
        padding: 0.25rem 0.5rem;
        transition: 0.1s all ease-out;

    }

    #submit:hover {
        color: white;
        transform: scale(1.2);
    }

    .hint-button {
        color: var(--a-yellow);
        position: absolute;
        right: 1.5rem;
        top: -3.5rem;

    }

    .hint-button:hover ~ .hint,
    .hint-button:focus ~ .hint {
        opacity: 1;
        visibility: visible;

    }

    .hint {
        color: white;
        position: absolute;
        right: 0;
        top: 0rem;
        display: block;
        border: 1px solid var(--a-light-blue);
        background: var(--space-blue-low);
        padding: 0.5rem 1rem;
        font-size: medium;
        max-width: 30rem;
        opacity: 0;
        transition: opacity 0.2s ease-out;
        box-shadow: 0px 4px 4px 4px rgba(0 0 0 /.6);
        visibility: hidden;
        z-index: 10;

        > li {
            margin-inline-start: 0.5rem;
            line-height: 1.8;
            list-style-type: "-";

            > span {
                background-color: rgba(255 255 255 / 0.1);
                padding: 0.1rem 0.25rem;
                border-radius: 0.25rem;

            }

            > em {
                font-weight: bold;
                color: oklab(from var(--a-yellow) 10 a b / 1);
            }

        }
    }


    @media screen and (min-width: 768px) {
        .hint {
            right: -18rem;
            top: 3rem;
        }

        .hint-button {
            color: var(--a-yellow);
            position: absolute;
            right: -3rem;
            top: 0.7rem;
            left: unset;
        }
    }


</style>