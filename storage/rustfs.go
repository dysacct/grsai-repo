package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RustFSConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Bucket        string
	Region        string
	Prefix        string
	PublicBaseURL string
	MaxBytes      int64
	PublicRead    bool
}

type RustFSClient struct {
	cfg           RustFSConfig
	httpClient    *http.Client
	ensureOnce    sync.Once
	ensureErr     error
	endpointHost  string
	endpointPath  string
	endpointProto string
}

type StoredImage struct {
	URL         string
	ObjectKey   string
	ContentType string
	Bytes       []byte
}

func NewRustFSClient(cfg RustFSConfig) (*RustFSClient, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "lyncr"
	}
	if cfg.MaxBytes <= 0 {
		cfg.MaxBytes = 50 * 1024 * 1024
	}
	endpoint, err := url.Parse(strings.TrimRight(cfg.Endpoint, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse rustfs endpoint: %w", err)
	}
	if endpoint.Scheme == "" || endpoint.Host == "" {
		return nil, fmt.Errorf("rustfs endpoint must include scheme and host")
	}
	return &RustFSClient{
		cfg:           cfg,
		httpClient:    &http.Client{Timeout: 180 * time.Second},
		endpointHost:  endpoint.Host,
		endpointPath:  strings.TrimRight(endpoint.Path, "/"),
		endpointProto: endpoint.Scheme,
	}, nil
}

func (c *RustFSClient) StoreImageURL(ctx context.Context, imageURL string) (*StoredImage, error) {
	body, contentType, err := c.download(ctx, imageURL)
	if err != nil {
		return nil, err
	}
	return c.StoreImageBytes(ctx, body, contentType)
}

func (c *RustFSClient) StoreImageBytes(ctx context.Context, body []byte, contentType string) (*StoredImage, error) {
	c.ensureOnce.Do(func() {
		c.ensureErr = c.ensureBucket(ctx)
	})
	if c.ensureErr != nil {
		return nil, c.ensureErr
	}

	if contentType == "" {
		contentType = http.DetectContentType(body)
	}
	key := c.objectKey(contentType)
	if err := c.putObject(ctx, key, body, contentType); err != nil {
		return nil, err
	}
	return &StoredImage{
		URL:         c.publicURL(key),
		ObjectKey:   key,
		ContentType: contentType,
		Bytes:       body,
	}, nil
}

func (c *RustFSClient) download(ctx context.Context, imageURL string) ([]byte, string, error) {
	if imageURL == "" {
		return nil, "", fmt.Errorf("image URL is empty")
	}
	if strings.HasPrefix(imageURL, "data:") {
		return nil, "", fmt.Errorf("data URLs are not supported by RustFS downloader")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create image download request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("download image: status %d", resp.StatusCode)
	}

	limit := c.cfg.MaxBytes + 1
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return nil, "", fmt.Errorf("read image body: %w", err)
	}
	if int64(len(body)) > c.cfg.MaxBytes {
		return nil, "", fmt.Errorf("image exceeds max size of %d bytes", c.cfg.MaxBytes)
	}

	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = http.DetectContentType(body)
	}
	return body, contentType, nil
}

func (c *RustFSClient) ensureBucket(ctx context.Context) error {
	if err := c.ensureBucketExists(ctx); err != nil {
		return err
	}
	if c.cfg.PublicRead {
		return c.ensurePublicReadPolicy(ctx)
	}
	return nil
}

func (c *RustFSClient) ensureBucketExists(ctx context.Context) error {
	headURL := c.bucketURL("")
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, headURL, nil)
	if err != nil {
		return fmt.Errorf("create bucket head request: %w", err)
	}
	c.sign(req, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("check rustfs bucket: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == http.StatusForbidden {
		// Some S3-compatible stores allow object writes but deny bucket metadata checks.
		// Treat 403 as "try object upload" so existing buckets still work.
		return nil
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("check rustfs bucket: status %d", resp.StatusCode)
	}

	putReq, err := http.NewRequestWithContext(ctx, http.MethodPut, headURL, nil)
	if err != nil {
		return fmt.Errorf("create bucket put request: %w", err)
	}
	c.sign(putReq, nil)
	putResp, err := c.httpClient.Do(putReq)
	if err != nil {
		return fmt.Errorf("create rustfs bucket: %w", err)
	}
	defer putResp.Body.Close()
	if putResp.StatusCode >= 200 && putResp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("create rustfs bucket: status %d", putResp.StatusCode)
}

