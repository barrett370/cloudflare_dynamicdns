package cloudflare

import "time"

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
type UpdateDNSRecordRequest struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	TTL     int      `json:"ttl"`
	Proxied bool     `json:"proxied"`
	Tags    []string `json:"tags"`
}

type UpdateDNSRecordResponse struct {
	Success  bool              `json:"success"`
	Errors   []ResponseMessage `json:"errors,omitempty"`
	Messages []ResponseMessage `json:"messages,omitempty"`
}

type Zone struct {
	ID                  string    `json:"id"`
	ActivatedOn         time.Time `json:"activated_on"`
	CreatedOn           time.Time `json:"created_on"`
	DevelopmentMode     int       `json:"development_mode"`
	ModifiedOn          time.Time `json:"modified_on"`
	Name                string    `json:"name"`
	OriginalDNShost     string    `json:"original_dnshost"`
	OriginalNameServers []string  `json:"original_name_servers"`
	OriginalRegistrar   string    `json:"original_registrar"`
}

type ListZonesResponse struct {
	Result     []Zone            `json:"result"`
	ResultInfo ResultInfo        `json:"result_info"`
	Success    bool              `json:"success"`
	Errors     []ResponseMessage `json:"errors"`
	Messages   []ResponseMessage `json:"messages"`
}
type Meta struct {
	Step                    int  `json:"step"`
	CustomCertificateQuota  int  `json:"custom_certificate_quota"`
	PageRuleQuota           int  `json:"page_rule_quota"`
	PhishingDetected        bool `json:"phishing_detected"`
	MultipleRailgunsAllowed bool `json:"multiple_railguns_allowed"`
}
type Owner struct {
	ID    interface{} `json:"id"`
	Type  string      `json:"type"`
	Email interface{} `json:"email"`
}
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Plan struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Price             int    `json:"price"`
	Currency          string `json:"currency"`
	Frequency         string `json:"frequency"`
	IsSubscribed      bool   `json:"is_subscribed"`
	CanSubscribe      bool   `json:"can_subscribe"`
	LegacyID          string `json:"legacy_id"`
	LegacyDiscount    bool   `json:"legacy_discount"`
	ExternallyManaged bool   `json:"externally_managed"`
}
type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

type ListDNSRecordsResponse struct {
	Errors     []ResponseMessage `json:"errors"`
	Messages   []ResponseMessage `json:"messages"`
	Result     []Record          `json:"result"`
	Success    bool              `json:"success"`
	ResultInfo ResultInfo        `json:"result_info"`
}
type Data struct {
	ZoneID    string
	EntryID   string
	CurrentIP string
}

type Metadata struct {
	AutoAdded bool   `json:"auto_added"`
	Source    string `json:"source"`
}
type Record struct {
	Comment    string    `json:"comment"`
	Content    string    `json:"content"`
	CreatedOn  time.Time `json:"created_on"`
	Data       Data      `json:"data"`
	ID         string    `json:"id"`
	Locked     bool      `json:"locked"`
	Meta       Metadata  `json:"meta"`
	ModifiedOn time.Time `json:"modified_on"`
	Name       string    `json:"name"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	Tags       []string  `json:"tags"`
	TTL        int       `json:"ttl"`
	Type       string    `json:"type"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
}
