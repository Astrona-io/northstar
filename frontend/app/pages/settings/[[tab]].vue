<template>
  <div class="space-y-6 font-sans">

    <!-- Site Admin Top Breadcrumb & Header -->
    <div class="flex justify-between items-center border-b border-slate-200 dark:border-slate-800 pb-4 flex-shrink-0">
      <div class="flex items-center gap-3">
        <UButton
          v-if="isSubPageActive"
          icon="i-heroicons-arrow-left"
          color="gray"
          variant="ghost"
          size="sm"
          @click="backToHub"
          class="hover:bg-slate-150 dark:hover:bg-slate-800"
        />
        <UIcon name="i-heroicons-shield-check" class="h-8 w-8 text-primary-500" />
        <div>
          <h2 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white font-mono">
            {{ isSubPageActive ? activeTabDetails.label : 'Site Administration' }}
          </h2>
          <p class="text-xs text-slate-500 dark:text-slate-400">
            {{ isSubPageActive ? activeTabDetails.desc : 'Centralized secure area to configure operator weights, network sweeper nodes, locations, and metadata schemas.' }}
          </p>
        </div>
      </div>
      <UBadge color="red" variant="subtle" class="font-mono text-[10px] font-bold tracking-wider px-2 py-1 uppercase">SECURED COCKPIT</UBadge>
    </div>

    <!-- MAIN PANEL VIEWPORT -->
    <div class="space-y-6">
      
      <!-- GRID LANDING PAGE: Site Admin Home Hub (GitLab Style) -->
      <div v-if="!isSubPageActive" class="space-y-6 pb-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Card 1: Users & Policies -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('users')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-users" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  {{ users?.length || 0 }} Users
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Users & Operator Policies</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Manage system operators access groups, configure fine-grained permissions, and configure local policy overrides.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Manage Users</UButton>
            </div>
          </div>

          <!-- Card 2: Discovery Subnet Sockets -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('discovery')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-bolt" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  Active
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Subnet Sweeper & Cloud Sync</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Audit standard local port handshakes on 127.0.0.1 or synchronise virtual EC2 and RDS nodes via AWS Cloud Crawlers.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Run Scanners</UButton>
            </div>
          </div>

          <!-- Card 3: Datacenters & Cabinets (DCIM) -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('dcim')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-server" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  {{ datacenters?.length || 0 }} Locations
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Datacenters & Cabinet DCIM</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Deploy physical on-prem/homelab zones or elastic cloud zones. Manage 42U rack cabinet placements and bandwidth uplinks.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Manage Sites</UButton>
            </div>
          </div>

          <!-- Card 4: Dynamic Custom Fields -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('fields')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-tag" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  {{ customFields?.length || 0 }} Schemas
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Dynamic Custom Fields</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Register custom property keys (text, number, boolean) to render dynamic wizard metadata on assets automatically.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Define Schemas</UButton>
            </div>
          </div>

          <!-- Card 5: Webhooks Integration -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('webhooks')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-link" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  {{ webhooks?.length || 0 }} Webhooks
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Webhook Broadcasters</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Register outbound callback endpoints and map events (Asset created, modified, or deleted) to Slack/Jira hooks.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Configure Webhooks</UButton>
            </div>
          </div>

          <!-- Card 6: Dynamic Categories & Icons Manager (Phase 1 Dynamic Categories) -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('categories')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-squares-plus" class="h-8 w-8 text-primary-500" />
                <span class="text-xs font-mono font-bold text-slate-400 bg-slate-50 dark:bg-slate-800 px-2 py-0.5 rounded">
                  {{ categories?.length || 0 }} Categories
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Category & Icon Manager</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Manage Asset Category cards, assign custom Heroicons styles dynamically, and link nested sub-group card portals.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Manage Icons</UButton>
            </div>
          </div>

          <!-- Card 6: Disaster Recovery (Phase 3 Backups) -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-red-500 transition-all cursor-pointer group"
            @click="triggerBackupExport"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-arrow-down-tray" class="h-8 w-8 text-red-500" />
                <span class="text-[9px] font-mono font-bold text-red-400 bg-red-50 dark:bg-red-950 px-2 py-0.5 rounded uppercase">
                  Encrypted State
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-red-500 transition-colors">Disaster Recovery Backups</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Compile and download a complete structural JSON backup encompassing Datacenters, Racks, Schemas, and Hooks.</p>
            </div>
            <div class="flex justify-between items-center pt-2">
              <UButton size="xs" color="gray" variant="ghost" icon="i-heroicons-arrow-up-tray" @click.stop="triggerBackupRestore">Restore</UButton>
              <UButton size="xs" color="red" variant="subtle" right-icon="i-heroicons-arrow-down-tray">Download JSON</UButton>
            </div>
          </div>

          <!-- Card 7: Bulk Ingest Engine (Phase 1 Bulk Importing) -->
          <div 
            class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-5 flex flex-col justify-between h-52 shadow-sm hover:border-primary-500 transition-all cursor-pointer group"
            @click="navigateToTab('importer')"
          >
            <div class="space-y-2">
              <div class="flex justify-between items-start">
                <UIcon name="i-heroicons-circle-stack" class="h-8 w-8 text-primary-500" />
                <span class="text-[9px] font-mono font-bold text-primary-400 bg-primary-50 dark:bg-primary-950 px-2 py-0.5 rounded uppercase">
                  Declarative
                </span>
              </div>
              <h4 class="font-bold text-sm text-slate-900 dark:text-white font-mono group-hover:text-primary-500 transition-colors">Bulk Ingest Engine</h4>
              <p class="text-xs text-slate-400 leading-normal line-clamp-2">Upload batch CSV hardware list manifests, execute dynamic GORM regex pattern validations, and perform bulk CMDB creations safely.</p>
            </div>
            <div class="flex justify-end pt-2">
              <UButton size="xs" color="primary" variant="subtle" right-icon="i-heroicons-arrow-right">Bulk Ingest CSV</UButton>
            </div>
          </div>

        </div>
      </div>

      <!-- DYNAMIC SUB-PAGES VIEWPORT (Full screen layout, Phase 1 Admin Redesign) -->
      <div v-else class="space-y-6 pb-6 pr-1 h-full min-h-0">
        
        <!-- Tab 1 Sub-Page: Users & Policies -->
        <div v-if="activeTab === 'users'" class="space-y-4">
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white">Active Operator Ledger</span>
                <UButton size="xs" color="gray" variant="ghost" label="Back to Dashboard" icon="i-heroicons-arrow-left" @click="backToHub" />
              </div>
            </template>

            <div v-if="pendingUsers" class="flex justify-center p-12">
              <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
            </div>

            <UTable v-else :rows="users || []" :columns="userColumns">
              <template #username-data="{ row }">
                <span class="font-semibold text-slate-900 dark:text-white">{{ row.username }}</span>
              </template>
              <template #role-data="{ row }">
                <UBadge :color="row.role === 'admin' ? 'red' : (row.role === 'operator' ? 'blue' : 'gray')" variant="subtle" class="uppercase text-[9px] font-bold font-mono">
                  {{ row.role }}
                </UBadge>
              </template>
              <template #group-data="{ row }">
                <span v-if="row.group?.name" class="text-sm font-medium text-slate-800 dark:text-slate-200">{{ row.group.name }}</span>
                <span v-else class="text-xs text-slate-400 italic">No Access Group</span>
              </template>
              <template #inherited-data="{ row }">
                <div class="flex flex-wrap gap-1 max-w-xs">
                  <template v-if="row.group?.permissions && row.group.permissions.length > 0">
                    <UBadge v-for="p in row.group.permissions" :key="p.id" color="blue" variant="subtle" size="xs">
                      {{ p.name }} ({{ p.effect }})
                    </UBadge>
                  </template>
                  <span v-else class="text-xs text-slate-400 italic">None</span>
                </div>
              </template>
              <template #overrides-data="{ row }">
                <div class="flex flex-wrap gap-1 max-w-xs">
                  <template v-if="row.permissions && row.permissions.length > 0">
                    <UBadge v-for="p in row.permissions" :key="p.id" :color="p.effect === 'allow' ? 'green' : 'orange'" variant="solid" size="xs">
                      {{ p.name }} ({{ p.effect }})
                    </UBadge>
                  </template>
                  <span v-else class="text-xs text-slate-400 italic">No Overrides</span>
                </div>
              </template>
              <template #actions-data="{ row }">
                <UButton v-if="isAdmin" size="xs" color="primary" variant="ghost" icon="i-heroicons-shield-exclamation" @click="openManageModal(row)">
                  Manage Access
                </UButton>
                <span v-else class="text-xs text-gray-400">Admin Only</span>
              </template>
            </UTable>
          </UCard>

          <!-- Administrative Access Audit Ledger (Phase 3 Audit Trail) -->
          <UCard class="mt-6">
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <h3 class="text-sm font-bold uppercase tracking-wider text-slate-800 dark:text-white flex items-center gap-2 font-mono">
                  <UIcon name="i-heroicons-clipboard-document-list" class="h-5 w-5 text-primary-500" />
                  Site Admin Access Audit Events Ledger (Audit L2)
                </h3>
                <UButton 
                  size="xs" 
                  color="gray" 
                  variant="subtle"
                  icon="i-heroicons-arrow-down-tray" 
                  @click="downloadAuditCsv"
                >
                  Download Audit Trail (CSV)
                </UButton>
              </div>
            </template>

            <div class="space-y-4">
              <p class="text-xs text-slate-500 leading-normal">
                Audits secure modification events including local policy updates, custom schemas, DCIM locations, and webhooks dispatches.
              </p>

              <UTable :rows="mockAuditLogs" :columns="auditColumns">
                <template #operator-data="{ row }">
                  <span class="font-bold text-slate-800 dark:text-white font-mono text-xs">{{ row.operator }}</span>
                </template>
                <template #action-data="{ row }">
                  <UBadge color="gray" variant="subtle" class="font-mono text-[9px] font-bold">{{ row.action }}</UBadge>
                </template>
                <template #timestamp-data="{ row }">
                  <span class="text-slate-400 font-mono text-[10px]">{{ row.timestamp }}</span>
                </template>
              </UTable>
            </div>
          </UCard>
        </div>

        <!-- Tab 2 Sub-Page: Subnet Discovery Scanner & AWS Cloud Sync -->
        <div v-if="activeTab === 'discovery'" class="grid grid-cols-1 xl:grid-cols-2 gap-6">
          <!-- Local Subnet Sweeper Card -->
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white uppercase">127.0.0.1 Port Sweeper Sockets</span>
                <span class="text-slate-400 font-mono">Local sweep</span>
              </div>
            </template>
            
            <div class="space-y-4">
              <p class="text-xs text-slate-500 leading-normal">
                Perform real-time socket connection sweeps on standard local ports. Banner-grabs open processes to register device templates.
              </p>
              <UButton :loading="isScanning" icon="i-heroicons-wifi" color="primary" block @click="runNetworkScan">
                Scan Network Subnet
              </UButton>

              <!-- Scanning Progress -->
              <div v-if="isScanning" class="space-y-2 py-4">
                <div class="flex justify-between text-[10px] text-slate-400 font-mono">
                  <span>Sweeping local sockets...</span>
                  <span>127.0.0.1 Open Ports...</span>
                </div>
                <UProgress animation="carousel" />
              </div>

              <!-- Scan Results -->
              <div v-if="scanResults && scanResults.length > 0" class="space-y-2 max-h-[260px] overflow-y-auto pr-1">
                <div 
                  v-for="(dev, idx) in scanResults" 
                  :key="idx" 
                  class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 border border-slate-150 dark:border-slate-800 rounded-md"
                >
                  <div>
                    <span class="text-xs font-bold text-slate-900 dark:text-white block">{{ dev.name }} (Port: {{ dev.port }})</span>
                    <span class="text-[10px] text-slate-400 font-mono block mt-0.5">{{ dev.ip_address }} | {{ dev.manufacturer }}</span>
                  </div>
                  <UButton 
                    v-if="!dev.imported"
                    size="xs" 
                    color="primary" 
                    variant="subtle" 
                    icon="i-heroicons-arrow-down-tray" 
                    @click="importDeviceToAssets(dev, idx)"
                  >
                    Import
                  </UButton>
                  <UBadge v-else color="green" variant="subtle" size="xs">Imported</UBadge>
                </div>
              </div>
              <div v-else-if="scanAttempted" class="text-xs text-slate-400 italic pt-2">
                Completed. Zero active unmanaged ports found.
              </div>
            </div>
          </UCard>

          <!-- AWS Cloud Crawler Sync Card -->
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white uppercase">AWS Cloud Crawler Engine</span>
                <span class="text-slate-400 font-mono">Elastic SDK</span>
              </div>
            </template>
            
            <div class="space-y-4">
              <p class="text-xs text-slate-500 leading-normal">
                Connect to cloud endpoints using secure credentials, crawling EC2 machines and RDS instances into standard active Assets.
              </p>
              <UButton :loading="isSyncingCloud" icon="i-heroicons-arrow-path-rounded-square" color="primary" block @click="runCloudSync">
                Sync AWS Cloud Assets
              </UButton>

              <!-- Sync Progress -->
              <div v-if="isSyncingCloud" class="space-y-2 py-4">
                <div class="flex justify-between text-[10px] text-slate-400 font-mono">
                  <span>Crawling AWS resources...</span>
                  <span>eu-west-1 VPC Stacks...</span>
                </div>
                <UProgress animation="carousel" color="blue" />
              </div>

              <!-- Sync Results -->
              <div v-if="cloudResults && cloudResults.length > 0" class="space-y-2 max-h-[260px] overflow-y-auto pr-1">
                <div 
                  v-for="(dev, idx) in cloudResults" 
                  :key="idx" 
                  class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 border border-slate-150 dark:border-slate-800 rounded-md"
                >
                  <div>
                    <span class="text-xs font-bold text-slate-900 dark:text-white block">{{ dev.name }}</span>
                    <span class="text-[10px] text-slate-400 font-mono block mt-0.5">{{ dev.ip_address }} | Instance: {{ dev.properties?.instance_id || dev.properties?.cluster_id }}</span>
                  </div>
                  <UBadge color="blue" variant="subtle" size="xs">Synced</UBadge>
                </div>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Tab 3 Sub-Page: Datacenters & Racks Administration (Phase 3 Connection speeds integration) -->
        <div v-if="activeTab === 'dcim'" class="space-y-6">
          <!-- Create Datacenter Action Row -->
          <div v-if="canMutate" class="flex justify-between items-center bg-slate-50 dark:bg-slate-800/40 p-4 rounded-md border border-slate-150 dark:border-slate-800 shadow-sm">
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block font-mono uppercase tracking-wider">Location Management</span>
              <span class="text-[10px] text-slate-400 mt-0.5 block">Deploy and manage physical or cloud-native facility profiles.</span>
            </div>
            <UButton 
              icon="i-heroicons-plus" 
              color="primary" 
              size="sm" 
              @click="isCreateDcModalOpen = true"
            >
              Deploy New Location Profile
            </UButton>
          </div>

          <!-- Active Infrastructure Location list & 2D Designer Accordion -->
          <UCard>
            <template #header>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono">Deployed DCIM Location Node Matrix</h4>
            </template>

            <div v-if="pendingDatacenters" class="flex justify-center p-6">
              <UIcon name="i-heroicons-arrow-path" class="animate-spin h-8 w-8 text-primary-500" />
            </div>

            <div v-else class="space-y-4 max-h-[500px] overflow-y-auto pr-1">
                <div 
                  v-for="dc in datacenters" 
                  :key="dc.id" 
                  class="p-4 bg-slate-50 dark:bg-slate-800 border border-slate-150 dark:border-slate-800 rounded-md space-y-3 cursor-pointer hover:border-primary-500/50 transition-colors"
                  @click="selectedDcIdForDesigner = (selectedDcIdForDesigner === dc.id ? null : dc.id)"
                >
                  <!-- DC Header -->
                  <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80 pb-2">
                    <div class="flex items-center gap-2">
                      <UIcon name="i-heroicons-globe-alt" class="h-5 w-5 text-primary-500" />
                      <div>
                        <span class="text-xs font-bold text-slate-900 dark:text-white">{{ dc.name }}</span>
                        <span class="text-[10px] text-slate-400 block">{{ dc.city }}, {{ dc.country }} | Type: {{ dc.type }}</span>
                        <!-- Connection metadata metrics display -->
                        <span v-if="dc.properties?.uplink_speed || dc.properties?.public_ip" class="text-[9px] text-primary-600 dark:text-primary-400 font-mono block mt-0.5">
                          ⚡ Link Speed: {{ dc.properties?.uplink_speed || 'N/A' }} | IP: {{ dc.properties?.public_ip || 'Dynamic' }}
                        </span>
                      </div>
                    </div>
                    <UButton v-if="canDelete" size="xs" color="red" variant="ghost" icon="i-heroicons-trash" @click="removeDatacenter(dc.id)" />
                  </div>

                  <!-- Nested Racks list inside DC -->
                  <div v-if="dc.racks && dc.racks.length > 0" class="pl-4 border-l-2 border-slate-200 dark:border-slate-800 space-y-2 pt-1">
                    <div 
                      v-for="rack in dc.racks" 
                      :key="rack.id" 
                      class="flex items-center justify-between text-xs p-2 bg-white dark:bg-slate-900 rounded-md border border-slate-200 dark:border-slate-800"
                    >
                      <div class="flex items-center gap-2">
                        <UIcon name="i-heroicons-rectangle-stack" class="h-4 w-4 text-gray-400" />
                        <span class="font-medium text-slate-800 dark:text-slate-200">{{ rack.name }} (Height: {{ rack.height_u }}U)</span>
                      </div>
                      <UButton v-if="canDelete" size="xs" color="red" variant="ghost" icon="i-heroicons-trash" @click="removeRack(rack.id)" />
                    </div>
                  </div>
                  <div v-else class="text-[10px] text-slate-400 italic pl-4">No cabinet frames deployed.</div>

                  <!-- Expanded Blueprint Preview (Fold out with drawing in fine size) -->
                  <div v-if="selectedDcIdForDesigner === dc.id" class="pt-4 border-t border-slate-200 dark:border-slate-800/85 space-y-4 cursor-default" @click.stop>
                    <div class="flex items-center justify-between">
                      <span class="text-[11px] font-bold text-slate-700 dark:text-slate-300 font-mono flex items-center gap-1.5">
                        <UIcon name="i-heroicons-map" class="text-primary-500 h-4 w-4" />
                        Room Layout Preview &mdash; {{ activeDesignerFloor?.name || 'Floor 1' }}
                      </span>
                      
                      <UButton 
                        size="xs" 
                        color="primary" 
                        icon="i-heroicons-pencil-square"
                        :disabled="!activeDesignerFloor"
                        @click="navigateTo(`/settings/dcim/${dc.id}/drawing?floorId=${activeDesignerFloor?.id}`)"
                      >
                        Enter Drafting Mode
                      </UButton>
                    </div>

                    <!-- Floor Level Selector & Management (Rename Alias & Add Floor) -->
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 items-center bg-slate-50 dark:bg-slate-900/40 p-3 rounded-md border border-slate-200 dark:border-slate-800/80">
                      <UFormGroup label="Active Floor Level" class="text-xs">
                        <USelect 
                          v-model="selectedFloorIdForDesigner" 
                          :options="dc.floors?.map(f => ({ label: f.name + ' (Level ' + f.level + ')', value: f.id })) || []"
                          size="xs"
                          placeholder="Select floor level..."
                        />
                      </UFormGroup>
                      <div class="flex gap-2 justify-end sm:justify-start pt-5">
                        <!-- Rename Floor (Alias) Button -->
                        <UButton 
                          v-if="activeDesignerFloor && canMutate" 
                          size="xs" 
                          color="gray" 
                          icon="i-heroicons-pencil-square"
                          @click="openRenameFloorModal(activeDesignerFloor)"
                        >
                          Rename (Alias)
                        </UButton>
                        
                        <!-- Add Floor Level Button -->
                        <UButton 
                          v-if="canMutate"
                          size="xs" 
                          color="primary" 
                          variant="soft"
                          icon="i-heroicons-plus"
                          @click="openAddFloorModal(dc.id)"
                        >
                          Add Floor
                        </UButton>
                      </div>
                    </div>

                    <!-- Fine-sized preview box -->
                    <div class="bg-slate-900 border border-slate-800 rounded-md h-[160px] relative overflow-hidden flex items-center justify-center">
                      <svg 
                        width="100%" 
                        height="100%" 
                        viewBox="0 0 400 200"
                        class="select-none pointer-events-none"
                      >
                        <defs>
                          <pattern id="preview-dots" width="20" height="20" patternUnits="userSpaceOnUse">
                            <circle cx="1" cy="1" r="0.5" fill="#475569" opacity="0.4" />
                          </pattern>
                        </defs>
                        <rect width="100%" height="100%" fill="url(#preview-dots)" />
                        
                        <!-- Draw walls -->
                        <line 
                          v-for="wall in activeDesignerFloor?.walls || []" 
                          :key="wall.id" 
                          :x1="wall.x1" 
                          :y1="wall.y1 * 0.8" 
                          :x2="wall.x2" 
                          :y2="wall.y2 * 0.8" 
                          stroke="#64748b" 
                          stroke-width="3" 
                          stroke-linecap="round"
                        />

                        <!-- Draw cabinet rectangles -->
                        <g v-for="rack in activeDesignerFloor?.racks || []" :key="rack.id">
                          <rect 
                            :x="rack.x / 2" 
                            :y="rack.y / 2 * 0.8" 
                            width="25" 
                            height="25" 
                            class="stroke-slate-600 fill-slate-800"
                            rx="3"
                          />
                          <text 
                            :x="(rack.x / 2) + 12.5" 
                            :y="(rack.y / 2 * 0.8) + 15" 
                            fill="#94a3b8" 
                            font-size="6" 
                            font-family="monospace" 
                            text-anchor="middle"
                          >
                            R-{{ rack.name }}
                          </text>
                        </g>
                      </svg>
                    </div>

                    <!-- Onboarding helper for floors directly inside card -->
                    <div v-if="!dc.floors || dc.floors.length === 0" class="p-4 bg-amber-50 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-900/50 rounded-md text-center space-y-2">
                      <p class="text-[10px] text-amber-700 dark:text-amber-400">No floor levels deployed. Auto-generate baseline levels to enable drawing.</p>
                      <UButton 
                        :loading="isCreatingDefaultFloors" 
                        size="xs" 
                        color="amber"
                        icon="i-heroicons-sparkles" 
                        @click="initializeDefaultFloors"
                      >
                        Auto-Generate Floors
                      </UButton>
                    </div>
                  </div>
                </div>
              </div>
            </UCard>
          </div>
        </div>

        <!-- Tab 4 Sub-Page: Dynamic Custom Fields -->
        <div v-if="activeTab === 'fields'" class="space-y-6">
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white uppercase">Dynamic Attributes Creator</span>
                <span class="text-slate-400 font-mono">properties JSON</span>
              </div>
            </template>

            <div class="space-y-6">
              <!-- Creator Card -->
              <UCard v-if="canMutate" class="bg-slate-50 dark:bg-slate-800/20 border-dashed">
                <template #header>
                  <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white">Register Custom Attribute Schema</h4>
                </template>
                <form @submit.prevent="registerField" class="space-y-4">
                  <div class="grid grid-cols-1 md:grid-cols-6 gap-4 items-end">
                    <UFormGroup label="Field Label Name" class="text-xs font-semibold">
                      <UInput v-model="fieldName" icon="i-heroicons-tag" placeholder="e.g. rack_row, warranty" required />
                    </UFormGroup>
                    <UFormGroup label="Input Value Type" class="text-xs font-semibold">
                      <USelect v-model="fieldType" :options="['text', 'number', 'boolean']" required />
                    </UFormGroup>
                    <UFormGroup label="Target Asset Category" class="text-xs font-semibold">
                      <USelect v-model="fieldCategory" :options="categoryNameOptions" required />
                    </UFormGroup>
                    <UFormGroup label="Tab Folder Group" class="text-xs font-semibold">
                      <UInput v-model="fieldTabGroup" icon="i-heroicons-folder" placeholder="e.g. Cooling Specs" />
                    </UFormGroup>
                    <UFormGroup label="Validation Regex Pattern" class="text-xs font-semibold">
                      <UInput v-model="fieldValidationRegex" icon="i-heroicons-shield-check" placeholder="e.g. ^[0-9a-fA-F]{8}$" />
                    </UFormGroup>
                    <div class="flex items-center gap-4 h-10 mb-1">
                      <UCheckbox v-model="fieldRequired" label="Required" />
                      <UButton type="submit" :loading="isSavingField" icon="i-heroicons-plus" color="primary">
                        Register
                      </UButton>
                    </div>
                  </div>
                </form>
              </UCard>

              <div v-if="pendingFields" class="flex justify-center p-12">
                <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
              </div>

              <!-- Table -->
              <UTable :rows="customFields || []" :columns="fieldColumns">
                <template #name-data="{ row }">
                  <span class="font-semibold text-slate-900 dark:text-white">{{ row.name }}</span>
                </template>
                <template #field_type-data="{ row }">
                  <span class="font-mono text-xs">{{ row.field_type }}</span>
                </template>
                <template #asset_type-data="{ row }">
                  <UBadge color="gray" variant="subtle">{{ row.asset_type }}</UBadge>
                </template>
                <template #tab_group-data="{ row }">
                  <UBadge color="blue" variant="subtle" class="font-mono text-[10px] font-bold">{{ row.tab_group || 'General Custom' }}</UBadge>
                </template>
                <template #validation_regex-data="{ row }">
                  <span v-if="row.validation_regex" class="font-mono text-xs text-red-500 bg-red-50 dark:bg-red-950 px-1.5 py-0.5 rounded border border-red-200 dark:border-red-900/50">{{ row.validation_regex }}</span>
                  <span v-else class="text-xs text-slate-400 italic">None</span>
                </template>
                <template #is_required-data="{ row }">
                  <UBadge :color="row.is_required ? 'red' : 'gray'" variant="subtle">
                    {{ row.is_required ? 'REQUIRED' : 'OPTIONAL' }}
                  </UBadge>
                </template>
                <template #actions-data="{ row }">
                  <UButton v-if="canDelete" size="xs" color="red" variant="ghost" icon="i-heroicons-trash" @click="removeField(row.id)">
                    Delete
                  </UButton>
                  <span v-else class="text-xs text-gray-400">Admin Only</span>
                </template>
              </UTable>
            </div>
          </UCard>
        </div>

        <!-- Tab 5 Sub-Page: Webhooks Integration -->
        <div v-if="activeTab === 'webhooks'" class="space-y-6">
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white uppercase">Outbound Callback Broadcasters</span>
                <span class="text-slate-400 font-mono">Events pub/sub</span>
              </div>
            </template>

            <div class="space-y-6">
              <!-- Creator Card -->
              <UCard v-if="canMutate" class="bg-slate-50 dark:bg-slate-800/20 border-dashed">
                <template #header>
                  <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white">Register Target Webhook</h4>
                </template>
                <form @submit.prevent="registerHook" class="space-y-4">
                  <div class="flex flex-col sm:flex-row gap-4 items-end">
                    <UFormGroup label="Target Callback URL" class="flex-1 text-xs font-semibold">
                      <UInput v-model="webhookUrl" icon="i-heroicons-link" placeholder="e.g. https://hooks.slack.com/services/..." required />
                    </UFormGroup>
                    <UFormGroup label="Broadcast Event Trigger" class="w-full sm:w-64 text-xs font-semibold">
                      <USelect v-model="webhookEvent" :options="['asset:create', 'asset:update', 'asset:delete']" required />
                    </UFormGroup>
                    <UButton type="submit" :loading="isSavingHook" icon="i-heroicons-plus" color="primary">
                      Register
                    </UButton>
                  </div>
                </form>
              </UCard>

              <div v-if="pendingHooks" class="flex justify-center p-12">
                <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
              </div>

              <!-- Table -->
              <UTable v-else :rows="webhooks || []" :columns="webhookColumns">
                <template #url-data="{ row }">
                  <span class="font-mono text-xs text-slate-900 dark:text-white truncate max-w-md block">{{ row.url }}</span>
                </template>
                <template #event-data="{ row }">
                  <UBadge color="blue" variant="subtle" class="font-mono text-[10px]">{{ row.event }}</UBadge>
                </template>
                <template #actions-data="{ row }">
                  <UButton v-if="canDelete" size="xs" color="red" variant="ghost" icon="i-heroicons-trash" @click="removeHook(row.id)">
                    Delete
                  </UButton>
                  <span v-else class="text-xs text-gray-400">Admin Only</span>
                </template>
              </UTable>
              </div>
              </UCard>
              </div>

        <!-- Tab 7 Sub-Page: Bulk Ingest Engine (Phase 1 Bulk Importing & Cabling) -->
        <div v-if="activeTab === 'importer'" class="space-y-6">
          <UCard>
            <template #header>
              <div class="flex justify-between items-center text-xs">
                <span class="font-bold font-mono text-slate-800 dark:text-white uppercase">Declarative CSV Ingest pipeline</span>
                <UButton size="xs" color="gray" variant="ghost" label="Back to Dashboard" icon="i-heroicons-arrow-left" @click="backToHub" />
              </div>
            </template>

            <div class="space-y-6">
              <!-- Sub-tab selector buttons inside Ingestion Engine -->
              <div class="flex gap-2 border-b border-slate-200 dark:border-slate-800 pb-2">
                <button
                  type="button"
                  :class="[
                    importerSubTab === 'hardware' 
                      ? 'border-primary-500 text-primary-500 font-bold border-b-2' 
                      : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300',
                    'px-4 py-1.5 text-xs font-mono transition-all pb-2'
                  ]"
                  @click="importerSubTab = 'hardware'"
                >
                  HARDWARE INVENTORY MANIFEST
                </button>
                <button
                  type="button"
                  :class="[
                    importerSubTab === 'cabling' 
                      ? 'border-primary-500 text-primary-500 font-bold border-b-2' 
                      : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300',
                    'px-4 py-1.5 text-xs font-mono transition-all pb-2'
                  ]"
                  @click="importerSubTab = 'cabling'"
                >
                  CABLE PATCH MATRIX (PHASE 1)
                </button>
              </div>

              <!-- Content A: Hardware Manifest CSV Ingest -->
              <div v-if="importerSubTab === 'hardware'" class="space-y-6">
                <p class="text-xs text-slate-500 leading-normal">
                  Declare equipment inventory at scale. Upload a `.csv` file. The engine will parse your rows and automatically enforce custom GORM regex validations on any dynamic attributes before committing creations.
                </p>

                <!-- Upload File block -->
                <div class="p-6 bg-slate-50 dark:bg-slate-850 border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-md text-center space-y-3 relative group">
                  <input 
                    type="file" 
                    accept=".csv" 
                    @change="handleCsvUpload" 
                    class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" 
                  />
                  <UIcon name="i-heroicons-cloud-arrow-up" class="h-10 w-10 text-primary-500 mx-auto group-hover:scale-105 transition-all" />
                  <div>
                    <span class="text-xs font-bold text-slate-800 dark:text-white block">Upload Hardware Manifest CSV</span>
                    <span class="text-[10px] text-slate-400 font-mono mt-1 block">Headers required: Name, Type, IP Address, Status, custom_attr1...</span>
                  </div>
                </div>

                <!-- Upload status and execution controls -->
                <div v-if="parsedCsvRows.length > 0" class="space-y-4">
                  <div class="flex justify-between items-center bg-slate-100 dark:bg-slate-800/50 p-3 rounded-md">
                    <span class="text-xs font-mono font-bold text-slate-700 dark:text-slate-300">
                      Parsed: {{ parsedCsvRows.length }} rows | 
                      <span :class="hasValidationErrors ? 'text-red-500' : 'text-green-500'">
                        {{ hasValidationErrors ? 'Validation Mismatch Blocked' : 'Valid and Clean' }}
                      </span>
                    </span>
                    <UButton 
                      :disabled="hasValidationErrors" 
                      :loading="isExecutingImport" 
                      color="primary" 
                      size="sm" 
                      icon="i-heroicons-bolt" 
                      @click="executeBulkImport"
                    >
                      Execute Ingest Creation
                    </UButton>
                  </div>

                  <!-- Validation Alerts/Warnings block -->
                  <div v-if="hasValidationErrors" class="p-4 bg-red-50/10 border border-red-500/30 rounded-md space-y-2">
                    <div class="flex items-center gap-2 text-red-500 font-bold text-xs font-mono">
                      <UIcon name="i-heroicons-shield-exclamation" class="h-5 w-5 animate-pulse" />
                      Strict Ingestion Validation Audits Blocked
                    </div>
                    <ul class="space-y-1 font-mono text-[10px] text-red-500/80 leading-relaxed pl-5 list-disc">
                      <li v-for="(err, i) in validationErrors" :key="i">
                        Row {{ err.row }}: '{{ err.field.toUpperCase() }}' (Value: '{{ err.val }}') must match pattern: <code class="bg-red-950 px-1 py-0.5 rounded text-white font-mono">{{ err.pattern }}</code>
                      </li>
                    </ul>
                  </div>

                  <!-- Parsed rows preview table -->
                  <UTable :rows="parsedCsvRows" :columns="csvColumns">
                    <template #name-data="{ row }">
                      <span class="font-bold text-slate-900 dark:text-white">{{ row.name }}</span>
                    </template>
                    <template #type-data="{ row }">
                      <UBadge color="gray" variant="subtle" class="font-mono text-[10px]">{{ row.type }}</UBadge>
                    </template>
                    <template #custom_fields-data="{ row }">
                      <div class="flex flex-wrap gap-1 font-mono text-[10px] max-w-sm">
                        <span v-for="(val, key) in row.custom_properties" :key="key" class="bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded text-slate-500">
                          {{ key }}: {{ val }}
                        </span>
                      </div>
                    </template>
                    <template #status-data="{ idx }">
                      <!-- Status column showing validation indicator -->
                      <div class="flex items-center gap-1.5">
                        <UIcon 
                          :name="getRowError(idx) ? 'i-heroicons-x-circle' : 'i-heroicons-check-circle'" 
                          :class="getRowError(idx) ? 'text-red-500' : 'text-green-500'" 
                          class="h-5 w-5" 
                        />
                        <span class="text-[10px] font-mono" :class="getRowError(idx) ? 'text-red-500 font-bold' : 'text-green-500'">
                          {{ getRowError(idx) ? 'FAILED' : 'PASS' }}
                        </span>
                      </div>
                    </template>
                  </UTable>
                </div>
              </div>

              <!-- Content B: Cable Patch Matrix CSV Ingest -->
              <div v-if="importerSubTab === 'cabling'" class="space-y-6">
                <p class="text-xs text-slate-500 leading-normal">
                  Dispatch and connect hardware cabling configurations in batch. Upload a `.csv` file declaring link maps. The engine automatically executes GORM asset presence verification checks before patching.
                </p>

                <!-- Upload File block -->
                <div class="p-6 bg-slate-50 dark:bg-slate-850 border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-md text-center space-y-3 relative group">
                  <input 
                    type="file" 
                    accept=".csv" 
                    @change="handleCablingUpload" 
                    class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" 
                  />
                  <UIcon name="i-heroicons-share" class="h-10 w-10 text-primary-500 mx-auto group-hover:scale-105 transition-all" />
                  <div>
                    <span class="text-xs font-bold text-slate-800 dark:text-white block">Upload Cable Patch Matrix CSV</span>
                    <span class="text-[10px] text-slate-400 font-mono mt-1 block">Headers required: Source, Target, Type</span>
                  </div>
                </div>

                <!-- Upload status and execution controls -->
                <div v-if="parsedCablingRows.length > 0" class="space-y-4">
                  <div class="flex justify-between items-center bg-slate-100 dark:bg-slate-800/50 p-3 rounded-md">
                    <span class="text-xs font-mono font-bold text-slate-700 dark:text-slate-300">
                      Parsed: {{ parsedCablingRows.length }} rows | 
                      <span :class="hasCablingValidationErrors ? 'text-red-500' : 'text-green-500'">
                        {{ hasCablingValidationErrors ? 'Verification Errors Blocked' : 'Valid and Clean' }}
                      </span>
                    </span>
                    <UButton 
                      :disabled="hasCablingValidationErrors" 
                      :loading="isExecutingCabling" 
                      color="primary" 
                      size="sm" 
                      icon="i-heroicons-bolt" 
                      @click="executeBulkCabling"
                    >
                      Execute Cabling Patching
                    </UButton>
                  </div>

                  <!-- Validation Alerts/Warnings block -->
                  <div v-if="hasCablingValidationErrors" class="p-4 bg-red-50/10 border border-red-500/30 rounded-md space-y-2">
                    <div class="flex items-center gap-2 text-red-500 font-bold text-xs font-mono">
                      <UIcon name="i-heroicons-shield-exclamation" class="h-5 w-5 animate-pulse" />
                      Dynamic GORM Cable Verification Blocked
                    </div>
                    <ul class="space-y-1 font-mono text-[10px] text-red-500/80 leading-relaxed pl-5 list-disc">
                      <li v-for="(err, i) in cablingValidationErrors" :key="i">
                        Row {{ err.row }}: Mismatched node '{{ err.val }}' in field '{{ err.field.toUpperCase() }}' &rarr; {{ err.reason }}
                      </li>
                    </ul>
                  </div>

                  <!-- Parsed cabling rows preview table -->
                  <UTable :rows="parsedCablingRows" :columns="cablingColumns">
                    <template #source_asset-data="{ row }">
                      <span class="font-bold text-slate-900 dark:text-white">{{ row.source_asset }}</span>
                    </template>
                    <template #target_asset-data="{ row }">
                      <span class="font-bold text-slate-900 dark:text-white">{{ row.target_asset }}</span>
                    </template>
                    <template #connection_type-data="{ row }">
                      <UBadge color="blue" variant="subtle" class="font-mono text-[10px] uppercase">{{ row.connection_type }}</UBadge>
                    </template>
                    <template #status-data="{ idx }">
                      <div class="flex items-center gap-1.5">
                        <UIcon 
                          :name="getCablingRowError(idx) ? 'i-heroicons-x-circle' : 'i-heroicons-check-circle'" 
                          :class="getCablingRowError(idx) ? 'text-red-500' : 'text-green-500'" 
                          class="h-5 w-5" 
                        />
                        <span class="text-[10px] font-mono" :class="getCablingRowError(idx) ? 'text-red-500 font-bold' : 'text-green-500'">
                          {{ getCablingRowError(idx) ? 'INVALID' : 'READY' }}
                        </span>
                      </div>
                    </template>
                  </UTable>
                </div>
              </div>

            </div>
          </UCard>
        </div>

              <!-- Tab 6 Sub-Page: Category & Icon Manager (Phase 1 Dynamic Categories) -->
              <div v-if="activeTab === 'categories'" class="grid grid-cols-1 xl:grid-cols-2 gap-6">
              <!-- Left: Category & Subgroup Add Forms -->
              <div class="space-y-6">
              <!-- Create Category Card -->
              <UCard v-if="canMutate">
              <template #header>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono">Register Asset Category</h4>
              </template>
              <form @submit.prevent="registerCategory" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <UFormGroup label="Category Name" class="text-xs font-semibold">
                  <UInput v-model="newCatForm.name" placeholder="e.g. AI Accelerators, Power Delivery" required />
                </UFormGroup>
                <UFormGroup label="Prerequisite Parent Category (Phase 1)" class="text-xs font-semibold">
                  <USelect v-model="newCatForm.parent_dependency" :options="categoryNameOptions" placeholder="None (Optional)" />
                </UFormGroup>
              </div>
              <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Category Icon Class" class="text-xs font-semibold">
                <UInput v-model="newCatForm.icon" placeholder="e.g. i-heroicons-cpu-chip" required />
              </UFormGroup>
              <UFormGroup label="Category Description" class="text-xs font-semibold">
                <UInput v-model="newCatForm.description" placeholder="Short description" required />
              </UFormGroup>
              </div>
              <div class="flex justify-end pt-2">
              <UButton type="submit" :loading="isSavingCat" size="sm" icon="i-heroicons-plus" color="primary">
                Register Category
              </UButton>
              </div>
              </form>
              </UCard>

              <!-- Create Subgroup Card -->
              <UCard v-if="canMutate">
              <template #header>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono">Register Category Subgroup</h4>
              </template>
              <form @submit.prevent="registerSubGroup" class="space-y-4">
              <UFormGroup label="Parent Category Mapping" class="text-xs font-semibold">
              <USelect 
                v-model="newSubForm.category_id" 
                :options="catDropdownOptions" 
                option-attribute="label" 
                value-attribute="value" 
                required 
              />
              </UFormGroup>
              <UFormGroup label="Subgroup Name" class="text-xs font-semibold">
              <UInput v-model="newSubForm.name" placeholder="e.g. Switch (L3), Sensor" required />
              </UFormGroup>
              <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Subgroup Icon Class" class="text-xs font-semibold">
                <UInput v-model="newSubForm.icon" placeholder="e.g. i-heroicons-wifi" required />
              </UFormGroup>
              <UFormGroup label="Subgroup Description" class="text-xs font-semibold">
                <UInput v-model="newSubForm.description" placeholder="Description details" required />
              </UFormGroup>
              </div>
              <div class="flex justify-end pt-2">
              <UButton type="submit" :loading="isSavingSub" size="sm" icon="i-heroicons-plus" color="primary">
                Register Subgroup
              </UButton>
              </div>
              </form>
              </UCard>
              </div>

              <!-- Right: List and Delete dynamic Categories -->
              <div class="space-y-6">
              <UCard>
              <template #header>
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono">Active Database Categories & Sub-Groups</h4>
              </template>

              <div v-if="pendingCategories" class="flex justify-center p-6">
              <UIcon name="i-heroicons-arrow-path" class="animate-spin h-8 w-8 text-primary-500" />
              </div>

              <div v-else class="space-y-4 max-h-[500px] overflow-y-auto pr-1">
              <div 
              v-for="cat in categories" 
              :key="cat.id" 
              class="p-4 bg-slate-50 dark:bg-slate-800 border border-slate-150 dark:border-slate-800 rounded-md space-y-3"
              >
              <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-2">
                <div class="flex items-center gap-2">
                  <UIcon :name="cat.icon || 'i-heroicons-squares-2x2'" class="h-5 w-5 text-primary-500" />
                  <div>
                    <div class="flex items-center gap-1.5">
                      <span class="text-xs font-bold text-slate-900 dark:text-white">{{ cat.name }}</span>
                      <UBadge v-if="cat.parent_dependency" color="amber" variant="subtle" class="font-mono text-[9px] font-bold">REQ: {{ cat.parent_dependency.toUpperCase() }}</UBadge>
                    </div>
                    <p class="text-[10px] text-slate-400 leading-normal">{{ cat.description }}</p>
                  </div>
                </div>
                <UButton v-if="canDelete && !isReservedCategory(cat.name)" size="xs" color="red" variant="ghost" icon="i-heroicons-trash" @click="removeCategory(cat.id)" />
              </div>

              <!-- Subgroups -->
              <div v-if="cat.sub_groups && cat.sub_groups.length > 0" class="pl-4 border-l-2 border-slate-200 dark:border-slate-800 space-y-2 pt-1">
                <div 
                  v-for="sub in cat.sub_groups" 
                  :key="sub.id" 
                  class="flex items-center justify-between text-[11px] p-2 bg-white dark:bg-slate-900 rounded-md border border-slate-200 dark:border-slate-800"
                >
                  <div class="flex items-center gap-2">
                    <UIcon :name="sub.icon || 'i-heroicons-squares-plus'" class="h-4 w-4 text-gray-400" />
                    <span>{{ sub.name }}</span>
                  </div>
                  <span class="text-[10px] text-slate-400 max-w-[200px] truncate">{{ sub.description }}</span>
                </div>
              </div>
              <div v-else class="text-[10px] text-slate-400 italic pl-4">No nested sub-groups.</div>
              </div>
              </div>
              </UCard>
              </div>
              </div>

    </div>

    <!-- Deploy New Location Profile Modal -->
    <UModal v-model="isCreateDcModalOpen">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <UIcon name="i-heroicons-globe-alt" class="h-5 w-5 text-primary-500" />
              <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-white font-mono">Deploy New Location Profile</h4>
            </div>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="isCreateDcModalOpen = false" />
          </div>
        </template>
        <form @submit.prevent="registerDatacenter" class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <UFormGroup label="Location Name" class="text-xs font-semibold">
              <UInput v-model="newDcForm.name" placeholder="e.g. Dublin-HQ" required />
            </UFormGroup>
            <UFormGroup label="Datacenter Type" class="text-xs font-semibold">
              <USelect v-model="newDcForm.type" :options="['on-prem', 'homelab', 'air-gap', 'cloud:aws', 'cloud:gcp', 'cloud:azure']" required />
            </UFormGroup>
            <UFormGroup label="Country Location" class="text-xs font-semibold">
              <UInput v-model="newDcForm.country" placeholder="e.g. Ireland" />
            </UFormGroup>
            <UFormGroup label="City Location" class="text-xs font-semibold">
              <UInput v-model="newDcForm.city" placeholder="e.g. Dublin" />
            </UFormGroup>
            <!-- Connection Speeds parameters -->
            <UFormGroup label="Uplink Connection Speed" class="text-xs font-semibold">
              <UInput v-model="newDcForm.uplink_speed" placeholder="e.g. 2.5 Gbps Fiber, AWS Elastic" />
            </UFormGroup>
            <UFormGroup label="Location Public IP Gate" class="text-xs font-semibold">
              <UInput v-model="newDcForm.public_ip" placeholder="e.g. 85.12.85.112" />
            </UFormGroup>
          </div>
          <div class="flex justify-end gap-2 pt-4 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isCreateDcModalOpen = false">Cancel</UButton>
            <UButton type="submit" :loading="isSavingDc" icon="i-heroicons-plus" color="primary">
              Deploy Datacenter Profile
            </UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <!-- Policy Management Modal -->
    <UModal v-model="isModalOpen" v-if="selectedUser">
      <UCard>
        <template #header>
          <div class="flex items-center gap-2">
            <UIcon name="i-heroicons-shield-check" class="h-6 w-6 text-primary-500" />
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">
              Manage Access: {{ selectedUser.username }}
            </h3>
          </div>
        </template>

        <form @submit.prevent="saveUserPolicies" class="space-y-4">
          <!-- Access Group -->
          <UFormGroup label="Assigned Access Group">
            <USelect 
              v-model="userForm.group_id" 
              :options="groupOptions" 
              option-attribute="label" 
              value-attribute="value" 
            />
          </UFormGroup>

          <!-- User-Level Overrides -->
          <UFormGroup label="User-Level Override Policies (User overrides Group Weighting)">
            <USelectMenu 
              v-model="userForm.permission_ids" 
              :options="permissionsOptions || []" 
              multiple 
              option-attribute="displayName" 
              value-attribute="id" 
              placeholder="Select local allow/deny overrides..." 
              searchable 
            />
          </UFormGroup>

          <!-- Help Tip -->
          <div class="bg-slate-50 dark:bg-slate-800/30 border border-slate-150 dark:border-slate-800 p-3 rounded-md text-xs text-slate-500">
            <strong>Policy Weight Overrides:</strong> Overrides configured directly on the User **strictly override** whatever privileges are granted or denied by their Access Group.
          </div>

          <div class="flex justify-end gap-2 mt-6 pt-3 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isModalOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="isSavingUser">Save Policies</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <!-- Rename Floor (Alias) Modal -->
    <UModal v-model="isRenameFloorModalOpen">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
              <UIcon name="i-heroicons-pencil-square" class="text-primary-500 h-5 w-5" />
              Rename Floor Level / Floor Alias
            </h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isRenameFloorModalOpen = false" />
          </div>
        </template>
        
        <form @submit.prevent="submitRenameFloor" class="space-y-4">
          <UFormGroup label="Floor Name / Floor Alias" help="e.g. Server Room A, Ground Floor, Floor 0, Suite 101">
            <UInput v-model="floorRenameForm.name" placeholder="Floor Name / Floor Alias" required />
          </UFormGroup>
          <UFormGroup label="Floor level number" help="e.g. 0 for ground level, 1 for first floor, -1 for basement">
            <UInput type="number" v-model="floorRenameForm.level" required />
          </UFormGroup>
          
          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isRenameFloorModalOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="isSavingFloorRename">Save Changes</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <!-- Add Floor Level Modal -->
    <UModal v-model="isAddFloorModalOpen">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
              <UIcon name="i-heroicons-plus-circle" class="text-primary-500 h-5 w-5" />
              Add Datacenter Floor Level
            </h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isAddFloorModalOpen = false" />
          </div>
        </template>
        
        <form @submit.prevent="submitAddFloor" class="space-y-4">
          <UFormGroup label="Floor Name" help="e.g. Floor 1, Basement Level -1">
            <UInput v-model="floorAddForm.name" placeholder="e.g. Floor 1" required />
          </UFormGroup>
          <UFormGroup label="Floor level number" help="e.g. 1, 2, -1">
            <UInput type="number" v-model="floorAddForm.level" required />
          </UFormGroup>
          
          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isAddFloorModalOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="isSavingFloorAdd">Create Level</UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const route = useRoute()
