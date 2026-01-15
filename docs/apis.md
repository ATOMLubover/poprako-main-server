# API 参考（简要）

以下为项目中 HTTP 接口的简要参考，路由基准路径为 `/api/v1`。统一响应格式为：

**统一响应格式**

```json
{
  "code": <uint16>,
  "msg": "<错误信息，如果成功则为空>",
  "data": <具体数据，成功时存在>
}
```

成功或失败均使用此包装；常见状态码与处理：

- 200: 成功并返回 `data`
- 201: 资源已创建（如创建翻译单元）
- 204: 无内容（操作成功但不返回 body，例如 PATCH 成功）
- 4xx/5xx: 错误，`msg` 包含错误描述

---

**认证**：除 `/api/v1/login` 外的路由需要在 `Authorization: Bearer <token>` 中携带有效 JWT。

---

**Auth（登录）**

- **POST** `/api/v1/login`
  - 认证：否
  - 请求体：`model.LoginArgs`（JSON）
  - 返回：登录结果与 token（见 `data`）
  - 响应体（data）：`model.LoginReply`

**Users（用户）**

- **GET** `/api/v1/users`
  - 认证：是
  - 查询参数：对应 `model.RetrieveUserOpt`（通过 query 读取）
  - 返回：用户列表（`data`）
  - 响应体（data）：`[]model.UserInfo`
- **GET** `/api/v1/users/{user_id}`
  - 认证：是
  - 路径参数：`user_id`（字符串）
  - 返回：单个用户信息（`data`）
  - 响应体（data）：`model.UserInfo`
- **POST** `/api/v1/users/invite`
  - 认证：是
  - 请求体：`model.InviteUserArgs`（JSON）
  - 返回：邀请操作结果（`data`）
  - 响应体（data）：`model.InviteUserReply`
- **PATCH** `/api/v1/users/{user_id}`

  - 认证：是
  - 路径参数：`user_id`
  - 请求体：`model.UpdateUserArgs`（JSON），要求 `user_id` 在路径与 body 中一致
  - 返回：204 No Content（成功）
  - 响应体：无（204）

- **PATCH** `/api/v1/users/{user_id}/role`
  - 认证：是
  - 路径参数：`user_id`
  - 请求体：`model.AssignUserRoleArgs`（JSON）
  - 要求：路径参数 `user_id` 与请求体内 `user_id` 字段一致
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Worksets（工作集）**

- **GET** `/api/v1/worksets`
  - 认证：是
  - 查询参数：`limit`（int, default=10）、`offset`（int, default=0）
  - 返回：工作集列表（`data`）
  - 响应体（data）：`[]model.WorksetInfo`
- **GET** `/api/v1/worksets/{workset_id}`
  - 认证：是
  - 路径参数：`workset_id`
  - 返回：工作集详情（`data`）
  - 响应体（data）：`model.WorksetInfo`
- **POST** `/api/v1/worksets`
  - 认证：是
  - 请求体：`model.CreateWorksetArgs`（JSON）
  - 返回：创建结果（`data`）
  - 响应体（data）：`model.CreateWorksetReply`
- **PATCH** `/api/v1/worksets/{workset_id}`
  - 认证：是
  - 路径参数：`workset_id`
  - 请求体：`model.UpdateWorksetArgs`（JSON），实现时会把 `args.ID` 设为路径参数
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Comics（漫画）**

- **GET** `/api/v1/comics`
  - 认证：是
  - 查询参数：对应 `model.RetrieveComicOpt`
  - 返回：漫画简要列表（`data`）
  - 响应体（data）：`[]model.ComicBrief`
- **GET** `/api/v1/comics/{comic_id}`
  - 认证：是
  - 路径参数：`comic_id`
  - 返回：漫画详情（`data`）
  - 响应体（data）：`model.ComicInfo`
- **POST** `/api/v1/comics`
  - 认证：是
  - 请求体：`model.CreateComicArgs`（JSON）
  - 返回：创建结果（`data`）
  - 响应体（data）：`model.CreateComicReply`
- **PATCH** `/api/v1/comics/{comic_id}`
  - 认证：是
  - 路径参数：`comic_id`
  - 请求体：`model.UpdateComicArgs`（JSON），`args.ID` 会被设为路径参数
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Workset -> Comics**

