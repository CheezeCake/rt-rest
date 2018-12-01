package rpc

import (
	"encoding/xml"
	"strings"
)

func EncodeRequest(method string, params []string) []byte {
	var builder strings.Builder
	builder.WriteString("<methodCall>")

	builder.WriteString("<methodName>")
	xml.EscapeText(&builder, []byte(method))
	builder.WriteString("</methodName>")

	builder.WriteString("<params>")
	buildParams(params, &builder)
	builder.WriteString("</params>")

	builder.WriteString("</methodCall>")
	return []byte(builder.String())
}

func buildParams(params []string, builder *strings.Builder) {
	for _, param := range params {
		builder.WriteString("<param><value><string>")
		xml.EscapeText(builder, []byte(param))
		builder.WriteString("</string></value></param>")
	}
}
