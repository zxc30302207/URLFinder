## URLFinder

URLFinder是一款快速、全面、易用的頁面資訊提取工具

此說明文件提供繁體中文（臺灣）翻譯，以提升在地使用體驗。

用於分析頁面中的js與url,查詢隱藏在其中的敏感資訊或未授權api介面

大致執行流程:

<img src="https://github.com/pingc0y/URLFinder/raw/master/img/process.png" width="85%"  />



有什麼需求或bug歡迎各位師傅提交issues

## 快速使用
單url
```
顯示全部狀態碼
URLFinder.exe -u http://www.baidu.com -s all -m 3

顯示200和403狀態碼
URLFinder.exe -u http://www.baidu.com -s 200,403 -m 3
```
批次url
```
結果分開儲存
匯出全部
URLFinder.exe -s all -m 3 -f url.txt -o .
只匯出html
URLFinder.exe -s all -m 3 -f url.txt -o res.html

結果統一儲存
URLFinder.exe -s all -m 3 -ff url.txt -o .
```
引數（更多引數使用 -i 配置）：
```
-a  自定義user-agent請求頭  
-b  自定義baseurl路徑  
-c  請求新增cookie  
-d  指定獲取的域名,支援正則表示式
-f  批次url抓取,需指定url文字路徑  
-ff 與-f區別：全部抓取的資料,視為同一個url的結果來處理（只打印一份結果 | 只會輸出一份結果） 
-h  幫助資訊   
-i  載入yaml配置檔案,可自定義請求頭、抓取規則等（不存在時,會在當前目錄建立一個預設yaml配置檔案）  
-m  抓取模式：
        1  正常抓取（預設）
        2  深入抓取 （URL深入一層 JS深入三層 防止抓偏）
        3  安全深入抓取（過濾delete,remove等敏感路由） 
-max 最大抓取數
-o  結果匯出到csv、json、html檔案,需指定匯出檔案目錄（.代表當前目錄）
-s  顯示指定狀態碼,all為顯示全部  
-t  設定執行緒數（預設50）
-time 設定超時時間（預設5,單位秒）
-u  目標URL  
-x  設定代理,格式: http://username:password@127.0.0.1:8877
-z  提取所有目錄對404連結進行fuzz(只對主域名下的連結生效,需要與 -s 一起使用）  
        1  目錄遞減fuzz  
        2  2級目錄組合fuzz
        3  3級目錄組合fuzz（適合少量連結使用）
```
## 使用截圖

