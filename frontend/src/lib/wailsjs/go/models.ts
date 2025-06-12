export namespace discord {
	
	export enum DiscordChannel {
	    STABLE = 0,
	    CANARY = 1,
	    PTB = 2,
	}
	export interface DiscordInstall {
	    corePath: string;
	    channel: DiscordChannel;
	    version: string;
	    isFlatpak: boolean;
	    isSnap: boolean;
	}

}

