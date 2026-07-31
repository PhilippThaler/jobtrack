package models

// Status values form a fixed workflow:
//
//	applied -> phone_screen -> interview -> offer -> accepted
//	          \-> rejected    \-> rejected  \-> rejected
//	withdrawn (from any state)
const (
	StatusApplied     = "applied"
	StatusPhoneScreen = "phone_screen"
	StatusInterview   = "interview"
	StatusOffer       = "offer"
	StatusAccepted    = "accepted"
	StatusRejected    = "rejected"
	StatusWithdrawn   = "withdrawn"
)

var AllStatuses = []string{
	StatusApplied,
	StatusPhoneScreen,
	StatusInterview,
	StatusOffer,
	StatusAccepted,
	StatusRejected,
	StatusWithdrawn,
}

// Application is a single job application.
type Application struct {
	ID         int64  `json:"id"`
	Company    string `json:"company"`
	Role       string `json:"role"`
	URL        string `json:"url,omitempty"`
	Status     string `json:"status"`
	Notes      string `json:"notes,omitempty"`
	IsFavorite bool   `json:"is_favorite"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Contact is a person at the company you're in touch with.
type Contact struct {
	ID            int64  `json:"id"`
	ApplicationID int64  `json:"application_id"`
	Name          string `json:"name"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Role          string `json:"role,omitempty"` // e.g. "recruiter", "hiring_manager"
}

// TimelineEvent is a dated entry in an application's history,
// e.g. "Applied", "Recruiter replied", "Follow-up sent".
type TimelineEvent struct {
	ID            int64  `json:"id"`
	ApplicationID int64  `json:"application_id"`
	Event         string `json:"event"`
	Note          string `json:"note,omitempty"`
	HappenedAt    string `json:"happened_at"`
}
