## 检查模块

### 接口：检查更新

- **URL**: `/check-update`
- **请求方法**: `GET`
- **请求头**:
  - `X-Client-App-Version` (字符串): 客户端应用程序的版本。

#### 参数 DTO

- **CheckVersionReply**:
  - `latest_version` (字符串): 应用程序的最新版本。
  - `title` (字符串): 更新的标题。
  - `description` (字符串): 更新的描述。
  - `allow_usage` (布尔值): 指示客户端是否允许继续使用应用程序。

## 漫画模块

### 接口：根据ID获取漫画信息

- **URL**: `/comics/{comic_id}`
- **请求方法**: `GET`
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。

#### 响应 DTO

- **ComicInfo**:
  - `id` (字符串): 漫画的唯一标识符。
  - `workset_id` (字符串): 所属工作集的ID。
  - `workset_index` (整数): 工作集中的索引。
  - `index` (整数): 漫画的索引。
  - `creator_id` (字符串): 创建者的ID。
  - `creator_nickname` (字符串): 创建者的昵称。
  - `author` (字符串): 作者名称。
  - `title` (字符串): 漫画标题。
  - `description` (字符串，可选): 漫画描述。
  - `comment` (字符串，可选): 漫画评论。
  - `page_count` (整数): 页数。
  - `translating_started_at` (整数，可选): 翻译开始时间戳。
  - `translating_completed_at` (整数，可选): 翻译完成时间戳。
  - `proofreading_started_at` (整数，可选): 校对开始时间戳。
  - `proofreading_completed_at` (整数，可选): 校对完成时间戳。
  - `typesetting_started_at` (整数，可选): 排版开始时间戳。
  - `typesetting_completed_at` (整数，可选): 排版完成时间戳。
  - `reviewing_completed_at` (整数，可选): 审核完成时间戳。
  - `uploading_completed_at` (整数，可选): 上传完成时间戳。

### 接口：根据工作集ID获取漫画简要信息

- **URL**: `/worksets/{workset_id}/comics`
- **请求方法**: `GET`
- **路径参数**:
  - `workset_id` (字符串): 工作集的唯一标识符。
- **查询参数**:
  - `limit` (整数，默认值: 10): 返回的最大记录数。
  - `offset` (整数，默认值: 0): 返回记录的偏移量。

#### 响应 DTO

- **ComicBrief**:
  - `id` (字符串): 漫画的唯一标识符。
  - `workset_id` (字符串): 所属工作集的ID。
  - `workset_index` (整数): 工作集中的索引。
  - `index` (整数): 漫画的索引。
  - `author` (字符串): 作者名称。
  - `title` (字符串): 漫画标题。
  - `page_count` (整数): 页数。
  - `translating_started_at` (整数，可选): 翻译开始时间戳。
  - `translating_completed_at` (整数，可选): 翻译完成时间戳。
  - `proofreading_started_at` (整数，可选): 校对开始时间戳。
  - `proofreading_completed_at` (整数，可选): 校对完成时间戳。
  - `typesetting_started_at` (整数，可选): 排版开始时间戳。
  - `typesetting_completed_at` (整数，可选): 排版完成时间戳。
  - `reviewing_completed_at` (整数，可选): 审核完成时间戳。
  - `uploading_completed_at` (整数，可选): 上传完成时间戳。

---

### 接口：导出漫画 (LabelPlus)

- **URL**: `/api/v1/comics/{comic_id}/export`
- **请求方法**: `GET`
- **认证**: 该接口受认证保护（需通过 `AuthMiddleware`），调用导出需要有效的认证令牌。
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。

#### 响应 DTO

- **ExportComicReply**:
  - `export_uri` (字符串): 导出文件的可访问相对 URI，形如 `/comics/export/{filename}`。

---

### 接口：检索漫画简要信息

- **URL**: `/comics`
- **请求方法**: `GET`
- **查询参数**:
  - `aut` (字符串，可选): 作者名称（模糊查询）。
  - `tit` (字符串，可选): 漫画标题（模糊查询）。
  - `widx` (字符串，可选): 工作集索引。
  - `idx` (字符串，可选): 漫画索引。
  - `tsl_pending` (布尔值，可选): 是否未开始翻译。
  - `tsl_wip` (布尔值，可选): 是否正在翻译中。
  - `tsl_fin` (布尔值，可选): 是否已完成翻译。
  - `pr_pending` (布尔值，可选): 是否未开始校对。
  - `pr_wip` (布尔值，可选): 是否正在校对中。
  - `pr_fin` (布尔值，可选): 是否已完成校对。
  - `tst_pending` (布尔值，可选): 是否未开始排版。
  - `tst_wip` (布尔值，可选): 是否正在排版中。
  - `tst_fin` (布尔值，可选): 是否已完成排版。
  - `rv_pending` (布尔值，可选): 是否未开始审核。
  - `rv_fin` (布尔值，可选): 是否已完成审核。
  - `ul_pending` (布尔值，可选): 是否未开始上传。
  - `ul_fin` (布尔值，可选): 是否已完成上传。
  - `auid` (字符串，可选): 分配的用户ID。
  - `offset` (整数): 偏移量。
  - `limit` (整数): 返回的最大记录数。

