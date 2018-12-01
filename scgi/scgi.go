package scgi

import "fmt"

func EncodeRequest(data []byte) []byte {
	scgiHeaders := fmt.Sprintf("CONTENT_LENGTH\x00%d\x00", len(data))
	scgiRequest := fmt.Sprintf("%d:%s,%s", len(scgiHeaders), scgiHeaders, data)
	return []byte(scgiRequest)
}
