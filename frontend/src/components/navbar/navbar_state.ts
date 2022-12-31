import {reactive} from "vue"

export class NavbarState {
    protected path: string
    constructor() {
        this.path = window.location.pathname
    }

    getClass(path: string) {
        return this.path === path
            ? 'is-active'
            : path !== '/' && this.path.search(path) > -1 ? 'is-active' : ''
    }

    setPath(path: string) {
        this.path = path
    }
}

export const navbarState = reactive(new NavbarState())