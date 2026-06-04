# REST API Endpoint Reference

This reference document outlines the REST API endpoints, payload contracts, authentication, and HTTP response schemas exposed by the Northstar Echo backend on port `8000`.

---

## Base Configuration & Headers

All REST operations utilize JSON bodies. Endpoints modifying the state (POST, PUT, DELETE) require a JWT bearer authentication token in the request header.

* **Base URL:** `http://localhost:8000`
* **Content-Type:** `application/json`
* **Authorization Header:** `Bearer <token>`

---

## 🔐 Authentication

### POST `/api/auth/login`
Generates a signed JWT session token valid for 24 hours.

#### Request Body
```json
{
  "username": "admin",
  "password": "admin"
}
```

#### Successful Response (`200 OK`)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "role": "admin"
}
```

---

## 🖥️ CMDB Assets

### GET `/api/assets`
Retrieve a list of registered hardware and cloud-native assets.

#### Query Parameters
* `q` (string, optional): Searches asset names and descriptions.
* `type` (string, optional): Filters by specific asset category (e.g., `Server`, `Database`).
* `status` (string, optional): Filters by asset operational state (e.g., `active`, `maintenance`).

#### Successful Response (`200 OK`)
```json
[
  {
    "id": "e2f75a6c-b248-4cb2-a7cf-3467bf2ea904",
    "name": "Production DB Server",
    "type": "Database",
    "status": "active",
    "ip_address": "10.0.1.5",
    "description": "Main PostgreSQL database",
    "properties": {
      "engine": "PostgreSQL",
      "version": "14.2"
    },
    "maintenance_status": "none"
  }
]
```

### POST `/api/assets`
Creates a new configuration asset in the database. Requires `asset:write` permission.

#### Request Body
```json
{
  "name": "Edge-Switch-A",
  "type": "Network",
  "status": "active",
  "ip_address": "10.0.1.200",
  "description": "Rack A core switch",
  "properties": {
    "sub_group": "Switch (L3)",
    "ports": 48
  }
}
```

#### Successful Response (`201 Created`)
```json
{
  "id": "2a75904c-3ca1-4091-bf99-0d1234901f4c",
  "name": "Edge-Switch-A",
  "type": "Network",
  "status": "active",
  "ip_address": "10.0.1.200",
  "properties": {
    "sub_group": "Switch (L3)",
    "ports": 48
  }
}
```

---

## 🏢 DCIM Datacenters

### GET `/api/datacenters`
Retrieves all registered physical location profiles including nested floor blueprint relations.

#### Successful Response (`200 OK`)
```json
[
  {
    "id": "86a08489-ef65-45b5-814e-90a6900ea458",
    "name": "Northstar-Dublin-HQ",
    "type": "homelab",
    "country": "Ireland",
    "city": "Dublin",
    "properties": {
      "uplink_speed": "2.5 Gbps Fiber WAN"
    },
    "floors": [
      {
        "id": "f2901a4c-bf12-4091-88fc-ef12a45c90f2",
        "datacenter_id": "86a08489-ef65-45b5-814e-90a6900ea458",
        "name": "Floor 1",
        "level": 1,
        "width": 20,
        "depth": 20
      }
    ]
  }
]
```
