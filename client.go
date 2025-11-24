package jetsms

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	defaultBaseURL = "https://smsapi.jetsms.com.tr"
)

// Config holds configuration for the Client.
type Config struct {
	// Username is the JetSMS API username.
	Username string
	// Password is the JetSMS API password.
	Password string
	// Originator is the default SMS sender title (originator).
	// It can be overridden per request.
	Originator string
	// BaseURL is the API base URL. If empty, the default JetSMS URL is used.
	BaseURL string
	// HTTPClient is an optional custom HTTP client.
	HTTPClient *http.Client
}

// Client is a JetSMS REST API client.
//
// Authentication is performed by sending Username and Password in the
// JSON request body for each call, as required by the JetSMS REST API.
type Client struct {
	resty      *resty.Client
	username   string
	password   string
	originator string
}

// NewClient creates a new JetSMS client using the given configuration.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	var httpClient *http.Client
	if cfg.HTTPClient != nil {
		httpClient = cfg.HTTPClient
	} else {
		httpClient = http.DefaultClient
	}

	r := resty.NewWithClient(httpClient)
	r.SetBaseURL(baseURL)

	return &Client{
		resty:      r,
		username:   cfg.Username,
		password:   cfg.Password,
		originator: cfg.Originator,
	}, nil
}

// SendOptions provides optional parameters for high level send helpers.
type SendOptions struct {
	// Originator overrides the client's default originator for this request.
	Originator string
	// Reference is an optional reference value.
	Reference string
	// Channel is the JetSMS channel code (VD, TRKC, ...).
	Channel string
	// StartDate is the scheduled start date (ddMMyyyyhhmmss or relative like s10/m15/h3).
	StartDate string
	// ExpireDate is the expiry date (ddMMyyyyhhmmss or relative).
	ExpireDate string
	// Realtime corresponds to the realtime parameter; when "1", the API
	// attempts to send immediately and return operator-level response.
	Realtime string
}

// SendOne sends a single SMS (1-1) to the given recipient.
func (c *Client) SendOne(ctx context.Context, recipient, message string, opts *SendOptions) (*SendSmsResponse, error) {
	if recipient == "" {
		return nil, fmt.Errorf("recipient is required")
	}
	if message == "" {
		return nil, fmt.Errorf("message is required")
	}

	req := SendSmsRequest{
		Originator: c.resolveOriginator(opts),
		Reference:  valueOrEmpty(opts, func(o *SendOptions) string { return o.Reference }),
		StartDate:  valueOrEmpty(opts, func(o *SendOptions) string { return o.StartDate }),
		ExpireDate: valueOrEmpty(opts, func(o *SendOptions) string { return o.ExpireDate }),
		Channel:    valueOrEmpty(opts, func(o *SendOptions) string { return o.Channel }),
		SmsMessages: []SmsMessage{
			{
				MessageText: message,
				Recipient:   recipient,
			},
		},
		Realtime: valueOrEmpty(opts, func(o *SendOptions) string { return o.Realtime }),
	}

	return c.SendSMS(ctx, req)
}

// SendMany sends the same SMS text to multiple recipients (1-n).
func (c *Client) SendMany(ctx context.Context, recipients []string, message string, opts *SendOptions) (*SendSmsResponse, error) {
	if len(recipients) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}
	if message == "" {
		return nil, fmt.Errorf("message is required")
	}

	smsMessages := make([]SmsMessage, 0, len(recipients))
	for _, rcp := range recipients {
		if rcp == "" {
			continue
		}
		smsMessages = append(smsMessages, SmsMessage{
			Recipient: rcp,
		})
	}
	if len(smsMessages) == 0 {
		return nil, fmt.Errorf("no valid recipients provided")
	}

	req := SendSmsRequest{
		Originator:       c.resolveOriginator(opts),
		Reference:        valueOrEmpty(opts, func(o *SendOptions) string { return o.Reference }),
		StartDate:        valueOrEmpty(opts, func(o *SendOptions) string { return o.StartDate }),
		ExpireDate:       valueOrEmpty(opts, func(o *SendOptions) string { return o.ExpireDate }),
		Channel:          valueOrEmpty(opts, func(o *SendOptions) string { return o.Channel }),
		BroadcastMessage: message,
		SmsMessages:      smsMessages,
		Realtime:         valueOrEmpty(opts, func(o *SendOptions) string { return o.Realtime }),
	}

	return c.SendSMS(ctx, req)
}

