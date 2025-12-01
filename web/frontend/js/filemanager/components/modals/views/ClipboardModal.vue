<template>
    <div class="modal-content fm-modal-clipboard">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.clipboard.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <template v-if="clipboard.type">
                <div class="d-flex justify-content-between">
                    <div class="w-75 text-truncate">
                        <span class="badge bg-light text-dark"
                            ><i class="fa-solid fa-hard-drive pe-1"></i> {{ clipboard.disk }}</span
                        >
                    </div>
                    <div class="text-right text-muted">
                        <span v-bind:title="`${lang.clipboard.actionType} - ${lang.clipboard[clipboard.type]}`">
                            <i v-if="clipboard.type === 'copy'" class="fa-regular fa-copy" />
                            <i v-else class="fa-solid fa-scissors" />
                        </span>
                    </div>
                </div>
                <hr />
                <div
                    class="d-flex justify-content-between"
                    v-for="(dir, index) in directories"
                    v-bind:key="`d-${index}`"
                >
                    <div class="w-75 text-truncate">
                        <span><i class="fa-regular fa-folder"></i> {{ dir.name }}</span>
                    </div>
                    <div class="text-right">
                        <button
                            type="button"
                            class="btn-close"
                            v-bind:title="lang.btn.delete"
                            v-on:click="deleteItem('directories', dir.path)"
                        ></button>
                    </div>
                </div>
                <div class="d-flex justify-content-between" v-for="(file, index) in files" v-bind:key="`f-${index}`">
                    <div class="w-75 text-truncate">
                        <span><i v-bind:class="file.icon" /> {{ file.name }}</span>
                    </div>
                    <div class="text-right">
                        <button
                            type="button"
                            class="btn-close"
                            v-bind:title="lang.btn.delete"
                            v-on:click="deleteItem('files', file.path)"
                        ></button>
                    </div>
                </div>
            </template>
            <template v-else>
                <span>{{ lang.clipboard.none }}</span>
            </template>
        </div>
        <div class="modal-footer">
            <button type="button" class="btn btn-danger rounded mr-2" v-bind:disabled="!clipboard.type" v-on:click="resetClipboardAction">
                {{ lang.btn.clear }}
            </button>
            <button type="button" class="btn btn-light rounded" v-on:click="hideModal">{{ lang.btn.cancel }}</button>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useHelper } from '../../../composables/useHelper.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const { lang } = useTranslate()
const { extensionToIcon } = useHelper()
const { hideModal } = useModal()

const clipboard = computed(() => fm.clipboard)

const directories = computed(() =>
    fm.clipboard.directories.map((item) => ({
        path: item,
        name: item.split('/').slice(-1)[0],
    }))
)

const files = computed(() =>
    fm.clipboard.files.map((item) => {
        const name = item.split('/').slice(-1)[0]
        return {
            path: item,
            name,
            icon: extensionToIcon(name.split('.').slice(-1)[0]),
        }
    })
)

function deleteItem(type, path) {
    fm.truncateClipboard({ type, path })
}

function resetClipboardAction() {
    fm.resetClipboard()
}
</script>

<style lang="scss">
.fm-modal-clipboard {
    .modal-body .far {
        padding-right: 0.5rem;
    }
}
</style>
