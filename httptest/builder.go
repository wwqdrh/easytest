package httptest

import (
	"encoding/json"
	"strings"
	"testing"
)

type PostmanSpecInfo struct {
	Info struct {
		Name string `json:"name"`
	} `json:"info"`
	Item []*PostmanItem `json:"item"`
}

type PostmanItem struct {
	Name  string `json:"name"`
	Event []struct {
		Listen string `json:"listen"`
		Script struct {
			Exec []string `json:"exec"`
			Type string   `json:"type"`
		} `json:"script"`
	} `json:"event"`
	Request struct {
		Auth struct {
			Type string `json:"type"`
		} `json:"auth"`
		Method string `json:"method"`
		Header []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"header"`
		Body struct {
			Mode    string `json:"mode"`
			Raw     string `json:"raw"`
			Options struct {
				Raw struct {
					Language string `json:"language"`
				} `json:"raw"`
			} `json:"options"`
		} `json:"body"`
		Url struct {
			Raw      string   `json:"raw"`
			Protocol string   `json:"protocol"`
			Host     []string `json:"host"`
			Port     string   `json:"port"`
			Path     []string `json:"path"`
		} `json:"url"`
	} `json:"request"`
	Response []interface{} `json:"response"`
}

type BasicSpecInfo []*BasicItem

type BasicParserSpecInfo BasicSpecInfo

type BasicItem struct {
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	Method      string   `json:"method"`
	Body        string   `json:"body"`
	ContentType string   `json:"content-type"`
	Header      []string `json:"header"`
	Expect      []string `json:"expect"`
	Event       []string `json:"event"`
}

func NewPostmanSpecInfo(data []byte, patch func(item *PostmanItem)) (*PostmanSpecInfo, error) {
	var res PostmanSpecInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if patch != nil {
		for _, item := range res.Item {
			patch(item)
		}
	}
	return &res, nil
}

func (s *PostmanSpecInfo) StartHandle(t *testing.T) error {
	ctx := NewHttpContext()
	for _, item := range s.Item {
		opt := s.specReq2option(item)
		ctx.Do(t, item.Name, opt)
	}
	return nil
}

func (s *PostmanSpecInfo) specReq2option(item *PostmanItem) *HandleOption {
	contentType := ""
	switch item.Request.Body.Options.Raw.Language {
	case "json":
		contentType = "application/json"
	}

	header := map[string]string{}
	for _, item := range item.Request.Header {
		header[item.Key] = item.Value
	}

	return &HandleOption{
		Url:         strings.Join(item.Request.Url.Host, "."),
		Method:      item.Request.Method,
		ContentType: contentType,
		Header:      header,
		Body:        strings.NewReader(item.Request.Body.Raw),
	}
}

func NewBasicSpecInfo(data []byte, patch func(item *BasicItem)) (*BasicSpecInfo, error) {
	var res BasicSpecInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if patch != nil {
		for _, item := range res {
			patch(item)
		}
	}
	return &res, nil
}

func (s *BasicSpecInfo) StartHandle(t *testing.T) error {
	ctx := NewHttpContext()
	for _, item := range *s {
		opt := s.specReq2option(item)
		ctx.Do(t, item.Name, opt)
	}
	return nil
}

func (s *BasicSpecInfo) specReq2option(item *BasicItem) *HandleOption {
	header := map[string]string{}
	for _, item := range item.Header {
		pairs := strings.Split(item, ":")
		if len(pairs) == 2 {
			header[strings.TrimSpace(pairs[0])] = strings.TrimSpace(pairs[1])
		}
	}

	return &HandleOption{
		Url:         item.Url,
		Method:      item.Method,
		ContentType: item.ContentType,
		Header:      header,
		Body:        strings.NewReader(item.Body),
		Expect:      item.Expect,
		Event:       item.Event,
	}
}

func NewBasicParserSpecInfo(data []byte, patch func(item *BasicItem)) (*BasicParserSpecInfo, error) {
	var res BasicParserSpecInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if patch != nil {
		for _, item := range res {
			patch(item)
		}
	}
	return &res, nil
}

func (s *BasicParserSpecInfo) StartHandle(t *testing.T) error {
	ctx := NewHttpContext()
	for _, item := range *s {
		opt := s.specReq2option(item)
		ctx.DoParser(t, item.Name, opt)
	}
	return nil
}

func (s *BasicParserSpecInfo) specReq2option(item *BasicItem) *HandleOption {
	header := map[string]string{}
	for _, item := range item.Header {
		pairs := strings.Split(item, ":")
		if len(pairs) == 2 {
			header[strings.TrimSpace(pairs[0])] = strings.TrimSpace(pairs[1])
		}
	}

	return &HandleOption{
		Url:         item.Url,
		Method:      item.Method,
		ContentType: item.ContentType,
		Header:      header,
		Body:        strings.NewReader(item.Body),
		Expect:      item.Expect,
		Event:       item.Event,
	}
}
