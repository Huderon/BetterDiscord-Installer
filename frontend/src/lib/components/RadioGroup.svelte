<script lang="ts">
    import type {Snippet} from "svelte";
    import {handleArrowKeys} from "$lib/utils/handlers";

    const {index = 0, children}: {index?: number, children: Snippet} = $props();

    let container: HTMLDivElement;
</script>

<div
    onkeydown={(e: KeyboardEvent) => handleArrowKeys(e, container)}
    bind:this={container}
    tabindex="0"
    data-selected-index={index}
    style:--index={index}
    class="radio-group"
    role="radiogroup"
>
    {@render children?.()}
    <div class="selection-indicator"></div>
</div>

<style>
    .radio-group {
        position: relative;
        display: flex;
        flex-direction: column;
        border-radius: 3px;
        transition: 150ms ease;
    }

    .selection-indicator {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 40px;
        background-color: var(--accent);
        border-radius: 3px;
        transition: 250ms ease;
        z-index: 0;
        transform: translateY(calc((100% + 12px) * var(--index)));
    }
</style>