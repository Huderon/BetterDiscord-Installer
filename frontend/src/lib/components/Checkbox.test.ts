import {fireEvent, render} from "@testing-library/svelte";
import Checkbox from "./Checkbox.svelte";
import {describe, expect, it} from "vitest";


describe("Checkbox", () => {
    it("renders the label and initial checked state", () => {
        const {getByRole, getByText} = render(Checkbox, {
            props: {
                checked: true,
                label: "Enable updates"
            }
        });

        expect(getByText("Enable updates")).toBeInTheDocument();
        expect(getByRole("checkbox")).toBeChecked();
    });

    it("toggles when pressing Enter on the label", async () => {
        const {getByRole, container} = render(Checkbox, {
            props: {
                checked: false,
                label: "Enable updates"
            }
        });

        const label = container.querySelector(".checkbox-container");
        expect(label).toBeTruthy();

        await fireEvent.keyPress(label as HTMLElement, {key: "Enter"});

        expect(getByRole("checkbox")).toBeChecked();
    });
});
