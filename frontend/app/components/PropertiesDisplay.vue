<template>
  <template v-if="asset.properties && Object.keys(asset.properties).length > 0">
    <template v-if="asset.type?.toLowerCase() === 'server'">
      <div class="sm:col-span-2 mt-4">
        <UDivider label="Hardware & System Details" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.manufacturer"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Manufacturer</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.manufacturer }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.model"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Model</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.model }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.serial_number"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Serial Number</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.serial_number }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.form_factor"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Form Factor</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.form_factor }}</dd></div>
      
      <div class="sm:col-span-2 mt-4" v-if="asset.properties.cpu_model || asset.properties.ram_capacity_gb">
        <UDivider label="Compute & Memory" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.cpu_model"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">CPU</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.cpu_model }} ({{ asset.properties.cpu_sockets }} Sockets, {{ asset.properties.cpu_cores }} Cores)</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.ram_capacity_gb"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">RAM</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.ram_capacity_gb }} GB ({{ asset.properties.ram_type }})</dd></div>

      <div class="sm:col-span-2 mt-4" v-if="asset.properties.storage_subsystem || asset.properties.pcie_cards">
        <UDivider label="Storage & Expansion" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.storage_subsystem"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Storage Subsystem</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.storage_subsystem }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.storage_controller"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Storage Controller</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.storage_controller }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.pcie_cards"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">PCI-e Expansion</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.pcie_cards }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.network_interfaces"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Network Interfaces</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.network_interfaces }}</dd></div>

      <div class="sm:col-span-2 mt-4" v-if="asset.properties.os_name">
        <UDivider label="Operating System" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.os_name"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">OS</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.os_name }} {{ asset.properties.os_version }}</dd></div>
    </template>

    <template v-else-if="asset.type?.toLowerCase() === 'router' || asset.type?.toLowerCase() === 'switch'">
      <div class="sm:col-span-2 mt-4">
        <UDivider label="Hardware & Identity" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.manufacturer"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Manufacturer</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.manufacturer }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.model"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Model</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.model }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.serial_number"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Serial Number</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.serial_number }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.form_factor"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Form Factor</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.form_factor }}</dd></div>

      <div class="sm:col-span-2 mt-4">
        <UDivider label="Software & Configuration" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.firmware"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Firmware/OS</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.firmware }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.management_ip"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Management IP</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.management_ip }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.mac_address"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">MAC Address</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.mac_address }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.network_role"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Role</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.network_role }}</dd></div>

      <div class="sm:col-span-2 mt-4">
        <UDivider label="Capacity & Interfaces" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.port_config"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Ports</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.port_config }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.throughput"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Throughput</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.throughput }}</dd></div>
    </template>

    <template v-else-if="asset.type?.toLowerCase() === 'database'">
      <div class="sm:col-span-2 mt-4">
        <UDivider label="Database Engine & Version" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.engine"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Engine</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.engine }} {{ asset.properties.version }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.edition"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Edition</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.edition }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.storage_engine"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Storage Engine</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.storage_engine }}</dd></div>

      <div class="sm:col-span-2 mt-4">
        <UDivider label="Deployment & Capacity" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.host_cluster"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Host/Cluster</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.host_cluster }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.environment"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Environment</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.environment }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.allocated_storage_gb"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Allocated Storage</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.allocated_storage_gb }} GB</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.ha_setup"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">HA Setup</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.ha_setup }}</dd></div>
    </template>

    <template v-else-if="asset.type?.toLowerCase() === 'application'">
      <div class="sm:col-span-2 mt-4">
        <UDivider label="Application Details" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.vendor"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Vendor/Publisher</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.vendor }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.version"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Version/Build</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.version }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.environment"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Environment</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.environment }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.criticality"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Criticality</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.criticality }}</dd></div>

      <div class="sm:col-span-2 mt-4">
        <UDivider label="Technical Stack" />
      </div>
      <div class="sm:col-span-1" v-if="asset.properties.framework"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Framework/Language</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.framework }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.frontend_tech"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Frontend Tech</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.frontend_tech }}</dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.url"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">URL</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white"><a :href="asset.properties.url" target="_blank" class="text-primary-500 hover:underline">{{ asset.properties.url }}</a></dd></div>
      <div class="sm:col-span-1" v-if="asset.properties.auth_method"><dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Authentication</dt><dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.properties.auth_method }}</dd></div>
    </template>

    <template v-else>
      <div class="sm:col-span-2 mt-4">
        <UDivider label="Specific Details" />
      </div>
      <div v-for="(value, key) in asset.properties" :key="key" class="sm:col-span-1">
        <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 capitalize">{{ key.replace(/_/g, ' ') }}</dt>
        <dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ value || 'N/A' }}</dd>
      </div>
    </template>
  </template>
</template>

<script setup>
defineProps({
  asset: {
    type: Object,
    required: true
  }
})
</script>