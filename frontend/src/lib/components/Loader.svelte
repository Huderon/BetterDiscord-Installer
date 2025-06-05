<script lang="ts">
    import Spinner from "$lib/components/Spinner.svelte";
</script>

<div class="loader">
    <Spinner />
</div>

<style>
    .loader {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>