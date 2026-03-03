<script>
    import SearchForm from '$lib/components/search-form.svelte';
    import SearchHint from '$lib/components/search-hint.svelte';

    import Flort from '$lib/components/flort.svelte';
    import Footer from '$lib/components/footer.svelte';
</script>

<div class="container">
    <!--    <FlashMessage>-->
    <!--            Welcome to the new Ploogle! <br> Read all about it <a href="#" target="_blank" rel="nofollow">here!</a>-->
    <!--    </FlashMessage>-->
    <main id="search">
        <div class="heading-container">
            <h1 class="heading">Ploogle</h1>
            <Flort/>
            <p class="sr-only">
                Ploogle - Search for <a href="https://humandomestication.guide/">HDG Stories</a>
            </p>
        </div>
        <SearchForm/>
    </main>
    <section id="faq">
        <ul>
            <h2>FAQ</h2>
            <li>
                <details >
                    <summary>
                        OwO What's this?
                    </summary>
                    <p>Ploogle finds HDG stories. Ploogle has it for you</p>
                </details>
            </li>

            <li>
                <details >
                    <summary>
                        Search options
                    </summary>
                    <SearchHint/>
                </details>
            </li>
            <li>
                <details >
                    <summary>
                        Which stories are indexed?
                    </summary>
                    <p>
                        New stories posted on AO3 under the Human Domestication Tag are picked up automatically about
                        once
                        per day. All #hdg stories on ReadOnlyMind.com are also indexed.
                    </p>
                </details>

            </li>
            <li>
                <details >
                    <summary>
                        I found a broken link/something inappropriate!
                    </summary>
                    <p>
                        Shoot me a message on tumblr or at vivi[at]humandomestication.guide and I'll look into it
                    </p>
                </details>
            </li>
            <li>
                <details >
                    <summary>
                        My story doesn't show up!
                    </summary>
                    <p>
                        Shoot me a message on tumblr or at vivi[at]humandomestication.guide and I'll look into it
                    </p>
                </details>
            </li>
            <li>
                <details >
                    <summary>
                        Why not just use the ao3 search?
                    </summary>
                    <p>
                        Ploogle searches the <em>full text</em> of all stories simultaneously. If you're looking for
                        specific tags or want see stories ordered by release date or other statistics, Ao3's search is
                        where it's at!
                    </p>
                </details>
            </li>

        </ul>

    </section>

    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="scroll-hint">
        <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 5.25 7.5 7.5 7.5-7.5m-15 6 7.5 7.5 7.5-7.5" />
    </svg>


    <Footer/>
</div>

<style>

    .heading-container {
        /*max-width: 1200px;*/
        width: 100%;
        display: flex;
        flex-direction: column;
        padding: 0rem;
    }

    .heading {
        font-family: Fugaz One, sans, sans-serif;
        font-size: 3rem;
        /*text-shadow: #DA1C5C 1px 1px 1px;*/
        margin: 0 auto;
        text-align: center;
        background-color: var(--affini-pink-main);
        width: 100%;
        --bg-color: var(--affini-pink-main);
    display: none;
        background: linear-gradient(-90deg, oklch(from var(--bg-color) l c -15deg) 0%, var(--bg-color) 100%);

    }

    h2 {
        --bg-color: var(--affini-pink-main);
        font-family: Fugaz One, sans, sans-serif;
        margin: 1rem 0 1rem 0;
        padding: 0.25rem .75rem;
        clip-path: polygon(0% 0%, 100% 0%, calc(100% - 4rem) 100%, 0% 100%);
        background: linear-gradient(-90deg, oklch(from var(--bg-color) l c -15deg) 0%, var(--bg-color) 100%);
    }


    .container {
        min-height: 100vh;
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: 1fr auto;
        /*background: url("/img/banner.gif") no-repeat center center / cover;*/

    }

    .container:has(.flash-message) {
        & > main {
            padding-top: 4rem;
        }
    }

    main {
        align-self: center;
        margin: auto;
        display: grid;
        grid-template-columns: 1fr;
        gap: 1rem;
        min-height: calc(100vh - 1px);
        background: url("/img/banner.gif") center center ;
        width: 100%;
        position: relative;
        border-bottom: 1px solid var(--affini-pink-main);

    }

    .scroll-hint {
        position: absolute;
        display: block;
        bottom: 1rem;
        left: 0;
        right: 0;
        width: 32px;
        margin: 0 auto;
        color: var(--affini-pink-main);
    }


    #faq {
        background: var(--content-bg);
        width: 100%;
        max-width: 1200px;
        padding: 3rem 1rem;
        display: flex;
        flex-direction: column;
        justify-content: center;
        margin: 0 auto;
        & > ul {
            list-style: none;
        }
    }

    #faq > ul {
        max-width: 40rem;
        width: 100%;
        margin: 0 auto;

    }

    details {
        margin-bottom: .5rem;

        summary {
            position: relative;
            cursor: pointer;
            border: 2px solid rgba(from var(--a-blue-low) r g b / 0.9);
            /*background: #6c757d;*/
            padding: 1rem 1rem 1rem 1rem;
            border-radius: 0.2rem;
            display: grid;
            grid-template-rows: 0fr;
            font-weight: 700;
            font-size: 1.15rem;
            + * {
                margin-top: 0;
            }

            &::after {
                content: "+";
                display: flex;
                flex-direction: column;
                justify-content: center;
                position: absolute;
                right: .5rem;
                top: 0;
                bottom: 0;
                margin: auto;
            }


        }

        &:open {
            summary {
                border-radius: 0.2rem 0.2rem 0 0;
                border-color: rgba(from var(--a-blue-low) r g b / 0.9);
                border-width: 2px 2px 0px 2px;

                /*background: #DA1C5C;*/
                &::after {
                    content: "-";
                    /*display: flex;*/
                    /*flex-direction: column;*/
                    /*justify-content: center;*/
                    /*position: absolute;*/
                    /*right: .5rem;*/
                    /*top: 0;*/
                    /*bottom: 0;*/
                    /*margin: auto;*/
                }
            }



            &::details-content {
                border-radius: 0 0 0.2rem 0.2rem;
                border-width: 0px 2px 2px 2px;
                border-color: rgba(from var(--a-blue-low) r g b / 0.9);
                border-style: solid;
                padding: 0.5rem 1rem 1.5rem 1rem;
                display: grid;
                grid-template-rows: 1fr;
                /*background: #6a4298;*/
            }
        }
    }

    @media screen and (min-width: 768px) {

        h2 {
            margin: 1rem -8rem 4rem -4rem;
            padding: 0.25rem 4.75rem 0.25rem 4.5rem;
        }

        .heading-container {
            padding: 4rem 4rem 2rem 4rem;
        }
        #faq > ol {
            padding-block: 2rem 4rem;
        }
    }
</style>