[![0.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/0.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/0.jpg)   
[![1.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/1.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/1.jpg)  
[![2.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/2.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/2.jpg)  
[![3.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/3.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/3.jpg)  
[![4.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/4.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/4.jpg)  
[![5.jpg](https://github.com/pingc0y/URLFinder/raw/master/img/5.jpg)](https://github.com/pingc0y/URLFinder/raw/master/img/5.jpg)

## 部分說明

fuzz功能是基於抓到的404目錄和路徑。將其當作字典,隨機組合並碰撞出有效路徑,從而解決路徑拼接錯誤的問題

結果會優先顯示輸入的url頂級域名,其他域名不做區分顯示在 other

結果會優先顯示200,按從小到大排序（輸入的域名最優先,就算是404也會排序在其他子域名的200前面）

為了更好的相容和防止漏抓連結,放棄了低誤報率,錯誤的連結會變多但漏抓機率變低,可透過 ‘-s 200’ 篩選狀態碼過濾無效的連結（但不推薦只看200狀態碼）
##  編譯
以下是在windows環境下,編譯出各平臺可執行檔案的命令

```
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" -o ./URLFinder-windows-amd64.exe

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -ldflags "-s -w" -o ./URLFinder-windows-386.exe

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o ./URLFinder-linux-amd64

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
go build -ldflags "-s -w" -o ./URLFinder-linux-arm64

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -ldflags "-s -w" -o ./URLFinder-linux-386

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o ./URLFinder-macos-amd64

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=arm64
go build -ldflags "-s -w" -o ./URLFinder-macos-arm64
```


## 更新說明
2023/9/9  
修復 -ff 重複驗證問題  
修復 自動識別協議bug  

2023/9/2  
修復 子目錄定位bug   

2023/8/30  
修復 -i 配置請求頭錯誤問題   
變化 支援自動識別http/https  
變化 -o 輸入完整檔名可只匯出指定型別  
變化 無 -s 引數時，連結改為無顏色方便使用管道符

2023/5/11  
變化 -i 配置檔案可自定義：執行緒數、抓取深度、敏感路由、超時時間、最大抓取數   
新增 -time 設定超時時間  
新增 -max 設定最大抓取數  
新增 新增版本更新提示  
修復 已知bug

2023/5/5   
修復 多個任務時html結果混亂  
新增 結果新增302跳轉資訊  
變化 未獲取到資料時不列印與輸出結果

2023/4/22   
修復 已知bug  
變化 -d 改為正則表示式  
變化 列印顯示抓取來源  
新增 敏感資訊增加Other  
新增 -ff 全部抓取的資料,視為同一個url的結果來處理（只打印一份結果 | 只會輸出一份結果）

2023/2/21   
修復 已知bug

2023/2/3   
新增 域名資訊展示  
變化 -i配置檔案可配置抓取規則等

2023/1/29  
新增 -b 設定baseurl路徑  
新增 -o json、html格式匯出  
新增 部分敏感資訊獲取  
新增 預設會進行簡單的js爆破  
變化 能抓到更多連結,但垃圾資料變多  
變化 代理設定方式變更

2022/10/25  
新增 -t 設定執行緒數(預設50)  
新增 -z 對主域名的404連結fuzz測試  
最佳化 部分細節

2022/10/6  
新增 -x http代理設定  
修改 多個相同域名匯出時覆蓋問題處理

2022/9/23  
新增 對base標籤的相容  
修復 正則bug

2022/9/16  
新增 -m 3 安全的深入抓取,過濾delete、remove等危險URL   
新增 -d 獲取指定域名資源  
新增 -o 匯出到檔案顯示獲取來源source  
修復 已知bug

2022/9/15  
修復 某種情況下的陣列越界

2022/9/12  
修復 linux與mac下的配置檔案生成錯誤  
修復 已知邏輯bug

2022/9/5  
新增 連結存在標題時,顯示標題  
新增 -i 引數,載入yaml配置檔案（目前只支援配置請求頭headers）  
修改 部分程式碼邏輯  
修復 當ip存在埠時,匯出會去除埠

2022/8/29  
新增 抓取url數量顯示  
最佳化 部分程式碼  
新增 提供各平臺可執行檔案

2022/8/27   
新增 -o 改為自定義檔案目錄  
新增 匯出檔案改為csv字尾,表格檢視更方便  
修復 已知正則bug

2022/8/19  
最佳化 加長超時時間避免誤判

2022/8/5  
新增 狀態碼過濾  
新增 狀態碼驗證顯示進度  
修復 域名帶埠輸出本地錯誤問題

2022/7/25   
最佳化 js規則  
最佳化 排序  
新增 根據狀態碼顯示彩色字型

2022/7/6   
完善 規則

2022/6/27   
最佳化 規則  
新增 提供linux成品程式

2022/6/21   
修改 獲取狀態碼從自動改為手動（-s）  
新增 顯示響應內容大小

2022/6/16   
最佳化 提取規則增強相容性  
修復 陣列越界錯誤處理

2022/6/14  
修復 部分網站返回空值的問題

2022/6/13  
新增 自定義user-agent請求頭功能  
新增 批次url抓取功能  
新增 結果匯出功能  
最佳化 過濾規則  
最佳化 結果排版

2022/6/8  
修復 忽略ssl證書錯誤

# 開發由來
致敬JSFinder！開發此工具的初衷是因為經常使用 JSFinder 時會返回空或連結不完整,而且作者已經很久沒有更新修復 bug 了。因此,萌生了自己開發一款類似工具的想法。

