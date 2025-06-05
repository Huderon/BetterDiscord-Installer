import {writable} from "svelte/store";
import readwritable from "./types/readwritable";
import type {DiscordChannel, InstallerAction, OSName} from "$lib/types";

export const status = readwritable("");
export const hasAgreed = writable(false);
export const platforms = writable<Record<DiscordChannel, boolean>>({stable: false, canary: false, ptb: false});
export const paths = writable<Record<DiscordChannel, string>>({stable: "", canary: "", ptb: ""});
export const os = readwritable<OSName>("windows");

export const progress = readwritable(0);
export const action = readwritable<InstallerAction>("install");
