<template>
<div class="p-4">
  <h2 class="text-xl font-bold mb-4">Analytics</h2>
  <SelectInput label="Chart Type" v-model="chartType" :options="[
    { label: 'Bandwidth', val: 'bandwidth' },
    { label: 'Total Requests', val: 'total-requests' },
    { label: 'Unique Visitors', val: 'unique-visitors' },
  ]"
      option-text="label"
      option-value="val" />

  <Spinner v-if="isLoading" />

  <div v-else-if="error" class="text-red-600">
    {{ error }}
  </div>

  <div v-else-if="!chartData">
    <p class="text-gray-500">No analytics data available.</p>
  </div>

  <div v-else class="chart-container">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
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

const props = defineProps({
  account: {
    type: Object,
    default: {}
  }
})

const isLoading = ref(false)
const error = ref(null)
const chartData = ref({
  labels: [],
  datasets: []
})
const chartType = ref('')

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      labels: {
        color: '#333'
      }
    }
  },
  scales: {
    x: {
      ticks: {
        color: '#666'
      }
    },
    y: {
      ticks: {
        color: '#666'
      }
    }
  },
}

watch(chartType, async () => {
  try {
    isLoading.value = true;
    console.log('Fetching analytics for', props.account.name)
    const res = await axios.get(`/accounts/${props.account.name}/analytics`, {
      params: { type: chartType.value },
    })
    console.log('Analytics response:', res.data)

    if (res.data.error) {
      error.value = res.data.error
      return
    }

    const data = res.data[chartType.value]
    if (!data || !data.values || data.values.length === 0) {
      error.value = 'No analytics data available.'
      return
    }

    chartData.value = {
      labels: Array.from({ length: data.values.length }, (_, i) => i + 1),
      datasets: [
        {
          label: chartType.value,
          data: data.values,
          borderColor: '#4F46E5',
          backgroundColor: 'rgba(79,70,229)',
          tension: 0.1,
        },
      ],
    }
  } catch (err) {
    console.error('Failed to load analytics', err)
    error.value = err.response?.data?.error || 'Failed to load analytics data.'
  } finally {
    isLoading.value = false
  }
})
chartType.value = 'bandwidth';
</script>

<style scoped>
.chart-container {
  height: 300px;
}
</style>
