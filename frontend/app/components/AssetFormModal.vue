<template>
  <UModal v-model="isOpen" :ui="{ width: 'sm:max-w-5xl' }">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <UIcon :name="step === 1 ? 'i-heroicons-squares-plus' : selectedTypeIcon" class="h-6 w-6 text-primary-500" />
            <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
              {{ isEditing ? 'Edit Asset' : (step === 1 ? 'Select Asset Type' : 'Configure Asset Details') }}
            </h3>
          </div>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isOpen = false" />
        </div>
      </template>

      <!-- Step 1: Category & Sub-Group Selection (Phase 1 Sub-Group Portals) -->
      <div v-if="step === 1" class="space-y-4 py-2">
        <!-- Guided Onboarding Catalog Check -->
        <div v-if="deviceCatalogEmpty && !isEditing" class="p-6 bg-slate-50 dark:bg-slate-800/30 rounded-md border border-dashed border-slate-200 dark:border-slate-800 text-center space-y-4 my-2">
          <UIcon name="i-heroicons-shield-exclamation" class="h-12 w-12 text-primary-500 mx-auto animate-pulse" />
          <h4 class="font-bold text-base text-gray-900 dark:text-white font-mono">Device Catalog Registration Required</h4>
          <p class="text-sm text-gray-500 dark:text-gray-400 max-w-lg mx-auto leading-normal">
            Welcome to Northstar! Before you can register physical or virtual network assets, you must first define standardized hardware specifications inside the Device Catalog so we can standardize your equipment.
          </p>
          <div class="flex justify-center gap-2 pt-2">
            <UButton to="/devices" color="primary" icon="i-heroicons-cpu-chip" @click="isOpen = false">
              Configure Device Catalog
            </UButton>
            <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
          </div>
        </div>

        <!-- Normal Selection Grids -->
        <template v-else>
          <!-- View 1.1: Root Categories Selection -->
          <div v-if="subGroupActive === ''" class="space-y-4">
            <p class="text-xs text-gray-500 dark:text-gray-400">Choose the hardware category or cloud service deploy structure to standard register:</p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 pb-2">
              <button 
                v-for="t in assetTypes" 
                :key="t.name"
                type="button"
                class="flex flex-col items-center p-4 rounded-md border border-slate-200 dark:border-slate-800 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-all text-center group min-h-[145px]"
                @click="onCategoryClick(t.name)"
              >
                <UIcon :name="t.icon" class="h-7 w-7 text-gray-400 group-hover:text-primary-500 mb-2 transition-colors" />
                <span class="text-xs font-semibold text-gray-900 dark:text-white group-hover:text-primary-500 transition-colors leading-tight">{{ t.name }}</span>
                <span class="text-[10px] text-gray-400 dark:text-gray-500 mt-1 line-clamp-3 leading-normal">{{ t.desc }}</span>
              </button>
            </div>
            <div class="flex justify-end pt-2 border-t border-slate-200 dark:border-slate-800">
              <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
            </div>
          </div>

          <!-- View 1.2: Nested Network Sub-Groups selection -->
          <div v-else-if="subGroupActive === 'Network'" class="space-y-4">
            <div class="flex justify-between items-center bg-blue-50/10 dark:bg-blue-950/10 p-3 rounded-md border border-blue-500/20">
              <span class="text-xs text-slate-600 dark:text-slate-300 font-mono">Select specific Network subtype to proceed:</span>
              <UButton size="xs" color="gray" variant="ghost" label="Back to Categories" icon="i-heroicons-arrow-left" @click="subGroupActive = ''" />
            </div>
            
            <div class="grid grid-cols-2 md:grid-cols-3 gap-4 pb-2">
              <button 
                v-for="sub in networkSubGroups" 
                :key="sub.name"
                type="button"
                class="flex flex-col items-center p-4 rounded-md border border-slate-200 dark:border-slate-800 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-all text-center group"
                @click="selectNetworkSubGroup(sub.name)"
              >
                <UIcon :name="sub.icon" class="h-6 w-6 text-gray-400 group-hover:text-primary-500 mb-2 transition-colors" />
                <span class="text-xs font-semibold text-gray-900 dark:text-white group-hover:text-primary-500 transition-colors">{{ sub.name }}</span>
                <span class="text-[9px] text-slate-400 dark:text-slate-500 mt-1 line-clamp-2 leading-tight">{{ sub.desc }}</span>
              </button>
            </div>
            
            <div class="flex justify-end pt-2 border-t border-slate-200 dark:border-slate-800">
              <UButton color="gray" variant="ghost" @click="subGroupActive = ''">Back</UButton>
            </div>
          </div>

          <!-- View 1.3: Nested IoT / Edge Sub-Groups selection -->
          <div v-else-if="subGroupActive === 'IoT / Edge'" class="space-y-4">
            <div class="flex justify-between items-center bg-teal-50/10 dark:bg-teal-950/10 p-3 rounded-md border border-teal-500/20">
              <span class="text-xs text-slate-600 dark:text-slate-300 font-mono">Select specific IoT Edge subtype to proceed:</span>
              <UButton size="xs" color="gray" variant="ghost" label="Back to Categories" icon="i-heroicons-arrow-left" @click="subGroupActive = ''" />
            </div>
            
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 pb-2">
              <button 
                v-for="sub in iotSubGroups" 
                :key="sub.name"
                type="button"
                class="flex flex-col items-center p-4 rounded-md border border-slate-200 dark:border-slate-800 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-all text-center group"
                @click="selectIotSubGroup(sub.name)"
              >
                <UIcon :name="sub.icon" class="h-6 w-6 text-gray-400 group-hover:text-primary-500 mb-2 transition-colors" />
                <span class="text-xs font-semibold text-gray-900 dark:text-white group-hover:text-primary-500 transition-colors">{{ sub.name }}</span>
                <span class="text-[9px] text-slate-400 dark:text-slate-500 mt-1 line-clamp-2 leading-tight">{{ sub.desc }}</span>
              </button>
            </div>
            
            <div class="flex justify-end pt-2 border-t border-slate-200 dark:border-slate-800">
              <UButton color="gray" variant="ghost" @click="subGroupActive = ''">Back</UButton>
            </div>
          </div>
        </template>
      </div>

      <!-- Step 2: Form Details -->
      <form v-else @submit.prevent="saveAsset" class="space-y-4 h-[420px] overflow-y-auto pr-1">
        <!-- Category Summary & Back Button -->
        <div class="flex items-center justify-between bg-slate-50 dark:bg-slate-800/30 p-3 rounded-md border border-slate-200 dark:border-slate-800/50 mb-2">
          <div class="flex items-center gap-2">
            <UIcon :name="selectedTypeIcon" class="h-5 w-5 text-primary-500" />
            <span class="text-sm font-medium text-slate-700 dark:text-slate-300">
              Category: <strong class="text-slate-900 dark:text-white">{{ form.type }}</strong>
              <span v-if="form.properties.network_subtype" class="text-slate-400"> &rarr; {{ form.properties.network_subtype }}</span>
              <span v-else-if="form.properties.iot_subtype" class="text-slate-400"> &rarr; {{ form.properties.iot_subtype }}</span>
            </span>
          </div>
          <UButton v-if="!isEditing" size="xs" variant="ghost" color="gray" icon="i-heroicons-arrow-left" @click="step = 1" />
        </div>

        <!-- Taxonomy Prerequisite Alert Banner (Phase 1 Taxonomy Chains) -->
        <div v-if="activeParentDependency" class="p-3 bg-amber-50/10 border border-amber-500/30 rounded-md text-[11px] text-amber-500/90 leading-relaxed font-mono flex items-start gap-2 mb-2">
          <UIcon name="i-heroicons-shield-exclamation" class="h-4 w-4 shrink-0 text-amber-500" />
          <div>
            <strong>Taxonomy Rule:</strong> This category strictly requires selecting a parent host/worker asset of type: <code class="bg-amber-950 px-1.5 py-0.5 rounded text-white">{{ activeParentDependency.toUpperCase() }}</code>
          </div>
        </div>

        <UFormGroup label="Name">
          <UInput v-model="form.name" placeholder="Server-01" required />
        </UFormGroup>
        
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup label="Status" class="col-span-2">
            <USelect v-model="form.status" :options="['active', 'inactive', 'maintenance']" required />
          </UFormGroup>
        </div>

        <UFormGroup label="Description">
          <UTextarea v-model="form.description" />
        </UFormGroup>

        <!-- ServerSpecs -->
        <div v-if="form.type === 'Server'" class="space-y-4 pt-2">
          <UDivider label="Server Properties" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="CPU Model">
              <UInput v-model="form.properties.cpu_model" placeholder="e.g. Intel Xeon" />
            </UFormGroup>
            <UFormGroup label="RAM (GB)">
              <UInput type="number" v-model="form.properties.ram_capacity_gb" placeholder="256" />
            </UFormGroup>
            <UFormGroup label="OS Name">
              <UInput v-model="form.properties.os_name" placeholder="e.g. Ubuntu Linux" />
            </UFormGroup>
            <UFormGroup label="Storage Subsystem">
              <UInput v-model="form.properties.storage_subsystem" placeholder="e.g. 4x 1.2TB NVMe" />
            </UFormGroup>
          </div>
        </div>

        <!-- Database Properties -->
        <div v-else-if="form.type === 'Database'" class="space-y-4 pt-2">
          <UDivider label="Database Properties" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Engine">
              <USelect v-model="form.properties.engine" :options="['PostgreSQL', 'MySQL', 'Oracle', 'SQL Server', 'MongoDB']" placeholder="Select..." />
            </UFormGroup>
            <UFormGroup label="Version">
              <UInput v-model="form.properties.version" placeholder="e.g. 14.2" />
            </UFormGroup>
            <UFormGroup label="Environment">
              <USelect v-model="form.properties.environment" :options="['Production', 'Staging', 'QA', 'Development']" placeholder="Select..." />
            </UFormGroup>
            <UFormGroup label="Allocated Storage (GB)">
              <UInput type="number" v-model="form.properties.allocated_storage_gb" placeholder="500" />
            </UFormGroup>
          </div>
        </div>

        <!-- Network Properties (Router/Switches Sub-groups) -->
        <div v-else-if="form.type === 'Network'" class="space-y-4 pt-2">
          <UDivider label="Network Device Specifications" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Standard Catalog Device Model" class="col-span-2">
              <USelectMenu 
                v-model="selectedDeviceModel" 
                :options="devices" 
                option-attribute="displayName" 
                placeholder="Select standardized hardware catalog model spec (autofills below)..." 
                searchable
                class="w-full"
              />
            </UFormGroup>

            <!-- Catalog revision details and manual upgrade -->
            <div v-if="selectedDeviceModel" class="col-span-2 bg-slate-50 dark:bg-slate-800/40 p-3 rounded border border-slate-150 dark:border-slate-800 flex justify-between items-center text-xs">
              <div>
                <span class="font-bold text-slate-800 dark:text-white block font-mono">Catalog Revision Locking</span>
                <span class="text-slate-500 mt-0.5 block font-mono">
                  Asset Pinned Revision: <strong>v{{ form.device_model_revision || 1 }}</strong> 
                  | Catalog Current Revision: <strong>v{{ selectedDeviceModel.revision || 1 }}</strong>
                </span>
              </div>
              <UButton 
                v-if="selectedDeviceModel.revision > (form.device_model_revision || 1)" 
                size="xs" 
                color="blue" 
                icon="i-heroicons-arrow-up-circle"
                @click="upgradeDeviceModelRevision"
              >
                Upgrade to v{{ selectedDeviceModel.revision }}
              </UButton>
              <span v-else class="text-[10px] text-green-600 dark:text-green-400 font-bold font-mono uppercase tracking-wider bg-green-50 dark:bg-green-950/30 px-2 py-0.5 rounded border border-green-200 dark:border-green-900/40">Latest Revision Pinned</span>
            </div>

            <UFormGroup label="Network Device Subtype">
              <USelect v-model="form.properties.network_subtype" :options="['Router', 'Switch (L2)', 'Switch (L3)', 'Access Point (AP)', 'Firewall', 'Load Balancer']" placeholder="Select subtype..." required />
            </UFormGroup>
            <UFormGroup label="Management IP Address">
              <UInput v-model="form.properties.management_ip" placeholder="e.g. 10.0.1.5" />
            </UFormGroup>
            <UFormGroup label="Manufacturer Brand">
              <UInput v-model="form.properties.manufacturer" placeholder="e.g. Cisco Systems, Ubiquiti" />
            </UFormGroup>
            <UFormGroup label="Hardware Model">
              <UInput v-model="form.properties.model" placeholder="e.g. Catalyst 9300, UniFi AP" />
            </UFormGroup>
            <UFormGroup label="Ports Config Descriptor" class="col-span-2">
              <UInput v-model="form.properties.port_config" placeholder="e.g. 48x 10/100/1000M RJ45 + 4x 10G SFP+" />
            </UFormGroup>
          </div>
          
          <!-- Topology Connections (Phase 8) -->
          <UDivider label="Topology Patch Cabling" />
          <div class="space-y-4">
            <UFormGroup label="Uplink Devices (Upstream Switch/Router)">
              <USelectMenu v-model="selectedUplinks" :options="allAssets" multiple option-attribute="name" value-attribute="id" placeholder="Select upstream devices..." searchable />
            </UFormGroup>
            <UFormGroup label="Downlink Devices (Downstream Switches/Servers)">
              <USelectMenu v-model="selectedDownlinks" :options="allAssets" multiple option-attribute="name" value-attribute="id" placeholder="Select downstream devices..." searchable />
            </UFormGroup>
          </div>
        </div>

        <!-- IoT Edge Properties -->
        <div v-else-if="form.type === 'IoT / Edge'" class="space-y-4 pt-2">
          <UDivider label="IoT Edge Device Specifications" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="IoT Device Sub-Stack">
              <USelect v-model="form.properties.iot_subtype" :options="['Edge Gateway (e.g. Pi/Jetson)', 'Industrial PLC / SCADA', 'Smart IP Camera', 'Environment IoT Sensor']" placeholder="Select IoT Sub-Stack..." required />
            </UFormGroup>
            <UFormGroup label="Active Firmware Version">
              <UInput v-model="form.properties.firmware_version" placeholder="e.g. v1.2.4-stable" />
            </UFormGroup>
            <UFormGroup label="Power Source Mode">
              <USelect v-model="form.properties.power_source" :options="['PoE (Power over Ethernet)', 'Battery Operated', 'External 12V DC', 'Wall Plug (110V/230V)']" placeholder="Select..." />
            </UFormGroup>
            <UFormGroup label="Operational Status Range">
              <UInput v-model="form.properties.op_range" placeholder="e.g. -10C to 60C, Indoor" />
            </UFormGroup>
            <!-- Conditional fields based on IoT subtype selection -->
            <UFormGroup v-if="form.properties.iot_subtype === 'Industrial PLC / SCADA'" label="Loop Voltage Calibration (V)" class="col-span-2">
              <UInput type="number" v-model="form.properties.loop_voltage" placeholder="24" />
            </UFormGroup>
            <UFormGroup v-else-if="form.properties.iot_subtype === 'Smart IP Camera'" label="Video Stream Encoding (Resolution)" class="col-span-2">
              <USelect v-model="form.properties.video_encoding" :options="['H.264 (1080p)', 'H.265 (4K UHD)', 'MJPEG (Substream)']" />
            </UFormGroup>
          </div>
        </div>

        <!-- Application Properties -->
        <div v-else-if="form.type === 'Application'" class="space-y-4 pt-2">
          <UDivider label="Application Properties" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Vendor / Publisher">
              <UInput v-model="form.properties.vendor" placeholder="e.g. Internal" />
            </UFormGroup>
            <UFormGroup label="Version / Build">
              <UInput v-model="form.properties.version" placeholder="e.g. v2.4.0" />
            </UFormGroup>
            <UFormGroup label="Framework / Language">
              <UInput v-model="form.properties.framework" placeholder="e.g. Java/Spring" />
            </UFormGroup>
            <UFormGroup label="URL / Endpoint">
              <UInput v-model="form.properties.url" placeholder="https://" />
            </UFormGroup>
          </div>
        </div>

        <!-- Kubernetes Deployment Form (Phase 14) -->
        <div v-else-if="form.type === 'Kubernetes Deployment'" class="space-y-4 pt-2">
          <UDivider label="Kubernetes Pod Specification" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Pod Container Image Tag">
              <UInput v-model="form.container_image" placeholder="e.g. postgres:15-alpine" required />
            </UFormGroup>
            <UFormGroup label="Service Port Map (IP Port)">
              <UInput v-model="form.container_port_mapping" placeholder="e.g. 5432:5432" />
            </UFormGroup>
          </div>
          <UFormGroup label="Hosting Kubernetes Worker Node (Server)">
            <USelectMenu v-model="form.host_asset_id" :options="serverAssets" option-attribute="name" value-attribute="id" placeholder="Select worker host server..." searchable />
          </UFormGroup>
        </div>

        <!-- Standalone Container Form (Docker/AWS Fargate, New Request) -->
        <div v-else-if="form.type === 'Container'" class="space-y-4 pt-2">
          <UDivider label="Standalone Container Task Specifications" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Base Image Name & Tag">
              <UInput v-model="form.container_image" placeholder="e.g. redis:7-alpine" required />
            </UFormGroup>
            <UFormGroup label="Container Ports Map">
              <UInput v-model="form.container_port_mapping" placeholder="e.g. 6379:6379" />
            </UFormGroup>
          </div>
          <UFormGroup label="Hosting Parent Node Server or AWS Fargate Instance">
            <USelectMenu v-model="form.host_asset_id" :options="serverAssets" option-attribute="name" value-attribute="id" placeholder="Select parent docker/ECS server host..." searchable />
          </UFormGroup>
        </div>

        <!-- Generic / Other properties form -->
        <div v-else-if="form.type === 'Other'" class="space-y-4 pt-2">
          <UDivider label="Other Spec / Custom Properties" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Model / Identifier">
              <UInput v-model="form.properties.model" placeholder="e.g. APC Smart-UPS" />
            </UFormGroup>
            <UFormGroup label="Serial Number">
              <UInput v-model="form.properties.serial_number" placeholder="Unique serial details" />
            </UFormGroup>
            <UFormGroup label="Custom Group">
              <UInput v-model="form.properties.custom_group" placeholder="e.g. Infrastructure, Facility" />
            </UFormGroup>
            <UFormGroup label="Capacity / Spec">
              <UInput v-model="form.properties.capacity_spec" placeholder="Custom measurements or details" />
            </UFormGroup>
          </div>
        </div>

        <!-- Dynamic Custom Fields Form Block (Phase 1 Custom Field Tabs) -->
        <div v-if="dynamicCustomFields && dynamicCustomFields.length > 0" class="space-y-4 pt-2">
          <UDivider label="Custom Metadata Fields" />
          
          <!-- Tab selector bar for custom attribute groups (Only show if multiple groups exist) -->
          <div v-if="customFieldTabGroups.length > 1" class="flex gap-2 border-b border-slate-200 dark:border-slate-800 pb-2">
            <button
              v-for="group in customFieldTabGroups"
              :key="group"
              type="button"
              :class="[
                activeCustomFieldsTab === group 
                  ? 'border-primary-500 text-primary-500 font-bold border-b-2' 
                  : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300',
                'px-3 py-1 text-xs font-mono transition-all pb-1.5'
              ]"
              @click="activeCustomFieldsTab = group"
            >
              {{ group.toUpperCase() }}
            </button>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup 
              v-for="field in filteredCustomFieldsForTab" 
              :key="field.id" 
              :label="field.name.replace(/_/g, ' ').toUpperCase()"
            >
              <UCheckbox 
                v-if="field.field_type === 'boolean'" 
                v-model="form.properties[field.name]" 
              />
              <UInput 
                v-else-if="field.field_type === 'number'" 
                type="number" 
                v-model="form.properties[field.name]" 
                :placeholder="`Enter ${field.name.replace(/_/g, ' ')}...`" 
              />
              <UInput 
                v-else 
                v-model="form.properties[field.name]" 
                :placeholder="`Enter ${field.name.replace(/_/g, ' ')}...`" 
              />
            </UFormGroup>
          </div>
        </div>

        <!-- Cabinet Mount positioning (Phase 6 DCIM Mapping) -->
        <div v-if="form.type !== 'Kubernetes Deployment' && form.type !== 'Container' && form.type !== 'Application'" class="space-y-4 pt-2 border-t border-slate-200 dark:border-slate-800">
          <UDivider label="Datacenter Cabinet Mounting (DCIM)" />
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Target Server Cabinet / Rack">
              <USelectMenu v-model="form.rack_id" :options="activeRacks" option-attribute="name" value-attribute="id" placeholder="Unmounted (Select rack...)" searchable />
            </UFormGroup>
            <UFormGroup label="Starting Rack Slot (U Height Location)">
              <UInput type="number" v-model="form.rack_position_u" placeholder="e.g. 12" min="1" max="100" />
            </UFormGroup>
          </div>
        </div>
        
        <div class="flex justify-end gap-2 mt-6 pt-3 border-t border-slate-200 dark:border-slate-800">
          <UButton color="gray" variant="ghost" @click="isOpen = false">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="isSaving">Save</UButton>
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
  asset: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isEditing = computed(() => !!props.asset)
