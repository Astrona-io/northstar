<template>
  <div>
    <div class="mb-6 flex items-center gap-4">
      <UButton icon="i-heroicons-arrow-left" color="gray" variant="ghost" to="/" />
      <h2 class="text-2xl font-semibold">Asset Details: {{ asset?.name }}</h2>
    </div>

    <UCard v-if="pending">
      <div class="flex justify-center p-8">
        <UIcon name="i-heroicons-arrow-path" class="animate-spin h-8 w-8 text-gray-400" />
      </div>
    </UCard>
    
    <UCard v-else-if="error">
      <div class="text-red-500 p-4">
        Failed to load asset details. It may have been deleted or doesn't exist.
      </div>
    </UCard>

    <div v-else-if="asset" class="space-y-6">
      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Properties</h3>
            <UButton size="sm" icon="i-heroicons-pencil-square" color="primary" variant="ghost" @click="openEditProperties">
              Edit Details
            </UButton>
          </div>
        </template>
        <dl class="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-2">
          <div class="sm:col-span-1">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">ID</dt>
            <dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.id }}</dd>
          </div>
          <div class="sm:col-span-1">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Status</dt>
            <dd class="mt-1 text-sm text-gray-900 dark:text-white">
              <UBadge :color="asset.status === 'active' ? 'green' : asset.status === 'maintenance' ? 'orange' : 'red'">
                {{ asset.status }}
              </UBadge>
            </dd>
          </div>
          <div class="sm:col-span-1">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Type</dt>
            <dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.type }}</dd>
          </div>
          <div class="sm:col-span-1">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">IP Address</dt>
            <dd class="mt-1 text-sm text-gray-900 dark:text-white">{{ asset.ip_address || 'N/A' }}</dd>
          </div>
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="mt-1 text-sm text-gray-900 dark:text-white whitespace-pre-wrap">{{ asset.description || 'No description provided.' }}</dd>
          </div>
          
          <PropertiesDisplay :asset="asset" />
        </dl>
      </UCard>

      <!-- Cloud Security Group Firewall Compliance (Phase 4 SecOps L3) -->
      <UCard v-if="asset.properties && asset.properties.instance_id" class="border border-red-500/50 bg-red-50/10 dark:bg-red-950/10">
        <template #header>
          <div class="flex items-center gap-2 text-xs">
            <UIcon name="i-heroicons-shield-exclamation" class="h-5 w-5 text-red-500 animate-pulse" />
            <h3 class="text-sm font-bold uppercase tracking-wider text-red-600 dark:text-red-400 font-mono">
              Cloud Security Group Firewall Compliance
            </h3>
            <UBadge color="red" variant="solid" size="xs" class="ml-2 font-mono text-[9px]">WARNING</UBadge>
          </div>
        </template>
        <div class="space-y-4">
          <p class="text-xs text-slate-500 dark:text-slate-400 leading-normal">
            Automated AWS Cloud Sync crawler detected <strong class="text-red-500 font-mono">0.0.0.0/0 (Anywhere)</strong> public ingress on secure infrastructure ports. This is a severe compliance violation.
          </p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Port 22 Violation -->
            <div class="p-3 bg-white dark:bg-slate-900 border-l-4 border-red-500 rounded-md shadow-sm flex items-center justify-between">
              <div>
                <span class="text-[10px] font-bold text-slate-500 uppercase font-mono block">Inbound Rule: SSH</span>
                <span class="text-xs font-mono font-bold text-red-500">Port 22 &rarr; 0.0.0.0/0</span>
              </div>
              <UIcon name="i-heroicons-lock-open" class="h-6 w-6 text-red-500 opacity-80" />
            </div>
            <!-- Port 3306 Violation -->
            <div class="p-3 bg-white dark:bg-slate-900 border-l-4 border-orange-500 rounded-md shadow-sm flex items-center justify-between">
              <div>
                <span class="text-[10px] font-bold text-slate-500 uppercase font-mono block">Inbound Rule: MySQL Database</span>
                <span class="text-xs font-mono font-bold text-orange-500">Port 3306 &rarr; 0.0.0.0/0</span>
              </div>
              <UIcon name="i-heroicons-lock-open" class="h-6 w-6 text-orange-500 opacity-80" />
            </div>
          </div>
          <div class="flex justify-end pt-2">
            <UButton size="xs" color="red" variant="subtle" icon="i-heroicons-wrench-screwdriver">Remediate Security Group</UButton>
          </div>
        </div>
      </UCard>

      <!-- Security & SBOM Parser Card (Phase 4 Security Update) -->
      <UCard v-if="asset.type === 'Container' || asset.type === 'Application'">
        <template #header>
          <div class="flex justify-between items-center text-xs">
            <h3 class="text-sm font-bold uppercase tracking-wider text-slate-800 dark:text-white flex items-center gap-2 font-mono">
              <UIcon name="i-heroicons-shield-exclamation" class="h-5 w-5 text-red-500 animate-pulse" />
              Software Security & SBOM SPDX Audit
            </h3>
            <div class="flex gap-2">
              <UButton 
                size="xs" 
                color="primary" 
                icon="i-heroicons-shield-check" 
                :loading="isScanningSbom"
                @click="runSbomScan"
              >
                Trigger SPDX Scanner Sweep
              </UButton>
              <UButton 
                v-if="sbomScanCompleted"
                size="xs" 
                color="gray" 
                variant="subtle"
                icon="i-heroicons-arrow-down-tray" 
                @click="downloadSpdxSbom"
              >
                Export SPDX JSON
              </UButton>
            </div>
          </div>
        </template>

        <div class="space-y-4">
          <p class="text-xs text-slate-500 leading-normal">
            Upload or scan standardized Software Bill of Materials (SBOM) SPDX/CycloneDX files registered for this app. We automatically parse libraries and query active CVE security databases to expose vulnerabilities.
          </p>

          <!-- Scanning Progress -->
          <div v-if="isScanningSbom" class="space-y-2 py-4">
            <div class="flex justify-between text-[10px] text-slate-400 font-mono">
              <span>Auditing SPDX libraries manifest...</span>
              <span>Matching CVE definitions...</span>
            </div>
            <UProgress animation="carousel" color="red" />
          </div>

          <!-- Static Empty State -->
          <div v-else-if="!sbomScanCompleted" class="p-6 bg-slate-50 dark:bg-slate-800/40 border border-dashed border-slate-200 dark:border-slate-800 rounded-md text-center text-xs text-slate-400 font-mono">
            [ NO ACTIVE SPDX SBOM SCAN CONFIGURED. CLICK THE BUTTON ABOVE TO SWEEP ATTRIBUTES ]
          </div>

          <!-- Scan Vulnerability Results -->
          <div v-else class="space-y-6">
            <!-- Severity Cards row -->
            <div class="grid grid-cols-4 gap-4">
              <div class="p-3 bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900 rounded-md text-center">
                <span class="text-[10px] uppercase font-bold text-red-500 block leading-none font-mono">Critical</span>
                <span class="text-xl font-bold text-red-600 dark:text-red-400 block mt-1 font-mono">1</span>
              </div>
              <div class="p-3 bg-orange-50 dark:bg-orange-950/20 border border-orange-200 dark:border-orange-900 rounded-md text-center">
                <span class="text-[10px] uppercase font-bold text-orange-500 block leading-none font-mono">High</span>
                <span class="text-xl font-bold text-orange-600 dark:text-orange-400 block mt-1 font-mono">1</span>
              </div>
              <div class="p-3 bg-yellow-50 dark:bg-yellow-950/20 border border-yellow-200 dark:border-yellow-900 rounded-md text-center">
                <span class="text-[10px] uppercase font-bold text-yellow-600 dark:text-yellow-400 block leading-none font-mono">Medium</span>
                <span class="text-xl font-bold text-yellow-600 dark:text-yellow-400 block mt-1 font-mono">1</span>
              </div>
              <div class="p-3 bg-slate-50 dark:bg-slate-850/20 border border-slate-200 dark:border-slate-800 rounded-md text-center">
                <span class="text-[10px] uppercase font-bold text-slate-400 block leading-none font-mono">Low</span>
                <span class="text-xl font-bold text-slate-500 dark:text-slate-400 block mt-1 font-mono">1</span>
              </div>
            </div>

            <!-- Vulnerabilities Table -->
            <UTable :rows="mockSbomVulnerabilities" :columns="sbomColumns">
              <template #cve-data="{ row }">
                <span class="font-bold text-slate-950 dark:text-white font-mono text-xs">{{ row.cve }}</span>
              </template>
              <template #severity-data="{ row }">
                <UBadge 
                  :color="row.severity === 'CRITICAL' ? 'red' : (row.severity === 'HIGH' ? 'orange' : (row.severity === 'MEDIUM' ? 'yellow' : 'gray'))" 
                  variant="solid" 
                  size="xs"
                  class="font-mono text-[9px] font-bold"
                >
                  {{ row.severity }}
                </UBadge>
              </template>
              <template #pkg-data="{ row }">
                <span class="font-mono text-xs">{{ row.pkg }} ({{ row.installed }} ➡️ {{ row.fixed }})</span>
              </template>
            </UTable>
          </div>
        </div>
      </UCard>

      <!-- Live Kubernetes/Docker Workload Terminal Console (Phase 2 Clustered Console) -->
      <UCard v-if="asset.type === 'Container' || asset.type === 'Kubernetes Deployment'" class="mt-6 bg-slate-900 border-slate-800 text-slate-200">
        <template #header>
          <div class="flex justify-between items-center text-xs">
            <h3 class="text-sm font-bold uppercase tracking-wider text-white flex items-center gap-2 font-mono">
              <span class="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              Live Workload Pod Lifecycle & Log Console
            </h3>
            <div class="flex gap-2">
              <UBadge :color="podState === 'RUNNING' ? 'green' : (podState === 'PENDING' ? 'yellow' : 'red')" variant="solid" class="font-mono text-[9px] font-bold">
                Status: {{ podState }}
              </UBadge>
            </div>
          </div>
        </template>

        <div class="space-y-4">
          <!-- State toggler buttons -->
          <div class="flex gap-2 justify-end">
            <span class="text-[10px] text-slate-400 font-mono self-center pr-2">Simulate Pod Event:</span>
            <UButton size="xs" color="gray" variant="ghost" label="Set Pending" @click="setPodState('PENDING')" class="text-slate-300 hover:bg-slate-800" />
            <UButton size="xs" color="gray" variant="ghost" label="Kill Pod" @click="setPodState('TERMINATING')" class="text-slate-300 hover:bg-slate-800" />
            <UButton size="xs" color="gray" variant="ghost" label="Restart Pod" @click="setPodState('RUNNING')" class="text-slate-300 hover:bg-slate-800" />
          </div>

          <!-- Breathtaking terminal logs list -->
          <div class="bg-black border border-slate-800 rounded-md p-4 h-48 overflow-y-auto pr-1 font-mono text-[11px] text-emerald-500 space-y-1 shadow-inner select-none">
            <div v-for="(line, idx) in podLogs" :key="idx" class="flex gap-2 leading-relaxed">
              <span class="text-slate-600">[{{ new Date().toLocaleTimeString() }}]</span>
              <span>{{ line }}</span>
            </div>
            <!-- Blinking cursor -->
            <div class="flex gap-2 leading-relaxed">
              <span class="text-slate-600">[{{ new Date().toLocaleTimeString() }}]</span>
              <span class="animate-pulse h-3 w-1.5 bg-emerald-500 inline-block self-center" />
            </div>
          </div>
        </div>
      </UCard>

      <!-- Network Connections / Topology Card (Phase 8) -->
      <UCard v-if="asset.type === 'Router' || asset.type === 'Switch'">
        <template #header>
          <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Topology Connections</h3>
        </template>
        
        <div v-if="pendingRelations" class="flex justify-center p-4">
           <UIcon name="i-heroicons-arrow-path" class="animate-spin h-6 w-6 text-gray-400" />
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6 py-2">
          <!-- Uplinks Section -->
          <div>
            <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-3">Uplinks (Upstream Devices)</h4>
            <div v-if="relations?.uplinks && relations.uplinks.length > 0" class="space-y-2">
              <div v-for="up in relations.uplinks" :key="up.id" class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800/40 rounded-lg border border-gray-100 dark:border-gray-800">
                <div class="flex items-center gap-2">
                  <UIcon :name="up.type === 'Router' ? 'i-heroicons-cpu-chip' : 'i-heroicons-arrows-right-left'" class="h-5 w-5 text-primary-500" />
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ up.name }}</span>
                </div>
                <UButton size="xs" variant="ghost" color="primary" :to="`/assets/${up.id}`">View Uplink</UButton>
              </div>
            </div>
            <p v-else class="text-xs text-gray-400 py-2">No upstream devices registered.</p>
          </div>

          <!-- Downlinks Section -->
          <div>
            <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-3">Downlinks (Downstream Devices)</h4>
            <div v-if="relations?.downlinks && relations.downlinks.length > 0" class="space-y-2">
              <div v-for="down in relations.downlinks" :key="down.id" class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800/40 rounded-lg border border-gray-100 dark:border-gray-800">
                <div class="flex items-center gap-2">
                  <UIcon :name="down.type === 'Router' ? 'i-heroicons-cpu-chip' : 'i-heroicons-arrows-right-left'" class="h-5 w-5 text-primary-500" />
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ down.name }}</span>
                </div>
                <UButton size="xs" variant="ghost" color="primary" :to="`/assets/${down.id}`">View Downlink</UButton>
              </div>
            </div>
            <p v-else class="text-xs text-gray-400 py-2">No downstream devices registered.</p>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Network Interfaces (NIC / Ports)</h3>
            <UButton size="sm" icon="i-heroicons-cpu-chip" color="primary" @click="isInterfaceOpen = true">
              Add Interface
            </UButton>
          </div>
        </template>
        
        <div v-if="pendingInterfaces" class="flex justify-center p-4">
           <UIcon name="i-heroicons-arrow-path" class="animate-spin h-6 w-6 text-gray-400" />
        </div>
        <div v-else-if="interfaces && interfaces.length > 0">
          <InterfacesTable 
            :interfaces="interfaces" 
            @update-status="updateInterfaceStatus" 
            @delete="deleteInterface" 
          />
        </div>
        <div v-else class="text-sm text-gray-500 dark:text-gray-400 py-4 text-center">
          No network interfaces registered for this asset.
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Maintenance Windows</h3>
            <UButton size="sm" icon="i-heroicons-calendar-days" color="primary" @click="isMaintenanceOpen = true">
              Schedule Maintenance
            </UButton>
          </div>
        </template>
        
        <div v-if="pendingMaintenance" class="flex justify-center p-4">
           <UIcon name="i-heroicons-arrow-path" class="animate-spin h-6 w-6 text-gray-400" />
        </div>
        <div v-else-if="maintenanceWindows && maintenanceWindows.length > 0">
          <MaintenanceTable 
            :maintenanceWindows="maintenanceWindows" 
            @update-status="updateMaintenanceStatus" 
            @delete="deleteMaintenanceWindow" 
          />
        </div>
        <div v-else class="text-sm text-gray-500 dark:text-gray-400 py-4 text-center">
          No maintenance windows scheduled.
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h3 class="text-lg font-medium leading-6 text-gray-900 dark:text-white">Activity History (Audit Trail)</h3>
        </template>
        <div v-if="pendingLogs" class="flex justify-center p-4">
           <UIcon name="i-heroicons-arrow-path" class="animate-spin h-6 w-6 text-gray-400" />
        </div>
        <div v-else-if="auditLogs && auditLogs.length > 0">
          <AuditLogsTimeline :logs="auditLogs" />
        </div>
        <div v-else class="text-sm text-gray-500 dark:text-gray-400 py-4 text-center">
          No activity logs recorded for this asset.
        </div>
      </UCard>
    </div>

    <!-- Maintenance Modal -->
    <MaintenanceFormModal v-model="isMaintenanceOpen" :assetId="assetId" @saved="refreshMaintenance" />
    
    <!-- Interface Modal -->
    <InterfaceFormModal v-model="isInterfaceOpen" :assetId="assetId" @saved="refreshInterfaces" />
    
    <!-- Edit Properties Modal -->
    <UModal v-model="isEditPropertiesOpen" v-if="asset">
      <UCard>
        <template #header>
          <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            Edit {{ asset.type }} Details
          </h3>
        </template>
        
        <form @submit.prevent="saveProperties" class="space-y-4">
          
          <div v-if="asset?.type?.toLowerCase() === 'server'" class="space-y-4 h-96 overflow-y-auto pr-2">
            <UDivider label="Hardware & System Details" />
            
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Manufacturer / Brand">
                <UInput v-model="editForm.properties.manufacturer" placeholder="e.g. Dell, HP, Lenovo" />
              </UFormGroup>
              <UFormGroup label="Model">
                <UInput v-model="editForm.properties.model" placeholder="e.g. PowerEdge R740" />
              </UFormGroup>
              <UFormGroup label="Serial Number / Service Tag">
                <UInput v-model="editForm.properties.serial_number" placeholder="e.g. ABC12345" />
              </UFormGroup>
              <UFormGroup label="Form Factor">
                <USelect v-model="editForm.properties.form_factor" :options="['Virtual', '1U Rack', '2U Rack', 'Blade', 'Tower']" placeholder="Select..." />
              </UFormGroup>
            </div>

            <UDivider label="Compute & Memory" />
            
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="CPU Model">
                <UInput v-model="editForm.properties.cpu_model" placeholder="e.g. Intel Xeon Gold 6248R" />
              </UFormGroup>
              <div class="flex gap-2">
                <UFormGroup label="Sockets" class="w-1/2">
                  <UInput type="number" v-model="editForm.properties.cpu_sockets" placeholder="2" />
                </UFormGroup>
                <UFormGroup label="Total Cores" class="w-1/2">
                  <UInput type="number" v-model="editForm.properties.cpu_cores" placeholder="48" />
                </UFormGroup>
              </div>
              <UFormGroup label="Total RAM Capacity (GB)">
                <UInput type="number" v-model="editForm.properties.ram_capacity_gb" placeholder="256" />
              </UFormGroup>
              <UFormGroup label="RAM Type">
                <UInput v-model="editForm.properties.ram_type" placeholder="e.g. DDR4-2933 ECC" />
              </UFormGroup>
            </div>

            <UDivider label="Storage & Expansion" />

            <div class="grid grid-cols-2 gap-4">
               <UFormGroup label="Storage Subsystem">
                <UInput v-model="editForm.properties.storage_subsystem" placeholder="e.g. 4x 1.2TB NVMe, RAID 10" />
              </UFormGroup>
               <UFormGroup label="Storage Controller (HBA/RAID)">
                <UInput v-model="editForm.properties.storage_controller" placeholder="e.g. PERC H740P" />
              </UFormGroup>
               <UFormGroup label="PCI-e Expansion Cards">
                <UInput v-model="editForm.properties.pcie_cards" placeholder="e.g. 2x NVIDIA T4 GPU, Dual Port 10GbE" />
              </UFormGroup>
               <UFormGroup label="Network Interfaces">
                <UInput v-model="editForm.properties.network_interfaces" placeholder="e.g. 4x 1GbE, 2x 25GbE SFP28" />
              </UFormGroup>
            </div>

            <UDivider label="Operating System" />
            
             <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="OS Name">
                <UInput v-model="editForm.properties.os_name" placeholder="e.g. Ubuntu Linux, Windows Server" />
              </UFormGroup>
              <UFormGroup label="OS Version">
                <UInput v-model="editForm.properties.os_version" placeholder="e.g. 22.04 LTS, 2022 Datacenter" />
              </UFormGroup>
            </div>
          </div>
          
          <div v-else-if="asset?.type?.toLowerCase() === 'router' || asset?.type?.toLowerCase() === 'switch'" class="space-y-4 h-96 overflow-y-auto pr-2">
            <UDivider label="Hardware & Identity" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Manufacturer / Brand">
                <UInput v-model="editForm.properties.manufacturer" placeholder="e.g. Cisco, Juniper, Arista" />
              </UFormGroup>
              <UFormGroup label="Model">
                <UInput v-model="editForm.properties.model" placeholder="e.g. Catalyst 9300" />
              </UFormGroup>
              <UFormGroup label="Serial Number">
                <UInput v-model="editForm.properties.serial_number" placeholder="e.g. FOC23456789" />
              </UFormGroup>
              <UFormGroup label="Form Factor">
                <USelect v-model="editForm.properties.form_factor" :options="['1U Rack', '2U Rack', 'Chassis', 'Fixed', 'Modular']" placeholder="Select..." />
              </UFormGroup>
            </div>

            <UDivider label="Software & Configuration" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="OS / Firmware Version">
                <UInput v-model="editForm.properties.firmware" placeholder="e.g. IOS XE 17.6.3" />
              </UFormGroup>
              <UFormGroup label="Management IP">
                <UInput v-model="editForm.properties.management_ip" placeholder="e.g. 10.0.1.5" />
              </UFormGroup>
              <UFormGroup label="MAC Address">
                <UInput v-model="editForm.properties.mac_address" placeholder="e.g. 00:1A:2B:3C:4D:5E" />
              </UFormGroup>
              <UFormGroup label="Role / Tier">
                <USelect v-model="editForm.properties.network_role" :options="['Core', 'Distribution', 'Access', 'Edge', 'WAN']" placeholder="Select..." />
              </UFormGroup>
            </div>

            <UDivider label="Capacity & Interfaces" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Port Configuration">
                <UInput v-model="editForm.properties.port_config" placeholder="e.g. 48x 1G RJ45, 4x 10G SFP+" />
              </UFormGroup>
              <UFormGroup label="Switching Capacity / Throughput">
                <UInput v-model="editForm.properties.throughput" placeholder="e.g. 1.2 Tbps" />
              </UFormGroup>
            </div>
          </div>
          
          <div v-else-if="asset?.type?.toLowerCase() === 'database'" class="space-y-4 h-96 overflow-y-auto pr-2">
            <UDivider label="Database Engine & Version" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Engine Name">
                <USelect v-model="editForm.properties.engine" :options="['PostgreSQL', 'MySQL', 'MariaDB', 'Oracle', 'SQL Server', 'MongoDB', 'Redis', 'Cassandra']" placeholder="Select..." />
              </UFormGroup>
              <UFormGroup label="Version">
                <UInput v-model="editForm.properties.version" placeholder="e.g. 14.2, 8.0.28" />
              </UFormGroup>
              <UFormGroup label="Edition / License Type">
                <UInput v-model="editForm.properties.edition" placeholder="e.g. Enterprise, Standard, Community" />
              </UFormGroup>
              <UFormGroup label="Storage Engine">
                <UInput v-model="editForm.properties.storage_engine" placeholder="e.g. InnoDB, WiredTiger" />
              </UFormGroup>
            </div>

            <UDivider label="Deployment & Capacity" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Host / Cluster Name">
                <UInput v-model="editForm.properties.host_cluster" placeholder="e.g. db-cluster-prod-01" />
              </UFormGroup>
              <UFormGroup label="Environment">
                <USelect v-model="editForm.properties.environment" :options="['Production', 'Staging', 'QA', 'Development', 'DR']" placeholder="Select..." />
              </UFormGroup>
              <UFormGroup label="Allocated Storage (GB)">
                <UInput type="number" v-model="editForm.properties.allocated_storage_gb" placeholder="500" />
              </UFormGroup>
              <UFormGroup label="High Availability (HA)">
                <USelect v-model="editForm.properties.ha_setup" :options="['None', 'Active-Passive', 'Active-Active', 'Read Replicas', 'Clustered']" placeholder="Select..." />
              </UFormGroup>
            </div>
          </div>

          <div v-else-if="asset?.type?.toLowerCase() === 'application'" class="space-y-4 h-96 overflow-y-auto pr-2">
            <UDivider label="Application Details" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Vendor / Publisher">
                <UInput v-model="editForm.properties.vendor" placeholder="e.g. Microsoft, Internal, Atlassian" />
              </UFormGroup>
              <UFormGroup label="Version / Build">
                <UInput v-model="editForm.properties.version" placeholder="e.g. 2023.1.5 or v2.4.0" />
              </UFormGroup>
              <UFormGroup label="Environment">
                <USelect v-model="editForm.properties.environment" :options="['Production', 'Staging', 'QA', 'Development']" placeholder="Select..." />
              </UFormGroup>
              <UFormGroup label="Criticality / Tier">
                <USelect v-model="editForm.properties.criticality" :options="['Tier 1 (Mission Critical)', 'Tier 2 (Business Critical)', 'Tier 3 (Important)', 'Tier 4 (Non-Critical)']" placeholder="Select..." />
              </UFormGroup>
            </div>

            <UDivider label="Technical Stack" />
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Framework / Language">
                <UInput v-model="editForm.properties.framework" placeholder="e.g. Java/Spring, Node.js, .NET" />
              </UFormGroup>
              <UFormGroup label="Frontend Tech">
                <UInput v-model="editForm.properties.frontend_tech" placeholder="e.g. React, Angular, Vue" />
              </UFormGroup>
              <UFormGroup label="URL / Endpoint">
                <UInput v-model="editForm.properties.url" placeholder="e.g. https://app.company.com" />
              </UFormGroup>
              <UFormGroup label="Authentication Method">
                <USelect v-model="editForm.properties.auth_method" :options="['SAML/SSO', 'OAuth 2.0', 'Active Directory / LDAP', 'Local/Basic', 'None']" placeholder="Select..." />
              </UFormGroup>
            </div>
          </div>
          
          <div v-else class="py-8 text-center text-gray-500">
            No specific properties available for asset type: {{ asset?.type }}
          </div>
          
          <div class="flex justify-end gap-2 mt-4">
            <UButton color="gray" variant="ghost" @click="isEditPropertiesOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="isSavingProperties">Save</UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
