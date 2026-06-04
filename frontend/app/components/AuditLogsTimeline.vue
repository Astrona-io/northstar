<template>
  <div class="flow-root max-h-96 overflow-y-auto pr-2">
    <ul role="list" class="-mb-8">
      <li v-for="(log, idx) in logs" :key="log.id">
        <div class="relative pb-8">
          <span v-if="idx !== logs.length - 1" class="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200 dark:bg-gray-800" aria-hidden="true" />
          <div class="relative flex space-x-3">
            <div>
              <span :class="[getActionColorClass(log.action), 'flex h-8 w-8 items-center justify-center rounded-full ring-8 ring-white dark:ring-gray-900']">
                <UIcon :name="getActionIcon(log.action)" class="h-5 w-5 text-white" />
              </span>
            </div>
            <div class="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
              <div>
                <p class="text-sm text-gray-800 dark:text-gray-200">
                  <span class="font-semibold capitalize text-gray-900 dark:text-white">{{ log.action }}</span> 
                  - {{ formatChanges(log.changes) }}
                </p>
              </div>
              <div class="whitespace-nowrap text-right text-xs text-gray-500 dark:text-gray-400">
                <time :datetime="log.timestamp">{{ formatTime(log.timestamp) }}</time>
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  logs: {
    type: Array,
    default: () => []
  }
})

const getActionColorClass = (action) => {
  switch (action) {
    case 'create':
      return 'bg-green-500'
    case 'update':
      return 'bg-blue-500'
    case 'delete':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

const getActionIcon = (action) => {
  switch (action) {
    case 'create':
      return 'i-heroicons-plus-circle'
    case 'update':
      return 'i-heroicons-pencil-square'
    case 'delete':
      return 'i-heroicons-trash'
    default:
      return 'i-heroicons-question-mark-circle'
  }
}

const formatChanges = (changesStr) => {
  try {
    const parsed = JSON.parse(changesStr)
    const items = []
    for (const [key, val] of Object.entries(parsed)) {
      if (key === 'properties') {
        items.push('Metadata updated')
      } else {
        items.push(`${key} set to "${val}"`)
      }
    }
    return items.join(', ') || 'Modified details'
  } catch {
    return changesStr
  }
}

const formatTime = (ts) => {
  return new Date(ts).toLocaleString()
}
</script>