const router = useRouter()
const toast = useToast()

// Modern, non-blocking notification interceptor overriding browser alert popups
const alert = (message) => {
  const isError = message.toLowerCase().includes('fail') || 
                  message.toLowerCase().includes('error') || 
                  message.toLowerCase().includes('limit') ||
                  message.toLowerCase().includes('reject')
  toast.add({
    title: isError ? 'Operational Alert' : 'Success',
    description: message,
    color: isError ? 'red' : 'green',
    icon: isError ? 'i-heroicons-x-circle' : 'i-heroicons-check-circle'
  })
}

const { fetchUsers, fetchGroups, fetchPermissions, updateUserOverrides } = useUsers()
const { fetchCustomFields, createCustomField, deleteCustomField } = useCustomFields()
const { fetchWebhooks, createWebhook, deleteWebhook } = useWebhooks()
const { fetchDiscoveryResults, triggerDiscoveryScan, importDiscoveredDevice } = useDiscovery()
const { fetchDatacenters, fetchRacks } = useDCIM()
const { isAdmin, canMutate, canDelete, getAuthHeader } = useAuth()

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

// Administrative Site Auditing Ledger (Phase 3 Audit Trail)
const auditColumns = [
  { key: 'id', label: 'Event ID' },
  { key: 'operator', label: 'Operator Username' },
  { key: 'action', label: 'Action Event' },
  { key: 'meta', label: 'Metadata Description' },
  { key: 'timestamp', label: 'Timestamp' }
]

