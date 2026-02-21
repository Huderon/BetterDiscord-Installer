export namespace types {
	
	export enum DiscordChannel {
	    STABLE = 0,
	    CANARY = 1,
	    PTB = 2,
	}
	export interface InstallOptions {
	    restartDiscord: boolean;
	}
	export interface RepairOptions {
	    disablePlugins: boolean;
	}
	export interface UninstallOptions {
	    fullUninstall: boolean;
	}

}

