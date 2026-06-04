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
              <!-- Select / Pointer Mode -->
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'select' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-cursor-arrow-rays" 
                @click="cadToolMode = 'select'; selectedElementId = null"
                class="justify-start font-medium"
              >
                Select &amp; Edit
              </UButton>
              <!-- Wall Line Mode -->
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'walls' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-pencil" 
                @click="cadToolMode = 'walls'"
                class="justify-start font-medium"
              >
                Wall Line
              </UButton>
              <!-- Wall Enclosure Mode -->
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'wall-rect' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-square-3-stack-3d" 
                @click="cadToolMode = 'wall-rect'"
                class="justify-start font-medium"
              >
                Wall Enclosure
              </UButton>
              <!-- Door swing Mode -->
              <UButton 
                block
                size="sm" 
                :color="cadToolMode === 'door' ? 'primary' : 'gray'" 
                variant="soft" 
                icon="i-heroicons-arrow-right-start-on-rectangle" 
                @click="cadToolMode = 'door'"
                class="justify-start font-medium"
              >
                Door Swing
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
              CAD Quick Guide
            </h3>
            <div class="space-y-2 text-[11px] text-slate-400 leading-normal">
              <p>⚡ <strong>Dragging:</strong> Click and hold any cabinet on the canvas to move it smoothly along the grid anytime.</p>
              <p>📐 <strong>Resizing:</strong> Left-click any wall or door to show blue sizing handles; drag those handles to stretch/resize!</p>
              <p>⚙️ <strong>Configure:</strong> Double-click any element to edit its thickness/type, or right-click to open actions menu.</p>
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
              <!-- Dotted subgrid with 10px spacing (3 small dots between meters) -->
              <pattern id="workspace-subgrid" width="10" height="10" patternUnits="userSpaceOnUse">
                <circle cx="0" cy="0" r="0.8" fill="#475569" opacity="0.4" />
              </pattern>
              <!-- Major 1-Meter Grid with 40px spacing and big corner dots -->
              <pattern id="workspace-dots" width="40" height="40" patternUnits="userSpaceOnUse">
                <rect width="40" height="40" fill="url(#workspace-subgrid)" />
                <circle cx="0" cy="0" r="2.2" fill="#64748b" opacity="0.9" />
              </pattern>

              <!-- Dynamic Wall Cutout Mask (Creates clean transparent openings for doors/windows) -->
              <mask id="wall-mask">
                <!-- White area allows standard walls to paint normally -->
                <rect width="100%" height="100%" fill="#ffffff" />
                <!-- Black lines represent door segments, cutting out doorways from the walls -->
                <line 
                  v-for="door in localWalls.filter(w => w.type === 'door' || w.type === 'door-double')" 
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
            <rect width="100%" height="100%" fill="url(#workspace-dots)" />
            
            <!-- Render custom drawn elements (Walls / Doors) -->
            <g v-for="wall in localWalls" :key="wall.id">
              <!-- Render Door Segment with Swing Arc (Single or Double Door) -->
              <g 
                v-if="wall.type === 'door' || wall.type === 'door-double'" 
                class="cursor-pointer group"
                @mousedown.stop="onWallMouseDown(wall, $event)"
                @dblclick="openElementSettings(wall)"
                @contextmenu.prevent="openContextMenu($event, wall, 'wall')"
              >
                <!-- If double/split door, render two half-size open door leaves -->
                <template v-if="wall.type === 'door-double'">
                  <line 
                    :x1="getDoorLeafCoords(wall).x1" :y1="getDoorLeafCoords(wall).y1" 
                    :x2="getDoorLeafCoords(wall).openX1" :y2="getDoorLeafCoords(wall).openY1" 
                    :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#f43f5e'" 
                    :stroke-width="(wall.thickness || 10) * 0.4" 
                    stroke-linecap="round"
                    class="hover:stroke-red-400 transition-all"
                    vector-effect="non-scaling-stroke"
                  />
                  <line 
                    :x1="getDoorLeafCoords(wall).x2" :y1="getDoorLeafCoords(wall).y2" 
                    :x2="getDoorLeafCoords(wall).openX2" :y2="getDoorLeafCoords(wall).openY2" 
                    :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#f43f5e'" 
                    :stroke-width="(wall.thickness || 10) * 0.4" 
                    stroke-linecap="round"
                    class="hover:stroke-red-400 transition-all"
                    vector-effect="non-scaling-stroke"
                  />
                </template>
                <!-- Standard single open door leaf -->
                <line 
                  v-else
                  :x1="getDoorLeafCoords(wall).x1" :y1="getDoorLeafCoords(wall).y1" 
                  :x2="getDoorLeafCoords(wall).openX" :y2="getDoorLeafCoords(wall).openY" 
                  :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#f43f5e'" 
                  :stroke-width="(wall.thickness || 10) * 0.4" 
                  stroke-linecap="round"
                  class="hover:stroke-red-400 transition-all"
                  vector-effect="non-scaling-stroke"
                />
                
                <!-- Shared Arc Swing Paths -->
                <path 
                  :d="getDoorArcPath(wall)" 
                  fill="none" 
                  :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#f43f5e'" 
                  stroke-width="1.5" 
                  stroke-dasharray="3" 
                  class="opacity-70 pointer-events-none"
                />
              </g>

              <!-- Render Room Enclosure (Rectangular Room Enclosure) -->
              <rect 
                v-else-if="wall.type === 'wall-rect'"
                :x="Math.min(wall.x1, wall.x2) * 2" 
                :y="Math.min(wall.y1, wall.y2) * 2" 
                :width="Math.abs(wall.x2 - wall.x1) * 2" 
                :height="Math.abs(wall.y2 - wall.y1) * 2" 
                fill="none"
                :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#64748b'" 
                :stroke-width="(wall.thickness || 20) * 0.4" 
                stroke-linejoin="round"
                class="cursor-pointer hover:stroke-primary-400 transition-all"
                mask="url(#wall-mask)"
                vector-effect="non-scaling-stroke"
                @mousedown.stop="onWallMouseDown(wall, $event)"
                @dblclick="openElementSettings(wall)"
                @contextmenu.prevent="openContextMenu($event, wall, 'wall')"
              />

              <!-- Render Standard Wall Segment -->
              <line 
                v-else
                :x1="wall.x1 * 2" 
                :y1="wall.y1 * 2" 
                :x2="wall.x2 * 2" 
                :y2="wall.y2 * 2" 
                :stroke="selectedElementId === wall.id ? '#0ea5e9' : '#64748b'" 
                :stroke-width="(wall.thickness || 20) * 0.4" 
                stroke-linecap="round" 
                class="cursor-pointer hover:stroke-primary-400 transition-all"
                mask="url(#wall-mask)"
                vector-effect="non-scaling-stroke"
                @mousedown.stop="onWallMouseDown(wall, $event)"
                @dblclick="openElementSettings(wall)"
                @contextmenu.prevent="openContextMenu($event, wall, 'wall')"
              />

              <!-- Interactive Sizing Handles (Rendered only when selected) -->
              <g v-if="selectedElementId === wall.id">
                <circle 
                  :cx="wall.x1 * 2" 
                  :cy="wall.y1 * 2" 
                  r="6" 
                  fill="#0ea5e9" 
                  stroke="#fff" 
                  stroke-width="2" 
                  class="cursor-pointer hover:fill-[#38bdf8] transition-all"
                  @mousedown.stop="startResizeHandle(wall, 'start', $event)"
                />
                <circle 
                  :cx="wall.x2 * 2" 
                  :cy="wall.y2 * 2" 
                  r="6" 
                  fill="#0ea5e9" 
                  stroke="#fff" 
                  stroke-width="2" 
                  class="cursor-pointer hover:fill-[#38bdf8] transition-all"
                  @mousedown.stop="startResizeHandle(wall, 'end', $event)"
                />
              </g>
            </g>

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
              v-for="rack in localRacks.filter(r => r.x > 0 || r.y > 0)" 
              :key="rack.id" 
              class="cursor-pointer" 
              @mousedown="onDesignerMouseDown(rack, $event)"
              @contextmenu.prevent="openContextMenu($event, rack, 'cabinet')"
            >
              <rect 
                :x="rack.x" 
                :y="rack.y" 
                width="60" 
                height="60" 
                :class="draggedRackId === rack.id ? 'stroke-primary-500 stroke-2 fill-primary-950/40' : 'stroke-slate-600 fill-slate-800 hover:stroke-primary-400'"
                rx="6"
              />
              <text 
                :x="rack.x + 30" 
                :y="rack.y + 36" 
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
        
        <!-- Deploy Action Button (Create Blueprint Style) -->
        <UButton 
          block 
          color="primary" 
          icon="i-heroicons-plus-circle"
          class="font-bold py-2 shadow-sm animate-pulse"
          @click="isDeployModalOpen = true"
        >
          Create Blueprint
        </UButton>

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

        <div class="space-y-4 flex-grow">
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider font-mono">Blueprint Inventory</span>
          
          <div class="space-y-2">
            <!-- Cabinets Group Title -->
            <div class="flex items-center gap-1.5 text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase font-mono border-b border-slate-800 pb-1.5">
              <UIcon name="i-heroicons-server" class="h-3.5 w-3.5" />
              <span>Cabinets</span>
            </div>
            
            <div class="space-y-1 max-h-[300px] overflow-y-auto pr-1">
              <div 
                v-for="rack in localRacks" 
                :key="rack.id"
                class="flex items-center justify-between p-2 bg-slate-950 border border-slate-800 rounded-md text-xs hover:border-slate-700"
              >
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-rectangle-stack" class="text-primary-500 h-4 w-4" />
                  <div class="flex flex-col">
                    <span class="font-mono font-semibold">Rack-{{ rack.name }}</span>
                    <span class="text-[10px] text-slate-500 font-mono">Capacity: {{ rack.height_u }}U</span>
                  </div>
                </div>
                
                <div class="flex items-center gap-1.5">
                  <template v-if="rack.x > 0 || rack.y > 0">
                    <span class="text-[10px] text-slate-400 font-mono">x:{{ Math.round(rack.x) }}, y:{{ Math.round(rack.y) }}</span>
                    <UButton 
                      size="xs" 
                      color="red" 
                      variant="ghost" 
                      icon="i-heroicons-minus-circle"
                      title="Unplace from Grid"
                      @click="unplaceCabinetFromGrid(rack)"
                    />
                  </template>
                  <template v-else>
                    <UBadge size="xs" color="gray" variant="soft" class="text-[9px] font-bold">Not Placed</UBadge>
                    <UButton 
                      size="xs" 
                      color="primary" 
                      variant="soft" 
                      icon="i-heroicons-map-pin"
                      title="Place on Grid"
                      @click="placeCabinetOnGrid(rack)"
                    />
                  </template>
                </div>
              </div>
              <div v-if="localRacks.length === 0" class="text-xs text-slate-500 italic p-3 text-center">
                No deployed cabinet frames.
              </div>
            </div>
          </div>
        </div>
      </aside>
    </div>

    <!-- Deploy Cabinet Enclosure Modal Dialog -->
    <UModal v-model="isDeployModalOpen">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
              <UIcon name="i-heroicons-plus-circle" class="text-primary-500 h-5 w-5 animate-pulse" />
              Deploy Cabinet Enclosure Frame
            </h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isDeployModalOpen = false" />
          </div>
        </template>
        
        <form @submit.prevent="deployNewCabinet" class="space-y-4">
          <UFormGroup label="Enclosure Tag/Name" help="e.g. Rack-01, Cabinet-A05">
            <UInput v-model="newCabinetForm.name" placeholder="e.g. Rack-01" required />
          </UFormGroup>
          <UFormGroup label="Height Capacity (U)" help="The standard cabinet vertical height (e.g. 42U is standard)">
            <UInput v-model="newCabinetForm.height_u" type="number" min="1" max="100" placeholder="42" required />
          </UFormGroup>
          
          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isDeployModalOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary">Deploy Frame</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <!-- Edit Element settings Modal (Wall/Door Thickness & Type) -->
    <UModal v-model="isElementSettingsOpen">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
              <UIcon name="i-heroicons-cog-6-tooth" class="text-primary-500 h-5 w-5 animate-spin" />
              Edit Blueprint Element Properties
            </h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isElementSettingsOpen = false" />
          </div>
        </template>
        
        <form @submit.prevent="saveElementSettings" class="space-y-4">
          <UFormGroup label="Element Type" help="Define if this element represents a wall, door, or window">
            <USelect 
              v-model="editingElementForm.type" 
              :options="[
                { label: 'Standard Wall', value: 'wall' },
                { label: 'Room Enclosure (Rect)', value: 'wall-rect' },
                { label: 'Single Door Swing', value: 'door' },
                { label: 'Double / Split Door', value: 'door-double' },
                { label: 'Window Opening', value: 'window' }
              ]" 
            />
          </UFormGroup>
          
          <UFormGroup v-if="editingElementForm.type.startsWith('door')" label="Door Swing Options" help="Reverse door leaf hinge side or swing direction">
            <UCheckbox v-model="editingElementForm.flipped" label="Flip Swing Direction (Invert Arc)" />
          </UFormGroup>
          
          <UFormGroup label="Thickness (cm)" help="Real-world wall thickness in centimeters (e.g. 20 is standard wall, 10 for thin door panels)">
            <UInput v-model="editingElementForm.thickness" type="number" min="1" max="100" placeholder="20" required />
          </UFormGroup>
          
          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200 dark:border-slate-800">
            <UButton color="gray" variant="ghost" @click="isElementSettingsOpen = false">Cancel</UButton>
            <UButton type="submit" color="primary">Save Changes</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <!-- Right-Click Context Menu -->
    <div 
      v-if="contextMenuState.isOpen" 
      :style="{ top: contextMenuState.y + 'px', left: contextMenuState.x + 'px' }" 
      class="fixed z-50 bg-slate-900 border border-slate-800 text-slate-200 rounded-md shadow-2xl py-1 text-xs font-mono w-44"
      @click.stop
    >
      <button 
        class="w-full text-left px-3 py-2 hover:bg-slate-800 flex items-center gap-2 text-slate-300"
        @click="configureContextTarget"
      >
        <UIcon name="i-heroicons-cog-8-tooth" class="h-4 w-4 text-primary-500" />
        <span>Configure...</span>
      </button>
      <button 
        class="w-full text-left px-3 py-2 hover:bg-slate-800 flex items-center gap-2 text-rose-400 hover:text-rose-300 border-t border-slate-800"
        @click="deleteContextTarget"
      >
        <UIcon name="i-heroicons-trash" class="h-4 w-4" />
        <span>Delete Element</span>
      </button>
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

