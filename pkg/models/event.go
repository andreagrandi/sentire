package models

import "time"

// Event represents a complete Sentry event with all fields
type Event struct {
	ID          string                 `json:"id"`
	EventID     string                 `json:"eventID"`
	ProjectID   string                 `json:"projectID"`
	GroupID     string                 `json:"groupID,omitempty"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Platform    string                 `json:"platform"`
	Type        string                 `json:"type"`
	DateCreated time.Time              `json:"dateCreated"`
	DateReceived time.Time             `json:"dateReceived"`
	Size        int64                  `json:"size"`
	Dist        string                 `json:"dist,omitempty"`
	Location    string                 `json:"location,omitempty"`
	Logger      string                 `json:"logger,omitempty"`
	Culprit     string                 `json:"culprit,omitempty"`
	
	// Core debugging information
	Entries     []Entry                `json:"entries"`
	Exception   *Exception             `json:"exception,omitempty"`
	Breadcrumbs *Breadcrumbs           `json:"breadcrumbs,omitempty"`
	Request     *Request               `json:"request,omitempty"`
	
	// Context and metadata
	Tags        []EventTag             `json:"tags"`
	User        *EventUser             `json:"user,omitempty"`
	Contexts    *Contexts              `json:"contexts,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Fingerprint []string               `json:"fingerprint"`
	
	// Release and SDK info
	Release     *EventRelease          `json:"release,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	SDK         *EventSDK              `json:"sdk,omitempty"`
	
	// Error handling
	Errors      []EventError           `json:"errors,omitempty"`
}

// Entry represents different types of entries in an event (exception, breadcrumbs, request, etc.)
type Entry struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Exception contains exception information with stack traces
type Exception struct {
	Values []ExceptionValue `json:"values"`
}

// ExceptionValue represents a single exception with stack trace
type ExceptionValue struct {
	Type       string       `json:"type"`
	Value      string       `json:"value"`
	Module     string       `json:"module,omitempty"`
	ThreadID   *int64       `json:"threadId,omitempty"`
	Mechanism  *Mechanism   `json:"mechanism,omitempty"`
	Stacktrace *Stacktrace  `json:"stacktrace,omitempty"`
	RawStacktrace *Stacktrace `json:"rawStacktrace,omitempty"`
}

// Mechanism describes how an exception was captured
type Mechanism struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Handled     *bool                  `json:"handled,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// Stacktrace contains stack frames
type Stacktrace struct {
	Frames         []StackFrame `json:"frames"`
	FramesOmitted  []int        `json:"framesOmitted,omitempty"`
	Registers      map[string]string `json:"registers,omitempty"`
	HasSystemFrames bool        `json:"hasSystemFrames,omitempty"`
}

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	Filename         string                 `json:"filename"`
	Function         string                 `json:"function"`
	Module           string                 `json:"module,omitempty"`
	LineNo           *int                   `json:"lineNo,omitempty"`
	ColNo            *int                   `json:"colNo,omitempty"`
	AbsPath          string                 `json:"absPath,omitempty"`
	ContextLine      string                 `json:"contextLine,omitempty"`
	PreContext       []string               `json:"preContext,omitempty"`
	PostContext      []string               `json:"postContext,omitempty"`
	InApp            *bool                  `json:"inApp,omitempty"`
	Vars             map[string]interface{} `json:"vars,omitempty"`
	Package          string                 `json:"package,omitempty"`
	Platform         string                 `json:"platform,omitempty"`
	ImageAddr        string                 `json:"imageAddr,omitempty"`
	InstructionAddr  string                 `json:"instructionAddr,omitempty"`
	AddrMode         string                 `json:"addrMode,omitempty"`
	SymbolAddr       string                 `json:"symbolAddr,omitempty"`
	Symbol           string                 `json:"symbol,omitempty"`
	Trust            string                 `json:"trust,omitempty"`
	Lock             map[string]interface{} `json:"lock,omitempty"`
}

// Breadcrumbs contains the breadcrumb trail
type Breadcrumbs struct {
	Values []Breadcrumb `json:"values"`
}

// Breadcrumb represents a single breadcrumb
type Breadcrumb struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	Category  string                 `json:"category"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Level     string                 `json:"level,omitempty"`
}

// Request contains HTTP request information
type Request struct {
	URL         string                 `json:"url"`
	Method      string                 `json:"method"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Data        interface{}            `json:"data,omitempty"`
	QueryString interface{}            `json:"queryString,omitempty"`
	Cookies     map[string]string      `json:"cookies,omitempty"`
	Env         map[string]string      `json:"env,omitempty"`
	Fragment    string                 `json:"fragment,omitempty"`
	InferredContentType string         `json:"inferredContentType,omitempty"`
}

// Contexts contains various context information
type Contexts struct {
	Browser  *BrowserContext  `json:"browser,omitempty"`
	Client   *ClientContext   `json:"client_os,omitempty"`
	Device   *DeviceContext   `json:"device,omitempty"`
	OS       *OSContext       `json:"os,omitempty"`
	Runtime  *RuntimeContext  `json:"runtime,omitempty"`
	App      *AppContext      `json:"app,omitempty"`
	GPU      *GPUContext      `json:"gpu,omitempty"`
	Monitor  *MonitorContext  `json:"monitor,omitempty"`
	Culture  *CultureContext  `json:"culture,omitempty"`
	Cloud    *CloudContext    `json:"cloud_resource,omitempty"`
	Trace    *TraceContext    `json:"trace,omitempty"`
}