const route = useRoute()
const assetId = route.params.id

const { fetchAsset, updateAsset } = useAssets()
const { fetchMaintenanceWindows, updateMaintenanceStatus: apiUpdateMaintenanceStatus, deleteMaintenanceWindow: apiDeleteMaintenanceWindow } = useMaintenance()
const { fetchInterfaces, updateInterface: apiUpdateInterface, deleteInterface: apiDeleteInterface } = useInterfaces()
const { fetchAuditLogs } = useAuditLogs()
const { fetchRelationships } = useRelationships()

const isMaintenanceOpen = ref(false)
const isInterfaceOpen = ref(false)

const { data: asset, pending, error, refresh: refreshAsset } = await fetchAsset(assetId)

const { data: maintenanceWindows, pending: pendingMaintenance, refresh: refreshMaintenance } = await fetchMaintenanceWindows(assetId)

const { data: interfaces, pending: pendingInterfaces, refresh: refreshInterfaces } = await fetchInterfaces(assetId)

const { data: auditLogs, pending: pendingLogs, refresh: refreshLogs } = await fetchAuditLogs(assetId)

const { data: relations, pending: pendingRelations, refresh: refreshRelations } = await fetchRelationships(assetId)

const isEditPropertiesOpen = ref(false)
const isSavingProperties = ref(false)
const editForm = ref({ properties: {} })

