@echo off
setlocal enabledelayedexpansion

:: 构建配置区（按需修改）
set APP_NAME=Caesar_GUI_Tools
set ICON_FILE=icon.png
set MAIN_FILE=main.go
set OUTPUT_DIR=dist
set RESOURCE_FILE=bundled.go

:: 初始化环境
if exist "%OUTPUT_DIR%" (
    echo [清理] 删除旧构建目录...
    rmdir /s /q "%OUTPUT_DIR%"
    if errorlevel 1 (
        echo [错误] 清理失败，请关闭占用 %OUTPUT_DIR% 的文件
        exit /b 1
    )
)
mkdir "%OUTPUT_DIR%" >nul 2>&1

:: 验证必要文件存在
if not exist "%ICON_FILE%" (
    echo [错误] 图标文件 %ICON_FILE% 不存在！
    exit /b 1
)
if not exist "%MAIN_FILE%" (
    echo [错误] 主程序文件 %MAIN_FILE% 不存在！
    exit /b 1
)

:: 资源打包阶段（网页2/网页5）
echo [阶段1] 打包图标资源...
fyne bundle -o "%RESOURCE_FILE%" "%ICON_FILE%"
if errorlevel 1 (
    echo [错误] 资源打包失败，请检查图标格式（需512x512 PNG）
    exit /b 1
)

:: 编译阶段（网页9）
echo [阶段2] 编译可执行文件...
go build -o "%OUTPUT_DIR%\%APP_NAME%.exe" "%MAIN_FILE%" "%RESOURCE_FILE%"
if errorlevel 1 (
    echo [错误] 编译失败，请检查Go代码
    exit /b 1
)

:: 应用打包阶段（网页2/网页5/网页9）
echo [阶段3] 生成分发包...
cd "%OUTPUT_DIR%"
fyne package -os windows -icon "..\%ICON_FILE%" -name "%APP_NAME%"
if errorlevel 1 (
    echo [错误] 打包失败，建议更新fyne工具：go install fyne.io/tools/cmd/fyne@latest
    exit /b 1
)
cd ..

:: 结果输出
echo [完成] 构建成功！输出文件：
dir /b "%OUTPUT_DIR%\*.exe"
echo.
echo 最终程序包位置：%CD%\%OUTPUT_DIR%\
pause