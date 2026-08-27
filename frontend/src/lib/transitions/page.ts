import app from "$lib/stores/state.svelte";
import {quartInOut} from "svelte/easing";


interface PTProps {
    delay: number;
    duration: number;
    easing(n: number): number;
    x: number;
    y: number;
    out: boolean;
}

export default function page(node: HTMLElement, {delay = 0, duration = 250, easing = quartInOut, x = 550, y = 0, out = false}: Partial<PTProps> = {}) {
    const style = getComputedStyle(node);
    const transform = style.transform === "none" ? "" : style.transform;

    const direction = out ? -1 : 1;
    x = direction * x;
    x = app.navigation.direction * x;

    return {
        delay,
        duration,
        easing,
        css: (t: number) => {
            return `transform: ${transform} translate(${(1 - t) * x}px, ${(1 - t) * y}px);`;
        }
    };
}