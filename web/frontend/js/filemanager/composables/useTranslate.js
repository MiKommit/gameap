import { computed } from 'vue'
import { useSettingsStore } from '../stores/useSettingsStore.js'

/**
 * Composable for translations
 * Replaces translate.js mixin
 */
export function useTranslate() {
    const settings = useSettingsStore()

    const lang = computed(() => {
        if (Object.prototype.hasOwnProperty.call(settings.translations, settings.lang)) {
            return settings.translations[settings.lang]
        }
        return settings.translations.en
    })

    return {
        lang,
    }
}
