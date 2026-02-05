# LabelPlus 格式编码器实现文档

本文档提供从零开始实现 LabelPlus 格式编码器所需的完整信息，适用于从数据库数据生成 LabelPlus 格式翻校文本文件。

## 1. 数据结构前提

### 1.1 核心数据模型

LabelPlus 编码器需要以下三层数据结构：

#### Project（项目）

```go
type PortProject struct {
    Author string      // 作者名
    Title  string      // 作品标题
    Pages  []PortPage  // 页面列表
}
```

#### Page（页面）

```go
type PortPage struct {
    ImageFilename string      // 图片文件名（含扩展名）
    Units         []PortUnit  // 翻校单元列表
}
```

#### Unit（翻校单元）

```go
type PortUnit struct {
    X               float64  // 坐标 X
    Y               float64  // 坐标 Y
    IndexInPage     uint32   // 在当前页面中的序号（从1开始）
    IsInbox         bool     // 是否为框内文字
    TranslatedText  *string  // 翻译文本（可选）
    ProovedText     *string  // 校对文本（可选）
    IsProoved       bool     // 是否已校对
    Comment         *string  // 注释（可选）
}
```

### 1.2 字段说明

- **坐标 (X, Y)**: 翻校单元在图片中的坐标位置
- **IndexInPage**: 单元在页面中的排序序号，用于导出时的排序
- **IsInbox**: `true` 表示框内文字（组ID=1），`false` 表示框外文字（组ID=2）
- **TranslatedText**: 翻译者的翻译内容
- **ProovedText**: 校对者的校对内容
- **Comment**: 翻译或校对注释

## 2. LabelPlus 文件头规范

LabelPlus 格式文件由**文件头**和**正文体**两部分组成。

### 2.1 文件头结构

文件头包含4个固定部分，严格按以下顺序：

```
1,0
-
框内
框外
-
Exported by PopRaKo Web
<空行>
```

#### 详细说明

1. **第1行**: 版本号 `1,0`（固定不变）
2. **第2行**: 分隔符 `-`
3. **第3-4行**: 组名列表
   - 第3行: `框内` (组ID=1，对应 `is_inbox=true`)
   - 第4行: `框外` (组ID=2，对应 `is_inbox=false`)
4. **第5行**: 分隔符 `-`
5. **第6行**: 用户备注/导出工具标识（可自定义，如 `Exported by PopRaKo Web`）
6. **第7行**: 空行

### 2.2 Go 实现示例

```go
func writeHeader(w io.Writer) error {
    _, err := fmt.Fprintf(w, "1,0\n-\n框内\n框外\n-\nExported by PopRaKo Web\n\n")
    return err
}
```

## 3. 正文体导出算法

### 3.1 整体流程

```
for each page in project.pages:
    1. 写入页面头
    2. 按 index_in_page 排序该页面的所有 units
    3. for each unit in sorted_units:
        a. 写入单元头
        b. 写入主文本内容
        c. 写入注释（如果有）
        d. 写入空行
```

### 3.2 步骤详解

#### 步骤1: 页面头格式

每个页面以以下格式开始：

```
<两个空行>
>>>>>>>>[图片文件名]<<<<<<<<
```

**格式规则**:

- 前导两个空行（`\n\n`）
- 8个大于号 `>`
- 方括号包裹的图片文件名
- 8个小于号 `<`
- 结尾换行

**Go 实现**:

```go
fmt.Fprintf(w, "\n\n>>>>>>>>[%s]<<<<<<<<\n", page.ImageFilename)
```

#### 步骤2: 单元排序

**关键要求**: 必须按 `index_in_page` 升序排序单元后再输出。

```go
sort.Slice(page.Units, func(i, j int) bool {
    return page.Units[i].IndexInPage < page.Units[j].IndexInPage
})
```

#### 步骤3: 单元头格式

每个单元以以下格式开始：

```
----------------[N]----------------[X,Y,G]
```

**格式规则**:

- 16个短横线 `-`
- 方括号包裹的单元序号 `N`（从1开始，按排序后的顺序递增）
- 16个短横线 `-`
- 方括号包裹的坐标和组信息 `[X,Y,G]`
- 结尾换行

**字段说明**:

