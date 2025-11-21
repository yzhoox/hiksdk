# 贡献指南

感谢你考虑为 HikSDK 做出贡献！

## 如何贡献

### 报告问题

如果你发现了 Bug 或有功能建议：

1. 在提交前先搜索现有的 [Issues](../../issues)
2. 如果没有找到相关的 Issue，创建一个新的
3. 清晰地描述问题，包括：
   - 问题的详细描述
   - 重现步骤
   - 期望的行为
   - 实际的行为
   - 环境信息（操作系统、Go 版本、设备型号等）

### 提交代码

1. **Fork 项目**
   ```bash
   # 点击右上角的 Fork 按钮
   ```

2. **克隆仓库**
   ```bash
   git clone https://github.com/你的用户名/hiksdk.git
   cd hiksdk
   ```

3. **创建分支**
   ```bash
   git checkout -b feature/你的功能名
   # 或
   git checkout -b fix/你的修复名
   ```

4. **进行修改**
   - 遵循项目的代码风格
   - 添加必要的测试
   - 更新相关文档

5. **运行测试**
   ```bash
   go test -v ./pkg/...
   go fmt ./...
   go vet ./...
   ```

6. **提交更改**
   ```bash
   git add .
   git commit -m "类型: 简短描述
   
   详细描述你的更改..."
   ```
   
   提交信息格式：
   - `feat:` 新功能
   - `fix:` Bug 修复
   - `docs:` 文档更新
   - `test:` 测试相关
   - `refactor:` 重构
   - `style:` 代码格式
   - `chore:` 构建/工具相关

7. **推送到 GitHub**
   ```bash
   git push origin feature/你的功能名
   ```

8. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 清晰描述你的更改
   - 关联相关的 Issue

## 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 使用 `go vet` 检查代码
- 变量和函数命名使用驼峰命名法
- 导出的函数和类型需要添加注释

### 测试要求

- 新功能必须包含单元测试
- 测试覆盖率应不低于 70%
- 测试应该是独立的，不依赖外部状态
- 使用环境变量进行集成测试配置

### 文档要求

- 所有公开的 API 都需要有注释
- 更新 README.md（如果适用）
- 更新 CHANGELOG.md
- 添加或更新示例代码（如果适用）

## 开发环境设置

### 前置条件

- Go 1.25 或更高版本
- Git
- CGO 编译器（Windows: MinGW-w64, Linux: gcc）

### 安装依赖

```bash
go mod download
```

### 运行测试

```bash
# 设置环境变量
# Linux/macOS
source test_env.sh
# Windows
.\test_env.ps1

# 运行所有测试
go test -v ./pkg/...

# 运行特定测试
go test -v ./pkg/ -run TestDeviceLogin

# 运行基准测试
go test -v ./pkg/ -bench=. -benchmem
```

### 代码检查

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 整理依赖
go mod tidy

# 一键运行所有检查
go fmt ./... && go vet ./... && go test ./pkg/...
```

## 项目结构

```
hiksdk/
├── pkg/            # SDK 核心代码
├── examples/       # 示例代码
├── include/        # C 头文件
├── lib/            # 动态链接库
├── docs/           # 文档（如果有）
├── README.md       # 项目说明
├── CHANGELOG.md    # 更新日志
└── CONTRIBUTING.md # 贡献指南（本文件）
```

## Pull Request 流程

1. 确保所有测试通过
2. 更新相关文档
3. 在 CHANGELOG.md 中记录更改
4. 提交 PR 后等待代码审查
5. 根据反馈进行修改
6. PR 被接受后会被合并

## 审查标准

代码审查时会关注：

- **功能性**：代码是否按预期工作
- **测试**：是否有充分的测试覆盖
- **代码质量**：代码是否清晰、可维护
- **性能**：是否有不必要的性能开销
- **安全性**：是否存在安全隐患
- **文档**：是否有适当的注释和文档

## 行为准则

### 我们的承诺

为了营造开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用性别化的语言或图像
- 人身攻击或侮辱性评论
- 公开或私下的骚扰
- 未经许可发布他人的私人信息
- 其他不道德或不专业的行为

## 需要帮助？

如果你有任何问题或需要帮助：

- 查看 [README.md](README.md)
- 搜索或创建 [Issue](../../issues)
- 查看现有的 [Pull Requests](../../pulls)

## 许可证

提交代码即表示你同意你的贡献将在 MIT 许可证下发布。

---

再次感谢你的贡献！🎉

