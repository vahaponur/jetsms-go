package jetsms

// SendSmsRequest models the JSON body for /api/SendSms
// excluding authentication fields (user, password).
type SendSmsRequest struct {
	Originator       string       `json:"originator"`
	Reference        string       `json:"reference,omitempty"`
	StartDate        string       `json:"startdate,omitempty"`
	ExpireDate       string       `json:"expiredate,omitempty"`
	BroadcastMessage string       `json:"broadcastmessage,omitempty"`
	BroadcastMessageParametric string `json:"broadcastmessageparametric,omitempty"`
	SmsMessages      []SmsMessage `json:"smsmessages"`
	ExclusionStartTime string     `json:"exclusionstarttime,omitempty"`
	ExclusionExpireTime string    `json:"exclusionexpiretime,omitempty"`
	Channel          string       `json:"channel,omitempty"`
	BlacklistFilter  string       `json:"blacklistfilter,omitempty"`
	DomesticFilter   string       `json:"domesticfilter,omitempty"`
	ForceEnglish     string       `json:"forceenglish,omitempty"`
	IYSFilter        string       `json:"iysfilter,omitempty"`
	IYSCode          string       `json:"iyscode,omitempty"`
	BrandCode        string       `json:"brandcode,omitempty"`
	RetailerCode     string       `json:"retailercode,omitempty"`
	RecipientType    string       `json:"recipienttype,omitempty"`
	State            int          `json:"state,omitempty"`
	OrigPID          string       `json:"origpid,omitempty"`
	TotalCount       string       `json:"totalcount,omitempty"`
	Sequence         string       `json:"sequence,omitempty"`
	Multipart        string       `json:"multipart,omitempty"`
	MultichannelType string       `json:"multichanneltype,omitempty"`
	Multichannels    []string     `json:"multichannels,omitempty"`
	Multioriginators []string     `json:"multioriginators,omitempty"`
	SimCheckInDay    string       `json:"simcheckinday,omitempty"`
	MNPCheckInDay    string       `json:"mnpcheckinday,omitempty"`
	BWFilter         string       `json:"bwfilter,omitempty"`
	X                string       `json:"x,omitempty"`
	Realtime         string       `json:"realtime,omitempty"`
}

// SmsMessage represents an individual message in SendSmsRequest.
type SmsMessage struct {
	MessageText   string   `json:"messagetext,omitempty"`
	Recipient     string   `json:"recipient"`
	MessageID     string   `json:"messageid,omitempty"`
	MessageParams []string `json:"messageparams,omitempty"`
	// MessageType corresponds to the messageType field in the API.
	// Values: 0:text,1:binary,2:wapPush,3:flashSMS,4:unicode,5:iot binary.
	MessageType   string   `json:"messageType,omitempty"`
	MessageHeader string   `json:"messageHeader,omitempty"`
}

// Internal payload including authentication for /api/SendSms.
type sendSmsRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	SendSmsRequest
}

// SendSmsResponse models the response of /api/SendSms.
type SendSmsResponse struct {
	ResponseCode        string                  `json:"responseCode"`
	ResponseMessage     string                  `json:"responseMessage"`
	ResponseGroupIDList []SendSmsGroupIDElement `json:"responseGroupIdArray"`
}

// SendSmsGroupIDElement represents a single entry in responseGroupIdArray.
type SendSmsGroupIDElement struct {
	GSMNumber string `json:"gsmNumber"`
	MessageID string `json:"messageid"`
	Status    string `json:"status"`
}

// CancelSmsRequest models the request for /api/CancelSms excluding auth fields.
type CancelSmsRequest struct {
	ProcessID string `json:"processid"`
	MessageID string `json:"messageid,omitempty"`
}

type cancelSmsRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	CancelSmsRequest
}

// CancelSmsResponse models the response for /api/CancelSms.
type CancelSmsResponse struct {
	Status string `json:"status"`
}

