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

export function tooltip (node: HTMLElement, {
        text = "",
        color = "default",
        position = "top",
        spacing = 3,
        x = 0,
        y = 0
    }: TooltipProps) {

    let isComponentRendered = false;
    let component: ReturnType<typeof Tooltip>;
    let tooltipDOM: HTMLElement;

    async function renderTooltip() {

        let tooltipsLayer = document.getElementById("tooltips-layer");

        // Create Component
        // component = Tooltip({
        //     target: node}, {
        //         text,
        //         color,
        //         position,
        //         x,
        //         y
        // });

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
                component.setCoords(
                    node.getBoundingClientRect().left + (node.offsetWidth / 2) - (tooltipDOM.offsetWidth / 2),
                    (node.getBoundingClientRect().top - tooltipDOM.offsetHeight - 5) - spacing
                );
            }
            else if (position === "bottom") {
                component.setCoords(
                    node.getBoundingClientRect().left + (node.offsetWidth / 2) - (tooltipDOM.offsetWidth / 2),
                    (node.getBoundingClientRect().bottom + 5) + spacing
                );
            }
            else if (position === "left") {
                component.setCoords(
                    (node.getBoundingClientRect().left - tooltipDOM.offsetWidth - 5) - spacing,
                    node.getBoundingClientRect().top + (node.offsetHeight / 2) - (tooltipDOM.offsetHeight / 2)
                );
            }
            else if (position === "right") {
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

        const tooltipsLayer = document.getElementById("tooltips-layer");

        // Check if component is already rendered to prevent warnings
        if (isComponentRendered) {

            // Remove component
            unmount(component);
            tooltipsLayer?.remove();

            // Tooltip is no longer rendered, update our check
            isComponentRendered = false;
        }
    }

    // Add listeners for rendering/unrendering
    node.addEventListener("mouseenter", renderTooltip);
    node.addEventListener("mouseleave", unmountTooltip);
    node.addEventListener("focus", renderTooltip);
    node.addEventListener("blur", unmountTooltip);
    node.childNodes.forEach(child => {
        child.addEventListener("focus", renderTooltip);
    });
    node.childNodes.forEach(child => {
        child.addEventListener("blur", unmountTooltip);
    });

    return {
        destroy() {
            node.removeEventListener("mouseenter", renderTooltip);
            node.removeEventListener("mouseleave", unmountTooltip);
            node.removeEventListener("focus", renderTooltip);
            node.childNodes.forEach(child => {
                child.removeEventListener("focus", renderTooltip);
            });
            node.childNodes.forEach(child => {
                child.removeEventListener("blur", unmountTooltip);
            });
        }
    };
}