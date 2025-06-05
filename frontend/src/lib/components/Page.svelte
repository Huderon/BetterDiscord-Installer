<script lang="ts">
    import type {Snippet} from "svelte";
    import PageHeader from "./PageHeader.svelte";
    import page from "$lib/transitions/page";


    interface Props {
        children: Snippet;
        icon?: Snippet;
        title?: string;
    }

    const {children, icon, title = "BetterDiscord"}: Props = $props();
</script>


<div style:display="contents">
    {#key title}
        <section class="page" in:page|global out:page|global={{out: true}}>
            <PageHeader {icon}>{title}</PageHeader>
            {@render children()}
        </section>
    {/key}
</div>


<style>
    .page {
        flex: 1 1 auto;
        overflow: visible;
        display: flex;
        flex-direction: column;
        position: absolute;
        width: 100%;
        height: 100%;
    }
</style>