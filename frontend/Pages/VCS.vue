<template>
<div class="vcs">
    <div class="page-header">
        <h1>Repos</h1>
        <p class="subtitle">
            Repos are managed per account — pick an account's provider on its Repos tab.
            This page uses the server-wide default provider.
        </p>
    </div>

    <div class="grid">
        <div class="card">
            <h3 class="card-title">Repositories</h3>
            <ul class="repo-list">
                <li v-for="r in repos" :key="r" class="repo-item">{{ r }}</li>
            </ul>
            <p v-if="repos.length === 0" class="empty">No repos found.</p>
            <p v-if="vcsError" class="error-msg">{{ vcsError }}</p>
        </div>

        <div v-if="vcsAvailable" class="card">
            <h3 class="card-title">Create Repo</h3>
            <form @submit.prevent="createRepo">
                <TextInput
                    label="Repository Name"
                    v-model="repoForm.name"
                    placeholder="my-repo"
                    :required="true" />
                <button type="submit" class="btn-primary" :disabled="!repoForm.name || repoForm.processing">Create</button>
            </form>
        </div>

        <div v-if="vcsAvailable" class="card">
            <h3 class="card-title">Add Deploy Key</h3>
            <form @submit.prevent="addKey">
                <SelectInput
                    label="Repository"
                    v-model="keyForm.repo"
                    :options="repos"
                    :required="true" />
                <TextInput
                    label="Key Name"
                    v-model="keyForm.name"
                    placeholder="deploy-key-1"
                    :required="true" />
                <textarea
                    v-model="keyForm.key"
                    class="key-textarea"
                    placeholder="ssh-rsa AAAA..."
                    rows="4"
                    required></textarea>
                <CheckboxInput
                    label="Read-only"
                    v-model:checked="keyForm.readOnly" />
                <button type="submit" class="btn-primary" :disabled="!keyForm.repo || !keyForm.name || !keyForm.key || keyForm.processing">Add Key</button>
            </form>
        </div>
    </div>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import SelectInput from "@/Components/SelectInput.vue";
import CheckboxInput from "@/Components/CheckboxInput.vue";
import { useForm } from '@inertiajs/vue3';

defineOptions({ layout: Layout });

const props = defineProps({
    vcsError: {
        type: String,
        default: "",
    },
    repos: {
        type: Array,
        default: [],
    },
    vcsAvailable: {
        type: Boolean,
        default: false,
    },
});

const repoForm = useForm({
    name: "",
});

const keyForm = useForm({
    repo: "",
    name: "",
    key: "",
    readOnly: true,
});

const createRepo = () => {
    repoForm.post("/vcs/repos", {
        onSuccess: (page) => { if (!page.props.error) repoForm.reset(); },
    });
};

const addKey = () => {
    keyForm.post("/vcs/repos/deploy-keys", {
        onSuccess: (page) => { if (!page.props.error) keyForm.reset(); },
    });
};
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.page-header {
    margin-bottom: 32px;

    h1 {
        font-size: 24px;
        font-weight: 700;
        letter-spacing: -0.02em;
    }

    .subtitle {
        font-size: 13px;
        color: $c-od-muted;
        margin-top: 6px;
    }
}

.grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;

    @media (max-width: 800px) {
        grid-template-columns: 1fr;
    }
}

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 24px;
}

.card-title {
    font-size: 15px;
    font-weight: 600;
    margin-bottom: 20px;
}

.repo-list {
    list-style: none;
}

.repo-item {
    padding: 12px 20px;
    font-size: 14px;
    color: $c-od-fg;
    border-bottom: 1px solid $c-od-border;
    font-family: monospace;

    &:last-child { border-bottom: none; }
}

.empty, .error-msg {
    padding: 20px;
    font-size: 14px;
    color: $c-od-muted;
}

.error-msg { color: $c-od-danger; }

form {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.key-textarea {
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    padding: 11px 14px;
    color: $c-od-fg;
    font-size: 13px;
    font-family: monospace;
    width: 100%;
    outline: none;
    resize: vertical;
    transition: border-color 0.15s;

    &:focus {
        border-color: $c-od-accent;
    }

    &::placeholder {
        color: $c-od-muted;
        opacity: 0.5;
    }
}

.btn-primary {
    display: inline-block;
    padding: 10px 20px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;
    align-self: flex-start;

    &:hover { opacity: 0.92; }

    &:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }
}
</style>
