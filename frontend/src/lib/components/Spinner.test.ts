import {render} from "@testing-library/svelte";
import Spinner from "./Spinner.svelte";
import {describe, expect, it} from "vitest";


describe("Spinner", () => {
    it("renders a progressbar with aria attributes", () => {
        const {getByRole} = render(Spinner);

        const spinner = getByRole("progressbar");
        expect(spinner).toHaveAttribute("aria-valuenow");
    });
});
