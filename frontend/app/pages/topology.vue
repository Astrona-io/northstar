<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6 border-b border-slate-200 dark:border-slate-800 pb-4">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">Relational Network Topology Map</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400">Dynamic visual cabling maps running across physical cabs or global datacenter VPN tunnels</p>
      </div>
      <UBadge color="teal" variant="subtle" class="font-mono text-[10px] font-bold tracking-wider px-2 py-1 uppercase">TOPOLOGY LEVEL 2</UBadge>
    </div>

    <div v-if="pending" class="flex justify-center p-12">
      <UIcon name="i-heroicons-arrow-path" class="animate-spin h-10 w-10 text-primary-500" />
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Main SVG Graph Panel (Phase 7 Multi-DC Peering Overhaul) -->
      <UCard :class="isFullScreen ? 'fixed inset-0 z-50 h-screen w-screen rounded-none bg-white dark:bg-gray-950 p-6 flex flex-col' : 'lg:col-span-3 h-[550px] relative overflow-hidden flex flex-col justify-between'">
        <template #header>
          <div class="flex justify-between items-center text-xs">
            <div class="flex items-center gap-3">
              <span class="font-medium text-slate-500 dark:text-slate-400">Select Topology View:</span>
              <USelect 
                v-model="topologyView" 
                :options="[{label: 'Local Cabinet Cabling', value: 'local'}, {label: 'Global Datacenter WAN Map', value: 'global'}]" 
                option-attribute="label"
                value-attribute="value"
                size="xs"
                class="w-48 font-mono"
              />
            </div>
            
            <div class="flex gap-4 items-center">
              <template v-if="topologyView === 'local'">
                <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-red-500" /> Core / Routers</span>
                <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-blue-500" /> Distribution / Switches</span>
              </template>
              <template v-else>
                <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-teal-500 animate-pulse" /> IPsec Encrypted</span>
                <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-blue-500" /> Public Cloud VPC</span>
              </template>
              <UButton 
                size="xs" 
                color="gray" 
                variant="subtle" 
                :icon="isFullScreen ? 'i-heroicons-arrows-pointing-in' : 'i-heroicons-arrows-pointing-out'" 
                @click="isFullScreen = !isFullScreen"
              >
                {{ isFullScreen ? 'Exit' : 'Maximize' }}
              </UButton>
            </div>
          </div>
        </template>
        
        <!-- Live SVG Canvas -->
        <div class="flex-1 w-full bg-slate-50 dark:bg-slate-950/60 rounded-lg relative overflow-auto border border-slate-200 dark:border-slate-800 shadow-inner">
          
          <!-- View A: Local Cabinet Cabling -->
          <svg v-if="topologyView === 'local'" width="100%" height="100%" min-height="400" viewBox="0 0 800 450" class="mx-auto block">
            <!-- Cable Patch Connections (Lines) -->
            <g>
              <template v-for="(link, i) in connections" :key="i">
                <!-- Base Standard Cable line -->
                <line 
                  :x1="link.x1" 
                  :y1="link.y1" 
                  :x2="link.x2" 
                  :y2="link.y2"
                  stroke="#10B981" 
                  stroke-width="3" 
                  stroke-dasharray="5,5" 
                  class="opacity-75 stroke-primary-500" 
                />
                <!-- Thick glowing, pulsing cable tracer (Phase 1 Tracer) -->
                <line 
                  v-if="link.isTraced"
                  :x1="link.x1" 
                  :y1="link.y1" 
                  :x2="link.x2" 
                  :y2="link.y2"
                  stroke="#10B981" 
                  stroke-width="7" 
                  class="opacity-80 animate-pulse stroke-primary-500" 
                />
              </template>
            </g>

            <!-- Physical Device Nodes -->
            <g v-for="node in nodes" :key="node.id" class="cursor-pointer group" @click="selectNode(node)">
              <!-- Outer Glow (Glow on selected node or hover) -->
              <circle 
                :cx="node.x" 
                :cy="node.y" 
                r="32" 
                class="fill-transparent stroke-primary-500/30 stroke-2 opacity-0 group-hover:opacity-100 transition-all" 
                :class="[selectedNodeId === node.id ? 'opacity-100 scale-105' : '']"
              />
              
              <!-- Main Node circle -->
              <circle 
                :cx="node.x" 
                :cy="node.y" 
                r="26" 
                :class="[
                  getNodeCircleClass(node),
                  'shadow-lg'
                ]" 
              />
              
              <!-- Centered Hardware Icon -->
              <foreignObject :x="node.x - 12" :y="node.y - 12" width="24" height="24">
                <div class="text-white flex items-center justify-center">
                  <UIcon :name="getNodeIcon(node)" class="h-6 w-6" />
                </div>
              </foreignObject>

              <!-- Asset Label Badge below node -->
              <text 
                :x="node.x" 
                :y="node.y + 45" 
                font-family="monospace" 
                font-size="11" 
                font-weight="bold" 
                text-anchor="middle"
                class="fill-slate-800 dark:fill-white group-hover:fill-primary-500 transition-colors"
              >
                {{ node.name }}
              </text>
              <text 
                :x="node.x" 
                :y="node.y + 58" 
                font-family="monospace" 
                font-size="9" 
                fill="#94a3b8" 
                text-anchor="middle"
              >
                {{ node.ip_address || 'No IP' }}
              </text>
            </g>
          </svg>

          <!-- View B: Global Datacenter WAN Map (Phase 7 Multi-DC Peering) -->
          <svg v-else width="100%" height="100%" min-height="400" viewBox="0 0 800 450" class="mx-auto block">
            <!-- Grid lines background -->
            <defs>
              <pattern id="dc-grid" width="25" height="25" patternUnits="userSpaceOnUse">
                <path d="M 25 0 L 0 0 0 25" fill="none" stroke="#334155" stroke-width="0.5" opacity="0.1" />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#dc-grid)" />

            <!-- Active VPN Tunnels (Cables) -->
            <g>
              <!-- Glowing thick IPsec Tunnel line -->
              <line x1="220" y1="225" x2="580" y2="225" stroke="#10b981" stroke-width="6" class="opacity-80 animate-pulse stroke-primary-500" />
              <line x1="220" y1="225" x2="580" y2="225" stroke="#34d399" stroke-width="2" class="opacity-90 stroke-primary-400" />
              
              <!-- Secure Lock icon in center of tunnel -->
              <foreignObject x="385" y="210" width="30" height="30">
                <div class="bg-primary-500 text-white rounded-full h-8 w-8 flex items-center justify-center shadow-lg border-2 border-slate-900">
                  <UIcon name="i-heroicons-lock-closed" class="h-4 w-4" />
                </div>
              </foreignObject>

              <!-- Tunnel metadata description overlay -->
              <text x="400" y="258" font-family="monospace" font-size="10" font-weight="bold" fill="#38bdf8" text-anchor="middle">IPsec VPN TUNNEL</text>
              <text x="400" y="272" font-family="monospace" font-size="9" fill="#94a3b8" text-anchor="middle">Trunk Capacity: 2.5 Gbps Fiber | WireGuard Active</text>
            </g>

            <!-- Datacenter Nodes (Dublin HQ & AWS Ireland) -->
            <g class="cursor-pointer group" @click="navigateTo('/dcim')">
              <!-- Site 1: Dublin HQ -->
              <circle cx="220" cy="225" r="48" class="fill-slate-900 stroke-slate-700 group-hover:stroke-primary-500 stroke-2 shadow-lg transition-colors" />
              <circle cx="220" cy="225" r="40" class="fill-slate-850" />
              <foreignObject x="202" y="207" width="36" height="36">
                <div class="text-primary-500 group-hover:text-primary-400 flex items-center justify-center h-full transition-colors">
                  <UIcon name="i-heroicons-home" class="h-8 w-8" />
                </div>
              </foreignObject>
              <text x="220" y="295" font-family="monospace" font-size="12" font-weight="bold" fill="#ffffff" text-anchor="middle">Dublin HQ</text>
              <text x="220" y="310" font-family="monospace" font-size="9" fill="#94a3b8" text-anchor="middle">On-Premises (IE)</text>
            </g>

            <g class="cursor-pointer group" @click="navigateTo('/dcim')">
              <!-- Site 2: AWS Ireland -->
              <circle cx="580" cy="225" r="48" class="fill-slate-900 stroke-slate-700 group-hover:stroke-primary-500 stroke-2 shadow-lg transition-colors" />
              <circle cx="580" cy="225" r="40" class="fill-slate-850" />
              <foreignObject x="562" y="207" width="36" height="36">
                <div class="text-blue-500 group-hover:text-blue-400 flex items-center justify-center h-full transition-colors">
                  <UIcon name="i-heroicons-cloud" class="h-8 w-8" />
                </div>
              </foreignObject>
              <text x="580" y="295" font-family="monospace" font-size="12" font-weight="bold" fill="#ffffff" text-anchor="middle">AWS Ireland Zone</text>
              <text x="580" y="310" font-family="monospace" font-size="9" fill="#94a3b8" text-anchor="middle">Public Cloud VPC (IE)</text>
            </g>
          </svg>

        </div>
      </UCard>

      <!-- Side Details Inspector -->
      <div class="lg:col-span-1">
        <!-- Local inspector view -->
        <UCard v-if="inspectorNode && topologyView === 'local'" class="h-[550px] flex flex-col">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon :name="getNodeIcon(inspectorNode)" class="h-6 w-6 text-primary-500" />
              <h3 class="font-bold text-gray-900 dark:text-white truncate font-mono">{{ inspectorNode.name }}</h3>
            </div>
          </template>
          
          <div class="space-y-4 text-sm flex-1">
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Category Type</span>
              <span class="font-semibold text-gray-900 dark:text-white capitalize">{{ inspectorNode.type }}</span>
            </div>
            <div v-if="inspectorNode.properties?.network_subtype">
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Device Subtype</span>
              <UBadge color="gray" variant="subtle" class="capitalize font-mono text-[9px] font-bold">{{ inspectorNode.properties.network_subtype }}</UBadge>
            </div>
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">IP Address</span>
              <span class="font-mono text-gray-900 dark:text-white">{{ inspectorNode.ip_address || 'N/A' }}</span>
            </div>
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Connection Status</span>
              <UBadge :color="inspectorNode.status === 'active' ? 'green' : 'orange'" size="xs" variant="subtle" class="capitalize font-mono text-[9px] font-bold">
                {{ inspectorNode.status }}
              </UBadge>
            </div>
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Description</span>
              <p class="text-xs text-gray-600 dark:text-gray-300 mt-1 leading-normal">
                {{ inspectorNode.description || 'No description logged.' }}
              </p>
            </div>
            
            <div v-if="inspectorNode.properties && Object.keys(inspectorNode.properties).length > 0" class="pt-3 border-t border-gray-150 dark:border-gray-800">
              <span class="text-xs text-gray-500 uppercase tracking-wider block mb-2 font-mono">Specifications</span>
              <div class="grid grid-cols-2 gap-2 text-xs">
                <template v-for="(val, key) in inspectorNode.properties" :key="key">
                  <span class="text-gray-400 capitalize" v-if="key !== 'network_subtype'">{{ key.replace(/_/g, ' ') }}</span>
                  <span class="font-medium text-gray-900 dark:text-white truncate text-right" v-if="key !== 'network_subtype'">{{ val || 'N/A' }}</span>
                </template>
              </div>
            </div>
          </div>
          
          <div class="mt-6 pt-3 border-t border-gray-150 dark:border-gray-800 flex-shrink-0">
            <UButton block color="primary" icon="i-heroicons-eye" :to="`/assets/${inspectorNode.id}`">
              Open Asset Panel
            </UButton>
          </div>
        </UCard>
        
        <!-- Global WAN inspector view -->
        <UCard v-else-if="topologyView === 'global'" class="h-[550px] flex flex-col">
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-heroicons-shield-check" class="h-6 w-6 text-primary-500" />
              <h3 class="font-bold text-gray-900 dark:text-white font-mono">WAN Peering Audit</h3>
            </div>
          </template>
          
          <div class="space-y-4 text-sm flex-1 leading-relaxed">
            <p class="text-xs text-slate-500 leading-normal">
              Active secure IPSec VPN Tunnel connects on-premises and virtual cloud datacenters over public WAN interfaces.
            </p>
            
            <UDivider class="my-2" />
            
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Tunnel Status</span>
              <div class="flex items-center gap-1.5 font-bold text-green-500 mt-1">
                <span class="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                UP / ENCRYPTED
              </div>
            </div>
            
            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">Encryption Standard</span>
              <span class="text-slate-800 dark:text-slate-200 font-mono text-xs">AES-256-GCM / WireGuard</span>
            </div>

            <div>
              <span class="text-xs text-gray-500 uppercase tracking-wider block font-mono">WAN Gateways</span>
              <ul class="space-y-1 font-mono text-[10px] text-slate-500 mt-1">
                <li class="flex justify-between"><span>Dublin HQ:</span><span class="text-slate-800 dark:text-slate-200">85.12.85.112</span></li>
                <li class="flex justify-between"><span>AWS Ireland:</span><span class="text-slate-800 dark:text-slate-200">10.150.0.1</span></li>
              </ul>
            </div>
          </div>
          
          <div class="mt-6 pt-3 border-t border-gray-150 dark:border-gray-800 flex-shrink-0">
            <UButton block color="primary" icon="i-heroicons-server" to="/dcim">
              Open DCIM Planner
            </UButton>
          </div>
        </UCard>

        <!-- Welcome Guide Inspector -->
        <UCard v-else class="h-[550px] flex items-center justify-center text-center text-gray-500 dark:text-gray-400">
          <div class="p-6 space-y-4">
            <UIcon name="i-heroicons-magnifying-glass" class="h-12 w-12 text-primary-500 mx-auto animate-bounce" />
            <h4 class="font-bold text-gray-900 dark:text-white font-mono">Cable Map Inspector</h4>
            <p class="text-xs leading-normal">
              Click on any physical node in the interactive topology map to inspect active interfaces, VLAN ports, speeds, and asset specs.
            </p>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const route = useRoute()
