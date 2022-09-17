:: Deletes dist folder
set folder=".\dist"
cd /d %folder%
for /F "delims=" %%i in ('dir /b') do (rmdir "%%i" /s/q || del "%%i" /s/q)
cd ..

:: Builds the executable game
go build --ldflags "-s"

:: Moving all needed files into the dist folder
copy .\gopherLand.exe .\dist\gopherLand.exe
copy .\README.md .\dist\README.md
xcopy .\data .\dist\data /E /I