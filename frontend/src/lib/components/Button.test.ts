import {createRawSnippet} from "svelte";
import {expect, it, describe} from "vitest";
import {render} from "@testing-library/svelte";
import Button from "./Button.svelte";


describe("Button", () => {
    it("renders content inside the button", () => {
        const label = createRawSnippet(() => ({
            render: () => "<span>Click me</span>"
        }));

        const {getByRole} = render(Button, {
            props: {
                children: label
            }
        });

        expect(getByRole("button")).toHaveTextContent("Click me");
    });

    it("uses the primary style when requested", () => {
        const label = createRawSnippet(() => ({
            render: () => "<span>Primary</span>"
        }));

        const {getByRole} = render(Button, {
            props: {
                style: "primary",
                children: label
            }
        });

        expect(getByRole("button")).toHaveClass("style-primary");
    });
});
