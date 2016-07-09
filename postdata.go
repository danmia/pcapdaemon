package main

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "mime/multipart"
  "net/textproto"
  "net/http"
  "os"
  "path/filepath"
  "crypto/tls"
)


// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName string, path string) (*http.Request, error)  {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    if err != nil {
        return nil, err
    }
    _ = writer.WriteField("filename", path)
    part, err := writer.CreateFormFile(paramName, filepath.Base(path))
    _, err = io.Copy(part, file)

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }

    err = writer.Close()
    if err != nil {
        return nil, err
    }

    myReq, myErr := http.NewRequest("POST", uri, body)
    myReq.Header.Add("Content-Type", "multipart/form-data; boundary=" + writer.Boundary())
    if myErr != nil {
        log.Fatal(myErr)
    }

    return myReq, myErr
}

func newBufferUploadRequest(uri string, params map[string]string, paramName string, buf bytes.Buffer, filename string) (*http.Request, error)  {
    var err error
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    if err != nil {
        return nil, err
    }
    mh := make(textproto.MIMEHeader)
    mh.Set("Content-Type", "application/octet-stream")
    mh.Set("Content-Disposition", "form-data; name=\"file\"; filename=\"" + filename + "\"")
    part, err := writer.CreatePart(mh)
    if nil != err {
        panic(err.Error())
    }
    _, err = io.Copy(part, &buf)

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }

    err = writer.Close()
    if err != nil {
        return nil, err
    }

    myReq, myErr := http.NewRequest("POST", uri, body)
    myReq.Header.Add("Content-Type", "multipart/form-data; boundary=" + writer.Boundary())
    if myErr != nil {
        log.Fatal(myErr)
    }

    return myReq, myErr
}

func postFileCloudshark(scheme string, host string, token string, fname string, tags string)  {

    extraParams := map[string]string{
        "additional_tags":        tags,
    }
    request, err := newfileUploadRequest(scheme + "://" + host + "/api/v1/" + token + "/upload", extraParams, "file", fname)
    if err != nil {
        log.Fatal(err)
    }

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    } else {
        body := &bytes.Buffer{}
        _, err := body.ReadFrom(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        resp.Body.Close()
        fmt.Println(resp.StatusCode)
        fmt.Println(resp.Header)
        fmt.Println(body)
    }
}

func postBufferCloudshark(scheme string, host string, token string, buf bytes.Buffer, filename string, tags string)  {

    extraParams := map[string]string{
        "additional_tags":        tags,
    }
    request, err := newBufferUploadRequest(scheme + "://" + host + "/api/v1/" + token + "/upload", extraParams, "file", buf, filename)
    if err != nil {
        log.Fatal(err)
    }

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    } else {
        body := &bytes.Buffer{}
        _, err := body.ReadFrom(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        resp.Body.Close()
        fmt.Println(resp.StatusCode)
        fmt.Println(resp.Header)
        fmt.Println(body)
    }
}