// ReportSmsDetailRequest models the request for /api/ReportSmsDetail excluding auth fields.
type ReportSmsDetailRequest struct {
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
	GSMNumber string `json:"gsmnumber,omitempty"`
	Originator string `json:"originator,omitempty"`
}

type reportSmsDetailRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	ReportSmsDetailRequest
}

// ReportSmsDetailResponse models the response for /api/ReportSmsDetail.
type ReportSmsDetailResponse struct {
	ReportSmsDetailed []ReportSmsDetailItem `json:"ReportSMSDetailed"`
}

type ReportSmsDetailItem struct {
	ID             string `json:"ID"`
	Phone          string `json:"Phone"`
	Status         string `json:"Status"`
	SendDate       string `json:"SendDate"`
	DeliveredDate  string `json:"DeliveredDate"`
	Originator     string `json:"Originator"`
	CreateDate     string `json:"CreateDate"`
	MessageText    string `json:"MessageText"`
	ResponseCode   string `json:"ResponseCode"`
	SendRetryCount string `json:"SendRetryCount"`
	SendedSMSCount string `json:"SendedSMSCount"`
	SMSValueToLength string `json:"SMSValueToLength"`
}

// ReportSmsRequest models the request for /api/ReportSms excluding auth fields.
type ReportSmsRequest struct {
	GroupID    string   `json:"groupid"`
	MessageIDs []string `json:"messageIds,omitempty"`
}

type reportSmsRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	ReportSmsRequest
}

// ReportSmsResponse models the response for /api/ReportSms.
type ReportSmsResponse struct {
	ResponseCode     string              `json:"ResponseCode"`
	ResponseMessage  string              `json:"ResponseMessage"`
	MessageReportArr []MessageReportItem `json:"MessageReportArray"`
}

type MessageReportItem struct {
	GSMNumber       string `json:"GsmNumber"`
	Status          string `json:"Status"`
	SentDate        string `json:"SentDate"`
	DeliveredDate   string `json:"DeliveredDate"`
	MessageID       string `json:"MessageId"`
	ResponseCode    string `json:"ResponseCode"`
	SendRetryCount  string `json:"SendRetryCount"`
	SMSValueToLength string `json:"SMSValueToLength"`
	SendedSMSCount  string `json:"SendedSMSCount"`
	Channel         string `json:"Channel"`
}

// ReportSmsSummaryRequest models the request for /api/ReportSmsSummary excluding auth fields.
type ReportSmsSummaryRequest struct {
	GroupID    string   `json:"groupid"`
	MessageIDs []string `json:"messageIds,omitempty"`
}

type reportSmsSummaryRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	ReportSmsSummaryRequest
}

// ReportSmsSummaryResponse models the response for /api/ReportSmsSummary.
type ReportSmsSummaryResponse struct {
	ProcessID        string `json:"ProcessId"`
	MessageCount     int    `json:"MessageCount"`
	MessageSuccessful int   `json:"MessageSuccessful"`
	MessageWaiting   int    `json:"MessageWaiting"`
	MessageFailed    int    `json:"MessageFailed"`
	MessageOverdue   int    `json:"MessageOverdue"`
	ResponseCode     int    `json:"ResponseCode"`
	// NOTE: The OpenAPI spec defines ResponseMessage as number,
	// but the description says it is a message string.
	// We keep it as string here; if the backend really uses a number
	// this can be adjusted later.
	ResponseMessage  string `json:"ResponseMessage"`
	ProcessStatusID  int    `json:"ProcessStatusId"`
}

// ReportOriginatorRequest models the request for /api/reportoriginator excluding auth fields.
type ReportOriginatorRequest struct{}

type reportOriginatorRequestPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
	ReportOriginatorRequest
}

// ReportOriginatorResponse models the response for /api/reportoriginator.
type ReportOriginatorResponse struct {
	Originators []OriginatorItem `json:"originators"`
}

type OriginatorItem struct {
	Channel      string `json:"channel"`
	MessageTitle string `json:"messagetitle"`
}
