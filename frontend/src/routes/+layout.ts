/* eslint-disable @typescript-eslint/no-explicit-any */
import {dev} from "$app/environment";
import {goto} from "$app/navigation";
import {LogTest} from "@backend/Backend";
import {EventsOn} from "@wails/runtime";

export const ssr = false;

if (dev) {
    (window as any).goto = goto;
    (window as any).backend = LogTest;
    (window as any).listen = EventsOn;
}