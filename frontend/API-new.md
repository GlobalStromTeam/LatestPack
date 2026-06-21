# LatestPack 新增 API 文档

Base URL: `/api`

所有请求和响应使用 `Content-Type: application/json`。

认证：所有接口需要 `Authorization: Bearer <token>` 请求头。

---

## 1. 渠道管理（Channels）

### GET /api/channels

获取所有渠道列表。

**Response 200:**

```json
[
  {
    "id": "ch_local",
    "name": "本地存储",
    "type": "local",
    "enabled": true,
    "weight": 0,
    "config": {}
  },
  {
    "id": "ch_webdav_1",
    "name": "WebDAV 备份",
    "type": "webdav",
    "enabled": true,
    "weight": 50,
    "config": {
      "endpoint": "https://dav.example.com",
      "path": "/backups",
      "accessKey": "user",
      "secretKey": "pass"
    }
  },
  {
    "id": "ch_s3_1",
    "name": "S3 存储",
    "type": "s3",
    "enabled": false,
    "weight": 30,
    "config": {
      "endpoint": "https://s3.amazonaws.com",
      "bucket": "my-bucket",
      "region": "us-east-1",
      "path": "packs",
      "accessKey": "AKIA...",
      "secretKey": "..."
    }
  }
]
```

**字段说明：**

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | 渠道唯一标识 |
| name | string | 渠道显示名称 |
| type | string | `"local"`、`"webdav"` 或 `"s3"` |
| enabled | boolean | 是否启用 |
| weight | number | 权重（1-100），本地渠道固定为 0 |
| config.endpoint | string | 远程地址（WebDAV / S3） |
| config.bucket | string | S3 Bucket 名称 |
| config.region | string | S3 区域 |
| config.path | string | 远程路径前缀 |
| config.accessKey | string | WebDAV 用户名 / S3 Access Key |
| config.secretKey | string | WebDAV 密码 / S3 Secret Key |

**不同渠道类型的 config 字段：**

| 字段 | local | webdav | s3 |
|---|---|---|---|
| endpoint | - | 必填 | 必填 |
| path | - | 可选 | 可选 |
| accessKey | - | 必填 | 必填 |
| secretKey | - | 必填 | 必填 |
| bucket | - | - | 必填 |
| region | - | - | 可选 |

---

### POST /api/channels

创建新渠道。本地渠道由系统自动创建，不可手动添加。

**Request Body:**

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 渠道名称 |
| type | string | 是 | `"webdav"` 或 `"s3"` |
| enabled | boolean | 否 | 默认 `true` |
| weight | number | 否 | 默认 `50`，范围 1-100 |
| config | object | 是 | 对应类型的配置 |

**请求示例（WebDAV）：**

```json
{
  "name": "WebDAV 备份",
  "type": "webdav",
  "enabled": true,
  "weight": 50,
  "config": {
    "endpoint": "https://dav.example.com",
    "path": "/backups",
    "accessKey": "user",
    "secretKey": "pass"
  }
}
```

**请求示例（S3）：**

```json
{
  "name": "S3 存储",
  "type": "s3",
  "enabled": true,
  "weight": 30,
  "config": {
    "endpoint": "https://s3.amazonaws.com",
    "bucket": "my-bucket",
    "region": "us-east-1",
    "path": "packs",
    "accessKey": "AKIA...",
    "secretKey": "..."
  }
}
```

**Response 201:** 返回创建的渠道对象。

**Response 409:**

```json
{
  "error": "Channel name already exists"
}
```

---

### PUT /api/channels/:id

更新渠道。本地渠道仅允许修改 `enabled` 字段。

**URL 参数：**

| 参数 | 类型 | 说明 |
|---|---|---|
| id | string | 渠道 ID |

**Request Body:** 传入需要更新的字段（部分更新）。

**请求示例：**

```json
{
  "name": "WebDAV 备份（更新）",
  "enabled": false,
  "weight": 70,
  "config": {
    "endpoint": "https://dav2.example.com",
    "path": "/new-backups",
    "accessKey": "newuser",
    "secretKey": "newpass"
  }
}
```

**Response 200:** 返回更新后的渠道对象。

**Response 404:**

```json
{
  "error": "Channel not found"
}
```

---

### DELETE /api/channels/:id

删除渠道。本地渠道不可删除。

**URL 参数：**

| 参数 | 类型 | 说明 |
|---|---|---|
| id | string | 渠道 ID |

**Response 200:**

```json
{
  "message": "Deleted"
}
```

**Response 403:**

```json
{
  "error": "Cannot delete local channel"
}
```

**Response 404:**

```json
{
  "error": "Channel not found"
}
```

---

## 2. 账户设置（Auth Settings）

### PUT /api/auth/username

修改当前用户的账户名。

**Request Body:**

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| username | string | 是 | 新账户名 |

**请求示例：**

```json
{
  "username": "newname"
}
```

**Response 200:**

```json
{
  "username": "newname"
}
```

**Response 409:**

```json
{
  "error": "Username already exists"
}
```

---

### PUT /api/auth/password

修改当前用户的密码。修改成功后当前 Token 失效，需重新登录。

**Request Body:**

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| currentPassword | string | 是 | 当前密码，用于验证 |
| newPassword | string | 是 | 新密码，最少 6 位 |

**请求示例：**

```json
{
  "currentPassword": "oldpass123",
  "newPassword": "newpass456"
}
```

**Response 200:**

```json
{
  "message": "Password updated"
}
```

**Response 400:**

```json
{
  "error": "Current password is incorrect"
}
```

---

## 3. 接口汇总

| 方法 | 接口 | 说明 |
|---|---|---|
| GET | /api/channels | 获取所有渠道 |
| POST | /api/channels | 创建渠道 |
| PUT | /api/channels/:id | 更新渠道 |
| DELETE | /api/channels/:id | 删除渠道 |
| PUT | /api/auth/username | 修改账户名 |
| PUT | /api/auth/password | 修改密码 |
