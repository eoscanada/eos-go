package boot

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/abourget/llerrgroup"
)

func (b *Boot) downloadReferences() error {
	if err := b.ensureCacheExists(); err != nil {
		return fmt.Errorf("error creating cache path: %s", err)
	}

	eg := llerrgroup.New(10)
	for _, contentRef := range b.bootSequence.Contents {
		if eg.Stop() {
			continue
		}

		contentRef := contentRef
		eg.Go(func() error {
			if err := b.DownloadURL(contentRef.URL, contentRef.Hash); err != nil {
				return fmt.Errorf("content %q: %s", contentRef.Name, err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (b *Boot) ensureCacheExists() error {
	return os.MkdirAll(b.cachePath, 0777)
}

func (b *Boot) DownloadURL(ref string, hash string) error {
	if hash != "" && b.isInCache(ref) {
		return nil
	}

	cnt, err := b.downloadRef(ref)
	if err != nil {
		return err
	}

	if hash != "" {
		h := sha256.New()
		_, _ = h.Write(cnt)
		contentHash := hex.EncodeToString(h.Sum(nil))

		if contentHash != hash {
			return fmt.Errorf("hash in boot sequence [%q] not equal to computed hash on downloaded file [%q]", hash, contentHash)
		}
	}

	zlog.Info("Caching content.", zap.String("ref", ref))
	if err := b.writeToCache(ref, cnt); err != nil {
		return err
	}

	return nil
}

func (b *Boot) downloadRef(ref string) ([]byte, error) {
	zlog.Info("Downloading content", zap.String("from", ref))
	if _, err := os.Stat(ref); err == nil {
		return b.downloadLocalFile(ref)
	}

	destURL, err := url.Parse(ref)
	if err != nil {
		return nil, fmt.Errorf("ref %q is not a valid URL: %s", ref, err)
	}

	switch destURL.Scheme {
	case "file":
		return b.downloadFileURL(destURL)
	case "http", "https":
		return b.downloadHTTPURL(destURL)
	default:
		return nil, fmt.Errorf("don't know how to handle scheme %q (from ref %q)", destURL.Scheme, destURL)
	}
}

func (b *Boot) downloadLocalFile(ref string) ([]byte, error) {
	return ioutil.ReadFile(ref)
}

func (b *Boot) downloadFileURL(destURL *url.URL) ([]byte, error) {
	fmt.Printf("Path %s, Raw path: %s\n", destURL.Path, destURL.RawPath)
	return []byte{}, nil
}

func (b *Boot) downloadHTTPURL(destURL *url.URL) ([]byte, error) {
	req, err := http.NewRequest("GET", destURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("download attempts failed")
	}
	defer resp.Body.Close()

	cnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		if len(cnt) > 50 {
			cnt = cnt[:50]
		}
		return nil, fmt.Errorf("couldn't get %s, return code: %d, server error: %q", destURL, resp.StatusCode, cnt)
	}

	return cnt, nil
}

func (b *Boot) writeToCache(ref string, content []byte) error {
	fileName := replaceAllWeirdities(ref)
	return ioutil.WriteFile(filepath.Join(b.cachePath, fileName), content, 0666)
}

func (b *Boot) isInCache(ref string) bool {
	fileName := filepath.Join(b.cachePath, replaceAllWeirdities(ref))

	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}

func (b *Boot) ReadFromCache(ref string) ([]byte, error) {
	fileName := replaceAllWeirdities(ref)
	return ioutil.ReadFile(filepath.Join(b.cachePath, fileName))
}

func (b *Boot) ReaderFromCache(ref string) (io.ReadCloser, error) {
	fileName := replaceAllWeirdities(ref)
	return os.Open(filepath.Join(b.cachePath, fileName))
}

func (b *Boot) FileNameFromCache(ref string) string {
	fileName := replaceAllWeirdities(ref)
	return filepath.Join(b.cachePath, fileName)
}
