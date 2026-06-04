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

        <!-- Multi-Category Selection (Phase 2 Catalog Update - UniFi Switch/Router hybrid support) -->
        <UFormGroup label="Device Categories (Multi-select, defaults to 'global')">
          <USelectMenu 
            v-model="selectedCategoryIDs" 
            :options="categories" 
            multiple 
            option-attribute="name" 
            value-attribute="id" 
            placeholder="Select hardware category tags..." 
            searchable 
          />
        </UFormGroup>

        <UFormGroup label="General Information / Hard Specifications">
          <UTextarea v-model="form.general_info" placeholder="Enter device capabilities, default interfaces, or manufacturer details..." />
        </UFormGroup>
        
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

const form = ref({
  manufacturer_id: null,
  model_name: '',
  general_info: ''
})

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
        general_info: props.device.general_info
      }
      selectedCategoryIDs.value = props.device.categories ? props.device.categories.map(c => c.id) : []
    } else {
      form.value = { manufacturer_id: null, model_name: '', general_info: '' }
      selectedCategoryIDs.value = []
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

    const payload = {
      ...form.value,
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
