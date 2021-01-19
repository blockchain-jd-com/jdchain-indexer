package response

import (
	"bufio"
	"bytes"
	"encoding/json"
)

type Response struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type SelfJsonSerializable interface {
	ToJSON(writer *bufio.Writer) error
}

func (res *Response) ToJSONNoError(writer *bufio.Writer) (e error) {
	_, e = writer.WriteString(`"success": true, "data": `)
	switch data := res.Data.(type) {
	//case *QueryResult:
	case SelfJsonSerializable:
		e = data.ToJSON(writer)
	default:
		bs, err := json.Marshal(data)
		if err != nil {
			e = err
			return
		}
		_, e = writer.Write(bs)
	}
	return
}

func (res *Response) ToJSONWithError(writer *bufio.Writer) (e error) {
	_, e = writer.WriteString(`"success": false, "error": `)
	bs, err := json.Marshal(res.Error)
	if err != nil {
		e = err
		return
	}
	_, e = writer.Write(bs)
	return
}

func (res *Response) MarshalJSON() (raw []byte, e error) {
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	_, e = writer.WriteString("{")

	if res.Error == nil {
		e = res.ToJSONNoError(writer)
	} else {
		e = res.ToJSONWithError(writer)
	}

	_, e = writer.WriteString("}")

	e = writer.Flush()
	raw = buf.Bytes()
	return
}

type ErrorInfo struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (ei *ErrorInfo) ToJSON(writer *bufio.Writer) (e error) {
	return
}

func NewResponse(data interface{}, err error) *Response {
	if err == nil {
		return NewSuccessResponse(data)
	} else {
		return NewFailedResponse(err.Error())
	}
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func NewFailedResponse(msg string, codes ...int) *Response {
	if len(codes) > 0 {
		return &Response{
			Success: false,
			Error: &ErrorInfo{
				ErrorCode:    codes[0],
				ErrorMessage: msg,
			},
		}
	} else {
		return &Response{
			Success: false,
			Error: &ErrorInfo{
				ErrorCode:    1,
				ErrorMessage: msg,
			},
		}
	}
}
