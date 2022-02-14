package common

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var writeFile = `{"execute":"guest-file-write", "arguments":{"handle":%d,"buf-b64":"%s"}}`

func TestBase64Encode(t *testing.T) {
	base64SSHKey := base64.StdEncoding.EncodeToString([]byte("new"))
	t.Log(base64SSHKey)

	command := fmt.Sprintf(writeFile, 100, base64SSHKey)
	t.Log(command)
}

func TestBase64(t *testing.T) {
	decodeString, err := base64.StdEncoding.DecodeString("eyJhY2NvdW50SUQiOiIzMTQ1MTQyMDk5MDkiLCJhY2NvdW50TmFtZSI6Ind4dGVzdCIsInJlZnJlc2hUaW1lIjoxMjAsImlzcyI6ImN1c29mdHdhcmUiLCJleHAiOjE2MjU4MDA5ODEsImFjY2Vzc1Rva2VuIjoiZTI5ZDQ4NmJiY2UwNGNjMTljZGJjMDlmN2YwZjNiY2MiLCJ1c2VyTmFtZSI6Ind4dGVzdCIsInVzZXJJRCI6IjMxNDUxNDIwOTkwOSJ9")
	if err != nil {

	}
	t.Log(string(decodeString))
}
