# ifavicon

ifavicon 是一个用于获取目标网站 favicon hash 的脚本，获取到的 hash 主要用于 fofa 和 shodan 的资产搜索

# Installation

```bash
go install github.com/yuukisec/ifavicon@latest
```

# Use Example

```bash
# 获取目标 URL 的 favicon hash / Get favicon hash from target url
ifavicon -url https://example.com/favicon.ico

# 下载目标 URL 的 favicon / Download favicon from url
ifavicon -download https://example.com/favicon.ico

# 解析文件的 favicon hash / Get favicon hash from target file
ifavicon -file example.com.favicon.ico

# 使用代理请求目标 URL / Use proxy request target URL
ifavicon -url https://example.com/favicon.ico -proxy 127.0.0.1:7890
```
