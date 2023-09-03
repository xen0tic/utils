package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xen0tic/utils"
)

type Options struct {
	Username     string
	Password     string
	Server       string
	Port         int
	Database     string
	MaxIdleTime  time.Duration
	MaxLifeTime  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

type MySQL struct {
	client *sql.DB
}

type InfoCount struct {
	cn uint64
}

type SyncInfo struct {
	Id         uint64         `json:"id"`
	UserId     sql.NullInt32  `json:"user"`
	DeviceImei string         `json:"deviceImei"`
	ExpireDate sql.NullString `json:"expireDate"`
	ExpireTime string         `json:"expireTime"`
	ModelName  string         `json:"modelName"`
	Model      uint32         `json:"modelCode"`
}

type PushToken struct {
	PushToken string `json:"pushToken"`
}

func New(option *Options) *MySQL {
	instance := new(MySQL)
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", option.Username, option.Password, option.Server, option.Port, option.Database))
	
	if err != nil {
		panic(err.Error())
	}
	
	db.SetConnMaxLifetime(option.MaxLifeTime)
	db.SetConnMaxIdleTime(option.MaxIdleTime)
	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetMaxOpenConns(option.MaxOpenConns)
	
	instance.client = db
	
	return instance
}

func (s *MySQL) insertQuery(query string) (int64, bool) {
	res, err := s.client.Exec(query)
	
	if err != nil {
		return 0, false
	}
	lasID, err1 := res.LastInsertId()
	if err1 != nil {
		return 0, false
	}
	return lasID, true
}

func (s *MySQL) Insert(table string, column []string, values []string) (int64, string, error) {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, strings.Join(column, ","), strings.Join(values, ","))
	res, err := s.client.Exec(query)
	
	if err != nil {
		return 0, query, err
	}
	lasID, err1 := res.LastInsertId()
	if err1 != nil {
		return 0, query, err
	}
	return lasID, query, nil
}

func (s *MySQL) Update(table string, column []string, where string) (int64, bool) {
	query := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(column, ","))
	if where != "" {
		query += fmt.Sprintf(" WHERE %s", where)
	}
	
	res, err := s.client.Exec(query)
	if err != nil {
		return 0, false
	}
	rowAffect, err := res.RowsAffected()
	if err != nil {
		return 0, false
	}
	return rowAffect, true
}

func (s *MySQL) Delete(table string, where string) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s", table)
	
	if where != "" {
		query += fmt.Sprintf(" WHERE %s", where)
	}
	
	if res, err := s.client.Exec(query); err != nil {
		return 0, nil
	} else {
		if rowAffect, err := res.RowsAffected(); err != nil {
			return 0, err
		} else {
			return rowAffect, nil
		}
	}
}

func (s *MySQL) GetDeviceInformationCount(deviceID uint64, infoType string) uint64 {
	query := fmt.Sprintf("SELECT COUNT(*) as cn from device_information_transmissions WHERE deviceId=%d AND "+
		"informationType='%s';", deviceID, infoType)
	
	row := s.client.QueryRow(query)
	var info InfoCount
	_ = row.Scan(&info.cn)
	return info.cn
}

func (s *MySQL) InsertDeviceInformation(deviceID uint64, infoType string, content string) (int64, bool) {
	d := utils.GetDateWithFormat()
	return s.insertQuery(fmt.Sprintf("INSERT INTO device_information_transmissions "+
		"(deviceId, informationType, content, created_at, updated_at) VALUES "+
		"(%d, '%s', '%s', '%s', '%s');", deviceID, infoType, content, d, d))
}

func (s *MySQL) UpdateDeviceInformation(deviceID uint64, infoType string, content string) bool {
	query := fmt.Sprintf("UPDATE device_information_transmissions SET content='%s', updated_at='%s' "+
		"WHERE deviceId=%d AND informationType='%s';", content, utils.GetDateWithFormat(), deviceID, infoType)
	row, _ := s.client.Prepare(query)
	res, _ := row.Exec()
	ar, _ := res.RowsAffected()
	defer row.Close()
	return ar == 1
}

func (s *MySQL) InsertAlarm(deviceID uint64, latitude, longitude, speed, gpsInformation, course, mcc, mnc, lac, cellID, terminalInformation,
	voltageLevel, gsmSignal, alarmMode, language, date, serialNumber string) (int64, bool) {
	d := utils.GetDateWithFormat()
	return s.insertQuery(fmt.Sprintf("INSERT INTO device_alarms "+
		"(deviceId, latitude, longitude, speed, gpsInformation, course, mcc, mnc, lac, cellId, terminalInformation, "+
		"voltageLevel, gsmSignal, alarmMode, language, date, serialNumber, created_at, updated_at) VALUES "+
		"(%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');",
		deviceID, latitude, longitude, speed, gpsInformation, course, mcc, mnc, lac, cellID, terminalInformation,
		voltageLevel, gsmSignal, alarmMode, language, date, serialNumber, d, d))
}

func (s *MySQL) GetDeviceInfoSyncRedis() []SyncInfo {
	query := "SELECT devices.id, devices.userId, devices.deviceImei, device_models.deviceProtocol, devices.expireDate, device_models.modelName, " +
		"device_manages.expireDate AS expireTime FROM devices LEFT JOIN device_models ON devices.deviceModel = device_models.id " +
		"LEFT JOIN device_manages ON devices.deviceImei = device_manages.deviceImei;"
	rows, _ := s.client.Query(query)
	
	defer rows.Close()
	
	var result []SyncInfo
	
	for rows.Next() {
		var tmp SyncInfo
		
		err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.DeviceImei, &tmp.Model, &tmp.ExpireDate, &tmp.ModelName, &tmp.ExpireTime)
		
		if err != nil {
			panic(err.Error())
		}
		
		result = append(result, tmp)
	}
	return result
}

func (s *MySQL) GetAlarmSetting(deviceID uint64, alarmType string) (string, string, uint) {
	query := fmt.Sprintf("SELECT alarms.title, alarms.en_title, device_alarm_settings.notif_set from alarms "+
		"LEFT JOIN device_alarm_settings ON alarms.id = device_alarm_settings.alarm_id WHERE "+
		"alarms.mode='%s' AND device_alarm_settings.device_id=%d;", alarmType, deviceID)
	
	row := s.client.QueryRow(query)
	var title string
	var notifySet uint
	var enTitle string
	_ = row.Scan(&title, &enTitle, &notifySet)
	return title, enTitle, notifySet
}

func (s *MySQL) GetDeviceInfo(deviceID uint64) (sql.NullInt32, string) {
	query := fmt.Sprintf("SELECT userId, deviceName FROM devices WHERE id=%d;", deviceID)
	
	row := s.client.QueryRow(query)
	var userId sql.NullInt32
	var deviceName string
	_ = row.Scan(&userId, &deviceName)
	return userId, deviceName
}

func (s *MySQL) GetUserPushTokens(userID uint64) []PushToken {
	query := fmt.Sprintf("SELECT pushToken FROM user_push_tokens WHERE userId=%d;", userID)
	rows, _ := s.client.Query(query)
	
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	
	var result []PushToken
	
	for rows.Next() {
		var tmp PushToken
		
		err := rows.Scan(&tmp.PushToken)
		
		if err != nil {
			panic(err.Error())
		}
		
		result = append(result, tmp)
	}
	return result
}

func (s *MySQL) GetUserLanguage(userID uint64) string {
	query := fmt.Sprintf("SELECT lang FROM users WHERE id=%d;", userID)
	row := s.client.QueryRow(query)
	var language string
	_ = row.Scan(&language)
	return language
}
