export type DiscordChannel = "stable" | "ptb" | "canary";
export type InstallerAction = "install" | "repair" | "uninstall";
export type OSName = "windows" | "darwin" | "linux";

export const channels = ["stable", "canary", "ptb"] as const;
export const labels = {stable: "Discord", ptb: "Discord PTB", canary: "Discord Canary"} as const;
export const actions: InstallerAction[] = ["install", "repair", "uninstall"] as const;

export enum NavDirection {FORWARDS = 1, BACKWARDS = -1};

export interface NavigationState {
    direction: NavDirection;
    next?: string;
    previous?: string;
    nextLabel?: string;
    previousLabel?: string;
    canGoNext?: boolean;
    canGoPrevious?: boolean;
    nextAction?(): void;
    previousAction?(): void;
}

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
    corePaths: Record<DiscordChannel, string>;
    channels: Record<DiscordChannel, boolean>;
    navigation: NavigationState;
}