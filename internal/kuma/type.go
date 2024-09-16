package kuma

import (
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Retry      int64
	Interval   time.Duration
	Token      string
	Auth       AuthStruct
}

type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

type Tag struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type MonitorTag struct {
	Name  string `json:"name,omitempty"`
	TagId int64  `json:"tag_id,omitempty"`
	Value string `json:"value,omitempty"`
}

type Monitor struct {
	ID                                  int64        `json:"id,omitempty"`
	Name                                string       `json:"name"`
	Description                         string       `json:"description,omitempty"`
	PathName                            string       `json:"pathName,omitempty"`
	Parent                              int64        `json:"parent,omitempty"`
	ChildrenIDs                         []int64      `json:"childrenIDs,omitempty"`
	Url                                 string       `json:"url"`
	Method                              string       `json:"method,omitempty"`
	Hostname                            string       `json:"hostname,omitempty"`
	Port                                int64        `json:"port,omitempty"`
	MaxRetries                          int64        `json:"maxretries,omitempty"`
	Weight                              int64        `json:"weight,omitempty"`
	Active                              bool         `json:"active,omitempty"`
	ForceInactive                       bool         `json:"forceInactive,omitempty"`
	Type                                string       `json:"type,omitempty"`
	Timeout                             int64        `json:"timeout,omitempty"`
	Interval                            int64        `json:"interval,omitempty"`
	RetryInterval                       int64        `json:"retryInterval,omitempty"`
	ResendInterval                      int64        `json:"resendInterval,omitempty"`
	Keyword                             string       `json:"keyword,omitempty"`
	InvertKeyword                       bool         `json:"invertKeyword,omitempty"`
	ExpiryNotification                  bool         `json:"expiryNotification,omitempty"`
	IgnoreTls                           bool         `json:"ignoreTls,omitempty"`
	UpsideDown                          bool         `json:"upsideDown,omitempty"`
	PacketSize                          int64        `json:"packetSize,omitempty"`
	MaxRedirects                        int64        `json:"maxredirects,omitempty"`
	AcceptedStatusCodes                 []string     `json:"accepted_statuscodes,omitempty"`
	DNSResolveType                      string       `json:"dns_resolve_type,omitempty"`
	DNSResolveServer                    string       `json:"dns_resolve_server,omitempty"`
	DNSLastResult                       string       `json:"dns_last_result,omitempty"`
	DockerContainer                     string       `json:"docker_container,omitempty"`
	DockerHost                          string       `json:"docker_host,omitempty"`
	ProxyID                             string       `json:"proxyId,omitempty"`
	NotificationIDList                  []int64      `json:"notificationIDList,omitempty"`
	Tags                                []MonitorTag `json:"tags,omitempty"`
	Maintenance                         bool         `json:"maintenance,omitempty"`
	MQTTTopic                           string       `json:"mqttTopic,omitempty"`
	MQTTSuccessMessage                  string       `json:"mqttSuccessMessage,omitempty"`
	DatabaseQuery                       string       `json:"databaseQuery,omitempty"`
	AuthMethod                          string       `json:"authMethod,omitempty"`
	GRPCURL                             string       `json:"grpcUrl,omitempty"`
	GRPCProtobuf                        string       `json:"grpcProtobuf,omitempty"`
	GRPCMethod                          string       `json:"grpcMethod,omitempty"`
	GRPCServiceName                     string       `json:"grpcServiceName,omitempty"`
	GRPCEnableTLS                       bool         `json:"grpcEnableTls,omitempty"`
	RADIUSCalledStationID               string       `json:"radiusCalledStationId,omitempty"`
	RADIUSCallingStationID              string       `json:"radiusCallingStationId,omitempty"`
	Game                                string       `json:"game,omitempty"`
	GamedigGivenPortOnly                bool         `json:"gamedigGivenPortOnly,omitempty"`
	HTTPBodyEncoding                    string       `json:"httpBodyEncoding,omitempty"`
	JSONPath                            string       `json:"jsonPath,omitempty"`
	ExpectedValue                       string       `json:"expectedValue,omitempty"`
	KafkaProducerTopic                  string       `json:"kafkaProducerTopic,omitempty"`
	KafkaProducerBrokers                string       `json:"kafkaProducerBrokers,omitempty"`
	KafkaProducerSSL                    bool         `json:"kafkaProducerSsl,omitempty"`
	KafkaProducerAllowAutoTopicCreation bool         `json:"kafkaProducerAllowAutoTopicCreation,omitempty"`
	KafkaProducerMessage                string       `json:"kafkaProducerMessage,omitempty"`
	Screenshot                          string       `json:"screenshot,omitempty"`
	Headers                             string       `json:"headers,omitempty"`
	Body                                string       `json:"body,omitempty"`
	GRPCBody                            string       `json:"grpcBody,omitempty"`
	GRPCMetadata                        string       `json:"grpcMetadata,omitempty"`
	BasicAuthUser                       string       `json:"basic_auth_user,omitempty"`
	BasicAuthPass                       string       `json:"basic_auth_pass,omitempty"`
	OAuthClientID                       string       `json:"oauth_client_id,omitempty"`
	OAuthClientSecret                   string       `json:"oauth_client_secret,omitempty"`
	OAuthTokenURL                       string       `json:"oauth_token_url,omitempty"`
	OAuthScopes                         string       `json:"oauth_scopes,omitempty"`
	OAuthAuthMethod                     string       `json:"oauth_auth_method,omitempty"`
	PushToken                           string       `json:"pushToken,omitempty"`
	DatabaseConnectionString            string       `json:"databaseConnectionString,omitempty"`
	RADIUSUsername                      string       `json:"radiusUsername,omitempty"`
	RADIUSPassword                      string       `json:"radiusPassword,omitempty"`
	RADIUSSecret                        string       `json:"radiusSecret,omitempty"`
	MQTTUsername                        string       `json:"mqttUsername,omitempty"`
	MQTTPassword                        string       `json:"mqttPassword,omitempty"`
	AuthWorkstation                     string       `json:"authWorkstation,omitempty"`
	AuthDomain                          string       `json:"authDomain,omitempty"`
	TLSCA                               string       `json:"tlsCa,omitempty"`
	TLSCert                             string       `json:"tlsCert,omitempty"`
	TLSKey                              string       `json:"tlsKey,omitempty"`
	KafkaProducerSaslOptions            string       `json:"kafkaProducerSaslOptions,omitempty"`
	IncludeSensitiveData                bool         `json:"includeSensitiveData,omitempty"`
}

type Notification struct {
	ID        int64  `json:"id"`
	UserId    int64  `json:"userId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	IsDefault bool   `json:"isDefault"`
}
