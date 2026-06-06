import { describe, it, expect, afterEach } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import { nextTick } from 'vue'
import AssetFormModal from './AssetFormModal.vue'

describe('AssetFormModal', () => {
  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('displays Server properties when type is Server', async () => {
    await mountSuspended(AssetFormModal, {
      props: {
        modelValue: true,
        asset: null
      }
    })

    // Click on Server category card to move to Step 2
    const buttons = document.querySelectorAll('button')
    const serverButton = Array.from(buttons).find(b => b.textContent?.includes('Server'))
    if (serverButton) {
      serverButton.click()
    }
    await nextTick()

    expect(document.body.innerHTML).toContain('Server Properties')
    expect(document.body.innerHTML).toContain('CPU Model')
    expect(document.body.innerHTML).not.toContain('Database Properties')
  })

  it('displays Database properties when type is changed to Database', async () => {
    await mountSuspended(AssetFormModal, {
      props: {
        modelValue: true,
        asset: {
          name: '',
          type: 'Database',
          status: 'active',
          properties: {}
        }
      }
    })

    expect(document.body.innerHTML).toContain('Database Properties')
    expect(document.body.innerHTML).toContain('Engine')
    expect(document.body.innerHTML).not.toContain('Server Properties')
  })

  it('autofills specifications when a device model is selected', async () => {
    const wrapper = await mountSuspended(AssetFormModal, {
      props: {
        modelValue: true,
        asset: {
          name: 'Core-Switch',
          type: 'Network',
          status: 'active',
          properties: {
            network_subtype: 'Switch (L3)'
          }
        }
      }
    })

    const mockDeviceModel = {
      id: 'dm-99',
      model_name: 'Catalyst 9300',
      manufacturer: { name: 'Cisco' },
      revision: 2,
      ports: { RJ45: 24, 'SFP+': 4 }
    }

    wrapper.vm.devices = [
      {
        ...mockDeviceModel,
        displayName: 'Cisco - Catalyst 9300 (v2)'
      }
    ]
    await nextTick()

    // Simulate selection
    wrapper.vm.selectedDeviceModel = wrapper.vm.devices[0]
    await nextTick()

    expect(wrapper.vm.form.device_model_id).toBe('dm-99')
    expect(wrapper.vm.form.device_model_revision).toBe(2)
    expect(wrapper.vm.form.properties.manufacturer).toBe('Cisco')
    expect(wrapper.vm.form.properties.model).toBe('Catalyst 9300')
    expect(wrapper.vm.form.properties.port_config).toBe('24x RJ45 + 4x SFP+')
  })
})
