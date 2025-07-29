# 上传插件单元测试说明

## 概述

本文档描述了上传插件的单元测试结构和使用方法。测试覆盖了四个主要的存储服务：
- 本地存储 (Local Storage)
- 阿里云OSS (Aliyun OSS)
- 腾讯云COS (Tencent COS)
- 七牛云存储 (Qiniu Storage)

## 测试结构

```
service/
├── local/test/
│   └── service_test.go      # 本地存储测试
├── aliyun/test/
│   └── service_test.go      # 阿里云存储测试
├── tencent/test/
│   └── service_test.go      # 腾讯云存储测试
└── qiniuyun/test/
    └── service_test.go      # 七牛云存储测试
```

## 测试内容

### 1. 本地存储测试 (service/local/test/)

**测试用例：**
- `TestNewService` - 服务创建测试
- `TestSetUserID` - 用户ID设置测试
- `TestSetUserIDFromRequestClaims` - 从请求中提取用户ID测试
- `TestSaveFormFile_ValidImage` - 有效图片文件上传测试
- `TestSaveFormFile_InvalidFileType` - 无效文件类型测试
- `TestSaveEncodeFile_Base64Image` - Base64编码图片上传测试
- `TestDeleteFile` - 文件删除测试
- `TestValidateMagicNumber` - 魔数验证测试
- `TestGenerateUploadName` - 文件名生成测试
- `TestRelativePath` - 相对路径处理测试
- `TestUploadDisabled` - 上传禁用测试

**特点：**
- 使用真实的文件系统操作
- 包含魔数验证
- 路径遍历攻击防护测试
- 文件大小计算验证

### 2. 阿里云存储测试 (service/aliyun/test/)

**测试用例：**
- `TestNewService` - 服务创建测试
- `TestSetUserID` - 用户ID设置测试
- `TestSetUserIDFromRequestClaims` - 用户ID提取测试
- `TestSaveFormFile_UploadDisabled` - 上传禁用测试
- `TestSaveEncodeFile_UploadDisabled` - Base64上传禁用测试
- `TestGenerateObjectKey` - 对象键生成测试
- `TestEncodeDataToFile_Base64` - Base64数据处理测试
- `TestDeleteFile_BucketCRUDDisabled` - 删除功能禁用测试
- `TestIntegration_RealAliyunOSS` - 集成测试（需要真实配置）

**特点：**
- 模拟OSS客户端操作
- 配置验证测试
- 错误处理测试

### 3. 腾讯云存储测试 (service/tencent/test/)

**测试用例：**
- `TestNewService` - 服务创建测试
- `TestSetUserID` - 用户ID设置测试
- `TestSetUserIDFromRequestClaims` - 用户ID提取测试
- `TestSaveFormFile_UploadDisabled` - 上传禁用测试
- `TestSaveEncodeFile_UploadDisabled` - Base64上传禁用测试
- `TestGenerateObjectKey` - 对象键生成测试
- `TestEncodeDataToFile_Base64` - Base64数据处理测试
- `TestDeleteFile_BucketCRUDDisabled` - 删除功能禁用测试
- `TestIntegration_RealTencentCOS` - 集成测试（需要真实配置）

**特点：**
- 模拟COS客户端操作
- 配置验证测试
- 错误处理测试

### 4. 七牛云存储测试 (service/qiniuyun/test/)

**测试用例：**
- `TestNewService` - 服务创建测试
- `TestSetUserID` - 用户ID设置测试
- `TestSetUserIDFromRequestClaims` - 用户ID提取测试
- `TestSaveFormFile_UploadDisabled` - 上传禁用测试
- `TestSaveEncodeFile_UploadDisabled` - Base64上传禁用测试
- `TestGenerateObjectKey` - 对象键生成测试
- `TestEncodeDataToFile_Base64` - Base64数据处理测试
- `TestDeleteFile` - 文件删除测试
- `TestQiniuSpecificFeatures` - 七牛云特有功能测试
- `TestIntegration_RealQiniu` - 集成测试（需要真实配置）

**特点：**
- 模拟七牛云客户端操作
- 特有功能测试（如CDN域名）
- 配置验证测试

## 运行测试

### 1. 运行所有测试

```bash
./run_tests.sh
```

### 2. 运行特定存储测试

```bash
# 本地存储测试
cd service/local/test
go test -v -cover ./...

# 阿里云存储测试
cd service/aliyun/test
go test -v -cover ./...

# 腾讯云存储测试
cd service/tencent/test
go test -v -cover ./...

# 七牛云存储测试
cd service/qiniuyun/test
go test -v -cover ./...
```

### 3. 运行特定测试用例

```bash
# 运行特定测试函数
go test -v -run TestSaveFormFile_ValidImage

# 运行包含特定关键词的测试
go test -v -run "TestSaveFormFile"
```

## 测试配置

### 环境变量

```bash
export GO_ENV=test
export TEST_MODE=true
```

### 配置文件

测试使用 `test_config.yaml` 配置文件，包含：
- 各存储服务的测试配置
- 测试数据定义
- 测试用例配置
- 性能测试参数

## 测试数据

### 有效文件数据

- **JPEG图片**: `[0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01]`
- **PNG图片**: `[0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]`
- **PDF文档**: `[0x25, 0x50, 0x44, 0x46]`

### 无效文件数据

- **无效数据**: `[0x00, 0x01, 0x02, 0x03, 0x04, 0x05]`

## Mock对象

### MockMongoDB

模拟MongoDB连接，用于测试数据库操作：

```go
type MockMongoDB struct {
    *mongo.Client
    collections map[string]*mongo.Collection
}
```

### 测试辅助函数

- `createTestConfig()` - 创建测试配置
- `createTestService()` - 创建测试服务
- `cleanupTestFiles()` - 清理测试文件

## 注意事项

### 1. 云存储测试

云存储测试（阿里云、腾讯云、七牛云）需要真实的配置才能进行集成测试。在测试环境中：

- 使用模拟的客户端
- 跳过需要真实配置的测试
- 主要测试配置验证和错误处理

### 2. 文件系统测试

本地存储测试会创建真实的文件和目录：

- 测试目录：`/tmp/upload_test`
- 测试完成后自动清理
- 包含路径遍历攻击防护测试

### 3. 并发安全测试

- 使用线程安全的随机数生成器
- 测试并发文件上传场景
- 验证资源竞争问题

## 测试覆盖率

运行测试覆盖率报告：

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 持续集成

测试脚本可以集成到CI/CD流程中：

```yaml
# GitHub Actions 示例
- name: Run Tests
  run: |
    chmod +x run_tests.sh
    ./run_tests.sh
```

## 故障排除

### 常见问题

1. **权限错误**: 确保测试目录有写权限
2. **端口冲突**: 检查MongoDB测试端口是否被占用
3. **配置错误**: 验证测试配置文件格式

### 调试模式

启用详细日志：

```bash
go test -v -debug ./...
```

## 扩展测试

### 添加新测试用例

1. 在对应的测试文件中添加测试函数
2. 遵循命名规范：`Test[功能名]_[场景]`
3. 包含必要的断言和错误检查

### 性能测试

添加性能测试：

```go
func BenchmarkUpload(b *testing.B) {
    // 性能测试代码
}
```

运行性能测试：

```bash
go test -bench=. ./...
``` 