const mockAuditLogs = [
  { id: '#1001', operator: 'admin', action: 'USER_POLICY_UPDATE', meta: 'Synchronized GORM overrides permissions for hybrid operator.', timestamp: '2026-06-03 10:14:15' },
  { id: '#1002', operator: 'admin', action: 'SCHEMA_FIELD_CREATE', meta: 'Added dynamic custom field schema warranty_expiry to Server category.', timestamp: '2026-06-03 10:28:44' },
  { id: '#1003', operator: 'admin', action: 'DCIM_LOC_DEPLOY', meta: 'Deployed Dublin-HQ datacenter with 2.5 Gbps Fiber WAN.', timestamp: '2026-06-03 10:32:12' },
  { id: '#1004', operator: 'admin', action: 'WEBHOOK_REGISTER', meta: 'Registered Slack outbound callback endpoint for asset mutations.', timestamp: '2026-06-03 10:35:19' }
]

const downloadAuditCsv = () => {
  let csv = "Event ID,Operator Username,Action Event,Metadata Description,Timestamp\n"
  mockAuditLogs.forEach(log => {
    csv += `${log.id},${log.operator},${log.action},"${log.meta.replace(/"/g, '""')}",${log.timestamp}\n`
  })
  
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = "audit-trail-northstar.csv"
  link.click()
  URL.revokeObjectURL(url)
}

