import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import GET from '../http/get.js'
import { useSettingsStore } from './useSettingsStore.js'
import { useMessagesStore } from './useMessagesStore.js'

export const useTreeStore = defineStore('fm-tree', () => {
    /**
     * Directory structure:
     * - id: (int) element id
     * - basename: (string) folder name
     * - dirname: (string) directory name
     * - path: (string) path to directory
     * - props: (object) directory properties
     *   - hasSubdirectories: (boolean) has child directories
     *   - subdirectoriesLoaded: (boolean) child directories loaded
     *   - showSubdirectories: (boolean) show or hide subdirectories branch
     * - parentId: (int) parent id
     */
    const directories = ref([])
    const counter = ref(1)
    const tempIndexArray = ref([])

    // Getters
    const findDirectoryIndex = (path) => directories.value.findIndex((el) => el.path === path)

    const visibleDirectories = computed(() => {
        const settings = useSettingsStore()
        if (settings.hiddenFiles) {
            return directories.value
        }
        return directories.value.filter((item) => item.basename.match(/^([^.]).*/i))
    })

    // Actions
    function cleanTree() {
        directories.value = []
        counter.value = 1
    }

    function addDirectories({ directories: dirs, parentId }) {
        dirs.forEach((directory) => {
            directory.id = counter.value
            directory.parentId = parentId
            directory.props.subdirectoriesLoaded = false
            directory.props.showSubdirectories = false
            counter.value += 1
            directories.value.push(directory)
        })
    }

    function replaceDirectories(dirs) {
        directories.value = dirs
    }

    function updateDirectoryProps({ index, props }) {
        for (const property in props) {
            if (Object.prototype.hasOwnProperty.call(props, property)) {
                directories.value[index].props[property] = props[property]
            }
        }
    }

    function addToTempArray(index) {
        tempIndexArray.value.push(index)
    }

    function clearTempArray() {
        tempIndexArray.value = []
    }

    async function initTree(disk) {
        const response = await GET.tree(disk, null)
        if (response.data.result.status === 'success') {
            if (directories.value.length) {
                cleanTree()
            }
            addDirectories({
                parentId: 0,
                directories: response.data.directories,
            })
        }
    }

    function addToTree({ parentPath, newDirectory }) {
        const messages = useMessagesStore()

        if (parentPath) {
            const parentDirectoryIndex = findDirectoryIndex(parentPath)

            if (parentDirectoryIndex !== -1) {
                addDirectories({
                    directories: newDirectory,
                    parentId: directories.value[parentDirectoryIndex].id,
                })
                updateDirectoryProps({
                    index: parentDirectoryIndex,
                    props: {
                        hasSubdirectories: true,
                        showSubdirectories: true,
                        subdirectoriesLoaded: true,
                    },
                })
            } else {
                messages.setError({ message: 'Directory not found' })
            }
        } else {
            addDirectories({
                directories: newDirectory,
                parentId: 0,
            })
        }
    }

    function subDirsFinder(parentId) {
        directories.value.forEach((item, index) => {
            if (item.parentId === parentId) {
                addToTempArray(index)
                if (item.props.hasSubdirectories) {
                    subDirsFinder(item.id)
                }
            }
        })
    }

    function deleteFromTree(dirsToDelete) {
        dirsToDelete.forEach((item) => {
            const directoryIndex = findDirectoryIndex(item.path)

            if (directoryIndex !== -1) {
                addToTempArray(directoryIndex)
                if (directories.value[directoryIndex].props.hasSubdirectories) {
                    subDirsFinder(directories.value[directoryIndex].id)
                }
            }
        })

        const temp = directories.value.filter((item, index) => {
            return tempIndexArray.value.indexOf(index) === -1
        })

        replaceDirectories(temp)
        clearTempArray()
    }

    async function getSubdirectories({ path, parentId, parentIndex }, selectedDisk) {
        const response = await GET.tree(selectedDisk, path)
        if (response.data.result.status === 'success') {
            addDirectories({
                parentId,
                directories: response.data.directories,
            })
            updateDirectoryProps({
                index: parentIndex,
                props: {
                    subdirectoriesLoaded: true,
                },
            })
        }
    }

    async function showSubdirectories(path, selectedDisk) {
        const messages = useMessagesStore()
        const parentDirectoryIndex = findDirectoryIndex(path)

        if (parentDirectoryIndex !== -1) {
            if (directories.value[parentDirectoryIndex].props.subdirectoriesLoaded) {
                updateDirectoryProps({
                    index: parentDirectoryIndex,
                    props: {
                        showSubdirectories: true,
                    },
                })
            } else {
                await getSubdirectories(
                    {
                        path: directories.value[parentDirectoryIndex].path,
                        parentId: directories.value[parentDirectoryIndex].id,
                        parentIndex: parentDirectoryIndex,
                    },
                    selectedDisk
                )
                updateDirectoryProps({
                    index: parentDirectoryIndex,
                    props: {
                        showSubdirectories: true,
                    },
                })
            }
        } else {
            messages.setError({ message: 'Directory not found' })
        }
    }

    function hideSubdirectories(path) {
        const messages = useMessagesStore()
        const parentDirectoryIndex = findDirectoryIndex(path)

        if (parentDirectoryIndex !== -1) {
            updateDirectoryProps({
                index: parentDirectoryIndex,
                props: {
                    showSubdirectories: false,
                },
            })
        } else {
            messages.setError({ message: 'Directory not found' })
        }
    }

    async function reopenPath(path, selectedDisk) {
        if (path) {
            const splitPath = path.split('/')
            for (let i = 0; i < splitPath.length; i += 1) {
                await showSubdirectories(splitPath.slice(0, i + 1).join('/'), selectedDisk)
            }
        }
    }

    return {
        // State
        directories,
        counter,
        tempIndexArray,
        // Getters
        findDirectoryIndex,
        visibleDirectories,
        // Actions
        cleanTree,
        addDirectories,
        replaceDirectories,
        updateDirectoryProps,
        addToTempArray,
        clearTempArray,
        initTree,
        addToTree,
        deleteFromTree,
        subDirsFinder,
        getSubdirectories,
        showSubdirectories,
        hideSubdirectories,
        reopenPath,
    }
})
