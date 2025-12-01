<template>
    <div class="modal-content fm-modal-audio-player">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">
                {{ lang.modal.audioPlayer.title }}
            </h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <audio ref="fmAudio" controls />
            <hr />
            <div
                class="d-flex justify-content-between py-2 px-2"
                v-bind:class="playingIndex === index ? 'bg-light' : ''"
                v-for="(item, index) in audioFiles"
                v-bind:key="index"
            >
                <div class="w-75 text-truncate">
                    <span class="text-muted pr-2">{{ index }}.</span>
                    {{ item.basename }}
                </div>
                <template v-if="playingIndex === index">
                    <div v-if="status === 'playing'">
                        <i v-on:click="togglePlay()" class="fa-solid fa-play active" />
                    </div>
                    <div v-else>
                        <i v-on:click="togglePlay()" class="fa-solid fa-pause" />
                    </div>
                </template>
                <template v-else>
                    <div>
                        <i v-on:click="selectTrack(index)" class="fa-solid fa-play" />
                    </div>
                </template>
            </div>
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

const fmAudio = ref(null)
const player = ref(null)
const playingIndex = ref(0)
const status = ref('paused')

const selectedDisk = computed(() => fm.selectedDisk)
const audioFiles = computed(() => fm.selectedItems)

function setSource(index) {
    player.value.source = {
        type: 'audio',
        title: audioFiles.value[index].filename,
        sources: [
            {
                src: `${settings.baseUrl}/stream-file?disk=${selectedDisk.value}&path=${encodeURIComponent(audioFiles.value[index].path)}`,
                type: `audio/${audioFiles.value[index].extension}`,
            },
        ],
    }
}

function selectTrack(index) {
    if (player.value.playing) {
        player.value.stop()
    }
    setSource(index)
    player.value.play()
    playingIndex.value = index
}

function togglePlay() {
    player.value.togglePlay()
}

onMounted(() => {
    player.value = new Plyr(fmAudio.value, {
        speed: {
            selected: 1,
            options: [0.5, 1, 1.5],
        },
    })

    setSource(playingIndex.value)

    player.value.on('play', () => {
        status.value = 'playing'
    })

    player.value.on('pause', () => {
        status.value = 'paused'
    })

    player.value.on('ended', () => {
        if (audioFiles.value.length > playingIndex.value + 1) {
            selectTrack(playingIndex.value + 1)
        }
    })
})

onBeforeUnmount(() => {
    if (player.value) {
        player.value.destroy()
    }
})
</script>

<style lang="scss">
@import 'plyr/plyr.scss';

.fm-modal-audio-player {
    .bi.bi-play-fill {
        color: gray;
        opacity: 0.1;
        cursor: pointer;

        &:hover {
            opacity: 0.5;
        }

        &.active {
            opacity: 1;
            color: deepskyblue;
        }
    }
}
</style>
