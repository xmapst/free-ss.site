package utils

import "strings"

func StringsContains(array []string, val string) (index bool) {
	index = false
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = true
			return
		}
	}
	return
}

func InSliceString(v string, sl []string) bool {
	for _, vv := range sl {
		if strings.ContainsAny(v, vv) {
			return true
		}
	}
	return false
}

func RemoveRepByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
