package main

import (
	"net/http"
	"strings"
)

func HttpRequestToResty(r *http.Request, packageName, funcName string) string {
	var b strings.Builder

	b.WriteString("package " + packageName + "\n\n")
	b.WriteString("import (\n")
	// b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"github.com/go-resty/resty/v2\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func " + funcName + "() (*resty.Response, error) {\n")
	b.WriteString("\tclient := resty.New()\n\n")

	b.WriteString("\tresp, err := client.R().\n")

	if len(r.URL.Query()) > 0 {
		b.WriteString("\t\tSetQueryParams(map[string]string{\n")
		for key, value := range r.URL.Query() {
			b.WriteString("\t\t\t\"" + key + "\"" + ": \"" + value[0] + "\",\n")
		}
		b.WriteString("\t\t}).\n")
	}

	// Set Headers
	b.WriteString("\t\tSetHeaders(map[string]string{\n")

	for key, value := range r.Header {
		b.WriteString("\t\t\t\"" + key + "\"" + ": \"" + value[0] + "\",\n")
	}
	b.WriteString("\t\t}).\n")
	// Set Headers

	b.WriteString("\t\tSetBody(Input{}).\n")
	b.WriteString("\t\tSetResult(Output{}).\n")

	b.WriteString("\t\tPost(\"" + r.Header.Get("Origin") + r.URL.Path + "\")\n\n")

	b.WriteString("\tbody := resp.Result().(*Output)\n\n")
	b.WriteString("\treturn resp, err\n\n")

	b.WriteString("}\n\n")

	return b.String()
}
