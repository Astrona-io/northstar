<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white">Device Catalog & Specs</h2>
      <UButton icon="i-heroicons-plus" color="primary" @click="openAddModal">
        Add Device Spec
      </UButton>
    </div>

    <!-- Search / Filter bar -->
    <div class="flex flex-col md:flex-row gap-4 mb-6 bg-white dark:bg-gray-900 p-4 rounded-lg shadow-sm border border-gray-100 dark:border-gray-800">
      <div class="flex-1">
        <UInput v-model="searchQuery" icon="i-heroicons-magnifying-glass" placeholder="Search catalog by manufacturer or model..." />
      </div>
      <div class="w-full md:w-48">
        <USelect v-model="typeFilter" :options="[{label: 'All Categories', value: ''}, {label: 'Server', value: 'Server'}, {label: 'Router', value: 'Router'}, {label: 'Switch', value: 'Switch'}, {label: 'Database', value: 'Database'}, {label: 'Application', value: 'Application'}, {label: 'Other', value: 'Other'}]" option-attribute="label" value-attribute="value" />
      </div>
    </div>

    <!-- Catalog Table -->
    <UCard>
      <UTable :rows="filteredDevices" :columns="columns">
        <template #manufacturer-data="{ row }">
          <span class="font-semibold text-gray-900 dark:text-white">{{ row.manufacturer?.name || 'Generic' }}</span>
        </template>
        <template #model_name-data="{ row }">
          <span class="font-mono text-sm text-primary-600 dark:text-primary-400">{{ row.model_name }}</span>
        </template>
        <template #categories-data="{ row }">
          <div class="flex flex-wrap gap-1">
            <UBadge v-for="cat in row.categories" :key="cat.id" color="gray" variant="subtle" class="capitalize">
              {{ cat.name }}
            </UBadge>
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

    <!-- Creation / Edit Modal -->
    <DeviceFormModal v-model="isOpen" :device="editingDevice" @saved="refresh" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const { fetchDevices, deleteDevice } = useDevices()

const isOpen = ref(false)
const editingDevice = ref(null)

const searchQuery = ref('')
const typeFilter = ref('')

const { data: devices, refresh } = await fetchDevices()

const filteredDevices = computed(() => {
  if (!devices.value) return []
  return devices.value.filter(d => {
    const matchesSearch = d.manufacturer?.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
                          d.model_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesType = !typeFilter.value || d.categories?.some(c => c.name === typeFilter.value)
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
  } catch (error) {
    console.error('Failed to delete device spec:', error)
  }
}

const columns = [
  { key: 'manufacturer', label: 'Manufacturer' },
  { key: 'model_name', label: 'Model' },
  { key: 'categories', label: 'Categories' },
  { key: 'general_info', label: 'General Specifications' },
  { key: 'actions', label: 'Actions' }
]

const items = (row) => [
  [{
    label: 'Edit Spec',
    icon: 'i-heroicons-pencil-square-20-solid',
    click: () => openEditModal(row)
  }],
  [{
    label: 'Remove Spec',
    icon: 'i-heroicons-trash-20-solid',
    click: () => removeDevice(row.id)
  }]
]
</script>
