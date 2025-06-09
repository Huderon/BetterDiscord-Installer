<script lang="ts">
    import type {Snippet} from "svelte";

    const {children}: {children: Snippet} = $props();
</script>


<div class="button-group">
    {@render children()}
</div>


<style>
    .button-group {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        margin-top: -5px;
        margin-right: -5px;
    }

    :global(.button-group > *) {
        margin-top: 5px;
        margin-right: 5px;
    }
</style>