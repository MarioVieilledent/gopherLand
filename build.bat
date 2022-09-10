go build --ldflags "-s"
MOVE gopherLand.exe dist/gopherLand.exe
cd data
xcopy map.txt dist/data
xcopy background.png dist/data
xcopy ressources.png dist/data