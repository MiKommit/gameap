<template>
    <div class="modal-content fm-modal-preview">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title w-75 text-truncate">
                {{ showCropperModule ? lang.modal.cropper.title : lang.modal.preview.title }}
                <small class="text-muted pl-3">{{ selectedItem?.basename }}</small>
            </h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body flex text-center justify-center items-center">
            <template v-if="showCropperModule">
                <cropper-module v-bind:imgSrc="imgSrc" v-bind:maxHeight="maxHeight" v-on:closeCropper="closeCropper" />
            </template>
            <transition v-else name="fade" mode="out-in">
                <div class="spinner-border spinner-border-lg text-muted my-2" v-if="!imgSrc">
                    <span class="visually-hidden">Loading...</span>
                </div>
                <img
                    v-else
                    v-bind:src="imgSrc"
                    v-bind:alt="selectedItem?.basename"
                    v-bind:style="{ 'max-height': maxHeight + 'px' }"
                />
            </transition>
        </div>
        <div v-if="showFooter" class="d-flex justify-content-between">
            <span class="d-block">
                <button
                    type="button"
                    class="btn btn-info rounded mr-2"
                    v-bind:title="lang.modal.cropper.title"
                    v-on:click="showCropperModule = true"
                >
                    <i class="fa-solid fa-crop-simple"></i>
                </button>
            </span>
            <span class="d-block">
                <button type="button" class="btn btn-light rounded" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
            </span>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import CropperModule from '../additions/CropperModule.vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useSettingsStore } from '../../../stores/useSettingsStore.js'
import { useModalStore } from '../../../stores/useModalStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'
import GET from '../../../http/get.js'

const fm = useFileManagerStore()
const settings = useSettingsStore()
const modal = useModalStore()
const { lang } = useTranslate()
const { hideModal } = useModal()

const showCropperModule = ref(false)
const imgSrc = ref(null)

const auth = computed(() => settings.authHeader)
const selectedDisk = computed(() => fm.selectedDisk)
const selectedItem = computed(() => fm.selectedItems[0])

const showFooter = computed(() => {
    if (!selectedItem.value?.extension) return false
    return canCrop(selectedItem.value.extension) && !showCropperModule.value
})

const maxHeight = computed(() => {
    if (modal.modalBlockHeight) {
        return modal.modalBlockHeight - 170
    }
    return 300
})

function canCrop(extension) {
    return settings.cropExtensions.includes(extension.toLowerCase())
}

function closeCropper() {
    showCropperModule.value = false
    loadImage()
}

function loadImage() {
    if (auth.value) {
        GET.preview(selectedDisk.value, selectedItem.value.path).then((response) => {
            const mimeType = response.headers['content-type'].toLowerCase()
            const imgBase64 = btoa(String.fromCharCode.apply(null, new Uint8Array(response.data)))
            imgSrc.value = `data:${mimeType};base64,${imgBase64}`
        })
    } else {
        imgSrc.value = `${settings.baseUrl}preview?disk=${selectedDisk.value}&path=${encodeURIComponent(selectedItem.value.path)}&v=${selectedItem.value.timestamp}`
    }
}

onMounted(() => {
    loadImage()
})
</script>

<style lang="scss">
.fm-modal-preview {
    .modal-body {
        padding: 0;

        img {
            max-width: 100%;
        }
    }

    & > .d-flex {
        padding: 1rem;
        border-top: 1px solid #e9ecef;
    }
}
</style>
