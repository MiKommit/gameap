<template>
    <div class="modal-content fm-modal-rename">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.rename.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div class="form-group">
                <label for="fm-input-rename">{{ lang.modal.rename.fieldName }}</label>
                <input
                    type="text"
                    class="form-control"
                    id="fm-input-rename"
                    ref="renameInput"
                    v-bind:class="{ 'is-invalid': checkName }"
                    v-model="name"
                    v-on:keyup="validateName"
                />
                <div class="invalid-feedback" v-show="checkName">
                    {{ lang.modal.rename.fieldFeedback }}
                    {{ directoryExist ? ` - ${lang.modal.rename.directoryExist}` : '' }}
                    {{ fileExist ? ` - ${lang.modal.rename.fileExist}` : '' }}
                </div>
            </div>
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-info rounded mr-2" v-bind:disabled="submitDisable" v-on:click="rename">
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

const name = ref('')
const directoryExist = ref(false)
const fileExist = ref(false)
const renameInput = ref(null)

const selectedItem = computed(() => fm.getSelectedList(activeManager.value)[0])
const checkName = computed(() => directoryExist.value || fileExist.value || !name.value)
const submitDisable = computed(() => checkName.value || name.value === selectedItem.value?.basename)

onMounted(() => {
    name.value = selectedItem.value?.basename || ''
    renameInput.value?.focus()
})

function validateName() {
    if (name.value !== selectedItem.value?.basename) {
        if (selectedItem.value?.type === 'dir') {
            directoryExist.value = fm.directoryExist(activeManager.value, name.value)
        } else {
            fileExist.value = fm.fileExist(activeManager.value, name.value)
        }
    }
}

function rename() {
    const newName = selectedItem.value.dirname
        ? `${selectedItem.value.dirname}/${name.value}`
        : name.value

    fm.rename({
        type: selectedItem.value.type,
        newName,
        oldName: selectedItem.value.path,
    }).then(() => {
        hideModal()
    })
}
</script>
