<template>
    <div class="fm-additions-cropper">
        <div class="row" v-bind:style="{ 'max-height': maxHeight + 'px' }">
            <div class="col-sm-9 cropper-block">
                <img v-bind:src="imgSrc" ref="fmCropper" v-bind:alt="selectedItem?.basename" />
            </div>
            <div class="col-sm-3 ps-0">
                <div class="cropper-preview"></div>
                <div class="cropper-data">
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataX">X</label>
                        <input v-model.number="x" type="text" class="form-control" id="dataX" />
                        <span class="input-group-text">px</span>
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataY">Y</label>
                        <input v-model.number="y" type="text" class="form-control" id="dataY" />
                        <span class="input-group-text">px</span>
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataWidth">Width</label>
                        <input v-model.number="width" type="text" class="form-control" id="dataWidth" />
                        <span class="input-group-text">px</span>
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataHeight">Height</label>
                        <input v-model.number="height" type="text" class="form-control" id="dataHeight" />
                        <span class="input-group-text">px</span>
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataRotate">Rotate</label>
                        <input v-model.number="rotate" type="text" class="form-control" id="dataRotate" />
                        <span class="input-group-text">deg</span>
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataScaleX">ScaleX</label>
                        <input v-model.number="scaleX" type="text" class="form-control" id="dataScaleX" />
                    </div>
                    <div class="input-group input-group-sm">
                        <label class="input-group-text" for="dataScaleY">ScaleY</label>
                        <input v-model.number="scaleY" type="text" class="form-control" id="dataScaleY" />
                    </div>
                    <div class="d-grid gap-2">
                        <button
                            v-on:click="setData()"
                            v-bind:title="lang.modal.cropper.apply"
                            type="button"
                            class="btn btn-block btn-sm btn-info mb-2"
                        >
                            <i class="fa-solid fa-check"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>
        <div class="d-flex justify-content-between">
            <div>
                <div class="btn-group me-2" role="group" aria-label="Scale">
                    <button v-on:click="cropMove(-10, 0)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-left"></i>
                    </button>
                    <button v-on:click="cropMove(10, 0)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-right"></i>
                    </button>
                    <button v-on:click="cropMove(0, -10)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-up"></i>
                    </button>
                    <button v-on:click="cropMove(0, 10)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-down"></i>
                    </button>
                </div>
                <div class="btn-group me-2" role="group" aria-label="Scale">
                    <button v-on:click="cropScaleX()" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrows-left-right"></i>
                    </button>
                    <button v-on:click="cropScaleY()" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrows-up-down"></i>
                    </button>
                </div>
                <div class="btn-group me-2" role="group" aria-label="Rotate">
                    <button v-on:click="cropRotate(-45)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-rotate-left"></i>
                    </button>
                    <button v-on:click="cropRotate(45)" type="button" class="btn btn-info">
                        <i class="fa-solid fa-arrow-rotate-right"></i>
                    </button>
                </div>
                <div class="btn-group me-2" role="group" aria-label="Rotate">
                    <button v-on:click="cropZoom(0.1)" type="button" class="btn btn-info">
                        <i class="fa-regular fa-circle-plus"></i>
                    </button>
                    <button v-on:click="cropZoom(-0.1)" type="button" class="btn btn-info">
                        <i class="fa-regular fa-circle-minus"></i>
                    </button>
                </div>
                <button
                    v-on:click="cropReset()"
                    v-bind:title="lang.modal.cropper.reset"
                    type="button"
                    class="btn btn-info me-2"
                >
                    <i class="fa-solid fa-rotate"></i>
                </button>
                <button
                    v-on:click="cropSave()"
                    v-bind:title="lang.modal.cropper.save"
                    type="button"
                    class="btn btn-danger me-2"
                >
                    <i class="fa-regular fa-floppy-disk"></i>
                </button>
            </div>
            <span class="d-block">
                <button v-on:click="emit('closeCropper')" type="button" class="btn btn-light">
                    {{ lang.btn.back }}
                </button>
            </span>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import Cropper from 'cropperjs'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'

