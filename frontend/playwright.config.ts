import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  retries: 0,
  workers: 1,
  reporter: 'list',
  globalTeardown: './e2e/global-teardown.ts',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: 'cd ../backend && rm -f cmdb_e2e.db cmdb_e2e.db-journal cmdb_e2e.db-shm cmdb_e2e.db-wal && DATABASE_URL=cmdb_e2e.db go run cmd/server/main.go',
      url: 'http://localhost:8000/api/categories/',
      reuseExistingServer: false,
      timeout: 120 * 1000,
    },
    {
      command: 'npm run dev',
      url: 'http://localhost:3000',
      reuseExistingServer: false,
      timeout: 120 * 1000,
    }
  ],
});
