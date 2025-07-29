#!/bin/bash

# 上传插件单元测试运行脚本

echo "=========================================="
echo "开始运行上传插件单元测试"
echo "=========================================="

# 设置测试环境变量
export GO_ENV=test
export TEST_MODE=true

# 创建测试目录
mkdir -p /tmp/upload_test

# 运行本地存储测试
echo "运行本地存储测试..."
cd service/local/test
go test -v -cover ./...
if [ $? -ne 0 ]; then
    echo "❌ 本地存储测试失败"
    exit 1
else
    echo "✅ 本地存储测试通过"
fi

# 运行阿里云存储测试
echo "运行阿里云存储测试..."
cd ../../aliyun/test
go test -v -cover ./...
if [ $? -ne 0 ]; then
    echo "❌ 阿里云存储测试失败"
    exit 1
else
    echo "✅ 阿里云存储测试通过"
fi

# 运行腾讯云存储测试
echo "运行腾讯云存储测试..."
cd ../../tencent/test
go test -v -cover ./...
if [ $? -ne 0 ]; then
    echo "❌ 腾讯云存储测试失败"
    exit 1
else
    echo "✅ 腾讯云存储测试通过"
fi

# 运行七牛云存储测试
echo "运行七牛云存储测试..."
cd ../../qiniuyun/test
go test -v -cover ./...
if [ $? -ne 0 ]; then
    echo "❌ 七牛云存储测试失败"
    exit 1
else
    echo "✅ 七牛云存储测试通过"
fi

# 清理测试目录
echo "清理测试文件..."
rm -rf /tmp/upload_test

echo "=========================================="
echo "所有测试完成！"
echo "==========================================" 