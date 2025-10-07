<template>
    <div class="account-analytics">
        <h3>Account Analytics</h3>

        <div v-if="loading" class="text-gray-500">Loading analytics...</div>
        <div v-else>
            <Line v-if="chartData" :data="chartData" :options="chartOptions" />
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js'
import axios from 'axios'

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement)

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
})

const chartData = ref(null)
const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
})
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await axios.get(`/accounts/${props.account.name}/analytics?type=visitors`)
    const data = res.data.visitors

    chartData.value = {
      labels: Array.from({ length: data.Values.length }, (_, i) => i + 1),
      datasets: [
        {
          label: 'Visitors',
          data: data.Values,
          borderColor: '#4F46E5',
          backgroundColor: 'rgba(79,70,229,0.1)',
          tension: 0.4,
        },
      ],
    }
  } catch (err) {
    console.error('Failed to load analytics', err)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.chart-container {
  height: 300px;
}
</style>