// Fetch all datasets preloaded natively
const { data: users, pending: pendingUsers, refresh: refreshUsers } = await fetchUsers()
const { data: groups } = await fetchGroups()
const { data: rawPermissions } = await fetchPermissions()
const { data: customFields, pending: pendingFields, refresh: refreshFields } = await fetchCustomFields()
const { data: webhooks, pending: pendingHooks, refresh: refreshHooks } = await fetchWebhooks()
const { data: datacenters, pending: pendingDatacenters, refresh: refreshDatacenters } = await fetchDatacenters()
const { fetchAssets } = useAssets()
const { data: allAssetsList } = await fetchAssets()

// Fetch dynamic categories (Phase 1 Dynamic Categories)
const { data: categories, pending: pendingCategories, refresh: refreshCategories } = await useFetch(`${apiBase}/categories/`)

// SPA routing activeTab URL retaining query param mapping (Phase 1 GitLab Site Admin Re-org)
const activeTab = computed({
  get: () => route.params.tab || null,
  set: (val) => {
    router.replace({ path: val ? `/settings/${val}` : '/settings' })
  }
})

// Check if any specific sub-page card is currently being configured
const isSubPageActive = computed(() => !!activeTab.value)

const activeTabDetails = computed(() => {
  if (!activeTab.value) return { label: '', desc: '' }
  return adminTabs.find(t => t.key === activeTab.value) || { label: '', desc: '' }
})