#### 响应 DTO

- **ComicBrief**: 同上。

---

### 接口：创建漫画

- **URL**: `/comics`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **CreateComicArgs**:
    - `workset_id` (字符串): 工作集的唯一标识符。
    - `author` (字符串): 作者名称。
    - `title` (字符串): 漫画标题。
    - `description` (字符串，可选): 漫画描述。
    - `comment` (字符串，可选): 漫画评论。
    - `pre_asgns` (数组，可选): 预分配信息，包含以下字段：
      - `assignee_id` (字符串): 分配的用户ID。
      - `is_translator` (布尔值，可选): 是否为翻译者。
      - `is_proofreader` (布尔值，可选): 是否为校对者。
      - `is_typesetter` (布尔值，可选): 是否为排版者。
      - `is_redrawer` (布尔值，可选): 是否为修图者。
      - `is_reviewer` (布尔值，可选): 是否为审核者。

#### 响应 DTO

- **CreateComicReply**:
  - `id` (字符串): 创建的漫画的唯一标识符。

---

### 接口：根据ID更新漫画

- **URL**: `/comics/{comic_id}`
- **请求方法**: `PUT`
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。
- **请求体 DTO**:
  - **UpdateComicArgs**:
    - `id` (字符串): 漫画的唯一标识符。
    - `author` (字符串，可选): 作者名称。
    - `title` (字符串，可选): 漫画标题。
    - `description` (字符串，可选): 漫画描述。
    - `comment` (字符串，可选): 漫画评论。
    - `translating_started` (布尔值，可选): 是否开始翻译。
    - `translating_completed` (布尔值，可选): 是否完成翻译。
    - `proofreading_started` (布尔值，可选): 是否开始校对。
    - `proofreading_completed` (布尔值，可选): 是否完成校对。
    - `typesetting_started` (布尔值，可选): 是否开始排版。
    - `typesetting_completed` (布尔值，可选): 是否完成排版。
    - `reviewing_completed` (布尔值，可选): 是否完成审核。
    - `uploading_completed` (布尔值，可选): 是否完成上传。

---

### 接口：根据ID删除漫画

- **URL**: `/comics/{comic_id}`
- **请求方法**: `DELETE`
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。

---

## 漫画分配模块

### 接口：根据ID获取分配信息

- **URL**: `/assignments/{asgn_id}`
- **请求方法**: `GET`
- **路径参数**:
  - `asgn_id` (字符串): 分配的唯一标识符。

#### 响应 DTO

- **ComicAsgnInfo**:
  - `id` (字符串): 分配的唯一标识符。
  - `comic_id` (字符串): 漫画的唯一标识符。
  - `user_id` (字符串): 用户的唯一标识符。
  - `user_nickname` (字符串): 用户昵称。
  - `assigned_translator_at` (整数，可选): 分配翻译者的时间戳。
  - `assigned_proofreader_at` (整数，可选): 分配校对者的时间戳。
  - `assigned_typesetter_at` (整数，可选): 分配排版者的时间戳。
  - `assigned_redrawer_at` (整数，可选): 分配修图者的时间戳。
  - `assigned_reviewer_at` (整数，可选): 分配审核者的时间戳。
  - `created_at` (整数): 创建时间戳。
  - `updated_at` (整数): 更新时间戳。

---

### 接口：根据漫画ID获取分配信息

- **URL**: `/comics/{comic_id}/assignments`
- **请求方法**: `GET`
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。
- **查询参数**:
  - `limit` (整数，默认值: 10): 返回的最大记录数。
  - `offset` (整数，默认值: 0): 返回记录的偏移量。

#### 响应 DTO

- **ComicAsgnInfo**: 同上。

---

### 接口：根据用户ID获取分配信息

- **URL**: `/users/{user_id}/assignments`
- **请求方法**: `GET`
- **路径参数**:
  - `user_id` (字符串): 用户的唯一标识符。
