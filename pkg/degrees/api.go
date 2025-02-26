package degrees

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sharmayajush/challenge2015/pkg/constants"
)

func getData(url string) ([]byte, error) {
	resp, err := http.Get(constants.Url + url)
	if err != nil {
		return nil, fmt.Errorf("unable to do get for %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d for %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the resp body for %s: %w", url, err)
	}

	return data, nil
}
