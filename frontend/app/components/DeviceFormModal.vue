<template>
  <UModal v-model="isOpen">
    <UCard>
      <template #header>
        <div class="flex items-center gap-3">
          <UIcon name="i-heroicons-cpu-chip" class="h-6 w-6 text-primary-500" />
          <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            {{ isEditing ? 'Edit Device Specification' : 'Add New Hardware Model' }}
          </h3>
        </div>
      </template>
      
      <form @submit.prevent="saveDeviceModel" class="space-y-4">
        <!-- Manufacturer selector with inline Add New capability (Phase 2 Catalog Update) -->
        <UFormGroup label="Manufacturer / Brand">
          <div class="flex gap-2">
            <USelectMenu 
              v-if="!isAddingBrand" 
              v-model="form.manufacturer_id" 
              :options="manufacturers" 
              option-attribute="name" 
              value-attribute="id" 
              placeholder="Select existing manufacturer brand..." 
              required 
              searchable 
              class="flex-1" 
            />
            <UInput 
              v-else 
              v-model="newBrandName" 
              placeholder="Type new brand name..." 
              required 
              class="flex-1" 
            />
            <UButton 
              :icon="isAddingBrand ? 'i-heroicons-x-mark' : 'i-heroicons-plus'" 
              color="gray" 
              :title="isAddingBrand ? 'Select existing' : 'Register new brand'"
              @click="toggleAddBrand" 
            />
          </div>
        </UFormGroup>
        
        <UFormGroup label="Model Name / Specification Code">
          <UInput v-model="form.model_name" placeholder="e.g. Catalyst 9300-48UX, PowerEdge R740" required />
        </UFormGroup>

        <!-- Equipment Classification (One Device can only be One thing!) -->
        <UFormGroup label="Equipment Classification" class="text-xs font-semibold">
          <USelect 
            v-model="form.subtype" 
            :options="[
              'Server',
              'Database',
              'Router',
              'Switch (L2)',
              'Switch (L3)',
              'Access Point (AP)',
              'Firewall',
              'Load Balancer',
              'Industrial PLC',
              'Edge Gateway',
              'Smart IP Camera',
              'Environment Sensor',
              'Other'
            ]" 
            required 
          />
        </UFormGroup>

        <UFormGroup label="General Information / Hard Specifications">
          <UTextarea v-model="form.general_info" placeholder="Enter device capabilities, default interfaces, or manufacturer details..." />
        </UFormGroup>

        <!-- Interface Ports Builder (Standard physical port mapping) -->
        <div class="space-y-3 border-t border-gray-100 dark:border-gray-800 pt-3">
          <div class="flex justify-between items-center">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono flex items-center gap-1.5">
              <UIcon name="i-heroicons-bolt" class="text-primary-500 h-4 w-4" />
              Standard Port Configuration
            </h4>
            <UButton size="xs" color="primary" variant="soft" icon="i-heroicons-plus" @click="addPortConfigRow">
              Add Port
            </UButton>
          </div>
          
          <div v-if="portsList.length === 0" class="text-xs text-gray-400 italic">
            No physical ports defined. Standard switches and routers should specify port profiles.
          </div>
          <div v-else class="space-y-2 max-h-40 overflow-y-auto pr-1">
            <div v-for="(row, idx) in portsList" :key="idx" class="flex gap-2 items-center">
              <USelect 
                v-model="row.type" 
                :options="['RJ45', 'SFP', 'SFP+', 'QSFP+', 'GBIC', 'Fiber (LC)', 'Fiber (SC)', 'Coaxial', 'Serial']" 
                size="xs"
                class="flex-1"
                placeholder="Port Type..."
              />
              <UInput 
                type="number" 
                v-model="row.count" 
                size="xs" 
                placeholder="Count" 
                min="1"
                class="w-24"
              />
              <UButton 
                size="xs" 
                color="red" 
                variant="ghost" 
                icon="i-heroicons-trash" 
                @click="removePortConfigRow(idx)" 
              />
            </div>
          </div>
        </div>
        
        <div class="flex justify-end gap-2 mt-6 pt-3 border-t border-gray-100 dark:border-gray-800">
          <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="isSaving">Save Spec</UButton>
        </div>
      </form>
    </UCard>
  </UModal>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  device: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isEditing = computed(() => !!props.device)

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase
const { createDevice, updateDevice } = useDevices()
const { getAuthHeader } = useAuth()

