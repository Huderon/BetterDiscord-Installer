<script lang="ts">
    import Button from "./Button.svelte";
    import ButtonGroup from "./ButtonGroup.svelte";
    import SocialLinks from "./SocialLinks.svelte";
    import {canGoForward, canGoBack, nextPage, state} from "../stores/navigation";
    import quit from "../actions/quit";
    import {goto, onNavigate} from "$app/navigation";
    import {page} from "$app/state";
    import {base} from "$app/paths";


    let nextButtonContent = "Next";

    async function goToNext() {
        state.direction = 1;
        if ($nextPage) goto(`${base}${$nextPage}`, page.state);
        else await quit();
    }

    function goBack() {
        state.direction = -1;
        window.history.back();
    }

    onNavigate(() => {
        if (window.location.pathname.startsWith("/actions/setup/")) {
            const action = window.location.pathname.slice(15);
            const actionText = action[0].toUpperCase() + action.slice(1);
            nextButtonContent = actionText;
        }
        else {
            nextButtonContent = "Next";
        }
    });

    function navigatePage(event: KeyboardEvent) {
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
        <Button style="secondary" disabled={!$canGoBack} onclick={goBack}>Back</Button>
        <Button style="primary" disabled={!$canGoForward} onclick={goToNext}>{#if $nextPage}{nextButtonContent}{:else}Close{/if}</Button>
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