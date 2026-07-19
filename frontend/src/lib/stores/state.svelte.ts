import {channels, createNavState, NavDirection, type AppState} from "$lib/types";
import {GetDiscordPath} from "@api";


const app = $state<AppState>({
    eulaAgreed: false,
    action: "install",
    devUnlocked: false,
    resourcePaths: {stable: "", ptb: "", canary: ""},
    channels: {stable: false, ptb: false, canary: false},
    options: {
        install: {restartDiscord: true, useDevBuild: false},
        repair: {
            disablePlugins: false,
            disableThemes: false,
            clearCustomCSS: false,
            clearWebpackCache: false,
            clearAddonStoreCache: false,
            resetSettings: false
        },
        uninstall: {fullUninstall: false, restartDiscord: true}
    },
    navigation: createNavState({
        direction: NavDirection.FORWARDS,
    })
});

try {
    for (const channel of channels) {
        // eslint-disable-next-line new-cap
        void GetDiscordPath(channel)
            .then(path => app.resourcePaths[channel] = path)
            .catch(() => {/* leave resourcePaths[channel] empty (e.g. not in Wails) */});
    }
}
catch {
    // GetDiscordPath binding is missing entirely (not in a Wails environment)
}

export default app;
