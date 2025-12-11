# zTools Updater 使用指南

本工具用于为 zTools Electron 应用提供无签名的热更新功能。它可以替换 `app.asar` 和 optional 的 `app.asar.unpacked` 目录，并在更新完成后重启应用。

## 命令行参数

| 参数 | 必填 | 描述 |
| :--- | :--- | :--- |
| `--asar-src` | 是 | 新的 `app.asar` 文件路径 |
| `--asar-dst` | 是 | 目标 `app.asar` 文件路径 (将被覆盖) |
| `--unpacked-src` | 否 | 新的 `app.asar.unpacked` 目录路径 |
| `--unpacked-dst` | 否 | 目标 `app.asar.unpacked` 目录路径 (将被覆盖) |
| `--app` | 是 | Electron 主程序可执行文件路径 (用于等待退出和重启) |

## 调用流程

1.  **Electron 应用下载更新**: 应用在运行下载新的 `app.asar` (以及 `app.asar.unpacked`) 到临时目录。
2.  **调用 Updater**: 使用 `spawn` 或 `exec` 启动 `ztools-updater`，传入上述参数。
3.  **Electron 应用退出**: 启动 Updater 后，Electron 应用应立即调用 `app.quit()` 退出。
4.  **Updater 执行**:
    *   Updater 等待 Electron 进程完全退出。
    *   备份旧文件。
    *   将新文件移动/复制到目标位置。
    *   如果失败，自动回滚。
    *   重启 Electron 应用。

## 示例

### macOS

假设您的应用安装在 `/Applications/zTools.app`，更新文件下载到了 `/tmp/update`。

```bash
./dist/mac-arm64/ztools-updater \
  --asar-src "/tmp/update/app.asar" \
  --asar-dst "/Applications/zTools.app/Contents/Resources/app.asar" \
  --unpacked-src "/tmp/update/app.asar.unpacked" \
  --unpacked-dst "/Applications/zTools.app/Contents/Resources/app.asar.unpacked" \
  --app "/Applications/zTools.app/Contents/MacOS/zTools"
```

### Windows

假设您的应用安装在 `C:\Program Files\zTools`，更新文件下载到了 `C:\Temp\update`。

```powershell
.\dist\win-amd64\ztools-updater.exe ^
  --asar-src "C:\Temp\update\app.asar" ^
  --asar-dst "C:\Program Files\zTools\resources\app.asar" ^
  --unpacked-src "C:\Temp\update\app.asar.unpacked" ^
  --unpacked-dst "C:\Program Files\zTools\resources\app.asar.unpacked" ^
  --app "C:\Program Files\zTools\zTools.exe"
```

## Node.js (Electron) 集成示例

```javascript
const { spawn } = require('child_process');
const path = require('path');
const { app } = require('electron');

function runUpdater(updateDir) {
  const isMac = process.platform === 'darwin';
  const updaterName = isMac ? 'ztools-updater' : 'ztools-updater.exe';
  // 假设 updater 因为某种原因打包在应用内，或者随更新下载
  // 这里假设 updater 二进制文件路径
  const updaterPath = path.join(process.resourcesPath, 'bin', updaterName); 

  const asarSrc = path.join(updateDir, 'app.asar');
  const asarDst = path.join(process.resourcesPath, 'app.asar');
  
  // App executable path
  // macOS: /Applications/zTools.app/Contents/MacOS/zTools
  // Windows: C:\Program Files\zTools\zTools.exe
  const appPath = process.execPath; 

  const args = [
    '--asar-src', asarSrc,
    '--asar-dst', asarDst,
    '--app', appPath
  ];

  // 如果有 unpacked
  // args.push('--unpacked-src', ...);
  // args.push('--unpacked-dst', ...);

  const subprocess = spawn(updaterPath, args, {
    detached: true,
    stdio: 'ignore'
  });

  subprocess.unref();
  app.quit();
}
```

## 注意事项

*   **权限**: 在 Windows 上，如果应用安装在 `Program Files`，Updater 可能需要管理员权限才能写入文件。通常建议将 Updater 嵌入到主程序中，或确保主程序有写入权限，或者请求提权运行 Updater。
*   **备份**: Updater 会在同一目录下生成 `.bak` 备份文件。更新成功后会自动删除。
