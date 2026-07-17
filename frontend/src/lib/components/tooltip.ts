import {mount, tick, unmount} from "svelte";
import Tooltip from "./Tooltip.svelte";

type TooltipColor = "default" | "danger" | "accent";
type TooltipPosition = "top" | "bottom" | "left" | "right";
export interface TooltipProps {
    text?: string;
    color?: TooltipColor;
    position?: TooltipPosition;
    spacing?: number;
    x?: number;
    y?: number;
}

export function tooltip(node: HTMLElement, {
    text = "",
    color = "default",
    position = "top",
    spacing = 3,
    x = 0,
    y = 0
}: TooltipProps) {
    if (!text) return;

    let isComponentRendered = false;
    let component: ReturnType<typeof Tooltip>;
    let tooltipDOM: HTMLDivElement;

    async function renderTooltip() {

        let tooltipsLayer = document.getElementById("tooltips-layer");

        component = mount(Tooltip, {
            target: node,
            props: {
                text,
                color,
                position,
                x,
                y
            }
        });

        // Need to await a tick in order for the component to be built and populated
        await tick();

        // eslint-disable-next-line @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-assignment
        tooltipDOM = component.getElement();

        // Tooltip container
        if (!tooltipsLayer) {
            tooltipsLayer = Object.assign(document.createElement("div"), {
                className: "layer-container",
                id: "tooltips-layer"
            });
            document.body.appendChild(tooltipsLayer);
        }

        tooltipsLayer.appendChild(tooltipDOM);

        // Tooltip Positioning
        if (component) {
            if (position === "top") {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-call
                component.setCoords(
                    node.getBoundingClientRect().left + (node.offsetWidth / 2) - (tooltipDOM.offsetWidth / 2),
                    (node.getBoundingClientRect().top - tooltipDOM.offsetHeight - 5) - spacing
                );
            }
            else if (position === "bottom") {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-call
                component.setCoords(
                    node.getBoundingClientRect().left + (node.offsetWidth / 2) - (tooltipDOM.offsetWidth / 2),
                    (node.getBoundingClientRect().bottom + 5) + spacing
                );
            }
            else if (position === "left") {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-call
                component.setCoords(
                    (node.getBoundingClientRect().left - tooltipDOM.offsetWidth - 5) - spacing,
                    node.getBoundingClientRect().top + (node.offsetHeight / 2) - (tooltipDOM.offsetHeight / 2)
                );
            }
            else if (position === "right") {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-call
                component.setCoords(
                    (node.getBoundingClientRect().left + node.offsetWidth + 5) + spacing,
                    node.getBoundingClientRect().top + (node.offsetHeight / 2) - (tooltipDOM.offsetHeight / 2)
                );
            }
        }

        // Indicate that our tooltip instance is now rendered
        isComponentRendered = true;
    }

    function unmountTooltip() {
        // Check if component is already rendered to prevent warnings
        if (isComponentRendered) {

            // Remove the tooltip's own DOM. The shared #tooltips-layer container
            // is intentionally left in place for other/future tooltips.
            void unmount(component);

            // Tooltip is no longer rendered, update our check
            isComponentRendered = false;
        }
    }

    // Add listeners for rendering/unrendering. Keep a stable handler reference
    // so it can actually be removed again in destroy().
    const showTooltip = () => void renderTooltip();
    node.addEventListener("mouseenter", showTooltip);
    node.addEventListener("mouseleave", unmountTooltip);
    node.addEventListener("focus", showTooltip);
    node.addEventListener("blur", unmountTooltip);
    node.childNodes.forEach(child => {
        child.addEventListener("focus", showTooltip);
        child.addEventListener("blur", unmountTooltip);
    });

    return {
        destroy() {
            node.removeEventListener("mouseenter", showTooltip);
            node.removeEventListener("mouseleave", unmountTooltip);
            node.removeEventListener("focus", showTooltip);
            node.removeEventListener("blur", unmountTooltip);
            node.childNodes.forEach(child => {
                child.removeEventListener("focus", showTooltip);
                child.removeEventListener("blur", unmountTooltip);
            });
        }
    };
}