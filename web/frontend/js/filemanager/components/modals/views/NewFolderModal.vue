<template>
    <div class="modal-content fm-modal-folder">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.newFolder.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div class="form-group">
                <label for="fm-folder-name">{{ lang.modal.newFolder.fieldName }}</label>
                <input
                    type="text"
                    class="form-control"
                    id="fm-folder-name"
                    ref="folderNameInput"
                    v-bind:class="{ 'is-invalid': directoryExist }"
                    v-model="directoryName"
                    v-on:keyup="validateDirName"
                />
                <div class="invalid-feedback" v-show="directoryExist">
                    {{ lang.modal.newFolder.fieldFeedback }}
                </div>
            </div>
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-info rounded mr-2" v-bind:disabled="!submitActive" v-on:click="addFolder">
                {{ lang.btn.submit }}
            </button>
            <button type="button" class="btn btn-light rounded mr-2" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const { lang } = useTranslate()
const { activeManager, hideModal } = useModal()

const directoryName = ref('')
const directoryExist = ref(false)
const folderNameInput = ref(null)

const submitActive = computed(() => directoryName.value && !directoryExist.value)

onMounted(() => {
    folderNameInput.value?.focus()
})

function validateDirName() {
    if (directoryName.value) {
        directoryExist.value = fm.directoryExist(activeManager.value, directoryName.value)
    } else {
        directoryExist.value = false
    }
}

function addFolder() {
    fm.createDirectory(directoryName.value).then((response) => {
        if (response.data.result.status === 'success') {
            hideModal()
        }
    })
}
</script>
