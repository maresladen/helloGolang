for /f "delims=" %%i in ('node -v') do set str=%%i

if "%str%" == "" (
    echo 未安装nodejs
    goto insNode
) else (
    echo 已安装nodejs
    echo 检查nodejs版本
    goto cheNode
)


: insNode
echo 开始安装node程序,请不要关闭
msiexec /i node.msi /qb 
echo %PATH%
echo 此处打印全局变量,用于动态加载npm模块
pause
goto normjob
exit

:cheNode
for /f "delims=" %%i in ('node -v') do set str=%%i

echo %str%|find "^v6.9"> nul
if %errorlevel% equ 0 (
    echo 请卸载现有nodejs程序，并重新运行本批处理命令
) else (
    echo 匹配到6.9.*版本
    echo 跳过nodejs安装
    goto normjob
)
exit


:normjob
npm config set registry `配置路径`
echo 开始安装构建工具,请等待
npm install -g yo 

echo 开始安装rainbowUI脚手架工具,请等待
npm install -g generator-rainbowui-cli 


mkdir RainbowUI
cd RainbowUI
echo 构建项目
yo rainbowui-cli
echo 请进入RainbowUI目录启动项目，命令为npm run dev
pause
exit


REM set /p input=请输入字母A或B： 
REM if "%input%"=="A" goto A 
REM if "%input%"=="B" goto B 
REM pause 
REM exit 

REM :A 
REM echo 您输入的字母是A 
REM pause 
REM exit 

REM :B 
REM echo 您输入的字母是B 
REM pause 
REM exit 