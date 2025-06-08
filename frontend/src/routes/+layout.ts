import {dev} from "$app/environment";
import {goto} from "$app/navigation";

export const ssr = false;

if (dev) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (window as any).goto = goto;
}