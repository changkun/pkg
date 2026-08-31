// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Doc: https://www.xfyun.cn/doc/asr/lfasr/API.html#%E6%8E%A5%E5%8F%A3%E8%B0%83%E7%94%A8%E6%B5%81%E7%A8%8B

const (
	defaultPartSize   = 10 * 1024 * 1024
	defaultRetryTimes = 3
	defaultUA         = "changkun/autran"
	defaultDomain     = "https://raasr.xfyun.cn/api"
)

// Conf config struct
type Conf struct {
	AppID      string
	SecretKey  string
	PartSize   int64
	RetryTimes int
	Ch         string
	UA         string
	Domain     string
}

func getDefaultConf() *Conf {
	conf := Conf{}
	conf.PartSize = defaultPartSize
	conf.RetryTimes = defaultRetryTimes
	conf.UA = defaultUA
	conf.Domain = defaultDomain

	return &conf
}

// Conn is a connection for accessing xfyun APIs.
type Conn struct {
	c    *http.Client
	conf *Conf
}

// Client ...
type Client struct {
	conn *Conn
	// sliceID is the odometer behind getNextSliceID. It used to be a package
	// level variable, so two clients uploading at once shared one counter and
	// handed the service overlapping slice identifiers.
	sliceID []byte
}

// RespInfo ...
type RespInfo struct {
	Ok     int    `json:"ok"`
	ErrNo  int    `json:"err_no"`
	Failed string `json:"failed"`
	Data   string `json:"data"`
}

// New ...
func New(appID, secretKey string) *Client {
	client := Client{}
	conf := getDefaultConf()
	conf.AppID = appID
	conf.SecretKey = secretKey

	conn := Conn{&http.Client{}, conf}
	client.conn = &conn
	// One before "aaaaaaaaaa", so the first identifier handed out is that.
	client.sliceID = []byte("aaaaaaaaa`")

	return &client
}

// UploadAudio ...
func (c *Client) UploadAudio(ctx context.Context, filename, language string) (taskid string, err error) {
	filesize, sliceNum, err := c.conn.getSizeAndSiceNum(filename)
	if err != nil {
		return
	}
	log.Printf("filesize: %v, sliceNum: %v", filesize, sliceNum)

	taskid, err = c.initSliceUpload(ctx, filename, language, filesize, sliceNum)
	if err != nil {
		return
	}

	log.Printf("taskid: %v", taskid)

	if err = c.performSliceUpload(ctx, filename, taskid, filesize, sliceNum); err != nil {
		return
	}

	log.Println("upload is complete.")

	if err = c.completeSliceUpload(ctx, taskid); err != nil {
		return
	}

	log.Println("merge is complete.")
	return
}

func (c *Client) initSliceUpload(ctx context.Context, filename, language string, filesize, sliceNum int64) (taskid string, err error) {
	var info RespInfo
	params := c.getBaseAuthParam("")
	params.Add("file_len", strconv.FormatInt(filesize, 10))
	params.Add("file_name", filename)
	params.Add("slice_num", strconv.FormatInt(sliceNum, 10))
	params.Add("has_participle", "false")
	params.Add("max_alternatives", "0")
	params.Add("speaker_number", "2")
	params.Add("has_seperate", "true") //nolint:misspell // the vendor API spells the field this way
	params.Add("language", language)
	params.Add("pd", "tech")

	resp, err := c.conn.httpDo(ctx, c.conn.conf.Domain+"/prepare", nil, params, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return
	}

	if info.Ok == 0 {
		taskid = info.Data
	} else {
		err = fmt.Errorf("info: %v, errno: %v", info.Failed, info.ErrNo)
	}

	return
}
func (c *Client) performSliceUpload(ctx context.Context, filename, taskid string, filesize, sliceNum int64) (err error) {
	log.Println("start uploading...")
	var info RespInfo
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	b := make([]byte, c.conn.conf.PartSize)
	for i := int64(1); i <= sliceNum; i++ {
		log.Printf("uploading slice %d...", i)
		offset := (i - 1) * c.conn.conf.PartSize
		if remaining := filesize - offset; len(b) > int(remaining) {
			b = make([]byte, remaining)
		}
		// ReadAt reads the slice in full or reports why it could not. A short
		// read used to go unnoticed and upload whatever the buffer still held
		// from the previous slice.
		if _, err := fi.ReadAt(b, offset); err != nil {
			return fmt.Errorf("cannot read slice %d at offset %d: %w", i, offset, err)
		}

		params := c.getBaseAuthParam(taskid)
		params.Add("slice_id", c.getNextSliceID())
		resp, err := c.conn.postMulti(ctx, c.conn.conf.Domain+"/upload", filename, b, params)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(resp, &info); err != nil {
			return err
		}

		if info.Ok != 0 {
			return fmt.Errorf("info: %v, errno: %v", info.Failed, info.ErrNo)
		}
	}
	return nil
}

