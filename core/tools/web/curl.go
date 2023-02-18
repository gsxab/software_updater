package web

import (
	"context"
	"github.com/tebeka/selenium"
	"io"
	"net/url"
	"os/exec"
)

func SeleniumCookiesToHeader(url *url.URL, cookies []selenium.Cookie) (result []string) {
	for _, selCookie := range cookies {
		if selCookie.Domain != url.Hostname() {
			continue
		}
		result = append(result, "-H", selCookie.Name+"="+selCookie.Value)
	}
	return
}

func CURL(ctx context.Context, url *url.URL, cookies []selenium.Cookie, output io.Writer) error {
	curlCookies := SeleniumCookiesToHeader(url, cookies)
	args := []string{"-s", "-L"}
	args = append(args, curlCookies...)
	args = append(args, url.String())
	cmd := exec.CommandContext(ctx, "curl", args...)
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()
	err = cmd.Start()
	if err != nil {
		return err
	}
	_, err = io.Copy(output, stream)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
