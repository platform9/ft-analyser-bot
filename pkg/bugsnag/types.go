package bugsnag

import "time"

type AllErrors []struct {
	URL string `json:"url"`
}

type ErrorDetails struct {
	EventsURL string `json:"events_url"`
}

type AllEventsOfError []struct {
	URL string `json:"url"`
}

type UIErrors struct {
	UserID     string
	ErrorFaced string
	Time       time.Time
}

type EventDetails struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	ProjectURL   string    `json:"project_url"`
	IsFullReport bool      `json:"is_full_report"`
	ErrorID      string    `json:"error_id"`
	ReceivedAt   time.Time `json:"received_at"`
	Exceptions   []struct {
		ErrorClass string `json:"error_class"`
		Message    string `json:"message"`
		Type       string `json:"type"`
		Stacktrace []struct {
			ColumnNumber      int         `json:"column_number"`
			InProject         interface{} `json:"in_project"`
			LineNumber        int         `json:"line_number"`
			Method            string      `json:"method"`
			File              string      `json:"file"`
			Type              interface{} `json:"type"`
			Code              interface{} `json:"code"`
			CodeFile          interface{} `json:"code_file"`
			AddressOffset     interface{} `json:"address_offset"`
			MachoUUID         interface{} `json:"macho_uuid"`
			RelativeAddress   interface{} `json:"relative_address"`
			FrameAddress      interface{} `json:"frame_address"`
			SourceControlLink interface{} `json:"source_control_link"`
			SourceControlName string      `json:"source_control_name"`
		} `json:"stacktrace"`
		Registers interface{} `json:"registers"`
	} `json:"exceptions"`
	Threads  interface{} `json:"threads"`
	MetaData struct {
		App struct {
			CustomerTier string `json:"customerTier"`
			DuVersion    string `json:"duVersion"`
		} `json:"App"`
		Device struct {
			DeviceMemory        int    `json:"deviceMemory"`
			HardwareConcurrency int    `json:"hardwareConcurrency"`
			AppCodeName         string `json:"appCodeName"`
			AppName             string `json:"appName"`
			AppVersion          string `json:"appVersion"`
			CookieEnabled       bool   `json:"cookieEnabled"`
			MaxTouchPoints      int    `json:"maxTouchPoints"`
			OnLine              bool   `json:"onLine"`
			Platform            string `json:"platform"`
			Product             string `json:"product"`
			ProductSub          string `json:"productSub"`
			Vendor              string `json:"vendor"`
			VendorSub           string `json:"vendorSub"`
			Webdriver           bool   `json:"webdriver"`
		} `json:"Device"`
		ResponseConfig struct {
			URL     string `json:"url"`
			Method  string `json:"method"`
			Headers struct {
				Accept        string `json:"Accept"`
				Authorization string `json:"Authorization"`
				XAuthToken    string `json:"X-Auth-Token"`
			} `json:"headers"`
			TransformRequest  []interface{} `json:"transformRequest"`
			TransformResponse []interface{} `json:"transformResponse"`
			Timeout           int           `json:"timeout"`
			XSRFCookieName    string        `json:"xsrfCookieName"`
			XSRFHeaderName    string        `json:"xsrfHeaderName"`
			MaxContentLength  int           `json:"maxContentLength"`
			MaxBodyLength     int           `json:"maxBodyLength"`
		} `json:"Response Config"`
		Request struct {
			BSU string `json:"BS~~U"`
			BSM string `json:"BS~~M"`
			BSS bool   `json:"BS~~S"`
		} `json:"Request"`
		ResponseData struct {
		} `json:"Response Data"`
		Stacktrace struct {
		} `json:"Stacktrace"`
	} `json:"metaData"`
	Request struct {
		URL      string `json:"url"`
		ClientIP string `json:"clientIp"`
		Headers  struct {
		} `json:"headers"`
	} `json:"request"`
	App struct {
		Version      string `json:"version"`
		ReleaseStage string `json:"releaseStage"`
		Duration     int    `json:"duration"`
	} `json:"app"`
	Device struct {
		ID             string    `json:"id"`
		OsName         string    `json:"osName"`
		BrowserName    string    `json:"browserName"`
		BrowserVersion string    `json:"browserVersion"`
		Orientation    string    `json:"orientation"`
		Locale         string    `json:"locale"`
		Time           time.Time `json:"time"`
	} `json:"device"`
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Breadcrumbs []struct {
		Timestamp time.Time `json:"timestamp"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		MetaData  struct {
			Status  int    `json:"status"`
			Request string `json:"request"`
		} `json:"metaData"`
	} `json:"breadcrumbs"`
	Context      string        `json:"context"`
	Severity     string        `json:"severity"`
	Unhandled    bool          `json:"unhandled"`
	FeatureFlags []interface{} `json:"feature_flags"`
}