const props = defineProps({
    imgSrc: { required: true },
    maxHeight: { type: Number, required: true },
})

const emit = defineEmits(['closeCropper'])

const fm = useFileManagerStore()
const { lang } = useTranslate()

const fmCropper = ref(null)
const cropper = ref(null)
const height = ref(0)
const width = ref(0)
const x = ref(0)
const y = ref(0)
const rotate = ref(0)
const scaleX = ref(1)
const scaleY = ref(1)

const selectedItem = computed(() => fm.selectedItems[0])

function cropMove(xVal, yVal) {
    cropper.value.move(xVal, yVal)
}

function cropScaleY() {
    cropper.value.scale(1, cropper.value.getData().scaleY === 1 ? -1 : 1)
}

function cropScaleX() {
    cropper.value.scale(cropper.value.getData().scaleX === 1 ? -1 : 1, 1)
}

function cropRotate(grade) {
    cropper.value.rotate(grade)
}

function cropZoom(ratio) {
    cropper.value.zoom(ratio)
}

function cropReset() {
    cropper.value.reset()
}

function setData() {
    cropper.value.setData({
        x: x.value,
        y: y.value,
        width: width.value,
        height: height.value,
        rotate: rotate.value,
        scaleX: scaleX.value,
        scaleY: scaleY.value,
    })
}

function cropSave() {
    cropper.value.getCroppedCanvas().toBlob(
        (blob) => {
            const formData = new FormData()
            formData.append('disk', fm.selectedDisk)
            formData.append('path', selectedItem.value.dirname)
            formData.append('file', blob, selectedItem.value.basename)

            fm.updateFile(formData).then((response) => {
                if (response.data.result.status === 'success') {
                    emit('closeCropper')
                }
            })
        },
        selectedItem.value.extension !== 'jpg' ? `image/${selectedItem.value.extension}` : 'image/jpeg'
    )
}

onMounted(() => {
    cropper.value = new Cropper(fmCropper.value, {
        preview: '.cropper-preview',
        crop: (e) => {
            x.value = Math.round(e.detail.x)
            y.value = Math.round(e.detail.y)
            height.value = Math.round(e.detail.height)
            width.value = Math.round(e.detail.width)
            rotate.value = typeof e.detail.rotate !== 'undefined' ? e.detail.rotate : ''
            scaleX.value = typeof e.detail.scaleX !== 'undefined' ? e.detail.scaleX : ''
            scaleY.value = typeof e.detail.scaleY !== 'undefined' ? e.detail.scaleY : ''
        },
    })
})

onBeforeUnmount(() => {
    if (cropper.value) {
        cropper.value.destroy()
    }
})
</script>

<style lang="scss">
@import 'cropperjs/dist/cropper.css';

.fm-additions-cropper {
    overflow: hidden;

    button > i {
        color: white;
        font-weight: bold;
    }

    & > .row {
        flex-wrap: nowrap;
    }

    .cropper-block {
        overflow: hidden;

        img {
            max-width: 100%;
        }
    }

    .col-sm-3 {
        overflow: auto;

        &::-webkit-scrollbar {
            display: none;
        }
    }

    .cropper-preview {
        margin-bottom: 1rem;
        overflow: hidden;
        height: 200px;

        img {
            max-width: 100%;
        }
    }

    .cropper-data {
        padding-left: 1rem;
        padding-right: 1rem;

        & > .input-group {
            margin-bottom: 0.5rem;
        }

        .input-group > .input-group-text:first-child {
            min-width: 4rem;
        }

        .input-group > .input-group-text:last-child {
            min-width: 3rem;
        }
    }

    & > .d-flex {
        padding: 1rem;
        border-top: 1px solid #e9ecef;
    }
}
</style>