const step = ref(1)
const subGroupActive = ref('') // '', 'Network', 'IoT / Edge' (Phase 1 Sub-Group Portals)
const dynamicCategories = ref([])

const activeCategoryRecord = computed(() => {
  if (!form.value.type || !dynamicCategories.value) return null
  return dynamicCategories.value.find(c => c.name === form.value.type)
})

const activeParentDependency = computed(() => {
  return activeCategoryRecord.value?.parent_dependency || ''
})

const fetchDynamicCategories = async () => {
  try {
    const data = await $fetch(`${apiBase}/categories`)
    if (data && data.length > 0) {
      dynamicCategories.value = data
    }
  } catch (err) {
    console.error('Failed to load dynamic categories:', err)
  }
}

const fallbackCategories = [
  { name: 'Server', icon: 'i-heroicons-server', desc: 'Compute server racks, blades, or VMs' },
  { name: 'Database', icon: 'i-heroicons-circle-stack', desc: 'Relational databases or storage engines' },
  { name: 'Network', icon: 'i-heroicons-arrows-right-left', desc: 'Routers, Switches, Access Points, Firewalls, or Balancers' },
  { name: 'Application', icon: 'i-heroicons-window', desc: 'Software apps, microservices, platforms' },
  { name: 'Kubernetes Deployment', icon: 'i-heroicons-square-3-stack-3d', desc: 'Clustered pods, replica specs, cloud workloads' },
  { name: 'Container', icon: 'i-heroicons-cube', desc: 'Standalone containers, AWS Fargate tasks, Docker' },
  { name: 'IoT / Edge', icon: 'i-heroicons-device-tablet', desc: 'Edge gateways, Raspberry Pi, Industrial PLCs, IP Cameras, or Sensors' },
  { name: 'Other', icon: 'i-heroicons-squares-2x2', desc: 'Generic asset types, custom equipment' }
]

