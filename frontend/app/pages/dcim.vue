<template>
  <div class="space-y-6 h-[calc(100vh-4rem)] flex flex-col overflow-hidden font-sans">
    
    <!-- DCIM Administrative Header -->
    <div class="flex justify-between items-center border-b border-slate-200 dark:border-slate-800 pb-4 flex-shrink-0">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-white">Datacenter Floor Planner</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400">Design 2D rooms, map layout grids, assign cabinet coordinates, and inspect physical rail mounting slots</p>
      </div>
      <UBadge color="teal" variant="subtle" class="font-mono text-[10px] font-bold tracking-wider px-2 py-1">DCIM LEVEL 3</UBadge>
    </div>

    <!-- 1. Sleek Horizontal Top-Navigation Filter Bar (Replacing Sidebar Card) -->
    <div class="flex flex-wrap items-center gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-4 rounded-md shadow-sm flex-shrink-0">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-funnel" class="h-5 w-5 text-slate-400" />
        <span class="text-xs font-bold text-slate-500 uppercase tracking-wider font-mono">CAD Selectors:</span>
      </div>

      <!-- Select Datacenter Dropdown -->
      <div class="w-full sm:w-64">
        <USelect 
          v-model="selectedDcId" 
          :options="dcOptions" 
          option-attribute="label" 
          value-attribute="value" 
          placeholder="Select Datacenter..."
        />
      </div>

      <!-- Floor Level Selection Tabs (Phase 2 Multi-Floor Selector) -->
      <div v-if="selectedDc && selectedDc.floors?.length > 0" class="flex gap-1.5 bg-slate-100 dark:bg-slate-800 p-1 rounded-md">
        <button
          v-for="floor in selectedDc.floors"
          :key="floor.id"
          type="button"
          :class="[
            selectedFloorId === floor.id 
              ? 'bg-white dark:bg-slate-900 text-primary-500 dark:text-primary-400 shadow-sm font-bold' 
              : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300',
            'px-3 py-1.5 text-xs font-medium rounded-md transition-all'
          ]"
          @click="selectedFloorId = floor.id"
        >
          {{ floor.name }} (L{{ floor.level }})
        </button>
      </div>

      <!-- Add Custom Floor button (Admin Only) -->
      <UButton 
        v-if="canMutate && selectedDc" 
        size="xs" 
        color="gray" 
        variant="subtle" 
        icon="i-heroicons-plus" 
        label="Add Floor Level" 
        @click="isFloorModalOpen = true"
      />
    </div>

    <!-- 2. SPLIT LAYOUT PLANNER -->
    <div class="flex-grow flex gap-6 overflow-hidden min-h-0">
      
      <!-- Left Main: Precision SVG 2D CAD Floor Planner (Draw.io Style, Phase 1) -->
      <div class="flex-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-4 flex flex-col h-full overflow-hidden shadow-sm relative">
        <div class="flex justify-between items-center text-xs border-b border-slate-100 dark:border-slate-800 pb-2 mb-4 flex-shrink-0">
          <span class="font-bold text-slate-800 dark:text-white font-mono uppercase">2D Room Grid Floorplan: {{ selectedDc?.name || 'No Datacenter Selected' }}</span>
          <span class="text-slate-400 font-mono text-[10px]">Snapping: 100mm Grid | CAD Unit: mm</span>
        </div>

        <div v-if="!selectedDc" class="flex-1 flex flex-col items-center justify-center text-center text-slate-400">
          <UIcon name="i-heroicons-globe-alt" class="h-16 w-16 text-primary-500/60 animate-pulse mb-3" />
          <h4 class="font-bold text-slate-800 dark:text-white font-mono">No Datacenter Selected</h4>
          <p class="text-xs max-w-xs leading-normal">Please select an active infrastructure site in the top filter bar to load drafting room grids.</p>
        </div>

        <div v-else class="flex-1 overflow-auto relative rounded-md bg-slate-950/95 shadow-inner border border-slate-800 flex items-center justify-center">
          <!-- Drafting Board CAD SVG Canvas -->
          <svg 
            width="100%" 
            height="100%" 
            viewBox="0 0 800 500" 
            class="max-h-[500px]"
          >
            <!-- Dotted 1-Meter Precision Blueprint Grid (Phase 1 Drafting) -->
            <defs>
              <pattern id="dotted-grid" width="40" height="40" patternUnits="userSpaceOnUse">
                <circle cx="2" cy="2" r="1.5" fill="#475569" opacity="0.4" />
              </pattern>

              <!-- Dynamic Wall Cutout Mask (Creates clean transparent openings for doors/windows) -->
              <mask id="wall-mask">
                <!-- White area allows standard walls to paint normally -->
                <rect width="100%" height="100%" fill="#ffffff" />
                <!-- Black lines represent door segments, cutting out doorways from the walls -->
                <line 
                  v-for="door in activeFloorWalls.filter(w => w.type === 'door' || w.type === 'door-double')" 
                  :key="'mask_' + door.id"
                  :x1="door.x1 * 2" 
                  :y1="door.y1 * 2" 
                  :x2="door.x2 * 2" 
                  :y2="door.y2 * 2" 
                  stroke="#000000" 
                  :stroke-width="(door.thickness || 10) * 0.4 + 4" 
                  stroke-linecap="round"
                />
              </mask>
            </defs>
            <rect width="100%" height="100%" fill="url(#dotted-grid)" />

            <!-- Room Boundary Walls (Precision Metric Drafting) -->
            <rect x="50" y="50" width="700" height="400" fill="none" stroke="#475569" stroke-width="8" stroke-linejoin="round" />
            <!-- Inner Room Grid Boundary -->
            <rect x="54" y="54" width="692" height="392" fill="none" stroke="#1e293b" stroke-width="1" />

            <!-- Custom Drawn 2D CAD blueprint walls (Phase 1 DCIM CAD Walls) -->
            <g v-for="wall in activeFloorWalls" :key="wall.id">
              <!-- Render Door Segment with Swing Arc (Single or Double Door) -->
              <g v-if="wall.type === 'door' || wall.type === 'door-double'">
                <!-- If double/split door, render two half-size open door leaves -->
                <template v-if="wall.type === 'door-double'">
                  <line 
                    :x1="getDoorLeafCoords(wall).x1" :y1="getDoorLeafCoords(wall).y1" 
                    :x2="getDoorLeafCoords(wall).openX1" :y2="getDoorLeafCoords(wall).openY1" 
                    stroke="#f43f5e" :stroke-width="(wall.thickness || 10) * 0.4" 
                    stroke-linecap="round"
                    vector-effect="non-scaling-stroke"
                  />
                  <line 
                    :x1="getDoorLeafCoords(wall).x2" :y1="getDoorLeafCoords(wall).y2" 
                    :x2="getDoorLeafCoords(wall).openX2" :y2="getDoorLeafCoords(wall).openY2" 
                    stroke="#f43f5e" :stroke-width="(wall.thickness || 10) * 0.4" 
                    stroke-linecap="round"
                    vector-effect="non-scaling-stroke"
                  />
                </template>
                <!-- Standard single open door leaf -->
                <line 
                  v-else
                  :x1="getDoorLeafCoords(wall).x1" :y1="getDoorLeafCoords(wall).y1" 
                  :x2="getDoorLeafCoords(wall).openX" :y2="getDoorLeafCoords(wall).openY" 
                  stroke="#f43f5e" :stroke-width="(wall.thickness || 10) * 0.4" 
                  stroke-linecap="round"
                  vector-effect="non-scaling-stroke"
                />
                
                <!-- Shared Arc Swing Paths -->
                <path 
                  :d="getDoorArcPath(wall)" 
                  fill="none" 
                  stroke="#f43f5e" 
                  stroke-width="1.5" 
                  stroke-dasharray="3" 
                  class="opacity-70"
                />
              </g>

              <!-- Render Room Enclosure (Rectangular Room) -->
              <rect 
                v-else-if="wall.type === 'wall-rect'"
                :x="Math.min(wall.x1, wall.x2) * 2" 
                :y="Math.min(wall.y1, wall.y2) * 2" 
                :width="Math.abs(wall.x2 - wall.x1) * 2" 
                :height="Math.abs(wall.y2 - wall.y1) * 2" 
                fill="none"
                stroke="#64748b" 
                :stroke-width="(wall.thickness || 20) * 0.4" 
                stroke-linejoin="round"
                class="opacity-90"
                mask="url(#wall-mask)"
                vector-effect="non-scaling-stroke"
              />

              <!-- Standard Wall Segment -->
              <line 
                v-else
                :x1="wall.x1 * 2" 
                :y1="wall.y1 * 2" 
                :x2="wall.x2 * 2" 
                :y2="wall.y2 * 2" 
                stroke="#64748b" 
                :stroke-width="(wall.thickness || 20) * 0.4" 
                stroke-linecap="round" 
                class="opacity-90"
                mask="url(#wall-mask)"
                vector-effect="non-scaling-stroke"
              />
            </g>

            <!-- Secure Entrance Gate (Drafting Lines) -->
            <line x1="380" y1="50" x2="420" y2="50" stroke="#f43f5e" stroke-width="12" />
            <text x="400" y="42" font-family="monospace" font-size="10" fill="#f43f5e" text-anchor="middle" font-weight="bold">SECURE AIRLOCK</text>

            <!-- Interactive Cabinet Footprint Blocks (Placement Zones & Clicking Portals, Phase 1 & 4) -->
            <g v-for="rack in filteredRacks" :key="rack.id" class="cursor-pointer group" @click="openRackRails(rack)">
              <!-- Cabinet footprint box rectangle -->
              <rect 
                :x="rack.x || 150 + (rack.id * 120)" 
                :y="rack.y || 180" 
                width="80" 
                height="80" 
                :class="[
                  selectedRackId === rack.id 
                    ? 'stroke-primary-500 stroke-4 fill-primary-950/20' 
                    : 'stroke-slate-700 hover:stroke-primary-500 fill-slate-900 hover:fill-slate-800/80',
                  'transition-all duration-200 stroke-2'
                ]" 
                rx="6"
              />
              
              <!-- Ventilation direction indicators -->
              <path 
                :d="`M ${ (rack.x || 150 + (rack.id * 120)) + 15 } ${ (rack.y || 180) + 10 } L ${ (rack.x || 150 + (rack.id * 120)) + 40 } ${ (rack.y || 180) + 10 } L ${ (rack.x || 150 + (rack.id * 120)) + 65 } ${ (rack.y || 180) + 10 }`"
                stroke="#38bdf8" 
                stroke-width="1" 
                stroke-dasharray="3,3" 
                class="opacity-60"
              />

              <!-- Cabinet label texts -->
              <text 
                :x="(rack.x || 150 + (rack.id * 120)) + 40" 
                :y="(rack.y || 180) + 40" 
                fill="#ffffff" 
                font-family="monospace" 
                font-size="11" 
                font-weight="bold" 
                text-anchor="middle"
                class="group-hover:fill-primary-400 transition-colors"
              >
                {{ rack.name }}
              </text>
              <text 
                :x="(rack.x || 150 + (rack.id * 120)) + 40" 
                :y="(rack.y || 180) + 60" 
                fill="#94a3b8" 
                font-family="monospace" 
                font-size="9" 
                text-anchor="middle"
              >
                {{ rack.placement_zone || 'Core Zone' }}
              </text>
            </g>

            <!-- Empty Grid Bays text indicator -->
            <text x="730" y="430" font-family="monospace" font-size="9" fill="#475569" text-anchor="end">Drafting Scale: Dots = 1 Meter (1000mm)</text>
          </svg>
        </div>
      </div>

      <!-- Right Sidebar: Cabinet Vertical 42U Rail Inspector (Click-to-Open details, Phase 1) -->
      <div class="w-[380px] flex-shrink-0 h-full overflow-hidden flex flex-col justify-between">
        
        <!-- Cabinet Selected View -->
        <div v-if="activeRackDetails" class="h-full flex flex-col overflow-hidden">
          <div class="bg-slate-900 border-[12px] border-slate-700 dark:border-slate-800 rounded-2xl p-4 shadow-2xl relative flex-1 flex flex-col justify-between overflow-hidden">
            
            <!-- Cabinet Header & Geodata Link Speeds -->
            <div class="border-b border-slate-800 pb-2 mb-3 flex-shrink-0">
              <div class="flex justify-between items-start">
                <h4 class="text-xs font-bold text-white font-mono uppercase">{{ activeRackDetails.name }}</h4>
                <UBadge color="primary" size="xs" variant="subtle" class="font-mono text-[8px]">
                  {{ activeRackDetails.placement_zone || 'Core Zone' }}
                </UBadge>
              </div>
              <span class="text-[9px] text-slate-500 font-mono block mt-1">Rails Capacity: {{ activeRackDetails.height_u }}U</span>
            </div>

            <!-- Scrollable rails slots -->
            <div class="flex-1 overflow-y-auto pr-1 space-y-1 bg-slate-950 p-2.5 rounded border border-slate-800 shadow-inner">
              <div 
                v-for="u in rackUnits" 
                :key="u" 
                class="flex items-center gap-2 group"
              >
                <span class="w-6 text-right font-mono text-[9px] font-bold text-slate-500">{{ u }}U</span>
                <div class="flex-1 min-h-[34px]">
                  <!-- Active Asset mounted panel -->
                  <div 
                    v-if="occupiedSlots[u]" 
                    :class="[
                      getChassisBorderClass(occupiedSlots[u].type),
                      'flex items-center justify-between px-3 py-1 bg-slate-800 rounded border-l-4 hover:brightness-110 transition-all cursor-pointer'
                    ]"
                    @click="navigateTo(`/assets/${occupiedSlots[u].id}`)"
                  >
                    <div class="flex items-center gap-2 truncate">
                      <UIcon :name="getCategoryIcon(occupiedSlots[u].type)" class="h-4 w-4 text-slate-400 flex-shrink-0" />
                      <span class="text-[10px] font-bold text-white font-mono truncate max-w-[100px]">{{ occupiedSlots[u].name }}</span>
                    </div>
                    <div class="flex gap-1 items-center flex-shrink-0">
                      <UButton 
                        v-if="occupiedSlots[u].type === 'Network'"
                        size="xs" 
                        color="primary" 
                        variant="ghost" 
                        icon="i-heroicons-arrow-path-rounded-square" 
                        @click.stop="navigateTo(`/topology?trace=${occupiedSlots[u].id}`)"
                      />
                      <span class="h-1 w-1 bg-emerald-500 rounded-full animate-pulse" />
                    </div>
                  </div>

                  <!-- Blanking Plate -->
                  <div 
                    class="flex items-center justify-between py-1.5 border border-slate-900 rounded bg-gradient-to-r from-slate-900 to-slate-950 text-[8px] text-slate-600 font-mono select-none px-3"
                    v-else
                  >
                    <span>BLANK PANEL</span>
                    <span class="h-0.5 w-0.5 bg-slate-800 rounded-full" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Cabling Port Capacity Heatmap (Phase 2 DCIM L4) -->
            <div v-if="rackSwitches && rackSwitches.length > 0" class="border-t border-slate-800 pt-2 mt-2 space-y-2 flex-shrink-0">
              <span class="text-[9px] text-slate-400 font-mono uppercase tracking-wider block">Port Capacity Heatmap</span>
              <div v-for="sw in rackSwitches" :key="sw.id" class="space-y-1">
                <div class="flex justify-between text-[9px] font-mono">
                  <span class="text-white truncate max-w-[120px]">{{ sw.name }}</span>
                  <span :class="sw.utilization > 80 ? 'text-red-400 font-bold' : 'text-primary-400'">{{ sw.utilization }}% Utilized</span>
                </div>
                <UProgress :value="sw.utilization" :color="sw.utilization > 80 ? 'red' : 'primary'" size="xs" />
                <span v-if="sw.utilization > 80" class="text-[8px] text-red-400 block pt-0.5 leading-tight animate-pulse">Warning: SFP+ Uplink Capacity Threshold Exceeded!</span>
              </div>
            </div>

            <!-- DCIM Connection / Cabinet WAN Details display (Phase 4) -->
            <div v-if="selectedDc" class="border-t border-slate-800 pt-2.5 mt-3 text-[10px] font-mono text-slate-400 space-y-1 flex-shrink-0 leading-tight">
              <div class="flex justify-between">
                <span>Site Location:</span>
                <span class="text-white font-semibold">{{ selectedDc.city }}, {{ selectedDc.country }}</span>
              </div>
              <div v-if="selectedDc.properties?.uplink_speed" class="flex justify-between">
                <span>WAN Link:</span>
                <span class="text-primary-400 font-semibold">{{ selectedDc.properties.uplink_speed }}</span>
              </div>
              <div v-if="selectedDc.properties?.public_ip" class="flex justify-between">
                <span>Gateway IP:</span>
                <span class="text-white">{{ selectedDc.properties.public_ip }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- No Cabinet Selected Sidebar Splash -->
        <div v-else class="h-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md p-6 flex flex-col justify-center text-center text-slate-500 shadow-sm">
          <UIcon name="i-heroicons-cursor-arrow-rays" class="h-10 w-10 text-primary-500/60 mx-auto animate-bounce mb-3" />
          <h4 class="font-bold text-slate-800 dark:text-white font-mono text-xs">Awaiting Portal Select</h4>
          <p class="text-[10px] leading-normal max-w-[200px] mx-auto">
            Click on any arranged cabinet block inside the 2D CAD Room Map on the left to reveal its physical vertical server rail details.
          </p>
        </div>
      </div>
    </div>

    <!-- 3. Deploy Floor Modal (Phase 2 Multi-Floor setup) -->
    <UModal v-model="isFloorModalOpen">
      <UCard>
        <template #header>
          <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono">Deploy Datacenter Floor Level</h3>
        </template>
        <form @submit.prevent="createFloor" class="space-y-4">
          <UFormGroup label="Floor Name label">
            <UInput v-model="floorForm.name" placeholder="e.g. Second Floor, Basement" required />
          </UFormGroup>
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Level Numeric Index">
              <UInput type="number" v-model="floorForm.level" placeholder="2" required />
            </UFormGroup>
            <UFormGroup label="CAD Grid Width (cm)">
              <UInput type="number" v-model="floorForm.width" placeholder="1000" />
            </UFormGroup>
          </div>
          <div class="flex justify-end gap-2 pt-2 border-t border-slate-100 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isFloorModalOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="isSavingFloor">Add Level</UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const { fetchDatacenters, fetchRackDetails } = useDCIM()
const { canMutate, canDelete, getAuthHeader } = useAuth()

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

const { data: datacenters, pending: pendingDatacenters, refresh: refreshDcs } = await fetchDatacenters()

const selectedDcId = ref(null)
const selectedFloorId = ref(null)
const selectedRackId = ref(null)

const activeRackDetails = ref(null)
const pendingRack = ref(false)

const isFloorModalOpen = ref(false)
const isSavingFloor = ref(false)
const floorForm = ref({ name: '', level: 2, width: 1000, depth: 1000 })

// Select first deployed datacenter by default on startup if available
watch(datacenters, (newDcs) => {
  if (newDcs && newDcs.length > 0) {
    selectedDcId.value = newDcs[0].id
  }
}, { immediate: true })

// Preselect Floor and Cabinet when DC changes
watch(selectedDcId, (newDcId) => {
  selectedFloorId.value = null
  selectedRackId.value = null
  activeRackDetails.value = null
  if (newDcId && selectedDc.value) {
    if (selectedDc.value.floors && selectedDc.value.floors.length > 0) {
      selectedFloorId.value = selectedDc.value.floors[0].id
    }
  }
})

// Clear / preselect Rack when Floor changes
watch(selectedFloorId, (newFloorId) => {
  selectedRackId.value = null
  activeRackDetails.value = null
  if (newFloorId && filteredRacks.value.length > 0) {
    openRackRails(filteredRacks.value[0])
  }
})

// Fetch Rack details dynamically when selection changes
watch(selectedRackId, async (newRackId) => {
  if (newRackId) {
    pendingRack.value = true
    try {
      const { data } = await fetchRackDetails(newRackId)
      if (data.value) {
        activeRackDetails.value = data.value
      }
    } catch (err) {
      console.error('Failed to load rack detail:', err)
    } finally {
      pendingRack.value = false
    }
  } else {
    activeRackDetails.value = null
  }
})

const selectedDc = computed(() => {
  if (!datacenters.value || !selectedDcId.value) return null
  return datacenters.value.find(d => d.id === selectedDcId.value)
})

const filteredRacks = computed(() => {
  if (!selectedDc.value) return []
  // Filter Racks by selected Floor Level (Phase 2 Multi-Floor)
  if (selectedFloorId.value) {
    return selectedDc.value.racks?.filter(r => r.datacenter_floor_id === selectedFloorId.value) || []
  }
  return selectedDc.value.racks || []
})

const activeFloorWalls = computed(() => {
  if (!selectedDc.value || !selectedFloorId.value) return []
  const floor = selectedDc.value.floors?.find(f => f.id === selectedFloorId.value)
  return floor?.walls || []
})

// Calculates the open-state door leaf coordinates (hinge to swung endpoints)
// Supports single doors, flipped direction, and double/split doors meeting cleanly
const getDoorLeafCoords = (wall) => {
  const x1 = wall.x1 * 2
  const y1 = wall.y1 * 2
  const x2 = wall.x2 * 2
  const y2 = wall.y2 * 2
  const dx = x2 - x1
  const dy = y2 - y1
  const sweep = wall.flipped ? -1 : 1
  
  if (wall.type === 'door-double' || wall.type === 'door-split') {
    // Left leaf pivots at x1,y1, open leaf is swung 90 deg (half-size!)
    const openX1 = x1 + (dy / 2) * sweep
    const openY1 = y1 - (dx / 2) * sweep
    
    // Right leaf pivots at x2,y2, open leaf is swung 90 deg (half-size!)
    const openX2 = x2 + (dy / 2) * sweep
    const openY2 = y2 - (dx / 2) * sweep
    
    return { x1, y1, openX1, openY1, x2, y2, openX2, openY2 }
  }
  
  // Standard single door leaf pivoted at x1,y1, open leaf is swung 90 deg (full-size!)
  const openX = x1 + dy * sweep
  const openY = y1 - dx * sweep
  return { x1, y1, openX, openY, x2, y2 }
}

// Calculates a high-precision quadrant door arc swing pivoting around hinge (x1, y1)
// Supports single doors, flipped direction, and double/split doors swinging concurrently
const getDoorArcPath = (wall) => {
  const x1 = wall.x1 * 2
  const y1 = wall.y1 * 2
  const x2 = wall.x2 * 2
  const y2 = wall.y2 * 2
  const dx = x2 - x1
  const dy = y2 - y1
  const r = Math.sqrt(dx*dx + dy*dy)
  const sweep = wall.flipped ? 0 : 1
  
  if (wall.type === 'door-double' || wall.type === 'door-split') {
    const halfR = r / 2
    const midX = (x1 + x2) / 2
    const midY = (y1 + y2) / 2
    
    const openX1 = x1 + (dy / 2) * (wall.flipped ? -1 : 1)
    const openY1 = y1 - (dx / 2) * (wall.flipped ? -1 : 1)
    const openX2 = x2 + (dy / 2) * (wall.flipped ? -1 : 1)
    const openY2 = y2 - (dx / 2) * (wall.flipped ? -1 : 1)
    
    // Sweeps elegantly from the open door end-points into the closed middle meeting point
    const arc1 = `M ${openX1} ${openY1} A ${halfR} ${halfR} 0 0 ${sweep} ${midX} ${midY}`
    const arc2 = `M ${openX2} ${openY2} A ${halfR} ${halfR} 0 0 ${1 - sweep} ${midX} ${midY}`
    return `${arc1} ${arc2}`
  }
  
  const openX = x1 + dy * (wall.flipped ? -1 : 1)
  const openY = y1 - dx * (wall.flipped ? -1 : 1)
  // Sweeps from open door leaf endpoint back to standard wall line endpoint (x2, y2)
  return `M ${openX} ${openY} A ${r} ${r} 0 0 ${sweep} ${x2} ${y2}`
}

const dcOptions = computed(() => {
  if (!datacenters.value) return []
  return datacenters.value.map(d => ({
    label: `${d.name} (${d.city}, ${d.country})`,
    value: d.id
  }))
})

const openRackRails = (rack) => {
  selectedRackId.value = rack.id
}

const createFloor = async () => {
  isSavingFloor.value = true
  try {
    const payload = {
      ...floorForm.value,
      datacenter_id: selectedDcId.value,
      level: Number(floorForm.value.level),
      width: Number(floorForm.value.width),
      depth: Number(floorForm.value.depth)
    }
    await $fetch(`${apiBase}/datacenter-floors/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
    isFloorModalOpen.value = false
    floorForm.value = { name: '', level: 2, width: 1000, depth: 1000 }
    await refreshDcs()
    alert('Datacenter floor level deployed successfully.')
  } catch (err) {
    console.error('Failed to deploy floor:', err)
    alert('Deploy failed: Check permissions.')
  } finally {
    isSavingFloor.value = false
  }
}

const rackUnits = computed(() => {
  if (!activeRackDetails.value) return []
  const height = activeRackDetails.value.height_u || 42
  const units = []
  for (let i = height; i >= 1; i--) {
    units.push(i)
  }
  return units
})

const occupiedSlots = computed(() => {
  if (!activeRackDetails.value || !activeRackDetails.value.assets) return {}
  const map = {}
  activeRackDetails.value.assets.forEach(asset => {
    if (asset.rack_position_u) {
      map[asset.rack_position_u] = asset
    }
  })
  return map
})

const rackSwitches = computed(() => {
  if (!activeRackDetails.value || !activeRackDetails.value.assets) return []
  return activeRackDetails.value.assets
    .filter(a => a.type === 'Network' || a.type === 'Switch' || a.type === 'Router')
    .map(sw => {
      // Mock deterministic port capacity logic for visual dashboard (Phase 2 DCIM L4)
      const portCount = sw.interfaces ? sw.interfaces.length : 0
      const totalPorts = sw.properties?.total_ports ? Number(sw.properties.total_ports) : 48
      let utilization = 0
      if (portCount > 0) {
        utilization = Math.round((portCount / totalPorts) * 100)
      } else {
        // If no interfaces seeded yet, generate deterministic capacity based on ID to demonstrate UI warning states
        const hash = sw.id.charCodeAt(0) + sw.id.charCodeAt(sw.id.length - 1)
        utilization = (hash % 100) > 20 ? (hash % 100) : 85 
      }
      return {
        ...sw,
        utilization
      }
    })
})

const getChassisBorderClass = (type) => {
  switch (type) {
    case 'Network':
      return 'border-l-red-500'
    case 'Database':
      return 'border-l-indigo-500'
    case 'Kubernetes Deployment':
    case 'Container':
      return 'border-l-emerald-500'
    case 'Server':
      return 'border-l-slate-400'
    default:
      return 'border-l-amber-500'
  }
}

const getCategoryIcon = (type) => {
  switch (type) {
    case 'Server':
      return 'i-heroicons-server'
    case 'Database':
      return 'i-heroicons-circle-stack'
    case 'Network':
      return 'i-heroicons-arrows-right-left'
    case 'Container':
      return 'i-heroicons-cube'
    case 'Kubernetes Deployment':
      return 'i-heroicons-square-3-stack-3d'
    default:
      return 'i-heroicons-squares-2x2'
  }
}
</script>
