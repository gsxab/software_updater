package web

import "C"
import (
	"github.com/andelf/go-curl"
	"github.com/tebeka/selenium"
	"io"
	"net/url"
	"software_updater/core/util/error_util"
)

func SeleniumCookiesToHeader(url *url.URL, cookies []selenium.Cookie) (result string) {
	for _, selCookie := range cookies {
		if selCookie.Domain != url.Hostname() {
			continue
		}
		result += selCookie.Name + "=" + selCookie.Value + ";"
	}
	return
}

func CURL(url *url.URL, cookies []selenium.Cookie, output io.Writer) error {
	c := curl.EasyInit()
	defer c.Cleanup()
	errs := error_util.NewCollector()

	errs.Collect(c.Setopt(curl.OPT_URL, url.String()))
	errs.Collect(c.Setopt(curl.OPT_HTTPHEADER, "Cookie: "+SeleniumCookiesToHeader(url, cookies)))
	errs.Collect(c.Setopt(curl.OPT_FOLLOWLOCATION, 1))
	errs.Collect(c.Perform())

	buffer := make([]byte, 1024)
	for {
		cnt, cErr := c.Recv(buffer)
		_, err := output.Write(buffer[:cnt])
		errs.Collect(err)
		if curlErr := cErr.(curl.CurlError); int(curlErr) != int(curl.E_AGAIN) {
			if cErr != nil {
				errs.Collect(cErr)
			}
			break
		}
	}

	return errs.ToError()
}
