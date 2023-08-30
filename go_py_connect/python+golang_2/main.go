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
	fmt.Printf("%q", string(output))
	fmt.Println()
	fmt.Println()

	// 將輸出拆分成二維矩陣
	matrixStr := string(output)
	matrixStr = strings.TrimSpace(matrixStr) // 清除字符串两端的空格
	rows := strings.Split(matrixStr, "\n")   // 使用换行符分隔行

	// 創建一個二維整數切片來儲存整數
	var intMatrix [][]int
	for _, rowStr := range rows {
		values := strings.Fields(rowStr) // 使用空格分隔值
		var intRow []int
		for _, valueStr := range values {
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				fmt.Printf("無法將值轉換成整數: %v\n", err)
				continue
			}
			intRow = append(intRow, value)
		}
		intMatrix = append(intMatrix, intRow)
	}

	// 输出二维矩阵
	fmt.Println("二維矩陣數據:")
	for _, row := range intMatrix {
		fmt.Println(row)
	}
}
