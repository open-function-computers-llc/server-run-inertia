<template>
<div class="modal">
    <div class="modal-contents">
        <slot />
    </div>
</div>
</template>

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
        width: 90vw;
        height: 90vh;
    }
}
</style>
