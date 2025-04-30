package utils

import (
    "fmt"
    "net/url"
    "strconv"
    "strings"
)

func GenerateEncodedPayloads(ip string) []string {
    var encoded []string
    encoded = append(encoded, "http://"+ip)
    encoded = append(encoded, url.QueryEscape("http://"+ip))

    if dec := toDecimalIP(ip); dec != "" {
        encoded = append(encoded, "http://"+dec)
    }

    if hex := toHexIP(ip); hex != "" {
        encoded = append(encoded, "http://"+hex)
    }

    if oct := toOctalIP(ip); oct != "" {
        encoded = append(encoded, "http://"+oct)
    }

    encoded = append(encoded, "http://[::ffff:"+ip+"]")
    return encoded
}

func toDecimalIP(ip string) string {
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return ""
    }
    var total int64
    for _, part := range parts {
        p, err := strconv.Atoi(part)
        if err != nil {
            return ""
        }
        total = total<<8 + int64(p)
    }
    return fmt.Sprintf("%d", total)
}

func toHexIP(ip string) string {
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return ""
    }
    hex := "0x"
    for _, part := range parts {
        p, err := strconv.Atoi(part)
        if err != nil {
            return ""
        }
        hex += fmt.Sprintf("%02x", p)
    }
    return hex
}

func toOctalIP(ip string) string {
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return ""
    }
    var octal []string
    for _, part := range parts {
        p, err := strconv.Atoi(part)
        if err != nil {
            return ""
        }
        octal = append(octal, fmt.Sprintf("0%o", p))
    }
    return strings.Join(octal, ".")
}