const assetTypes = computed(() => {
  if (dynamicCategories.value.length > 0) {
    return dynamicCategories.value.map(c => ({
      name: c.name,
      icon: c.icon || 'i-heroicons-squares-2x2',
      desc: c.description || ''
    }))
  }
  return fallbackCategories
})

const fallbackNetworkSubGroups = [
  { name: 'Router', icon: 'i-heroicons-cpu-chip', desc: 'BGP edge routers, core routers, or gateways' },
  { name: 'Switch (L2)', icon: 'i-heroicons-arrows-right-left', desc: 'Layer-2 switching distribution frames' },
  { name: 'Switch (L3)', icon: 'i-heroicons-arrows-pointing-out', desc: 'Layer-3 switching with IP routing protocols' },
  { name: 'Access Point (AP)', icon: 'i-heroicons-wifi', desc: 'Wireless access points or radio transmitters' },
  { name: 'Firewall', icon: 'i-heroicons-shield-check', desc: 'Hardware security gates and inspection blocks' },
  { name: 'Load Balancer', icon: 'i-heroicons-scale', desc: 'Reverse proxy and network load distribution appliances' }
]

const networkSubGroups = computed(() => {
  const net = dynamicCategories.value.find(c => c.name === 'Network')
  if (net && net.sub_groups && net.sub_groups.length > 0) {
    return net.sub_groups.map(s => ({
      name: s.name,
      icon: s.icon || 'i-heroicons-squares-2x2',
      desc: s.description || ''
    }))
  }
  return fallbackNetworkSubGroups
})

