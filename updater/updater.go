package updater

import (
	"fmt"
	"os"
)

// UpdateConfig 保存更新配置
type UpdateConfig struct {
	AsarSrc     string
	AsarDst     string
	UnpackedSrc string // 可选
	UnpackedDst string // 可选
	AppPath     string
}

// Update 执行更新流程
func Update(cfg UpdateConfig) error {
	// 1. 验证源文件
	if !FileExists(cfg.AsarSrc) {
		return fmt.Errorf("源 asar 文件未找到: %s", cfg.AsarSrc)
	}
	if cfg.UnpackedSrc != "" && !DirExists(cfg.UnpackedSrc) {
		return fmt.Errorf("源 unpacked 目录未找到: %s", cfg.UnpackedSrc)
	}

	// 2. 等待应用退出
	fmt.Println("正在等待应用退出...")
	if err := WaitForExit(cfg.AppPath); err != nil {
		return fmt.Errorf("应用未退出: %v", err)
	}

	// 3. 备份
	fmt.Println("正在备份现有文件...")
	backupAsar := cfg.AsarDst + ".bak"
	backupUnpacked := cfg.UnpackedDst + ".bak"

	// 辅助函数：如果更新成功则清理备份
	defer func() {
		// 成功后删除备份以节省空间
		os.Remove(backupAsar)
		if cfg.UnpackedDst != "" {
			os.RemoveAll(backupUnpacked)
		}
	}()

	// 备份 Asar
	if FileExists(cfg.AsarDst) {
		if err := os.Rename(cfg.AsarDst, backupAsar); err != nil {
			return fmt.Errorf("备份 asar 失败: %v", err)
		}
	}

	// 备份 Unpacked
	if cfg.UnpackedDst != "" && DirExists(cfg.UnpackedDst) {
		if err := os.Rename(cfg.UnpackedDst, backupUnpacked); err != nil {
			// 如果备份 unpacked 失败，尝试恢复 asar
			os.Rename(backupAsar, cfg.AsarDst)
			return fmt.Errorf("备份 unpacked 资源失败: %v", err)
		}
	}

	// 4. 执行更新 (替换文件)
	fmt.Println("正在应用更新...")
	err := performReplacement(cfg)
	if err != nil {
		fmt.Printf("更新失败: %v。正在回滚...\n", err)
		rollbackErr := performRollback(cfg, backupAsar, backupUnpacked)
		if rollbackErr != nil {
			return fmt.Errorf("更新失败: %v; 回滚也失败了: %v", err, rollbackErr)
		}
		return fmt.Errorf("更新失败并已回滚: %v", err)
	}

	// 5. 重启应用
	fmt.Println("正在重启应用...")
	if err := StartApp(cfg.AppPath); err != nil {
		return fmt.Errorf("重启应用失败: %v", err)
	}

	return nil
}

func performReplacement(cfg UpdateConfig) error {
	// 移动或复制 Asar
	// 优先尝试重命名，如果跨卷则回退到复制
	err := moveOrCopy(cfg.AsarSrc, cfg.AsarDst)
	if err != nil {
		return err
	}

	if cfg.UnpackedSrc != "" && cfg.UnpackedDst != "" {
		err := moveOrCopyUnpacked(cfg.UnpackedSrc, cfg.UnpackedDst)
		if err != nil {
			return err
		}
	}

	return nil
}

func moveOrCopy(src, dst string) error {
	// 优先尝试重命名
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	// 如果重命名失败，则复制
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	// os.Remove(src) // 可选：删除源文件
	return nil
}

func moveOrCopyUnpacked(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	if err := CopyDir(src, dst); err != nil {
		return err
	}
	// os.RemoveAll(src)
	return nil
}

func performRollback(cfg UpdateConfig, backupAsar, backupUnpacked string) error {
	var errs []error

	// 恢复 Asar
	if FileExists(backupAsar) {
		// 如果目标位置存在部分文件，先清理
		os.Remove(cfg.AsarDst)
		if err := os.Rename(backupAsar, cfg.AsarDst); err != nil {
			errs = append(errs, fmt.Errorf("恢复 asar: %v", err))
		}
	}

	// 恢复 Unpacked
	if cfg.UnpackedDst != "" && DirExists(backupUnpacked) {
		os.RemoveAll(cfg.UnpackedDst)
		if err := os.Rename(backupUnpacked, cfg.UnpackedDst); err != nil {
			errs = append(errs, fmt.Errorf("恢复 unpacked: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}
