<template>
    <div class="pubkey">
        <h3>Account SSH PubKey</h3>
        <!-- <div class="pubkey-viewer" :class="{ copying: temporaryCopyAnimationShowing }" @click="sshPubkeyToClipboard(account.sshPubKey)"> -->
        <div class="pubkey-viewer" :class="{ copying: temporaryCopyAnimationShowing }" @click="sshPubkeyToClipboard(tempPubKey)">
            <!-- {{ account.sshPubKey }} -->
            {{ tempPubKey }}
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

const tempPubKey = 'Harcoded for testing copying functionality, yo!'

const temporaryCopyAnimationShowing = ref(false)

function sshPubkeyToClipboard(key) {
    navigator.clipboard.writeText(key)
    temporaryCopyAnimationShowing.value = true
    setTimeout(() => {
        temporaryCopyAnimationShowing.value = false
    }, 1000)
}
</script>
