package main

import (
	"fmt"
	"os/exec"
	"strconv" //處理字串型別轉換
	"strings"
)

func main() {
	// 執行Python檔案並呼叫指定的函式
	var i, j string = "1", "5"
	cmd := exec.Command("python", "your_python_file.py", i, j)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("執行失敗：", err)
		return
	}

	fmt.Println("massage sent to golang:")
	fmt.Printf("%q", string(output)) //顯示含有轉義字符的字串形式
	fmt.Println()
	fmt.Println()

	// ---------------------------------------
	// 將輸出拆分
	matrixStr := string(output)
	matrixStr = strings.TrimSpace(matrixStr) // 清除字符串两端的空格
	matrixValues1 := strings.Split(matrixStr, " ")

	// 將字符串切片轉換為整數切片
	var ans1, ans2, ans3 int = 0, 0, 0
	for i, valueStr := range matrixValues1 {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			fmt.Printf("無法將字符轉換成整數: %v\n", err)
			continue
		}
		if i == 0 {
			ans1 = value
		}
		if i == 1 {
			ans2 = value
		}
		if i == 2 {
			ans3 = value
		}
	}

	fmt.Printf("ans1 = %v\n", ans1)
	fmt.Printf("ans2 = %v\n", ans2)
	fmt.Printf("ans3 = %v\n", ans3)
}
