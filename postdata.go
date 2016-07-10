package main

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "strconv"
  "mime/multipart"
  "net/textproto"
  "net/http"
  "crypto/tls"
)


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

func postBufferCloudshark(scheme string, host string, port int, token string, buf bytes.Buffer, filename string, tags string)  {

    var url string
    extraParams := map[string]string{
        "additional_tags":        tags,
    }
    if(port != 80 && port != 443)  {
        url = scheme + "://" + host + ":" + strconv.Itoa(port) + "/api/v1/" + token + "/upload"
    } else  {
        url = scheme + "://" + host + "/api/v1/" + token + "/upload"
    }

    request, err := newBufferUploadRequest(url, extraParams, "file", buf, filename)
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

