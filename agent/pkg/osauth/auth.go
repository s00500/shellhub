package osauth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/sirupsen/logrus"
)

var DefaultShadowFilename = "/etc/shadow"

type ShadowEntry struct {
	Username    string
	Password    string
	Lastchanged int
	Minimum     int
	Maximum     int
	Warn        int // The number of days before password is to expire that user is warned that his/her password must be changed
	Inactive    int // The number of days after password expires that account is disabled
	Expire      int
}

func AuthUser(username, passwd string) bool {
	shadowFile, err := os.Open(DefaultShadowFilename)
	if err != nil {
		// TODO: log error
		logrus.Error("Could not open /etc/shadow")
		return false
	}

	defer shadowFile.Close()

	entries, err := parseShadowReader(shadowFile)
	if err != nil {
		logrus.Println("Could not parse shadowfile %v", err)
		return false
	}

	if entry, ok := entries[username]; ok {
		if entry.Password != "" {
			crypt := crypt.NewFromHash(entry.Password)
			if crypt == nil {
				logrus.Error("Could not detect password crypto algorithm from shadow entry")
				return false
			}

			err := crypt.Verify(entry.Password, []byte(passwd))
			return err == nil
		}
		logrus.Error("Password entry is empty")
		return false
	}

	logrus.Warn("User not found")
	return false
}

func parseShadowReader(r io.Reader) (map[string]ShadowEntry, error) {
	lines := bufio.NewReader(r)
	entries := make(map[string]ShadowEntry)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			break
		}

		if len(line) == 0 || strings.HasPrefix(string(line), "#") {
			continue
		}

		entry, err := parseShadowLine(string(line))
		if err != nil {
			return nil, err
		}

		entries[entry.Username] = entry
	}
	return entries, nil
}

func parseShadowLine(line string) (ShadowEntry, error) {
	result := ShadowEntry{}
	parts := strings.Split(strings.TrimSpace(line), ":")
	if len(parts) != 9 {
		return result, fmt.Errorf("Shadow line had wrong number of parts %d != 9", len(parts))
	}
	result.Username = strings.TrimSpace(parts[0])
	result.Password = strings.TrimSpace(parts[1])
	/*
	   	lastchanged, err := strconv.Atoi(parts[2])
	   	if err != nil {
	   		return result, fmt.Errorf("Shadow line had badly formatted lastchanged %s", parts[2])
	   	}
	   	result.Lastchanged = lastchanged

	   	if parts[3] != "" {
	   	min, err := strconv.Atoi(strings.TrimSpace(parts[3]))
	   	if err != nil {
	   		return result, fmt.Errorf("Shadow line had badly formatted min %s", parts[3])
	   	}
	   	result.Minimum = min
	   }

	   if parts[4] != "" {
	   	max, err := strconv.Atoi(strings.TrimSpace(parts[4]))
	   	if err != nil {
	   		return result, fmt.Errorf("Shadow line had badly formatted max %s", parts[4])
	   	}
	   	result.Maximum = max
	   }

	   	if parts[5] != "" {
	   	warn, err := strconv.Atoi(strings.TrimSpace(parts[5]))
	   	if err != nil {
	   		return result, fmt.Errorf("Shadow line had badly formatted warn %s", parts[5])
	   	}
	   	result.Warn = warn
	   }

	   	if parts[6] != "" {
	   		inactive, err := strconv.Atoi(strings.TrimSpace(parts[6]))
	   		if err != nil {
	   			return result, fmt.Errorf("Shadow line had badly formatted inactive %s", parts[6])
	   		}
	   		result.Inactive = inactive
	   	}

	   	if parts[7] != "" {
	   	expire, err := strconv.Atoi(strings.TrimSpace(parts[7]))
	   	if err != nil {
	   		return result, fmt.Errorf("Shadow line had badly formatted expire %s", parts[7])
	   	}
	   	result.Expire = expire
	   	}
	*/

	return result, nil
}