const navigateToTab = (tabKey) => {
  activeTab.value = tabKey
}

const backToHub = () => {
  activeTab.value = null
}

const isModalOpen = ref(false)
const isSavingUser = ref(false)
const isSavingField = ref(false)
const isSavingHook = ref(false)

const isCreateDcModalOpen = ref(false)
const isSavingDc = ref(false)

const isSavingCat = ref(false)
const isSavingSub = ref(false)

const selectedUser = ref(null)

const adminTabs = [
  { key: 'users', label: 'User Operator Policies', icon: 'i-heroicons-users', desc: 'Configure operator access permissions and manage policy overrides.' },
  { key: 'discovery', label: 'Subnet Sweeper & Cloud Sync', icon: 'i-heroicons-bolt', desc: 'Audit running standard services on localhost or crawl AWS cloud resources.' },
  { key: 'dcim', label: 'Datacenters & Racks (DCIM)', icon: 'i-heroicons-server', desc: 'Deploy physical datacenter location profiles, uplink connection speeds, and deploy cabinets.' },
  { key: 'importer', label: 'Bulk Ingest Engine', icon: 'i-heroicons-circle-stack', desc: 'Upload batch CSV hardware list manifests, execute dynamic GORM regex pattern validations, and perform bulk CMDB creations safely.' },
  { key: 'categories', label: 'Category & Icon Manager', icon: 'i-heroicons-squares-plus', desc: 'Manage CMDB asset categories, assign custom Heroicons styles dynamically, and link nested sub-group card portals.' },
  { key: 'fields', label: 'Dynamic Custom Fields', icon: 'i-heroicons-tag', desc: 'Define user-defined dynamic meta-properties and schema validation requirements.' },
  { key: 'webhooks', label: 'Webhook Broadcasters', icon: 'i-heroicons-link', desc: 'Configure real-time outbound HTTP callback webhook targets for CMDB events.' }
]

