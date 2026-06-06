import { describe, it, expect, vi } from 'vitest'
import { mountSuspended, mockNuxtImport } from '@nuxt/test-utils/runtime'
import { ref } from 'vue'
import TabPage from './[[tab]].vue'

// Mock auto-imported useRoute so it returns the dcim tab
mockNuxtImport('useRoute', () => {
  return () => ({
    params: { tab: 'dcim' }
  })
})

// Mock auto-imported useAuth to grant full administrative privileges
mockNuxtImport('useAuth', () => {
  return () => ({
    isAuthenticated: ref(true),
    isAdmin: ref(true),
    canMutate: ref(true),
    canDelete: ref(true),
    getAuthHeader: () => ({})
  })
})

describe('Settings Tab Page - Floor Management UI', () => {
  it('renders Floor Level Selector, Rename Alias and Add Floor controls', async () => {
    // Mount page suspended
    const wrapper = await mountSuspended(TabPage)

    // Force values
    wrapper.vm.activeTab = 'dcim'
    wrapper.vm.datacenters = [
      {
        id: 'dc-1',
        name: 'Northstar Dublin HQ',
        city: 'Dublin',
        country: 'Ireland',
        floors: [
          { id: 'floor-0', name: 'Floor 0', level: 0 }
        ]
      }
    ]
    wrapper.vm.selectedDcIdForDesigner = 'dc-1'

    await wrapper.vm.$nextTick()

    const html = wrapper.html()
    expect(html).toContain('Active Floor Level')
    expect(html).toContain('Rename (Alias)')
    expect(html).toContain('Add Floor')
  })

  it('opens edit modal with datacenter details when edit button is clicked', async () => {
    const wrapper = await mountSuspended(TabPage)

    wrapper.vm.activeTab = 'dcim'
    const dcObj = {
      id: 'dc-1',
      name: 'Northstar Dublin HQ',
      type: 'on-prem',
      city: 'Dublin',
      country: 'Ireland',
      properties: {
        uplink_speed: '10 Gbps',
        public_ip: '1.2.3.4'
      },
      floors: []
    }
    wrapper.vm.datacenters = [dcObj]
    await wrapper.vm.$nextTick()

    wrapper.vm.openEditDcModal(dcObj)
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.isEditDcModalOpen).toBe(true)
    expect(wrapper.vm.editDcForm.name).toBe('Northstar Dublin HQ')
    expect(wrapper.vm.editDcForm.uplink_speed).toBe('10 Gbps')
  })
})
