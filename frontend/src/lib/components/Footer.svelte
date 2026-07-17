<script lang="ts">
    import Button from "./Button.svelte";
    import ButtonGroup from "./ButtonGroup.svelte";
    import SocialLinks from "./SocialLinks.svelte";
    import {goto} from "$app/navigation";
    import {page} from "$app/state";
    import {base} from "$app/paths";
    import app from "$lib/stores/state.svelte";
    import {NavDirection, type NavigationState} from "$lib/types";

    const {
        next, previous,
        nextLabel = "Next",
        previousLabel = "Back",
        canGoNext = true,
        canGoPrevious = true,
        nextAction, previousAction
    }: NavigationState = $derived(app.navigation);


    const nextDisabled: boolean = $derived.by(() => {
        if (!canGoNext) return true;
        if (!nextAction && !next) return true;
        return false;
    });

    const previousDisabled: boolean = $derived.by(() => {
        if (!canGoPrevious) return true;
        if (!previousAction && !previous) return true;
        return false;
    });


    async function goToNext() {
        app.navigation.direction = NavDirection.FORWARDS;
        if (typeof nextAction === "function") return nextAction();
        // eslint-disable-next-line svelte/no-navigation-without-resolve
        if (next) await goto(`${base}${next}`, page.state);
    }

    async function goBack() {
        app.navigation.direction = NavDirection.BACKWARDS;
        if (typeof previousAction === "function") return previousAction();
        // eslint-disable-next-line svelte/no-navigation-without-resolve
        if (previous) await goto(`${base}${previous}`, page.state);
    }

    async function navigatePage(event: KeyboardEvent) {
        if ((event.key === "ArrowRight" && event.ctrlKey) && canGoNext) {
            await goToNext();
        }
        else if ((event.key === "ArrowLeft" && event.ctrlKey) && canGoPrevious) {
            await goBack();
        }
    }

    async function tryGoNext(event: KeyboardEvent) {
        if (!canGoNext) return;
        if (event.key !== "Enter" && event.key !== " ") return;
        event.preventDefault();
         await goToNext();
    }

    async function tryGoPrevious(event: KeyboardEvent) {
        if (!canGoPrevious) return;
        if (event.key !== "Enter" && event.key !== " ") return;
        event.preventDefault();
         await goBack();
    }

</script>

<svelte:window onkeydown={navigatePage} />

<footer class="install-footer">
    <SocialLinks />
    <ButtonGroup>
        <Button style="secondary" disabled={previousDisabled} onclick={goBack} onkeypress={tryGoPrevious}>
            {previousLabel}
        </Button>
        <Button style="primary" disabled={nextDisabled} onclick={goToNext} onkeypress={tryGoNext}>
            {nextLabel}
        </Button>
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