@echo off
setlocal enabledelayedexpansion

for /r %%F in (*.go1) do (
    rem 完整路径（不含扩展名）
    set "BASE=%%~dpnF"

    rem 如果 .go 存在，删除
    if exist "!BASE!.go" (
        echo 删除 "!BASE!.go"
        del /f /q "!BASE!.go"
    )

    rem 使用 move 重命名（最关键）
    echo 重命名 "%%F" -> "!BASE!.go"
    move /y "%%F" "!BASE!.go" >nul
)

echo 完成
pause