- **查询参数**:
  - `limit` (整数，默认值: 10): 返回的最大记录数。
  - `offset` (整数，默认值: 0): 返回记录的偏移量。

#### 响应 DTO

- **ComicAsgnInfo**: 同上。

---

### 接口：创建分配

- **URL**: `/assignments`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **CreateComicAsgnArgs**:
    - `comic_id` (字符串): 漫画的唯一标识符。
    - `assignee_id` (字符串): 被分配的用户ID。
    - `is_translator` (布尔值，可选): 是否为翻译者。
    - `is_proofreader` (布尔值，可选): 是否为校对者。
    - `is_typesetter` (布尔值，可选): 是否为排版者。
    - `is_redrawer` (布尔值，可选): 是否为修图者。
    - `is_reviewer` (布尔值，可选): 是否为审核者。

#### 响应 DTO

- **ComicAsgnInfo**: 同上。

---

### 接口：根据ID更新分配

- **URL**: `/assignments/{asgn_id}`
- **请求方法**: `PUT`
- **路径参数**:
  - `asgn_id` (字符串): 分配的唯一标识符。
- **请求体 DTO**:
  - **UpdateComicAsgnArgs**:
    - `id` (字符串): 分配的唯一标识符。
    - `is_translator` (布尔值，可选): 是否为翻译者。
    - `is_proofreader` (布尔值，可选): 是否为校对者。
    - `is_typesetter` (布尔值，可选): 是否为排版者。
    - `is_redrawer` (布尔值，可选): 是否为修图者。
    - `is_reviewer` (布尔值，可选): 是否为审核者。

---

### 接口：根据ID删除分配

- **URL**: `/assignments/{asgn_id}`
- **请求方法**: `DELETE`
- **路径参数**:
  - `asgn_id` (字符串): 分配的唯一标识符。

---

## 漫画页面模块

### 接口：根据ID获取页面信息

- **URL**: `/pages/{page_id}`
- **请求方法**: `GET`
- **路径参数**:
  - `page_id` (字符串): 页面唯一标识符。

#### 响应 DTO

- **ComicPageInfo**:
  - `id` (字符串): 页面唯一标识符。
  - `comic_id` (字符串): 所属漫画的唯一标识符。
  - `index` (整数): 页面索引。
  - `oss_url` (字符串): 页面存储的OSS URL。
  - `uploaded` (布尔值): 页面是否已上传。
  - `inbox_unit_count` (整数): 收件箱单元数量。
  - `outbox_unit_count` (整数): 发件箱单元数量。
  - `translated_unit_count` (整数): 已翻译单元数量。
  - `proved_unit_count` (整数): 已校对单元数量。

---

### 接口：根据漫画ID获取页面信息

- **URL**: `/comics/{comic_id}/pages`
- **请求方法**: `GET`
- **路径参数**:
  - `comic_id` (字符串): 漫画的唯一标识符。

#### 响应 DTO

- **ComicPageInfo**: 同上。

---

### 接口：创建页面

- **URL**: `/pages`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **CreateComicPageArgs**:
    - `comic_id` (字符串): 所属漫画的唯一标识符。
    - `index` (整数): 页面索引。
    - `image_ext` (字符串): 图片扩展名。

#### 响应 DTO

- **CreateComicPageReply**:
  - `id` (字符串): 创建的页面唯一标识符。
  - `oss_url` (字符串): 页面存储的OSS URL。

---

### 接口：重新创建页面

- **URL**: `/pages/recreate`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **RecreateComicPageArgs**:
    - `id` (字符串): 页面唯一标识符。
    - `image_ext` (字符串): 图片扩展名。

#### 响应 DTO

- **CreateComicPageReply**:
  - `id` (字符串): 页面唯一标识符。
  - `oss_url` (字符串): 页面存储的OSS URL。

---

### 接口：根据ID更新页面

- **URL**: `/pages/{page_id}`
- **请求方法**: `PUT`
- **路径参数**:
  - `page_id` (字符串): 页面唯一标识符。
- **请求体 DTO**:
  - **PatchComicPageArgs**:
    - `id` (字符串): 页面唯一标识符。
    - `image_ext` (字符串，可选): 图片扩展名。
    - `uploaded` (布尔值，可选): 页面是否已上传。

---

### 接口：根据ID删除页面

- **URL**: `/pages/{page_id}`
- **请求方法**: `DELETE`
- **路径参数**:
  - `page_id` (字符串): 页面唯一标识符。

---

## 漫画翻译单元模块