// ==========================================
// Advanced CAD Architectural State & Helpers
// ==========================================
const selectedElementId = ref(null)
const resizingWallId = ref(null)
const resizingHandleType = ref(null) // 'start' or 'end'
const draggedWallId = ref(null)
const wallDragOffset = ref({ dx1: 0, dy1: 0, dx2: 0, dy2: 0 })

const isElementSettingsOpen = ref(false)
const editingElement = ref(null)
const editingElementForm = ref({
  thickness: 8,
  type: 'wall'
})

const contextMenuState = ref({
  isOpen: false,
  x: 0,
  y: 0,
  target: null,
  targetType: null
})

// Proximity-based magnetic snapping: Snaps to major 1-meter grid dots (40px step)
// if coordinates are within a 10cm (4px) threshold, otherwise snaps to standard 1cm steps
const snapToCanvasGrid = (val, step = 0.4, majorStep = 40, threshold = 4) => {
  const standardSnap = Math.round(val / step) * step
  const nearestMajor = Math.round(val / majorStep) * majorStep
  if (Math.abs(standardSnap - nearestMajor) <= threshold) {
    return nearestMajor
  }
  return parseFloat(standardSnap.toFixed(1))
}

// Proximity-based magnetic snapping for database scale (20px major step, 0.2px centimeter step)
const snapToDatabaseGrid = (val, step = 0.2, majorStep = 20, threshold = 2) => {
  const standardSnap = Math.round(val / step) * step
  const nearestMajor = Math.round(val / majorStep) * majorStep
  if (Math.abs(standardSnap - nearestMajor) <= threshold) {
    return nearestMajor
  }
  return parseFloat(standardSnap.toFixed(1))
}

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

