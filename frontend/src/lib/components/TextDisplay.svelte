<script lang="ts">
    import {tick} from "svelte";
    import Button from "./Button.svelte";
    import LoadingPage from "./Loader.svelte";


    interface Props {
        value: string;
        autoscroll?: boolean;
    }

    const {value, autoscroll = false}: Props = $props();

    let scroller: HTMLDivElement | undefined = $state();
    let element: HTMLElement | undefined = $state();
    let copyInputContainer: HTMLDivElement | undefined = $state();
    let copyButtonActive = $state(false);
    let copyButtonVisible = $state(false);
    let shouldAutoscroll: boolean = autoscroll;

    // Copy button
    function copyDisplayContents() {
        if (!element) return;
        copyButtonActive = true;
        const range = document.createRange();
        range.selectNode(element);
        window.getSelection()?.addRange(range);
        document.execCommand("Copy");
        document.getSelection()?.removeAllRanges();
        setTimeout(() => {
            copyButtonActive = false;
        }, 500);
    }

    function handleKeyboardCopyToggle(event: KeyboardEvent) {
        if (event.key === "Enter" || event.key === " ") copyDisplayContents();
    }

    // Autoscroll
    $effect.pre(() => {
        if (!value || !autoscroll) return;
        shouldAutoscroll = !!(scroller && (scroller.offsetHeight + scroller.scrollTop) > (scroller.scrollHeight - 20));
        if (!shouldAutoscroll) return;
        tick().then(() => {
            scroller?.scrollTo(0, scroller?.scrollHeight);
        });
    });
</script>

{#if value}
    <article
        bind:this={element}
        onmousemove={() => copyButtonVisible = true}
        onmouseleave={() => copyButtonVisible = false}
        class="text-display{value ? "" : " loading"}"
    >
        <div bind:this={scroller} onscroll={() => copyButtonVisible = false} class="display-inner" tabindex="0" role="textbox">
            {value}
        </div>
        <div bind:this={copyInputContainer} class="copy-input" class:visible={copyButtonVisible}>
            {#if copyButtonActive}
                <Button tabindex="0" style="primary" onkeypress={handleKeyboardCopyToggle} onclick={copyDisplayContents}>Copied!</Button>
            {:else}
                <Button tabindex="0" style="secondary" onkeypress={handleKeyboardCopyToggle} onclick={copyDisplayContents}>Copy</Button>
            {/if}
        </div>
    </article>
{:else}
    <LoadingPage />
{/if}

<style>
    .text-display {
        position: relative;
        display: flex;
        flex: 1;
        min-height: 0;
        margin-bottom: 10px;
        background: var(--bg3);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        border-radius: 2px;
    }

    .text-display .display-inner {
        color: var(--text-normal);
        font-size: 13px;
        line-height: 1.5;
        word-wrap: normal;
        white-space: pre-wrap;
        user-select: text;
        height: 100%;
        width: 100%;
        overflow: auto;
        padding: 12px;
        border-radius: inherit;
    }

    .text-display.loading {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    /* Copy Button */

    .copy-input {
        position: absolute;
        bottom: 8px;
        right: 8px;
    }

    :global(.copy-input .button) {
        border: none !important;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }

    :global(.copy-input .button.type-secondary) {
        background-color: var(--bg4) !important;
    }

    :global(.copy-input .button:hover) {
        color: var(--text-light) !important;
    }

    :global(.copy-input .button:not(.type-primary)) {
        opacity: 0;
    }

    :global(.copy-input.visible .button),
    :global(.display-inner.focus-visible + .copy-input .button),
    :global(.copy-input .button.focus-visible) {
        opacity: 1;
    }
</style>