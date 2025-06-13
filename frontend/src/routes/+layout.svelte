<script lang="ts">
    import "focus-visible";
    import "$lib/styles/theme.css";
    import "$lib/styles/global.css";
    import Titlebar from "$lib/components/Titlebar.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import {onMount} from "svelte";


    const {children} = $props();

    const isMac = $derived(window.navigator.platform.startsWith("Mac"));

    onMount(async () => {
        // Disable user zooming
        window.addEventListener("keydown", (e) => {
            if ((e.code === "Minus" || e.code === "Equal") && (e.ctrlKey || e.metaKey)) {
                e.preventDefault();
            }
        });
    });
</script>


<div class="main-window" class:platform-darwin={isMac}>
    <Titlebar macButtons={isMac} />
    <main class="installer-body">
        <div class="sections">
             <!-- {#key window.location.pathname} -->
             {@render children()}
             <!-- {/key} -->
        </div>
        <Footer />
    </main>
</div>


<style>
    .main-window {
        display: flex;
        flex-direction: column;
        /* border-radius: 3px; */
        overflow: hidden;
        contain: strict;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.25);
        /* margin: 11.5px 7.5px; */
        width: 100%;
        height: 100%;
        word-break: break-word;
    }

    .main-window.platform-darwin {
        border-radius: 0;
        box-shadow: none;
        width: 100%;
        height: 100%;
        margin: 0;
    }

    .installer-body {
        overflow: hidden;
        position: relative;
        display: flex;
        flex-direction: column;
        z-index: 1;
        padding: 20px;
        background: radial-gradient(var(--bg2) 50%, var(--bg2-alt));
        flex: 1;
    }

    .installer-body::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-image: url("$lib/assets/images/background.png");
        background-size: 60px;
        background-repeat: repeat;
        background-position: center;
        z-index: -1;
        opacity: 0.35;
        pointer-events: none;
        mask: radial-gradient(transparent, #000);
        -webkit-mask: radial-gradient(transparent, #000);
    }

    .sections {
        flex: 1 1 auto;
        overflow: visible;
        position: relative;
    }
</style>