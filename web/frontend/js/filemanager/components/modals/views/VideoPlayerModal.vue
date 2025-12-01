<template>
    <div class="modal-content fm-modal-video-player">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title w-75 text-truncate">
                {{ lang.modal.videoPlayer.title }} <small class="text-muted ps-3">{{ videoFile?.basename }}</small>
            </h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <video ref="fmVideo" controls />
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import Plyr from 'plyr'
import { useFileManagerStore } from '../../../stores/useFileManagerStore.js'
import { useSettingsStore } from '../../../stores/useSettingsStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const fm = useFileManagerStore()
const settings = useSettingsStore()
const { lang } = useTranslate()
const { hideModal } = useModal()

const fmVideo = ref(null)
const player = ref(null)

const selectedDisk = computed(() => fm.selectedDisk)
const videoFile = computed(() => fm.selectedItems[0])

onMounted(() => {
    player.value = new Plyr(fmVideo.value)

    player.value.source = {
        type: 'video',
        title: videoFile.value.filename,
        sources: [
            {
                src: `${settings.baseUrl}/stream-file?disk=${selectedDisk.value}&path=${encodeURIComponent(videoFile.value.path)}`,
                type: `audio/${videoFile.value.extension}`,
            },
        ],
    }
})

onBeforeUnmount(() => {
    if (player.value) {
        player.value.destroy()
    }
})
</script>

<style lang="scss">
@import 'plyr/plyr.scss';

.fm-modal-video-player {
}
</style>
