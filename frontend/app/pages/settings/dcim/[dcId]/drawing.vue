<template>
  <div class="flex flex-col h-screen bg-slate-950 text-slate-100 overflow-hidden font-sans">
    
    <!-- Top Bar: Workspace Controls -->
    <header class="flex items-center justify-between px-6 py-4 bg-slate-900 border-b border-slate-800 flex-shrink-0 shadow-md">
      <div class="flex items-center gap-4">
        <UButton 
          icon="i-heroicons-arrow-left" 
          color="gray" 
          variant="ghost" 
          size="sm" 
          @click="leaveWorkspace"
        >
          Back to DCIM Settings
        </UButton>
        <div class="h-6 w-px bg-slate-800" />
        <div>
          <h1 class="text-sm font-bold text-white font-mono flex items-center gap-2">
            <UIcon name="i-heroicons-wrench-screwdriver" class="text-primary-500 h-4 w-4" />
            Northstar CAD Blueprint Workspace
          </h1>
          <p class="text-[11px] text-slate-400 mt-0.5 font-sans" v-if="activeDatacenter">
            {{ activeDatacenter.name }} &mdash; {{ activeFloor?.name || 'Loading Floor...' }}
          </p>
        </div>
      </div>

      <!-- Saving Status & Buttons -->
      <div class="flex items-center gap-4">
        <!-- Status Indicator -->
        <div class="flex items-center gap-2 px-3 py-1 bg-slate-850 rounded-full border border-slate-800 text-[10px] font-mono">
          <span 
            class="h-2 w-2 rounded-full" 
            :class="[
              isSaving ? 'bg-blue-500 animate-pulse' : (isDirty ? 'bg-amber-500 animate-ping' : 'bg-emerald-500')
            ]"
          />
          <span v-if="isSaving" class="text-blue-400">Saving changes...</span>
          <span v-else-if="isDirty" class="text-amber-400">Unsaved changes (idle auto-save active)</span>
          <span v-else class="text-emerald-400">All changes saved to GORM</span>
        </div>

        <UButton 
          color="primary" 
          size="xs" 
          icon="i-heroicons-cloud-arrow-up" 
          :loading="isSaving" 
          :disabled="!isDirty"
          @click="saveChanges"
          class="font-bold"
        >
          Save Blueprint
        </UButton>
      </div>
    </header>

    <!-- Main Workspace Area -->
    <div class="flex flex-grow min-h-0 relative">
      
      <!-- Left Panel: Mode Selector -->
      <aside class="bg-slate-900 flex flex-col justify-between flex-shrink-0 w-64 border-r border-slate-800 p-4">
        <div class="space-y-6">
          <div class="space-y-2">
            <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider font-mono">Drafting Modes</span>
            <div class="space-y-1">
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'racks' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-server" 
                @click="cadToolMode = 'racks'"
                class="justify-start font-medium"
              >
                Cabinet Mode
              </UButton>
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'walls' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-pencil" 
                @click="cadToolMode = 'walls'"
                class="justify-start font-medium"
              >
                Draw Walls Mode
              </UButton>
            </div>
          </div>

          <div class="space-y-2">
            <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider font-mono">Active Level</span>
            <USelect 
              v-model="selectedFloorId" 
              :options="activeDatacenter?.floors?.map(f => ({ label: f.name, value: f.id })) || []" 
              size="xs"
              class="w-full"
            />
          </div>

          <div class="h-px bg-slate-800" />

          <!-- Instructions Card -->
          <div class="bg-slate-950 p-4 rounded-md border border-slate-800 space-y-3">
            <h3 class="text-xs font-bold text-white font-mono flex items-center gap-1.5">
              <UIcon name="i-heroicons-information-circle" class="text-primary-500 h-4 w-4" />
              Quick Guide
            </h3>
            <div class="space-y-2 text-[11px] text-slate-400 leading-normal" v-if="cadToolMode === 'racks'">
              <p>📍 <strong>Cabinet Mode:</strong> Drag cabinet frames directly onto the dotted grid canvas.</p>
              <p>💾 Snap coordinates are automatically buffered in memory and committed on idle.</p>
            </div>
            <div class="space-y-2 text-[11px] text-slate-400 leading-normal" v-else>
              <p>✏️ <strong>Wall Mode:</strong> Click and drag anywhere on the grid canvas to draw a wall barrier.</p>
              <p>❌ Click on any existing wall segment on the canvas to delete it instantly.</p>
            </div>
          </div>
        </div>

        <div class="text-[10px] text-slate-600 font-mono text-center">
          Northstar v1.0.0 &middot; astrona.io
        </div>
      </aside>

      <!-- Center Panel: The Expanded SVG Canvas -->
      <main class="flex-grow bg-slate-950 p-6 flex items-center justify-center overflow-auto relative">
        <!-- Floating Zoom Controls (Phase 2 CAD Draw Updates) -->
        <div class="absolute top-4 right-4 z-20 flex items-center gap-1.5 bg-slate-900/90 backdrop-blur border border-slate-800 p-1.5 rounded-md shadow-lg select-none">
          <UButton 
            icon="i-heroicons-minus" 
            color="gray" 
            variant="ghost" 
            size="xs" 
            @click="decreaseZoom"
            :disabled="zoomScale <= 0.5"
            class="hover:bg-slate-800 text-slate-300"
          />
          <span class="text-[10px] font-mono text-slate-400 w-12 text-center select-none font-bold">
            {{ Math.round(zoomScale * 100) }}%
          </span>
          <UButton 
            icon="i-heroicons-plus" 
            color="gray" 
            variant="ghost" 
            size="xs" 
            @click="increaseZoom"
            :disabled="zoomScale >= 2.5"
            class="hover:bg-slate-800 text-slate-300"
          />
          <div class="h-4 w-px bg-slate-800 mx-1" />
          <UButton 
            icon="i-heroicons-arrows-pointing-out" 
            color="gray" 
            variant="ghost" 
            size="xs" 
            @click="resetZoom"
            title="Reset to 100%"
            class="hover:bg-slate-800 text-slate-300"
          />
        </div>

        <div 
          class="bg-slate-900 rounded-lg border-2 border-slate-800 shadow-2xl relative overflow-hidden transition-all duration-150 flex-shrink-0"
          :style="{
            width: '100%',
            height: '100%',
            maxWidth: '1200px',
            maxHeight: '750px',
            transform: `scale(${zoomScale})`,
            transformOrigin: 'center center'
          }"
        >
          <!-- Grid Canvas SVG (Double scale for precision drawing workspace) -->
          <svg 
            ref="canvasSvg"
            width="100%" 
            height="100%" 
            viewBox="0 0 800 500"
            @mousedown="onSvgCanvasMouseDown"
            @mousemove="onDesignerMouseMove" 
            @mouseup="onDesignerMouseUp" 
            @mouseleave="onDesignerMouseUp"
          >
            <defs>
              <pattern id="workspace-dots" width="40" height="40" patternUnits="userSpaceOnUse">
                <circle cx="2" cy="2" r="1.5" fill="#475569" opacity="0.6" />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#workspace-dots)" />
            
            <!-- Render custom drawn wall segments -->
            <line 
              v-for="wall in localWalls" 
              :key="wall.id" 
              :x1="wall.x1 * 2" 
              :y1="wall.y1 * 2" 
              :x2="wall.x2 * 2" 
              :y2="wall.y2 * 2" 
              stroke="#64748b" 
              stroke-width="8" 
              stroke-linecap="round" 
              class="cursor-pointer hover:stroke-red-500 hover:stroke-6 transition-all"
              @click="removeLocalWall(wall.id)"
            />

            <!-- Active dragging wall draft segment preview -->
            <line 
              v-if="wallDraftStart && wallDraftCurrent" 
              :x1="wallDraftStart.x * 2" 
              :y1="wallDraftStart.y * 2" 
              :x2="wallDraftCurrent.x * 2" 
              :y2="wallDraftCurrent.y * 2" 
              stroke="#0ea5e9" 
              stroke-width="6" 
              stroke-dasharray="6" 
              class="animate-pulse"
            />

            <!-- Render active racks -->
            <g 
              v-for="rack in localRacks" 
              :key="rack.id" 
              class="cursor-pointer" 
              @mousedown="onDesignerMouseDown(rack, $event)"
            >
              <rect 
                :x="rack.x || (20 + (rack.id * 70))" 
                :y="rack.y || 40" 
                width="60" 
                height="60" 
                :class="draggedRackId === rack.id ? 'stroke-primary-500 stroke-2 fill-primary-950/40' : 'stroke-slate-600 fill-slate-800 hover:stroke-primary-400'"
                rx="6"
              />
              <text 
                :x="(rack.x || (20 + (rack.id * 70))) + 30" 
                :y="(rack.y || 40) + 36" 
                fill="#fff" 
                font-size="10" 
                font-family="monospace" 
                text-anchor="middle"
                class="pointer-events-none select-none font-bold"
              >
                R-{{ rack.name }}
              </text>
            </g>
          </svg>
        </div>
      </main>

      <!-- Right Panel: Items Legend & Stats -->
      <aside class="w-80 bg-slate-900 border-l border-slate-800 p-4 flex flex-col gap-6 flex-shrink-0 overflow-y-auto">
        <div class="space-y-2">
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider font-mono">Location Metrics</span>
          <div class="p-3 bg-slate-950 rounded-md border border-slate-800 space-y-2 text-xs">
            <div class="flex justify-between">
              <span class="text-slate-400">Total Walls:</span>
              <span class="font-mono text-white font-bold">{{ localWalls.length }} segments</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Deployed Cabinets:</span>
              <span class="font-mono text-white font-bold">{{ localRacks.length }} frames</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Unsaved Changes:</span>
              <span class="font-mono font-bold" :class="isDirty ? 'text-amber-500' : 'text-emerald-500'">
                {{ isDirty ? 'Yes (Pending Auto-Save)' : 'Clean' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Deploy Cabinet Enclosure Panel (Phase 2 CAD Draw Updates) -->
        <div class="space-y-3 p-4 bg-slate-950 border border-slate-800 rounded-md">
          <span class="text-[10px] font-bold text-primary-500 uppercase tracking-wider font-mono block">Deploy Cabinet Enclosure</span>
          <form @submit.prevent="deployNewCabinet" class="space-y-3">
            <UFormGroup label="Enclosure Tag/Name" class="text-xs font-semibold">
              <UInput v-model="newCabinetForm.name" size="xs" placeholder="e.g. Rack-01" required />
            </UFormGroup>
            <UFormGroup label="Height Capacity (U)" class="text-xs font-semibold">
              <UInput v-model="newCabinetForm.height_u" type="number" min="1" max="100" size="xs" placeholder="42" required />
            </UFormGroup>
            <UButton type="submit" block size="xs" icon="i-heroicons-plus" color="primary" class="font-bold">
              Deploy Cabinet Frame
            </UButton>
          </form>
        </div>

        <div class="space-y-3 flex-grow">
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider font-mono">Blueprint Inventory</span>
          <div class="space-y-1 max-h-[300px] overflow-y-auto pr-1">
            <div 
              v-for="rack in localRacks" 
              :key="rack.id"
              class="flex items-center justify-between p-2.5 bg-slate-950 border border-slate-800 rounded-md text-xs hover:border-slate-700"
            >
              <div class="flex items-center gap-2">
                <UIcon name="i-heroicons-rectangle-stack" class="text-primary-500 h-4 w-4" />
                <span class="font-mono font-semibold">{{ rack.name }}</span>
              </div>
              <span class="text-[10px] text-slate-400 font-mono">x:{{ Math.round(rack.x) }}, y:{{ Math.round(rack.y) }}</span>
            </div>
            <div v-if="localRacks.length === 0" class="text-xs text-slate-500 italic p-3 text-center">
              No deployed cabinet frames.
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDCIM } from '#imports'
import { useAuth } from '#imports'

const route = useRoute()
const router = useRouter()
const { getAuthHeader } = useAuth()
const { fetchDatacenters } = useDCIM()

const runtimeConfig = useRuntimeConfig()
const apiBase = runtimeConfig.public.apiBase

const dcId = route.params.dcId || ''
const selectedFloorId = ref(route.query.floorId || '')

const { data: datacenters, refresh: refreshDatacenters } = await fetchDatacenters()

const activeDatacenter = computed(() => {
  if (!datacenters.value || !dcId) return null
  return datacenters.value.find(d => d.id === dcId)
})

const activeFloor = computed(() => {
  const dc = activeDatacenter.value
  if (!dc || !dc.floors || dc.floors.length === 0) return null
  if (selectedFloorId.value) {
    return dc.floors.find(f => f.id === selectedFloorId.value) || dc.floors[0]
  }
  return dc.floors[0]
})

// Local buffered state for drawing and auto-saving
const localWalls = ref([])
const localRacks = ref([])
const pendingDeletions = ref([]) // track database walls deleted locally

const newCabinetForm = ref({
  name: '',
  height_u: 42
})

const syncLocalStateFromDB = () => {
  const floor = activeFloor.value
  if (floor) {
    localWalls.value = JSON.parse(JSON.stringify(floor.walls || []))
    localRacks.value = JSON.parse(JSON.stringify(floor.racks || []))
    pendingDeletions.value = []
  }
}

// Watch activeFloor to load values when switching levels
watch(activeFloor, (newFloor) => {
  if (newFloor) {
    if (selectedFloorId.value !== newFloor.id) {
      selectedFloorId.value = newFloor.id
    }
    syncLocalStateFromDB()
  }
}, { immediate: true })

// Drawing state variables
const canvasSvg = ref(null)
const draggedRackId = ref(null)
const dragOffset = ref({ x: 0, y: 0 })
const cadToolMode = ref('racks') // 'racks' or 'walls'
const wallDraftStart = ref(null)
const wallDraftCurrent = ref(null)

// Zoom controls (Phase 2 CAD Draw Updates)
const zoomScale = ref(1.0)

const increaseZoom = () => {
  if (zoomScale.value < 2.5) {
    zoomScale.value = parseFloat((zoomScale.value + 0.1).toFixed(1))
  }
}

const decreaseZoom = () => {
  if (zoomScale.value > 0.5) {
    zoomScale.value = parseFloat((zoomScale.value - 0.1).toFixed(1))
  }
}

const resetZoom = () => {
  zoomScale.value = 1.0
}

// Auto-save & idle timers
const saveTimer = ref(null)
const isDirty = ref(false)
const isSaving = ref(false)

const markDirty = () => {
  isDirty.value = true
  resetAutoSaveTimer()
}

const resetAutoSaveTimer = () => {
  if (saveTimer.value) {
    clearTimeout(saveTimer.value)
  }
  saveTimer.value = setTimeout(() => {
    if (isDirty.value && !isSaving.value) {
      saveChanges()
    }
  }, 5000) // 5 seconds idle
}

const onDesignerMouseDown = (rack, event) => {
  if (cadToolMode.value !== 'racks') return
  draggedRackId.value = rack.id
  
  const svgElement = canvasSvg.value || event.currentTarget?.ownerSVGElement
  if (!svgElement) return
  const pt = svgElement.createSVGPoint()
  pt.x = event.clientX
  pt.y = event.clientY
  const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
  
  const currentX = rack.x || (20 + (rack.id * 70))
  const currentY = rack.y || 40
  
  dragOffset.value = {
    x: svgP.x - currentX,
    y: svgP.y - currentY
  }
}

const onSvgCanvasMouseDown = (event) => {
  if (cadToolMode.value !== 'walls' || !activeFloor.value) return
  
  const svgElement = canvasSvg.value || event.currentTarget
  if (!svgElement) return
  const pt = svgElement.createSVGPoint()
  pt.x = event.clientX
  pt.y = event.clientY
  const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
  
  // Snap coordinates in larger workspace
  const snappedX = Math.round(svgP.x / 40) * 20
  const snappedY = Math.round(svgP.y / 40) * 20
  
  wallDraftStart.value = { x: snappedX, y: snappedY }
  wallDraftCurrent.value = { x: snappedX, y: snappedY }
}

const onDesignerMouseMove = (event) => {
  const svgElement = canvasSvg.value || event.currentTarget
  if (!svgElement) return

  // Dragging cabinets
  if (draggedRackId.value && cadToolMode.value === 'racks') {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    const newX = svgP.x - dragOffset.value.x
    const newY = svgP.y - dragOffset.value.y
    
    // Snap to 40-unit grid boundaries
    const snappedX = Math.round(newX / 40) * 40
    const snappedY = Math.round(newY / 40) * 40
    
    const rack = localRacks.value.find(r => r.id === draggedRackId.value)
    if (rack) {
      if (rack.x !== snappedX || rack.y !== snappedY) {
        rack.x = snappedX
        rack.y = snappedY
        markDirty()
      }
    }
    return
  }
  
  // Drafting wall preview
  if (wallDraftStart.value && cadToolMode.value === 'walls') {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    const snappedX = Math.round(svgP.x / 40) * 20
    const snappedY = Math.round(svgP.y / 40) * 20
    
    wallDraftCurrent.value = { x: snappedX, y: snappedY }
  }
}

const onDesignerMouseUp = () => {
  // Release cabinet drag
  if (draggedRackId.value && cadToolMode.value === 'racks') {
    draggedRackId.value = null
    return
  }
  
  // Release wall drawing (Buffer to memory)
  if (wallDraftStart.value && wallDraftCurrent.value && cadToolMode.value === 'walls' && activeFloor.value) {
    const x1 = wallDraftStart.value.x
    const y1 = wallDraftStart.value.y
    const x2 = wallDraftCurrent.value.x
    const y2 = wallDraftCurrent.value.y
    
    wallDraftStart.value = null
    wallDraftCurrent.value = null
    
    if (x1 === x2 && y1 === y2) return
    
    // Insert new wall locally in list
    localWalls.value.push({
      id: 'temp_' + Date.now() + '_' + Math.floor(Math.random() * 1000),
      datacenter_floor_id: activeFloor.value.id,
      x1: x1 / 2,
      y1: y1 / 2,
      x2: x2 / 2,
      y2: y2 / 2
    })
    markDirty()
  }
}

const removeLocalWall = (wallId) => {
  localWalls.value = localWalls.value.filter(w => w.id !== wallId)
  if (!String(wallId).startsWith('temp_')) {
    pendingDeletions.value.push(wallId)
  }
  markDirty()
}

const saveChanges = async () => {
  if (!isDirty.value || isSaving.value || !activeFloor.value) return
  isSaving.value = true
  if (saveTimer.value) {
    clearTimeout(saveTimer.value)
  }
  
  try {
    // 1. Commit deletions to the database
    for (const wallId of pendingDeletions.value) {
      await $fetch(`${apiBase}/dcim/walls/${wallId}`, {
        method: 'DELETE',
        headers: getAuthHeader()
      })
    }
    pendingDeletions.value = []

    // 2. Commit new wall creations
    for (const wall of localWalls.value) {
      if (String(wall.id).startsWith('temp_')) {
        await $fetch(`${apiBase}/dcim/walls`, {
          method: 'POST',
          body: {
            datacenter_floor_id: activeFloor.value.id,
            x1: wall.x1,
            y1: wall.y1,
            x2: wall.x2,
            y2: wall.y2
          },
          headers: getAuthHeader()
        })
      }
    }

    // 3. Commit cabinet coordinate updates
    for (const rack of localRacks.value) {
      await $fetch(`${apiBase}/racks/${rack.id}/coordinates`, {
        method: 'PUT',
        body: { x: rack.x, y: rack.y },
        headers: getAuthHeader()
      })
    }

    // 4. Fully refresh and recheck state
    await refreshDatacenters()
    syncLocalStateFromDB()
    isDirty.value = false
    console.log('Saved CAD blueprint changes successfully')
  } catch (err) {
    console.error('Failed to auto-save blueprint changes:', err)
  } finally {
    isSaving.value = false
  }
}

const deployNewCabinet = async () => {
  if (!activeFloor.value) return
  try {
    const payload = {
      datacenter_id: dcId,
      name: newCabinetForm.value.name,
      height_u: Number(newCabinetForm.value.height_u),
      x: 40,
      y: 40
    }
    await $fetch(`${apiBase}/racks/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
    newCabinetForm.value = { name: '', height_u: 42 }
    await refreshDatacenters()
    syncLocalStateFromDB()
    console.log('Cabinet deployed successfully in drawing tool')
  } catch (err) {
    console.error('Failed to deploy cabinet from workspace:', err)
  }
}

const leaveWorkspace = async () => {
  if (isDirty.value) {
    const confirmLeave = confirm('You have unsaved blueprint changes. Save now before returning to settings?')
    if (confirmLeave) {
      await saveChanges()
    }
  }
  navigateTo(`/settings/dcim`)
}

// Clear auto-save timer when exiting drawing page
onBeforeUnmount(() => {
  if (saveTimer.value) {
    clearTimeout(saveTimer.value)
  }
})
</script>