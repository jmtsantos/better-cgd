package cgdmeals

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

const (
	// CgdBreakLogin1URL caixabreak login page
	CgdBreakLogin1URL = "https://portalprepagos.cgd.pt/portalprepagos/login.seam?sp_param=PrePago"
	// CgdBreakLogin2URL caixabreak login page
	CgdBreakLogin2URL = "https://portalprepagos.cgd.pt/portalprepagos/auth/forms/login.fcc"
	// CgdBreakMovementsURL caixabreak movements page
	CgdBreakMovementsURL = "https://portalprepagos.cgd.pt/portalprepagos/private/saldoMovimentos.seam"
)

// CGDBreak main
type CGDBreak struct {
	client   *http.Client
	Username string
	Password string
}

// New returns a new client
func New(u, p string) *CGDBreak {
	return &CGDBreak{
		client:   getClient(),
		Username: u,
		Password: p,
	}
}

// CheckBreakBalance logins and checks mealcard balance
func (c *CGDBreak) CheckBreakBalance() (decimal.Decimal, error) {
	var (
		page         *goquery.Document
		balanceParse string
		balance      decimal.Decimal
		err          error
	)

	// Send account data
	log.Println("[observer][cgd][break] login stage 1")
	if err = loginCGD(c.client, c.generateBreakPostDataPreLogin(), CgdBreakLogin1URL, false); err != nil {
		return decimal.Decimal{}, err
	}

	// Send account data
	log.Println("[observer][cgd][break] login stage 2")
	if err = loginCGD(c.client, c.generateBreakPostData(), CgdBreakLogin2URL, false); err != nil {
		return decimal.Decimal{}, err
	}

	log.Println("[observer][cgd][break] getting account movements")
	if page, err = getCGDPage(c.client, CgdBreakMovementsURL, false); err != nil {
		return decimal.Decimal{}, err
	}

	// Parse the data
	page.Find(".valor").Each(func(i int, s *goquery.Selection) {
		balanceParse = strings.Replace(strings.Replace(strings.TrimSpace(s.Text()), "EUR", "", 1), ",", ".", 1)

		if balance, err = decimal.NewFromString(balanceParse); err != nil {
			log.WithError(err).Errorln("[observer][cgd][break] error parsing balance float")
			return
		}
	})

	return balance, err

}

func (c *CGDBreak) generateBreakPostData() url.Values {
	postData := url.Values{}
	postData.Set("target", "/portalprepagos/private/home.seam")
	postData.Add("username", fmt.Sprintf("PPP%s", c.Username))
	postData.Add("userInput", c.Username)
	postData.Add("passwordInput", "*****")
	postData.Add("loginForm:submit", "Entrar")
	postData.Add("password", c.Password)

	return postData
}

func (c *CGDBreak) generateBreakPostDataPreLogin() url.Values {
	postData := url.Values{}
	postData.Set("login_btn_1PPP", "OK")
	postData.Add("USERNAME", c.Username)

	return postData
}
