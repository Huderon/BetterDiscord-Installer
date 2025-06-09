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

let i = 0;
export const handleArrowKeys = (event: KeyboardEvent, container: HTMLDivElement) => {
    container.focus();
    if (container.hasAttribute("data-selected-index")) i = parseInt(container.getAttribute("data-selected-index")!);
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