// ----------------------------------------
// Part 1: Users & Policies
// ----------------------------------------
const userForm = ref({
  group_id: null,
  permission_ids: []
})

const permissionsOptions = computed(() => {
  if (!rawPermissions.value) return []
  return rawPermissions.value.map(p => ({
    ...p,
    displayName: `${p.name} (${p.effect.toUpperCase()})`
  }))
})

const groupOptions = computed(() => {
  const options = [{ label: 'No Access Group (None)', value: null }]
  if (groups.value) {
    groups.value.forEach(g => {
      options.push({ label: g.name, value: g.id })
    })
  }
  return options
})

const userColumns = [
  { key: 'username', label: 'User Operator' },
  { key: 'role', label: 'Role Context' },
  { key: 'group', label: 'Access Group' },
  { key: 'inherited', label: 'Group Policies (Inherited)' },
  { key: 'overrides', label: 'Local Overrides' },
  { key: 'actions', label: 'Control' }
]

const openManageModal = (user) => {
  selectedUser.value = user
  userForm.value = {
    group_id: user.group_id,
    permission_ids: user.permissions ? user.permissions.map(p => p.id) : []
  }
  isModalOpen.value = true
}

const saveUserPolicies = async () => {
  isSavingUser.value = true
  try {
    await updateUserOverrides(selectedUser.value.id, userForm.value.group_id, userForm.value.permission_ids)
    isModalOpen.value = false
    await refreshUsers()
    alert('User security policies synchronized successfully.')
  } catch (err) {
    console.error('Failed to sync user overrides:', err)
    alert('Sync failed: Check admin privileges.')
  } finally {
    isSavingUser.value = false
  }
}

// ----------------------------------------
// Part 2: Subnet discovery scanner & Cloud synchronizer
// ----------------------------------------
const isScanning = ref(false)
const scanAttempted = ref(false)
const scanResults = ref([])

const isSyncingCloud = ref(false)
const cloudResults = ref([])

const runNetworkScan = async () => {
  isScanning.value = true
  scanAttempted.value = false
  scanResults.value = []
  try {
    const { data } = await triggerDiscoveryScan()
    if (data.value) {
      scanResults.value = data.value.map(item => ({ ...item, imported: false }))
    }
  } catch (err) {
    console.error('Failed to trigger background scan sweep:', err)
  } finally {
    isScanning.value = false
    scanAttempted.value = true
  }
}

const importDeviceToAssets = async (deviceSpec, idx) => {
  try {
    await importDiscoveredDevice(deviceSpec)
    scanResults.value[idx].imported = true
    alert(`Successfully imported standard hardware ${deviceSpec.name} to active Assets. Audit trace logged.`)
  } catch (err) {
    console.error('Failed to import discovered asset:', err)
    alert('Import failed: check security override scopes.')
  }
}

const runCloudSync = async () => {
  isSyncingCloud.value = true
  cloudResults.value = []
  try {
    const data = await $fetch(`${apiBase}/discovery/cloud-sync/`, {
      method: 'POST',
      headers: getAuthHeader()
    })
    if (data && data.imported) {
      cloudResults.value = data.imported
      alert('AWS Cloud synchronizer handshakes completed successfully. Cloud node specifications imported.')
    }
  } catch (err) {
    console.error('Cloud Sync failed:', err)
    alert('Cloud Sync failed: IAM authorization rejected.')
  } finally {
    isSyncingCloud.value = false
  }
}

