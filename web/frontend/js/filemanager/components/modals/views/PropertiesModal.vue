<template>
    <div class="modal-content fm-modal-properties">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.properties.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <div class="grid grid-cols-3 gap-4 my-3">
                <div><strong>{{ lang.modal.properties.disk }}:</strong></div>
                <div>{{ selectedDisk }}</div>
                <div class="text-right">
                    <i
                        v-on:click="copyToClipboard(selectedDisk)"
                        v-bind:title="lang.clipboard.copy"
                        class="fa-regular fa-copy"
                    />
                </div>
            </div>
            <div class="grid grid-cols-3 gap-4 my-3">
                <div><strong>{{ lang.modal.properties.name }}:</strong></div>
                <div>{{ selectedItem.basename }}</div>
                <div class="text-right">
                    <i
                        v-on:click="copyToClipboard(selectedItem.basename)"
                        v-bind:title="lang.clipboard.copy"
                        class="fa-regular fa-copy"
                    />
                </div>
            </div>
            <div class="grid grid-cols-3 gap-4 my-3">
                <div><strong>{{ lang.modal.properties.path }}:</strong></div>
                <div>{{ selectedItem.path }}</div>
                <div class="text-right">
                    <i
                        v-on:click="copyToClipboard(selectedItem.path)"
                        v-bind:title="lang.clipboard.copy"
                        class="fa-regular fa-copy"
                    />
                </div>
            </div>
            <template v-if="selectedItem.type === 'file'">
                <div class="grid grid-cols-3 gap-4 my-3">
                    <div><strong>{{ lang.modal.properties.size }}:</strong></div>
                    <div>{{ bytesToHuman(selectedItem.size) }}</div>
                    <div class="text-right">
                        <i
                            v-on:click="copyToClipboard(bytesToHuman(selectedItem.size))"
                            v-bind:title="lang.clipboard.copy"
                            class="fa-regular fa-copy"
                        />
                    </div>
                </div>
            </template>
            <template v-if="selectedItem.hasOwnProperty('timestamp')">
                <div class="grid grid-cols-3 gap-4 my-3">
                    <div><strong>{{ lang.modal.properties.modified }}:</strong></div>
                    <div>{{ timestampToDate(selectedItem.timestamp) }}</div>
                    <div class="text-right">
                        <i
                            v-on:click="copyToClipboard(timestampToDate(selectedItem.timestamp))"
                            v-bind:title="lang.clipboard.copy"
                            class="fa-regular fa-copy"
                        />
                    </div>
                </div>
            </template>
            <template v-if="selectedItem.hasOwnProperty('acl')">
                <div class="grid grid-cols-3 gap-4 my-3">
                    <div>{{ lang.modal.properties.access }}:</div>
                    <div>{{ lang.modal.properties['access_' + selectedItem.acl] }}</div>
                </div>
            </template>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useHelper } from '../../../composables/useHelper.js'
import { useModal } from '../../../composables/useModal.js'
import { notification } from '@/parts/dialogs.js'

const fm = useFileManagerStore()
const { lang } = useTranslate()
const { bytesToHuman, timestampToDate } = useHelper()
const { hideModal } = useModal()

const selectedDisk = computed(() => fm.selectedDisk)
const selectedItem = computed(() => fm.selectedItems[0])

function copyToClipboard(text) {
    const copyInputHelper = document.createElement('input')
    copyInputHelper.className = 'copyInputHelper'
    document.body.appendChild(copyInputHelper)
    copyInputHelper.value = text
    copyInputHelper.select()
    document.execCommand('copy')
    document.body.removeChild(copyInputHelper)

    notification({
        content: lang.value.notifications.copyToClipboard,
        type: 'success',
    })
}
</script>

<style lang="scss">
.fm-modal-properties .modal-body {
    .row {
        margin-bottom: 0.3rem;
        padding-top: 0.3rem;
        padding-bottom: 0.3rem;

        &:hover {
            background-color: #f8f9fa;
        }
    }
}
</style>
