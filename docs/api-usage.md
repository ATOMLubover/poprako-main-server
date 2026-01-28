# API 使用文档

## 全局说明

**Base URL**: `/api/v1`

**认证**: 大部分端点需要 `Authorization: Bearer <token>` 请求头或 `Authorization` Cookie。登录成功后返回 token 并自动设置 Cookie。

**响应格式**: 统一包装为 `{"code": <uint16>, "msg": "<string>", "data": <T|null>}`

- 成功时 `data` 有值，`msg` 为空
- 失败时 `msg` 有值，`data` 为 null

---

## 常见 HTTP 状态码

| 状态码 | 含义                                    |
| ------ | --------------------------------------- |
| 200    | 成功（OK）                              |
| 201    | 创建成功（Created）                     |
| 204    | 成功但无返回内容（No Content）          |
| 400    | 请求参数错误（Bad Request）             |
| 401    | 未认证或令牌无效（Unauthorized）        |
| 403    | 权限不足（Forbidden）                   |
| 404    | 资源不存在（Not Found）                 |
| 500    | 服务器内部错误（Internal Server Error） |

**注**: 响应 JSON 中的 `code` 字段与 HTTP 状态码一致，`msg` 字段包含具体错误信息（中文）。

---

## JSON 字段类型

- `string`: 字符串
- `int` / `int64`: 整数
- `bool`: 布尔值
- `float64`: 浮点数
- `array`: 数组
- `object`: 对象

**时间戳**: 所有时间戳字段（如 `created_at`、`updated_at`、`*_started_at`、`*_completed_at`）均为 **整数，单位为秒（Unix timestamp）**。

---

## PATCH（部分更新）语义

**规则**:

- **为 null 或不携带的字段不会被更新**
- **一旦携带（即在请求 JSON 中出现），该字段必然会被更新为新值**
- **前端绝对不允许在字段值未发生改变的情况下携带冗余字段，以防竞态条件导致意外覆盖**

**示例**: 若要清空可选字符串字段，发送空字符串 `""`；若要保持不变，则完全省略该字段或发送 `null`。

**适用端点**: `PATCH /users/{id}`, `PATCH /worksets/{id}`, `PATCH /comics/{id}`, `PATCH /pages/{id}`, `PATCH /pages/{page_id}/units`, `PATCH /assignments/{id}`

---

## 端点列表

### 1. 客户端版本检查

**`GET /api/v1/check-update`** — 检查客户端版本是否需要更新  
**认证**: 无  
**Headers**: `X-Client-App-Version: string` (必需)  
**响应**: `{code, msg, data: {latest_version, title, description, allow_usage}}`

```bash
curl -X GET "http://localhost/api/v1/check-update" -H "X-Client-App-Version: 1.2.3"
```

---

### 2. 用户认证与管理

#### **`POST /api/v1/login`** — 用户登录

**认证**: 无  
**Body**: `{qq: string, password: string, nickname?: string, invitation_code?: string}`  
**响应**: `{code, msg, data: {token: string}}` + 设置 `Authorization` Cookie

```bash
curl -X POST "http://localhost/api/v1/login" -H "Content-Type: application/json" -d '{"qq":"12345","password":"pwd"}'
```

#### **`GET /api/v1/users/me`** — 获取当前用户信息

**认证**: 必需  
**响应**: `{code, msg, data: UserInfo}`

```bash
curl -X GET "http://localhost/api/v1/users/me" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/users`** — 查询用户列表

**认证**: 必需  
**Query**: `offset=0&limit=10&nn=<nickname>&qq=<qq>&ia=<bool>&itsl=<bool>&ipr=<bool>...`  
**响应**: `{code, msg, data: [UserInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/users?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/users/{user_id}`** — 获取指定用户信息

**认证**: 必需  
**Path**: `user_id: string`  
**响应**: `{code, msg, data: UserInfo}` | `404`

```bash
curl -X GET "http://localhost/api/v1/users/u123" -H "Authorization: Bearer <token>"
```

#### **`PATCH /api/v1/users/{user_id}`** — 更新用户信息

