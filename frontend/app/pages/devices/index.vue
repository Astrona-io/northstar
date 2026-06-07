<template>
  <div class="space-y-6 font-sans">
    
    <!-- Header -->
    <div class="flex justify-between items-center pb-4 border-b border-gray-200 dark:border-gray-800">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white font-mono uppercase">
          Device Catalog & Specs
        </h2>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Pre-define standardized hardware specifications and standard interface layouts.
        </p>
      </div>
      <div class="flex items-center gap-2">
        <UButton 
          icon="i-heroicons-cloud-arrow-down" 
          color="gray" 
          variant="subtle"
          @click="goToCatalogExplorer"
        >
          Load Existing Devices
        </UButton>
        <UButton icon="i-heroicons-plus" color="primary" @click="openAddModal">
          Add Device Spec
        </UButton>
      </div>
    </div>

    <!-- VIEW: Standard Table list -->
    <div class="space-y-6">
      <!-- Search / Filter bar -->
      <div class="flex flex-col md:flex-row gap-4 bg-white dark:bg-gray-900 p-4 rounded-lg shadow-sm border border-gray-100 dark:border-gray-800">
        <div class="flex-1">
          <UInput v-model="searchQuery" icon="i-heroicons-magnifying-glass" placeholder="Search catalog by manufacturer or model..." />
        </div>
        <div class="w-full md:w-48">
          <USelect v-model="typeFilter" :options="[{label: 'All Categories', value: ''}, {label: 'Server', value: 'Server'}, {label: 'Router', value: 'Router'}, {label: 'Switch (L2)', value: 'Switch (L2)'}, {label: 'Switch (L3)', value: 'Switch (L3)'}, {label: 'Database', value: 'Database'}, {label: 'Application', value: 'Application'}, {label: 'Other', value: 'Other'}]" option-attribute="label" value-attribute="value" />
        </div>
      </div>

      <!-- Catalog Table -->
      <UCard>
        <UTable :rows="filteredDevices" :columns="columns">
          <template #origin-data="{ row }">
            <UBadge v-if="row.is_imported" color="emerald" variant="subtle" class="font-mono text-[9px] font-bold uppercase flex items-center gap-1 w-fit">
              <UIcon name="i-heroicons-lock-closed" class="h-3 w-3" />
              Standard Spec
            </UBadge>
            <UBadge v-else color="gray" variant="subtle" class="font-mono text-[9px] font-bold uppercase flex items-center gap-1 w-fit">
              <UIcon name="i-heroicons-user" class="h-3 w-3" />
              Manual Spec
            </UBadge>
          </template>
          <template #manufacturer-data="{ row }">
            <span class="font-semibold text-gray-900 dark:text-white">{{ row.manufacturer?.name || 'Generic' }}</span>
          </template>
          <template #model_name-data="{ row }">
            <span class="font-mono text-sm text-primary-600 dark:text-primary-400 font-bold">{{ row.model_name }}</span>
          </template>
          <template #subtype-data="{ row }">
            <UBadge color="gray" variant="subtle" class="font-mono text-[10px] font-bold">{{ row.subtype || 'Other' }}</UBadge>
          </template>
          <template #revision-data="{ row }">
            <UBadge color="blue" variant="subtle" class="font-mono text-xs">v{{ row.revision || 1 }}</UBadge>
          </template>
          <template #ports-data="{ row }">
            <div class="text-xs font-mono max-w-xs truncate">
              <span v-if="row.ports && Object.keys(row.ports).length > 0">
                {{ Object.entries(row.ports).map(([type, count]) => `${count}x ${type}`).join(', ') }}
              </span>
              <span v-else class="text-gray-400 italic">No ports</span>
            </div>
          </template>
          <template #general_info-data="{ row }">
            <span class="text-xs text-gray-500 dark:text-gray-400 truncate max-w-xs block">
              {{ row.general_info || 'No description provided.' }}
            </span>
          </template>
          <template #actions-data="{ row }">
            <UDropdown :items="items(row)">
              <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
            </UDropdown>
          </template>
        </UTable>
      </UCard>
    </div>

    <!-- Creation / Edit Modal -->
    <DeviceFormModal v-model="isOpen" :device="editingDevice" @saved="refresh" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const { fetchDevices, deleteDevice } = useDevices()

// Modern, non-blocking notification toast helper
const toast = useToast()
const notify = (msg, isError = false) => {
  toast.add({
    title: isError ? 'Operation Failed' : 'Success',
    description: msg,
    color: isError ? 'red' : 'green',
    icon: isError ? 'i-heroicons-x-circle' : 'i-heroicons-check-circle'
  })
}

const isOpen = ref(false)
const editingDevice = ref(null)

const searchQuery = ref('')
const typeFilter = ref('')

const { data: devices, refresh } = await fetchDevices()

const goToCatalogExplorer = () => {
  navigateTo('/devices/catalog')
}

const filteredDevices = computed(() => {
  if (!devices.value) return []
  return devices.value.filter(d => {
    const matchesSearch = d.manufacturer?.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
                          d.model_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesType = !typeFilter.value || d.subtype === typeFilter.value
    return matchesSearch && matchesType
  })
})

const openAddModal = () => {
  editingDevice.value = null
  isOpen.value = true
}

const openEditModal = (dev) => {
  editingDevice.value = dev
  isOpen.value = true
}

const removeDevice = async (id) => {
  if (!confirm('Are you sure you want to delete this device specification from the catalog?')) return
  try {
    await deleteDevice(id)
    await refresh()
    notify('Device specification deleted successfully.')
  } catch (error) {
    console.error('Failed to delete device spec:', error)
    notify('Failed to delete device spec.', true)
  }
}

const columns = [
  { key: 'origin', label: 'Origin' },
  { key: 'manufacturer', label: 'Manufacturer' },
  { key: 'model_name', label: 'Model' },
  { key: 'subtype', label: 'Classification' },
  { key: 'revision', label: 'Revision' },
  { key: 'ports', label: 'Ports' },
  { key: 'general_info', label: 'General Specifications' },
  { key: 'actions', label: 'Actions' }
]

const items = (row) => {
  const actions = []
  
  if (row.is_imported) {
    actions.push([{
      label: 'Standard (Immutable)',
      icon: 'i-heroicons-lock-closed',
      disabled: true
    }])
  } else {
    actions.push([{
      label: 'Edit Spec',
      icon: 'i-heroicons-pencil-square-20-solid',
      click: () => openEditModal(row)
    }])
  }
  
  actions.push([{
    label: 'Remove Spec',
    icon: 'i-heroicons-trash-20-solid',
    click: () => removeDevice(row.id)
  }])
  
  return actions
}
</script>