const fallbackIotSubGroups = [
  { name: 'Edge Gateway (e.g. Pi/Jetson)', icon: 'i-heroicons-cpu-chip', desc: 'Raspberry Pi, Jetson Nano, industrial gateways' },
  { name: 'Industrial PLC / SCADA', icon: 'i-heroicons-bolt', desc: 'Programmable Logic Controllers and telemetry lines' },
  { name: 'Smart IP Camera', icon: 'i-heroicons-video-camera', desc: 'Surveillance cameras, streams, security loops' },
  { name: 'Environment IoT Sensor', icon: 'i-heroicons-sun', desc: 'Temperature loop, humidity, air quality sensors' }
]

const iotSubGroups = computed(() => {
  const iot = dynamicCategories.value.find(c => c.name === 'IoT / Edge')
  if (iot && iot.sub_groups && iot.sub_groups.length > 0) {
    return iot.sub_groups.map(s => ({
      name: s.name,
      icon: s.icon || 'i-heroicons-squares-2x2',
      desc: s.description || ''
    }))
  }
  return fallbackIotSubGroups
})

const selectedTypeIcon = computed(() => {
  switch (form.value.type) {
    case 'Server':
      return 'i-heroicons-server'
    case 'Database':
      return 'i-heroicons-circle-stack'
    case 'Network':
      return 'i-heroicons-arrows-right-left'
    case 'Application':
      return 'i-heroicons-window'
    case 'Kubernetes Deployment':
      return 'i-heroicons-square-3-stack-3d'
    case 'Container':
      return 'i-heroicons-cube'
    case 'IoT / Edge':
      return 'i-heroicons-device-tablet'
    case 'Other':
      return 'i-heroicons-squares-2x2'
    default:
      return 'i-heroicons-question-mark-circle'
  }
})

