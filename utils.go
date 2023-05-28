package utils

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/xen0tic/utils/devices/concox"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ConstDefaultDateFormat = "2006-01-02 15:04:05"
	RedisDevices           = "devices"
	RedisLocations         = "locations"
	RedisOnlines           = "onlines"
	RedisHeartBeat         = "heart_beats"
	RedisX3InGt            = "X3InGt"
	RedisGtInX3            = "GtInX3"
	RedisLocationList      = "locList"
	RedisLbsLocation       = "lbs_locations"
)

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}

func setOutput(ws zapcore.WriteSyncer, conf zap.Config) zap.Option {
	var enc zapcore.Encoder
	switch conf.Encoding {
	case "json":
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	case "console":
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	default:
		panic("unknown encoding")
	}

	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(enc, ws, conf.Level)
	})
}

func getWriteSyncer(logName string) zapcore.WriteSyncer {
	var ioWriter = &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    100, // MB
		MaxBackups: 3,   // number of backups
		MaxAge:     28,  // days
		LocalTime:  true,
		Compress:   true, // disabled by default
	}
	var sw = WriteSyncer{
		ioWriter,
	}
	return sw
}

func InitLogger(filePath string) (*zap.Logger, error) {
	var cfg zap.Config

	pwd, _ := os.Getwd()
	directory := fmt.Sprintf("%s/log", pwd)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		_ = os.Mkdir(directory, 0777)
	}

	filePath = fmt.Sprintf("%s/%s", directory, filePath)

	cfg = zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Mon, 2006-01-02 03:04:05 MST")
	cfg.OutputPaths = []string{filePath}
	sw := getWriteSyncer(filePath)

	return cfg.Build(setOutput(sw, cfg))
}

func FailOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func GenerateIsoDate(dateTime string) time.Time {
	l, _ := time.LoadLocation("UTC")
	tehran, _ := time.LoadLocation("Asia/Tehran")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, tehran)
	return t.In(l)
}

func GetDateDiff(t time.Time) int64 {
	t1 := GetLocalizedTime()
	return t.Sub(t1).Milliseconds()
}

func GetLocalizedTime() time.Time {
	location, err := GetIranLocation()
	if err != nil {
		return time.Time{}
	}
	return time.Now().In(location)
}

func GetIranLocation() (*time.Location, error) {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func GetDateWithFormat(format ...string) string {
	now := GetLocalizedTime()
	if len(format) == 1 {
		return now.Format(format[0])
	}
	return now.Format(ConstDefaultDateFormat)
}

func ConvertByteToHex(input []byte) []string {
	arr := make([]string, 0)
	for _, item := range input {
		arr = append(arr, fmt.Sprintf("%02x", item))
	}
	return arr
}

func ConvertByteToString(input []byte) string {
	result := ""
	for _, item := range ConvertByteToHex(input) {
		result += item
	}
	return result
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func IsConcox(packet []byte) bool {
	return (concox.IsNormalPackage(packet) || concox.IsLongPackage(packet)) && !concox.IsOnlineCommandRequest(packet)
}

func IsOnlineRequest(packet []byte) bool {
	return concox.IsOnlineCommandRequest(packet) && (!concox.IsNormalPackage(packet) && !concox.IsLongPackage(packet))
}

func IsLogin(packet []byte) bool {
	return concox.GetPackageType(packet) == concox.ParserLogin
}

func GetDate(format ...string) (string, error) {
	mode := format[0]
	dFormat := ""
	if len(format) > 2 || len(format) == 0 {
		return "", errors.New("input parameter not valid")
	}
	if len(format) == 1 {
		dFormat = ConstDefaultDateFormat
	} else {
		dFormat = format[1]
	}
	return AddDate(mode).Format(dFormat), nil
}

func AddDate(mode string) time.Time {
	dateMod := strings.Split(mode, "_")
	now := GetLocalizedTime()
	value, _ := strconv.Atoi(dateMod[0])
	switch dateMod[1] {
	case "year", "years":
		now = now.AddDate(value, 0, 0)
	case "month", "months":
		now = now.AddDate(0, value, 0)
	case "day", "days":
		now = now.AddDate(0, 0, value)
	}
	return now
}

func ConvertHexToByte(input []string) []byte {
	arr1 := make([]byte, 0)
	for _, item := range input {
		val, _ := strconv.ParseInt(item, 16, 64)
		arr1 = append(arr1, byte(val))
	}
	return arr1
}

func FormatDate(time time.Time) string {
	return time.Format(ConstDefaultDateFormat)
}

func GetDistance(p1 []float64, p2 []float64) float64 {
	R := 6378137
	dLat := Rad(p2[0] - p1[0])
	dLong := Rad(p2[1] - p1[1])
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(Rad(p1[0]))*math.Cos(Rad(p2[0]))*math.Sin(dLong/2)*math.Sin(dLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float64(R) * c
}

func Rad(x float64) float64 {
	return (x * math.Pi) / 180
}

func GetRemoteAddr(remoteAddr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(remoteAddr, ":", ""), ".", "")
}

func GetPackageSn(input []byte) string {
	arr := ConvertByteToHex(input)
	result := ""
	tmp := arr[12:14]
	for _, item := range tmp {
		result += item
	}
	return result
}

func GetPackageSnInNumber(input []byte) int {
	val, _ := strconv.Atoi(GetPackageSn(input))
	return val
}

func SplitPackage(input []byte, array [][]byte) [][]byte {

	if len(input) < 10 {
		return array
	}

	if !concox.IsNormalPackage(input) && !concox.IsLongPackage(input) {
		return array
	}

	iLen := int(input[2])
	if concox.IsLongPackage(input) {
		fd, _ := strconv.ParseInt(ConvertByteToString(input[2:4]), 16, 64)
		iLen = int(fd) + 1
	}

	iLen += 5

	if len(input) == iLen {
		array = append(array, input)
	} else if len(input) > iLen {
		if concox.ValidatePackage(input[:iLen]) {
			array = append(array, input[:iLen])
		}
		return SplitPackage(input[iLen:], array)
	}

	return array
}

func SortArray(array [][]byte) [][]byte {

	sort.Slice(array, func(i, j int) bool {
		a := GetPackageSnInNumber(array[i])
		b := GetPackageSnInNumber(array[j])
		return a < b
	})

	return array
}

func DecodeLocation(latitude float64) float64 {
	return ToFixed(latitude/60.0/30000.0, 6)
}

func GetDeviceImei(input []byte) string {
	result := ""
	tmp := input[4:12]
	for _, item := range tmp {
		result = result + fmt.Sprintf("%02x", int(item))
	}
	return LeftPad2Len(result[1:], "0", 15)
}