// SendSMS executes the generic /api/SendSms call using the provided request.
//
// Note: HTTP-level errors are returned as Go errors. For application-level
// errors you must inspect resp.ResponseCode and resp.ResponseMessage.
func (c *Client) SendSMS(ctx context.Context, req SendSmsRequest) (*SendSmsResponse, error) {
	if req.Originator == "" {
		if c.originator == "" {
			return nil, fmt.Errorf("originator is required")
		}
		req.Originator = c.originator
	}

	payload := sendSmsRequestPayload{
		User:           c.username,
		Password:       c.password,
		SendSmsRequest: req,
	}

	var respBody SendSmsResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/SendSms")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

// CancelSMS cancels SMS messages for the given process (group) id and optional message id.
func (c *Client) CancelSMS(ctx context.Context, req CancelSmsRequest) (*CancelSmsResponse, error) {
	if req.ProcessID == "" {
		return nil, fmt.Errorf("process id is required")
	}

	payload := cancelSmsRequestPayload{
		User:             c.username,
		Password:         c.password,
		CancelSmsRequest: req,
	}

	var respBody CancelSmsResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/CancelSms")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

// ReportSmsDetail retrieves detailed SMS reports for the given filters.
func (c *Client) ReportSmsDetail(ctx context.Context, req ReportSmsDetailRequest) (*ReportSmsDetailResponse, error) {
	if req.StartDate == "" || req.EndDate == "" {
		return nil, fmt.Errorf("start date and end date are required")
	}

	payload := reportSmsDetailRequestPayload{
		User:                    c.username,
		Password:                c.password,
		ReportSmsDetailRequest:  req,
	}

	var respBody ReportSmsDetailResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/ReportSmsDetail")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

// ReportSms retrieves message-level report information for a given group id.
func (c *Client) ReportSms(ctx context.Context, req ReportSmsRequest) (*ReportSmsResponse, error) {
	if req.GroupID == "" {
		return nil, fmt.Errorf("group id is required")
	}

	payload := reportSmsRequestPayload{
		User:              c.username,
		Password:          c.password,
		ReportSmsRequest:  req,
	}

	var respBody ReportSmsResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/ReportSms")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

// ReportSmsSummary retrieves summary report information for a given group id.
func (c *Client) ReportSmsSummary(ctx context.Context, req ReportSmsSummaryRequest) (*ReportSmsSummaryResponse, error) {
	if req.GroupID == "" {
		return nil, fmt.Errorf("group id is required")
	}

	payload := reportSmsSummaryRequestPayload{
		User:                     c.username,
		Password:                 c.password,
		ReportSmsSummaryRequest:  req,
	}

	var respBody ReportSmsSummaryResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/ReportSmsSummary")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

// ListOriginators lists SMS originators (titles) for the authenticated user.
func (c *Client) ListOriginators(ctx context.Context) (*ReportOriginatorResponse, error) {
	payload := reportOriginatorRequestPayload{
		User: c.username,
		Password: c.password,
		ReportOriginatorRequest: ReportOriginatorRequest{},
	}

	var respBody ReportOriginatorResponse
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&respBody).
		Post("/api/reportoriginator")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, httpError(resp)
	}

	return &respBody, nil
}

func (c *Client) resolveOriginator(opts *SendOptions) string {
	if opts != nil && opts.Originator != "" {
		return opts.Originator
	}
	return c.originator
}

func valueOrEmpty[T any](opts *SendOptions, getter func(*SendOptions) T) T {
	var zero T
	if opts == nil {
		return zero
	}
	return getter(opts)
}

