<template>
<div class="pubkey">
    <h3 class="section-title">SSH Public Key</h3>
    <p class="hint">Click to copy to clipboard.</p>
    <div
        class="pubkey-viewer"
        :class="{ copying: temporaryCopyAnimationShowing }"
        @click="sshPubkeyToClipboard(account.sshPubKey)">
        {{ account.sshPubKey }}
        <span class="copied-badge">Copied!</span>
    </div>
</div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
})

const temporaryCopyAnimationShowing = ref(false)

function sshPubkeyToClipboard(key) {
    navigator.clipboard.writeText(key)
    temporaryCopyAnimationShowing.value = true
    setTimeout(() => {
        temporaryCopyAnimationShowing.value = false
    }, 1000)
}
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.pubkey {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
}

.hint {
    font-size: 13px;
    color: $c-od-muted;
}

.pubkey-viewer {
    background: $c-od-bg;
    color: $c-od-muted;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    height: 300px;
    font-family: monospace;
    font-size: 12px;
    padding: 16px;
    overflow-wrap: anywhere;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: border-color 0.2s, color 0.2s;

    &:hover {
        border-color: $c-od-accent;
        color: $c-od-fg;
    }

    &.copying {
        border-color: $c-od-success;

        .copied-badge {
            transform: translateY(0);
            opacity: 1;
        }
    }
}

.copied-badge {
    position: absolute;
    bottom: 12px;
    right: 12px;
    background: $c-od-success;
    color: oklch(18% 0.015 275);
    font-size: 12px;
    font-weight: 600;
    padding: 4px 10px;
    border-radius: 6px;
    transform: translateY(8px);
    opacity: 0;
    transition: transform 0.2s, opacity 0.2s;
}
</style>
