package barred

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ReadBarredFile - read the barred number file into an array
func ReadBarredFile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("error opening file: %s", fileName)
		return nil, err
	}
	defer f.Close()
	var barred []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		barred = append(barred, s)
	}
	return barred, nil
}

// CheckBarred - check if telephone is barred
func CheckBarred(telephone string, barred []string) bool {
	// NOTE: telephone includes country code i.e. 07123456789 => 447123456789
	// barred telephone list does not include country code but starts with 07...
	tel := strings.Replace(telephone, "447", "07", 1)
	for i := 0; i < len(barred); i++ {
		if strings.HasPrefix(tel, barred[i]) {
			return true
		}
	}
	return false
}
