<script lang="ts">
    import {page} from '$app/state';
    import Footer from '$lib/components/footer.svelte'
    import {resolve} from "$app/paths";
</script>

<main>
    <div>
        <h1>{page.error.message}</h1>
        <p>
            <a href={resolve("/")}>Back to home page</a>
        </p>
    </div>

    <Footer/>
</main>


<style>
    div {
        padding: 2rem;
    }
    p {
        padding: 2rem 0;
    }
    h1 {
        font-size: 3rem;
    }

    main {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        background: var(--content-bg);

    }
</style>