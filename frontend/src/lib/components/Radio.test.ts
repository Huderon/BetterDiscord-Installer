import {createRawSnippet} from "svelte";
import {render} from "@testing-library/svelte";
import Radio from "./Radio.svelte";
import {describe, expect, it} from "vitest";


describe("Radio", () => {
    it("renders radio input with value and content", () => {
        const label = createRawSnippet(() => ({
            render: () => "<span>Stable</span>"
        }));

        const {getByDisplayValue, getByText} = render(Radio, {
            props: {
                group: "channel",
                value: "stable",
                children: label
            }
        });

        expect(getByDisplayValue("stable")).toBeInTheDocument();
        expect(getByText("Stable")).toBeInTheDocument();
    });
});
