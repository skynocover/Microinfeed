// main
package main

import (
	//"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/zserge/lorca"
)

func main() {
	ui, err := lorca.New("", "", 520, 560)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("download", func(filename string, magn string) {
		magnum, _ := strconv.ParseFloat(magn, 64)
		num := 5
		switch magnum {
		case 1000:
			num = 1
		case 100:
			num = 2
		case 10:
			num = 3
		case 1:
			num = 4
		}

		file := readf(filename)
		data := readata(file)
		name1 := getname(2, file)
		name2 := getname(1, file)
		arr1 := getarr(1, data)
		show1 := name1 + "     " + "\\n "
		valuearr1 := dataarr(arr1)
		for i := 0; i < len(valuearr1); i++ {
			show1 = show1 + strconv.FormatFloat(valuearr1[i]*magnum, 'f', num, 64) + "   " + "\\n "
		}
		show1 = show1 + "maxerror" + "\\n "
		show1 = show1 + strconv.FormatFloat(arrerr(valuearr1, 0.001)*magnum, 'f', num, 64) + "\\n "
		ui.Eval(`document.querySelector('#done1').innerText= '` + show1 + `'`) //若要用在div則是.value
		if name2 != "" {
			arr2 := getarr(2, data)
			show2 := name2 + "     " + "\\n "
			valuearr2 := dataarr(arr2)
			for i := 0; i < len(valuearr2); i++ {
				show2 = show2 + strconv.FormatFloat(valuearr2[i]*magnum, 'f', num, 64) + "   " + "\\n "
			}
			show2 = show2 + "maxerror" + "\\n "
			show2 = show2 + strconv.FormatFloat(arrerr(valuearr2, 0.001)*magnum, 'f', num, 64) + "\\n "
			ui.Eval(`document.querySelector('#done2').innerText= '` + show2 + `'`)
		}

	})
	// Load HTML after Go functions are bound to JS
	ui.Load("data:text/html," + url.PathEscape(`
		<html>
		<head>
		<title>Microinfeed</title>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		</head>
			<body>			
				<div class="field half">
				<label for="name" style="font-size:20px;">請輸入檔案名稱</label>
					<input id="URL" type="text" value="test.CSV"  SIZE=20  height="35" style="font-size:20px;">
					<input type="button" onclick="download(document.querySelector('#URL').value,document.querySelector('#mag').value)" style="width:70px;height:30px;font-size:16px;" value="讀取">
				</div>
				<div class="field half">
				<label id="name" style="font-size:20px;">倍率</label>				
					<select id="mag" type="text" value="1000"  SIZE=1  height="35" style="font-size:20px;">
						　<option value="1000">1000</option>
		　					  <option value="100">100</option>
		　					  <option value="10">10</option>
		　					  <option value="1">1</option>
					</select>
				</div>
				
				<div><br></div>
				<div class="field half">
				<label id="done1" for="name" style="display:block;font-size:16px;position: absolute;"></label>
				</div>
				<div class="field half">
				<label id="done2" for="name" style="display:block;font-size:16px;position: absolute;"></label>
				</div>			
			</body>
			<style type="text/css">
			#done1{			
				padding:10px;
				border:2px blue solid;
				margin-left:50px;
				float:left;
				width:70px;
				height:440px;
			}
			#done2{			
				padding:10px;
				border:2px green solid;
				margin-left:280px;
				float:left;
				width:70px;
				height:440px;
			}
			#name{
				margin-left:100px;
			}
			</style>
		</html>
	`))
	<-ui.Done()
}

//算出最大移動誤差
func arrerr(data []float64, target float64) float64 {
	var infeed []float64
	err := 0.0
	for i := 1; i < len(data)-1; i++ {
		infeed = append(infeed, math.Abs(data[i]-data[i-1]))
	}
	for i := 0; i < len(infeed); i++ {
		if math.Abs(infeed[i]-target) > err {
			err = math.Abs(infeed[i] - target)
		}
	}
	return err
}

//抓出階梯的陣列
func dataarr(data []string) []float64 {
	var numdata []float64 //newdara
	var result []float64

	var sw bool = false
	a, b := 0, 0

	var i = 0
	for i := 0; i < len(data); i++ {
		num, _ := strconv.ParseFloat(data[i], 64)
		numdata = append(numdata, num)
	}

	for i+400 < len(numdata) {
		if math.Abs(numdata[i]-numdata[i+399]) > 0.0003 { //不要的片段
			if sw != true { //如果上一段是要的片段則執行
				b = i
				arr := numdata[a:b]
				avg := arravg(arr)
				result = append(result, avg)
				i = i + 400
			}
			sw = false
			i = i + 400
		} else { //要的片段
			if sw == false {
				a = i
				sw = true
			}
			sw = true
			i++
		}
	}
	b = i
	arr := numdata[a:b]
	avg := arravg(arr)
	result = append(result, avg)

	return result

}
func arravg(arr []float64) float64 {
	sum := 0.0
	for i := 0; i < len(arr); i++ {
		sum = sum + arr[i]
	}
	avg := sum / float64(len(arr))
	return avg
}

//取得名稱
func getname(v int, data []string) string {
	var namearr []string
	var name string
	for i := 0; i < len(data); i++ {
		namearr = strings.Split(data[i], ",")
		if namearr[0] == "Waveform Name" {
			namearr[len(namearr)-v] = strings.Replace(namearr[len(namearr)-v], "\r", "", -1)
			name = strings.Replace(namearr[len(namearr)-v], "\"", "", -1)
			break
		}
	}
	return name
}

//取得V1或V2陣列
func getarr(v int, data []string) []string {
	var arr []string
	for i := 0; i < len(data); i++ {
		dataarr := strings.Split(data[i], ",")
		dataarr[v+1] = strings.Replace(dataarr[v+1], "\r", "", -1)
		value := dataarr[v+1]
		arr = append(arr, value)
	}
	return arr
}

//讀取文件
func readf(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

//將要的DATA從文件中取出
func readata(str []string) []string {
	var start, end int
	i := 0

	for {
		chat := strings.Split(str[i], ",")
		if chat[0] == "#EndHeader" {
			start = i + 1
		} else if chat[0] == "#BeginMark" {
			end = i - 1
			break
		}
		i++
	}
	newstr := str[start:end]
	return newstr
}
