# LatestPack API Documentation

Base URL: `/api`

All requests and responses use `Content-Type: application/json`, except file upload which uses `multipart/form-data`.

Authentication: all endpoints except login require `Authorization: Bearer <token>` header.

---

## 1. Authentication

### POST /api/auth/login

User login, returns JWT token.

**Request Body:**

| Field | Type | Required | Description |
|---|---|---|---|
| username | string | Yes | Username |
| password | string | Yes | Password |

**Response 200:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "username": "admin"
}
```

**Response 401:**

```json
{
  "error": "Invalid credentials"
}
```

---

## 2. Dashboard

### GET /api/dashboard/stats

Returns today's statistics.

**Response 200:**

```json
{
  "launches": {
    "value": 372,
    "change": 8.1
  },
  "updates": {
    "value": 56,
    "change": 3.4
  },
  "traffic": {
    "value": 2847,
    "unit": "MB",
    "change": 15.7
  }
}
```

| Field | Type | Description |
|---|---|---|
| launches.value | number | Today's launch count |
| launches.change | number | Percentage change compared to yesterday (positive = increase) |
| updates.value | number | Today's update count |
| updates.change | number | Percentage change |
| traffic.value | number | Traffic in the given unit |
| traffic.unit | string | Unit, always `"MB"` |
| traffic.change | number | Percentage change |

### GET /api/dashboard/latest-version

Returns the latest version announcement.

**Response 200:**

```json
{
  "version": "v1.2.0",
  "date": "2026-06-20",
  "size": "24.3 MB",
  "notes": [
    "优化启动速度，冷启动耗时降低 40%",
    "修复了在低内存设备上的崩溃问题",
    "新增增量更新机制，减少下载流量",
    "改进日志系统，支持远程上报"
  ]
}
```

| Field | Type | Description |
|---|---|---|
| version | string | Semantic version string |
| date | string | Release date (ISO 8601) |
| size | string | Package size (human-readable) |
| notes | string[] | Changelog entries |

---

## 3. Versions

### GET /api/versions

Returns paginated version list, newest first.

**Query Parameters:**

| Parameter | Type | Default | Description |
|---|---|---|---|
| page | number | 1 | Page number (1-based) |
| pageSize | number | 5 | Items per page |

**Response 200:**

```json
{
  "items": [
    {
      "version": "v1.2.0",
      "date": "2026-06-20",
      "size": "24.3 MB",
      "notes": [
        "优化启动速度，冷启动耗时降低 40%",
        "修复了在低内存设备上的崩溃问题",
        "新增增量更新机制，减少下载流量",
        "改进日志系统，支持远程上报"
      ]
    }
  ],
  "total": 11,
  "page": 1,
  "pageSize": 5
}
```

| Field | Type | Description |
|---|---|---|
| items | Version[] | Version entries for the current page |
| items[].version | string | Semantic version string |
| items[].date | string | Release date (ISO 8601) |
| items[].size | string | Package size (human-readable) |
| items[].notes | string[] | Changelog entries |
| total | number | Total number of versions |
| page | number | Current page number |
| pageSize | number | Items per page |

### POST /api/versions

Create a new version. Only the latest version can exist at the top; newly created versions become the latest automatically.

**Request Body:**

| Field | Type | Required | Description |
|---|---|---|---|
| version | string | Yes | Semantic version string, e.g. `"v1.3.0"` |
| notes | string[] | Yes | Changelog entries |

**Request Example:**

```json
{
  "version": "v1.3.0",
  "notes": [
    "新增自动回滚功能",
    "优化下载速度"
  ]
}
```

**Response 201:**

```json
{
  "version": "v1.3.0",
  "date": "2026-06-21",
  "size": "0 MB",
  "notes": [
    "新增自动回滚功能",
    "优化下载速度"
  ]
}
```

**Response 400:**

```json
{
  "error": "Version already exists"
}
```

### DELETE /api/versions/:version

Delete a version. Only the latest version can be deleted.

**URL Parameters:**

| Parameter | Type | Description |
|---|---|---|
| version | string | Version string, e.g. `v1.2.0` |

**Response 200:**

```json
{
  "message": "Deleted"
}
```

**Response 403:**

```json
{
  "error": "Only the latest version can be deleted"
}
```

**Response 404:**

```json
{
  "error": "Version not found"
}
```

---

## 4. Files

### GET /api/files

List files and folders in a directory.

**Query Parameters:**

| Parameter | Type | Default | Description |
|---|---|---|---|
| path | string | `""` | Directory path, empty string for root |

**Response 200:**

```json
{
  "path": "",
  "items": [
    {
      "name": "releases",
      "type": "folder",
      "size": "-",
      "date": "2026-06-20"
    },
    {
      "name": "README.md",
      "type": "file",
      "size": "2.4 KB",
      "date": "2026-06-01"
    }
  ]
}
```

| Field | Type | Description |
|---|---|---|
| path | string | Current directory path |
| items | FileItem[] | Items in the directory, folders first sorted alphabetically |
| items[].name | string | File or folder name |
| items[].type | string | `"folder"` or `"file"` |
| items[].size | string | Human-readable size, `"-"` for folders |
| items[].date | string | Last modified date (ISO 8601) |

### POST /api/files/folder

Create a new folder.

**Request Body:**

| Field | Type | Required | Description |
|---|---|---|---|
| path | string | Yes | Parent directory path, empty string for root |
| name | string | Yes | New folder name |

**Request Example:**

```json
{
  "path": "",
  "name": "logs"
}
```

**Response 201:**

```json
{
  "name": "logs",
  "type": "folder",
  "size": "-",
  "date": "2026-06-21"
}
```

**Response 409:**

```json
{
  "error": "Item already exists"
}
```

### POST /api/files/upload

Upload a file to a directory.

**Request:** `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| path | string | Yes | Target directory path, empty string for root |
| file | file | Yes | The file to upload |

