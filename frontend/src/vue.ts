import { library } from "@fortawesome/fontawesome-svg-core"
import { fas } from "@fortawesome/free-solid-svg-icons"
import type { App } from "vue"

export default (app: App) => {
  library.add(fas)
}
