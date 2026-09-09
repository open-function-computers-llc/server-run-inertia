<template>
<div class="spinner" role="status">
    <span class="sr-only">Loading...</span>
</div>
</template>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.spinner {
    width: 28px;
    height: 28px;
    border: 3px solid $c-od-border;
    border-top-color: $c-od-accent;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

.sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
}
</style>