const startResizeHandle = (wall, handleType, event) => {
  resizingWallId.value = wall.id
  resizingHandleType.value = handleType
  
  if (!String(wall.id).startsWith('temp_')) {
    pendingDeletions.value.push(wall.id)
    wall.id = 'temp_edited_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
  }
}

const onWallMouseDown = (wall, event) => {
  if (cadToolMode.value !== 'select') return
  selectedElementId.value = wall.id
  draggedWallId.value = wall.id
  
  if (!String(wall.id).startsWith('temp_')) {
    pendingDeletions.value.push(wall.id)
    wall.id = 'temp_edited_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
    selectedElementId.value = wall.id
    draggedWallId.value = wall.id
  }
  
  const svgElement = canvasSvg.value || event.currentTarget?.ownerSVGElement
  if (!svgElement) return
  const pt = svgElement.createSVGPoint()
  pt.x = event.clientX
  pt.y = event.clientY
  const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
  
  wallDragOffset.value = {
    dx1: svgP.x - wall.x1 * 2,
    dy1: svgP.y - wall.y1 * 2,
    dx2: svgP.x - wall.x2 * 2,
    dy2: svgP.y - wall.y2 * 2
  }
}

const openElementSettings = (wall) => {
  if (!String(wall.id).startsWith('temp_')) {
    pendingDeletions.value.push(wall.id)
    wall.id = 'temp_edited_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
  }
  editingElement.value = wall
  editingElementForm.value = {
    thickness: wall.thickness || 20,
    type: wall.type || 'wall',
    flipped: wall.flipped || false
  }
  isElementSettingsOpen.value = true
}

