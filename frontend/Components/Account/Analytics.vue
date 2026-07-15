<template>
<div class="analytics">
    <h2 class="section-title">Analytics</h2>
    <SelectInput width="auto" label="Chart Type" v-model="chartType" :options="[
        { label: 'Bandwidth', val: 'bandwidth' },
        { label: 'Total Requests', val: 'total-requests' },
        { label: 'Unique Visitors', val: 'unique-visitors' },
    ]"
        option-text="label"
        option-value="val" />

    <Spinner v-if="isLoading" />

    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else-if="!chartData || !chartData.datasets || chartData.datasets.length === 0" class="empty-msg">
        No analytics data available.
    </div>

    <div v-else class="chart-container">
        <component :is="currentChartComponent" :data="chartData" :options="currentChartOptions" />
    </div>
</div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useTheme } from '../../composables/useTheme.js'
import axios from 'axios'
import { Bar, Line } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    LineElement,
    CategoryScale,
    LinearScale,
    PointElement
} from 'chart.js'
import Spinner from '../Spinner.vue';
import SelectInput from '../SelectInput.vue';

ChartJS.register(Title, Tooltip, Legend, BarElement, LineElement, CategoryScale, LinearScale, PointElement)

const { theme } = useTheme()

// Chart.js takes JS color values — compute them from the active theme
const C_MUTED  = computed(() => theme.value === 'dark' ? '#8892a4' : '#64748b')
const C_BORDER = computed(() => theme.value === 'dark' ? '#404558' : '#e2e8f0')
const C_FG     = computed(() => theme.value === 'dark' ? '#f5f5f8' : '#1e293b')
const C_ACCENT = '#7c6cf7'
const C_GREEN  = '#4ade80'
const C_AMBER  = '#f59e0b'

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
})

const isLoading = ref(false)
const error = ref(null)
const chartData = ref({ labels: [], datasets: [] })
const chartType = ref('')

const parseDate = (dateStr) => {
    const [day, month, year] = dateStr.split('/')
    const monthMap = { 'Jan': 0, 'Feb': 1, 'Mar': 2, 'Apr': 3, 'May': 4, 'Jun': 5, 'Jul': 6, 'Aug': 7, 'Sep': 8, 'Oct': 9, 'Nov': 10, 'Dec': 11 }
    return new Date(year, monthMap[month], parseInt(day))
}

const formatDate = (date) => {
    const day = String(date.getDate()).padStart(2, '0')
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    return `${day}/${months[date.getMonth()]}/${date.getFullYear()}`
}

const generateDateLabels = (startDateStr, endDateStr, count) => {
    const startDate = parseDate(startDateStr)
    const labels = []
    for (let i = 0; i < count; i++) {
        const date = new Date(startDate)
        date.setDate(startDate.getDate() + i)
        labels.push(formatDate(date))
    }
    return labels
}

const currentChartComponent = computed(() => chartType.value === 'bandwidth' ? Line : Bar)

const chartTitles = {
    'bandwidth': { text: 'Bandwidth Usage Over Time', yLabel: 'Bandwidth (MB)' },
    'total-requests': { text: 'Total Requests', yLabel: 'Requests' },
    'unique-visitors': { text: 'Unique Visitors', yLabel: 'Visitors' },
}

const currentChartOptions = computed(() => {
    const muted = C_MUTED.value
    const border = C_BORDER.value
    const fg = C_FG.value
    const meta = chartTitles[chartType.value]
    return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: false },
            title: meta ? { display: true, text: meta.text, color: fg, font: { size: 14 } } : { display: false }
        },
        scales: {
            x: {
                ticks: { color: muted, maxRotation: 45, minRotation: 45 },
                grid: { color: border },
                title: meta ? { display: true, text: 'Date', color: muted } : { display: false }
            },
            y: {
                ticks: {
                    color: muted,
                    ...(chartType.value === 'bandwidth' ? { callback: v => v.toLocaleString('en-US') + ' MB' } : {})
                },
                grid: { color: border },
                title: meta ? { display: true, text: meta.yLabel, color: muted } : { display: false },
                beginAtZero: true
            }
        }
    }
})

watch(chartType, async () => {
    if (!chartType.value) return
    try {
        isLoading.value = true
        error.value = null
        const res = await axios.get(`/api/accounts/${props.account.name}/analytics`, { params: { type: chartType.value } })
        if (res.data.error) { error.value = res.data.error; return }

        const data = res.data[chartType.value]
        if (!data || !data.values || data.values.length === 0) {
            error.value = 'No analytics data available.'
            chartData.value = { labels: [], datasets: [] }
            return
        }

        const colors = {
            'bandwidth':       { border: C_ACCENT, bg: C_ACCENT + '33' },
            'total-requests':  { border: C_GREEN,  bg: C_GREEN  + '99' },
            'unique-visitors': { border: C_AMBER,  bg: C_AMBER  + '99' },
        }
        const c = colors[chartType.value] ?? colors['bandwidth']

        chartData.value = {
            labels: generateDateLabels(data.start, data.end, data.values.length),
            datasets: [{ label: chartType.value, data: data.values, borderColor: c.border, backgroundColor: c.bg }]
        }
    } catch (err) {
        error.value = err.response?.data?.error || err.response?.data?.message || 'Failed to load analytics data.'
        chartData.value = { labels: [], datasets: [] }
    } finally {
        isLoading.value = false
    }
})

chartType.value = 'bandwidth'
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.analytics {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
}

.chart-container {
    height: 500px;
}

.error-msg {
    font-size: 14px;
    color: $c-od-danger;
}

.empty-msg {
    font-size: 14px;
    color: $c-od-muted;
}
</style>
