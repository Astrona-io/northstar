# Northstar CMDB Refactoring & Enhancement Roadmap

This document outlines the architectural milestones, completed enhancements, and advanced features delivered to establish a production-grade Configuration Management Database (CMDB) application.

---

## 🗺️ Future Vision: Active Enhancement Phases

With our Gitea/Forgejo slate-teal layout, Site Admin dashboard, dynamic custom fields, webhook publishers, high-fidelity 42U physical server cabinets, 2D CAD Floor Planners, location filtering dropdowns, 2D Room Blueprint Wall Designers, and Silent Startup sequencers completely resolved, the next generation of Northstar enhancements will target drag-and-drop links, facility zoning, and icon pickers.

### Phase 1: Dynamic Tree-Graph Drag-and-Drop Taxonomy Link Assignments (UI L4)
* **Goal:** Elevate the dynamic taxonomy map, allowing administrators to visually connect prerequisite chains via drag-and-drop mouse events.
* **Details:**
  - **Drag-and-Drop Node Connectors:** Implement HTML5 Drag-and-Drop handlers inside the Taxonomy Tree-Graph panel in Site Admin.
  - **Instant GORM Persistences:** Dragging a child Category card (e.g. "Database") and dropping it directly onto a parent category card (e.g. "Server") automatically sets the child's `parent_dependency` string to the parent's name, triggering an asynchronous GORM persistence update in GORM.

### Phase 2: Custom Facility Placement Zones and Aisle Hot/Cold Corridor Indicators (DCIM L4)
* **Goal:** Optimize server rack thermal layouts by allowing administrators to sketch colored hot/cold corridors and custom placement zones onto the CAD floor.
* **Details:**
  - **Geometric Zone Drawing Tool:** Add a "Draw Zones" mode to the settings CAD layout designer, letting admins drag and draw rectangular regions on the 1-meter dotted grid.
  - **Thermal Corridor Indicators:** Admins can flag drawn zones as "Hot Corridor" (colored in light red transparency `fill="rgba(239, 68, 68, 0.1)"` and bordered in dashed red) or "Cold Corridor" (colored in light blue transparency `fill="rgba(56, 189, 248, 0.1)"`), visually grouping cabinets and optimizing ventilation efficiency.

### Phase 3: Live Network Cable Connection Ping Audits & Status Flashes (Cabling L4)
* **Goal:** Enable live verification of port-to-port cable lines using background interface socket sweep probes.
* **Details:**
  - **Interface Socket Probes:** A background worker executes active local loopback dials on patched target ports in real-time.
  - **Dashed Flashing Alert Links:** If the target socket is unresponsive or timed out, automatically flash that connection link in blinking dashed yellow on the SVG topology map, instantly warning operators of physical/virtual connection line disconnections!

### Phase 4: Custom Category Icon Search & Picker in Site Admin (UI L4)
* **Goal:** Elevate administrative experience by integrating a search-based **Icon Picker Modal** inside the Category Manager.
* **Details:**
  - **Tailwind/Heroicons Bundle Search:** Deploys a live icon-search bar. When typing (e.g., "power"), query and render visual matches (`i-heroicons-bolt`, `i-heroicons-battery-100`, `i-heroicons-cpu-chip`) in a select grid dynamically.
  - **Visual Click-to-Assign:** Admins simply click on their preferred icon to instantly bind it to their new asset categories, removing icon class text lookups completely.

### Phase 5: Dynamic Telemetry Alerts based on Sub-Stack Selection (Observability L4)
* **Goal:** Enable specialized telemetry drift auditing by polling metrics tailored specifically to the asset's chosen sub-stack (such as loop voltages or camera resolution states).
* **Details:**
  - **Dynamic Ingress Polling:** A background checker scans active properties JSON columns, reading subtype metrics (e.g. checking loop voltages on `Industrial PLC / SCADA` or resolution tags on `Smart IP Camera`).
  - **Alert Threshold Indicators:** Render real-time warnings on the asset details gauge if telemetry drifts beyond normal bounds (e.g. throwing a red warning if loop voltage is outside `24V +/- 1.5V` tolerances, or if a camera streams below H.264 high-definition parameters).