**Response 201:**

```json
{
  "name": "update.json",
  "type": "file",
  "size": "0.5 KB",
  "date": "2026-06-21"
}
```

**Response 409:**

```json
{
  "error": "File already exists"
}
```

### PUT /api/files/rename

Rename a file or folder.

**Request Body:**

| Field | Type | Required | Description |
|---|---|---|---|
| path | string | Yes | Parent directory path |
| oldName | string | Yes | Current name |
| newName | string | Yes | New name |

**Request Example:**

```json
{
  "path": "config",
  "oldName": "update.json",
  "newName": "update-v2.json"
}
```

**Response 200:**

```json
{
  "name": "update-v2.json",
  "type": "file",
  "size": "0.5 KB",
  "date": "2026-06-21"
}
```

**Response 409:**

```json
{
  "error": "Item already exists"
}
```

**Response 404:**

```json
{
  "error": "Item not found"
}
```

### DELETE /api/files

Delete a file or folder. Folders are deleted recursively.

**Request Body:**

| Field | Type | Required | Description |
|---|---|---|---|
| path | string | Yes | Parent directory path |
| name | string | Yes | Name of the file or folder to delete |

**Request Example:**

```json
{
  "path": "",
  "name": "old-reports"
}
```

**Response 200:**

```json
{
  "message": "Deleted"
}
```

**Response 404:**

```json
{
  "error": "Item not found"
}
```

---

## 5. Error Response Format

All error responses follow a consistent format:

```json
{
  "error": "Description of the error"
}
```

Common HTTP status codes:

| Code | Meaning |
|---|---|
| 400 | Bad request (validation error, missing required field) |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden (operation not allowed, e.g. deleting non-latest version) |
| 404 | Resource not found |
| 409 | Conflict (duplicate name, version already exists) |
| 500 | Internal server error |

---

## 6. Endpoint Summary

| Method | Endpoint | Description |
|---|---|---|
| POST | /api/auth/login | User login |
| GET | /api/dashboard/stats | Today's statistics |
| GET | /api/dashboard/latest-version | Latest version announcement |
| GET | /api/versions | Paginated version list |
| POST | /api/versions | Create new version |
| DELETE | /api/versions/:version | Delete latest version |
| GET | /api/files | List directory contents |
| POST | /api/files/folder | Create folder |
| POST | /api/files/upload | Upload file (multipart) |
| PUT | /api/files/rename | Rename file or folder |
| DELETE | /api/files | Delete file or folder |
