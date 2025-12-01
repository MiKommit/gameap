<template>
    <div class="modal-content fm-modal-text-edit">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title w-75 text-truncate">
                {{ lang.modal.editor.title }}
                <small class="text-muted pl-3">{{ selectedItem?.basename }}</small>
            </h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div v-if="codeLoaded">
                <codemirror
                    ref="fmCodeEditor"
                    v-model="code"
                    :style="{ height: editorHeight + 'px' }"
                    :extensions="extensions"
                    v-on:change="onChange"
                />
            </div>
            <div class="p-5" v-else :style="{ height: editorHeight + 'px' }">
                <div class="d-flex justify-content-center">
                    <div class="spinner-border spinner-border-big" role="status"></div>
                </div>
            </div>
        </div>
        <div class="modal-footer mt-2">
            <button type="button" class="btn btn-info mr-2 rounded" v-on:click="updateFile">
                {{ lang.btn.submit }}
            </button>
            <button type="button" class="btn btn-light rounded" v-on:click="hideModal">
                {{ lang.btn.cancel }}
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { javascript } from '@codemirror/lang-javascript'
import { xml } from '@codemirror/lang-xml'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'

import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useModalStore } from '../../../stores/useModalStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const modal = useModalStore()
const { lang } = useTranslate()
const { hideModal } = useModal()

const code = ref('')
const extensions = [javascript(), xml(), json(), oneDark]
const editedCode = ref('')
const codeLoaded = ref(false)
const fmCodeEditor = ref(null)

const selectedDisk = computed(() => fm.selectedDisk)
const selectedItem = computed(() => fm.selectedItems[0])

const editorHeight = computed(() => {
    if (modal.modalBlockHeight) {
        return modal.modalBlockHeight - 200
    }
    return 300
})

function onChange(value) {
    editedCode.value = value
}

function updateFile() {
    const formData = new FormData()
    formData.append('disk', selectedDisk.value)
    formData.append('path', selectedItem.value.dirname)
    formData.append('file', new Blob([editedCode.value]), selectedItem.value.basename)

    fm.updateFile(formData).then((response) => {
        if (response.data.result.status === 'success') {
            hideModal()
        }
    })
}

onMounted(() => {
    fm.getFile({
        disk: selectedDisk.value,
        path: selectedItem.value.path,
    })
        .then((response) => {
            if (selectedItem.value.extension === 'json') {
                code.value = JSON.stringify(response.data, null, 4)
            } else {
                code.value = response.data
            }
            codeLoaded.value = true
        })
        .catch(() => {
            hideModal()
        })
})
</script>

<style lang="scss">

.fm-modal-text-edit {
    .modal-body {
        padding: 0;
    }
}

.spinner-border-big {
    width: 3rem;
    height: 3rem;
}
</style>
