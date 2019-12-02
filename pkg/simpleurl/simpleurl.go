package simpleurl

import (
	"fmt"
	"net/url"
)

func Builder(u string, query map[string]string) string {
	baseUrl, err := url.Parse(u)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return ""
	}

	params := url.Values{}
	for key, value := range query {
		params.Add(key, value)
	}

	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}
