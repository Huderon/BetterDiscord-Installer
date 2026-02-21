import {describe, expect, it, vi} from "vitest";
import {checkItem, handleArrowKeys, handleKeyboardToggle} from "./handlers";


const createCheckbox = (checked = false) => {
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = checked;
    return checkbox;
};

const createContainer = (count: number) => {
    const container = document.createElement("div");
    for (let index = 0; index < count; index += 1) {
        const wrapper = document.createElement("div");
        wrapper.appendChild(createCheckbox(index === 0));
        container.appendChild(wrapper);
    }
    container.setAttribute("data-selected-index", "0");
    return container;
};

describe("handlers", () => {
    it("checkItem toggles and dispatches change", () => {
        const checkbox = createCheckbox(false);
        const changeHandler = vi.fn();

        checkbox.addEventListener("change", changeHandler);
        checkItem(checkbox);

        expect(checkbox.checked).toBe(true);
        expect(changeHandler).toHaveBeenCalledTimes(1);
    });

    it("handleKeyboardToggle toggles on Enter", () => {
        const checkbox = createCheckbox(false);
        const event = new KeyboardEvent("keydown", {key: "Enter"});

        handleKeyboardToggle(event, checkbox);

        expect(checkbox.checked).toBe(true);
    });

    it("handleArrowKeys cycles with ArrowDown", () => {
        const container = createContainer(3);
        const event = new KeyboardEvent("keydown", {key: "ArrowDown"});

        handleArrowKeys(event, container);

        const selected = container.children[1].children[0] as HTMLInputElement;
        expect(selected.checked).toBe(true);
    });
});
