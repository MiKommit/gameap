<template>
    <figure class="fm-thumbnail" ref="thumbnail">
        <transition name="fade" mode="out-in">
            <i v-if="!src" class="far fa-file-image" />
            <img v-else v-bind:src="src" v-bind:alt="file.filename" class="img-thumbnail" />
        </transition>
    </figure>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useSettingsStore } from '../../stores/useSettingsStore.js'
import GET from '../../http/get.js'

const props = defineProps({
    disk: {
        type: String,
        required: true,
    },
    file: {
        type: Object,
        required: true,
    },
})

const settings = useSettingsStore()

const src = ref('')
const thumbnail = ref(null)

const auth = computed(() => settings.authHeader)

function loadImage() {
    if (auth.value) {
        GET.thumbnail(props.disk, props.file.path).then((response) => {
            const mimeType = response.headers['content-type'].toLowerCase()
            const imgBase64 = btoa(String.fromCharCode.apply(null, new Uint8Array(response.data)))
            src.value = `data:${mimeType};base64,${imgBase64}`
        })
    } else {
        src.value = `${settings.baseUrl}thumbnails?disk=${props.disk}&path=${encodeURIComponent(props.file.path)}&v=${props.file.timestamp}`
    }
}

watch(
    () => props.file.timestamp,
    () => {
        loadImage()
    }
)

onMounted(() => {
    if (window.IntersectionObserver) {
        const observer = new IntersectionObserver(
            (entries, obs) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        loadImage()
                        obs.unobserve(thumbnail.value)
                    }
                })
            },
            {
                root: null,
                threshold: 0.5,
            }
        )

        observer.observe(thumbnail.value)
    } else {
        loadImage()
    }
})
</script>

<style lang="scss">
.fm-thumbnail {
    .img-thumbnail {
        width: 88px;
        height: 88px;
    }

    .fade-enter-active,
    .fade-leave-active {
        transition: opacity 0.3s;
    }

    .fade-enter,
    .fade-leave-to {
        opacity: 0;
    }
}
</style>
