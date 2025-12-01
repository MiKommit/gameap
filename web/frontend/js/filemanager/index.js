import FileManager from './FileManager.vue'

export default {
    install: (app) => {
        app.component('file-manager', FileManager)
    },
}
