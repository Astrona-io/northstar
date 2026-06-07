<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-semibold">Assets Overview</h2>
      <UButton v-if="canMutate" icon="i-heroicons-plus" color="primary" @click="openAddModal">
        Add Asset
      </UButton>
    </div>

    <!-- Search & Filtering Bar -->
    <div class="flex flex-col xl:flex-row gap-4 mb-6 bg-white dark:bg-slate-900 p-4 rounded-md shadow-sm border border-slate-200 dark:border-slate-800">
      <div class="flex-1">
        <UInput v-model="searchQuery" icon="i-heroicons-magnifying-glass" placeholder="Search assets by name or description..." />
      </div>
      <!-- Datacenter Location filter (Phase 1 UI L3) -->
      <div class="w-full xl:w-52">
        <USelect 
          v-model="selectedDatacenterId" 
          :options="dcFilterOptions" 
          option-attribute="label" 
          value-attribute="value" 
        />
      </div>
      <!-- Mapping compliance state filter (Orphan detection, Phase 1 UI L3) -->
      <div class="w-full xl:w-48">
        <USelect 
          v-model="selectedMappingState" 
          :options="[
            {label: 'All Mapping States', value: ''}, 
            {label: 'Housed (Mapped to DC)', value: 'mapped'}, 
            {label: 'Orphaned / Unassigned', value: 'orphaned'}
          ]" 
          option-attribute="label" 
          value-attribute="value" 
        />
      </div>
      <div class="w-full xl:w-40">
        <USelect v-model="selectedType" :options="[{label: 'All Types', value: ''}, {label: 'Server', value: 'Server'}, {label: 'Router', value: 'Router'}, {label: 'Switch', value: 'Switch'}, {label: 'Database', value: 'Database'}, {label: 'Application', value: 'Application'}, {label: 'Kubernetes Deployment', value: 'Kubernetes Deployment'}]" option-attribute="label" value-attribute="value" />
      </div>
      <div class="w-full xl:w-40">
        <USelect v-model="selectedStatus" :options="[{label: 'All Statuses', value: ''}, {label: 'Active', value: 'active'}, {label: 'Inactive', value: 'inactive'}, {label: 'Maintenance', value: 'maintenance'}]" option-attribute="label" value-attribute="value" />
      </div>
      <div>
        <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="resetFilters" v-if="searchQuery || selectedType || selectedStatus || selectedDatacenterId || selectedMappingState">
          Clear Filters
        </UButton>
      </div>
    </div>

    <AssetTable :assets="assets" @edit="openEditModal" @delete="deleteAsset" />

    <AssetFormModal v-model="isOpen" :asset="editingAsset" @saved="refresh" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const { canMutate } = useAuth()
const { fetchAssets, deleteAsset: apiDeleteAsset } = useAssets()
const { fetchManufacturers, fetchDevices } = useDevices()
const { fetchRacks, fetchDatacenters } = useDCIM()

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

// Preload datacenters for search filtering (Phase 1 UI L3)
const { data: datacenters } = await fetchDatacenters()

const isOpen = ref(false)
const editingAsset = ref(null)

const searchQuery = ref('')
const selectedType = ref('')
const selectedStatus = ref('')
const selectedDatacenterId = ref('')
const selectedMappingState = ref('')

const queryParams = computed(() => {
  const p = {}
  if (searchQuery.value) p.search = searchQuery.value
  if (selectedType.value) p.type = selectedType.value
  if (selectedStatus.value) p.status = selectedStatus.value
  if (selectedDatacenterId.value) p.datacenter_id = selectedDatacenterId.value
  if (selectedMappingState.value) p.mapping_state = selectedMappingState.value
  return p
})

const dcFilterOptions = computed(() => {
  const options = [{ label: 'All Datacenters', value: '' }]
  if (datacenters.value) {
    datacenters.value.forEach(dc => {
      options.push({ label: dc.name, value: dc.id })
    })
  }
  return options
})

const resetFilters = () => {
  searchQuery.value = ''
  selectedType.value = ''
  selectedStatus.value = ''
  selectedDatacenterId.value = ''
  selectedMappingState.value = ''
}

const openAddModal = () => {
  editingAsset.value = null
  isOpen.value = true
}

const openEditModal = (asset) => {
  navigateTo(`/assets/${asset.id}?edit=true`)
}

const { data: assets, refresh } = await fetchAssets(queryParams)

const deleteAsset = async (id) => {
  if (!confirm('Are you sure you want to delete this asset?')) return
  try {
    await apiDeleteAsset(id)
    await refresh()
  } catch (error) {
    console.error('Failed to delete asset:', error)
  }
}
</script>