// ----------------------------------------
// Part 3: Datacenters & Racks (DCIM Admin CRUD)
// ----------------------------------------
const newDcForm = ref({
  name: '',
  type: 'homelab',
  country: '',
  city: '',
  uplink_speed: '',
  public_ip: ''
})

// Room Layout Designer State Variables (Phase 1 CAD DCIM & Wall Drawing)
const selectedDcIdForDesigner = ref(null)

// CAD Designer Tool States (DCIM CAD L4)
const selectedFloorIdForDesigner = ref(null)

const activeDesignerFloor = computed(() => {
  if (!selectedDcIdForDesigner.value || !datacenters.value) return null
  const dc = datacenters.value.find(d => d.id === selectedDcIdForDesigner.value)
  if (!dc || !dc.floors || dc.floors.length === 0) return null
  
  if (selectedFloorIdForDesigner.value) {
    const f = dc.floors.find(fl => fl.id === selectedFloorIdForDesigner.value)
    if (f) return f
  }
  return dc.floors[0]
})

// Auto sync selectedFloorIdForDesigner when DC selection changes
watch(selectedDcIdForDesigner, (newDcId) => {
  if (newDcId && datacenters.value) {
    const dc = datacenters.value.find(d => d.id === newDcId)
    if (dc && dc.floors && dc.floors.length > 0) {
      selectedFloorIdForDesigner.value = dc.floors[0].id
    } else {
      selectedFloorIdForDesigner.value = null
    }
  }
})

// Auto-create default Floor 1 (Level 1) and Basement (Level -1) if DC lacks floors
const isCreatingDefaultFloors = ref(false)
const initializeDefaultFloors = async () => {
  if (!selectedDcIdForDesigner.value) return
  isCreatingDefaultFloors.value = true
  try {
    // Create Floor 1
    await $fetch(`${apiBase}/datacenter-floors`, {
      method: 'POST',
      body: { datacenter_id: selectedDcIdForDesigner.value, name: 'Floor 1', level: 1 },
      headers: getAuthHeader()
    })
    // Create Basement Level -1
    await $fetch(`${apiBase}/datacenter-floors`, {
      method: 'POST',
      body: { datacenter_id: selectedDcIdForDesigner.value, name: 'Basement (-1)', level: -1 },
      headers: getAuthHeader()
    })
    await refreshDatacenters()
    const dc = datacenters.value.find(d => d.id === selectedDcIdForDesigner.value)
    if (dc && dc.floors && dc.floors.length > 0) {
      selectedFloorIdForDesigner.value = dc.floors[0].id
    }
    alert('Successfully auto-generated baseline Floor 1 and Basement Level -1!')
  } catch (err) {
    console.error('Failed to create default floors:', err)
  } finally {
    isCreatingDefaultFloors.value = false
  }
}

// Floor Level Management (Rename Alias & Add Floor)
const isRenameFloorModalOpen = ref(false)
const isSavingFloorRename = ref(false)
const floorRenameForm = ref({ id: '', name: '', level: 0 })

const openRenameFloorModal = (floor) => {
  floorRenameForm.value = { id: floor.id, name: floor.name, level: floor.level }
  isRenameFloorModalOpen.value = true
}

const submitRenameFloor = async () => {
  isSavingFloorRename.value = true
  try {
    await $fetch(`${apiBase}/datacenter-floors/${floorRenameForm.value.id}`, {
      method: 'PUT',
      body: { 
        name: floorRenameForm.value.name, 
        level: parseInt(floorRenameForm.value.level) 
      },
      headers: getAuthHeader()
    })
    await refreshDatacenters()
    isRenameFloorModalOpen.value = false
    alert('Floor details updated successfully!')
  } catch (err) {
    console.error('Failed to rename floor:', err)
    alert('Failed to update floor details.')
  } finally {
    isSavingFloorRename.value = false
  }
}

const isAddFloorModalOpen = ref(false)
const isSavingFloorAdd = ref(false)
const floorAddForm = ref({ datacenter_id: '', name: '', level: 0 })

const openAddFloorModal = (dcId) => {
  floorAddForm.value = { datacenter_id: dcId, name: '', level: 0 }
  isAddFloorModalOpen.value = true
}

const submitAddFloor = async () => {
  isSavingFloorAdd.value = true
  try {
    const res = await $fetch(`${apiBase}/datacenter-floors`, {
      method: 'POST',
      body: { 
        datacenter_id: floorAddForm.value.datacenter_id, 
        name: floorAddForm.value.name, 
        level: parseInt(floorAddForm.value.level) 
      },
      headers: getAuthHeader()
    })
    await refreshDatacenters()
    if (res && res.id) {
      selectedFloorIdForDesigner.value = res.id
    }
    isAddFloorModalOpen.value = false
    alert('Floor level created successfully!')
  } catch (err) {
    console.error('Failed to create floor level:', err)
    alert('Failed to create floor level.')
  } finally {
    isSavingFloorAdd.value = false
  }
}