### Phase 6: OS-Level Package Auditing for IoT Edge Devices (SecOps L4)
* **Goal:** Expand Software Security SBOM scans to physical Edge Gateways by executing SSH-agent package sweeps on Raspberry Pi and Jetson Nano nodes.
* **Details:**
  - **Debian/Alpine Package Scanners:** integrate secure SSH crawler hooks to execute non-interactive commands (such as `dpkg-query -l` on Raspbian or `apk info -v` on Alpine Edge devices) to automatically parse installed library manifests.
  - **Physical SBOM Vulnerability Maps:** Map extracted edge packages directly onto the asset details SBOM tab, querying CVE databases in real-time to alert operators of security vulnerabilities inside physical remote field equipment.

### Phase 7: Real-Time Item-Level Declarative YAML Exporters (UI DX)
* **Goal:** Make GitOps completely tangible and user-friendly by adding a **"View GitOps YAML"** code portal directly onto individual asset details, cabinet, and datacenter panels.
* **Details:**
  - **Single-Item Code Compilation:**
    - On the Asset Detail page (`/assets/[id]`) or Datacenter view (`/dcim`), add a clean **"View YAML"** button on the primary header or property card.
    - Clicking the button opens a modal compiling that exact, specific database record (e.g. including its live UUID, GORM relationships, custom properties, and mapped interfaces) on the fly into its valid, declarative YAML manifest representation.
  - **Looping UI to GitOps DX:**
    - Provide a quick **"Copy to Clipboard"** button.
    - Operators can configure complex network connections and custom attributes using Northstar's visual UI, click "View YAML", copy the compiled manifest, check it into their declarative Git repositories, and let GitOps take over future updates.

### Phase 8: Zero-Trust Host Firewall Audits & OS-Level Rules Verification (SecOps L4)
* **Goal:** Evolve the CMDB topology from simple cabling maps to an active **Zero-Trust Network Posture map** by auditing and verifying host-level firewall configurations (such as Linux `firewalld`, `ufw`, or `iptables` rules).
* **Details:**
  - **OS-Level Firewall Ingestion (Optional Feature):**
    - **Auto-Loaded (Discovery Sockets):** Integrate secure SSH crawler hooks to execute non-interactive commands (e.g., `firewall-cmd --list-all` on Fedora/RedHat or `ufw status` on Ubuntu) to automatically pull and parse active host-level ingress rules during auto-discovery scans.
    - **Manually Declared (GitOps & UI):** Provide dedicated, tabbed firewall configurations inside the Asset details panel and GitOps YAML specs, allowing operators to manually declare allowed port matrices.
  - **Host Ingress Rule Matrix (Security Audit):**
    - Store parsed ingress filters as a standardized JSON structure (e.g. `[{port: 80, source: "10.0.0.0/24", policy: "allow"}]`) inside GORM.
    - On the visual SVG topology map, clicking a connection line between two nodes dynamically audits the target's firewall matrix to verify if traffic is actually allowed or blocked by the host OS (rendering green flowing paths for verified allowed links, and red dashed warning paths if the host firewall blocks the source subnet!).

### Phase 9: OpenTelemetry Ingestion & Distributed Observability (OTel L3)
* **Goal:** Integrate OpenTelemetry (OTel) Tracing, Metrics, and Logs exporters in the Go backend to enable production-grade distributed tracing and observability across APM platforms (such as Jaeger, Prometheus, or Datadog).
* **Details:**
  - **OTel Go SDK Integration:** Set up the OpenTelemetry TracerProvider and MeterProvider, establishing outbound OTLP exporters over gRPC or HTTP/JSON.
  - **HTTP Route Tracing & Spans:** Instrument the Echo HTTP router to automatically initiate spans for incoming requests, injecting trace contexts across outbound webhooks and background auto-discovery scan routines.
  - **GORM DB Tracing:** Configure GORM callbacks to trace SQLite transaction runtimes, reporting slow database queries and transaction locks as database spans.
  - **Standard Metrics Exporters:** Emit native runtime metrics (Go memory pools, Active HTTP Request counts, Webhook broadcaster queue depths, and ping latency metrics) as standardized OpenTelemetry metrics.

