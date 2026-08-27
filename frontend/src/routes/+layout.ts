/* eslint-disable @typescript-eslint/no-explicit-any */
import {dev} from "$app/environment";
import {goto} from "$app/navigation";
import {EventsOn} from "@wails/runtime";

export const ssr = false;

if (dev) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    (window as any).goto = goto;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    (window as any).listen = EventsOn;
}