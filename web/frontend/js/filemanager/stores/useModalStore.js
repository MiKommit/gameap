import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useModalStore = defineStore('fm-modal', () => {
    const showModal = ref(false)
    const modalName = ref(null)
    const modalBlockHeight = ref(0)

    // Actions
    function setModalState({ show, modalName: name }) {
        showModal.value = show
        modalName.value = name
    }

    function clearModal() {
        showModal.value = false
        modalName.value = null
    }

    function setModalBlockHeight(height) {
        modalBlockHeight.value = height
    }

    return {
        // State
        showModal,
        modalName,
        modalBlockHeight,
        // Actions
        setModalState,
        clearModal,
        setModalBlockHeight,
    }
})