- **GET** `/api/v1/worksets/{workset_id}/comics`
  - 认证：是
  - 路径参数：`workset_id`
  - 查询参数：`limit`/`offset`（默认见实现）
  - 返回：该工作集下的漫画简要列表（`data`）
  - 响应体（data）：`[]model.ComicBrief`

**Pages（页面）**

- **GET** `/api/v1/pages/{page_id}`
  - 认证：是
  - 路径参数：`page_id`
  - 返回：页面详情（`data`）
  - 响应体（data）：`model.ComicPageInfo`
- **POST** `/api/v1/pages`
  - 认证：是
  - 请求体：数组 `[]model.CreateComicPageArgs`（JSON），至少一项
  - 返回：创建结果（`data`），成功返回 200 或 201（实现中为 accept）
  - 响应体（data）：`[]model.CreateComicPageReply`
- **PATCH** `/api/v1/pages/{page_id}`
  - 认证：是
  - 路径参数：`page_id`
  - 请求体：`model.PatchComicPageArgs`（JSON），`args.ID` 会被设为路径参数
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Comic -> Pages**

- **GET** `/api/v1/comics/{comic_id}/pages`
  - 认证：是
  - 路径参数：`comic_id`
  - 返回：该漫画下的页面列表（`data`）
  - 响应体（data）：`[]model.ComicPageInfo`

**Units（翻译单元）**

- **GET** `/api/v1/pages/{page_id}/units`
  - 认证：是
  - 路径参数：`page_id`
  - 返回：该页面的翻译单元（`data`）
  - 响应体（data）：`[]model.ComicUnitInfo`
- **POST** `/api/v1/pages/{page_id}/units`
  - 认证：是
  - 路径参数：`page_id`
  - 请求体：数组 `[]model.NewComicUnitArgs`（JSON），每项 `PageID` 必须与路径一致，非空列表
  - 返回：201 Created（实现使用 `StatusCreated`）
  - 响应体：无（201）
- **PATCH** `/api/v1/pages/{page_id}/units`
  - 认证：是
  - 请求体：数组 `[]model.PatchComicUnitArgs`（JSON），非空列表
  - 返回：204 No Content（成功）
  - 响应体：无（204）
- **DELETE** `/api/v1/pages/{page_id}/units`
  - 认证：是
  - 请求体：数组 `[]string`（unit ID 列表），非空列表
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Assignments（分配/任务）**

- **GET** `/api/v1/assignments/{asgn_id}`
  - 认证：是
  - 路径参数：`asgn_id`
  - 返回：任务详情（`data`）
  - 响应体（data）：`model/po.BasicComicAsgn`
- **POST** `/api/v1/assignments`
  - 认证：是
  - 请求体：`model.CreateComicAsgnArgs`（JSON）
  - 返回：创建结果（`data`）
  - 响应体（data）：`string`（新分配的 ID）
- **PATCH** `/api/v1/assignments/{asgn_id}`
  - 认证：是
  - 路径参数：`asgn_id`
  - 请求体：`model/po.PatchComicAsgn`（JSON），实现会将 `ID` 设为路径参数
  - 返回：204 No Content（成功）
  - 响应体：无（204）

**Comic -> Assignments / User -> Assignments**

- **GET** `/api/v1/comics/{comic_id}/assignments`
  - 认证：是
  - 路径参数：`comic_id`
  - 查询参数：`limit`/`offset`
  - 返回：该漫画的任务列表（`data`）
  - 响应体（data）：`[]model/po.BasicComicAsgn`
- **GET** `/api/v1/users/{user_id}/assignments`
  - 认证：是
  - 路径参数：`user_id`
  - 查询参数：`limit`/`offset`
  - 返回：该用户的任务列表（`data`）
  - 响应體（data）：`[]model/po.BasicComicAsgn`

---

**DTO（请求/响应 数据模型）**

下面列出接口中常用的请求/响应结构（来自 `internal/model`）：

- `UserInfo` (响应)

  - `user_id` (string)
  - `qq` (string)
  - `nickname` (string)
  - `assigned_translator_at` (int64)
  - `assigned_proofreader_at` (int64)
  - `assigned_typesetter_at` (int64)
  - `assigned_redrawer_at` (int64)
  - `assigned_reviewer_at` (int64)
  - `assigned_uploader_at` (int64)
  - `is_admin` (bool)
  - `created_at` (int64)

