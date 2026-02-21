# PopRaKo JSON 导出格式规范（服务端重建版）

本文档基于当前工作区实现整理，目标是让 Go 服务端可以**完整、无歧义**地解析 `.poprako.json`，并正确处理 `unit.is_local` 语义。

## 1. 适用范围

- 适用于文件名后缀为 `.poprako.json` 的导出文件。
- 当前客户端导出时，PopRaKo JSON 通过 `serde_json::to_writer_pretty` 写出（UTF-8，带缩进）。
- 本文不定义 LabelPlus 文本格式；LabelPlus 见 `docs/labelplus-export-format.md`。

## 2. 导出产物与命名

### 2.1 目录导出（Dir）

导出到目录时会生成：

1. `【{author}】{title}.poprako.json`
2. `【{author}】{title}.labelplus.txt`
3. 图片文件（由 `pages[].image_filename` 指向）

### 2.2 压缩导出（Zip）

导出到 zip 时会打包：

1. `【{author}】{title}.poprako.json`
2. `【{author}】{title}.labelplus.txt`
3. 每一页对应图片文件（取 `image_filename` 的文件名部分）

zip 名称为：`【{author}】{title}.zip`

## 3. JSON 顶层结构

```json
{
  "author": "string",
  "title": "string",
  "pages": [
    {
      "image_filename": "string",
      "units": [
        {
          "id": "string",
          "x": 0.0,
          "y": 0.0,
          "index_in_page": 1,
          "is_inbox": true,
          "translated_text": "string?",
          "prooved_text": "string?",
          "is_prooved": false,
          "comment": "string?",
          "is_local": false
        }
      ]
    }
  ]
}
```

## 4. 字段定义（严格）

## 4.1 Project

- `author: string`：作者名。
- `title: string`：标题。
- `pages: PortPage[]`：页面数组，顺序有意义。

## 4.2 Page

- `image_filename: string`
  - 导出时通常为图片文件名（例如 `001.jpg`）。
  - 服务端若只处理 JSON，不强依赖文件存在；若处理 zip/目录重建，应校验同名图片文件。
- `units: PortUnit[]`
  - 原始顺序保留；业务顺序应优先用 `index_in_page`。

## 4.3 Unit

- `id: string`
  - 单元业务 ID。
  - 可能是服务器下发 ID，也可能是本地生成 ID（需结合 `is_local` 判断来源）。
- `x: number(float64)`：横向相对坐标。
- `y: number(float64)`：纵向相对坐标。
- `index_in_page: uint32`：页内序号（建议 1..N）。
- `is_inbox: bool`：`true`=框内，`false`=框外。
- `translated_text?: string`
- `prooved_text?: string`
- `is_prooved: bool`
- `comment?: string`
- `is_local: bool`
  - **关键语义**：`true` 表示该 unit 是本地创建；`false` 表示该 unit 来自服务端。

## 5. 可省略字段规则（必须兼容）

客户端对以下字段使用了 `skip_serializing_if`：

- `translated_text`
- `prooved_text`
- `comment`

规则是：当字段为 `null` 或空字符串 `""` 时，导出 JSON 会**省略该 key**。

因此服务端反序列化必须支持：

1. key 不存在
2. key 存在但值为空字符串
3. key 存在且为普通字符串

并统一规范化为业务上的 “可空字符串” 表达。

## 6. `is_local` 与 `id` 的联合语义（重点）

这是服务端必须遵守的判断方式：

- 仅凭 `id` 不能判断来源。
- 必须使用联合键：`(id, is_local)`。

建议判定逻辑：

1. `is_local == false`
   - 视为“来自服务端的既有 unit”。
   - `id` 应匹配服务端已有记录 ID。
   - 上传时走更新流程（update/patch）。

2. `is_local == true`
   - 视为“客户端本地新建 unit”。
   - `id` 仅在客户端本地上下文有意义，服务端不应直接覆盖既有远端 ID 空间。
   - 上传时走新建流程（create），并返回服务端正式 ID（若有）。

## 7. 服务端重建建议（Go）

## 7.1 Go DTO（建议）

```go
type PortProject struct {
		Author string     `json:"author"`
		Title  string     `json:"title"`
		Pages  []PortPage `json:"pages"`
}

type PortPage struct {
		ImageFilename string     `json:"image_filename"`
		Units         []PortUnit `json:"units"`
}

type PortUnit struct {
		ID          string  `json:"id"`
		X           float64 `json:"x"`
		Y           float64 `json:"y"`
		IndexInPage uint32  `json:"index_in_page"`
		IsInbox     bool    `json:"is_inbox"`

		TranslatedText *string `json:"translated_text,omitempty"`
		ProovedText    *string `json:"prooved_text,omitempty"`
		IsProoved      bool    `json:"is_prooved"`
		Comment        *string `json:"comment,omitempty"`

		IsLocal bool `json:"is_local"`
}
```

说明：`*string` 可以天然兼容字段缺失。

## 7.2 反序列化后规范化

建议在入库前做一次 normalize：

1. `translated_text` / `prooved_text` / `comment`
   - 若指针非 nil 且 `trim == ""`，转为 nil（可选）。
2. `index_in_page`
   - 建议按 `index_in_page` 排序并校验唯一性。
3. `(id, is_local)`
   - 构建去重索引，防止重复单元。

## 7.3 校验规则（最低要求）

建议至少校验：

- `author` 非空
- `title` 非空
- `pages` 可空数组，但不得为 null
- 每个 page：`image_filename` 非空
- 每个 unit：
  - `id` 非空
  - `index_in_page >= 1`
  - `x`、`y` 为有限数值（非 NaN/Inf）

## 8. 兼容性约定

- 当前 JSON 无显式 `version` 字段。
- 服务端应采取“向前兼容”解析策略：
  - 忽略未知字段
  - 对可选文本字段容错缺失
- 若未来新增字段，不应破坏现有字段语义，尤其不能改变 `is_local` 的定义。

## 9. 与当前客户端行为对齐的注意事项

1. 导出时 `is_local` 直接从本地数据库原样写出。
2. 从 LabelPlus 导入生成的 unit，在当前实现中默认为 `is_local = false`。
3. 通过导入流程写入本地库时，当前实现会新建本地 unit 并写成 `is_local = true`。

第 2、3 条意味着：同一批文本在不同流程下，`is_local` 可能变化。服务端处理时应始终以“当前上送 JSON 中的 `id + is_local`”为准，不要推断历史来源。

## 10. 最小示例

```json
{
  "author": "Alice",
  "title": "Chapter 01",
  "pages": [
    {
      "image_filename": "001.jpg",
      "units": [
        {
          "id": "srv-unit-1001",
          "x": 0.4123,
          "y": 0.2874,
          "index_in_page": 1,
          "is_inbox": true,
          "translated_text": "你好",
          "is_prooved": false,
          "is_local": false
        },
        {
          "id": "local-temp-abc",
          "x": 0.7331,
          "y": 0.8452,
          "index_in_page": 2,
          "is_inbox": false,
          "is_prooved": false,
          "comment": "本地补充气泡",
          "is_local": true
        }
      ]
    }
  ]
}
```

上例中第二个 unit 没有 `translated_text` 与 `prooved_text`，这在当前格式下是合法的。
