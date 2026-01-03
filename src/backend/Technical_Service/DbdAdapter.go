package Technical_Service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var globalDbdAdpterIntance *DbdAdpter

type DbdAdpter struct {
}

func (adapter *DbdAdpter) newDbdAdapter() {

}

var IncomingHeader = []string{
	"Content-Type",
	"Authorization",
	"correlationId",
}

func (adapter *DbdAdpter) StaticGetDbdAdpterInstance() *DbdAdpter {
	if globalAdapterInstance == nil {
		globalDbdAdpterIntance = new(DbdAdpter)
		globalDbdAdpterIntance.newDbdAdapter()
	}

	return globalDbdAdpterIntance
}

func (adapter *DbdAdpter) NewRegistration(requestForm interface{}) error {
	// Declaration
	var err error

	// Defination
	err = nil

	// Alogorithm
	err = errors.New("not implement")

	// Return
	return err
}

func (adapter *DbdAdpter) ChangeRegistration(requestForm interface{}) error {
	// Declaration
	var err error

	// Defination
	err = nil

	// Alogorithm
	err = errors.New("not implement")

	// Return
	return err
}

func (adapter *DbdAdpter) CancelRegistration(requestForm interface{}) error {
	// Declaration
	var err error

	// Defination
	err = nil

	// Alogorithm
	err = errors.New("not implement")

	// Return
	return err
}

func (adapter *DbdAdpter) LPSNewRegistration(c *gin.Context) (*http.Response, []byte, error) {
	// Declaration
	var err error
	var request *http.Request
	var b []byte

	// Defination
	err = nil
	request = nil

	// Alogorithm
	request, err = http.NewRequestWithContext(c.Request.Context(), c.Request.Method, "https://lps-internal-ms-dbd-service-dev.np.baac.aella.tech/addApiKey", c.Request.Body)

	if err != nil {
		return nil, nil, err
	}

	copyHeaders(request.Header, c.Request.Header)

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Return
	return resp, b, nil
}

func (adapter *DbdAdpter) LPSChangeRegistration(c *gin.Context) (*http.Response, []byte, error) {
	// Declaration
	var err error
	var request *http.Request
	var b []byte

	// Defination
	err = nil
	request = nil

	// Alogorithm
	request, err = http.NewRequestWithContext(c.Request.Context(), c.Request.Method, "https://lps-internal-ms-dbd-service-dev.np.baac.aella.tech/addApiKey", c.Request.Body)

	if err != nil {
		return nil, nil, err
	}

	copyHeaders(request.Header, c.Request.Header)

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Return
	return resp, b, nil
}

func (adapter *DbdAdpter) LPSCancelRegistration(c *gin.Context) (*http.Response, []byte, error) {
	// Declaration
	var err error
	var request *http.Request
	var b []byte

	// Defination
	err = nil
	request = nil

	// Alogorithm
	request, err = http.NewRequestWithContext(c.Request.Context(), c.Request.Method, "lol, lmao even", c.Request.Body)

	if err != nil {
		return nil, nil, err
	}

	copyHeaders(request.Header, c.Request.Header)

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Return
	return resp, b, nil
}

func (adapter *DbdAdpter) LPSGetPreferredInfo(c *gin.Context) (*http.Response, []byte, error) {
	// Declaration
	var err error
	var request *http.Request
	var b []byte

	// Defination
	err = nil
	request = nil

	// Alogorithm
	request, err = http.NewRequestWithContext(c.Request.Context(), c.Request.Method, "lol, lmao even", c.Request.Body)

	if err != nil {
		return nil, nil, err
	}

	copyHeaders(request.Header, c.Request.Header)

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, err
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Return
	return resp, b, nil
}

func (adapter *DbdAdpter) LPSGetSecuredTransaction(c *gin.Context, tb []byte) (*http.Response, []byte, error) {
	// Declaration
	var err error
	var request *http.Request
	var ob []byte

	// Defination
	err = nil
	request = nil

	// Alogorithm
	request, err = http.NewRequestWithContext(c.Request.Context(), c.Request.Method, "https://lps-internal-ms-dbd-service-sit2.np.baac.aella.tech/addApiKeyAndConvertToXML", bytes.NewReader(tb))

	if err != nil {
		return nil, nil, err
	}

	copyHeaders(request.Header, c.Request.Header)

	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, err
	}

	ob, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Return
	return resp, ob, nil
}

// #region JSON Config For Mapping Configuration
// #endregion

func copyHeaders(dst, src http.Header) {
	hopByHop := map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"TE":                  true,
		"Trailer":             true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
		"Content-Type":        false,
		"Authorization":       true,
		"correlationId":       false,
		"X-API-KEY":           true,

		/* // Usually exclude these too:
		   "Host":           true,
		   "Content-Length": true, */
	}

	for k, vv := range src {
		if hopByHop[http.CanonicalHeaderKey(k)] {
			continue
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (dbadapter *DbdAdpter) Holiday(jsonData []byte, c *gin.Context) (*http.Response, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST",
		"https://172.26.82.46:8443/WsSOAR-miniLnInquiry/v1/cert/PNNBDLIST",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "error creating request"})
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("serviceName", "PNNBDLIST")
	req.Header.Set("requestDateTime", "2019-06-20T12:21:30.45+07:00")
	req.Header.Set("clientIpAddress", "172.31.5.66")
	req.Header.Set("clientWorkstationName", "DRMPBTEST")
	req.Header.Set("clientUserId", "1801554")
	req.Header.Set("uuid", "59QAV9C4YKWSSJM358HW7B7B9U")
	req.Header.Set("transactionReferenceId", "LOSL20190304083944978")
	req.Header.Set("channel", "LOS")
	req.Header.Set("apikey", "2cB49mfwQTQVcLtN1FxtXTe1RHTbOY4S")

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error MiniCifList": resp.Status})
		return nil, err
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}

	return resp, nil
}
