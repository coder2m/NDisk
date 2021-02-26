/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:36
 **/
package recaptcha

import (
	"github.com/go-resty/resty/v2"
	"github.com/coder2m/component/pkg/xjson"
	"github.com/coder2m/component/xlog"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	apiEndpoint    = "https://www.recaptcha.net/recaptcha/api/siteverify"
	defaultTimeout = 15 * time.Second
)

var (
	// ErrMissingInputSecret when the secret parameter is missing.
	ErrMissingInputSecret = errors.New("missing-input-secret")
	// ErrInvalidInputSecret when the secret parameter is invalid or malformed.
	ErrInvalidInputSecret = errors.New("invalid-input-secret")
	// ErrMissingInputResponse when the response parameter is missing.
	ErrMissingInputResponse = errors.New("missing-input-response")
	// ErrInvalidInputResponse when the response parameter is invalid or malformed.
	ErrInvalidInputResponse = errors.New("invalid-input-response")
	// ErrBadRequest when the request is invalid or malformed.
	ErrBadRequest = errors.New("bad-request")
	// ErrUnsucceeded when response success value equals to false.
	ErrUnsucceeded = errors.New("unsucceeded status")

	ErrDefaultClientNil = errors.New("defaultClient is nil")

	ErrCaptchaNil = errors.New("Captcha is nil")

	errorsMap = map[string]error{
		"missing-input-secret":   ErrMissingInputSecret,
		"invalid-input-secret":   ErrInvalidInputSecret,
		"missing-input-response": ErrMissingInputResponse,
		"invalid-input-response": ErrInvalidInputResponse,
		"bad-request":            ErrBadRequest,
	}

	defaultClient *Client

	nilResponse = new(Response)
)

type (
	// Response holds data from reCaptcha API response.
	Response struct {
		Success     bool        `json:"success"`
		ChallengeTs challengeTs `json:"challenge_ts"`
		Hostname    string      `json:"hostname"`
		ErrorCodes  []string    `json:"error-codes"`
		Score       float64     `json:"score"`
		Action      string      `json:"action"`
	}
	// Config
	Config struct {
		Secret  string        `mapStructure:"secret"`
		Timeout time.Duration `mapStructure:"timeout"`
		Debug   bool          `mapStructure:"debug"`
	}

	// Client struct to verify captcha.
	Client struct {
		secret     string
		httpClient *resty.Client
	}

	challengeTs time.Time
)

func (t *challengeTs) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	parsedTime := parseTime(s)
	*t = challengeTs(parsedTime)
	return nil
}

func (t *challengeTs) String() string {
	return time.Time(*t).String()
}

func parseTime(s string) time.Time {
	if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", s); err == nil {
		return parsedTime
	}

	if parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", s); err == nil {
		return parsedTime
	}

	return time.Now()
}

func DefaultConfig() *Config {
	return &Config{
		Secret:  "",
		Timeout: defaultTimeout,
	}
}

type Option func(*Client)

func SetTimeout(timeout time.Duration) Option {
	return func(cli *Client) {
		cli.httpClient.SetTimeout(timeout)
	}
}

func SerDebug(debug bool) Option {
	return func(cli *Client) {
		cli.httpClient.SetDebug(debug)
	}
}

func NewDefault(c *Config, options ...Option) {
	defaultClient = New(c, options...)
	return
}

func Verify(captcha string) *Response {
	if defaultClient != nil {
		return defaultClient.Verify(captcha)
	}
	return nilResponse
}

func VerifyE(captcha string) (*Response, error) {
	if defaultClient != nil {
		return defaultClient.VerifyE(captcha)
	}
	return nilResponse, ErrDefaultClientNil
}

func VerifyWithIP(captcha, remoteIP string) *Response {
	if defaultClient != nil {
		return defaultClient.VerifyWithIP(captcha, remoteIP)
	}
	return nilResponse
}

func VerifyWithIPE(captcha, remoteIP string) (*Response, error) {
	if defaultClient != nil {
		return defaultClient.VerifyWithIPE(captcha, remoteIP)
	}
	return nilResponse, ErrDefaultClientNil
}

// New returns initialized Client.
func New(c *Config, options ...Option) *Client {
	if c.Secret == "" {
		panic("recaptcha is not nil")
	}
	cli := Client{
		secret: c.Secret,
		httpClient: resty.
			New().
			SetTimeout(c.Timeout).
			SetLogger(xlog.GetDefaultLogger()).
			SetDebug(c.Debug),
	}

	for _, option := range options {
		option(&cli)
	}
	return &cli
}

// Verify verifies reCaptcha response received from frontend.
func (cli *Client) Verify(captcha string) *Response {
	r, err := cli.verify(captcha, "")
	if err != nil {
		xlog.Warn("Verify reCaptcha", xlog.FieldErr(err))
	}
	return r
}

func (cli *Client) VerifyE(captcha string) (*Response, error) {
	return cli.verify(captcha, "")
}

// VerifyWithIP verifies reCaptcha response received from frontend with optional remoteip parameter.
func (cli *Client) VerifyWithIP(captcha, remoteIP string) *Response {
	r, err := cli.verify(captcha, remoteIP)
	if err != nil {
		xlog.Warn("Verify reCaptcha", xlog.FieldErr(err))
	}
	return r
}

func (cli *Client) VerifyWithIPE(captcha, remoteIP string) (*Response, error) {
	return cli.verify(captcha, remoteIP)
}

func (cli *Client) verify(captcha string, remoteIP string) (*Response, error) {
	if captcha == "" {
		return nilResponse, ErrCaptchaNil
	}
	resp, err := cli.httpClient.
		R().
		SetFormData(map[string]string{
			"secret":   cli.secret,
			"response": captcha,
			"remoteip": remoteIP,
		}).
		Post(apiEndpoint)
	if err != nil {
		return nilResponse, errors.Wrap(err, "send request")
	}

	response := new(Response)
	if err := xjson.Unmarshal(resp.Body(), response); err != nil && err != io.EOF {
		return nilResponse, errors.Wrap(err, "unmarshal api response")
	}

	for _, errCode := range response.ErrorCodes {
		if err, ok := errorsMap[errCode]; ok {
			return response, err
		}
	}

	if !response.Success {
		return response, ErrUnsucceeded
	}

	return response, nil
}
