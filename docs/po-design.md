# PO 设计指南

本文档说明如何将 SQL 表结构映射为项目中的 PO（Persistence Objects）结构体，采用与 `internal/model/po/user.go` 相同的风格与约定。

**目的**

- 统一数据库映射风格，便于 ORM（GORM）使用、序列化与单元测试。

**总体约定**

- 每张表对应一个或多个 PO 结构体（例如用于创建、返回、局部更新的变体）。
- 在每个 PO 文件中，使用常量保存表名（例如 `USER_TABLE`），并为需要的结构体实现 `TableName()` 方法返回该常量。
- 结构体字段使用 GORM 风格的标签（`gorm:"column_name;..."`），至少标注主键 `primaryKey` 与列名（如需要）。

**命名规则**

- 表名常量：大写下划线，例如 `USER_TABLE`。
- 创建用结构体：`New<Target>`（例如 `NewUser`）。
- 读取用结构体：`Basic<Target>`（例如 `BasicUser`）。
- 局部更新结构体：`Patch<Target>`，所有可更新字段均使用指针类型以支持“忽略零值”的语义（例如 `PatchUser`）。必选字段不用指针。

**字段类型与标签**

- 使用与数据库列直接对应的 Go 类型（例如 `string` 对应 `text`/`varchar`；时间戳使用 `time.Time`，由项目约定决定）。
- 在需要表达“该字段是主键”时，使用 `gorm:"...;primaryKey"` 标签。
- 不要在 PO 结构体中放置业务逻辑方法；PO 仅作数据承载。

**Patch 结构体**

- `Patch` 结构体中每个可选字段都应为指针类型（例如 `*string`、`*int64`），以便在应用更新时：
  - `nil` 表示“未提供该字段——不要更新数据库列”。
  - 非 `nil`（即使指向零值）表示明确更新为该值。
- 这样可在服务层统一判断哪些字段需要更新，避免误把零值当作要写入的值。

**外键与关联**

- 外键对应的结构体必须以指针方式嵌入（例如 `Author *BasicUser`），理由：
  - 指针可为 `nil`，表明该关联未被预加载或不存在。
  - 避免值拷贝导致的结构体复制，减小内存与循环引用风险。
- 建议同时保留原生外键列字段（例如 `AuthorID string`）与指针型关联字段（例如 `Author *BasicUser`）。
- 关联字段的 GORM tag 可根据需要加上 `foreignKey`/`references` 等指示，但基本约定是保留外键列并在需要时预加载指针结构体。

**TableName 方法**

- 对于需要指定表名的 PO 结构体，实现：

```go
func (*BasicUser) TableName() string { return USER_TABLE }
```

- 统一返回文件中声明的表名常量，避免硬编码字符串散落在代码中。

**示例（参考 `internal/model/po/user.go`）**

- `NewUser`（用于创建）:

```go
type NewUser struct {
    ID           string `gorm:"id;primaryKey"`
    Email        string `gorm:"email"`
    PasswordHash string `gorm:"password_hash"`
}
```

- `BasicUser`（用于返回/查询）:

```go
type BasicUser struct {
    ID        string `gorm:"id;primaryKey"`
    Email     string `gorm:"email"`
    Nickname  string `gorm:"nickname"`
    CreatedAt int64  `gorm:"created_at"`
}
```

- `PatchUser`（用于局部更新，注意指针字段）:

```go
type PatchUser struct {
    ID       string `gorm:"id;primaryKey"`
    Email    *string `gorm:"email"`
    Nickname *string `gorm:"nickname"`
}
```

- 外键示例（假设 `comic` 表有 `author_id` 外键指向 `user`）:

```go
type Comic struct {
    ID       string    `gorm:"id;primaryKey"`
    Title    string    `gorm:"title"`
    AuthorID string    `gorm:"author_id;foreignKey:AuthorID;references:ID"`
    Author   *BasicUser
}
```

> 说明：具体 `gorm` 标签（例如 `foreignKey`, `references`）可按项目中使用的 GORM 版本与查询习惯加入，上面示例重点展示结构与指针约定。

**检查清单（每次新增 PO 时）**

- 是否有表名常量与 `TableName()`？
- 创建/查询/更新是否使用了合适的结构体变体（`New*` / `Basic*` / `Patch*`）？
- `Patch` 结构体字段是否为指针？
- 外键是否同时保留原生外键列与指针型关联？
- 字段标签是否标注了必要的 `primaryKey` 或列名？