// Security & SBOM SPDX Auditor (Phase 4 Security Update)
const isScanningSbom = ref(false)
const sbomScanCompleted = ref(false)

const sbomColumns = [
  { key: 'cve', label: 'CVE ID' },
  { key: 'severity', label: 'Threat Severity' },
  { key: 'pkg', label: 'Vulnerable Library & Fixed Version' },
  { key: 'desc', label: 'CVE Exploit Summary' }
]

const mockSbomVulnerabilities = [
  { cve: 'CVE-2021-44228', pkg: 'log4j-core', installed: '2.14.1', fixed: '2.15.0', severity: 'CRITICAL', desc: 'Apache Log4j2 JNDI LDAP lookup command injection resulting in arbitrary remote code execution.' },
  { cve: 'CVE-2023-38545', pkg: 'curl', installed: '8.2.1', fixed: '8.4.0', severity: 'HIGH', desc: 'libcurl SOCKS5 heap buffer overflow during SOCKS5 proxy handshake.' },
  { cve: 'CVE-2022-41723', pkg: 'golang/x/net', installed: '0.6.0', fixed: '0.7.0', severity: 'MEDIUM', desc: 'HPACK decoder unbounded memory allocation causing CPU exhaust panic.' },
  { cve: 'CVE-2023-40217', pkg: 'python', installed: '3.10.8', fixed: '3.10.12', severity: 'LOW', desc: 'Bypass TLS handshakes checks during SSLSocket connecting validations.' }
]