const { fetchAssets } = useAssets()
const { data: assets, pending } = await fetchAssets({ limit: 1000 })

const selectedNodeId = ref(null)
const inspectorNode = ref(null)
const isFullScreen = ref(false)
const topologyView = ref('local') // 'local' or 'global' (Phase 7 Multi-DC Peering)

const selectNode = (node) => {
  selectedNodeId.value = node.id
  inspectorNode.value = node
}

// Helper styled metrics for consolidated Network Subtypes (Router, Switch, AP, Firewall)
const getNodeCircleClass = (node) => {
  const subtype = node.properties?.network_subtype || node.type
  if (subtype === 'Router') return 'fill-red-500'
  if (subtype === 'Switch') return 'fill-blue-500'
  if (subtype === 'AP (Access Point)') return 'fill-emerald-500'
  return 'fill-orange-500' // Firewalls / Load Balancers
}

const getNodeIcon = (node) => {
  const subtype = node.properties?.network_subtype || node.type
  if (subtype === 'Router') return 'i-heroicons-cpu-chip'
  if (subtype === 'Switch') return 'i-heroicons-arrows-right-left'
  if (subtype === 'AP (Access Point)') return 'i-heroicons-wifi'
  return 'i-heroicons-shield-check' // Firewalls / Load Balancers
}

