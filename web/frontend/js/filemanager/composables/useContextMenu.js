import { computed } from 'vue'
import HTTP from '../http/get.js'
import { useFileManagerStore } from '../stores/useFileManagerStore.js'
import { useSettingsStore } from '../stores/useSettingsStore.js'
import { useModalStore } from '../stores/useModalStore.js'

/**
 * Composable for context menu functionality
 * Replaces contextMenu.js, contextMenuActions.js, and contextMenuRules.js mixins
 */
export function useContextMenu() {
    const fm = useFileManagerStore()
    const settings = useSettingsStore()
    const modal = useModalStore()

    // Computed properties
    const selectedDisk = computed(() => fm.selectedDisk)
    const selectedItems = computed(() => fm.selectedItems)
    const selectedDiskDriver = computed(() => fm.disks[selectedDisk.value]?.driver)
    const multiSelect = computed(() => selectedItems.value.length > 1)
    const firstItemType = computed(() => selectedItems.value[0]?.type)

    // Ability checks
    function canView(extension) {
        if (!extension) return false
        return settings.imageExtensions.includes(extension.toLowerCase())
    }

    function canEdit(extension) {
        if (!extension) return false
        return Object.keys(settings.textExtensions).includes(extension.toLowerCase())
    }

    function canAudioPlay(extension) {
        if (!extension) return false
        return settings.audioExtensions.includes(extension.toLowerCase())
    }

    function canVideoPlay(extension) {
        if (!extension) return false
        return settings.videoExtensions.includes(extension.toLowerCase())
    }

    function isZip(extension) {
        if (!extension) return false
        return extension.toLowerCase() === 'zip'
    }

    // Rules (show/hide menu items)
    const openRule = computed(() => !multiSelect.value && firstItemType.value === 'dir')

    const audioPlayRule = computed(() =>
        selectedItems.value.every((elem) => elem.type === 'file') &&
        selectedItems.value.every((elem) => canAudioPlay(elem.extension))
    )

    const videoPlayRule = computed(() =>
        !multiSelect.value && canVideoPlay(selectedItems.value[0]?.extension)
    )

    const viewRule = computed(() =>
        !multiSelect.value &&
        firstItemType.value === 'file' &&
        canView(selectedItems.value[0]?.extension)
    )

    const editRule = computed(() =>
        !multiSelect.value &&
        firstItemType.value === 'file' &&
        canEdit(selectedItems.value[0]?.extension)
    )

    const selectRule = computed(() =>
        !multiSelect.value && firstItemType.value === 'file' && fm.fileCallback
    )

    const downloadRule = computed(() => !multiSelect.value && firstItemType.value === 'file')

    const copyRule = computed(() => true)
    const cutRule = computed(() => true)
    const renameRule = computed(() => !multiSelect.value)
    const pasteRule = computed(() => !!fm.clipboard.type)
    const zipRule = computed(() => selectedDiskDriver.value === 'local')

    const unzipRule = computed(() =>
        selectedDiskDriver.value === 'local' &&
        !multiSelect.value &&
        firstItemType.value === 'file' &&
        isZip(selectedItems.value[0]?.extension)
    )

    const deleteRule = computed(() => true)
    const propertiesRule = computed(() => !multiSelect.value)

    // Actions
    function openAction() {
        fm.selectDirectory(fm.activeManager, {
            path: selectedItems.value[0].path,
            history: true,
        })
    }

    function audioPlayAction() {
        modal.setModalState({ show: true, modalName: 'AudioPlayerModal' })
    }

    function videoPlayAction() {
        modal.setModalState({ show: true, modalName: 'VideoPlayerModal' })
    }

    function viewAction() {
        modal.setModalState({ show: true, modalName: 'PreviewModal' })
    }

    function editAction() {
        modal.setModalState({ show: true, modalName: 'TextEditModal' })
    }

    async function selectAction() {
        const response = await fm.url({
            disk: selectedDisk.value,
            path: selectedItems.value[0].path,
        })
        if (response.data.result.status === 'success') {
            fm.fileCallback(response.data.url)
        }
    }

    async function downloadAction() {
        const tempLink = document.createElement('a')
        tempLink.style.display = 'none'
        tempLink.setAttribute('download', selectedItems.value[0].basename)

        const response = await HTTP.download(selectedDisk.value, selectedItems.value[0].path)
        tempLink.href = window.URL.createObjectURL(new Blob([response.data]))
        document.body.appendChild(tempLink)
        tempLink.click()
        document.body.removeChild(tempLink)
    }

    function copyAction() {
        fm.toClipboard('copy')
    }

    function cutAction() {
        fm.toClipboard('cut')
    }

    function renameAction() {
        modal.setModalState({ show: true, modalName: 'RenameModal' })
    }

    function pasteAction() {
        fm.paste()
    }

    function zipAction() {
        modal.setModalState({ show: true, modalName: 'ZipModal' })
    }

    function unzipAction() {
        modal.setModalState({ show: true, modalName: 'UnzipModal' })
    }

    function deleteAction() {
        modal.setModalState({ show: true, modalName: 'DeleteModal' })
    }

    function propertiesAction() {
        modal.setModalState({ show: true, modalName: 'PropertiesModal' })
    }

    // Get rule by name
    function getRule(name) {
        const rules = {
            open: openRule,
            audioPlay: audioPlayRule,
            videoPlay: videoPlayRule,
            view: viewRule,
            edit: editRule,
            select: selectRule,
            download: downloadRule,
            copy: copyRule,
            cut: cutRule,
            rename: renameRule,
            paste: pasteRule,
            zip: zipRule,
            unzip: unzipRule,
            delete: deleteRule,
            properties: propertiesRule,
        }
        return rules[name]?.value ?? false
    }

    // Get action by name
    function getAction(name) {
        const actions = {
            open: openAction,
            audioPlay: audioPlayAction,
            videoPlay: videoPlayAction,
            view: viewAction,
            edit: editAction,
            select: selectAction,
            download: downloadAction,
            copy: copyAction,
            cut: cutAction,
            rename: renameAction,
            paste: pasteAction,
            zip: zipAction,
            unzip: unzipAction,
            delete: deleteAction,
            properties: propertiesAction,
        }
        return actions[name]
    }

    return {
        // Computed
        selectedDisk,
        selectedItems,
        selectedDiskDriver,
        multiSelect,
        firstItemType,
        // Ability checks
        canView,
        canEdit,
        canAudioPlay,
        canVideoPlay,
        isZip,
        // Rules
        openRule,
        audioPlayRule,
        videoPlayRule,
        viewRule,
        editRule,
        selectRule,
        downloadRule,
        copyRule,
        cutRule,
        renameRule,
        pasteRule,
        zipRule,
        unzipRule,
        deleteRule,
        propertiesRule,
        getRule,
        // Actions
        openAction,
        audioPlayAction,
        videoPlayAction,
        viewAction,
        editAction,
        selectAction,
        downloadAction,
        copyAction,
        cutAction,
        renameAction,
        pasteAction,
        zipAction,
        unzipAction,
        deleteAction,
        propertiesAction,
        getAction,
    }
}
