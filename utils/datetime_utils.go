package utils

// ToExtendedFormat Convert datetime formatted in ISO 8601 basic format to extended format
func ToExtendedFormat(datetime string) string {
	if len(datetime) != 13 {
		return datetime
	}

	year := datetime[0:4]
	month := datetime[4:6]
	day := datetime[6:8]
	hour := datetime[9:11]
	minute := datetime[11:13]

	return year + "-" + month + "-" + day + "T" + hour + ":" + minute
}
