package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/twmb/murmur3"
)

type Config struct {
	URL      string
	File     string
	Download bool
}

func getContentFromURL(requestURL string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getContentFromFile(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	content, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	h32.Write(raw)
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func standardBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

func main() {
	config := Config{}
	flag.StringVar(&config.File, "file", "", "Get favicon hash from target url")
	flag.StringVar(&config.URL, "url", "", "Get favicon hash from target file")
	flag.BoolVar(&config.Download, "download", false, "Download favicon from url")
	flag.Parse()

	if config.URL == "" && config.File == "" {
		flag.Usage()
		fmt.Println("Example:")
		fmt.Println("  ifavicon -url https://example.com/favicon.ico")
		fmt.Println("  ifavicon -download https://example.com/favicon.ico")
		fmt.Println("  ifavicon -file example.com.favicon.ico")
		os.Exit(0)
	}

	if config.URL != "" {
		if !strings.HasSuffix(config.URL, "favicon.ico") {
			config.URL = config.URL + "/favicon.ico"
		}
		content, err := getContentFromURL(config.URL)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		if config.Download {
			u, err := url.Parse(config.URL)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			err = ioutil.WriteFile(u.Host+".favicon.ico", content, 0644)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
		}
		hash := mmh3Hash32(standardBase64(content))
		output(hash)
		os.Exit(0)
	}

	if config.File != "" {
		content, err := getContentFromFile(config.File)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		hash := mmh3Hash32(standardBase64(content))
		output(hash)
		os.Exit(0)
	}
}

func output(hash string) {
	fmt.Printf("FOFA:\n")
	fmt.Printf("  icon_hash=\"%s\"\n", hash)
	fmt.Printf("  link: https://fofa.info/result?qbase64=%s\n", base64.StdEncoding.EncodeToString([]byte("icon_hash="+hash)))
	fmt.Printf("Shodan:\n")
	fmt.Printf("  http.favicon.hash:%s\n", hash)
	fmt.Printf("  link: https://www.shodan.io/search?query=http.favicon.hash%%3A%s\n", hash)
}
