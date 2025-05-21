<template>
<div class="modal">
    <div class="modal-contents">
        {{ slot }}
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
    background: rgba($c-backgroundDark, 0.75);
    display: flex;
    justify-content: center;
    align-items: center;

    .modal-contents {
        background-color: $c-white;
        padding: 1rem;
        border-radius: 25px;
        overflow: auto;
        width: 90vw;
        height: 90vh;
    }
}
</style>
