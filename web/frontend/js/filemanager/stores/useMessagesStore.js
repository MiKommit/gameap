import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useMessagesStore = defineStore('fm-messages', () => {
    const actionResult = ref({
        status: null,
        message: null,
    })
    const actionProgress = ref(0)
    const loadingCount = ref(0)
    const errors = ref([])

    // Getters
    const loading = computed(() => loadingCount.value > 0)

    // Actions
    function setActionResult({ status, message }) {
        actionResult.value.status = status
        actionResult.value.message = message
    }

    function clearActionResult() {
        actionResult.value.status = null
        actionResult.value.message = null
    }

    function setProgress(progress) {
        actionProgress.value = progress
    }

    function clearProgress() {
        actionProgress.value = 0
    }

    function addLoading() {
        loadingCount.value += 1
    }

    function subtractLoading() {
        loadingCount.value -= 1
    }

    function clearLoading() {
        loadingCount.value = 0
    }

    function setError(error) {
        errors.value.push(error)
    }

    function clearErrors() {
        errors.value = []
    }

    return {
        // State
        actionResult,
        actionProgress,
        loadingCount,
        errors,
        // Getters
        loading,
        // Actions
        setActionResult,
        clearActionResult,
        setProgress,
        clearProgress,
        addLoading,
        subtractLoading,
        clearLoading,
        setError,
        clearErrors,
    }
})
