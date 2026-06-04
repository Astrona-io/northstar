<template>
  <UTable :rows="maintenanceWindows" :columns="maintenanceColumns">
    <template #start_time-data="{ row }">
      {{ new Date(row.start_time).toLocaleString() }}
    </template>
    <template #end_time-data="{ row }">
      {{ new Date(row.end_time).toLocaleString() }}
    </template>
      <template #status-data="{ row }">
      <UBadge :color="row.status === 'completed' ? 'green' : row.status === 'cancelled' ? 'red' : 'yellow'" variant="subtle">
        {{ row.status }}
      </UBadge>
    </template>
    <template #actions-data="{ row }">
      <UDropdown :items="maintenanceItems(row)">
        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
      </UDropdown>
    </template>
  </UTable>
</template>

<script setup>
const props = defineProps({
  maintenanceWindows: {
    type: Array,
    required: true
  }
})
const emit = defineEmits(['update-status', 'delete'])

const maintenanceColumns = [
  { key: 'title', label: 'Title' },
  { key: 'start_time', label: 'Start Time' },
  { key: 'end_time', label: 'End Time' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions' }
]

const maintenanceItems = (row) => [
  [{
    label: 'Mark Completed',
    icon: 'i-heroicons-check-circle',
    disabled: row.status !== 'scheduled',
    click: () => emit('update-status', row.id, 'completed')
  },
  {
    label: 'Cancel Window',
    icon: 'i-heroicons-x-circle',
    disabled: row.status !== 'scheduled',
    click: () => emit('update-status', row.id, 'cancelled')
  }],
  [{
    label: 'Delete',
    icon: 'i-heroicons-trash-20-solid',
    click: () => emit('delete', row.id)
  }]
]
</script>