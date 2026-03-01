<script>
    import SearchForm from '$lib/components/search-form.svelte';
    import Flort from '$lib/components/flort.svelte';
    import Footer from '$lib/components/footer.svelte';
    import FlashMessage from '$lib/components/flash-message.svelte';

</script>

<div class="container">
<!--    <FlashMessage>-->
<!--            Welcome to the new Ploogle! <br> Read all about it <a href="#" target="_blank" rel="nofollow">here!</a>-->
<!--    </FlashMessage>-->
    <main>
        <!--        <h1 class="heading">Ploogle</h1>-->
        <p class="sr-only">
            Ploogle - Search for <a href="https://humandomestication.guide/">HDG Stories</a>
        </p>
        <Flort/>
        <SearchForm/>
    </main>
    <Footer/>
</div>

<style>

    .heading {
        font-family: Fugaz One, sans, sans-serif;
        font-size: 3rem;
        text-shadow: #DA1C5C 1px 1px 1px;
        margin: 0 auto;

    }

    .container {
        min-height: 100vh;
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: 1fr auto;
        /*background: url("/img/banner.gif") no-repeat center center / cover;*/

        background: var(--space-blue-low);
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
    }


</style>