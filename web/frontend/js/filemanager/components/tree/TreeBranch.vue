<template>
    <ul class="list-unstyled fm-tree-branch">
        <li v-for="(directory, index) in subDirectories" v-bind:key="index">
            <p
                class="unselectable"
                v-bind:class="{ selected: isDirectorySelected(directory.path) }"
                v-on:click="selectDirectoryAction(directory.path)"
            >
                <i
                    class="fa-regular"
                    v-if="directory.props.hasSubdirectories"
                    v-on:click.stop="showSubdirectories(directory.path, directory.props.showSubdirectories)"
                    v-bind:class="[arrowState(index) ? 'fa-square-minus' : 'fa-square-plus']"
                />
                <i class="fa-solid fa-minus" v-else />
                {{ directory.basename }}
            </p>

            <transition name="fade-tree">
                <tree-branch
                    v-show="arrowState(index)"
                    v-if="directory.props.hasSubdirectories"
                    v-bind:parent-id="directory.id"
                >
                </tree-branch>
            </transition>
        </li>
    </ul>
</template>

<script setup>
import { computed } from 'vue'
import { useFileManagerStore } from '../../stores/useFileManagerStore.js'
import { useTreeStore } from '../../stores/useTreeStore.js'

const props = defineProps({
    parentId: { type: Number, required: true },
})

const fm = useFileManagerStore()
const tree = useTreeStore()

const subDirectories = computed(() =>
    tree.directories.filter((item) => item.parentId === props.parentId)
)

function isDirectorySelected(path) {
    return fm.left.selectedDirectory === path
}

function arrowState(index) {
    return subDirectories.value[index].props.showSubdirectories
}

function showSubdirectories(path, showState) {
    if (showState) {
        tree.hideSubdirectories(path)
    } else {
        tree.showSubdirectories(path)
    }
}

function selectDirectoryAction(path) {
    if (!isDirectorySelected(path)) {
        fm.selectDirectory('left', { path, history: true })
    }
}
</script>

<style lang="scss">
.fm-tree-branch {
    display: table;
    width: 100%;
    padding-left: 1rem;

    li > p {
        margin-bottom: 0;
        padding: 0.4rem 0.4rem;
        white-space: nowrap;
        cursor: pointer;

        &:hover,
        &.selected {
            background-color: #f8f9fa;
        }
    }
}

.fade-tree-enter-active,
.fade-tree-leave-active {
    transition: all 0.3s ease;
}

.fade-tree-enter,
.fade-tree-leave-to {
    transform: translateX(20px);
    opacity: 0;
}
</style>
