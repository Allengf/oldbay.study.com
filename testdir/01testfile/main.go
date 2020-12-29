package main

import "fmt"

// UserInfo 用户信息
// type UserInfo struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// func main() {
// 	u1 := UserInfo{Name: "q1mi", Age: 18}

// 	b, _ := json.Marshal(&u1)
// 	fmt.Println(b)
// 	var m map[string]interface{}
// 	_ = json.Unmarshal(b, &m)
// 	for k, v := range m {
// 		fmt.Printf("key:%v value:%v\n", k, v)
// 	}
// }

func main() {
	ld := map[string]interface{}{
		"data":   string("测试"),
		"data01": string("测试01"),
	}
	fmt.Println(ld)
}
