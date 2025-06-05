import {paths, os} from "$lib/stores/installation";
import {GetPlatform as platform} from "@backend/Backend";
import {GetDiscordPath as getDiscordPath} from "@backend/Paths";


// window.refresh = () => window.location.href = `http://${window.location.host}/`;


// Disable user zooming
window.addEventListener("keydown", (e) => {
    if ((e.code === "Minus" || e.code === "Equal") && (e.ctrlKey || e.metaKey)) {
        e.preventDefault();
    }
});

try {
    platform().then(osName => os.set(osName));

    const channels = ["stable", "ptb", "canary"];

    for (const channel of channels) {
        getDiscordPath(channel).then(path => {
            paths.update(current => ({...current, [channel]: path}));
        });
    }
}
catch {
    // Not in a wails environment
}

export const ssr = false;
// export const prerender = true;