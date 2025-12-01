<template>
    <div class="modal-content fm-modal-unzip">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.unzip.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div class="d-flex justify-content-between">
                <div>
                    <strong>{{ lang.modal.unzip.fieldRadioName }}</strong>
                </div>
                <div class="form-check form-check-inline">
                    <input
                        class="form-check-input"
                        id="unzipRadio1"
                        type="radio"
                        v-bind:checked="!createFolder"
                        v-on:change="createFolder = false"
                    />
                    <label class="form-check-label" for="unzipRadio1">
                        {{ lang.modal.unzip.fieldRadio1 }}
                    </label>
                </div>
                <div class="form-check form-check-inline">
                    <input
                        class="form-check-input"
                        id="unzipRadio2"
                        type="radio"
                        v-bind:checked="createFolder"
                        v-on:change="createFolder = true"
                    />
                    <label class="form-check-label" for="unzipRadio2">
                        {{ lang.modal.unzip.fieldRadio2 }}
                    </label>
                </div>
            </div>
            <hr />
            <div v-if="createFolder" class="form-group">
                <label for="fm-folder-name">{{ lang.modal.unzip.fieldName }}</label>
                <input
                    type="text"
                    class="form-control"
                    id="fm-folder-name"
                    ref="folderInput"
                    v-bind:class="{ 'is-invalid': directoryExist }"
                    v-model="directoryName"
                    v-on:keyup="validateDirName"
                />
                <div class="invalid-feedback" v-show="directoryExist">
                    {{ lang.modal.unzip.fieldFeedback }}
                </div>
            </div>
            <span v-else class="text-danger">{{ lang.modal.unzip.warning }}</span>
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-info rounded mr-2" v-bind:disabled="!submitActive" v-on:click="unpackArchive">
                {{ lang.btn.submit }}
            </button>
            <button type="button" class="btn btn-light rounded" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const { lang } = useTranslate()
const { activeManager, hideModal } = useModal()

const createFolder = ref(false)
const directoryName = ref('')
const directoryExist = ref(false)
const folderInput = ref(null)

const submitActive = computed(() => {
    if (createFolder.value) {
        return directoryName.value && !directoryExist.value
    }
    return true
})

watch(createFolder, (newVal) => {
    if (newVal) {
        setTimeout(() => folderInput.value?.focus(), 0)
    }
})

function validateDirName() {
    if (directoryName.value) {
        directoryExist.value = fm.directoryExist(activeManager.value, directoryName.value)
    } else {
        directoryExist.value = false
    }
}

function unpackArchive() {
    fm.unzip(createFolder.value ? directoryName.value : null).then(() => {
        hideModal()
    })
}
</script>