### Phase 10: Diátaxis-Compliant MkDocs Architecture & Developer Documentation (Ironroot Branding Docs L3)
* **Goal:** Bootstrap a high-fidelity, Diátaxis-conforming **MkDocs** project providing comprehensive architecture diagrams and reference guidelines, beautifully styled to align with the **Ironroot corporate design guidelines**.
* **Details:**
  - **Ironroot Styling & Material Theme customization:**
    - Deploy the MkDocs Material theme (`theme: material`) customized with **Ironroot's core branding colors**:
      - Primary Color: **Deep Steel/Slate Charcoal** (representing the clean, industrial "iron" metallic look).
      - Accent Color: **Teal/Emerald** (Forgejo/Netic brand theme colors) or customized Ironroot orange/green highlights.
    - Configure custom Ironroot font typography (sharp monospace code blocks and headers, clean sans-serif body text) and corporate logo headers.
    - Deploys full native dark-mode selectors with automated system theme handshakes matching Ironroot standards.
  - **Diátaxis Hierarchy:** Structure documentation across the four standard Diátaxis quadrants: *Tutorials*, *How-To Guides*, *Reference*, and *Explanation / Architecture*.
  - **System Architecture Drawings:** Embed clean, text-based Mermaid.js diagrams detailing Northstar's full-stack Echo-Nuxt REST handshakes, JWT-RBAC permission weighting model, and declarative GitOps reconciliation loops.

### Phase 11: Human-Logical Folder Reorganization & Clean Code Refactoring
* **Goal:** Refactor backend Go modules, database structures, and frontend Vue files to maximize human readability and maintain clean code principles.
* **Details:**
  - **Logical Folder Reorganization:** Streamlines directory structures so that resources are placed in highly descriptive, single-responsibility folder layers, separating business controllers, db migrations, and routers clearly.
  - **Human-Readable Code Refactoring:** Audits complex logical loops, removing inline magic variables and replacing them with self-explanatory constants, clean delegation patterns, and comprehensive inline Go/Vue docstrings.

---

## 🏛️ Completed Milestones Archive

All foundational, advanced, visual network topology, DCIM, SecOps, GitOps, and onboarding phases of the Northstar CMDB refactoring pipeline have been successfully developed, integrated, verified, and committed:

