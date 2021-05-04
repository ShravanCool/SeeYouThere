package helper

import (
    "net/http"
    "strings"
)

func GetParams(r *http.Request) string {
    p := strings.Split(r.URL.Path, "/")

    return p[len(p)-1]
}