// BrowserContext contains browser information
type BrowserContext struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type,omitempty"`
}

// ClientContext contains client OS information  
type ClientContext struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Build        string `json:"build,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty"`
	Type         string `json:"type,omitempty"`
}

// DeviceContext contains device information
type DeviceContext struct {
	Name         string  `json:"name"`
	Family       string  `json:"family"`
	Model        string  `json:"model,omitempty"`
	ModelID      string  `json:"modelId,omitempty"`
	Arch         string  `json:"arch,omitempty"`
	BatteryLevel *float64 `json:"batteryLevel,omitempty"`
	Charging     *bool   `json:"charging,omitempty"`
	LowMemory    *bool   `json:"lowMemory,omitempty"`
	Online       *bool   `json:"online,omitempty"`
	Orientation  string  `json:"orientation,omitempty"`
	Simulator    *bool   `json:"simulator,omitempty"`
	MemorySize   *int64  `json:"memorySize,omitempty"`
	FreeMemory   *int64  `json:"freeMemory,omitempty"`
	UsableMemory *int64  `json:"usableMemory,omitempty"`
	StorageSize  *int64  `json:"storageSize,omitempty"`
	FreeStorage  *int64  `json:"freeStorage,omitempty"`
	ExternalStorageSize *int64 `json:"externalStorageSize,omitempty"`
	ExternalFreeStorage *int64 `json:"externalFreeStorage,omitempty"`
	BootTime     *time.Time `json:"bootTime,omitempty"`
	ProcessorCount *int    `json:"processorCount,omitempty"`
	ProcessorFrequency *int64 `json:"processorFrequency,omitempty"`
	CpuDescription string   `json:"cpuDescription,omitempty"`
	Type         string   `json:"type,omitempty"`
}

// OSContext contains operating system information
type OSContext struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Build        string `json:"build,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty"`
	Rooted       *bool  `json:"rooted,omitempty"`
	Type         string `json:"type,omitempty"`
}

// RuntimeContext contains runtime information
type RuntimeContext struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Build   string `json:"build,omitempty"`
	Type    string `json:"type,omitempty"`
}

// AppContext contains application information
type AppContext struct {
	AppName        string     `json:"appName,omitempty"`
	AppVersion     string     `json:"appVersion,omitempty"`
	AppBuild       string     `json:"appBuild,omitempty"`
	AppIdentifier  string     `json:"appIdentifier,omitempty"`
	AppStartTime   *time.Time `json:"appStartTime,omitempty"`
	DeviceAppHash  string     `json:"deviceAppHash,omitempty"`
	BuildType      string     `json:"buildType,omitempty"`
	AppMemory      *int64     `json:"appMemory,omitempty"`
	Type           string     `json:"type,omitempty"`
}

// GPUContext contains GPU information
type GPUContext struct {
	Name                 string `json:"name"`
	ID                   *int   `json:"id,omitempty"`
	VendorID             string `json:"vendorId,omitempty"`
	VendorName           string `json:"vendorName,omitempty"`
	MemorySize           *int64 `json:"memorySize,omitempty"`
	APIType              string `json:"apiType,omitempty"`
	MultiThreadedRendering *bool `json:"multiThreadedRendering,omitempty"`
	Version              string `json:"version,omitempty"`
	NpotSupport          string `json:"npotSupport,omitempty"`
	Type                 string `json:"type,omitempty"`
}

// MonitorContext contains monitor information
type MonitorContext struct {
	DPI    *int `json:"dpi,omitempty"`
	Height *int `json:"height,omitempty"`
	Width  *int `json:"width,omitempty"`
	Type   string `json:"type,omitempty"`
}

// CultureContext contains culture information
type CultureContext struct {
	Calendar      string `json:"calendar,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Is24HourFormat *bool `json:"is24HourFormat,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
	Type          string `json:"type,omitempty"`
}

// CloudContext contains cloud resource information
type CloudContext struct {
	Provider         string `json:"provider,omitempty"`
	AccountID        string `json:"accountId,omitempty"`
	Region           string `json:"region,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	MachineType      string `json:"machineType,omitempty"`
	ProjectID        string `json:"projectId,omitempty"`
	Type             string `json:"type,omitempty"`
}

// TraceContext contains tracing information
type TraceContext struct {
	TraceID      string                 `json:"traceId,omitempty"`
	SpanID       string                 `json:"spanId,omitempty"`
	ParentSpanID string                 `json:"parentSpanId,omitempty"`
	Op           string                 `json:"op,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Tags         map[string]interface{} `json:"tags,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Type         string                 `json:"type,omitempty"`
}

// EventTag represents a tag on an event
type EventTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// EventUser represents user information in an event
type EventUser struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	IPAddr   string `json:"ip_address,omitempty"`
}

// EventRelease represents release information
type EventRelease struct {
	Version      string `json:"version"`
	ShortVersion string `json:"shortVersion,omitempty"`
}

// EventSDK represents SDK information
type EventSDK struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// EventError represents processing errors
type EventError struct {
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}