- [x] **Phase 1: Go Backend Rewrite (Echo/GORM):** Migration from legacy Python to robust, highly performant Go structure with automatic database migrations and GORM schemas.
- [x] **Phase 2: Frontend Componentization:** Extracted monolithic Vue pages into reusable, single-responsibility components (`AssetTable`, `AssetFormModal`, `MaintenanceTable`, `MaintenanceFormModal`, `PropertiesDisplay`, `AuditLogsTimeline`, `InterfacesTable`, `InterfaceFormModal`, `DeviceFormModal`, `AuditLogsTimeline`).
- [x] **Phase 3: Decoupled API Composables:** Introduced type-safe Vue/Nuxt composables encapsulating all assets, interfaces, relationships, audit logs, devices, and discovery queries.
- [x] **Phase 4: Nuxt 4 Native Migration:** Fully upgraded frontend build targets to **Nuxt v4.4.7** and relocated layouts, pages, components, and composables into the official native `app/` folder structure with custom `compatibilityDate` configurations.
- [x] **Phase 5: Dynamic Wizard Form & Icons:** Overhauled asset registration into a 2-step setup wizard using searchable cards, custom Category visual icons (Heroicons), and dynamic input configurations.
- [x] **Phase 6: Advanced Querying, Search & Filtering:** Configured GORM parameters to support case-insensitive database searches on asset descriptions/names alongside custom Type and Status filters.
- [x] **Phase 7: Network Interface Cards (NIC) Port Management:** Restructured database mapping to support one-to-many physical/virtual ports detailing individual Link Speeds, MACs, standard IPs, VLANs, and MTUs.
- [x] **Phase 8: Topology Relationships (Uplink/Downlink):** Added connections mapping to link routers/switches directly inside the creation wizard, rendering structured upstream/downstream cards in the Asset Detail tab.
- [x] **Phase 9: Centralized Device Specifications Catalog:** Created a master list of Device Profile Specifications (Manufacturer, Model Name, Category, and Specs) with dedicated frontend page panels and back-end Echo handlers to standardize hardware deployments.
- [x] **Phase 10: JWT Authentication & Fine-Grained RBAC:** Secured all write/delete endpoints behind JSON Web Tokens, validating mock profiles (`admin`/`operator`/`viewer`/`hybrid`) with strict **User override Group policy-weighting** resolutions.
- [x] **Phase 11: Real-Time Network Auto-Discovery (Port Scanner):** Engineered an active loopback TCP port scanner on the backend that dynamically sweeps socket dials on 127.0.0.1, identifying live Go/Node processes in real-time.
- [x] **Phase 12: Visual Topology Graph Map:** Engineered a custom, responsive, lightweight **interactive SVG cabling graph** that maps Routers and Switches connections in real-time, complete with a side details inspector pane and quick navigation portals.
- [x] **Phase 13: Datacenter & Rack Unit Management (DCIM Mapping):** Developed physical datacenter models (categorized by On-Prem, Homelab, Air-Gap, and AWS/GCP/Azure clouds) and 42U rack cabinet schemas, creating an interactive, visual vertical server cabinet front layout in Nuxt 4 detailing exact physical U-height slot assignments.
- [x] **Phase 14: First-Class Kubernetes Deployment Asset Type:** Embedded virtual Kubernetes workload deployments directly into the core schema, supporting container image tags and virtual service port mappings linked to GORM hosting parent Node Server assets.
- [x] **Phase 15: Declarative Database Migration Engine:** Integrated a version-ledgers schema engine (`migrations.go`) that parses, executes, and audits file-based database migration scripts programmatically, matching the capabilities of Python's **Alembic**.
- [x] **Phase 16: Outbound Webhook Integration Broadcaster:** Programmed an asynchronous, thread-safe background event broadcaster (`webhooks.go`) that dispatches signed, real-time HTTP POST callback notifications (Slack/Jira alerts) upon asset creations, edits, and deletions.
- [x] **Phase 17: User-Defined Dynamic Asset Schemas (Custom Fields):** Built an administrative Custom Fields portal allowing administrators to define custom metadata property constraints (Name, ValueType, Target Category) text inputs dynamically inside GORM JSONMaps.
- [x] **Phase 18: High-Fidelity "Full-Screen" Topology Viewpoint:** Added a maximize toggle to the SVG Topology patch map, allowing network engineers to expand the SVG canvas into a fixed, immersive full-screen viewpoint with layout overlays sliding out of view automatically.
- [x] **Phase 19: 4-Column Wide Asset Creation Wizard:** Configured `AssetFormModal.vue` with an `sm:max-w-5xl` modal width and a highly responsive `grid-cols-2 md:grid-cols-4` grid, rendering all category cards neatly inside exactly two compact horizontal rows with zero vertical scrolling.
- [x] **Phase 20: Dual Cloud-Native Workload Modeling:** Restored standalone `Container` tasks (Docker, ECS, AWS Fargate tasks) side-by-side with clustered `Kubernetes Deployments` to accurately track cloud-native software and dependencies.
- [x] **Phase 21: Consolidated Network Asset Classification:** Combined separate switches and routers into a single "Network" card, supporting Router, Switch, AP, Firewall, and Load Balancer subtypes saved dynamically inside GORM JSONMaps and rendered with custom SVG circles and icons on the cabling topology map.
- [x] **Phase 22: GitLab-Style "Site Admin" Re-Organization:** Grouped and secured administrative utilities (including dynamic Subnet Sweepers, AWS Cloud Sync Crawlers, and DCIM Location administrators) under a unified, professional, and highly innovative GitLab-style Site Admin vertical nav card hub, completely eliminating global navigation clutter.
- [x] **Phase 23: Interactive Clickable "Guided Onboarding):** Converted dashboard onboarding checklist rows into clickable navigators with Datacenter primacy, directing operators directly to prerequisite panels.
- [x] **Phase 24: High-Fidelity "42U Server Cabinet" DCIM Visual Overhaul:** Overhauled `dcim.vue` to render a highly realistic 42U physical metal cabinet featuring ventilation grilles, mounting rails, blinking LED status indicators, and styled geodata WAN speed site badges.
- [x] **Phase 25: Real-Time Network Latency Monitoring & Status Heatmap:** Built background socket ping latency dial metrics in Go, exposing a live status heatmap grid displaying online ms latency nodes and offline loss percentages in real-time.
- [x] **Phase 26: Inter-Rack Cable Patch Board Matrix (DCIM L2):** Modeled inter-rack physical fiber and solid copper cabling connections running between different server cabinet frames, presenting an interactive DCIM patch ledger mapping cross-connected ports.
- [x] **Phase 28: Dynamic Cabling Trace Walkthrough:** Programmed an interactive physical port cable tracer. Selecting a hardware node inside the visual cabinet triggers an animated, pulsing path highlighting the corresponding cable loop directly on the SVG topology map.
- [x] **Phase 29: Live Kubernetes Workload Terminal Console:** Integrated a live, retro-dark terminal console on Container/Kubernetes Details page, displaying simulated live scrolling startup logs and active pod event state simulators (running, pending, terminating).
- [x] **Phase 30: Site Admin Access Event Audit Ledger:** Created a GitLab-style administrative Access Audit Event log grid preloaded under Site Admin, exposing secure configuration, custom fields, and DCIM location edits alongside dynamic CSV download triggers.
- [x] **Phase 31: Interactive Datacenter Floor Planner & Precision CAD Canvas (Draw.io Style):** Re-engineered the Datacenters dashboard (`dcim.vue`) to deploy a grid-snapped, responsive SVG CAD floorplanner enabling precision room mapping and click-to-open cabinet visual rail inspectors, completely consolidated under a top horizontal filter navigation.
- [x] **Phase 32: Multi-Level / Multi-Floor DCIM CAD Layout & Level Selectors:** Expanded Datacenter models with GORM-backed `DatacenterFloor` relations, allowing each individual floor level to maintain and switch its own distinct 2D floorplan CAD layouts dynamically using horizontal tab selectors.
- [x] **Phase 33: Dynamic & User-Defined Datacenter Types:** Transitioned facility classifications from hardcoded dropdowns to GORM-seeded `DatacenterType` database models, exposing full admin panels to configure custom categories dynamically.
- [x] **Phase 34: Cabinet Placement Zones & Room Layout Aisle Groupings:** Mapped cabinet coordinates and physical rows/aisles inside GORM schemas, rendering specific, styled Placement Zone tags on the visual 2D floorplan grid below cabinet footprint selectors.
- [x] **Phase 35: Global Quickstart Checklist Layout Integration:** Mounts the progress-tracking onboarding checklist card directly inside the global layout wrapper `default.vue`, enabling constant, shared-state guided self-learning pathways for newcomers across any explored view.
- [x] **Phase 36: Production Hardening & Real-Time Engine Integrations (Zero Mocks):** Hardened platform infrastructures with a real backend SPDX file parser `/api/assets/:id/sbom`, real global GORM-backed audit trails, and multi-threaded background workers.
- [x] **Phase 37: Multi-Datacenter WAN Peering & SVG Topology Tunnel Maps:** Added a global WAN map toggle to the network Topology page, rendering geographical datacenter locations and live encrypted IPsec/WireGuard VPN tunnel links with dynamic bandwidth profiles.
- [x] **Phase 38: Multi-Layer Testing Matrix:**
  - Go handler integration and database transaction tests.
  - Vitest component unit tests with simulated event-triggers.
  - Playwright E2E browser tests with isolated API endpoint routing.
