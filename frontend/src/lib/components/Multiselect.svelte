<script lang="ts">
    import Button from "./Button.svelte";
    import {handleKeyboardToggle} from "$lib/utils/handlers";
    import type {Snippet} from "svelte";
    import {tooltip} from "./tooltip";

    interface Props {
        description: string;
        disabled?: boolean
        checked?: boolean;
        onclick?: (event: MouseEvent) => void;
        onchange?: (event: Event) => void;
        icon?: Snippet;
        children?: Snippet;
        showTooltip?: boolean;
        showDescription?: boolean;
    }

    // eslint-disable-next-line prefer-const
    let {children, icon, description, disabled = false, checked = $bindable(), onclick, onchange, showTooltip = false, showDescription = true}: Props = $props();

    let checkbox: HTMLInputElement;

</script>

<label class="check-container" use:tooltip={{text: showTooltip ? description : ""}}>
    <input bind:this={checkbox} bind:checked type="checkbox" hidden {disabled} {onchange} />
    <div onkeypress={(e: KeyboardEvent) => handleKeyboardToggle(e, checkbox)} tabindex="0" class="check-item" class:disabled role="listbox">
        <div class="icon">
            {@render icon?.()}
        </div>
        <div class="content">
            <h5>
                {@render children?.()}
            </h5>
            {#if showDescription}
                <span title={description}>{description}</span>
            {/if}
        </div>
        <div class="controls" role="presentation" onkeypress={e => e.stopPropagation()}>
            {#if onclick}
                <Button style="secondary" {onclick}>Browse</Button>
            {:else}
                <div class="checkmark-wrap" class:checked>
                    <svg class="checkmark" viewBox="0 0 24 24">
                        <path d="M0.73, 11.91 8.1,19.28 22.79,4.59" fill="none" />
                    </svg>
                </div>
            {/if}
        </div>
    </div>
</label>

<style>
    .check-item {
        display: flex;
        align-items: center;
        border-radius: 3px;
        background-color: var(--bg3);
        padding: 12px;
        user-select: none;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        cursor: pointer;
        transition: 100ms ease;
        flex-wrap: nowrap;
        position: relative;
        overflow: hidden;
    }

    .check-container {
        margin-bottom: 12px;
    }

    .check-container:last-child {
        margin: 0;
    }

    .check-item.disabled {
        background-color: var(--bg2-alt);
        cursor: not-allowed;
    }

    .check-container input:checked + .check-item {
        background-color: var(--accent);
    }

    .check-item.disabled .content,
    .check-item.disabled .icon {
        opacity: 0.5;
        pointer-events: none;
    }

    .controls,
    .icon {
        flex: 0 0 auto;
    }

    :global(.icon img) {
        width: 32px;
        height: 32px;
    }

    .content {
        display: flex;
        flex-direction: column;
        margin: 0 10px;
        overflow: hidden;
        flex: 1 1 auto;
    }

    .content span,
    .content h5 {
        transition: 100ms ease;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        line-height: normal;
    }

    .content span {
        color: var(--text-muted);
        font-size: 12px;
        font-weight: 400;
    }

    .content h5 {
        color: var(--text-normal);
        font-weight: 600;
        font-size: 13px;
        margin: 0;
    }

    .check-item:not(.disabled):hover .content h5 {
        color: var(--text-light);
    }

    .check-item:not(.disabled):hover .content span {
        color: var(--text-normal);
    }

    .check-container input:checked + .check-item .content h5,
    .check-container input:checked + .check-item .content span {
        color: #fff;
    }

    :global(.check-container input:checked + .check-item .button) {
        background-color: #fff;
        border-color: transparent !important;
        color: var(--accent);
    }

    :global(.check-container input:checked + .check-item .button:active) {
        opacity: 0.8;
    }

    .checkmark {
        width: 24px;
        height: 24px;
        color: #fff;
    }

    .checkmark path {
        transform: scale(0.8);
        transform-origin: center;
        stroke: currentColor;
        stroke-width: 2.45;
        stroke-dasharray: 32;
        stroke-dashoffset: 32;
    }

    .checked .checkmark path {
        transition: 250ms cubic-bezier(0.55, 0, 0, 1) stroke-dashoffset;
        stroke-dashoffset: 0;
    }

    .controls {
        margin-right: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .checkmark-wrap {
        width: 24px;
        height: 24px;
        border-radius: 4px;
        /* background-color: rgba(255, 255, 255, 0.2); */
        display: flex;
        justify-content: center;
        align-items: center;
        border: 2px solid rgba(255, 255, 255, 0.2);
        transition: 250ms cubic-bezier(0.55, 0, 0, 1) border-color;
    }

    .checkmark-wrap.checked {
        border-color: transparent;
    }
</style>