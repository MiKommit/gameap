<template>
    <div class="fm-additions-file-list">
        <div class="d-flex justify-content-between" v-for="(item, index) in selectedItems" v-bind:key="index">
            <div class="w-75 text-truncate">
                <span v-if="item.type === 'dir'"> <i class="fa-regular fa-folder"></i>{{ item.basename }} </span>
                <span v-else>
                    <i v-bind:class="extensionToIcon(item.extension)" /> {{ item.basename }}
                </span>
            </div>
            <div class="text-end" v-if="item.type === 'file'">
                {{ bytesToHuman(item.size) }}
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useHelper } from '../../../composables/useHelper.js'

const fm = useFileManagerStore()
const { bytesToHuman, extensionToIcon } = useHelper()

const selectedItems = computed(() => fm.selectedItems)
</script>

<style lang="scss">
.fm-additions-file-list {
    .fa {
        padding-right: 0.5rem;
    }
}
</style>
