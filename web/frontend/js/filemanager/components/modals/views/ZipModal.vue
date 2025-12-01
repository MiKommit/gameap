<template>
    <div class="modal-content fm-modal-zip">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.zip.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <label for="fm-zip-name">{{ lang.modal.zip.fieldName }}</label>
            <div class="input-group mb-3">
                <input
                    type="text"
                    class="form-control"
                    id="fm-zip-name"
                    ref="archiveInput"
                    v-bind:class="{ 'is-invalid': archiveExist }"
                    v-model="archiveName"
                    v-on:keyup="validateArchiveName"
                />
                <div class="input-group-append">
                    <span class="input-group-text">.zip</span>
                </div>
                <div class="invalid-feedback" v-show="archiveExist">
                    {{ lang.modal.zip.fieldFeedback }}
                </div>
            </div>
            <hr />
            <selected-file-list />
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-info rounded mr-2" v-bind:disabled="!submitActive" v-on:click="createArchive">
                {{ lang.btn.submit }}
            </button>
            <button type="button" class="btn btn-light rounded" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import SelectedFileList from '../additions/SelectedFileList.vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const { lang } = useTranslate()
const { activeManager, hideModal } = useModal()

const archiveName = ref('')
const archiveExist = ref(false)
const archiveInput = ref(null)

const submitActive = computed(() => archiveName.value && !archiveExist.value)

onMounted(() => {
    archiveInput.value?.focus()
})

function validateArchiveName() {
    if (archiveName.value) {
        archiveExist.value = fm.fileExist(activeManager.value, `${archiveName.value}.zip`)
    } else {
        archiveExist.value = false
    }
}

function createArchive() {
    fm.zip(`${archiveName.value}.zip`).then(() => {
        hideModal()
    })
}
</script>