const onCategoryClick = (categoryName) => {
  if (categoryName === 'Network' || categoryName === 'IoT / Edge') {
    subGroupActive.value = categoryName
  } else {
    selectType(categoryName)
  }
}

const selectNetworkSubGroup = (subtypeName) => {
  form.value.type = 'Network'
  form.value.properties = {
    network_subtype: subtypeName
  }
  step.value = 2
}

const selectIotSubGroup = (subtypeName) => {
  form.value.type = 'IoT / Edge'
  form.value.properties = {
    iot_subtype: subtypeName
  }
  step.value = 2
}

const selectType = (typeName) => {
  form.value.type = typeName
  form.value.properties = {}
  step.value = 2
}

const { createAsset, updateAsset, fetchAssets } = useAssets()
const { fetchRelationships, syncRelationships } = useRelationships()
const { fetchRacks } = useDCIM()
const { fetchCustomFields } = useCustomFields()
const { fetchDevices } = useDevices()

const isSaving = ref(false)

const allAssets = ref([])
const serverAssets = ref([])
const activeRacks = ref([])
const dynamicCustomFields = ref([])

// Tabbed Custom Fields Logic (Phase 1 Custom Field Tabs)
const activeCustomFieldsTab = ref('General Custom')

const customFieldTabGroups = computed(() => {
  if (dynamicCustomFields.value.length === 0) return []
  const groups = new Set()
  dynamicCustomFields.value.forEach(field => {
    groups.add(field.tab_group || 'General Custom')
  })
  return Array.from(groups)
})

