<template>
  <div class="space-y-6 font-sans">
    
    <!-- Header -->
    <div class="flex justify-between items-center pb-4 border-b border-gray-200 dark:border-gray-800">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white font-mono uppercase">
          Hardware Catalog Explorer
        </h2>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Browse, inspect, and instant-import standard Ubiquiti and MikroTik product specifications directly from disk.
        </p>
      </div>
      <div>
        <UButton 
          icon="i-heroicons-arrow-left" 
          color="gray" 
          variant="ghost"
          @click="backToDevices"
        >
          Back to Device Catalog
        </UButton>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isCatalogLoading" class="flex justify-center p-12">
      <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
    </div>



    <!-- VIEW 1: Brand Cards Grid -->
    <div v-else-if="!selectedBrand" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div 
        v-for="b in diskCatalog" 
        :key="b.brand"
        class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-lg p-6 flex items-center justify-between shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
        @click="selectedBrand = b"
      >
        <div class="flex items-center gap-4">
          <!-- Stylized Brand Logo Cards -->
          <div v-if="b.logo === 'ubiquiti'" class="h-16 w-16 bg-blue-50 dark:bg-blue-950/40 rounded-lg flex items-center justify-center border border-blue-200 dark:border-blue-900/60 shadow-inner group-hover:scale-105 transition-all">
            <svg class="h-10 w-10 text-blue-500 animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div v-else-if="b.logo === 'mikrotik'" class="h-16 w-16 bg-red-50 dark:bg-red-950/40 rounded-lg flex items-center justify-center border border-red-200 dark:border-red-900/60 shadow-inner group-hover:scale-105 transition-all">
            <svg class="h-10 w-10 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 5h10a2 2 0 012 2v10a2 2 0 01-2 2H7a2 2 0 01-2-2V7a2 2 0 012-2z" />
            </svg>
          </div>
          <div class="h-16 w-16 bg-slate-50 dark:bg-slate-850 rounded-lg flex items-center justify-center border border-slate-200 dark:border-slate-800 group-hover:scale-105 transition-all" v-else>
            <UIcon name="i-heroicons-rectangle-group" class="h-8 w-8 text-gray-400" />
          </div>

          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">{{ b.brand }}</h3>
            <p class="text-xs text-slate-400">{{ b.website }}</p>
          </div>
        </div>

        <div class="text-right">
          <span class="text-xs font-mono font-bold text-primary-500 bg-primary-50 dark:bg-primary-950 px-2 py-1 rounded">
            {{ b.models?.length || 0 }} Models Available
          </span>
        </div>
      </div>
    </div>

    <!-- VIEW 2: Models List inside Selected Brand Explorer -->
    <div v-else class="space-y-6">
      <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
        <div class="flex items-center gap-3">
          <UButton 
            size="xs" 
            color="gray" 
            variant="subtle" 
            icon="i-heroicons-arrow-left" 
            @click="selectedBrand = null"
          >
            Back to Brands
          </UButton>
          <span class="text-sm font-mono font-bold text-slate-900 dark:text-white">
            Browsing: {{ selectedBrand.brand }} Catalog
          </span>
        </div>
        <span class="text-xs text-slate-400 font-mono">{{ selectedBrand.website }}</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <UCard v-for="m in selectedBrand.models" :key="m.model">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white font-mono text-primary-500">{{ m.model }}</h4>
                  
                  <!-- Sync Verification Badges (Import Status & Newer version checks) -->
                  <UBadge v-if="m.is_imported && m.update_available" color="amber" variant="subtle" class="font-mono text-[9px] font-bold uppercase animate-pulse">
                    v{{ m.imported_revision }} out of date
                  </UBadge>
                  <UBadge v-else-if="m.is_imported" color="emerald" variant="subtle" class="font-mono text-[9px] font-bold uppercase">
                    Imported v{{ m.imported_revision }}
                  </UBadge>
                  <UBadge v-else color="gray" variant="subtle" class="font-mono text-[9px] font-bold uppercase">
                    Not Imported
                  </UBadge>
                </div>
                <p class="text-[10px] text-slate-400 leading-normal font-mono">{{ m.general_info }}</p>
              </div>
              <UBadge color="gray" variant="subtle" class="font-mono text-[9px] uppercase font-bold">{{ m.subtype }}</UBadge>
            </div>
          </template>

          <div class="space-y-4">
            <!-- Ports layout specs -->
            <div class="flex items-center gap-2">
              <UIcon name="i-heroicons-cpu-chip" class="h-4 w-4 text-slate-400" />
              <span class="text-xs font-mono font-semibold">Standard Ports Layout:</span>
            </div>
            <div class="flex flex-wrap gap-2 pl-6">
              <UBadge 
                v-for="(count, type) in m.ports" 
                :key="type" 
                color="blue" 
                variant="subtle" 
                class="font-mono text-xs font-bold"
              >
                {{ count }}x {{ type }}
              </UBadge>
              <span v-if="!m.ports || Object.keys(m.ports).length === 0" class="text-xs text-slate-400 italic">No static ports</span>
            </div>

            <div class="flex justify-end pt-3 border-t border-slate-100 dark:border-slate-800">
              <!-- Case 1: Out of Date (Requires Update, will bump GORM Revision!) -->
              <UButton 
                v-if="m.is_imported && m.update_available"
                size="xs" 
                color="amber" 
                icon="i-heroicons-arrow-path" 
                :loading="isImporting[m.model]"
                @click="updateModelSpec(m)"
              >
                Update Spec (v{{ m.imported_revision + 1 }} Ready)
              </UButton>
              
              <!-- Case 2: Imported and up to date -->
              <UButton 
                v-else-if="m.is_imported"
                size="xs" 
                color="emerald" 
                variant="subtle"
                icon="i-heroicons-check-circle" 
                disabled
              >
                Catalog Up to Date
              </UButton>

              <!-- Case 3: Needs initial import -->
              <UButton 
                v-else
                size="xs" 
                color="primary" 
                icon="i-heroicons-arrow-down-tray" 
                :loading="isImporting[m.model]"
                @click="importModelSpec(m)"
              >
                Import Model Spec
              </UButton>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const route = useRoute()
