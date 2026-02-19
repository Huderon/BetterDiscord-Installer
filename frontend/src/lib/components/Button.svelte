<script lang="ts">
    import type {Snippet} from "svelte";

    interface Props {
        style?: "primary" | "secondary";
        tabindex?: number|string;
        disabled?: boolean;
        onkeypress?: (event: KeyboardEvent) => void;
        onclick?: (event: MouseEvent) => void;
        children: Snippet;
    }

    const {style = "secondary", onkeypress, onclick, children, tabindex, disabled = false}: Props = $props();

    const styles = ["primary", "secondary"];
</script>

<button
    class="button {styles.includes(style) ? `style-${style}` : "style-secondary"}"
    type="button"
    {onkeypress}
    {onclick}
    tabindex={tabindex as number}
    {disabled}
>
    <span>
        {@render children()}
    </span>
</button>

<style>
    .button {
        width: auto;
        border: none;
        display: flex;
        align-items: center;
        justify-content: center;
        text-decoration: none;
        transition: 150ms ease;
        cursor: pointer;
        padding: 0 12px;
        height: 28px;
        white-space: nowrap;
        font-size: 12px;
        font-weight: 400;
        border-radius: 2px;
    }

    .button[disabled] {
        opacity: .5;
        pointer-events: none;
    }

    .button span {
        max-width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .button.style-primary {
        border: 1px solid transparent;
        background-color: var(--accent);
        color: #ffffff;
    }

    .button.style-primary:hover {
        background-color: var(--accent-hover);
    }

    .button.style-secondary {
        background-color: transparent;
        border: 1px solid rgba(255, 255, 255, 0.05);
        color: var(--text-normal);
    }

    .button.style-secondary:hover {
        border-color: rgba(255, 255, 255, 0.1);
    }
</style>