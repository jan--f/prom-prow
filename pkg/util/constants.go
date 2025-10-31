package util

const (
	LabelReviewLGTM        = "review/lgtm"
	LabelReviewNeedsReview = "review/needs-review"
	LabelBlockedHold       = "blocked/hold"
	// Review state constants from different GitHub sources
	// Note: webhook events use lowercase state values, while ListReviews API returns uppercase
	ReviewStateApprovedWebhook = "approved" // webhook payloads
	ReviewStateApprovedAPI     = "APPROVED" // REST API responses
)