- `N`: 单元在当前页面中的序号（1-based，即第1个为1，第2个为2...）
- `X`: 坐标X，保留4位小数
- `Y`: 坐标Y，保留4位小数
- `G`: 组ID，`1`=框内，`2`=框外

**坐标格式化要求**:

- 必须保留恰好4位小数
- 避免科学计数法（如 `1.2345e+02` 是错误的）
- 示例: `123.4567` → `123.4567`，`0.5` → `0.5000`

**Go 实现**:

```go
n := i + 1  // 序号从1开始
g := 2      // 默认框外
if unit.IsInbox {
    g = 1
}

fmt.Fprintf(w, "----------------[%d]----------------[%.4f,%.4f,%d]\n",
    n, unit.X, unit.Y, g)
```

#### 步骤4: 主文本内容选择

主文本内容遵循**优先级规则**：

```
if prooved_text 非空:
    使用 prooved_text
else if translated_text 非空:
    使用 translated_text
else:
    不输出任何文本
```

**Go 实现**:

```go
var mainText string

if unit.ProovedText != nil && *unit.ProovedText != "" {
    mainText = *unit.ProovedText
} else if unit.TranslatedText != nil && *unit.TranslatedText != "" {
    mainText = *unit.TranslatedText
}

if mainText != "" {
    fmt.Fprintf(w, "%s\n", mainText)
}
```

#### 步骤5: 注释处理

如果 `comment` 字段非空，则在主文本后追加注释。

**格式规则**:

```
<空行>
#[翻校注释]：<注释内容>
```

**Go 实现**:

```go
if unit.Comment != nil && *unit.Comment != "" {
    fmt.Fprintf(w, "\n#[翻校注释]：%s\n", *unit.Comment)
}
```

#### 步骤6: 单元结尾

每个单元输出完成后，必须添加一个空行作为分隔。

```go
fmt.Fprintf(w, "\n")
```

## 4. 完整算法伪代码

```
function encodeLabelPlus(project: PortProject, output: Writer):
    // 写入文件头
    write("1,0\n")
    write("-\n")
    write("框内\n")
    write("框外\n")
    write("-\n")
    write("Exported by PopRaKo Web\n")
    write("\n")

    // 遍历所有页面
    for page in project.pages:
        // 页面头
        write("\n\n>>>>>>>>[" + page.image_filename + "]<<<<<<<<\n")

        // 排序单元
        sorted_units = sort(page.units, by: index_in_page)

        // 遍历单元
        for i, unit in enumerate(sorted_units):
            n = i + 1
            g = 1 if unit.is_inbox else 2

            // 单元头
            write("----------------[" + n + "]----------------[")
            write(format(unit.x, ".4f") + "," + format(unit.y, ".4f") + "," + g + "]\n")

            // 主文本
            main_text = ""
            if unit.prooved_text is not empty:
                main_text = unit.prooved_text
            else if unit.translated_text is not empty:
                main_text = unit.translated_text

            if main_text is not empty:
                write(main_text + "\n")

            // 注释
            if unit.comment is not empty:
                write("\n#[翻校注释]：" + unit.comment + "\n")

            // 单元结尾空行
            write("\n")
```

## 5. Go 完整实现示例

```go
package labelplus

import (
    "fmt"
    "io"
    "sort"
)

func EncodeLabelPlus(w io.Writer, project *PortProject) error {
    // 写入文件头
    if _, err := fmt.Fprintf(w, "1,0\n-\n框内\n框外\n-\nExported by PopRaKo Web\n\n"); err != nil {
        return err
    }

    // 遍历所有页面
    for _, page := range project.Pages {
        // 写入页面头
        if _, err := fmt.Fprintf(w, "\n\n>>>>>>>>[%s]<<<<<<<<\n", page.ImageFilename); err != nil {
            return err
        }

        // 排序单元（按 IndexInPage）
        sortedUnits := make([]PortUnit, len(page.Units))
        copy(sortedUnits, page.Units)
        sort.Slice(sortedUnits, func(i, j int) bool {
            return sortedUnits[i].IndexInPage < sortedUnits[j].IndexInPage
        })

        // 遍历并输出每个单元
        for i, unit := range sortedUnits {
            n := i + 1
            g := 2
            if unit.IsInbox {
                g = 1
            }

            // 写入单元头
            if _, err := fmt.Fprintf(w, "----------------[%d]----------------[%.4f,%.4f,%d]\n",
                n, unit.X, unit.Y, g); err != nil {
                return err
            }

            // 选择主文本内容
            var mainText string
            if unit.ProovedText != nil && *unit.ProovedText != "" {
                mainText = *unit.ProovedText
            } else if unit.TranslatedText != nil && *unit.TranslatedText != "" {
                mainText = *unit.TranslatedText
            }

            // 写入主文本
            if mainText != "" {
                if _, err := fmt.Fprintf(w, "%s\n", mainText); err != nil {
                    return err
                }
            }

            // 写入注释（如果有）
            if unit.Comment != nil && *unit.Comment != "" {
                if _, err := fmt.Fprintf(w, "\n#[翻校注释]：%s\n", *unit.Comment); err != nil {
                    return err
                }
            }

            // 单元结尾空行
            if _, err := fmt.Fprintf(w, "\n"); err != nil {
                return err
            }
        }
    }

    return nil
}
```

