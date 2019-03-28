package cgdmeals

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

// loginCGD run CGD account login steps
func loginCGD(client *http.Client, postData url.Values, url string, print bool) error {
	var (
		req *http.Request
		rep *http.Response
		err error
	)

	// Create request
	if req, err = http.NewRequest("POST", url, strings.NewReader(postData.Encode())); err != nil {
		log.WithError(err).Errorln("[observer][cgd] loginCGD - error creating request")
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0")

	// Execute request
	if rep, err = client.Do(req); err != nil {
		log.WithError(err).Errorln("[observer][cgd] loginCGD - unable to reach the server.")
		return err
	}
	defer rep.Body.Close()

	if print {
		body, err := ioutil.ReadAll(rep.Body)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		fmt.Printf("%s\n", body)
	}

	return nil
}

func getCGDPage(client *http.Client, url string, print bool) (*goquery.Document, error) {
	var (
		req *http.Request
		rep *http.Response
		doc *goquery.Document
		err error
	)

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.WithError(err).Errorln("[observer][cgd] getCGDPage - error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0")

	// Execute request
	if rep, err = client.Do(req); err != nil {
		log.WithError(err).Errorln("[observer][cgd] getCGDPage - unable to reach the server.")
		return nil, err
	}
	defer rep.Body.Close()

	if print {
		body, err := ioutil.ReadAll(rep.Body)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
		fmt.Printf("%s\n", body)
	}

	if doc, err = goquery.NewDocumentFromReader(rep.Body); err != nil {
		log.WithError(err).Errorln("[observer][cgd] getCGDPage - unable to create new document from reply")
		return nil, err
	}

	return doc, nil
}

// Setup http client
func getClient() *http.Client {
	jar, _ := cookiejar.New(nil)

	// Setup http client
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives:   true,
			MaxIdleConns:        2,
			MaxIdleConnsPerHost: 2,
		},
		Jar: jar,
	}
}
