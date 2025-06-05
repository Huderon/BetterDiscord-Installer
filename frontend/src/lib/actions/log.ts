import logs from "../stores/logs";
import {action, progress, status} from "../stores/installation";

export function log(entry: string) {
    logs.update(a => {
        a.push(entry);
        return a;
    });
}

export function lognewline(entry: string) {
    logs.update(a => {
        a.push("");
        a.push(entry);
        return a;
    });
}

export function succeed() {
    const name = action.value;
    log("");
    log(`${name.charAt(0).toUpperCase() + name.slice(1)} completed!`);
    progress.set(100);
    status.set("success");
}

export async function reset() {
    logs.set([]);
    progress.set(0);
    status.set("");
    await new Promise(r => setTimeout(r, 500));
}

const discordURL = "https://betterdiscord.app/invite";

export function fail() {
    log("");
    log(`The ${action.value} seems to have failed. If this problem is recurring, join our discord community for support. ${discordURL}`);
    status.set("error");
}