**认证**: 必需  
**Path**: `user_id: string`  
**Body**: `{id: string, qq?: string, nickname?: string}` — `id` 必须与路径匹配  
**可更新字段**: `qq`, `nickname` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/users/u1" -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"id":"u1","nickname":"new"}'
```

#### **`PATCH /api/v1/users/{user_id}/roles`** — 分配/取消用户角色

**认证**: 必需  
**Path**: `user_id: string`  
**Body**: `{id: string, roles: [{role: string, assigned: bool}, ...]}`  
**可用角色**: `translator`, `proofreader`, `typesetter`, `redrawer`, `reviewer`, `uploader`  
**响应**: `204 No Content` | `403`

```bash
curl -X PATCH "http://localhost/api/v1/users/u1/roles" -H "Authorization: Bearer <token>" -d '{"id":"u1","roles":[{"role":"translator","assigned":true}]}'
```

---

### 3. 邀请管理

#### **`POST /api/v1/users/invitations`** — 创建邀请码

**认证**: 必需  
**Body**: `{invitee_qq: string, assign_translator?: bool, assign_proofreader?: bool, ...}`  
**响应**: `{code, msg, data: {invitation_code: string}}`

```bash
curl -X POST "http://localhost/api/v1/users/invitations" -H "Authorization: Bearer <token>" -d '{"invitee_qq":"99999"}'
```

#### **`GET /api/v1/users/invitations`** — 获取我创建的邀请列表

**认证**: 必需  
**响应**: `{code, msg, data: [InvitationInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/users/invitations" -H "Authorization: Bearer <token>"
```

---

### 4. 工作集（Workset）管理

#### **`GET /api/v1/worksets`** — 获取工作集列表

**认证**: 必需  
**Query**: `offset=0&limit=10`  
**响应**: `{code, msg, data: [WorksetInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/worksets?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/worksets/{workset_id}`** — 获取工作集详情

**认证**: 必需  
**Path**: `workset_id: string`  
**响应**: `{code, msg, data: WorksetInfo}` | `404`

```bash
curl -X GET "http://localhost/api/v1/worksets/w1" -H "Authorization: Bearer <token>"
```

#### **`POST /api/v1/worksets`** — 创建工作集

**认证**: 必需（仅管理员）  
**Body**: `{name: string, description?: string}`  
**响应**: `{code, msg, data: {id: string}}`

```bash
curl -X POST "http://localhost/api/v1/worksets" -H "Authorization: Bearer <token>" -d '{"name":"NewWS"}'
```

#### **`PATCH /api/v1/worksets/{workset_id}`** — 更新工作集

**认证**: 必需  
**Path**: `workset_id: string`  
**Body**: `{id?: string, description?: string}`  
**可更新字段**: `description` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/worksets/w1" -H "Authorization: Bearer <token>" -d '{"description":"new desc"}'
```

#### **`DELETE /api/v1/worksets/{workset_id}`** — 删除工作集

**认证**: 必需  
**Path**: `workset_id: string`  
**响应**: `204 No Content` | `404`

```bash
curl -X DELETE "http://localhost/api/v1/worksets/w1" -H "Authorization: Bearer <token>"
```

---

### 5. 漫画（Comic）管理

#### **`GET /api/v1/comics`** — 搜索/查询漫画列表

**认证**: 必需  
**Query**: `offset=0&limit=10&aut=<author>&tit=<title>&wid=<workset_id>&auid=<assigned_user_id>&tsl_pending=<bool>&tsl_wip=<bool>...`  
**响应**: `{code, msg, data: [ComicBrief, ...]}`

```bash
curl -X GET "http://localhost/api/v1/comics?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/comics/{comic_id}`** — 获取漫画详情

**认证**: 必需  
**Path**: `comic_id: string`  
**响应**: `{code, msg, data: ComicInfo}` | `404`

```bash
curl -X GET "http://localhost/api/v1/comics/c1" -H "Authorization: Bearer <token>"
```

#### **`POST /api/v1/comics`** — 创建漫画

**认证**: 必需  
**Body**: `{workset_id: string, author: string, title: string, description?: string, comment?: string, pre_asgns?: [...]}`  
**响应**: `{code, msg, data: {id: string}}`

```bash
curl -X POST "http://localhost/api/v1/comics" -H "Authorization: Bearer <token>" -d '{"workset_id":"w1","author":"a","title":"t"}'
```

#### **`PATCH /api/v1/comics/{comic_id}`** — 更新漫画元数据

**认证**: 必需  
**Path**: `comic_id: string`  
**Body**: `{id?: string, author?: string, title?: string, description?: string, comment?: string}`  
**可更新字段**: `author`, `title`, `description`, `comment` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/comics/c1" -H "Authorization: Bearer <token>" -d '{"title":"NewTitle"}'
```

#### **`DELETE /api/v1/comics/{comic_id}`** — 删除漫画

**认证**: 必需  
**Path**: `comic_id: string`  
**响应**: `204 No Content` | `404`

```bash
curl -X DELETE "http://localhost/api/v1/comics/c1" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/worksets/{workset_id}/comics`** — 获取工作集下的漫画列表

**认证**: 必需  
**Path**: `workset_id: string`  
**Query**: `offset=0&limit=10`  
**响应**: `{code, msg, data: [ComicBrief, ...]}`

```bash
curl -X GET "http://localhost/api/v1/worksets/w1/comics?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

---

### 6. 页面（Page）管理

#### **`GET /api/v1/pages/{page_id}`** — 获取页面详情

**认证**: 必需  
**Path**: `page_id: string`  
**响应**: `{code, msg, data: ComicPageInfo}`

```bash
curl -X GET "http://localhost/api/v1/pages/p1" -H "Authorization: Bearer <token>"
```

#### **`POST /api/v1/pages`** — 批量创建页面

**认证**: 必需  
**Body**: `[{comic_id: string, index: int64, image_ext: string}, ...]`  
**响应**: `{code, msg, data: [{id: string, oss_url: string}, ...]}`

```bash
curl -X POST "http://localhost/api/v1/pages" -H "Authorization: Bearer <token>" -d '[{"comic_id":"c1","index":1,"image_ext":"jpg"}]'
```

#### **`PATCH /api/v1/pages/{page_id}`** — 更新页面状态

**认证**: 必需  
**Path**: `page_id: string`  
**Body**: `{id?: string, uploaded?: bool}`  
**可更新字段**: `uploaded` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/pages/p1" -H "Authorization: Bearer <token>" -d '{"uploaded":true}'
```

#### **`DELETE /api/v1/pages/{page_id}`** — 删除页面

**认证**: 必需  
**Path**: `page_id: string`  
**响应**: `204 No Content`

```bash
curl -X DELETE "http://localhost/api/v1/pages/p1" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/comics/{comic_id}/pages`** — 获取漫画的页面列表

**认证**: 必需  
**Path**: `comic_id: string`  
**响应**: `{code, msg, data: [ComicPageInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/comics/c1/pages" -H "Authorization: Bearer <token>"
```

---

### 7. 翻译单元（Unit）管理

#### **`GET /api/v1/pages/{page_id}/units`** — 获取页面的翻译单元列表

**认证**: 必需  
**Path**: `page_id: string`  
**响应**: `{code, msg, data: [ComicUnitInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/pages/p1/units" -H "Authorization: Bearer <token>"
```

#### **`POST /api/v1/pages/{page_id}/units`** — 批量创建翻译单元

**认证**: 必需  
**Path**: `page_id: string`  
**Body**: `[{page_id: string, index: int64, x_coordinate: float64, y_coordinate: float64, is_in_box: bool, translated_text?: string, ...}, ...]`  
**响应**: `201 Created`

```bash
curl -X POST "http://localhost/api/v1/pages/p1/units" -H "Authorization: Bearer <token>" -d '[{"page_id":"p1","index":1,"x_coordinate":1.0,"y_coordinate":2.0,"is_in_box":true}]'
```

#### **`PATCH /api/v1/pages/{page_id}/units`** — 批量更新翻译单元

**认证**: 必需  
**Path**: `page_id: string`  
**Body**: `[{id: string, index?: int64, x_coordinate?: float64, y_coordinate?: float64, is_in_box?: bool, translated_text?: string, translator_comment?: string, proved_text?: string, proved?: bool, proofreader_comment?: string}, ...]`  
**可更新字段**: `index`, `x_coordinate`, `y_coordinate`, `is_in_box`, `translated_text`, `translator_comment`, `proved_text`, `proved`, `proofreader_comment` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/pages/p1/units" -H "Authorization: Bearer <token>" -d '[{"id":"u1","translated_text":"new"}]'
```

#### **`DELETE /api/v1/pages/{page_id}/units`** — 批量删除翻译单元

**认证**: 必需  
**Path**: `page_id: string`  
**Body**: `["unit_id1", "unit_id2", ...]`  
**响应**: `204 No Content`

```bash
curl -X DELETE "http://localhost/api/v1/pages/p1/units" -H "Authorization: Bearer <token>" -d '["u1","u2"]'
```

---

### 8. 任务分配（Assignment）管理

#### **`GET /api/v1/assignments/{asgn_id}`** — 获取分配详情

**认证**: 必需  
**Path**: `asgn_id: string`  
**响应**: `{code, msg, data: ComicAsgnInfo}` | `404`

```bash
curl -X GET "http://localhost/api/v1/assignments/a1" -H "Authorization: Bearer <token>"
```

#### **`POST /api/v1/assignments`** — 创建任务分配

**认证**: 必需  
**Body**: `{comic_id: string, assignee_id: string, is_translator?: bool, is_proofreader?: bool, ...}`  
**响应**: `{code, msg, data: {...}}`

```bash
curl -X POST "http://localhost/api/v1/assignments" -H "Authorization: Bearer <token>" -d '{"comic_id":"c1","assignee_id":"u1"}'
```

#### **`PATCH /api/v1/assignments/{asgn_id}`** — 更新分配角色

**认证**: 必需  
**Path**: `asgn_id: string`  
**Body**: `{id?: string, is_translator?: bool, is_proofreader?: bool, is_typesetter?: bool, is_redrawer?: bool, is_reviewer?: bool}`  
**可更新字段**: `is_translator`, `is_proofreader`, `is_typesetter`, `is_redrawer`, `is_reviewer` (PATCH语义：省略/null=不更新；携带=更新)  
**响应**: `204 No Content`

```bash
curl -X PATCH "http://localhost/api/v1/assignments/a1" -H "Authorization: Bearer <token>" -d '{"is_translator":true}'
```

#### **`DELETE /api/v1/assignments/{asgn_id}`** — 删除分配

**认证**: 必需  
**Path**: `asgn_id: string`  
**响应**: `204 No Content`

```bash
curl -X DELETE "http://localhost/api/v1/assignments/a1" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/comics/{comic_id}/assignments`** — 获取漫画的分配列表

**认证**: 必需  
**Path**: `comic_id: string`  
**Query**: `offset=0&limit=10`  
**响应**: `{code, msg, data: [ComicAsgnInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/comics/c1/assignments?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

#### **`GET /api/v1/users/{user_id}/assignments`** — 获取用户的分配列表

**认证**: 必需  
**Path**: `user_id: string`  
**Query**: `offset=0&limit=10`  
**响应**: `{code, msg, data: [ComicAsgnInfo, ...]}`

```bash
curl -X GET "http://localhost/api/v1/users/u1/assignments?offset=0&limit=10" -H "Authorization: Bearer <token>"
```

---

## 数据模型示例

### UserInfo

```json
{
  "id": "string",
  "qq": "string",
  "nickname": "string",
  "assigned_translator_at": int64 | null,
  "assigned_proofreader_at": int64 | null,
  "assigned_typesetter_at": int64 | null,
  "assigned_redrawer_at": int64 | null,
  "assigned_reviewer_at": int64 | null,
  "assigned_uploader_at": int64 | null,
  "is_admin": bool,
  "created_at": int64
}
```

### WorksetInfo

```json
{
  "id": "string",
  "index": int64,
  "name": "string",
  "comic_count": int64,
  "description": "string" | null,
  "creator_id": "string",
  "creator_nickname": "string",
  "created_at": int64,
  "updated_at": int64
}
```

### ComicInfo

```json
{
  "id": "string",
  "workset_id": "string",
  "workset_index": int,
  "index": int64,
  "creator_id": "string",
  "creator_nickname": "string",
  "author": "string",
  "title": "string",
  "description": "string" | null,
  "comment": "string" | null,
  "page_count": int64,
  "translating_started_at": int64 | null,
  "translating_completed_at": int64 | null,
  "proofreading_started_at": int64 | null,
  "proofreading_completed_at": int64 | null,
  "typesetting_started_at": int64 | null,
  "typesetting_completed_at": int64 | null,
  "reviewing_completed_at": int64 | null,
  "uploading_completed_at": int64 | null,
  "created_at": int64,
  "updated_at": int64
}
```

### ComicPageInfo

```json
{
  "id": "string",
  "comic_id": "string",
  "index": int64,
  "oss_url": "string",
  "uploaded": bool,
  "inbox_unit_count": int64,
  "outbox_unit_count": int64,
  "translated_unit_count": int64,
  "proved_unit_count": int64
}
```

### ComicUnitInfo

```json
{
  "id": "string",
  "page_id": "string",
  "index": int64,
  "x_coordinate": float64,
  "y_coordinate": float64,
  "is_in_box": bool,
  "translated_text": "string" | null,
  "translator_id": "string" | null,
  "translator_comment": "string" | null,
  "proved_text": "string" | null,
  "proved": bool,
  "proofreader_id": "string" | null,
  "proofreader_comment": "string" | null,
  "creator_id": "string" | null,
  "created_at": int64,
  "updated_at": int64
}
```
