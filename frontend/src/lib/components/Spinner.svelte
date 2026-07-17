<script lang="ts">
    // Purely presentational, indeterminate spinner — no props or state.
</script>

<svg
    class="spinner"
    width="32"
    height="32"
    viewBox="0 0 16 16"
    role="progressbar"
    aria-label="Loading"
    aria-busy="true"
>
    <circle class="spinner-ring" cx="50%" cy="50%" r="7" />
    <circle class="spinner-fill" cx="50%" cy="50%" r="7" stroke-dasharray="3" />
</svg>

<style>
    @keyframes spinner-indeterminate {
        0% {
            stroke-dasharray: 0.01px 43.97px;
            transform: rotate(0deg);
        }
        50% {
            transform: rotate(450deg);
            stroke-dasharray: 21.99px 21.99px;
        }
        100% {
            stroke-dasharray: 0.01px 43.97px;
            transform: rotate(1080deg);
        }
    }

    .spinner {
        width: 32px;
        height: 32px;
    }

    .spinner circle {
        fill: none;
        transform: rotate(-90deg);
        transition: all 0.2s ease-in-out;
        stroke-width: 2;
        stroke-linecap: round;
        transform-origin: 50% 50%;
    }

    .spinner-ring {
        stroke: var(--bg3);
    }

    .spinner .spinner-fill {
        stroke: var(--accent);
        stroke-dasharray: 43.75;
        animation: spinner-indeterminate 2s linear infinite;
    }
</style>
