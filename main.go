package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zing/ztools-updater/updater"
)

func main() {
	asarSrc := flag.String("asar-src", "", "新 app.asar 文件路径")
	asarDst := flag.String("asar-dst", "", "目标 app.asar 文件路径")
	unpackedSrc := flag.String("unpacked-src", "", "新 app.asar.unpacked 目录路径 (可选)")
	unpackedDst := flag.String("unpacked-dst", "", "目标 app.asar.unpacked 目录路径 (可选)")
	appPath := flag.String("app", "", "Electron 主程序可执行文件路径")

	flag.Parse()

	// 参数校验
	if *asarSrc == "" || *asarDst == "" || *appPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg := updater.UpdateConfig{
		AsarSrc:     *asarSrc,
		AsarDst:     *asarDst,
		UnpackedSrc: *unpackedSrc,
		UnpackedDst: *unpackedDst,
		AppPath:     *appPath,
	}

	fmt.Println("正在启动 zTools Updater...")
	fmt.Printf("目标应用: %s\n", cfg.AppPath)
	fmt.Printf("替换: %s -> %s\n", cfg.AsarSrc, cfg.AsarDst)

	if err := updater.Update(cfg); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("更新成功完成！")
}
