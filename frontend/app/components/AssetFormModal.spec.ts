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
})
