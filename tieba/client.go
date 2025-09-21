package tieba

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// API URLs
const (
	urlTbs  = "https://tieba.baidu.com/dc/common/tbs"
	urlLike = "https://c.tieba.baidu.com/c/f/forum/like"
	urlSign = "https://c.tieba.baidu.com/c/c/forum/sign"
)

type Client struct {
	// 百度使用的 token BDUSS
	bduss string

	client *http.Client

	log *slog.Logger

	// TODO 邮件通知等？
}

func WithLog(log *slog.Logger) Option {
	return func(c *Client) {
		c.log = log
	}
}

type Option func(c *Client)

func NewClient(bduss string, opts ...Option) (*Client, error) {
	if bduss == "" {
		return nil, errors.New("bduss is required")
	}

	client := &Client{
		bduss: bduss,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		log: slog.Default(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) header() http.Header {
	h := http.Header{}

	h.Add("Host", "tieba.baidu.com")
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	h.Add("Content-Type", "application/x-www-form-urlencoded")

	cookie := http.Cookie{
		Name:  "BDUSS",
		Value: c.bduss,
	}

	h.Add("Cookie", cookie.String())

	return h
}

func (c *Client) doWithJSON(req *http.Request, point any) error {
	for k, values := range c.header() {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	// 错误检查
	var check Error

	err = json.Unmarshal(body, &check)
	if err != nil {
		return fmt.Errorf("json decode check %s , err: %w", req.Host, err)
	}

	if check.Msg != "" {
		return &check
	}

	err = json.Unmarshal(body, point)
	if err != nil {
		return fmt.Errorf("json decode point %s , err: %w", req.Host, err)
	}

	return nil
}

func (c *Client) Tbs(ctx context.Context, request *TbsRequest) (*TbsResponse, error) {
	c.log.InfoContext(ctx, "Tbs Start")
	defer c.log.InfoContext(ctx, "Tbs Finished")

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, urlTbs, nil)
	if err != nil {
		return nil, err
	}

	var response TbsResponse

	err = c.doWithJSON(r, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Favorite
// 获取关注的贴吧列表.
func (c *Client) Favorite(ctx context.Context, request *FavoriteRequest) (*FavoriteResponse, error) {
	c.log.InfoContext(ctx, "Favorite Start")
	defer c.log.InfoContext(ctx, "Favorite Finished")

	request.pageNo = 0

	response := FavoriteResponse{
		HasMore: "1",
	}

	// 处理分页
	for response.HasMore == "1" {
		time.Sleep(time.Second * 3)

		request.pageNo++

		res, err := c.favorite(ctx, request)
		if err != nil {
			c.log.ErrorContext(ctx, "Favorite", "HasMore Request", request, "err", err)

			break
		}

		c.log.DebugContext(ctx, "Favorite", "HasMore res", res)

		// 合并
		res.ForumList.NonGconForum = append(response.ForumList.NonGconForum, res.ForumList.NonGconForum...)
		res.ForumList.GconForum = append(response.ForumList.GconForum, res.ForumList.GconForum...)

		response = *res
	}

	return &response, nil
}

func (c *Client) favorite(ctx context.Context, request *FavoriteRequest) (*FavoriteResponse, error) {
	data := map[string]string{
		"_client_type":    "2",
		"_client_id":      "wappc_1534235498291_488",
		"_client_version": "9.7.8.0",
		"_phone_imei":     "000000000000000",
		"from":            "1008621y",
		"model":           "MI+5",
		"net_type":        "1",
		"vcode_tag":       "11",

		"BDUSS": c.bduss,

		"page_size": Itoa(request.PageSize),
		"page_no":   Itoa(request.pageNo),

		"timestamp": Timestamp(),
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, urlLike, c.urlEncode(data))
	if err != nil {
		return nil, err
	}

	var response FavoriteResponse

	err = c.doWithJSON(r, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Sign
// 贴吧签到.
func (c *Client) Sign(ctx context.Context, request *SignRequest) (*SignResponse, error) {
	data := map[string]string{
		"_client_type":    "2",
		"_client_id":      "wappc_1534235498291_488",
		"_client_version": "9.7.8.0",
		"_phone_imei":     "000000000000000",
		"from":            "1008621y",
		"model":           "MI+5",
		"net_type":        "1",
		"vcode_tag":       "11",

		"BDUSS": c.bduss,

		"tbs": request.Tbs,
		"fid": request.Fid,
		"kw":  request.KW,

		"timestamp": Timestamp(),
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, urlSign, c.urlEncode(data))
	if err != nil {
		return nil, err
	}

	var response SignResponse

	err = c.doWithJSON(r, &response)
	if err != nil {
		var check *Error
		if errors.As(err, &check) && check.Code == codeSignRepeat {
			return &response, nil
		}

		return nil, err
	}

	c.log.DebugContext(ctx, "Sign", "request", request)
	c.log.DebugContext(ctx, "Sign", "response", response)

	return &response, nil
}

// 根据规则，编码数据.
func (c *Client) urlEncode(data map[string]string) io.Reader {
	x := ""

	// 1. 字段 key 排序
	// 2. 拼接 k1=v1k2=v2...
	// 3. 生成 md5签名
	// 4. 转换为 表单数据
	values := url.Values{}

	for _, k := range slices.Sorted(maps.Keys(data)) {
		v := data[k]
		x += k + "=" + v
		values.Add(k, v)
	}

	hash := md5.Sum([]byte(x + "tiebaclient!!!"))
	sign := strings.ToUpper(hex.EncodeToString(hash[:]))

	values.Add("sign", sign)

	return strings.NewReader(values.Encode())
}