- `LoginArgs` (请求)

  - `qq` (string)
  - `password` (string)
  - `nickname` (string, optional)
  - `invitation_code` (string, optional)

- `LoginReply` (响应)

  - `token` (string)

- `UpdateUserArgs` (请求)

  - `user_id` (string)
  - `qq` (\*string, optional)
  - `nickname` (\*string, optional)
  - `is_admin` (\*bool, optional)
  - role assignment flags: `assign_translator`, `assign_proofreader`, `assign_typesetter`, `assign_redrawer`, `assign_reviewer`, `assign_uploader` (all \*bool, optional)

- `InviteUserArgs` (请求)
  - `invitee_id` (string)
- `InviteUserReply` (响应)

  - `invitation_code` (string)

- `WorksetInfo` (响应)

  - `id` (string)
  - `index` (int64)
  - `name` (string)
  - `comic_count` (int64)
  - `description` (\*string)
  - `creator_id` (string)
  - `creator_nickname` (string)
  - `created_at` (int64)
  - `updated_at` (int64)

- `CreateWorksetArgs` (请求)
  - `name` (string)
  - `description` (\*string, optional)
- `CreateWorksetReply` (响应)

  - `id` (string)

- `ComicBrief` / `ComicInfo` (响应)

  - `id`, `workset_id`, `workset_index`, `index` (id/index fields)
  - `author` (string), `title` (string)
  - 状态时间戳字段（`translating_started_at`, `translating_completed_at`, `proofreading_started_at`, 等，皆为 \*int64）
  - `creator_id`, `creator_nickname` (仅 `ComicInfo`)
  - `description`, `comment` (\*string, optional, `ComicInfo`)
  - `created_at`, `updated_at` (int64, `ComicInfo`)

- `RetrieveComicOpt` (查询参数)

  - 若干模糊匹配、状态过滤和分页字段（`author`, `title`, `workset_id`, `index`, `offset`, `limit` 等）

- `CreateComicArgs` (请求)
  - `workset_id` (string), `author` (string), `title` (string)
  - `description`, `comment` (\*string, optional)
  - `pre_asgns` ([]PreAsgnArgs, optional)
- `CreateComicReply` (响应)

  - `id` (string)

- `ComicPageInfo` (响应)
  - `id` (string), `comic_id` (string), `index` (int64), `oss_url` (string), `uploaded` (bool)
- `CreateComicPageArgs` (请求)
  - `comic_id` (string), `index` (int64), `image_ext` (string)
- `CreateComicPageReply` (响应)
  - `id` (string), `oss_url` (string)
- `PatchComicPageArgs` (请求)

  - `id` (string), `uploaded` (\*bool, optional)

- `ComicUnitInfo` (响应)
  - `id` (string), `page_id` (string), `index` (int64)
  - 坐标: `x_coordinate`, `y_coordinate` (float64)
  - `is_in_box` (bool)
  - 翻译/校对字段: `translated_text` (*string), `translator_id` (*string), `translator_comment` (\*string)
  - 校对结果: `proved_text` (*string), `proved` (bool), `proofreader_id` (*string), `proofreader_comment` (\*string)
  - `created_at`, `updated_at` (int64)
- `NewComicUnitArgs` (请求，批量创建)
  - `page_id`, `index`, `x_coordinate`, `y_coordinate`, `is_in_box`
  - 可选文本/注释字段（`translated_text`, `translator_comment`, `proved_text`, `proved`, `proofreader_comment`）
- `PatchComicUnitArgs` (请求，批量更新)

  - `id` (string)，以及可选的 `index`, `x_coordinate`, `y_coordinate`, `is_in_box`, `translated_text`, `translator_comment`, `proved_text`, `proved`, `proofreader_comment`

- `CreateComicAsgnArgs` / `PreAsgnArgs` (请求)
  - `comic_id` (仅 `CreateComicAsgnArgs`), `assignee_id` (string)
  - 角色标志：`is_translator`, `is_proofreader`, `is_typesetter`, `is_redrawer`, `is_reviewer` (all \*bool, optional)