const filteredCustomFieldsForTab = computed(() => {
  return dynamicCustomFields.value.filter(field => {
    const groupName = field.tab_group || 'General Custom'
    return groupName === activeCustomFieldsTab.value
  })
})

// Auto pre-select first custom field tab group when loaded
watch(customFieldTabGroups, (newGroups) => {
  if (newGroups && newGroups.length > 0) {
    activeCustomFieldsTab.value = newGroups[0]
  } else {
    activeCustomFieldsTab.value = 'General Custom'
  }
}, { immediate: true })
const deviceCatalogEmpty = ref(false)
const devices = ref([])
const selectedDeviceModelId = ref(null)

const loadDevicesCatalog = async () => {
  try {
    const { data } = await fetchDevices()
    devices.value = (data.value || []).map(d => ({
      ...d,
      displayName: `${d.manufacturer?.name || 'Generic'} - ${d.model_name} (v${d.revision || 1})`
    }))
  } catch (error) {
    console.error('Failed to load device catalog:', error)
  }
}

const selectedDeviceModel = computed({
  get: () => {
    if (!selectedDeviceModelId.value || devices.value.length === 0) return null
    return devices.value.find(d => d.id === selectedDeviceModelId.value) || null
  },
  set: (val) => {
    if (val) {
      selectedDeviceModelId.value = val.id
      form.value.device_model_id = val.id
      form.value.properties.manufacturer = val.manufacturer?.name || ''
      form.value.properties.model = val.model_name
      if (val.ports && Object.keys(val.ports).length > 0) {
        form.value.properties.port_config = Object.entries(val.ports)
          .map(([type, count]) => `${count}x ${type}`)
          .join(' + ')
      } else {
        form.value.properties.port_config = ''
      }
      if (!isEditing.value || !form.value.device_model_revision) {
        form.value.device_model_revision = val.revision || 1
      }
      
      // Dynamic Type & Subtype autofilling based on catalog hardware classification
      if (val.subtype) {
        const subtype = val.subtype
        if (subtype === 'Switch (L2)' || subtype === 'Switch (L3)' || subtype === 'Router' || subtype === 'Access Point (AP)' || subtype === 'Firewall' || subtype === 'Load Balancer') {
          form.value.type = 'Network'
          form.value.properties.network_subtype = subtype
        } else if (subtype === 'Server') {
          form.value.type = 'Server'
        } else if (subtype === 'Database') {
          form.value.type = 'Database'
        } else if (subtype === 'Industrial PLC' || subtype === 'Edge Gateway' || subtype === 'Smart IP Camera' || subtype === 'Environment Sensor') {
          form.value.type = 'IoT / Edge'
          form.value.properties.network_subtype = subtype
        } else {
          form.value.type = 'Other'
        }
      }
    } else {
      selectedDeviceModelId.value = null
      form.value.device_model_id = null
      form.value.device_model_revision = null
    }
  }
})

