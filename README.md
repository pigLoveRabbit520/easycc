# 使用
```
func main() {
	var concurrency uint = 40 // 并发量

	data := Person{
		Name:  "meng",
		Birth: "1993-09-21",
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println("json data is " + string(jsonData)) // body使用json数据

	req := &easycc.CCRequest{
		Method: "POST",
		URL:    BASE_URL+"/user/add/seat",
		Body:   jsonData,
		Headers: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}
	respArr := easycc.CCTest(req, concurrency)

	fmt.Println("all requests finished")

	for i := 0; i < len(respArr); i++ {
		resp := respArr[i]
		if resp.Err != nil {
			fmt.Println("this request encouter error " + resp.Err.Error())
		} else {
			fmt.Println("this request result is " + string(resp.Body))
		}
	}
}
```
