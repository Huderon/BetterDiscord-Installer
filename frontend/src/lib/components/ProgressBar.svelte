<script lang="ts">
    import type {ClassValue} from "svelte/elements";

    interface Props {
        value?: number;
        max?: number;
        indeterminate?: boolean;
        class: ClassValue;
    }

    // eslint-disable-next-line prefer-const
    let {value = 0, max = 100, indeterminate = false, class: className}: Props = $props();
</script>

<div class="progress-bar {className}">
    {#if indeterminate}
        <div class="progress-fill increase"></div>
        <div class="progress-fill decrease"></div>
    {:else}
        <div class="progress-fill" style:width="{(value / max) * 100}%"></div>
    {/if}
</div>

<style>
    .progress-bar {
        width: 100%;
        height: 6px;
        border-radius: 6px;
        overflow: hidden;
        background-color: var(--bg3);
        position: relative;
    }

    .progress-fill {
        height: 100%;
        border-radius: 6px;
        overflow: hidden;
        transition: 500ms ease width, 150ms ease background-color;
        transform-origin: left;
        background-color: var(--accent);
        position: absolute;
    }

    :global(.progress-bar.error .progress-fill) {
        background-color: var(--danger);
    }

    :global(.progress-bar.success .progress-fill) {
        background-color: #3ac172;
    }

    /* indeterminate styles */
    @keyframes decrease {
        from {
            margin-left: -80%;
            width: 80%;
        }

        to {
            width: 10%;
            margin-left: 110%;
        }
    }

    @keyframes increase {
        from {
            margin-left: -5%;
            width: 5%;
        }

        to {
            width: 100%;
            margin-left: 130%;
        }
    }

    .progress-bar .progress-fill.increase {
        animation: increase 2s ease infinite;
        width: 0%;
    }

    .progress-bar .progress-fill.decrease {
        animation: decrease 2s 0.5s ease infinite;
        width: 0%;
    }
</style>