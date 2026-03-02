<script>
    import SocialIcons from '@rodneylab/svelte-social-icons';
</script>

<footer class="footer">

    <div class="container">
        <section>
            <h3>About</h3>
            <p>by <a href="https://archiveofourown.org/users/immaterial_vivi">immaterial vivi</a>
            <p>and the <a href="https://humandomestication.guide/authors">hd.g contributors</a>
        </section>
        <section>
            <h3>Links</h3>
            <a href="https://www.tumblr.com/immaterial-vivi"><SocialIcons alt="Tumblr" network="tumblr" fgColor="#15a5ff" bgColor="#11111100" /></a>
            <a href="https://www.github.com/immaterial-vivi"><SocialIcons alt="GitHub" network="github" fgColor="#15a5ff" bgColor="#11111100" /></a>

        </section>
        <section>
            <h3>&nbsp;</h3>
            <span>
            <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24" stroke-width="1.5"
                 stroke="currentColor" class="heart">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"/>
            </svg> Ploogle is <a
                    href="https://www.github.com/immaterial-vivi/ploogleV2"> free and open source software </a>
            </span>
        </section>
    </div>

</footer>

<style>

    .footer {
        background: url("/img/footer.gif") no-repeat 0 0  / cover;
        width: 100%;
        border-top: 1px var(--affini-pink-main) solid;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
        position: relative;
    }

    .footer::before {
        content: "";
        display: block;
        position: absolute;
        top: -50px;
        right: 12px;
        width: 64px;
        height: 50px;
        background: url("/img/bnuy.svg") no-repeat 0 0 / 64px 64px;
    }


    .container {
        position: relative;
        max-width: 1200px;
        display: grid;
        grid-template-columns: 1fr;
        align-items: top;

        padding: 1rem;
        margin: 0 auto;
        & > section {
            padding: .5rem 0;
        }
    }

    @media screen and (min-width: 768px) {
        .container {
            grid-template-columns: 1fr 1fr 1fr;
        }
    }

    h3 {
        font-family: Fugaz One, sans, sans-serif;
        font-weight: bold;
        color: var(--affini-pink-main);
    }

    .heart {
        height: 1.2rem;
        display: inline-block;
        position: relative;
        vertical-align: middle;
        color: var(--affini-pink-main);
    }

</style>