package generics

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Online struct {
	Id         uint64 `json:"id"`
	IP         string `json:"ip"`
	ServerIP   string `json:"serverIP"`
	ServerPort int    `json:"serverPort"`
	UpdateAt   string `json:"updateAt"`
}

type Event struct {
	Event    string `json:"event"`
	Protocol int32  `json:"protocol"`
}

type AlarmData struct {
	Alarm  Alarm `json:"alarm"`
	UserID int64 `json:"userID"`
}

type AlarmNotification struct {
	AlarmID   string `json:"alarm_id"`
	AlarmMode string `json:"alarm_mode"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Token     string `json:"token"`
}

type Device struct {
	Id         uint64      `json:"id"`
	ExpireDate null.String `json:"expireDate"`
	ExpireTime string      `json:"expireTime"`
	Model      string      `json:"model"`
	ModelCode  uint32      `json:"modelCode"`
	Expired    bool        `json:"expired"`
	User       null.Int    `json:"user"`
	Updated    string      `json:"updated"`
}

type Location struct {
	DeviceId     uint64    `json:"deviceId" bson:"deviceId"`
	Lng          string    `json:"lng" bson:"longitude"`
	Lat          string    `json:"lat" bson:"latitude"`
	Speed        string    `json:"speed" bson:"speed"`
	CellId       int64     `json:"cellId" bson:"cellId"`
	Mcc          uint16    `json:"mcc" bson:"mcc"`
	SerialNumber uint16    `json:"serialNumber" bson:"serialNumber"`
	Gps          string    `json:"gps" bson:"gpsInformation"`
	Mnc          uint      `json:"mnc" bson:"mnc"`
	Course       string    `json:"course" bson:"course"`
	AccOff       bool      `json:"accOff" bson:"accOff"`
	Lac          uint16    `json:"lac" bson:"lac"`
	Date         string    `json:"pdate" bson:"date"`
	CreatedAt    string    `json:"created_at" bson:"created_at"`
	UpdatedAt    string    `json:"updated_at" bson:"updated_at"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp" `
	Nanoseconds  int64     `json:"nanoseconds" bson:"nanoseconds"`
}

type StringNumber struct {
	String string `json:"string"`
	Number int    `json:"number"`
}

type HeartBeatTerminalInfo struct {
	Status      bool         `json:"status"`
	Ignition    bool         `json:"ignition"`
	Charging    bool         `json:"charging"`
	Alarm       StringNumber `json:"alarm"`
	GpsTracking bool         `json:"gps_tracking"`
	RelayState  bool         `json:"relay_state"`
}

type HeartBeat struct {
	TerminalInfo HeartBeatTerminalInfo `json:"terminal_info"`
	Voltage      StringNumber          `json:"voltage"`
	GsmSignal    StringNumber          `json:"gsm_signal"`
	SerialNumber string                `json:"serial_number"`
	CreatedAt    string                `json:"created_at"`
}

type Information struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type LbsLocation struct {
	DeviceId     uint64 `json:"device_id"`
	MCC          uint16 `json:"mcc"`
	MNC          uint   `json:"mnc"`
	LAC          uint16 `json:"lac"`
	CellId       int64  `json:"cell_id"`
	RSSI         int    `json:"rssi"`
	NLAC1        uint16 `json:"nlac_1"`
	NCellId1     int64  `json:"n_cell_id_1"`
	NRSSI1       int    `json:"nrssi_1"`
	NLAC2        uint16 `json:"nlac_2"`
	NCellId2     int64  `json:"n_cell_id_2"`
	NRSSI2       int    `json:"nrssi_2"`
	NLAC3        uint16 `json:"nlac_3"`
	NCellId3     int64  `json:"n_cell_id_3"`
	NRSSI3       int    `json:"nrssi_3"`
	NLAC4        uint16 `json:"nlac_4"`
	NCellId4     int64  `json:"n_cell_id_4"`
	NRSSI4       int    `json:"nrssi_4"`
	NLAC5        uint16 `json:"nlac_5"`
	NCellId5     int64  `json:"n_cell_id_5"`
	NRSSI5       int    `json:"nrssi_5"`
	NLAC6        uint16 `json:"nlac_6"`
	NCellId6     int64  `json:"n_cell_id_6"`
	NRSSI6       int    `json:"nrssi_6"`
	TimeAdvance  int    `json:"time_advance"`
	Language     string `json:"language"`
	SerialNumber string `json:"serial_number"`
	Date         string `json:"date"`
	CreatedAt    string `json:"created_at"`
}

type Alarm struct {
	DeviceID            uint64    `json:"deviceId" db:"column:deviceId" bson:"deviceId"`
	Latitude            string    `json:"latitude" db:"column:latitude" bson:"latitude"`
	Longitude           string    `json:"longitude" db:"column:longitude" bson:"longitude"`
	Speed               string    `json:"speed" db:"column:speed" bson:"speed"`
	GpsInformation      string    `json:"gpsInformation" db:"column:gpsInformation" bson:"gpsInformation"`
	Course              string    `json:"course" db:"column:course" bson:"course"`
	Mcc                 string    `json:"mcc" db:"column:mcc" bson:"mcc"`
	Mnc                 string    `json:"mnc" db:"column:mnc" bson:"mnc"`
	Lac                 string    `json:"lac" db:"column:lac" bson:"lac"`
	CellID              string    `json:"cellId" db:"column:cellId" bson:"cellId"`
	TerminalInformation string    `json:"terminalInformation" db:"column:terminalInformation" bson:"terminalInformation"`
	VoltageLevel        string    `json:"voltageLevel" db:"column:voltageLevel" bson:"voltageLevel"`
	GsmSignal           string    `json:"gsmSignal" db:"column:gsmSignal" bson:"gsmSignal"`
	AlarmMode           string    `json:"alarmMode" db:"column:alarmMode" bson:"alarmMode"`
	Language            string    `json:"language" db:"column:language" bson:"language"`
	Date                string    `json:"date" db:"column:date" bson:"date"`
	SerialNumber        string    `json:"serialNumber" db:"column:serialNumber" bson:"serialNumber"`
	CreatedAt           string    `json:"created_at" db:"column:created_at" bson:"created_at"`
	UpdatedAt           string    `json:"updated_at" db:"column:updated_at" bson:"updated_at"`
	Timestamp           time.Time `json:"timestamp" bson:"timestamp" `
}

type ParseResponse struct {
	HasDeviceResponse bool      `json:"has_device_response,omitempty"`
	HasSiteResponse   bool      `json:"has_site_response,omitempty"`
	Message           []byte    `json:"message,omitempty"`
	Event             Event     `json:"event,omitempty"`
	Device            Device    `json:"device,omitempty"`
	ResponseTime      time.Time `json:"response_time,omitempty"`
}

type ParseRequest struct {
	Data   []byte `json:"data,omitempty"`
	Imei   string `json:"imei,omitempty"`
	Device Device `json:"device,omitempty"`
}

type DeviceType int32

const (
	DEVICE_TYPE_UNSPECIFIED DeviceType = iota
	DEVICE_TYPE_GT06
	DEVICE_TYPE_X3
	DEVICE_TYPE_MO_PLUS
	DEVICE_TYPE_COOBAN
	DEVICE_TYPE_Q_BIT
	DEVICE_TYPE_V_TRACK
)
