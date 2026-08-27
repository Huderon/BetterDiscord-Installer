import {render} from "@testing-library/svelte";
import ProgressBar from "./ProgressBar.svelte";
import {describe, expect, it} from "vitest";


describe("ProgressBar", () => {
    it("renders a determinate width from value and max", () => {
        const {container} = render(ProgressBar, {
            props: {
                "value": 25,
                "max": 100,
                "class": ""
            }
        });

        const fill = container.querySelector(".progress-fill");
        expect(fill).toHaveStyle({width: "25%"});
    });

    it("renders two fills in indeterminate mode", () => {
        const {container} = render(ProgressBar, {
            props: {
                "indeterminate": true,
                "class": ""
            }
        });

        const fills = container.querySelectorAll(".progress-fill");
        expect(fills.length).toBe(2);
        expect(fills[0]).toHaveClass("increase");
        expect(fills[1]).toHaveClass("decrease");
    });
});
