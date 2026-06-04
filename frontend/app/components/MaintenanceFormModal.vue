<template>
  <UModal v-model="isOpen">
    <UCard>
      <template #header>
        <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
          Schedule Maintenance Window
        </h3>
      </template>
      
      <form @submit.prevent="saveMaintenanceWindow" class="space-y-4">
        <UFormGroup label="Title">
          <UInput v-model="maintenanceForm.title" placeholder="Quarterly Patching" required />
        </UFormGroup>
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="Start Time">
            <UInput type="datetime-local" v-model="maintenanceForm.start_time" required />
          </UFormGroup>
          <UFormGroup label="End Time">
            <UInput type="datetime-local" v-model="maintenanceForm.end_time" required />
          </UFormGroup>
        </div>
        <UFormGroup label="Description">
          <UTextarea v-model="maintenanceForm.description" />
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

const { createMaintenanceWindow } = useMaintenance()

const isSaving = ref(false)

const maintenanceForm = ref({
  title: '',
  description: '',
  start_time: '',
  end_time: ''
})

const saveMaintenanceWindow = async () => {
  isSaving.value = true
  try {
    const payload = {
      ...maintenanceForm.value,
      start_time: new Date(maintenanceForm.value.start_time).toISOString(),
      end_time: new Date(maintenanceForm.value.end_time).toISOString()
    }
    await createMaintenanceWindow(props.assetId, payload)
    isOpen.value = false
    maintenanceForm.value = { title: '', description: '', start_time: '', end_time: '' }
    emit('saved')
  } catch (error) {
    console.error('Failed to save maintenance window:', error)
    alert('Failed to save maintenance window')
  } finally {
    isSaving.value = false
  }
}
</script>