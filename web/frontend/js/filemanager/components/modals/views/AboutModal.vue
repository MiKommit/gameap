<template>
    <div class="modal-content fm-modal-about">
        <div class="modal-header grid grid-cols-2">
            <h5 class="modal-title">{{ lang.modal.about.title }}</h5>
            <button type="button" class="btn-close" aria-label="Close" v-on:click="hideModal">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div class="modal-body">
            <strong>{{ lang.modal.about.name }}</strong>
            <hr />
            <div class="grid grid-cols-2">
                <div class="my-3">{{ lang.modal.about.version }}:</div>
                <div class="my-3">{{ version }}</div>

                <div class="my-3">{{ lang.modal.about.developer }}:</div>
                <div class="my-3">Aleksandr Manekin (alexusmai@gmail.com)</div>

                <div class="my-3">{{ lang.modal.about.developer }}:</div>
                <div class="my-3">Nikita Kuznetsov (nikita.hldm@gmail.com)</div>

                <div class="my-3">GitHub:</div>
                <div class="my-3">
                    <a class="font-medium text-blue-600 dark:text-blue-500 hover:underline" href="https://github.com/gameap/gameap-file-manager" target="_blank">
                        gameap/gameap-file-manager
                    </a>
                    <br />
                    <a class="font-medium text-blue-600 dark:text-blue-500 hover:underline" href="gameap/vue-gameap-file-manager" target="_blank">
                        gameap/vue-laravel-file-manager
                    </a>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useSettingsStore } from '../../../stores/useSettingsStore.js'
import { useTranslate } from '../../../composables/useTranslate.js'
import { useModal } from '../../../composables/useModal.js'

const settings = useSettingsStore()
const { lang } = useTranslate()
const { hideModal } = useModal()

const version = computed(() => settings.version)
</script>
