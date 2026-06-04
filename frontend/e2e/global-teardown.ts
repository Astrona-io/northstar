import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default async function globalTeardown() {
  const dbPaths = [
    path.resolve(__dirname, '../../backend/cmdb_e2e.db'),
    path.resolve(__dirname, '../../backend/cmdb_e2e.db-journal'),
    path.resolve(__dirname, '../../backend/cmdb_e2e.db-shm'),
    path.resolve(__dirname, '../../backend/cmdb_e2e.db-wal'),
  ];

  for (const dbPath of dbPaths) {
    if (fs.existsSync(dbPath)) {
      try {
        fs.unlinkSync(dbPath);
        console.log(`[Global Teardown] Cleaned up E2E database file: ${dbPath}`);
      } catch (err) {
        console.error(`[Global Teardown] Error deleting ${dbPath}:`, err);
      }
    }
  }
}
