package main

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	user = "j.nelson34"
	pass = "spanktron"
)

func main() {
	body := strings.NewReader(fmt.Sprintf(`auth_user=%v&auth_pass=%v&accept=Continue`, user, pass))
	req, err := http.NewRequest("POST", "https://studentwireless.clark.edu:8003/index.php?zone=stuwrls", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}
