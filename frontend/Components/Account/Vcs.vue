<template>
<div class="account-vcs">
    <h3 class="section-title">Repos</h3>

    <p v-if="!vcsAvailable" class="empty-msg">
        Version control isn't configured on this server. Set GITEA_URL and GITEA_TOKEN (the default provider),
        or GITHUB_TOKEN, then restart server-run.
    </p>

    <template v-else>
        <div class="provider-row">
            <SelectInput
                label="VCS Provider"
                width="auto"
                v-model="selectedProvider"
                :options="providerOptions"
                option-text="label"
                option-value="val" />
            <button class="btn-ghost" :disabled="!providerChanged || savingProvider" @click="saveProvider">Save</button>
        </div>

        <div class="vcs-grid">
            <div class="panel">
                <h4 class="sub-title">Repositories</h4>
                <ul v-if="repos.length > 0" class="repo-list">
                    <li v-for="repo in repos" :key="repo">{{ repo }}</li>
                </ul>
                <p v-else class="empty-msg">{{ reposError || "No repos found." }}</p>
                <button class="btn-ghost" @click="loadRepos">Refresh</button>
            </div>

            <div class="panel">
                <h4 class="sub-title">Create Repo</h4>
                <form @submit.prevent="createRepo">
                    <TextInput label="Repository Name" v-model="repoName" placeholder="my-repo" :required="true" />
                    <button type="submit" class="btn-primary" :disabled="!repoName || submitting">Create</button>
                </form>
            </div>

            <div class="panel">
                <h4 class="sub-title">Add Deploy Key</h4>
                <form @submit.prevent="addKey">
                    <SelectInput label="Repository" v-model="keyRepo" :options="repos" :required="true" placeholder="Choose a repo" />
                    <TextInput label="Key Name" v-model="keyName" placeholder="deploy-key-1" :required="true" />
                    <textarea
                        v-model="keyValue"
                        class="key-textarea"
                        placeholder="ssh-rsa AAAA..."
                        rows="4"
                        required></textarea>
                    <CheckboxInput label="Read-only" v-model:checked="keyReadOnly" />
                    <button type="submit" class="btn-primary" :disabled="!keyRepo || !keyName || !keyValue || submitting">Add Key</button>
                </form>
            </div>
        </div>
    </template>
</div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { router } from '@inertiajs/vue3'
import axios from 'axios'
import TextInput from '@/Components/TextInput.vue'
import SelectInput from '@/Components/SelectInput.vue'
import CheckboxInput from '@/Components/CheckboxInput.vue'

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    },
    vcsAvailable: {
        type: Boolean,
        default: false,
    },
    providers: {
        type: Array,
        default: [],
    },
})

const selectedProvider = ref(props.account.vcsProvider || "default")
const savingProvider = ref(false)
const submitting = ref(false)

const repos = ref([])
const reposError = ref("")

const repoName = ref("")
const keyRepo = ref("")
const keyName = ref("")
const keyValue = ref("")
const keyReadOnly = ref(true)

const providerOptions = computed(() => [
    { label: "Server default", val: "default" },
    ...props.providers.map(p => ({ label: p, val: p })),
])

const providerChanged = computed(() => selectedProvider.value !== (props.account.vcsProvider || "default"))

const loadRepos = async () => {
    reposError.value = ""
    try {
        const { data } = await axios.get(`/accounts/${props.account.name}/vcs/repos`)
        repos.value = data.repos ?? []
    } catch (e) {
        repos.value = []
        reposError.value = e.response?.data?.error ?? "Couldn't reach the VCS provider."
    }
}

const saveProvider = () => {
    savingProvider.value = true
    router.patch(`/accounts/${props.account.name}/vcs-provider`, {
        provider: selectedProvider.value === "default" ? "" : selectedProvider.value,
    }, {
        onFinish: () => { savingProvider.value = false },
    })
}

const createRepo = () => {
    submitting.value = true
    router.post(`/accounts/${props.account.name}/vcs/repos`, {
        name: repoName.value,
    }, {
        onFinish: () => { submitting.value = false },
    })
}

const addKey = () => {
    submitting.value = true
    router.post(`/accounts/${props.account.name}/vcs/repos/deploy-keys`, {
        repo: keyRepo.value,
        name: keyName.value,
        key: keyValue.value,
        readOnly: keyReadOnly.value,
    }, {
        onFinish: () => { submitting.value = false },
    })
}

onMounted(() => {
    if (props.vcsAvailable) loadRepos()
})
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
    margin-bottom: 16px;
}

.sub-title {
    font-size: 13px;
    font-weight: 600;
    color: $c-od-muted;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin-bottom: 12px;
}

.provider-row {
    display: flex;
    align-items: flex-end;
    gap: 10px;
    margin-bottom: 24px;
}

.vcs-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;

    @media (max-width: 900px) {
        grid-template-columns: 1fr;
    }
}

.panel {
    border: 1px solid $c-od-border;
    border-radius: 10px;
    padding: 16px;
}

.repo-list {
    list-style: none;
    margin-bottom: 12px;

    li {
        font-family: monospace;
        font-size: 13px;
        padding: 6px 0;
        border-bottom: 1px solid $c-od-border;

        &:last-child { border-bottom: none; }
    }
}

.empty-msg {
    font-size: 14px;
    color: $c-od-muted;
    margin-bottom: 12px;
}

form {
    display: flex;
    flex-direction: column;
    gap: 12px;
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

    &:focus {
        border-color: $c-od-accent;
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

.btn-ghost {
    display: inline-block;
    padding: 10px 20px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
    align-self: flex-start;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }

    &:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }
}
</style>
