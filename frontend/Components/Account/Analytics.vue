<template>
<h2 class="text-xl font-bold mb-4">Analytics</h2>
<SelectInput width="auto" label="Chart Type" v-model="chartType" :options="[
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

<div v-else-if="!chartData || !chartData.datasets || chartData.datasets.length === 0">
  <p class="text-gray-500">No analytics data available.</p>
</div>

<div v-else class="chart-container">
  <component :is="currentChartComponent" :data="chartData" :options="currentChartOptions" />
</div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
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

// Helper function to parse date from "DD/MMM/YYYY" format
const parseDate = (dateStr) => {
  const [day, month, year] = dateStr.split('/')
  const monthMap = {
    'Jan': 0, 'Feb': 1, 'Mar': 2, 'Apr': 3, 'May': 4, 'Jun': 5,
    'Jul': 6, 'Aug': 7, 'Sep': 8, 'Oct': 9, 'Nov': 10, 'Dec': 11
  }
  return new Date(year, monthMap[month], parseInt(day))
}

// Helper function to format date to "DD/MMM/YYYY"
const formatDate = (date) => {
  const day = String(date.getDate()).padStart(2, '0')
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const month = months[date.getMonth()]
  const year = date.getFullYear()
  return `${day}/${month}/${year}`
}

// Helper function to generate date range
const generateDateLabels = (startDateStr, endDateStr, count) => {
  const startDate = parseDate(startDateStr)
  const endDate = parseDate(endDateStr)
  const labels = []

  for (let i = 0; i < count; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    labels.push(formatDate(date))
  }

  return labels
}

// Computed property to determine which chart component to use
const currentChartComponent = computed(() => {
  return chartType.value === 'bandwidth' ? Line : Bar
})

// Computed property for chart-specific options
const currentChartOptions = computed(() => {
  const baseOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
        labels: {
          color: '#333'
        }
      }
    },
    scales: {
      x: {
        ticks: {
          color: '#666',
          maxRotation: 45,
          minRotation: 45
        }
      },
      y: {
        ticks: {
          color: '#666'
        }
      }
    }
  }

  // Customize options based on chart type
  switch (chartType.value) {
    case 'bandwidth':
      return {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: 'Bandwidth Usage Over Time',
            color: '#333',
            font: {
              size: 16
            }
          }
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: {
              display: true,
              text: 'Date',
              color: '#666'
            }
          },
          y: {
            ...baseOptions.scales.y,
            title: {
              display: true,
              text: 'Bandwidth (MB)',
              color: '#666'
            },
            ticks: {
              ...baseOptions.scales.y.ticks,
              callback: function (value) {
                return value.toLocaleString('en-US') + ' MB'
              }
            }
          }
        }
      }

    case 'total-requests':
      return {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: 'Total Requests',
            color: '#333',
            font: {
              size: 16
            }
          }
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: {
              display: true,
              text: 'Date',
              color: '#666'
            }
          },
          y: {
            ...baseOptions.scales.y,
            title: {
              display: true,
              text: 'Requests',
              color: '#666'
            },
            beginAtZero: true
          }
        }
      }

    case 'unique-visitors':
      return {
        ...baseOptions,
        plugins: {
          ...baseOptions.plugins,
          title: {
            display: true,
            text: 'Unique Visitors',
            color: '#333',
            font: {
              size: 16
            }
          }
        },
        scales: {
          ...baseOptions.scales,
          x: {
            ...baseOptions.scales.x,
            title: {
              display: true,
              text: 'Date',
              color: '#666'
            }
          },
          y: {
            ...baseOptions.scales.y,
            title: {
              display: true,
              text: 'Visitors',
              color: '#666'
            },
            beginAtZero: true
          }
        }
      }

    default:
      return baseOptions
  }
})

watch(chartType, async () => {
  if (!chartType.value) return;

  try {
    isLoading.value = true;
    error.value = null;
    console.log('Fetching analytics for', props.account.name, 'type:', chartType.value)

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
      chartData.value = { labels: [], datasets: [] }
      return
    }

    // Generate date labels from start to end
    const dateLabels = generateDateLabels(data.start, data.end, data.values.length)

    // Customize colors based on chart type
    let borderColor, backgroundColor
    switch (chartType.value) {
      case 'bandwidth':
        borderColor = '#4F46E5'  // Indigo
        backgroundColor = 'rgba(79, 70, 229, 0.2)'
        break
      case 'total-requests':
        borderColor = '#10B981'  // Green
        backgroundColor = 'rgba(16, 185, 129, 0.6)'
        break
      case 'unique-visitors':
        borderColor = '#F59E0B'  // Amber
        backgroundColor = 'rgba(245, 158, 11, 0.6)'
        break
      default:
        borderColor = '#4F46E5'
        backgroundColor = 'rgba(79, 70, 229, 0.2)'
    }

    chartData.value = {
      labels: dateLabels,  // Use date labels instead of numeric indices
      datasets: [
        {
          label: chartType.value,
          data: data.values,
          borderColor: borderColor,
          backgroundColor: backgroundColor,
          // tension: chartType.value === 'bandwidth' ? 0.4 : 0,  // Smooth line for bandwidth only
        },
      ],
    }
  } catch (err) {
    console.error('Failed to load analytics', err)
    console.error('Error details:', err.response)
    error.value = err.response?.data?.error || err.response?.data?.message || 'Failed to load analytics data.'
    chartData.value = { labels: [], datasets: [] }
  } finally {
    isLoading.value = false
  }
})

chartType.value = 'bandwidth';
</script>

<style scoped>
.chart-container {
  height: 600px;
}
</style>