const router = useRouter()
const { getAuthHeader } = useAuth()

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

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

// State parameters for catalog explorer
const isCatalogLoading = ref(true)
const diskCatalog = ref([])
const isImporting = ref({})



const selectedBrand = computed({
  get: () => {
    const brandQuery = route.query.brand
    if (!brandQuery || diskCatalog.value.length === 0) return null
    return diskCatalog.value.find(b => b.brand.toLowerCase().replace(/[^a-z0-9]/g, '') === brandQuery.toLowerCase()) || null
  },
  set: (val) => {
    const brandQueryVal = val ? val.brand.toLowerCase().replace(/[^a-z0-9]/g, '') : undefined
    router.replace({
      path: route.path,
      query: { ...route.query, brand: brandQueryVal }
    })
  }
})

const loadCatalog = async () => {
  isCatalogLoading.value = true
  try {
    const res = await $fetch(`${apiBase}/devices/catalog`)
    diskCatalog.value = res || []
  } catch (err) {
    console.error('Failed to load disk catalog:', err)
    notify('Failed to load hardware catalog data from disk.', true)
  } finally {
    isCatalogLoading.value = false
  }
}

onMounted(() => {
  loadCatalog()
})

const importModelSpec = async (m) => {
  isImporting.value[m.model] = true
  try {
    // Look up or auto-register manufacturer name inside GORM first (Robust Alphanumeric Case-Insensitive Match)
    let manufacturers = await $fetch(`${apiBase}/manufacturers`, {
      headers: getAuthHeader()
    })
    let manufacturer = manufacturers.find(brand => {
      const bName = brand?.name || ''
      const selName = selectedBrand.value?.brand || ''
      return bName.toLowerCase().replace(/[^a-z0-9]/g, '') === selName.toLowerCase().replace(/[^a-z0-9]/g, '')
    })
    
    if (!manufacturer) {
      manufacturer = await $fetch(`${apiBase}/manufacturers`, {
        method: 'POST',
        body: { name: selectedBrand.value.brand, website: selectedBrand.value.website },
        headers: getAuthHeader()
      })
    }

    // Save model spec to catalog securely (Set is_imported to true)
    await $fetch(`${apiBase}/devices`, {
      method: 'POST',
      body: {
        manufacturer_id: manufacturer.id,
        model_name: m.model,
        subtype: m.subtype,
        general_info: m.general_info,
        ports: m.ports,
        is_imported: true
      },
      headers: getAuthHeader()
    })

    notify(`${m.model} successfully imported into your hardware catalog!`)
    await loadCatalog() // Refresh to update imported states and badges!
  } catch (err) {
    console.error('Failed to import model spec:', err)
    notify(`Failed to import specification: Check if model already exists.`, true)
  } finally {
    isImporting.value[m.model] = false
  }
}

const updateModelSpec = async (m) => {
  isImporting.value[m.model] = true
  try {
    // Perform GORM PUT update which revision-pins and revision-increments used models!
    await $fetch(`${apiBase}/devices/${m.id}`, {
      method: 'PUT',
      body: {
        subtype: m.subtype,
        general_info: m.general_info,
        ports: m.ports
      },
      headers: getAuthHeader()
    })

    notify(`${m.model} successfully updated in your GORM catalog! Bumped to newer revision.`)
    await loadCatalog() // Refresh to update imported states and badges!
  } catch (err) {
    console.error('Failed to update model spec:', err)
    notify(`Failed to update specification.`, true)
  } finally {
    isImporting.value[m.model] = false
  }
}

const backToDevices = () => {
  navigateTo('/devices')
}
</script>
