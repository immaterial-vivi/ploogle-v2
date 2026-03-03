<svg class="bnuy" width="64" height="64" viewBox="0 0 128 128" fill="#DA1C5C" xmlns="http://www.w3.org/2000/svg">
    <title>bunnies go hop hop hop c:</title>
    <path class="bnuy-body" d="M115.496 87.1653C114.972 104.225 89.6369 107.687 82.7073 108.5H35.6263C27.6775 108.5 7.1771 105.664 7.49999 87.1653C7.49999 75.5836 21.5 56.6115 30.5 56C39.5 55.3885 50.9999 55.3885 65.5336 55.3885H65.5338C65.5338 55.3885 116.5 54.5 115.496 87.1653Z" fill="#DA1C5C" />
    <path class="bnuy-ear-b" d="M54.6995 12.0062C58.2844 12.1656 59.0942 26.0675 56.7059 38.4217C56.1635 41.2272 51.3899 65.6342 50.5 69C50.5 70.4919 30 69 27.5 69C26.9048 66.2622 27.3582 61.5631 28.5795 56C28.5795 36.5 44.2616 11.542 54.6995 12.0062Z" fill="#DA1C5C" />
    <path class="bnuy-ear-f" d="M96.5 19.8941C98.4857 21.8071 98.9857 30.3071 88.9438 41.8383C80.9044 49.8141 66.5934 67.6832 64.5 69.1059C63.8388 69.5553 38 68.5 38 68.5C38 65.5 48.4857 43.8071 60.2253 33.8071C69.9503 25.4951 89.3757 16.0193 96.5 19.8941Z" fill="#DA1C5C" />
</svg>

<style>
    .bnuy {
        --perk: cubic-bezier(0.516, 0, 0.074, 2);
    }

    .bnuy-body {
        transform-origin: bottom center;
        transform: scaleX(1) scaleY(1);
        transition: all 150ms var(--perk);
    }

    .bnuy-ear-f {
        transform-origin: center center;
        transform: rotateZ(13deg);
        transition: all 150ms 15ms cubic-bezier(0.516, 0, 0.252, 1.423);
    }

    .bnuy-ear-b {
        transform-origin: center;
        transform: rotateZ(7deg);
        transition: all 150ms 30ms var(--perk);
    }

    .bnuy:hover {
        & > .bnuy-ear-b {
            transform: rotateZ(-15deg) translateY(-3px);
        }

        & > .bnuy-ear-f {
            transform-origin: center left;
            transform: rotateZ(-5deg);
        }

        & > .bnuy-body {
            transform: scaleX(0.95) scaleY(1.05);
        }
    }
</style>