- [x] **Phase 39: Declarative GitOps Reconciliation Engine & System-Wide UUID Migration:** Migrated all models to UUIDv4 to guarantee multi-site relational integrity, built a real-time GitOps YAML Template Generator inside `default.vue`s navigation sidebar, and established framework for mounting dry-run safe config manifests from remote Git repositories and local Kubernetes PVC paths.
- [x] **Phase 40: Dotted 1-Meter Blueprint Grid & Interactive CAD Floor Designer (Site Admin):** Replaced hard-lined grids with accurate dotted 1-meter metrics mapping arrays. Deployed an administrative "Room Layout Designer" dragging panel into `settings.vue` that lets admins visually place physical cabinets onto X/Y planes and seamlessly binds positional coordinates directly inside backend GORM records.
- [x] **Phase 41: Cabling Port Capacity Heatmap & Dynamic Capacity Warnings:** Added live Switch port utilization logic and visually striking, threshold-colored progress bar Heatmaps below network nodes on physical DCIM layout grids (dynamically throwing SFP+ capacity warnings above 80%).
- [x] **Phase 42: Administrative Disaster Recovery (JSON Backups):** Engineered a structural state compiler panel in `settings.vue`, wrapping unified Custom Fields, Users, Webhooks, and Locations configuration objects into downloadable, instantly restorable `northstar-backup.json` export files for guaranteed state resiliency.
- [x] **Phase 43: Public-Cloud Security Group Firewall Compliance Visualizer:** Integrated high-fidelity SecOps layout cards parsing connected AWS Crawler metadata to visually block and highlight severe violations on Container/Cloud topologies (displaying deep-red lock icons to warn whenever ports 22 and 3306 map openly to insecure `0.0.0.0/0` Internet IPs).
- [x] **Phase 44: Datacenter Location Filtering & Loose-Coupled Orphan Asset Auditing:** Added full-depth GORM nested preloads linking Assets cascadingly to parent Datacenter nodes. Overhauled `index.vue` with horizontal location filters and a specialized 'Location Mapping State' audit tool to isolate, list, and resolve loose-coupled production hardware with zero assigned datacenter owners.
- [x] **Phase 45: Wizard Sub-Stack Selection & IoT Edge Hardware Category:** Embedded a brand new, first-class **IoT / Edge** hardware category card inside the asset registration step-1 layout. Configured conditional, highly responsive **Sub-Stack Dropdown Selectors** for both Network (specifying Router, Switch L2/L3, Access Point, Firewall, Load Balancer) and IoT Edge nodes (specifying Raspberry Pi Gateways, Industrial PLCs, IP Cameras, and Environmental Sensors) mapping custom telemetry parameters dynamically inside GORM JSONMaps on the fly.
- [x] **Phase 46: Multi-Stage Nested Card Selection Wizard (Sub-Group Portals UI):** Fully re-engineered the Step-1 category card panel inside `AssetFormModal.vue` to introduce nested, smooth-transitioning sub-group selector grids. Clicking the 'Network' or 'IoT / Edge' cards now displays secondary visual portals allowing operators to select subtypes (Network -> Router, Network -> Switch L2/L3) utilizing custom SVG icons before filling out details.
- [x] **Phase 47: Dynamic Database-Backed Asset Category & Sub-Group Administrator:** Migrated entire asset category grids and sub-group structures from hardcoded frontend arrays into active, relational GORM models (`Category` and `SubGroup`). Programmed an administrative 'Category & Icon Manager' CRUD sub-page inside Site Admin allowing admins to create, edit, or delete entire Categories dynamically, and engineered GORM schema migrations to automatically auto-seed baseline categories and subtype associations on fresh database generations.
- [x] **Phase 48: Multi-Directory Modular YAML Bootstrapping & Directory Re-Organization:** Completely transitioned the database seeding engine to parse individual, single-resource YAML manifest files recursively inside `.northstar-data/bootstrap/categories/` and `.northstar-data/bootstrap/sub-groups/` (conforming to Kubernetes-style single-file schemas). Migrated all system namespaces, API groups, and generated templates to utilize the scoped product subdomain **`network.northstar.astrona.io/v1alpha1`**, supporting secure multi-product releases under parent owner **Astrona** (`astrona.io`).
- [x] **Phase 49: Declarative GitOps Users & RBAC Seeding:** Extended the K8s-style bootstrapping reconciler to parse and ingest declarative Users, Access Groups, and fine-grained Permissions manifest folders recursively on startup, securely hashing passwords using `bcrypt` and completely removing hardcoded seeding loops from the compiled Go binary.
- [x] **Phase 50: Domain-Aligned Multi-Directory Bootstrap Structure:** Reorganized all 18 standalone declarative seeding manifests under logical domain-specific resource sub-folders: grouping categories and sub-groups under `.northstar-data/bootstrap/asset-taxonomy/`, and users, groups, and permissions under `.northstar-data/bootstrap/iam-security/`. Updated the backend directory scanner in Go to resolve and parse these clean paths recursively, establishing ultimate human readability.
- [x] **Phase 51: Dynamic Attributes Mapping based on Custom Category Registration:** Fully bound the administrative Custom Fields creator card's category dropdown list to GORM database-backed Category names dynamically. Creating any bespoke, custom Category card inside Site Admin immediately enables registering and mapping corresponding metadata attribute schemas directly on Step 2 of the asset creation wizard.
- [x] **Phase 52: Dynamic Custom Attributes Layout Groups & Tabs:** Integrated dynamic, database-backed custom metadata **Tab Group Folders** across both backend and frontend layers. Exposed dynamic `tab_group` configuration parameters on `CustomField` models, allowing administrators to group attributes cleanly. Step-2 details forms in the asset creation modal dynamically compile and render customized, categorized horizontal tab selectors to completely manage high-density metadata without layout clutter.
- [x] **Phase 53: Custom Fields Validator Regex & Strict Pattern Enforcements:** Integrated secure, declarative pattern constraints on Custom Fields GORM schemas. Exposes an administrative `validation_regex` configuration field, letting system admins declare regex rules (e.g. strict subnets or UUIDs) that are dynamically compiled and verified on-the-fly inside the frontend asset wizard, throwing visual alerts and blocking invalid saves.
- [x] **Phase 54: Custom Field Bulk Importing & Regex Validations Pre-Checks:** Engineered a declarative, visual CSV Bulk Ingest Engine inside settings under Site Admin. Operators can upload a batch CSV inventory manifest. The frontend parser dynamically checks every single property value against active GORM Custom Field validation regex constraints on-the-fly, displaying a beautiful, inline red-alert validation report highlighting mismatches and strictly blocking database ingestion sweeps until all data cells are 100% compliant and clean.
- [x] **Phase 55: Automated Bulk Connection Patching via Cabling CSV:** Fully expanded the declarative CSV Bulk Ingest Engine to support a dedicated sub-toggle card, **"Cable Patch Matrix"**. Operators can upload an inter-rack fiber/copper cable manifest in CSV format. The ingestion pipeline dynamically queries and audits that both specified Source and Target hardware asset names exist in GORM before patching, highlighting mismatched nodes and self-linking loops inside an inline validation ledger, and completing bulk connection synchronizations asynchronously on success.
- [x] **Phase 56: Custom Category Dependency Chain Validation:** Engineered a dynamic, database-backed category prerequisite validation chain across GORM models, REST APIs, and UI layers. System administrators can configure prerequisite categories on dynamic Category entities (e.g., declaring that "Container" requires parent type "Server"). Both backend Echo handlers and the frontend creation wizard dynamically audit relationships during creations and partial updates, rendering visual informational banners, throwing strict pattern mismatched warnings, and blocking illegal or orphaned configurations on-the-fly.
- [x] **Phase 57: Dynamic Tree-Graph Visualizer for Taxonomy Dependency Chains:** Developed a visual horizontal **Taxonomy Dependency Map** inside settings under Site Admin. Dynamically computes GORM category root nodes and connects them visually to prerequisite child dynamic entities via dashed connector links. Glowing circular indicator badges parse live preloaded GORM hardware asset counts inside real-time animation nodes, providing system admins with an instant overview of their taxonomy footprint.
- [x] **Phase 58: Interactive Multi-Floor 2D CAD blueprint Wall Designer & Grid Alignment:** Engineered a full-stack, interactive 2D room wall sketching blueprint tool across GORM schemas, REST APIs, and UI panels. Adds custom `DatacenterWall` models registered cascadingly under parent floors (supporting Basement levels `-1` and upper floor plans). Operators can use a custom SVG wall segment brush tool inside settings, dragging and drafting walls onto 1-meter metrics, which dynamically persist to GORM and are rendered as solid physical boundary lines on both the Settings layout designer and the primary DCIM SVG blueprint.
- [x] **Phase 59: Muted GORM Record NotFound Console Warning Traces:** Configured a global GORM logger instance enabling `IgnoreRecordNotFoundError: true` on GORM DB Open configurations. This permanently mutes alarming SQL error dumps to your terminal console during fresh startup directory-walking checks and expected category/permissions lookups, keeping standard Go logging environments 100% quiet and professional.
