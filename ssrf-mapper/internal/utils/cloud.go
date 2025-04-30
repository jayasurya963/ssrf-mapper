package utils

import (
    "net/http"
    "strings"
)

func DetectCloud(resp *http.Response) string {
    headers := resp.Header

    if strings.Contains(headers.Get("Server"), "EC2") || headers.Get("X-Amzn-Trace-Id") != "" {
        return "AWS"
    }

    if headers.Get("Metadata-Flavor") == "Google" {
        return "GCP"
    }

    if headers.Get("Server") == "Microsoft-IIS" || headers.Get("x-ms-request-id") != "" {
        return "Azure"
    }

    return "Unknown"
}
