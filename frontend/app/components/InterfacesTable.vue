<template>
  <UTable :rows="interfaces || []" :columns="columns">
    <template #status-data="{ row }">
      <UBadge :color="row.status === 'up' ? 'green' : 'red'" variant="subtle">
        {{ row.status?.toUpperCase() }}
      </UBadge>
    </template>
    <template #mtu-data="{ row }">
      <span class="font-mono text-xs">{{ row.mtu }} B</span>
    </template>
    <template #speed-data="{ row }">
      <span class="text-xs">{{ row.speed || 'N/A' }}</span>
    </template>
    <template #vlan-data="{ row }">
      <UBadge v-if="row.vlan" color="blue" variant="subtle">
        {{ row.vlan }}
      </UBadge>
      <span v-else class="text-gray-400 text-xs">None</span>
    </template>
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
      </UDropdown>
    </template>
  </UTable>
</template>

<script setup>
const props = defineProps({
  interfaces: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update-status', 'delete'])

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'type', label: 'Type' },
  { key: 'mac_address', label: 'MAC Address' },
  { key: 'ip_address', label: 'IP Address' },
  { key: 'speed', label: 'Speed' },
  { key: 'mtu', label: 'MTU' },
  { key: 'vlan', label: 'VLAN' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions' }
]

const items = (row) => [
  [{
    label: row.status === 'up' ? 'Set Down' : 'Set Up',
    icon: row.status === 'up' ? 'i-heroicons-arrow-down-circle' : 'i-heroicons-arrow-up-circle',
    click: () => emit('update-status', row.id, row.status === 'up' ? 'down' : 'up')
  }],
  [{
    label: 'Delete Port',
    icon: 'i-heroicons-trash-20-solid',
    click: () => emit('delete', row.id)
  }]
]
</script>
