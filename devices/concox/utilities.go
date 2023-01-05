package concox

import (
	"bytes"

	"github.com/snksoft/crc"
	"golang.org/x/exp/slices"
)

var serialNumber uint16 = 0

const (
	ParserLogin                        = 0x01
	ParserLocationX3                   = 0x22
	ParserLocationGT06                 = 0x12
	ParserStatus                       = 0x13
	ParserAlarmX3                      = 0x26
	ParserAlarmX3V2                    = 0x27
	ParserAlarmGT06                    = 0x16
	ParserWifiInformation              = 0x2c
	ParserLBSLocationX3                = 0x28
	ParserLBSLocationGT06              = 0x18
	ParserInformation                  = 0x94
	ParserTimeCalibration              = 0x8a
	ParserOnlineCommand                = 0xaf
	ParserOnlineCommandResponse        = 0x15
	ParserOnlineCommandLongResponse    = 0x21
	ParserOnlineCommandRequestStartBit = 0x70
	ParserStartBit                     = 0x78
	ParserLongStartBit                 = 0x79
	ParserEndBitFirst                  = 0x0d
	ParserEndBitEnd                    = 0x0a
	ParserTimeCalibrationLength        = 0x0b
	ParserResponseLength               = 0x05
	ParserOnlineCommandProtocol        = 0x80
)

const (
	ParserLoginStr                     = "Login"
	ParserLocationStr                  = "Location"
	ParserStatusStr                    = "Status"
	ParserAlarmStr                     = "Alarm"
	ParserWifiInformationStr           = "Wifi Information"
	ParserLBSLocationStr               = "LBS Location"
	ParserInformationStr               = "Information"
	ParserTimeCalibrationStr           = "Time Calibration"
	ParserOnlineCommandResponseStr     = "Online Command Response"
	ParserOnlineCommandLongResponseStr = "Online Command Response Long"
)

func validPackage() []byte {
	return []byte{
		ParserLogin,
		ParserStatus,
		ParserLocationGT06,
		ParserAlarmGT06,
		ParserTimeCalibration,
		ParserOnlineCommandResponse,
		ParserLBSLocationGT06,
		ParserInformation,
		ParserLocationX3,
		ParserOnlineCommandLongResponse,
		ParserAlarmX3,
		ParserLBSLocationX3,
	}
}

func GeneratePackageSN() (byte, byte) {
	if serialNumber++; serialNumber > 65535 {
		serialNumber = 1
	}
	return byte(serialNumber >> 8), byte(serialNumber)
}

func GenerateCrc(data []byte) (byte, byte) {
	ccittCrc := crc.CalculateCRC(crc.X25, data)
	return byte(ccittCrc >> 8), byte(ccittCrc)
}

func CrcChecker(data []byte) bool {
	crack := data[2 : len(data)-4]
	dataC := data[len(data)-4 : len(data)-2]
	crc1, crc2 := GenerateCrc(crack)
	return bytes.Equal(dataC, []byte{crc1, crc2})
}

func GetStartBytes(input []byte) []byte {
	return input[:2]
}

func GetEndBytes(input []byte) []byte {
	return input[len(input)-2:]
}

func ValidateStartBytes(input []byte) bool {
	return IsNormalPackage(input) || IsLongPackage(input) || IsOnlineCommandRequest(input)
}

func IsNormalPackage(input []byte) bool {
	return bytes.Equal(GetStartBytes(input), []byte{ParserStartBit, ParserStartBit})
}

func IsLongPackage(input []byte) bool {
	return bytes.Equal(GetStartBytes(input), []byte{ParserLongStartBit, ParserLongStartBit})
}

func IsOnlineCommandRequest(input []byte) bool {
	return bytes.Equal(GetStartBytes(input), []byte{
		ParserOnlineCommandRequestStartBit,
		ParserOnlineCommandRequestStartBit,
	})
}

func IsPacketFromDevice(input []byte) bool {
	return (IsNormalPackage(input) || IsLongPackage(input)) && !IsOnlineCommandRequest(input)
}

func IsPacketFromSite(input []byte) bool {
	return IsOnlineCommandRequest(input) && (!IsNormalPackage(input) && !IsLongPackage(input))
}

func ValidateEndBytes(input []byte) bool {
	return bytes.Equal(GetEndBytes(input), []byte{ParserEndBitFirst, ParserEndBitEnd})
}

func ValidatePackage(input []byte) bool {
	return ValidateStartBytes(input) && ValidateEndBytes(input) &&
		slices.Contains(validPackage(), GetPackageType(input)) && CrcChecker(input)
}

func GetPackageType(input []byte) byte {
	if IsLongPackage(GetStartBytes(input)) {
		return input[4]
	}
	
	return input[3]
}

func CreatePackageForDevice(message string) []byte {
	result := make([]byte, len(message)+15)
	msg := []byte(message)
	
	result[0] = ParserStartBit
	result[1] = ParserStartBit
	result[2] = byte(len(message) + 10)
	result[3] = ParserOnlineCommandProtocol
	result[len(result)-6], result[len(result)-5] = GeneratePackageSN()
	result[4] = byte(len(message) + 4)
	
	copy(result[9:9+len(msg)], msg)
	
	result[len(result)-4], result[len(result)-3] = GenerateCrc(result[2 : len(result)-4])
	
	result[len(result)-2] = ParserEndBitFirst
	result[len(result)-1] = ParserEndBitEnd
	
	return result
}

func IsLogin(input []byte) bool {
	return GetPackageType(input) == ParserLogin
}
