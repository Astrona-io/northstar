<template>
  <UModal v-model="isOpen">
    <UCard>
      <template #header>
        <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Add Network Interface (NIC / Port)
        </h3>
      </template>
      
      <form @submit.prevent="saveInterface" class="space-y-4">
        <UFormGroup label="Interface Name">
          <UInput v-model="form.name" placeholder="e.g. eth0, GigabitEthernet0/1" required />
        </UFormGroup>
        
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="Interface Type">
            <USelect v-model="form.type" :options="['ethernet', 'management', 'console', 'sfp+', 'fiber']" required />
          </UFormGroup>
          <UFormGroup label="Link Speed">
            <UInput v-model="form.speed" placeholder="e.g. 1 Gbps, 10 Gbps" />
          </UFormGroup>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="IP Address">
            <UInput v-model="form.ip_address" placeholder="e.g. 192.168.1.15" />
          </UFormGroup>
          <UFormGroup label="MAC Address">
            <UInput v-model="form.mac_address" placeholder="e.g. 00:1A:2B:3C:4D:5E" />
          </UFormGroup>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="VLAN">
            <UInput v-model="form.vlan" placeholder="e.g. VLAN 100" />
          </UFormGroup>
          <UFormGroup label="MTU (Bytes)">
            <UInput type="number" v-model="form.mtu" placeholder="1500" />
          </UFormGroup>
        </div>

        <UFormGroup label="Initial Link Status">
          <USelect v-model="form.status" :options="['up', 'down']" required />
        </UFormGroup>
        
        <div class="flex justify-end gap-2 mt-4">
          <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="isSaving">Save</UButton>
        </div>
      </form>
    </UCard>
  </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'

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

const { createInterface } = useInterfaces()

const isSaving = ref(false)

const form = ref({
  name: '',
  type: 'ethernet',
  speed: '1 Gbps',
  ip_address: '',
  mac_address: '',
  vlan: '',
  mtu: 1500,
  status: 'up'
})

const saveInterface = async () => {
  isSaving.value = true
  try {
    const payload = {
      ...form.value,
      mtu: Number(form.value.mtu) || 1500
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
      status: 'up'
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
