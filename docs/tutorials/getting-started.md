# Getting Started with Northstar

This tutorial is a hands-on, learning-oriented introduction to Northstar. You will learn how to bootstrap a clean local workspace, spin up the Go backend, run the Nuxt 3 dashboard, and register your first physical server asset.

---

## 🎯 Learning Objectives

By the end of this tutorial, you will have:
1. Checked out and set up the correct runtime environments (Go + Node).
2. Bootstrapped a fresh local database seeded with Astrona taxonomy models.
3. Authenticated your user operator session and browsed the 2D CAD canvas.
4. Created and verified a physical server asset linked to a network interface.

---

## 🛠️ Prerequisites

To complete this tutorial, you need to install the following on your machine:
- **Golang 1.20+**
- **Node.js 24+ & npm**
- **Git**

---

## Step 1: Clone the Codebase and Run Setup

First, clone your repository and navigate to the project root directory.

Run the package setups for both layers of the stack:

=== "Go Backend Setup"
    ```bash
    cd backend
    go mod download
    cd ..
    ```

=== "Nuxt Frontend Setup"
    ```bash
    cd frontend
    npm install
    cd ..
    ```

---

## Step 2: Start the Go REST Server

The Go backend handles automatic GORM schema migrations, processes SBOM SPDX uploads, and performs active subnet discovery scans.

Start the backend server:
```bash
cd backend
go run cmd/server/main.go
```

Verify that you see the following initialization logs in your terminal:
```text
[Migration Engine] Applying automatic GORM schema AutoMigrate updates...
[Bootstrap Engine] All modular K8s-style bootstrapping manifests reconciled successfully.
Starting Go server on port 8000...
```

---

## Step 3: Run the Nuxt 3 Dashboard

In a separate terminal window, start the frontend development server:
```bash
cd frontend
npm run dev
```

Once compiled, navigate your web browser to **`http://localhost:3000`**. You will be automatically redirected to the secure Sign-In page.

---

## Step 4: Sign In and Register Your First Asset

1. **Sign In:** Use the pre-seeded operator credentials:
   - **Username:** `admin`
   - **Password:** `admin`
2. **Access Wizard:** Click the glowing green **"Add Asset"** button in the top right corner of the dashboard.
3. **Choose Category:** Select **Server** on Step 1.
4. **Enter Attributes:**
   - **Name:** `Northstar-Edge-01`
   - **IP Address:** `10.0.1.15`
   - **CPU Model:** `Intel Xeon D-1541` (Fill under Custom Properties tab)
5. **Save and Verify:** Click **Save**. The dashboard list will instantly populate with your new asset, displaying an active monitoring card and scheduling hooks.

🎉 **Congratulations!** You have successfully set up, run, and registered your first physical hardware asset on the Northstar platform!
