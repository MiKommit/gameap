<template>
    <div class="modal-content fm-modal-folder">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.newFile.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div class="form-group">
                <label for="fm-file-name">{{ lang.modal.newFile.fieldName }}</label>
                <input
                    type="text"
                    class="form-control"
                    id="fm-file-name"
                    ref="fileNameInput"
                    v-bind:class="{ 'is-invalid': fileExist }"
                    v-model="fileName"
                    v-on:keyup="validateFileName"
                />
                <div class="invalid-feedback" v-show="fileExist">
                    {{ lang.modal.newFile.fieldFeedback }}
                </div>
            </div>
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-info rounded mr-2" v-bind:disabled="!submitActive" v-on:click="addFile">
                {{ lang.btn.submit }}
            </button>
            <button type="button" class="btn btn-light rounded" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
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

const fileName = ref('')
const fileExist = ref(false)
const fileNameInput = ref(null)

const submitActive = computed(() => fileName.value && !fileExist.value)

onMounted(() => {
    fileNameInput.value?.focus()
})

function validateFileName() {
    if (fileName.value) {
        fileExist.value = fm.fileExist(activeManager.value, fileName.value)
    } else {
        fileExist.value = false
    }
}

function addFile() {
    fm.createFile(fileName.value).then((response) => {
        if (response.data.result.status === 'success') {
            hideModal()
        }
    })
}
</script>
