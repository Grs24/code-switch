<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import { Browser } from "@wailsio/runtime";
import { marked } from "marked";
import { checkForUpdates } from "../../services/version";
import type { ReleaseInfo } from "../../../bindings/codeswitch/models";

const { t } = useI18n();
const releaseInfo = ref<ReleaseInfo | null>(null);
const isChecking = ref(false);
const showModal = ref(false);
const dismissed = ref(false);
let intervalId: number | null = null;

const hasUpdate = computed(() => releaseInfo.value?.has_update ?? false);
const latestVersion = computed(() => releaseInfo.value?.version ?? "");
const releaseUrl = computed(() => releaseInfo.value?.url ?? "");
const releaseBody = computed(() => releaseInfo.value?.body ?? "");
const parsedReleaseBody = computed(() => {
  if (!releaseBody.value) return "";
  return marked.parse(releaseBody.value);
});
const publishedDate = computed(() => {
  if (!releaseInfo.value?.published_at) return "";
  try {
    const date = new Date(releaseInfo.value.published_at);
    return date.toLocaleDateString();
  } catch {
    return "";
  }
});

const checkUpdate = async () => {
  if (isChecking.value) return;
  isChecking.value = true;
  try {
    const result = await checkForUpdates();
    releaseInfo.value = result;

    // Only show modal if it wasn't dismissed and there's an update
    if (!dismissed.value && result?.has_update) {
      showModal.value = true;
    }
  } catch (error) {
    console.error("Failed to check for updates:", error);
  } finally {
    isChecking.value = false;
  }
};

const openReleaseUrl = () => {
  if (releaseUrl.value) {
    Browser.OpenURL(releaseUrl.value).catch(console.error);
  }
};

const downloadUpdate = () => {
  openReleaseUrl();
  closeModal();
};

const closeModal = () => {
  showModal.value = false;
  dismissed.value = true;
};

onMounted(() => {
  checkUpdate();
  // Check every hour
  intervalId = window.setInterval(checkUpdate, 60 * 60 * 1000);
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});

defineExpose({
  checkUpdate,
  hasUpdate,
  showModal,
});
</script>

<template>
  <!-- Modal Backdrop -->
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="hasUpdate && showModal"
        class="modal-backdrop"
        @click.self="closeModal"
      >
        <div class="modal-dialog" @click.stop>
          <div class="modal-header">
            <h2 class="modal-title">
              {{ t("components.general.update.available") }}
            </h2>
            <button
              class="modal-close"
              @click="closeModal"
              :aria-label="t('components.general.update.dismiss')"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path
                  d="M6 18L18 6M6 6l12 12"
                  stroke-width="2"
                  stroke-linecap="round"
                />
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <div class="version-header">
              <span class="new-badge">NEW</span>
              <h3 class="version-title">New Version {{ latestVersion }}</h3>
            </div>

            <div v-if="parsedReleaseBody" class="release-notes">
              <div class="release-content" v-html="parsedReleaseBody"></div>
            </div>

            <div v-if="publishedDate" class="publish-date">
              {{ publishedDate }}
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn-secondary" @click="closeModal">
              {{ t("components.general.update.dismiss") }}
            </button>
            <button class="btn-primary" @click="downloadUpdate">
              {{ t("components.general.update.download") }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.modal-dialog {
  background: var(--mac-surface);
  border-radius: 16px;
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.3);
  width: min(560px, 100%);
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--mac-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--mac-text);
}

.modal-close {
  border: none;
  background: transparent;
  color: var(--mac-text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.15s;
}

.modal-close:hover {
  background: rgba(0, 0, 0, 0.05);
}

html.dark .modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
}

.modal-close svg {
  width: 18px;
  height: 18px;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.version-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.new-badge {
  background: #ef4444;
  color: white;
  font-size: 0.7rem;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 6px;
  letter-spacing: 0.05em;
}

.version-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--mac-text);
}

.release-notes {
  margin-bottom: 16px;
}

.release-content {
  color: var(--mac-text);
  line-height: 1.6;
  font-size: 0.95rem;
}

.release-content :deep(h1) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-weight: 600;
  font-size: 1.5rem;
}

.release-content :deep(h2) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-weight: 600;
  font-size: 1.25rem;
}

.release-content :deep(h3) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-weight: 600;
  font-size: 1.1rem;
}

.release-content :deep(ul),
.release-content :deep(ol) {
  margin: 0.5em 0;
  padding-left: 1.5em;
}

.release-content :deep(li) {
  margin: 0.25em 0;
}

.release-content :deep(code) {
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
  font-family: "SFMono-Regular", Menlo, Consolas, monospace;
}

html.dark .release-content :deep(code) {
  background: rgba(255, 255, 255, 0.1);
}

.publish-date {
  font-size: 0.85rem;
  color: var(--mac-text-secondary);
  margin-top: 16px;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--mac-border);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-secondary,
.btn-primary {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
  border: none;
}

.btn-secondary {
  background: transparent;
  color: var(--mac-text-secondary);
  border: 1px solid var(--mac-border);
  width: auto;
  height: auto;
  padding: 6px 20px;
}

.btn-secondary:hover {
  background: rgba(0, 0, 0, 0.05);
}

html.dark .btn-secondary:hover {
  background: rgba(255, 255, 255, 0.05);
}

.btn-primary {
  background: #3b82f6;
  color: white;
  width: auto;
  height: auto;
  padding: 6px 20px;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-primary:active {
  transform: scale(0.98);
}

/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-dialog,
.modal-leave-active .modal-dialog {
  transition: transform 0.2s ease;
}

.modal-enter-from .modal-dialog,
.modal-leave-to .modal-dialog {
  transform: scale(0.95);
}
</style>
