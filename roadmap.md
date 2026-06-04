# Northstar CMDB Refactoring & Enhancement Roadmap

This document outlines the architectural milestones, completed achievements, and advanced future roadmap features designed to establish a production-grade Configuration Management Database (CMDB) application.

---

## 🗺️ Future Vision: Active Enhancement Phases

The next generation of Northstar enhancements will target visual taxonomy linkages, thermal facility zoning, live port audits, and cloud-native observability.

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
  - **Debian/Alpine Package Scanners:** Integrate secure SSH crawler hooks to execute non-interactive commands (such as `dpkg-query -l` on Raspbian or `apk info -v` on Alpine Edge devices) to automatically parse installed library manifests.
  - **Physical SBOM Vulnerability Maps:** Map extracted edge packages directly onto the asset details SBOM tab, querying CVE databases in real-time to alert operators of security vulnerabilities inside physical remote field equipment.

### Phase 7: Real-Time Item-Level Declarative YAML Exporters (UI DX)
* **Goal:** Make GitOps completely tangible and user-friendly by adding a **"View GitOps YAML"** code portal directly onto individual asset details, cabinet, and datacenter panels.
* **Details:**
  - **Single-Item Code Compilation:** On the Asset Detail page (`/assets/[id]`) or Datacenter view (`/dcim`), add a clean **"View YAML"** button on the primary header or property card. Clicking the button opens a modal compiling that exact, specific database record on the fly into its valid, declarative YAML manifest representation.
  - **Looping UI to GitOps DX:** Provide a quick **"Copy to Clipboard"** button, allowing operators to design assets visually and copy their YAML specs into their GitOps repositories instantly.

### Phase 8: Zero-Trust Host Firewall Audits & OS-Level Rules Verification (SecOps L4)
* **Goal:** Evolve the CMDB topology from simple cabling maps to an active **Zero-Trust Network Posture map** by auditing and verifying host-level firewall configurations (such as Linux `firewalld`, `ufw`, or `iptables` rules).
* **Details:**
  - **OS-Level Firewall Ingestion:** Integrate secure SSH crawler hooks to execute non-interactive commands to pull active host-level ingress rules during auto-discovery scans or let operators declare port rules in GORM.
  - **Host Ingress Rule Matrix:** On the visual SVG topology map, clicking a connection line dynamically audits the target's firewall matrix to verify if traffic is actually allowed or blocked (rendering green flowing paths for verified allowed links, and red dashed warning paths if blocked).

### Phase 12: Declarative API Sandbox and Mock-free Live API Docs (API DX)
* **Goal:** Streamline developer integrations by serving a live interactive API playground directly from the admin dashboard.
* **Details:**
  - **Interactive OpenAPI v3 Playgrounds:** Inject Swagger/OpenAPI routing directly into Echo, allowing developers to test OAuth-scoped CRUD operations against the active database inside their browser with zero local configuration.

### Phase 13: Multi-Region Geographical Cable Mapping & Satellite Fiber Links (GIS WAN)
* **Goal:** Visualize global corporate networks by mapping cross-continental fiber paths and satellite WAN uplinks on a world map.
* **Details:**
  - **Geo-Mapping Canvas:** Integrate an interactive Leaflet.js or Mapbox world map. Draw geographical coordinates connecting datacenters, displaying active paths, undersea fiber cable outages, and Starlink gateway telemetry in real-time.

---

## 🏛️ Completed Milestones Summary

The Northstar platform has been successfully developed, refactored, and verified to production grade across 59 distinct evolutionary phases:

1. **Robust Go Backend & Relational Database (GORM & SQLite):** Migrated from legacy Python to a high-performance Go Echo REST API. Structured clean, human-logical modular folder structures separating routing handlers (`internal/router/router.go`) and business controllers (`internal/handlers/`) completely (Phase 11). Engineered version-led programmatic database migrations (`migrations.go`), and silent startup directory walking routines with muted `RecordNotFound` warnings.
2. **Kubernetes-Style Declarative Seeding (GitOps L2):** Built a recursive, modular K8s-style bootstrapping engine. Ingests categories, sub-groups, users, permissions, and group linkages declaratively from YAML files under the `network.northstar.astrona.io/v1alpha1` namespace.
3. **Advanced UX/UI & Nuxt 4 Architecture:** Migrated the frontend to a decoupled Nuxt 4 structure (`app/`). Engineered an elegant GitLab-style "Site Admin" utility hub, horizontal datacenter location filters, multi-stage 4-column creation wizards, interactive on-screen onboarding guides, and an **interactive HTML5 Drag-and-Drop Taxonomy visualizer** inside Site Admin (Phase 1) supporting drag-and-drop category prerequisite assignments with live, asynchronous GORM persistence updates and circular dependency prevention guards.
4. **Visual CAD Blueprinting & 3D Cabinet DCIM Canvas:** Designed interactive vertical 42U physical server cabinets featuring mounting rails and blinking status LEDs. Programmed SVG CAD floor plans with multi-floor selectors, precision 2D wall segment drawing tools, cross-cabinet copper/fiber patch ledgers, and automated Switch SFP+ port capacity heatmaps.
5. **Advanced RBAC, Dynamic Custom Fields & SecOps Auditing:** Secured write/delete APIs with JWT authentication. Designed dynamic category attributes and tabs, CSV bulk inventory and cabling patch importers with inline regex pre-checks, active background TCP subnet scanners, asynchronous thread-safe webhook publishers (Slack/Jira), and AWS security compliance group visualizers.
6. **Multi-Layer Verification, Real-Data E2E Test Suite & Diátaxis Developer Documentation:** Fully verified the codebase utilizing Go handler integration tests, Vitest component units, and a robust Playwright E2E browser test suite running against a dedicated E2E test database (`cmdb_e2e.db`). Bootstrapped a high-fidelity, Diátaxis-conforming MkDocs Material developer documentation portal (Phase 10) providing comprehensive tutorials, goal-oriented how-to recipes, REST API specifications, and Mermaid.js system architecture flowcharts fully aligned with Astrona corporate branding guidelines.
7. **OpenTelemetry APM & Distributed Observability (OTel L3):** Engineered complete OpenTelemetry (OTel) Tracing SDK initialization into the server startup lifecycle (Phase 9). Built a custom, lightweight, highly optimized GORM database tracing plugin (`GormOTelPlugin`) that automatically hooks SQLite query and transaction execution callbacks, and enabled standard Echo HTTP request tracing middleware to map incoming endpoint traffic to backend trace Spans.
