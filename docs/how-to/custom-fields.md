# How-To Guides: Custom Fields & Bulk Ingest

These guides provide problem-oriented, practical recipes to help you solve specific operational tasks within the Northstar platform.

---

## How to Declare a Custom Field Schema with Regex Constraints

### Problem
You want to enforce a strict formatting constraint (such as requiring IP addresses, MAC addresses, or precise serial numbers) on a dynamic custom attribute added to your assets.

### Solution
Use the administrative **Custom Fields Creator** under Site Admin to register the field with a custom compiled regular expression validator.

1. Navigate to **Settings** and select the **Site Administration** tab (GitLab-Style Sidebar).
2. Click on the **Custom Schema Attributes** panel.
3. Fill out the creation form:
    * **Field Name:** `rack_u_index`
    * **Value Type:** `number`
    * **Category Scope:** `Server`
    * **Validation Regex Pattern:** `^[1-9][0-9]?$` (Restricts values to numbers between 1 and 99)
    * **Form Tab Group:** `Hardware Metrics`
4. Click **Deploy Attribute Schema**.

The creation wizard will dynamically compile the regex rule, blocking non-compliant inputs inside the wizard on Step 2 with visual validation errors.

---

## How to Bulk Ingest Inventory CSV Files

### Problem
You have hundreds of newly acquired server rack chassis and network routers that you want to ingest into Northstar without registering them individually in the UI.

### Solution
Upload your CSV asset manifest via the **CSV Bulk Ingest Engine**, which runs pre-compliance checks before committing anything to GORM.

1. Ensure your CSV file is formatted correctly with the standard headers:
    ```csv
    name,type,status,ip_address,description,cpu_model,memory_gb
    Dublin-Router-01,Router,active,10.0.5.1,Core backhaul router,Broadcom DNX,64
    Dublin-Edge-02,Server,active,10.0.5.2,Edge server cluster,AMD EPYC,128
    ```
2. Navigate to **Site Admin** &rarr; **CSV Bulk Ingest Engine**.
3. Under **Target Data Type**, choose **CMDB Assets & Hardware**.
4. Drag and drop your `.csv` file into the upload zone.
5. Click **Run Pre-Compliance Checks**.
6. Review the inline **Verification Report**:
    * Green rows represent compliant records.
    * Red rows indicate cells violating active custom fields validation regex rules (e.g. invalid IPs or out-of-bounds metrics).
7. Resolve any highlights in your local spreadsheet, re-upload, and click **Commit Data Ingest Sweep** to write records to GORM.

---

## How to Bulk Connect Cabling Patches

### Problem
You want to map the physical copper/fiber cable matrix connecting routers and switch ports inside a 42U rack.

### Solution
Use the **Cabling Matrix** upload selector inside the CSV Bulk Ingest panel.

1. Create a `.csv` file detailing source and destination nodes:
    ```csv
    source_asset_name,source_port,target_asset_name,target_port,cable_type,link_speed_gbps
    Dublin-Router-01,eth1,Dublin-Switch-02,eth24,fiber,10
    Dublin-Router-01,eth2,Dublin-Switch-03,eth24,fiber,10
    ```
2. Upload the file under the **Cable Patch Matrix** toggle.
3. The parser will verify that both source and target hardware exist in GORM, automatically linking them in the relational mapping and rendering pulsing cabling tracers on the interactive SVG topology map!
