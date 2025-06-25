package service

import (
	"GonIO/internal/domain"
	"errors"
	"fmt"
	"strconv"
)

// Validate funcion checks all requirements for bucket name
func Validate(name string) error {
	countDots, err := CheckName(name)
	if err != nil {
		return err
	}
	if countDots == 3 {
		if CheckIp(name) {
			return domain.ErrNameLikeIpAdress
		}
	}
	if CheckConsecutive(name) {
		return domain.ErrNamePeriodOrDash
	}
	return nil
}

func CheckName(s string) (int, error) {
	countDots := 0
	if len(s) < 3 || len(s) > 63 {
		return 0, domain.ErrNameLenght
	}
	if s[0] == '-' || s[len(s)-1] == '-' {
		return 0, domain.ErrNameHyphen
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '_' || s[i] == '`' {
			continue
		}

		if s[i] == '.' {
			countDots++
			continue
		}

		if s[i] >= '0' && s[i] <= '9' {
			continue
		}

		if s[i] >= 'A' && s[i] <= 'Z' {
			continue
		}
		if s[i] >= 'a' && s[i] <= 'z' {
			continue
		}
		return 0, errors.New("name error character: " + string(s[i]))
	}
	return countDots, nil
}

func CheckConsecutive(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 45 {
			if i+1 < len(s) {
				if s[i+1] == 45 {
					return true
				}
			}
		}
		if s[i] == 46 {
			if i+1 < len(s) {
				if s[i+1] == 46 {
					return true
				}
			}
		}
	}
	return false
}

func CheckIp(s string) bool {
	if s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	start := 0
	for i := 0; i < len(s); i++ {
		if i == len(s)-1 {
			num, err := strconv.Atoi(string(s[start:]))
			fmt.Println(s[start:])
			if err != nil {
				return false
			}
			if num < 0 || num > 256 {
				return false
			}
		}
		if s[i] == '.' {
			num, err := strconv.Atoi(string(s[start:i]))
			if err != nil {
				return false
			}
			if num < 0 || num > 256 {
				return false
			}
			if i+1 < len(s) {
				start = i + 1
			}
		}

	}
	return true
}
