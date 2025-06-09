import {channels, createNavState, NavDirection, type AppState} from "$lib/types";
import {GetDiscordPath} from "@backend/Paths";


const app = $state<AppState>({
    eulaAgreed: false,
    action: "install",
    corePaths: {stable: "", ptb: "", canary: ""},
    channels: {stable: false, ptb: false, canary: false},
    navigation: createNavState({
        direction: NavDirection.FORWARDS,
    })
});

try {
    for (const channel of channels) {
        GetDiscordPath(channel).then(path => app.corePaths[channel] = path);
    }
}
catch {
    // Not in wails environment
}

export default app;