## 6. 格式严格性要求总结

| 元素         | 格式要求        | 错误示例          | 正确示例                                   |
| ------------ | --------------- | ----------------- | ------------------------------------------ |
| 版本号       | `1,0`           | `1.0`             | `1,0`                                      |
| 分隔符       | `-`             | `--`              | `-`                                        |
| 页面头大于号 | 恰好8个         | `>>>>>[...]<<<<<` | `>>>>>>>>[...]<<<<<<<<`                    |
| 页面头小于号 | 恰好8个         | 同上              | 同上                                       |
| 单元头短横线 | 恰好16个        | `-[1]-[...]`      | `----------------[1]----------------[...]` |
| 坐标小数位   | 恰好4位         | `1.2`             | `1.2000`                                   |
| 坐标格式     | 普通小数        | `1.2e+02`         | `120.0000`                                 |
| 组ID         | `1` 或 `2`      | `0`, `3`          | `1`, `2`                                   |
| 注释标记     | `#[翻校注释]：` | `#注释:`          | `#[翻校注释]：`                            |

## 7. 示例输出

假设有一个包含1个页面、2个单元的项目：

```
1,0
-
框内
框外
-
Exported by PopRaKo Web


>>>>>>>>[page_001.jpg]<<<<<<<<
----------------[1]----------------[123.4500,67.8900,1]
这是第一个框内文字

----------------[2]----------------[456.7800,234.5600,2]
这是第二个框外文字

#[翻校注释]：这里需要注意语气

```

## 8. 常见问题与注意事项

### 8.1 文本内容包含换行符怎么处理？

保持原样输出。LabelPlus 格式支持多行文本，解析器通过单元头分隔符来区分不同单元。

### 8.2 如果 `translated_text` 和 `prooved_text` 都为空？

单元头仍然输出，但主文本部分直接跳过（只输出单元头后的空行）。

### 8.3 `IndexInPage` 可以不连续吗？

可以。算法只关心相对大小顺序，不要求连续。例如 `[1, 3, 5]` 是合法的。

### 8.4 是否需要转义特殊字符？

不需要。LabelPlus 格式不需要转义，所有文本内容原样输出。

### 8.5 文件编码是什么？

推荐使用 **UTF-8 without BOM**。

### 8.6 Windows 和 Unix 换行符？

建议使用 `\n` (LF)，但 `\r\n` (CRLF) 也兼容。保持一致性即可。

## 9. 测试检查清单

实现后应验证以下要点：

- [ ] 文件头格式完全正确（版本、分隔符、组名、空行）
- [ ] 页面头的大于号/小于号数量为8个
- [ ] 单元头的短横线数量为16个
- [ ] 坐标格式化为4位小数且无科学计数法
- [ ] 单元按 `IndexInPage` 正确排序
- [ ] 组ID映射正确（框内=1，框外=2）
- [ ] 主文本优先级正确（prooved_text > translated_text）
- [ ] 注释格式正确且仅在非空时输出
- [ ] 每个单元结尾有空行
- [ ] 文件使用 UTF-8 编码

---

**文档版本**: 1.0  
**最后更新**: 2026-02-03  
**来源**: PopRaKo Native App (Rust 实现)
