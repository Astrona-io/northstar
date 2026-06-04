<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white">Active Network Auto-Discovery</h2>
    </div>

    <!-- Scanner Config Card -->
    <UCard>
      <template #header>
        <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Subnet Scan Setup</h3>
      </template>
      <div class="flex flex-col sm:flex-row gap-4 items-end">
        <UFormGroup label="Target Subnet CIDR Range" class="flex-1">
          <UInput v-model="subnet" icon="i-heroicons-globe-alt" placeholder="e.g. 192.168.1.0/24" />
        </UFormGroup>
        <UButton :loading="isScanning" icon="i-heroicons-bolt" color="primary" size="md" @click="runDiscoveryScan">
          Trigger SNMP Sweep
        </UButton>
      </div>
    </UCard>

    <!-- Loading Overlay -->
    <UCard v-if="isScanning">
      <div class="flex flex-col items-center justify-center p-12 space-y-4">
        <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Scanning network range...</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400">Pinging subnets and mapping hardware SNMP profiles on {{ subnet }}...</p>
      </div>
    </UCard>

    <!-- Discovered Devices Results Table -->
    <div v-else-if="discoveredDevices.length > 0" class="space-y-4">
      <h3 class="text-lg font-medium text-gray-900 dark:text-white flex items-center gap-2">
        <UIcon name="i-heroicons-magnifying-glass-circle" class="h-6 w-6 text-primary-500" />
        Discovered Equipment ({{ discoveredDevices.length }})
      </h3>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <UCard v-for="(dev, idx) in discoveredDevices" :key="idx" class="flex flex-col">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon :name="getDeviceIcon(dev.type)" class="h-6 w-6 text-primary-500" />
                <h4 class="font-bold text-gray-900 dark:text-white">{{ dev.name }}</h4>
              </div>
              <UBadge color="gray" variant="subtle" size="xs">{{ dev.type }}</UBadge>
            </div>
          </template>
          
          <div class="space-y-2 text-sm text-gray-600 dark:text-gray-300 flex-1">
            <p>{{ dev.description }}</p>
            <div class="grid grid-cols-2 gap-x-2 gap-y-1.5 pt-2 text-xs border-t border-gray-100 dark:border-gray-800">
              <span class="text-gray-500 font-medium">IP Address</span>
              <span class="font-mono">{{ dev.ip_address }}</span>
              <span class="text-gray-500 font-medium">MAC Address</span>
              <span class="font-mono text-[10px]">{{ dev.mac_address }}</span>
              <span class="text-gray-500 font-medium">Hardware</span>
              <span>{{ dev.manufacturer }} {{ dev.model }}</span>
            </div>
          </div>
          
          <div class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-800 flex justify-end">
            <UButton 
              v-if="canMutate" 
              size="xs" 
              icon="i-heroicons-plus-circle" 
              color="primary" 
              @click="importToCMDB(dev, idx)"
            >
              Import to CMDB
            </UButton>
            <span v-else class="text-xs text-gray-400 italic">Login as Admin/Operator to Import</span>
          </div>
        </UCard>
      </div>
    </div>
    
    <div v-else-if="hasScanned && discoveredDevices.length === 0" class="text-center py-12 text-gray-500 dark:text-gray-400">
      <UIcon name="i-heroicons-check-circle" class="h-12 w-12 text-green-500 mb-2" />
      <p>No new or unmanaged assets discovered on subnet sweep. All items are in sync!</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const subnet = ref('192.168.1.0/24')
const isScanning = ref(false)
const hasScanned = ref(false)
const discoveredDevices = ref([])

const { canMutate } = useAuth()
const { scanSubnet, importDevice } = useDiscovery()

const runDiscoveryScan = async () => {
  isScanning.value = true
  hasScanned.value = false
  discoveredDevices.value = []
  
  // Launch scan endpoint
  try {
    const { data } = await scanSubnet(subnet.value)
    // Wait for the simulated delay
    await new Promise(resolve => setTimeout(resolve, 1200))
    if (data.value) {
      discoveredDevices.value = data.value
    }
  } catch (error) {
    console.error('Scan failed:', error)
  } finally {
    isScanning.value = false
    hasScanned.value = true
  }
}

const importToCMDB = async (deviceSpec, idx) => {
  try {
    await importDevice(deviceSpec)
    // Remove from active list
    discoveredDevices.value.splice(idx, 1)
    alert(`Successfully imported ${deviceSpec.name} to active CMDB Assets. Audit trace logged.`)
  } catch (error) {
    console.error('Failed to import discovered device:', error)
    alert('Import failed: Check permissions.')
  }
}

const getDeviceIcon = (type) => {
  switch (type) {
    case 'Server':
      return 'i-heroicons-server'
    case 'Database':
      return 'i-heroicons-circle-stack'
    case 'Router':
      return 'i-heroicons-cpu-chip'
    case 'Switch':
      return 'i-heroicons-arrows-right-left'
    default:
      return 'i-heroicons-squares-2x2'
  }
}
</script>