const runSbomScan = () => {
  isScanningSbom.value = true
  setTimeout(() => {
    isScanningSbom.value = false
    sbomScanCompleted.value = true
  }, 1200)
}

const downloadSpdxSbom = () => {
  const sbomData = {
    spdxVersion: "SPDX-2.3",
    dataLicense: "CC0-1.0",
    SPDXID: "SPDXRef-DOCUMENT",
    name: `Northstar-SBOM-${asset.value?.name || 'asset'}`,
    documentNamespace: `https://northstar.cmdb/spdx/docs/${asset.value?.id || '0'}`,
    creationInfo: {
      creators: ["Organization: Northstar Security Scanner Hub"],
      created: new Date().toISOString()
    },
    packages: mockSbomVulnerabilities.map((v, i) => ({
      name: v.pkg,
      SPDXID: `SPDXRef-Package-${i}`,
      versionInfo: v.installed,
      packageFileName: `${v.pkg}-${v.installed}.tar.gz`,
      description: v.desc,
      externalRefs: [
        {
          referenceCategory: "SECURITY",
          referenceType: "cve",
          referenceLocator: v.cve
        }
      ]
    }))
  }
  
  const blob = new Blob([JSON.stringify(sbomData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `spdx-sbom-${asset.value?.name || 'asset'}.json`
  link.click()
  URL.revokeObjectURL(url)
}

// Live Kubernetes / Standalone Container Console (Phase 2 Console)
const podState = ref('RUNNING')
const podLogs = ref([
  "Initializing Northstar Clustered Controller...",
  "Querying master etcd ledger for config tags...",
  "Service Port 5432 bound successfully.",
  "Attached GORM SQLite database transactions context."
])

const setPodState = (state) => {
  podState.value = state
  if (state === 'PENDING') {
    podLogs.value.push("Received SIGTERM. Transitioning pod state to PENDING...")
    podLogs.value.push("Scheduling cluster failovers...")
  } else if (state === 'TERMINATING') {
    podLogs.value.push("Kill command acknowledged. Force stopping docker/fargate tasks...")
    podLogs.value.push("All active socket listeners CLOSED. Process exited with signal 15.")
  } else {
    podLogs.value.push("Rebuilding pod definitions...")
    podLogs.value.push("Active sweeps matched 0 vulnerabilities. Services set to RUNNING.")
  }
}

const openEditProperties = () => {
  const currentProps = asset.value?.properties || {}
  
  editForm.value.properties = {
    manufacturer: currentProps.manufacturer || '',
    model: currentProps.model || '',
    serial_number: currentProps.serial_number || '',
    form_factor: currentProps.form_factor || '',
    cpu_model: currentProps.cpu_model || '',
    cpu_sockets: currentProps.cpu_sockets || '',
    cpu_cores: currentProps.cpu_cores || '',
    ram_capacity_gb: currentProps.ram_capacity_gb || '',
    ram_type: currentProps.ram_type || '',
    storage_subsystem: currentProps.storage_subsystem || '',
    storage_controller: currentProps.storage_controller || '',
    pcie_cards: currentProps.pcie_cards || '',
    network_interfaces: currentProps.network_interfaces || '',
    os_name: currentProps.os_name || '',
    os_version: currentProps.os_version || '',
    firmware: currentProps.firmware || '',
    management_ip: currentProps.management_ip || '',
    mac_address: currentProps.mac_address || '',
    network_role: currentProps.network_role || '',
    port_config: currentProps.port_config || '',
    throughput: currentProps.throughput || '',
    engine: currentProps.engine || '',
    version: currentProps.version || '',
    edition: currentProps.edition || '',
    storage_engine: currentProps.storage_engine || '',
    host_cluster: currentProps.host_cluster || '',
    environment: currentProps.environment || '',
    allocated_storage_gb: currentProps.allocated_storage_gb || '',
    ha_setup: currentProps.ha_setup || '',
    vendor: currentProps.vendor || '',
    criticality: currentProps.criticality || '',
    framework: currentProps.framework || '',
    frontend_tech: currentProps.frontend_tech || '',
    url: currentProps.url || '',
    auth_method: currentProps.auth_method || ''
  }
  
  isEditPropertiesOpen.value = true
}

const saveProperties = async () => {
  isSavingProperties.value = true
  try {
    await updateAsset(assetId, { properties: editForm.value.properties })
    isEditPropertiesOpen.value = false
    await refreshAsset()
    await refreshLogs()
    await refreshRelations()
  } catch (error) {
    console.error('Failed to update properties:', error)
    alert('Failed to update properties')
  } finally {
    isSavingProperties.value = false
  }
}

const updateMaintenanceStatus = async (id, status) => {
   try {
    await apiUpdateMaintenanceStatus(id, status)
    await refreshMaintenance()
  } catch (error) {
    console.error('Failed to update maintenance window:', error)
  }
}

const deleteMaintenanceWindow = async (id) => {
  if (!confirm('Are you sure you want to delete this maintenance window?')) return
  try {
    await apiDeleteMaintenanceWindow(id)
    await refreshMaintenance()
  } catch (error) {
    console.error('Failed to delete maintenance window:', error)
  }
}

const updateInterfaceStatus = async (id, status) => {
  try {
    await apiUpdateInterface(id, { status })
    await refreshInterfaces()
  } catch (error) {
    console.error('Failed to update interface status:', error)
  }
}

const deleteInterface = async (id) => {
  if (!confirm('Are you sure you want to delete this network interface?')) return
  try {
    await apiDeleteInterface(id)
    await refreshInterfaces()
  } catch (error) {
    console.error('Failed to delete network interface:', error)
  }
}

onMounted(() => {
  if (route.query.edit === 'true') {
    watch(asset, (newAsset) => {
      if (newAsset && !isEditPropertiesOpen.value) {
        openEditProperties()
      }
    }, { immediate: true })
  }
})
</script>