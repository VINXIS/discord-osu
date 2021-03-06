package tools

import "time"

import "strings"

// TimeParse tries to parses a date / time
func TimeParse(datetime string) (timestamp time.Time, err error) {
	timeFormats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	customDate := []string{
		"january 2 2006",
		"jan 2 2006",
		"1 2 2006",
		"01 2 2006",
		"1 02 2006",
		"01 02 2006",

		"2 january 2006",
		"2 jan 2006",
		"2 1 2006",
		"2 01 2006",
		"02 1 2006",
		"02 01 2006",

		"2006 january 2",
		"2006 jan 2",
		"2006 1 2",
		"2006 01 2",
		"2006 1 02",
		"2006 01 02",
	}

	customTime := []string{
		"15:04:05",
		"3:04 PM",
		"03:04 PM",
		"3:04PM",
		"03:04PM",
		"3 PM",
		"03 PM",
	}

	customZone := []string{
		"MST",
		"-0700",
		"-07",
		"-07:00",
		"-7:00",
	}

	for _, timeFormat := range timeFormats {
		timestamp, err = time.Parse(timeFormat, datetime)
		if err == nil {
			return timestamp, nil
		}
	}

	// Run custom formats only if none of the default formats work
	for _, date := range customDate {
		timestamp, err = time.Parse(date, datetime)
		if err == nil {
			return timestamp, nil
		}

		for _, timer := range customTime {
			timestamp, err = time.Parse(timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(date+" "+timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(timer+" "+date, datetime)
			if err == nil {
				return timestamp, nil
			}

			for _, zone := range customZone {
				timestamp, err = time.Parse(date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}
			}
		}
	}

	// Run with dashed date now if none of the non-dashed work
	for _, date := range customDate {
		date = dashed(date)

		timestamp, err = time.Parse(date, datetime)
		if err == nil {
			return timestamp, nil
		}

		for _, timer := range customTime {
			timestamp, err = time.Parse(timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(date+" "+timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(timer+" "+date, datetime)
			if err == nil {
				return timestamp, nil
			}

			for _, zone := range customZone {
				timestamp, err = time.Parse(date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}
			}
		}
	}

	return timestamp, err
}

func dashed(date string) string {
	return strings.Replace(date, " ", "-", -1)
}