const upgradeDeviceModelRevision = () => {
  if (selectedDeviceModel.value) {
    form.value.device_model_revision = selectedDeviceModel.value.revision
    if (selectedDeviceModel.value.ports && Object.keys(selectedDeviceModel.value.ports).length > 0) {
      form.value.properties.port_config = Object.entries(selectedDeviceModel.value.ports)
        .map(([type, count]) => `${count}x ${type}`)
        .join(' + ')
    } else {
      form.value.properties.port_config = ''
    }
    alert(`Successfully bumped asset to catalog revision v${selectedDeviceModel.value.revision}`)
  }
}

const selectedUplinks = ref([])
const selectedDownlinks = ref([])

const form = ref({
  name: '',
  type: 'Server',
  ip_address: '',
  status: 'active',
  description: '',
  properties: {},
  rack_id: null,
  rack_position_u: null,
  host_asset_id: null,
  container_image: '',
  container_port_mapping: '',
  device_model_id: null,
  device_model_revision: null
})

const auditDeviceCatalog = async () => {
  try {
    const { data } = await fetchDevices()
    deviceCatalogEmpty.value = !data.value || data.value.length === 0
  } catch (error) {
    console.error('Failed to audit device catalog:', error)
  }
}

const loadCategoryCustomFields = async () => {
  if (form.value.type) {
    try {
      const { data } = await fetchCustomFields(form.value.type)
      dynamicCustomFields.value = data.value || []
      dynamicCustomFields.value.forEach(field => {
        if (form.value.properties[field.name] === undefined) {
          form.value.properties[field.name] = field.field_type === 'boolean' ? false : ''
        }
      })
    } catch (err) {
      console.error('Failed to load custom fields:', err)
    }
  }
}

watch(() => form.value.type, () => {
  loadCategoryCustomFields()
})

const loadNetworkAssets = async () => {
  try {
    const { data } = await fetchAssets({ limit: 1000 })
    if (data.value) {
      allAssets.value = data.value.filter(a => 
        (a.type === 'Network') && 
        (!props.asset || a.id !== props.asset.id)
      )
      serverAssets.value = data.value.filter(a => a.type === 'Server')
    }
    
    const { data: racksData } = await fetchRacks()
    activeRacks.value = racksData.value || []
  } catch (error) {
    console.error('Failed to load network assets for mapping:', error)
  }
}

