<script>
    let {children} = $props();
    let dismissed = $state(false);
</script>
{#if !dismissed}
<div class="flash-message">
        <span>
            <p>
                {@render children?.()}
            </p>
            <button onclick={(_)=> dismissed = true} class="flash-dismiss"> <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                                stroke-width="1.5" stroke="currentColor" class="size-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12"/>
            </svg>
            </button>
        </span>
</div>
{/if}

<style>
    .flash-message {
        position: absolute;
        width: 100%;
        padding: 2rem;
        text-align: center;

        & > span {
            display: flex;
            justify-content: center;
            align-items: center;
            max-width: 22rem;
            background: #6c757d;
            margin: 0 auto;
            padding: 0.5rem 1rem;
            background: rgba(from var(--a-light-blue) r g b / 0.1);
            border: 1px solid var(--a-light-blue);
            color: var(--a-light-blue);
            border-radius: 0.5rem;
            gap: 1rem;

            & > a {
                display: contents;
                text-decoration: none;
                color: white;
            }
        }

        & > * > svg {
            height: 2rem;
        }
    }

    .flash-dismiss {
        background: none;
        border: none;
        color: var(--a-light-blue);
        cursor: pointer;

        & > svg {
            height: 2rem;
        }
    }
</style>