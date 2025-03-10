package telexcom

// Date structure
type Date struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Descriptions structure
type Descriptions struct {
	AppDescription  string `json:"app_description"`
	AppLogo         string `json:"app_logo"`
	AppName         string `json:"app_name"`
	AppURL          string `json:"app_url"`
	BackgroundColor string `json:"background_color"`
}

// MonitoringUser structure
type MonitoringUser struct {
	AlwaysOnline bool   `json:"always_online"`
	DisplayName  string `json:"display_name"`
}

// Permissions structure
type Permissions struct {
	MonitoringUser MonitoringUser `json:"monitoring_user"`
}

// Setting structure
type Setting struct {
	Label    string `json:"label"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  string `json:"default"`
}

type Data struct {
	Date                Date         `json:"date"`
	Descriptions        Descriptions `json:"descriptions"`
	IntegrationCategory string       `json:"integration_category"`
	IntegrationType     string       `json:"integration_type"`
	IsActive            bool         `json:"is_active"`
	KeyFeatures         []string     `json:"key_features"`
	Author              string       `json:"author"`
	Settings            []Setting    `json:"settings"`
	TickURL             string       `json:"tick_url"`
	TargetURL           string       `json:"target_url"`
}

// IntegrationJson main structure
type Integration struct {
	Data Data `json:"data"`
}

var IntegrationJson = Integration{
	Data: Data{
		Date: Date{
			CreatedAt: "2025-03-05",
			UpdatedAt: "2025-03-05",
		},
		Descriptions: Descriptions{
			AppDescription:  "This integration listens in channels of organisations that use it, for Frequently Asked Questions about the organisation, its documentation or any document granted access to the AI integration for training and personification.",
			AppLogo:         "https://my-portfolio-343207.web.app/MyLogo4.png",
			AppName:         "CustomerServiceAIChatbot",
			AppURL:          "https://96ac-102-216-183-26.ngrok-free.app/integration",
			BackgroundColor: "#fff",
		},
		IntegrationCategory: "AI & Machine Learning",
		IntegrationType:     "output",
		IsActive:            true,
		KeyFeatures:         []string{"Gives apt responses to FAQs.", "Sends users response based on data stored in database."},
		Author:              "Samuel Ikoli",
		Settings: []Setting{
			{
				Label:    "File-path",
				Type:     "text",
				Required: true,
				Default:  "",
			},
			{
				Label:    "Webhook",
				Type:     "text",
				Required: true,
				Default:  "",
			},
		},
		TargetURL: "https://96ac-102-216-183-26.ngrok-free.app/target",
	},
}

type MonitorPayload struct {
	ChannelID string        `json:"channel_id,omitempty"`
	ReturnURL string        `json:"return_url,omitempty"`
	Settings  []interface{} `json:"settings,omitempty"`
}

type PromptPayload struct {
	PathURL  string `json:"path_url,omitempty"`
	Question string `json:"question,omitempty"`
}

type TelexPromptPayload struct {
	Message  string    `json:"message,omitempty"`
	Settings []Setting `json:"settings,omitempty"`
}

type GeminiResponse struct {
	Response string `json:"response"`
}
