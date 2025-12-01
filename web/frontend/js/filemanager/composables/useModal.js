import { computed } from 'vue'
import { useFileManagerStore } from '../stores/useFileManagerStore.js'
import { useModalStore } from '../stores/useModalStore.js'

/**
 * Modal composable for common modal functionality
 * Replaces modals/mixins/modal.js
 */
export function useModal() {
    const fm = useFileManagerStore()
    const modal = useModalStore()

    const activeManager = computed(() => fm.activeManager)

    function hideModal() {
        modal.setModalState({
            modalName: null,
            show: false,
        })
    }

    // Vue 3 directive for autofocus
    const vFocus = {
        mounted(el) {
            el.focus()
        },
    }

    return {
        activeManager,
        hideModal,
        vFocus,
    }
}