func (c *RustFSClient) ensurePublicReadPolicy(ctx context.Context) error {
	policy := bucketPolicy{
		Version: "2012-10-17",
		Statement: []bucketPolicyStatement{
			{
				Effect:    "Allow",
				Principal: "*",
				Action:    []string{"s3:GetObject"},
				Resource:  []string{fmt.Sprintf("arn:aws:s3:::%s/*", c.cfg.Bucket)},
			},
		},
	}
	body, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("marshal public read bucket policy: %w", err)
	}
	policyURL := c.bucketURL("")
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, policyURL+"?policy=", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create bucket policy put request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.sign(req, body)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("set rustfs public read bucket policy: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("set rustfs public read bucket policy: status %d: %s", resp.StatusCode, string(errBody))
	}
	return nil
}

type bucketPolicy struct {
	Version   string                  `json:"Version"`
	Statement []bucketPolicyStatement `json:"Statement"`
}

type bucketPolicyStatement struct {
	Effect    string   `json:"Effect"`
	Principal string   `json:"Principal"`
	Action    []string `json:"Action"`
	Resource  []string `json:"Resource"`
}

func (c *RustFSClient) putObject(ctx context.Context, key string, body []byte, contentType string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.bucketURL(key), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create object put request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	c.sign(req, body)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload image to rustfs: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("upload image to rustfs: status %d: %s", resp.StatusCode, string(errBody))
	}
	return nil
}

func (c *RustFSClient) objectKey(contentType string) string {
	ext := extensionForContentType(contentType)
	now := time.Now().UTC()
	randomUUID := uuid.New().String()
	filename := fmt.Sprintf("%d_%s%s", now.UnixNano(), randomUUID, ext)
	key := fmt.Sprintf("%s/%s/%s", strings.Trim(c.cfg.Prefix, "/"), now.Format("2006/01/02"), filename)
	return strings.TrimLeft(path.Clean(key), "/")
}

func extensionForContentType(contentType string) string {
	if contentType == "" {
		return ".bin"
	}
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}

func (c *RustFSClient) bucketURL(key string) string {
	u := url.URL{
		Scheme: c.endpointProto,
		Host:   c.endpointHost,
		Path:   path.Join(c.endpointPath, c.cfg.Bucket, key),
	}
	if key == "" {
		u.Path = path.Join(c.endpointPath, c.cfg.Bucket)
	}
	return u.String()
}

func (c *RustFSClient) publicURL(key string) string {
	base := c.cfg.PublicBaseURL
	if base == "" {
		base = c.cfg.Endpoint
	}
	base = strings.TrimRight(base, "/")
	return base + "/" + strings.Trim(strings.Trim(c.cfg.Bucket, "/")+"/"+strings.TrimLeft(key, "/"), "/")
}

func (c *RustFSClient) sign(req *http.Request, body []byte) {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	payloadHash := sha256Hex(body)

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	canonicalURI := req.URL.EscapedPath()
	canonicalQuery := canonicalQueryString(req.URL.Query())
	signedHeaders := []string{"host", "x-amz-content-sha256", "x-amz-date"}
	canonicalHeaders := "host:" + req.URL.Host + "\n" +
		"x-amz-content-sha256:" + payloadHash + "\n" +
		"x-amz-date:" + amzDate + "\n"
	if req.Header.Get("Content-Type") != "" {
		signedHeaders = append([]string{"content-type"}, signedHeaders...)
		canonicalHeaders = "content-type:" + req.Header.Get("Content-Type") + "\n" + canonicalHeaders
	}

	scope := dateStamp + "/" + c.cfg.Region + "/s3/aws4_request"
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		strings.Join(signedHeaders, ";"),
		payloadHash,
	}, "\n")
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")
	signingKey := signingKey(c.cfg.SecretKey, dateStamp, c.cfg.Region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))
	credential := c.cfg.AccessKey + "/" + scope
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+credential+", SignedHeaders="+strings.Join(signedHeaders, ";")+", Signature="+signature)
}

func canonicalQueryString(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0)
	for _, key := range keys {
		vals := values[key]
		sort.Strings(vals)
		for _, value := range vals {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}
	return strings.Join(parts, "&")
}

func sha256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func signingKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	return hmacSHA256(kService, []byte("aws4_request"))
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
