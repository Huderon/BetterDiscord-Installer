import type {types} from "@models";

export type DiscordChannel = "stable" | "ptb" | "canary";
export type InstallerAction = "install" | "repair" | "uninstall";
export type OSName = "windows" | "darwin" | "linux";

export const channels = ["stable", "canary", "ptb"] as const;
export const labels = {stable: "Discord", ptb: "Discord PTB", canary: "Discord Canary"} as const;
export const actions: InstallerAction[] = ["install", "repair", "uninstall"] as const;

export enum NavDirection {FORWARDS = 1, BACKWARDS = -1}

export interface NavigationState {
    direction: NavDirection;
    next?: string;
    previous?: string;
    nextLabel?: string;
    previousLabel?: string;
    canGoNext?: boolean;
    canGoPrevious?: boolean;
    nextAction?: () => void;
    previousAction?: () => void;
}

export interface ActionOptionsMap {
    install: types.InstallOptions;
    repair: types.RepairOptions;
    uninstall: types.UninstallOptions;
}

export type Action = keyof ActionOptionsMap;
export type ActionOptions = ActionOptionsMap[Action];

export const createNavState = (state?: Partial<NavigationState>) => {
    return {
        next: undefined,
        previous: undefined,
        nextLabel: "Next",
        previousLabel: "Back",
        canGoNext: true,
        canGoPrevious: true,
        nextAction: undefined,
        previousAction: undefined,
        ...state
    } as NavigationState;
};

export interface AppState {
    eulaAgreed: boolean;
    action: InstallerAction;
    devUnlocked: boolean;
    resourcePaths: Record<DiscordChannel, string>;
    channels: Record<DiscordChannel, boolean>;
    options: ActionOptionsMap;
    navigation: NavigationState;
}
