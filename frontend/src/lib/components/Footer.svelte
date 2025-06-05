<script lang="ts">
    // const electron = require("electron");

    import Button from "./Button.svelte";
    import ButtonGroup from "./ButtonGroup.svelte";
    import SocialLinks from "./SocialLinks.svelte";
    import {canGoForward, canGoBack, nextPage, state} from "../stores/navigation";
    // import {push, pop, location} from "svelte-spa-router";
    import quit from "../actions/quit";
    import {goto, pushState} from "$app/navigation";
    import {page} from "$app/state";

    let nextButtonContent = "Next";

    async function goToNext() {
        state.direction = 1;
        if ($nextPage) goto($nextPage, page.state);
        else await quit();
    }

    function goBack() {
        state.direction = -1;
        window.history.back();
    }

    $: if (window.location.pathname.startsWith("/setup/")) {
        const action = window.location.pathname.slice(7);
        const actionText = action[0].toUpperCase() + action.slice(1);
        nextButtonContent = actionText;
    }
    else {
        nextButtonContent = "Next";
    }

    function navigatePage() {
        if ((event.key === "ArrowRight" && event.ctrlKey) && $canGoForward) {
            goToNext();
        }
        else if ((event.key === "ArrowLeft" && event.ctrlKey) && $canGoBack) {
            goBack();
        }
    }

</script>

<svelte:window on:keydown={navigatePage} />

<footer class="install-footer">
    <SocialLinks />
    <ButtonGroup>
        <Button type="secondary" disabled={!$canGoBack} on:click={goBack}>Back</Button>
        <Button type="primary" disabled={!$canGoForward} on:click={goToNext}>{#if $nextPage}{nextButtonContent}{:else}Close{/if}</Button>
    </ButtonGroup>
</footer>

<style>
    .install-footer {
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: flex-end;
        justify-content: space-between;
        flex: 0 0 auto;
        margin-top: 10px;
    }
</style>