const loadCurrentRelationships = async () => {
  if (isEditing.value) {
    try {
      const { data } = await fetchRelationships(props.asset.id)
      if (data.value) {
        selectedUplinks.value = data.value.uplinks ? data.value.uplinks.map(u => u.id) : []
        selectedDownlinks.value = data.value.downlinks ? data.value.downlinks.map(d => d.id) : []
      }
    } catch (error) {
      console.error('Failed to load relationships:', error)
    }
  } else {
    selectedUplinks.value = []
    selectedDownlinks.value = []
  }
}

watch(() => props.modelValue, (isOpenVal) => {
  if (isOpenVal) {
    subGroupActive.value = ''
    fetchDynamicCategories()
    auditDeviceCatalog()
    loadNetworkAssets()
    loadCurrentRelationships()
    loadCategoryCustomFields()
    loadDevicesCatalog()

    // Explicitly reset form states and wizard step on consecutive creates
    if (!props.asset) {
      form.value = { 
        name: '', 
        type: 'Server', 
        ip_address: '', 
        status: 'active', 
        description: '', 
        properties: {},
        rack_id: null,
        rack_position_u: null,
        host_asset_id: null,
        container_image: '',
        container_port_mapping: '',
        device_model_id: null,
        device_model_revision: null
      }
      selectedDeviceModelId.value = null
      step.value = 1
    }
  }
})

watch(() => props.asset, (newAsset) => {
  if (newAsset) {
    let fixedType = newAsset.type
    let fixedSubtype = newAsset.properties?.network_subtype || ''
    if (newAsset.type === 'Router' || newAsset.type === 'Switch') {
      fixedType = 'Network'
      fixedSubtype = newAsset.type
    }

    form.value = { 
      ...newAsset,
      type: fixedType,
      rack_id: newAsset.rack_id || null,
      rack_position_u: newAsset.rack_position_u || null,
      host_asset_id: newAsset.host_asset_id || null,
      container_image: newAsset.container_image || '',
      container_port_mapping: newAsset.container_port_mapping || '',
      device_model_id: newAsset.device_model_id || null,
      device_model_revision: newAsset.device_model_revision || null,
      properties: {
        ...newAsset.properties,
        network_subtype: fixedSubtype
      }
    }
    selectedDeviceModelId.value = newAsset.device_model_id || null
    step.value = 2
  } else {
    form.value = { 
      name: '', 
      type: 'Server', 
      ip_address: '', 
      status: 'active', 
      description: '', 
      properties: {},
      rack_id: null,
      rack_position_u: null,
      host_asset_id: null,
      container_image: '',
      container_port_mapping: '',
      device_model_id: null,
      device_model_revision: null
    }
    selectedDeviceModelId.value = null
    step.value = 1
  }
}, { immediate: true })

const saveAsset = async () => {
  isSaving.value = true
  
  // Category Parent Dependency Chain Validation (Phase 1 Taxonomy Chains)
  if (activeParentDependency.value) {
    if (!form.value.host_asset_id) {
      alert(`Taxonomy Compliance Error: This category type strictly requires assigning an active parent host/worker asset of type: ${activeParentDependency.value.toUpperCase()}`)
      isSaving.value = false
      return
    }
    const parentAsset = allAssets.value?.find(a => a.id === form.value.host_asset_id)
    if (!parentAsset || parentAsset.type !== activeParentDependency.value) {
      alert(`Taxonomy Compliance Error: The assigned parent host asset must be of type: ${activeParentDependency.value.toUpperCase()}`)
      isSaving.value = false
      return
    }
  }

  // Custom Field Regex Pattern Enforcements (Phase 1 Custom Field Regex)
  if (dynamicCustomFields.value && dynamicCustomFields.value.length > 0) {
    for (const field of dynamicCustomFields.value) {
      if (field.validation_regex) {
        const val = form.value.properties[field.name] || ''
        const regex = new RegExp(field.validation_regex)
        if (!regex.test(String(val))) {
          alert(`Validation Pattern Error: Custom Attribute '${field.name.replace(/_/g, ' ').toUpperCase()}' must strictly match pattern: ${field.validation_regex}`)
          isSaving.value = false
          return
        }
      }
    }
  }

  try {
    const payload = {
      ...form.value,
      rack_id: form.value.rack_id ? form.value.rack_id : null,
      rack_position_u: form.value.rack_position_u ? Number(form.value.rack_position_u) : null,
      host_asset_id: form.value.host_asset_id ? form.value.host_asset_id : null,
    }

    let savedAsset = null
    if (isEditing.value) {
      savedAsset = await updateAsset(props.asset.id, payload)
    } else {
      savedAsset = await createAsset(payload)
    }
    
    if (form.value.type === 'Network') {
      const targetId = isEditing.value ? props.asset.id : savedAsset.id
      await syncRelationships(targetId, selectedUplinks.value, selectedDownlinks.value)
    }
    
    isOpen.value = false
    emit('saved')
  } catch (error) {
    console.error('Failed to save asset:', error)
  } finally {
    isSaving.value = false
  }
}
</script>