### 接口：根据页面ID获取翻译单元

- **URL**: `/pages/{page_id}/units`
- **请求方法**: `GET`
- **路径参数**:
  - `page_id` (字符串): 页面唯一标识符。

#### 响应 DTO

- **ComicUnitInfo**:
  - `id` (字符串): 翻译单元的唯一标识符。
  - `page_id` (字符串): 所属页面的唯一标识符。
  - `index` (整数): 翻译单元的索引。
  - `x_coordinate` (浮点数): X 坐标。
  - `y_coordinate` (浮点数): Y 坐标。
  - `is_in_box` (布尔值): 是否在文本框内。
  - `translated_text` (字符串，可选): 翻译后的文本。
  - `translator_id` (字符串，可选): 翻译者的唯一标识符。
  - `translator_comment` (字符串，可选): 翻译者的评论。
  - `proved_text` (字符串，可选): 校对后的文本。
  - `proved` (布尔值): 是否已校对。
  - `proofreader_id` (字符串，可选): 校对者的唯一标识符。
  - `proofreader_comment` (字符串，可选): 校对者的评论。
  - `creator_id` (字符串，可选): 创建者的唯一标识符。
  - `created_at` (整数): 创建时间戳。
  - `updated_at` (整数): 更新时间戳。

---

### 接口：创建翻译单元

- **URL**: `/pages/{page_id}/units`
- **请求方法**: `POST`
- **路径参数**:
  - `page_id` (字符串): 页面唯一标识符。
- **请求体 DTO**:
  - **NewComicUnitArgs**:
    - `page_id` (字符串): 所属页面的唯一标识符。
    - `index` (整数): 翻译单元的索引。
    - `x_coordinate` (浮点数): X 坐标。
    - `y_coordinate` (浮点数): Y 坐标。
    - `is_in_box` (布尔值): 是否在文本框内。
    - `translated_text` (字符串，可选): 翻译后的文本。
    - `translator_comment` (字符串，可选): 翻译者的评论。
    - `proved_text` (字符串，可选): 校对后的文本。
    - `proved` (布尔值): 是否已校对。
    - `proofreader_comment` (字符串，可选): 校对者的评论。

---

### 接口：更新翻译单元

- **URL**: `/units`
- **请求方法**: `PUT`
- **请求体 DTO**:
  - **PatchComicUnitArgs**:
    - `id` (字符串): 翻译单元的唯一标识符。
    - `index` (整数，可选): 翻译单元的索引。
    - `x_coordinate` (浮点数，可选): X 坐标。
    - `y_coordinate` (浮点数，可选): Y 坐标。
    - `is_in_box` (布尔值，可选): 是否在文本框内。
    - `translated_text` (字符串，可选): 翻译后的文本。
    - `translator_comment` (字符串，可选): 翻译者的评论。
    - `proved_text` (字符串，可选): 校对后的文本。
    - `proved` (布尔值，可选): 是否已校对。
    - `proofreader_comment` (字符串，可选): 校对者的评论。

---

### 接口：删除翻译单元

- **URL**: `/units`
- **请求方法**: `DELETE`
- **请求体 DTO**:
  - `unit_ids` (数组): 要删除的翻译单元ID列表。

---

## 用户模块

### 接口：获取当前用户信息

- **URL**: `/users/me`
- **请求方法**: `GET`

#### 响应 DTO

- **UserInfo**:
  - `id` (字符串): 用户的唯一标识符。
  - `qq` (字符串): 用户的QQ号。
  - `nickname` (字符串): 用户昵称。
  - `assigned_translator_at` (整数，可选): 分配为翻译者的时间戳。
  - `assigned_proofreader_at` (整数，可选): 分配为校对者的时间戳。
  - `assigned_typesetter_at` (整数，可选): 分配为排版者的时间戳。
  - `assigned_redrawer_at` (整数，可选): 分配为修图者的时间戳。
  - `assigned_reviewer_at` (整数，可选): 分配为审核者的时间戳。
  - `assigned_uploader_at` (整数，可选): 分配为上传者的时间戳。
  - `is_admin` (布尔值): 是否为管理员。
  - `created_at` (整数): 用户创建时间戳。

---

### 接口：根据ID获取用户信息

- **URL**: `/users/{user_id}`
- **请求方法**: `GET`
- **路径参数**:
  - `user_id` (字符串): 用户的唯一标识符。

#### 响应 DTO

- **UserInfo**: 同上。

---

### 接口：更新用户信息

