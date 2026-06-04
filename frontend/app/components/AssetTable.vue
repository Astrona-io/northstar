<template>
  <UCard>
    <UTable :rows="assets || []" :columns="columns">
      <template #name-data="{ row }">
        <NuxtLink :to="`/assets/${row.id}`" class="text-primary-500 hover:text-primary-600 hover:underline font-medium">
          {{ row.name }}
        </NuxtLink>
      </template>
      <!-- Physical Location details slot (Orphan auditing, Phase 1 UI L3) -->
      <template #location-data="{ row }">
        <span v-if="row.rack?.datacenter" class="text-xs font-medium text-slate-800 dark:text-slate-200">
          📍 {{ row.rack.datacenter.name }} 
          <span class="text-slate-400 font-mono text-[10px]">({{ row.rack.name }} @ {{ row.rack_position_u || '?' }}U)</span>
        </span>
        <span v-else-if="row.host_asset?.name" class="text-xs font-medium text-slate-700 dark:text-slate-300 italic">
          🖥️ Hosted: {{ row.host_asset.name }}
        </span>
        <span v-else class="text-[10px] uppercase font-bold tracking-wider px-1.5 py-0.5 rounded bg-amber-50 dark:bg-amber-950/30 text-amber-600 dark:text-amber-500 font-mono border border-dashed border-amber-300 dark:border-amber-800">
          LOOSE / UNASSIGNED
        </span>
      </template>
      <template #maintenance_status-data="{ row }">
        <UBadge v-if="row.maintenance_status === 'in-progress'" color="red" variant="subtle">In Progress</UBadge>
        <UBadge v-else-if="row.maintenance_status === 'overdue'" color="orange" variant="subtle">Overdue</UBadge>
        <UBadge v-else-if="row.maintenance_status === 'coming'" color="yellow" variant="subtle">Coming</UBadge>
        <span v-else class="text-gray-400 text-sm">None</span>
      </template>
      <template #actions-data="{ row }">
        <UDropdown :items="items(row)">
          <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
        </UDropdown>
      </template>
    </UTable>
  </UCard>
</template>

<script setup>
const { canMutate, canDelete } = useAuth()

const props = defineProps({
  assets: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['edit', 'delete'])

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'name', label: 'Name' },
  { key: 'type', label: 'Type' },
  { key: 'location', label: 'Physical Location' }, // Phase 1 Location mapping column
  { key: 'ip_address', label: 'IP Address' },
  { key: 'maintenance_status', label: 'Maintenance' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions' }
]

const items = (row) => {
  const dropdown = [
    [{
      label: 'View Details',
      icon: 'i-heroicons-eye-20-solid',
      click: () => navigateTo(`/assets/${row.id}`)
    }]
  ]
  
  if (canMutate.value) {
    dropdown.push([{
      label: 'Edit',
      icon: 'i-heroicons-pencil-square-20-solid',
      click: () => emit('edit', row)
    }])
  }
  
  if (canDelete.value) {
    dropdown.push([{
      label: 'Delete',
      icon: 'i-heroicons-trash-20-solid',
      click: () => emit('delete', row.id)
    }])
  }
  
  return dropdown
}
</script>
