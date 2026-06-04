<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-semibold">Assets Overview</h2>
      <UButton v-if="canMutate" icon="i-heroicons-plus" color="primary" @click="openAddModal">
        Add Asset
      </UButton>
    </div>

    <!-- Live Network Latency Heatmap Grid (Phase 1 Latency Operations) -->
    <UCard class="mb-6">
      <template #header>
        <div class="flex justify-between items-center text-xs">
          <h3 class="text-sm font-bold uppercase tracking-wider text-slate-800 dark:text-white flex items-center gap-2 font-mono">
            <span class="h-2.5 w-2.5 rounded-full bg-emerald-500 animate-ping" />
            Live Network Latency Monitoring Matrix
          </h3>
          <UButton 
            size="xs" 
            color="gray" 
            variant="subtle" 
            icon="i-heroicons-arrow-path" 
            :loading="pendingPings"
            @click="refreshPings"
          >
            Poll Ping Latency
          </UButton>
        </div>
      </template>

      <!-- The Heatmap Grid -->
      <div v-if="pendingPings" class="flex justify-center p-6">
        <UIcon name="i-heroicons-arrow-path" class="animate-spin h-8 w-8 text-primary-500" />
      </div>

      <div v-else-if="pingMatrix && pingMatrix.length > 0" class="flex flex-wrap gap-4">
        <div 
          v-for="node in pingMatrix" 
          :key="node.id"
          class="flex items-center gap-3 px-4 py-3 border rounded-md shadow-sm cursor-pointer select-none transition-all hover:scale-[1.02] hover:shadow-md"
          :class="[
            node.status === 'offline' 
              ? 'border-red-250 dark:border-red-900 bg-red-50/10 dark:bg-red-950/10 hover:border-red-500' 
              : (node.status === 'maintenance' 
                  ? 'border-yellow-250 dark:border-yellow-900 bg-yellow-50/10 dark:bg-yellow-950/10 hover:border-yellow-500' 
                  : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 hover:border-primary-500')
          ]"
          @click="navigateTo(`/assets/${node.id}`)"
        >
          <!-- Status LED indicator dot -->
          <span 
            class="h-2.5 w-2.5 rounded-full flex-shrink-0"
            :class="[
              node.status === 'offline' ? 'bg-red-500 animate-pulse' : (node.status === 'maintenance' ? 'bg-yellow-500 animate-bounce' : 'bg-emerald-500 animate-pulse')
            ]"
          />

          <!-- Hostname and IP details -->
          <div class="truncate max-w-[120px]">
            <span class="text-xs font-bold text-slate-800 dark:text-white block truncate leading-tight">{{ node.name }}</span>
            <span class="text-[10px] text-slate-400 font-mono block truncate mt-0.5">{{ node.ip_address || 'No IP' }}</span>
          </div>

          <!-- Real-Time latency badge -->
          <div class="text-right flex-shrink-0 pl-1">
            <span 
              v-if="node.status === 'offline'" 
              class="text-[9px] font-bold font-mono text-red-600 dark:text-red-400 uppercase bg-red-100/50 dark:bg-red-950/50 px-1.5 py-0.5 rounded"
            >
              Loss: 100%
            </span>
            <span 
              v-else 
              class="text-[10px] font-bold font-mono px-1.5 py-0.5 rounded uppercase font-semibold"
              :class="[
                node.status === 'maintenance' ? 'text-yellow-600 dark:text-yellow-400 bg-yellow-100/50 dark:bg-yellow-950/50' : 'text-primary-600 dark:text-primary-400 bg-primary-100/30 dark:bg-primary-950/30'
              ]"
            >
              {{ node.latency_ms }} ms
            </span>
          </div>
        </div>
      </div>

      <div v-else class="text-xs text-slate-400 italic py-4 text-center font-mono">
        [ NO LIVE IP NODE ATTRIBUTES LOGGED FOR PING SWEEPS ]
      </div>
    </UCard>

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

    <AssetFormModal v-model="isOpen" :asset="editingAsset" @saved="refreshAndRecheck" />
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
const { data: pingMatrix, pending: pendingPings, refresh: refreshPings } = await useFetch(`${apiBase}/monitoring/ping`)

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

const refreshAndRecheck = async () => {
  await refresh()
  await refreshPings()
}

const deleteAsset = async (id) => {
  if (!confirm('Are you sure you want to delete this asset?')) return
  try {
    await apiDeleteAsset(id)
    await refreshAndRecheck()
  } catch (error) {
    console.error('Failed to delete asset:', error)
  }
}
</script>