const saveElementSettings = () => {
  if (editingElement.value) {
    editingElement.value.thickness = Number(editingElementForm.value.thickness)
    editingElement.value.type = editingElementForm.value.type
    editingElement.value.flipped = editingElementForm.value.flipped
    markDirty()
  }
  isElementSettingsOpen.value = false
  editingElement.value = null
}

const openContextMenu = (event, target, type) => {
  contextMenuState.value = {
    isOpen: true,
    x: event.clientX,
    y: event.clientY,
    target: target,
    targetType: type
  }
}

const closeContextMenu = () => {
  contextMenuState.value.isOpen = false
}

const deleteContextTarget = () => {
  const { target, targetType } = contextMenuState.value
  if (!target) return
  if (targetType === 'wall') {
    removeLocalWall(target.id)
  } else if (targetType === 'cabinet') {
    unplaceCabinetFromGrid(target)
  }
  closeContextMenu()
}

const configureContextTarget = () => {
  const { target, targetType } = contextMenuState.value
  if (!target) return
  if (targetType === 'wall') {
    openElementSettings(target)
  } else if (targetType === 'cabinet') {
    newCabinetForm.value = {
      name: target.name,
      height_u: target.height_u
    }
    isDeployModalOpen.value = true
  }
  closeContextMenu()
}

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

