import { test, expect } from '@playwright/test';
import crypto from 'crypto';

// Generate a mathematically valid JWT token signed with the backend's key
// to ensure perfect backend compatibility for RBAC middleware.
function generateRealToken(username = 'admin', role = 'admin') {
  const header = { alg: 'HS256', typ: 'JWT' };
  const payload = {
    username,
    role,
    exp: Math.floor(Date.now() / 1000) + 24 * 60 * 60,
  };
  
  const base64UrlEncode = (obj: any) => {
    return Buffer.from(JSON.stringify(obj))
      .toString('base64')
      .replace(/=/g, '')
      .replace(/\+/g, '-')
      .replace(/\//g, '_');
  };

  const encodedHeader = base64UrlEncode(header);
  const encodedPayload = base64UrlEncode(payload);

  const secret = 'cmdb-secret-key-2026';
  const signature = crypto
    .createHmac('sha256', secret)
    .update(`${encodedHeader}.${encodedPayload}`)
    .digest('base64')
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');

  return `${encodedHeader}.${encodedPayload}.${signature}`;
}

test.describe('UX Testing: CMDB Application', () => {
  test.beforeEach(async ({ context }) => {
    const token = generateRealToken('operator', 'operator');
    
    // Seed authenticated cookie state as operator so pages are accessible immediately (Phase 2 Auth E2E Stability)
    await context.addCookies([
      {
        name: 'cmdb-session',
        value: JSON.stringify({
          token,
          username: 'operator',
          role: 'operator'
        }),
        domain: 'localhost',
        path: '/'
      }
    ]);
  });

  test('Homepage loads correctly with Asset Overview', async ({ page }) => {
    // Navigate to frontend URL
    await page.goto('/');
    
    // Wait for Nuxt to hydrate
    await page.waitForTimeout(2000);

    // Check main title
    await expect(page.locator('h2')).toContainText('Assets Overview');
    
    // Check Add Asset button exists (UX check)
    const addAssetBtn = page.getByRole('button', { name: 'Add Asset' });
    await expect(addAssetBtn).toBeVisible();

    // Click Add Asset opens modal
    await addAssetBtn.click();
    await expect(page.locator('h3').filter({ hasText: 'Select Asset Type' })).toBeVisible();

    // Close modal
    await page.getByRole('button', { name: 'Cancel' }).click();
    await expect(page.locator('h3').filter({ hasText: 'Select Asset Type' })).not.toBeVisible();
  });

  test('Successful Sign-in as admin triggers license popup and accepts after scroll', async ({ context, page }) => {
    // Clear pre-seeded cookies to simulate starting unauthenticated
    await context.clearCookies();

    // Navigate directly to login
    await page.goto('/login');
    await page.waitForTimeout(1000);

    // Enter admin credentials
    await page.getByPlaceholder('admin, operator, or viewer').fill('admin');
    await page.getByPlaceholder('Enter password...').fill('admin');

    // Click submit
    await page.getByRole('button', { name: 'Sign In' }).click();
    await page.waitForTimeout(2000);

    // Verify redirection to home and that License Modal pops up immediately
    await expect(page).toHaveURL('/');
    const modalTitle = page.locator('h3').filter({ hasText: 'GNU Affero General Public License v3' });
    await expect(modalTitle).toBeVisible();

    // Verify "I Agree" button is disabled by default
    const agreeBtn = page.getByRole('button', { name: 'I Agree & Accept AGPLv3 Terms' });
    await expect(agreeBtn).toBeDisabled();

    // Locate the scrollable license text container and scroll it down to 100%
    const scrollContainer = page.locator('div.p-6.overflow-y-auto');
    await scrollContainer.evaluate((el) => {
      el.scrollTop = el.scrollHeight;
      el.dispatchEvent(new Event('scroll'));
    });
    await page.waitForTimeout(1000);

    // Verify "I Agree" button is now enabled
    await expect(agreeBtn).toBeEnabled();

    // Click on I Agree to accept
    await agreeBtn.click();
    await page.waitForTimeout(2000);

    // Verify modal is closed and they can see the signed admin profile
    await expect(modalTitle).not.toBeVisible();
    await expect(page.locator('span').filter({ hasText: 'Signed Operator' })).toBeVisible();
    await expect(page.locator('span').filter({ hasText: 'admin' }).first()).toBeVisible();
  });

  test('Settings page /settings/dcim loads correctly and displays room layout designer', async ({ context, page }) => {
    // Authenticate explicitly as admin for this admin-only settings view
    const token = generateRealToken('admin', 'admin');
    await context.addCookies([
      {
        name: 'cmdb-session',
        value: JSON.stringify({
          token,
          username: 'admin',
          role: 'admin'
        }),
        domain: 'localhost',
        path: '/'
      }
    ]);

    // Go to settings/dcim
    await page.goto('/settings/dcim');
    await page.waitForTimeout(2000);

    // If the license modal is visible, scroll and accept it first
    const modalTitle = page.locator('h3').filter({ hasText: 'GNU Affero General Public License v3' });
    if (await modalTitle.isVisible()) {
      const scrollContainer = page.locator('div.p-6.overflow-y-auto');
      await scrollContainer.evaluate((el) => {
        el.scrollTop = el.scrollHeight;
        el.dispatchEvent(new Event('scroll'));
      });
      await page.waitForTimeout(1000);
      await page.getByRole('button', { name: 'I Agree & Accept AGPLv3 Terms' }).click();
      await page.waitForTimeout(1000);
    }

    // Click on Deploy New Location Profile button to trigger the Modal
    const openModalBtn = page.getByRole('button', { name: 'Deploy New Location Profile' });
    await expect(openModalBtn).toBeVisible();
    await openModalBtn.click();
    await page.waitForTimeout(1000);

    // Verify Modal title is visible
    await expect(page.locator('h4').filter({ hasText: 'Deploy New Location Profile' })).toBeVisible();

    // Close the Modal
    await page.getByRole('button', { name: 'Cancel' }).click();
    await page.waitForTimeout(1000);
    
    // Click on Northstar-Dublin-HQ card in the location list
    await page.locator('div.cursor-pointer', { hasText: 'Northstar-Dublin-HQ' }).first().click();
    await page.waitForTimeout(1000);

    // Verify Room Layout Preview label becomes visible
    await expect(page.locator('span').filter({ hasText: 'Room Layout Preview' })).toBeVisible();

    // Click Enter Drafting Mode button inside the folded out datacenter card
    await page.getByText('Enter Drafting Mode').first().click();
    await page.waitForTimeout(2500);

    // Verify redirection to settings/dcim/[id]/drawing and workspace heading is visible
    await expect(page).toHaveURL(/\/settings\/dcim\/[a-zA-Z0-9-]+\/drawing/);
    await expect(page.locator('h1').filter({ hasText: 'Northstar CAD Blueprint Workspace' })).toBeVisible();
  });
});
