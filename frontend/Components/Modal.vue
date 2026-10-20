<template>
<div class="modal">
    <div class="modal-contents" :class="'size-' + size">
        <slot />
    </div>
</div>
</template>

<script setup>
defineProps({
    // 'full' (default) keeps the historical 90vw/90vh sheet used by the rest of
    // the app; 'compact' shrinks to fit content — right for small confirmations
    // where a full-screen sheet is overwhelming.
    size: {
        type: String,
        default: 'full',
        validator: (v) => ['full', 'compact'].includes(v)
    }
})
</script>

<style scoped lang="scss">
@import "../scss/variables.scss";

.modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: var(--c-od-overlay);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;

    .modal-contents {
        background: $c-od-surface;
        border: 1px solid $c-od-border;
        border-radius: 14px;
        padding: 1.5rem;
        overflow: auto;

        &.size-full {
            width: 90vw;
            height: 90vh;
        }

        &.size-compact {
            width: auto;
            height: auto;
            max-width: min(520px, 92vw);
            max-height: min(80vh, 92vh);
        }
    }
}
</style>