const isSaving = ref(false)
const isAddingBrand = ref(false)
const newBrandName = ref('')

const manufacturers = ref([])
const categories = ref([])
const selectedCategoryIDs = ref([])
const portsList = ref([])

const form = ref({
  manufacturer_id: null,
  model_name: '',
  general_info: '',
  subtype: 'Server'
})

const addPortConfigRow = () => {
  portsList.value.push({ type: 'RJ45', count: 24 })
}

const removePortConfigRow = (idx) => {
  portsList.value.splice(idx, 1)
}

const toggleAddBrand = () => {
  isAddingBrand.value = !isAddingBrand.value
  newBrandName.value = ''
  form.value.manufacturer_id = null
}

const loadFormMetadata = async () => {
  try {
    const brands = await $fetch(`${apiBase}/manufacturers`)
    manufacturers.value = brands || []
    
    const cats = await $fetch(`${apiBase}/categories`)
    categories.value = cats || []
  } catch (err) {
    console.error('Failed to load catalog metadata:', err)
  }
}

watch(() => props.modelValue, (isOpenVal) => {
  if (isOpenVal) {
    loadFormMetadata()
    isAddingBrand.value = false
    newBrandName.value = ''
    if (props.device) {
      form.value = {
        manufacturer_id: props.device.manufacturer_id,
        model_name: props.device.model_name,
        general_info: props.device.general_info,
        subtype: props.device.subtype || 'Server'
      }
      selectedCategoryIDs.value = props.device.categories ? props.device.categories.map(c => c.id) : []
      if (props.device.ports) {
        portsList.value = Object.entries(props.device.ports).map(([type, count]) => ({
          type,
          count: Number(count)
        }))
      } else {
        portsList.value = []
      }
    } else {
      form.value = { manufacturer_id: null, model_name: '', general_info: '', subtype: 'Server' }
      selectedCategoryIDs.value = []
      portsList.value = []
    }
  }
})

const saveDeviceModel = async () => {
  isSaving.value = true
  try {
    // 1. Strict Brand Registration Rule (Phase 2 Catalog)
    // If user is registering a brand that doesn't exist, save it first!
    if (isAddingBrand.value && newBrandName.value) {
      try {
        const newBrand = await $fetch(`${apiBase}/manufacturers`, {
          method: 'POST',
          body: { name: newBrandName.value },
          headers: getAuthHeader()
        })
        form.value.manufacturer_id = newBrand.id
      } catch (err) {
        console.error('Failed to register new brand:', err)
        alert('Failed to register brand: Check permissions or name uniqueness.')
        isSaving.value = false
        return
      }
    }

    if (!form.value.manufacturer_id) {
      alert('Please select or register a Manufacturer brand first.')
      isSaving.value = false
      return
    }

    const portsMap = {}
    portsList.value.forEach(row => {
      if (row.type && row.count) {
        portsMap[row.type] = Number(row.count)
      }
    })

    const payload = {
      ...form.value,
      ports: portsMap,
      category_ids: selectedCategoryIDs.value
    }

    if (isEditing.value) {
      await updateDevice(props.device.id, payload)
    } else {
      await createDevice(payload)
    }
    isOpen.value = false
    emit('saved')
  } catch (error) {
    console.error('Failed to save device profile:', error)
    alert('Failed to save device profile')
  } finally {
    isSaving.value = false
  }
}
</script>