func (c *Client) completeSliceUpload(ctx context.Context, taskid string) (err error) {
	log.Println("perform merge action...")

	params := c.getBaseAuthParam(taskid)
	resp, err := c.conn.httpDo(ctx, c.conn.conf.Domain+"/merge", nil, params, nil)
	if err != nil {
		return
	}
	var info RespInfo
	if err = json.Unmarshal(resp, &info); err != nil {
		return
	}

	if info.Ok != 0 {
		return fmt.Errorf("info: %v, errno: %v", info.Failed, info.ErrNo)
	}

	return nil
}

// GetProgress ...
func (c *Client) GetProgress(ctx context.Context, taskid string) (status int, err error) {
	params := c.getBaseAuthParam(taskid)
	resp, err := c.conn.httpDo(ctx, c.conn.conf.Domain+"/getProgress", nil, params, nil)
	if err != nil {
		return
	}
	var info RespInfo
	if err = json.Unmarshal(resp, &info); err != nil {
		return
	}

	type progressData struct {
		Description string `json:"desc"`
		Status      int    `json:"status"`
	}

	if info.Ok != 0 {
		return 0, fmt.Errorf("info: %v, errno: %d", info.Failed, info.ErrNo)
	}

	p := progressData{}
	err = json.Unmarshal([]byte(info.Data), &p)
	if err != nil {
		return
	}

	return p.Status, nil
}

// GetResult ...
func (c *Client) GetResult(ctx context.Context, taskid string) (content string, err error) {
	params := c.getBaseAuthParam(taskid)
	resp, err := c.conn.httpDo(ctx, c.conn.conf.Domain+"/getResult", nil, params, nil)
	if err != nil {
		return
	}
	var info RespInfo
	if err = json.Unmarshal(resp, &info); err != nil {
		return
	}

	if info.Ok != 0 {
		return info.Data, fmt.Errorf("info: %v, errno: %v", info.Failed, info.ErrNo)
	}

	return info.Data, nil
}

func (c *Conn) postMulti(ctx context.Context, uri, filename string, content []byte, params url.Values) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, ferr := writer.CreateFormFile("content", filename+params.Get("slice_id"))
	if ferr != nil {
		return nil, ferr
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		return nil, fmt.Errorf("cannot write the slice into the form: %w", err)
	}

	for key, val := range params {
		if err := writer.WriteField(key, val[0]); err != nil {
			return nil, fmt.Errorf("cannot write field %s: %w", key, err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (c *Conn) httpDo(ctx context.Context, url string, body []byte, params url.Values, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Conn) getSizeAndSiceNum(filename string) (filesize, num int64, err error) {
	filesize, err = fileSize(filename)
	if err != nil {
		return
	}
	num = int64(math.Ceil(float64(filesize) / float64(c.conf.PartSize)))
	return
}

func (c *Client) getBaseAuthParam(taskid string) url.Values {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha1.New, []byte(c.conn.conf.SecretKey))
	strByte := []byte(c.conn.conf.AppID + ts)
	strMd5Byte := md5.Sum(strByte)
	strMd5 := hex.EncodeToString(strMd5Byte[:])
	mac.Write([]byte(strMd5))
	signa := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	params := url.Values{}
	params.Add("app_id", c.conn.conf.AppID)
	params.Add("signa", signa)
	params.Add("ts", ts)
	if len(taskid) > 0 {
		params.Add("task_id", taskid)
	}

	return params
}

// getNextSliceID returns the identifier of the next slice. The service
// numbers slices with a lowercase base 26 odometer, so the sequence runs
// aaaaaaaaaa, aaaaaaaaab and carries leftwards when a position passes z.
func (c *Client) getNextSliceID() string {
	for i := len(c.sliceID) - 1; i >= 0; i-- {
		if c.sliceID[i] != 'z' {
			c.sliceID[i]++
			break
		}
		c.sliceID[i] = 'a'
	}
	return string(c.sliceID)
}

func fileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return 0, err
	}
	return info.Size(), nil
}
