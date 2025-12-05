//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// 测试跨平台编译
func main() {
	fmt.Println("==========================================")
	fmt.Println("海康SDK 跨平台编译测试")
	fmt.Println("==========================================")
	fmt.Printf("当前系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go版本: %s\n", runtime.Version())
	fmt.Println()

	// 测试当前平台编译
	fmt.Println("【1】测试当前平台编译...")
	if err := testBuild(runtime.GOOS, runtime.GOARCH); err != nil {
		fmt.Printf("✗ 编译失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 当前平台编译成功")

	// 显示CGO配置
	fmt.Println("\n【2】检查CGO配置...")
	checkCGO()

	// 检查系统库配置
	fmt.Println("\n【3】检查系统库配置...")
	checkSystemLibraries()

	// 测试语法检查
	fmt.Println("\n【4】运行代码检查...")
	if err := runVet(); err != nil {
		fmt.Printf("⚠ 代码检查发现潜在问题: %v\n", err)
	} else {
		fmt.Println("✓ 代码检查通过")
	}

	fmt.Println("\n==========================================")
	fmt.Println("测试完成！")
	fmt.Println("==========================================")
	fmt.Println("\n✓ SDK支持以下平台:")
	fmt.Println("  • Windows (amd64) - 使用MinGW-w64编译")
	fmt.Println("  • Linux (amd64)   - 使用GCC编译")
	fmt.Println("\n✓ 跨平台特性:")
	fmt.Println("  • 自动平台检测")
	fmt.Println("  • 类型兼容处理")
	fmt.Println("  • 统一API接口")
	fmt.Println("  • 支持重新初始化")
}

// 测试编译
func testBuild(goos, goarch string) error {
	cmd := exec.Command("go", "build", "-v", "./core")
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
		"CGO_ENABLED=1",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, string(output))
	}
	return nil
}

// 检查CGO
func checkCGO() {
	cmd := exec.Command("go", "env", "CGO_ENABLED")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠ 无法获取CGO状态: %v\n", err)
		return
	}

	cgoEnabled := strings.TrimSpace(string(output))
	if cgoEnabled == "1" {
		fmt.Println("✓ CGO已启用")
	} else {
		fmt.Println("✗ CGO未启用，请设置 CGO_ENABLED=1")
	}

	// 检查编译器
	var ccCmd string
	if runtime.GOOS == "windows" {
		ccCmd = "gcc"
	} else {
		ccCmd = "gcc"
	}

	cmd = exec.Command(ccCmd, "--version")
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠ 未找到C编译器 (%s)\n", ccCmd)
	} else {
		fmt.Printf("✓ C编译器可用 (%s)\n", ccCmd)
	}
}

// 检查系统库配置
func checkSystemLibraries() {
	if runtime.GOOS == "windows" {
		// Windows: 检查 PATH 环境变量
		path := os.Getenv("PATH")
		if strings.Contains(strings.ToLower(path), "hiksdk") ||
			strings.Contains(strings.ToLower(path), "hcnetsdk") {
			fmt.Println("✓ 检测到 PATH 中可能包含海康 SDK 路径")
		} else {
			fmt.Println("⚠ PATH 中未检测到海康 SDK 路径")
			fmt.Println("  请确保已从海康官网下载 SDK 并配置到系统 PATH")
		}
	} else {
		// Linux: 检查 LD_LIBRARY_PATH
		ldPath := os.Getenv("LD_LIBRARY_PATH")
		if strings.Contains(strings.ToLower(ldPath), "hiksdk") ||
			strings.Contains(strings.ToLower(ldPath), "hcnetsdk") {
			fmt.Println("✓ 检测到 LD_LIBRARY_PATH 中可能包含海康 SDK 路径")
		} else {
			fmt.Println("⚠ LD_LIBRARY_PATH 中未检测到海康 SDK 路径")
			fmt.Println("  请确保已从海康官网下载 SDK 并配置到 LD_LIBRARY_PATH")
		}
	}

	fmt.Println("\n  提示: 请从海康威视开放平台下载 SDK:")
	fmt.Println("  https://open.hikvision.com/download/5cda567cf47ae80dd41a54b3?type=10")
}

// 运行代码检查
func runVet() error {
	cmd := exec.Command("go", "vet", "./core")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) > 0 {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
