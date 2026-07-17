export function checkItem(item: HTMLInputElement) {
    item.checked = !item.checked;
    const changeEvent = new Event("change");
    item.dispatchEvent(changeEvent);
}

export const handleKeyboardToggle = (event: KeyboardEvent, checkbox: HTMLInputElement) => {
    if ((event.key === "Enter" || event.key === " ") && !checkbox.disabled) {
        checkItem(checkbox);
    }
};

export const handleArrowKeys = (event: KeyboardEvent, container: HTMLDivElement) => {
    container.focus();
    // Read the current index from the DOM on each call so multiple radio groups
    // don't share a single module-level cursor.
    let i = container.hasAttribute("data-selected-index") ? parseInt(container.getAttribute("data-selected-index")!) : 0;
    if (event.key === "ArrowDown") {
        if (i < (container.children.length - 2)) i++;
        else i = 0;
        checkItem(container.children[i].children[0] as HTMLInputElement);
    }
    if (event.key === "ArrowUp") {
        if (i > 0) i--;
        else i = container.children.length - 2;
        checkItem(container.children[i].children[0] as HTMLInputElement);
    }
};