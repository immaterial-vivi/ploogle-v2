<script lang="ts">
    import SearchIcon from '$lib/components/search-icon.svelte';

    let {query = ''} = $props();

</script>

<form method="POST" action="?/search">
    <section id="form-input">

        <label>
            <span class="sr-only">search query</span>
            <input name="query" type="text" value={query} placeholder="Search for HDG stories..."/>
        </label>

        <button type="submit" id="submit" class="search-button">
            <SearchIcon  class="button-icon"/>
            <span class="sr-only">Search</span>
        </button>
        <!--	<button formaction="?/plucky">I'm feeling plucky!</button>-->

<!--        <div class="hint-button" tabindex="0">-->
<!--            <svg xmlns="http://www.w3.org/2000/svg" fill="#040B20" viewBox="0 0 24 24" stroke-width="1.5"-->
<!--                 stroke="currentColor"-->
<!--                 class="icon" width="24" height="24">-->
<!--                <path stroke-linecap="round" stroke-linejoin="round"-->
<!--                      d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"/>-->
<!--            </svg>-->
<!--        </div>-->

<!--        <ul class="hint">-->
<!--            <span class="sr-only">Tips for more accurate results:</span>-->
<!--            <li>-->
<!--                <span>"quoted text"</span> to search for phrases-->
<!--            </li>-->
<!--            <li>-->
<!--                <span>unquoted text</span> to find stories containing both 'unquoted' <em>and</em> 'text'-->
<!--            </li>-->
<!--            <li>-->
<!--                put <span>OR</span> between terms to find stories containing <em>either</em> term-->
<!--            </li>-->
<!--            <li>-->
<!--                <span>-word</span> exclude stories containing 'word'-->
<!--            </li>-->
<!--        </ul>-->
    </section>
    <section id="form-actions">
        <button id="plucky" formaction="?/plucky" type="submit" class="action-button">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z" />
            </svg>

            <span>I'm feeling plucky</span>
        </button>

        <button  type="submit"  class="action-button">

            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="icon">
                <path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
            </svg>
            <span>Search plampts!!</span>
        </button>
    </section>



</form>

<style>
    form {
        /*max-width: 600px;*/
        margin: 0 auto;
    }
    #form-input {
        display: grid;
        grid-template-columns: 1fr auto;
        align-items: center;
        justify-content: center;
        /*background: #1aa179;*/
        background: rgba(250 250 250 /  0.1);
        border-radius: 0.25rem;
        backdrop-filter: blur(20px);
        box-shadow: inset 0px 2px 2px 0px rgba(0 0 0 /.15), 0px 4px 4px 4px rgba(0 0 0 /.1);
        transition: all .1s ease-out;
        &:focus-within {
            background: rgba(250 250 250 /  0.10);
            box-shadow: inset 0px 2px 2px 0px rgba(0 0 0 /.15), 0px 4px 16px 4px rgba(0 0 0 /.3);

        }
    }

    input[type=text] {
        position: relative;
        border: none;
        background: none;
        width: 100%;
        font-size: 1.25rem;
        color: var(--a-yellow);
        margin: 0.5rem 0.5rem 0.5rem 1.5rem;
        &::placeholder {
            color: rgba(from var(--a-yellow) r g b / 0.5);
            font-size: 1.25rem;

        }
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

    #form-actions {

        display: flex;
        flex-direction: column-reverse;
        grid-template-columns: 1fr;
        max-width: 600px;
        margin: 0 auto;
        gap: 2rem;
        padding: 2rem 0;
    }

    @media screen and (min-width: 768px) {
        #form-actions {
            display: grid;
            grid-template-columns: 1fr 1fr;
        }

        input[type=text] {
            width: 600px;
            min-height: 2.5rem;
            font-size: 1.25rem;
        }
    }

    .action-button{
        position: relative;
        display: grid;
        grid-template-columns: auto 1fr;
        border-radius: 4px;
        font-size: larger;
        padding: 0.75rem 1rem;
        gap: .5rem;
        cursor: pointer;
        background: rgba(250 250 250 /  0.05);

        border: none;

        box-shadow: inset 0px 2px 2px 0px rgba(0 0 0 /.15), 0px 0px 0px 0px rgba(0 0 0 /.2);

        transition: 0.1s all ease-out;
        & > svg {
            height: 1.5rem;
            color: var(--a-yellow);
        }
        &:hover, &:focus{
            & > svg {
                color: white;
            }
            background: rgba(250 250 250 /  0.10);
            transform: scale(1.1);
            box-shadow: inset 0px 2px 2px 0px rgba(0 0 0 /.15), 0px 4px 16px 4px rgba(0 0 0 /.3);

        }
    }

    /*.search-button {*/
    /*    opacity: 0.0;*/
    /*}*/

</style>