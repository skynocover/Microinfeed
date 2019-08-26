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
			//show := valuearr1[i] * magn
			show1 = show1 + strconv.FormatFloat(valuearr1[i]*magnum, 'f', num, 64) + "   " + "\\n "
		}
		show1 = show1 + "maxerror" + "\\n "
		show1 = show1 + strconv.FormatFloat(arrerr(valuearr1, 0.001)*magnum, 'f', num, 64) + "\\n "
		//ui.Eval(`document.querySelector('#done1').value= '` + show1 + `'`)
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

//<input id="mag" type="text" value="1000"  SIZE=5  height="35" style="font-size:20px;">
//			<div><textarea id="done4" rows="25" cols="7" style="position: absolute; "></textarea></div>
//<div><textarea id="done3" rows="25" cols="7" style="position: absolute; "></textarea></div>

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

// 單行讀取
/*
func readf(filename string) []string {
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Println("檔名或路徑錯誤")
		//return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	var str []string
	i := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		str[i] = string(a)
		i++
	}
	return str
}
*/

/*
func getfilename() string {
	fmt.Println("Please enter filename: ")
	var filename string
	fmt.Scanln(&filename)
	return filename
}
*/

/*
func dataarr2(data []string) []float64 {
	//宣告容器
	var valuearr []float64
	value := -999.9
	i, j := 0, 0
	//將每一千點抓出來若頭尾相差太大不採用
	for {
		if i+2999 > len(data) {
			break
		}

		newdata := data[i : i+2999]
		result := thoundpoints(newdata)
		if j < 25 {
			if result != 999999 {
				if math.Abs(result-value) > 0.0005 {
					valuearr = append(valuearr, result)
					value = valuearr[j]
					j++
				}
			}

		}

		i = i + 3000
	}
	//若與前一個數值相差太小不採用
	return valuearr
}

//陣列頭尾相差太大不採用
func thoundpoints(data []string) float64 {
	var newdata [50]float64
	for i := 0; i < len(newdata); i++ {
		newdata[i], _ = strconv.ParseFloat(data[i], 64)
	}

	data0, _ := strconv.ParseFloat(data[0], 64)
	dataend, _ := strconv.ParseFloat(data[len(data)-1], 64)
	dif := math.Abs(data0 - dataend)
	if dif > 0.0005 {
		return 999999
	} else {
		return arrayCountValues(newdata)[0]
		//return data0
	}

}
func arrayCountValues(args [50]float64) (MaxValue []float64) {
	/*【2】求出每个值对应出现的次数，例:[值:次数,值:次数]
	newMap := make(map[float64]float64)
	for _, value := range args {
		if newMap[value] != 0 {
			newMap[value]++
		} else {
			newMap[value] = 1
		}
	}

	/*【3】求出出现最多的次数
	var allCount []float64 //所有的次数
	var maxCount float64   //出现最多的次数
	for _, value := range newMap {
		allCount = append(allCount, value)
	}
	maxCount = allCount[0]
	for i := 0; i < len(allCount); i++ {
		if maxCount < allCount[i] {
			maxCount = allCount[i]
		}
	}

	/*【4】求数组中出现次数最多的值，例：[8,9]这个两个值出现的次数一样多
	var maxValue []float64
	for key, value := range newMap {
		if value == maxCount {
			maxValue = append(maxValue, key)
		}
	}
	return maxValue
}

func dataarr3(data []string) []float64 {
	var numdata []float64
	var classin []float64
	var class [][]float64
	var result []float64
	var i = 0
	for i := 0; i < len(data); i++ {
		num, _ := strconv.ParseFloat(data[i], 64)
		numdata = append(numdata, num)
	}

	for i+1000 < len(numdata) {
		if math.Abs(numdata[i]-numdata[i+999]) > 0.0003 {
			//fmt.Println(i)
			//fmt.Println(numdata[i])
			//fmt.Println(numdata[i+399])
			//fmt.Println(math.Abs(numdata[i] - numdata[i+399]))
			class = append(class, classin)
			i = i + 1000
		} else {
			classin = append(classin, numdata[i])
			i++
		}
	}
	//fmt.Println(class[1])
	fmt.Println(len(class))

	for j := 0; j < len(class); j++ {
		sum := 0.0
		for l := 0; l < len(class[j]); l++ {
			sum = sum + class[j][l]
		}
		result = append(result, sum/float64(len(class[j])))
	}

	//fmt.Println(class)

	return result

}
*/
