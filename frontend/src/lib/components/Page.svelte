<script lang="ts">
    import type {Snippet} from "svelte";
    import PageHeader from "./PageHeader.svelte";
    import page from "$lib/transitions/page";
    import {createNavState, type NavigationState} from "$lib/types";
    import app from "$lib/stores/state.svelte";


    interface Props extends Partial<NavigationState> {
        children: Snippet;
        icon?: Snippet;
        title?: string;
    }

    // TODO: remove when the PR lands in eslint-plugin-svelte
    // the below is not actually an issue with our code
    // eslint-disable-next-line svelte/valid-compile
    const {children, icon, title = "BetterDiscord", ...footer}: Props = $props();


    // Allow page users to easily setup the navigation state
    $effect(() => {
        const nav: Omit<NavigationState, "direction"> = createNavState(footer);
        for (const n in nav) {
            const key = n as keyof Omit<NavigationState, "direction">;
            // @ts-expect-error it doesn't make sense
            app.navigation[key] = nav[key];
        }
    });
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