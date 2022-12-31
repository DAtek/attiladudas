export function getElementById(id: string): HTMLElement {
    const element = document.getElementById(id)
    if (!element) throw `Element with '${id}' not found.`
    return element
}
