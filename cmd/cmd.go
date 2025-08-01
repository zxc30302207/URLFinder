package cmd

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
)

var Update = "2023.9.9"
var XUpdate string

var (
	H  bool
	I  bool
	M  int
	S  string
	U  string
	D  string
	C  string
	A  string
	B  string
	F  string
	FF string
	O  string
	X  string
	T  = 50
	TI = 5
	MA = 99999
	Z  int
)

func init() {
	flag.StringVar(&A, "a", "", "set user-agent\n設定user-agent請求頭")
	flag.StringVar(&B, "b", "", "set baseurl\n設定baseurl路徑")
	flag.StringVar(&C, "c", "", "set cookie\n設定cookie")
	flag.StringVar(&D, "d", "", "set domainName\n指定獲取的域名,支援正則表示式")
	flag.StringVar(&F, "f", "", "set urlFile\n批次抓取url,指定檔案路徑")
	flag.StringVar(&FF, "ff", "", "set urlFile one\n與-f區別：全部抓取的資料,視為同一個url的結果來處理（只打印一份結果 | 只會輸出一份結果）")
	flag.BoolVar(&H, "h", false, "this help\n幫助資訊")
	flag.BoolVar(&I, "i", false, "set configFile\n載入yaml配置檔案（不存在時,會在當前目錄建立一個預設yaml配置檔案）")
	flag.IntVar(&M, "m", 1, "set mode\n抓取模式 \n   1 normal\n     正常抓取（預設） \n   2 thorough\n     深入抓取（預設url深入一層,js深入三層,-i可以自定義） \n   3 security\n     安全深入抓取（過濾delete,remove等敏感路由.-i可自定義） ")
	flag.IntVar(&MA, "max", 99999, "set maximum\n最大抓取連結數")
	flag.StringVar(&O, "o", "", "set outFile\n結果匯出到csv、json、html檔案,需指定匯出檔案目錄,可填寫完整檔名只匯出一種型別（.代表當前目錄）")
	flag.StringVar(&S, "s", "", "set Status\n顯示指定狀態碼,all為顯示全部（多個狀態碼用,隔開）")
	flag.IntVar(&T, "t", 50, "set Thread\n設定執行緒數（預設50）")
	flag.IntVar(&TI, "time", 5, "set Timeout\n設定超時時間（預設5,單位秒）")
	flag.StringVar(&U, "u", "", "set Url\n目標URL")
	flag.StringVar(&X, "x", "", "set Proxy\n設定代理,格式: http://username:password@127.0.0.1:8809")
	flag.IntVar(&Z, "z", 0, "set Fuzz\n對404連結進行fuzz(只對主域名下的連結生效,需要與 -s 一起使用） \n   1 decreasing\n     目錄遞減fuzz \n   2 2combination\n     2級目錄組合fuzz（適合少量連結使用） \n   3 3combination\n     3級目錄組合fuzz（適合少量連結使用） ")

	// 改變預設的 Usage
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `Usage: URLFinder [-a user-agent] [-b baseurl] [-c cookie] [-d domainName] [-f urlFile] [-ff urlFile one]  [-h help]  [-i configFile]  [-m mode] [-max maximum] [-o outFile]  [-s Status] [-t thread] [-time timeout] [-u url] [-x proxy] [-z fuzz]

Options:
`)
	flag.PrintDefaults()
}

func Parse() {
	color.LightCyan.Printf("         __   __   ___ _           _           \n /\\ /\\  /__\\ / /  / __(_)_ __   __| | ___ _ __ \n/ / \\ \\/ \\/// /  / _\\ | | '_ \\ / _` |/ _ \\ '__|\n\\ \\_/ / _  \\ /___ /   | | | | | (_| |  __/ |   \n \\___/\\/ \\_\\____\\/    |_|_| |_|\\__,_|\\___|_|     \n\nBy: pingc0y\nUpdate: %s | %s\nGithub: https://github.com/pingc0y/URLFinder \n\n", Update, XUpdate)
	flag.Parse()
}