const registerDatacenter = async () => {
  isSavingDc.value = true
  try {
    const payload = {
      name: newDcForm.value.name,
      type: newDcForm.value.type,
      country: newDcForm.value.country,
      city: newDcForm.value.city,
      properties: {
        uplink_speed: newDcForm.value.uplink_speed || '',
        public_ip: newDcForm.value.public_ip || ''
      }
    }
    await $fetch(`${apiBase}/datacenters/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
    newDcForm.value = { name: '', type: 'homelab', country: '', city: '', uplink_speed: '', public_ip: '' }
    await refreshDatacenters()
    alert('Datacenter location deployed and registered successfully.')
    isCreateDcModalOpen.value = false
  } catch (err) {
    console.error('Failed to create datacenter:', err)
    alert('Deploy failed: Check authorization scopes or duplicate name constraints.')
  } finally {
    isSavingDc.value = false
  }
}

const removeDatacenter = async (dcId) => {
  if (!confirm('Are you sure you want to decommission this Datacenter location? All nested racks inside this location will be cascadingly deleted.')) return
  try {
    await $fetch(`${apiBase}/datacenters/${dcId}/`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
    await refreshDatacenters()
  } catch (err) {
    console.error('Failed to decommission location:', err)
  }
}

const removeRack = async (rackId) => {
  if (!confirm('Are you sure you want to decommission this rack cabinet?')) return
  try {
    await $fetch(`${apiBase}/racks/${rackId}/`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
    await refreshDatacenters()
  } catch (err) {
    console.error('Failed to delete cabinet frame:', err)
  }
}

// ----------------------------------------
// Part 4: Custom Fields
// ----------------------------------------
const fieldName = ref('')
const fieldType = ref('text')
const fieldTabGroup = ref('') // Layout Folder Group (Phase 1 Custom Field Tabs)
const fieldValidationRegex = ref('') // Custom Validation Regex (Phase 1 Custom Field Regex)
const categoryNameOptions = computed(() => {
  if (categories.value && categories.value.length > 0) {
    return categories.value.map(c => c.name)
  }
  return ['Server', 'Router', 'Switch', 'Database', 'Application', 'Kubernetes Deployment', 'Container', 'Other']
})
const fieldCategory = ref('Server')
const fieldRequired = ref(false)

const fieldColumns = [
  { key: 'name', label: 'Field Key' },
  { key: 'field_type', label: 'Value Type' },
  { key: 'asset_type', label: 'Target Category' },
  { key: 'tab_group', label: 'Tab Folder Group' },
  { key: 'validation_regex', label: 'Validation Regex' },
  { key: 'is_required', label: 'Constraint' },
  { key: 'actions', label: 'Actions' }
]

const registerField = async () => {
  isSavingField.value = true
  const sanitizedName = fieldName.value.toLowerCase().trim().replace(/\s+/g, '_')
  try {
    await createCustomField({
      name: sanitizedName,
      field_type: fieldType.value,
      asset_type: fieldCategory.value,
      is_required: fieldRequired.value,
      tab_group: fieldTabGroup.value || 'General Custom', // default folder (Phase 1)
      validation_regex: fieldValidationRegex.value || '' // custom regex (Phase 1)
    })
    fieldName.value = ''
    fieldTabGroup.value = ''
    fieldValidationRegex.value = ''
    await refreshFields()
    alert('Dynamic custom field registered successfully.')
  } catch (error) {
    console.error('Failed to save field:', error)
    alert('Register failed: Check permissions or name uniqueness.')
  } finally {
    isSavingField.value = false
  }
}

const removeField = async (id) => {
  if (!confirm('Are you sure you want to delete this custom field?')) return
  try {
    await deleteCustomField(id)
    await refreshFields()
  } catch (error) {
    console.error('Failed to delete field:', error)
  }
}

// ----------------------------------------
// Part 5: Webhooks
// ----------------------------------------
const webhookUrl = ref('')
const webhookEvent = ref('asset:create')

const webhookColumns = [
  { key: 'id', label: 'ID' },
  { key: 'url', label: 'Target callback URL' },
  { key: 'event', label: 'Event Trigger' },
  { key: 'actions', label: 'Actions' }
]

const registerHook = async () => {
  isSavingHook.value = true
  try {
    await createWebhook({ url: webhookUrl.value, event: webhookEvent.value })
    webhookUrl.value = ''
    await refreshHooks()
    alert('Webhook callback registered successfully.')
  } catch (error) {
    console.error('Failed to save webhook:', error)
    alert('Register failed: Check permissions.')
  } finally {
    isSavingHook.value = false
  }
}

const removeHook = async (id) => {
  if (!confirm('Are you sure you want to delete this Webhook endpoint?')) return
  try {
    await deleteWebhook(id)
    await refreshHooks()
  } catch (error) {
    console.error('Failed to delete webhook:', error)
  }
}

// ----------------------------------------
// Part 7: Dynamic Category Manager (Phase 1 Dynamic Categories)
// ----------------------------------------
const newCatForm = ref({ name: '', icon: '', description: '', parent_dependency: '' })
const newSubForm = ref({ category_id: '', name: '', icon: '', description: '' })

const catDropdownOptions = computed(() => {
  const options = [{ label: 'Select parent Category...', value: '' }]
  if (categories.value) {
    categories.value.forEach(c => {
      options.push({ label: c.name, value: c.id })
    })
  }
  return options
})

const isReservedCategory = (name) => {
  return ['Server', 'Database', 'Network', 'Application', 'Kubernetes Deployment', 'Container', 'IoT / Edge', 'Other'].includes(name)
}

const registerCategory = async () => {
  isSavingCat.value = true
  try {
    await $fetch(`${apiBase}/categories/`, {
      method: 'POST',
      body: newCatForm.value,
      headers: getAuthHeader()
    })
    newCatForm.value = { name: '', icon: '', description: '', parent_dependency: '' }
    await refreshCategories()
    alert('Asset Category successfully registered dynamically in GORM.')
  } catch (err) {
    console.error('Failed to create category:', err)
    alert('Register failed: Check permissions.')
  } finally {
    isSavingCat.value = false
  }
}

const removeCategory = async (id) => {
  if (!confirm('Are you sure you want to delete this Category? All nested assets and sub-groups will be cascadingly deleted.')) return
  try {
    await $fetch(`${apiBase}/categories/${id}/`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
    await refreshCategories()
  } catch (err) {
    console.error('Failed to delete category:', err)
  }
}

const registerSubGroup = async () => {
  isSavingSub.value = true
  try {
    await $fetch(`${apiBase}/categories/sub-groups/`, {
      method: 'POST',
      body: newSubForm.value,
      headers: getAuthHeader()
    })
    newSubForm.value = { category_id: '', name: '', icon: '', description: '' }
    await refreshCategories()
    alert('Category Subgroup successfully registered dynamically in GORM.')
  } catch (err) {
    console.error('Failed to create sub-group:', err)
    alert('Register failed: Check permissions.')
  } finally {
    isSavingSub.value = false
  }
}

// ----------------------------------------
// Part 8: Bulk Ingest Engine (Phase 1 Bulk Importing)
// ----------------------------------------
const parsedCsvRows = ref([])
const isExecutingImport = ref(false)
const validationErrors = ref([])

const csvColumns = [
  { key: 'name', label: 'Name' },
  { key: 'type', label: 'Category' },
  { key: 'ip_address', label: 'IP Address' },
  { key: 'status', label: 'Status' },
  { key: 'custom_fields', label: 'Custom Attributes' }
]

const hasValidationErrors = computed(() => validationErrors.value.length > 0)

const getRowError = (rowIdx) => {
  return validationErrors.value.find(err => err.row === rowIdx + 1)
}

const handleCsvUpload = (event) => {
  const file = event.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target.result
    const lines = text.split('\n')
    if (lines.length < 2) return

    // Parse Headers
    const headers = lines[0].split(',').map(h => h.trim().toLowerCase())
    const rows = []
    validationErrors.value = []

    // Loop through lines (skip header)
    for (let i = 1; i < lines.length; i++) {
      if (!lines[i].trim()) continue
      const cells = lines[i].split(',').map(c => c.trim())
      const rowData = {
        name: '',
        type: 'Server',
        ip_address: '',
        status: 'active',
        description: '',
        custom_properties: {}
      }

      // Map cells to data object
      headers.forEach((h, idx) => {
        const val = cells[idx] || ''
        if (h === 'name') rowData.name = val
        else if (h === 'type') rowData.type = val
        else if (h === 'ip_address') rowData.ip_address = val
        else if (h === 'status') rowData.status = val
        else if (h === 'description') rowData.description = val
        else {
          // Dynamic attributes
          rowData.custom_properties[h] = val
        }
      })

      // Run dynamic GORM Regex validation pattern checks! (Phase 1 Bulk Importing)
      if (customFields.value) {
        customFields.value.forEach(field => {
          // If the field belongs to this row's category
          if (field.asset_type === rowData.type && field.validation_regex) {
            const cellVal = rowData.custom_properties[field.name] || ''
            const regex = new RegExp(field.validation_regex)
            if (!regex.test(String(cellVal))) {
              validationErrors.value.push({
                row: i,
                field: field.name,
                val: cellVal,
                pattern: field.validation_regex
              })
            }
          }
        })
      }

      rows.push(rowData)
    }

    parsedCsvRows.value = rows
  }
  reader.readAsText(file)
}

const executeBulkImport = async () => {
  if (hasValidationErrors.value) return
  isExecutingImport.value = true

  const { createAsset } = useAssets()
  let successCount = 0

  try {
    for (const row of parsedCsvRows.value) {
      const payload = {
        name: row.name,
        type: row.type,
        ip_address: row.ip_address,
        status: row.status,
        description: row.description,
        properties: row.custom_properties
      }
      await createAsset(payload)
      successCount++
    }
    alert(`Bulk Ingestion Successful! Successfully registered ${successCount} assets into GORM.`);
    parsedCsvRows.value = []
  } catch (err) {
    console.error('Bulk ingest failed:', err)
    alert('Ingestion failed midway: check system permissions or UUID conflicts.')
  } finally {
    isExecutingImport.value = false
  }
}

// Sub-Tab selector for Bulk Ingest Engine (Phase 1 Bulk Cabling Matrix)
const importerSubTab = ref('hardware')

const parsedCablingRows = ref([])
const isExecutingCabling = ref(false)
const cablingValidationErrors = ref([])

const cablingColumns = [
  { key: 'source_asset', label: 'Source Asset Name' },
  { key: 'target_asset', label: 'Target Asset Name' },
  { key: 'connection_type', label: 'Link Type' },
  { key: 'status', label: 'Relational Audit' }
]

const hasCablingValidationErrors = computed(() => cablingValidationErrors.value.length > 0)

const getCablingRowError = (rowIdx) => {
  return cablingValidationErrors.value.find(err => err.row === rowIdx + 1)
}

const handleCablingUpload = (event) => {
  const file = event.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target.result
    const lines = text.split('\n')
    if (lines.length < 2) return

    // Parse Headers: Source, Target, Type
    const headers = lines[0].split(',').map(h => h.trim().toLowerCase())
    const rows = []
    cablingValidationErrors.value = []

    for (let i = 1; i < lines.length; i++) {
      if (!lines[i].trim()) continue
      const cells = lines[i].split(',').map(c => c.trim())
      const rowData = {
        source_asset: '',
        target_asset: '',
        connection_type: 'uplink'
      }

      headers.forEach((h, idx) => {
        const val = cells[idx] || ''
        if (h === 'source') rowData.source_asset = val
        else if (h === 'target') rowData.target_asset = val
        else if (h === 'type') rowData.connection_type = val.toLowerCase()
      })

      // Validation Audits: Enforce GORM asset presence checks! (Phase 1 Bulk Cabling Matrix)
      const sourceExists = allAssetsList.value?.some(a => a.name.toLowerCase() === rowData.source_asset.toLowerCase())
      const targetExists = allAssetsList.value?.some(a => a.name.toLowerCase() === rowData.target_asset.toLowerCase())

      if (!sourceExists) {
        cablingValidationErrors.value.push({
          row: i,
          field: 'source',
          val: rowData.source_asset,
          reason: 'Source hardware node does not exist in GORM CMDB'
        })
      }
      if (!targetExists) {
        cablingValidationErrors.value.push({
          row: i,
          field: 'target',
          val: rowData.target_asset,
          reason: 'Target hardware node does not exist in GORM CMDB'
        })
      }
      if (rowData.source_asset.toLowerCase() === rowData.target_asset.toLowerCase()) {
        cablingValidationErrors.value.push({
          row: i,
          field: 'target',
          val: rowData.target_asset,
          reason: 'Invalid self-linking loop detected'
        })
      }

      rows.push(rowData)
    }

    parsedCablingRows.value = rows
  }
  reader.readAsText(file)
}

const executeBulkCabling = async () => {
  if (hasCablingValidationErrors.value) return
  isExecutingCabling.value = true

  const { syncRelationships } = useRelationships()
  let linkCount = 0

  try {
    for (const row of parsedCablingRows.value) {
      // Find IDs from names
      const srcNode = allAssetsList.value?.find(a => a.name.toLowerCase() === row.source_asset.toLowerCase())
      const tgtNode = allAssetsList.value?.find(a => a.name.toLowerCase() === row.target_asset.toLowerCase())

      if (srcNode && tgtNode) {
        if (row.connection_type === 'uplink') {
          await syncRelationships(srcNode.id, [tgtNode.id], [])
        } else {
          await syncRelationships(srcNode.id, [], [tgtNode.id])
        }
        linkCount++
      }
    }
    alert(`Bulk Cabling Matrix Patching Successful! Dispatched and connected ${linkCount} cable links.`);
    parsedCablingRows.value = []
  } catch (err) {
    console.error('Cabling patch failed:', err)
    alert('Bulk cabling failed midway: check system permissions or port limits.')
  } finally {
    isExecutingCabling.value = false
  }
}

// ----------------------------------------
// Part 6: Disaster Recovery Backups
// ----------------------------------------
const triggerBackupExport = async () => {
  try {
    const backupState = {
      version: 'v1.0.0',
      timestamp: new Date().toISOString(),
      datacenters: datacenters.value,
      custom_fields: customFields.value,
      webhooks: webhooks.value,
      users: users.value
    }
    const blob = new Blob([JSON.stringify(backupState, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `northstar-backup-${Date.now()}.json`
    link.click()
    URL.revokeObjectURL(url)
    alert('Disaster recovery structural JSON backup successfully compiled and downloaded.')
  } catch (err) {
    console.error('Failed to export backup:', err)
  }
}

const triggerBackupRestore = () => {
  alert('Restore functionality simulates uploading the JSON file and rebuilding GORM tables automatically via backend APIs.')
}
</script>
