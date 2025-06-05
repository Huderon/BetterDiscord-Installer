<script lang="ts">
    import type {Snippet} from "svelte";


    interface Props {
        type: "header" | "paragraph" | "subtext";
        hasMargin?: boolean;
        children: Snippet
    }

    const {hasMargin = false, type = "paragraph", children}: Props = $props();
</script>

{#if type === "header"}
    <h3 class="header" class:has-margin={hasMargin}>
        {@render children()}
    </h3>
    {:else if type === "paragraph"}
    <p class="paragraph" class:has-margin={hasMargin}>
        {@render children()}
    </p>
    {:else if type === "subtext"}
    <span class="subtext" class:has-margin={hasMargin}>
        {@render children()}
    </span>
{/if}

<style>
    .header,
    .paragraph,
    .subtext {
        margin: 0;
        line-height: normal;
    }

    .has-margin {
        margin-bottom: 10px;
    }

    .header {
        color: var(--text-light);
        font-weight: 600;
        font-size: 16px;
    }

    .paragraph {
        color: var(--text-normal);
        font-weight: 400;
        font-size: 14px;
    }

    .subtext {
        color: var(--text-muted);
        font-weight: 400;
        font-size: 12px;
    }
</style>