const isDeployModalOpen = ref(false)
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
const cadToolMode = ref('select') // 'select' (standard pointer/dragging allowed) or 'walls' (drawing walls)
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
  const allowedModes = ['walls', 'wall-rect', 'door']
  if (!allowedModes.includes(cadToolMode.value) || !activeFloor.value) return
  
  const svgElement = canvasSvg.value || event.currentTarget
  if (!svgElement) return
  const pt = svgElement.createSVGPoint()
  pt.x = event.clientX
  pt.y = event.clientY
  const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
  
  // Snap coordinates with magnetic major-grid priority (10cm threshold, 1cm step)
  const snappedX = snapToDatabaseGrid(svgP.x / 2)
  const snappedY = snapToDatabaseGrid(svgP.y / 2)
  
  wallDraftStart.value = { x: snappedX, y: snappedY }
  wallDraftCurrent.value = { x: snappedX, y: snappedY }
}

const onDesignerMouseMove = (event) => {
  const svgElement = canvasSvg.value || event.currentTarget
  if (!svgElement) return

  // 1. Resizing Wall Segment Handles
  if (resizingWallId.value) {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    // Snap handle with magnetic major-grid priority (10cm threshold, 1cm step)
    const snappedX = snapToDatabaseGrid(svgP.x / 2)
    const snappedY = snapToDatabaseGrid(svgP.y / 2)
    
    const wall = localWalls.value.find(w => w.id === resizingWallId.value)
    if (wall) {
      if (resizingHandleType.value === 'start') {
        wall.x1 = snappedX
        wall.y1 = snappedY
      } else {
        wall.x2 = snappedX
        wall.y2 = snappedY
      }
      markDirty()
    }
    return
  }

  // 1.5 Dragging the entire wall/door element (Phase 4 CAD Upgrades)
  if (draggedWallId.value) {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    const newX1 = svgP.x - wallDragOffset.value.dx1
    const newY1 = svgP.y - wallDragOffset.value.dy1
    const newX2 = svgP.x - wallDragOffset.value.dx2
    const newY2 = svgP.y - wallDragOffset.value.dy2
    
    // Snap with magnetic major-grid priority (10cm threshold, 1cm step)
    const snappedX1 = snapToCanvasGrid(newX1)
    const snappedY1 = snapToCanvasGrid(newY1)
    const snappedX2 = snapToCanvasGrid(newX2)
    const snappedY2 = snapToCanvasGrid(newY2)
    
    const wall = localWalls.value.find(w => w.id === draggedWallId.value)
    if (wall) {
      // Scale back to 400 database coordinates (divide by 2)
      wall.x1 = snappedX1 / 2
      wall.y1 = snappedY1 / 2
      wall.x2 = snappedX2 / 2
      wall.y2 = snappedY2 / 2
      markDirty()
    }
    return
  }

  // 2. Dragging cabinets
  if (draggedRackId.value) {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    const newX = svgP.x - dragOffset.value.x
    const newY = svgP.y - dragOffset.value.y
    
    // Snap with magnetic major-grid priority (10cm threshold, 1cm step)
    const snappedX = snapToCanvasGrid(newX)
    const snappedY = snapToCanvasGrid(newY)
    
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
  
  // 3. Drafting wall / door / enclosure preview
  if (wallDraftStart.value && ['walls', 'wall-rect', 'door'].includes(cadToolMode.value)) {
    const pt = svgElement.createSVGPoint()
    pt.x = event.clientX
    pt.y = event.clientY
    const svgP = pt.matrixTransform(svgElement.getScreenCTM().inverse())
    
    // Snap coordinates with magnetic major-grid priority (10cm threshold, 1cm step)
    const snappedX = snapToDatabaseGrid(svgP.x / 2)
    const snappedY = snapToDatabaseGrid(svgP.y / 2)
    
    wallDraftCurrent.value = { x: snappedX, y: snappedY }
  }
}

const onDesignerMouseUp = () => {
  // Release cabinet drag
  if (draggedRackId.value) {
    draggedRackId.value = null
    return
  }

  // Release wall dragging element (Phase 4 CAD Upgrades)
  if (draggedWallId.value) {
    draggedWallId.value = null
    return
  }

  // Release wall resizing handle
  if (resizingWallId.value) {
    resizingWallId.value = null
    resizingHandleType.value = null
    return
  }
  
  // Release wall/door drawing (Buffer to memory)
  if (wallDraftStart.value && wallDraftCurrent.value && activeFloor.value) {
    const x1 = wallDraftStart.value.x
    const y1 = wallDraftStart.value.y
    const x2 = wallDraftCurrent.value.x
    const y2 = wallDraftCurrent.value.y
    const floorId = activeFloor.value.id
    
    wallDraftStart.value = null
    wallDraftCurrent.value = null
    
    if (x1 === x2 && y1 === y2) return
    
    if (cadToolMode.value === 'walls') {
      // Standard single line wall
      const newWall = {
        id: 'temp_' + Date.now() + '_' + Math.floor(Math.random() * 1000),
        datacenter_floor_id: floorId,
        x1, y1, x2, y2,
        thickness: 20, // 20cm thickness
        type: 'wall'
      }
      localWalls.value.push(newWall)
      selectedElementId.value = newWall.id // Auto-select on creation
      cadToolMode.value = 'select' // Switch back to Select & Edit mode
      markDirty()
    } else if (cadToolMode.value === 'door') {
      // Door segment (thinner red line with door type)
      const newDoor = {
        id: 'temp_door_' + Date.now() + '_' + Math.floor(Math.random() * 1000),
        datacenter_floor_id: floorId,
        x1, y1, x2, y2,
        thickness: 10, // 10cm thickness
        type: 'door'
      }
      localWalls.value.push(newDoor)
      selectedElementId.value = newDoor.id // Auto-select on creation
      cadToolMode.value = 'select' // Switch back to Select & Edit mode
      markDirty()
    } else if (cadToolMode.value === 'wall-rect') {
      // Wall enclosure (creates a single rectangular room element!)
      const newRect = {
        id: 'temp_rect_' + Date.now() + '_' + Math.floor(Math.random() * 1000),
        datacenter_floor_id: floorId,
        x1, y1, x2, y2,
        thickness: 20, // 20cm thickness
        type: 'wall-rect'
      }
      localWalls.value.push(newRect)
      selectedElementId.value = newRect.id // Auto-select on creation
      cadToolMode.value = 'select' // Switch back to Select & Edit mode
      markDirty()
    }
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
            y2: wall.y2,
            thickness: wall.thickness || 8,
            type: wall.type || 'wall'
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
  if (!activeFloor.value || !selectedFloorId.value) return
  try {
    const payload = {
      datacenter_id: dcId,
      datacenter_floor_id: selectedFloorId.value, // Link newly created cabinet directly to active floor level
      name: newCabinetForm.value.name,
      height_u: Number(newCabinetForm.value.height_u),
      x: 0, // Starts as Not Placed
      y: 0  // Starts as Not Placed
    }
    await $fetch(`${apiBase}/racks/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
    newCabinetForm.value = { name: '', height_u: 42 }
    isDeployModalOpen.value = false
    await refreshDatacenters()
    syncLocalStateFromDB()
    console.log('Cabinet deployed successfully in drawing tool')
  } catch (err) {
    console.error('Failed to deploy cabinet from workspace:', err)
  }
}

const placeCabinetOnGrid = (rack) => {
  rack.x = 40
  rack.y = 40
  markDirty()
}

const unplaceCabinetFromGrid = (rack) => {
  rack.x = 0
  rack.y = 0
  markDirty()
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

onMounted(() => {
  window.addEventListener('click', closeContextMenu)
})

// Clear auto-save timer when exiting drawing page
onBeforeUnmount(() => {
  window.removeEventListener('click', closeContextMenu)
  if (saveTimer.value) {
    clearTimeout(saveTimer.value)
  }
})
</script>