- **URL**: `/users/{user_id}`
- **请求方法**: `PUT`
- **路径参数**:
  - `user_id` (字符串): 用户的唯一标识符。
- **请求体 DTO**:
  - **UpdateUserArgs**:
    - `id` (字符串): 用户的唯一标识符。
    - `qq` (字符串，可选): 用户的QQ号。
    - `nickname` (字符串，可选): 用户昵称。

---

### 接口：邀请用户

- **URL**: `/invitations`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **CreateInvitationArgs**: 请参考邀请模块的文档。

---

### 接口：获取邀请信息

- **URL**: `/invitations`
- **请求方法**: `GET`

#### 响应 DTO

- **InvitationInfo**: 请参考邀请模块的文档。

---

### 接口：用户登录

- **URL**: `/login`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **LoginArgs**:
    - `qq` (字符串): 用户的QQ号。
    - `password` (字符串): 用户密码。
    - `nickname` (字符串，可选): 用户昵称。
    - `invitation_code` (字符串，可选): 邀请码。

#### 响应 DTO

- **LoginReply**:
  - `token` (字符串): 用户的认证令牌。

---

### 接口：检索用户信息

- **URL**: `/users`
- **请求方法**: `GET`
- **查询参数**:
  - `nn` (字符串，可选): 用户昵称（模糊查询）。
  - `qq` (字符串，可选): 用户的QQ号。
  - `ia` (布尔值，可选): 是否为管理员。
  - `itsl` (布尔值，可选): 是否为翻译者。
  - `ipr` (布尔值，可选): 是否为校对者。
  - `itst` (布尔值，可选): 是否为排版者。
  - `ird` (布尔值，可选): 是否为修图者。
  - `irv` (布尔值，可选): 是否为审核者。
  - `iul` (布尔值，可选): 是否为上传者。
  - `offset` (整数): 偏移量。
  - `limit` (整数): 返回的最大记录数。

#### 响应 DTO

- **UserInfo**: 同上。

---

### 接口：分配用户角色

- **URL**: `/users/{user_id}/roles`
- **请求方法**: `PUT`
- **路径参数**:
  - `user_id` (字符串): 用户的唯一标识符。
- **请求体 DTO**:
  - **AssignUserRoleArgs**:
    - `id` (字符串): 用户的唯一标识符。
    - `roles` (数组): 用户角色分配列表，包含以下字段：
      - `role` (字符串): 角色名称。
      - `assigned` (布尔值): 是否分配该角色。

---

## 工作集模块

### 接口：根据ID获取工作集信息

- **URL**: `/worksets/{workset_id}`
- **请求方法**: `GET`
- **路径参数**:
  - `workset_id` (字符串): 工作集的唯一标识符。

#### 响应 DTO

- **WorksetInfo**:
  - `id` (字符串): 工作集的唯一标识符。
  - `index` (整数): 工作集的索引。
  - `name` (字符串): 工作集名称。
  - `comic_count` (整数): 包含的漫画数量。
  - `description` (字符串，可选): 工作集描述。
  - `creator_id` (字符串): 创建者的唯一标识符。
  - `creator_nickname` (字符串): 创建者的昵称。
  - `created_at` (整数): 创建时间戳。
  - `updated_at` (整数): 更新时间戳。

---

### 接口：检索工作集

- **URL**: `/worksets`
- **请求方法**: `GET`
- **查询参数**:
  - `limit` (整数，默认值: 10): 返回的最大记录数。
  - `offset` (整数，默认值: 0): 返回记录的偏移量。

#### 响应 DTO

- **WorksetInfo**: 同上。

---

### 接口：创建工作集

- **URL**: `/worksets`
- **请求方法**: `POST`
- **请求体 DTO**:
  - **CreateWorksetArgs**:
    - `name` (字符串): 工作集名称。
    - `description` (字符串，可选): 工作集描述。

#### 响应 DTO

- **CreateWorksetReply**:
  - `id` (字符串): 创建的工作集唯一标识符。

---

### 接口：根据ID更新工作集

- **URL**: `/worksets/{workset_id}`
- **请求方法**: `PUT`
- **路径参数**:
  - `workset_id` (字符串): 工作集的唯一标识符。
- **请求体 DTO**:
  - **UpdateWorksetArgs**:
    - `id` (字符串): 工作集的唯一标识符。
    - `description` (字符串，可选): 工作集描述。

---

### 接口：根据ID删除工作集

- **URL**: `/worksets/{workset_id}`
- **请求方法**: `DELETE`
- **路径参数**:
  - `workset_id` (字符串): 工作集的唯一标识符。

---