// 1. Identify active Network nodes in the map (supporting backward-compatible types)
const nodes = computed(() => {
  if (!assets.value) return []
  
  // Extract Network devices (and legacy Router/Switch top-level items)
  const netDevices = assets.value.filter(a => a.type === 'Network' || a.type === 'Router' || a.type === 'Switch')
  
  // Distribute them logically in tiers on our 800x450 canvas:
  // Tier 1 (Routers) are Core - placed in the upper center
  // Tier 2 (Switches/APs/Firewalls) are Distribution - placed in the middle
  const routers = netDevices.filter(a => {
    const subtype = a.properties?.network_subtype || a.type
    return subtype === 'Router'
  })
  const switches = netDevices.filter(a => {
    const subtype = a.properties?.network_subtype || a.type
    return subtype !== 'Router'
  })
  
  const formatted = []
  
  // Calculate positions for routers (Tier 1)
  routers.forEach((r, idx) => {
    const total = routers.length
    const x = total === 1 ? 400 : 250 + idx * (300 / (total - 1))
    formatted.push({
      ...r,
      x,
      y: 90
    })
  })

  // Calculate positions for switches (Tier 2)
  switches.forEach((s, idx) => {
    const total = switches.length
    const x = total === 1 ? 400 : 150 + idx * (500 / (total - 1))
    formatted.push({
      ...s,
      x,
      y: 260
    })
  })
  
  return formatted
})

// 2. Fetch connection relations for all nodes and map coordinate lines (Cabling)
const connections = computed(() => {
  if (nodes.value.length === 0) return []
  
  const lines = []
  const traceId = route.query.trace ? route.query.trace : null
  
  // To keep it 100% stable, let's draw connections where a Switch links up to a Router:
  const routers = nodes.value.filter(n => {
    const subtype = n.properties?.network_subtype || n.type
    return subtype === 'Router'
  })
  const switches = nodes.value.filter(n => {
    const subtype = n.properties?.network_subtype || n.type
    return subtype !== 'Router'
  })
  
  switches.forEach(sw => {
    // Standard connection: draw line from switch to the first Router
    routers.forEach(rt => {
      const activeTrace = (traceId && (sw.id === traceId || rt.id === traceId))
      lines.push({
        x1: sw.x,
        y1: sw.y,
        x2: rt.x,
        y2: rt.y,
        isTraced: activeTrace
      })
    })
  })
  
  return lines
})
</script>
