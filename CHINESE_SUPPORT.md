# Chinese Font Support for Termshot

This branch adds full Chinese (CJK) character support to termshot by embedding Noto Sans CJK fonts.

## 新特性 / New Features

- ✅ **完整的中文字符支持** - 支持简体中文、繁体中文、日文和韩文字符
- ✅ **自动字体切换** - 根据字符类型自动在英文字体和中文字体之间切换
- ✅ **嵌入式字体** - 字体直接打包在二进制文件中，无需额外安装
- ✅ **保持原有功能** - 完全兼容原有的所有 termshot 功能

## 实现细节 / Implementation Details

### 字体选择 / Font Selection

使用 **Noto Sans CJK SC (思源黑体)** 作为中文字体：
- Regular: NotoSansCJK-Regular.ttc
- Bold: NotoSansCJK-Bold.ttc

### 技术架构 / Technical Architecture

1. **字体嵌入模块** (`internal/font/cjk.go`)
   - 使用 Go embed 将字体文件嵌入到二进制中
   - 支持 TrueType Collection (TTC) 格式解析
   - 提供与现有字体接口兼容的 API

2. **智能字体切换** (`internal/img/output.go`)
   - 新增 `isCJKChar()` 函数判断字符类型
   - 在渲染时根据字符自动选择合适的字体
   - 支持中英文混合显示

3. **CJK 字符检测范围**:
   - Unicode Han (CJK统一汉字)
   - Hiragana (日文平假名)
   - Katakana (日文片假名)
   - Hangul (韩文)
   - CJK Symbols and Punctuation (0x3000-0x303F)
   - Halfwidth and Fullwidth Forms (0xFF00-0xFFEF)

## 使用示例 / Usage Examples

### 基本使用 / Basic Usage

```bash
# 从文件读取包含中文的内容
./termshot --raw-read input_with_chinese.txt -f output.png

# 执行包含中文输出的命令
./termshot echo "你好，世界！"

# 捕获中文命令输出
./termshot -c -- bash -c "echo '测试中文'"
```

### 测试文件 / Test Files

项目包含以下测试文件：
- `test_input.txt` - 包含中文的示例输入
- `test_chinese.sh` - 中文测试脚本
- `test_chinese_output.png` - 测试生成的 PNG 示例

### 运行测试 / Run Tests

```bash
# 生成中文测试截图
./termshot --raw-read test_input.txt -f test_chinese_output.png

# 查看生成的文件
file test_chinese_output.png
```

## 构建 / Build

```bash
# 下载依赖
go mod tidy

# 编译
go build -v ./cmd/termshot

# 编译后的二进制文件包含所有嵌入的字体（约 40MB）
```

## 字体文件 / Font Files

字体文件位于 `internal/font/assets/`:
- `NotoSansCJK-Regular.ttc` (~19MB)
- `NotoSansCJK-Bold.ttc` (~19MB)

这些字体在编译时通过 `//go:embed` 指令嵌入到二进制文件中。

## 许可证 / License

- Termshot: MIT License (原项目许可证)
- Noto Sans CJK: SIL Open Font License 1.1

## 致谢 / Acknowledgments

- 原项目: [homeport/termshot](https://github.com/homeport/termshot)
- 字体: [Noto CJK Fonts](https://github.com/notofonts/noto-cjk) by Google
