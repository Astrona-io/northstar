<template>
  <UModal v-model="isOpen">
    <UCard>
      <template #header>
        <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Add Network Interface (NIC / Port)
        </h3>
      </template>
      
      <form @submit.prevent="saveInterface" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="Interface Name">
            <UInput v-model="form.name" placeholder="e.g. NIC 1, Intel Quad-Port" required />
          </UFormGroup>
          <UFormGroup label="Interface Type">
            <USelectMenu 
              v-model="selectedProfile" 
              :options="profiles" 
              option-attribute="type" 
              placeholder="Select port type..." 
              required
              class="w-full"
            />
          </UFormGroup>
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="Link Speed">
            <USelect v-model="form.speed" :options="speedOptions" :disabled="!selectedProfile" placeholder="Select speed..." required />
          </UFormGroup>
          <UFormGroup label="Port Count (Interfaces on this NIC)">
            <UInput type="number" v-model="form.port_count" placeholder="1" min="1" required />
          </UFormGroup>
        </div>
        
        <div class="flex justify-end gap-2 mt-4">
          <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="isSaving">Save</UButton>
        </div>
      </form>
    </UCard>
  </UModal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  assetId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

const { createInterface } = useInterfaces()

const isSaving = ref(false)
const profiles = ref([])
const selectedProfileId = ref(null)

const loadInterfaceProfiles = async () => {
  try {
    const data = await $fetch(`${apiBase}/port-types`)
    profiles.value = data || []
  } catch (err) {
    console.error('Failed to load standard interface profiles:', err)
  }
}

watch(() => props.modelValue, (isOpenVal) => {
  if (isOpenVal) {
    loadInterfaceProfiles()
    selectedProfileId.value = null
  }
})

const speedOptions = computed(() => {
  if (!selectedProfile.value) return []
  try {
    return JSON.parse(selectedProfile.value.speeds) || []
  } catch (err) {
    return []
  }
})

const selectedProfile = computed({
  get: () => {
    if (!selectedProfileId.value || profiles.value.length === 0) return null
    return profiles.value.find(p => p.id === selectedProfileId.value) || null
  },
  set: (val) => {
    if (val) {
      selectedProfileId.value = val.id
      form.value.type = val.type
      try {
        const speeds = JSON.parse(val.speeds)
        if (speeds && speeds.length > 0) {
          form.value.speed = speeds[0]
        } else {
          form.value.speed = ''
        }
      } catch (err) {
        form.value.speed = ''
      }
    } else {
      selectedProfileId.value = null
      form.value.type = 'ethernet'
      form.value.speed = '1 Gbps'
    }
  }
})

const form = ref({
  name: '',
  type: 'ethernet',
  speed: '1 Gbps',
  ip_address: '',
  mac_address: '',
  vlan: '',
  mtu: 1500,
  status: 'up',
  port_count: 1
})

const saveInterface = async () => {
  isSaving.value = true
  try {
    const payload = {
      ...form.value,
      nic_name: form.value.name,
      mtu: Number(form.value.mtu) || 1500,
      port_count: Number(form.value.port_count) || 1
    }
    await createInterface(props.assetId, payload)
    isOpen.value = false
    form.value = {
      name: '',
      type: 'ethernet',
      speed: '1 Gbps',
      ip_address: '',
      mac_address: '',
      vlan: '',
      mtu: 1500,
      status: 'up',
      port_count: 1
    }
    emit('saved')
  } catch (error) {
    console.error('Failed to save interface:', error)
    alert('Failed to save interface')
  } finally {
    isSaving.value = false
  }
}
</script>
