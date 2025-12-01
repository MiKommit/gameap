<template>
    <div class="fm-content">
        <bread-crumb v-bind:manager="manager" />
        <div class="fm-content-body">
            <table-view v-if="viewType === 'table'" v-bind:manager="manager" />
            <grid-view v-else v-bind:manager="manager" />
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useManager } from '../../composables/useManager.js'

import DiskList from './DiskList.vue'
import BreadCrumb from './BreadCrumb.vue'
import TableView from './TableView.vue'
import GridView from './GridView.vue'

const props = defineProps({
    manager: { type: String, required: true },
})

const { viewType } = useManager(props.manager)
</script>

<style lang="scss">
.fm-content {
    display: flex;
    flex-direction: column;
    height: 100%;

    .fm-content-body {
        flex: 1 1 0;
        min-height: 0;
        overflow-y: auto;
    }
}
</style>
