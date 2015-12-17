package main

type RedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  interface{} `json:"modhash"`
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Domain     string      `json:"domain"`
				BannedBy   interface{} `json:"banned_by"`
				MediaEmbed struct {
				} `json:"media_embed"`
				Subreddit     string        `json:"subreddit"`
				SelftextHTML  string        `json:"selftext_html"`
				Selftext      string        `json:"selftext"`
				Likes         interface{}   `json:"likes"`
				SuggestedSort interface{}   `json:"suggested_sort"`
				UserReports   []interface{} `json:"user_reports"`
				SecureMedia   interface{}   `json:"secure_media"`
				LinkFlairText interface{}   `json:"link_flair_text"`
				ID            string        `json:"id"`
				FromKind      interface{}   `json:"from_kind"`
				Gilded        int           `json:"gilded"`
				Archived      bool          `json:"archived"`
				Clicked       bool          `json:"clicked"`
				ReportReasons interface{}   `json:"report_reasons"`
				Author        string        `json:"author"`
				Media         interface{}   `json:"media"`
				Score         int           `json:"score"`
				ApprovedBy    interface{}   `json:"approved_by"`
				Over18        bool          `json:"over_18"`
				Hidden        bool          `json:"hidden"`
				Preview       struct {
					Images []struct {
						Source struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"source"`
						Resolutions []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"resolutions"`
						Variants struct {
						} `json:"variants"`
						ID string `json:"id"`
					} `json:"images"`
				} `json:"preview"`
				NumComments         int         `json:"num_comments"`
				Thumbnail           string      `json:"thumbnail"`
				SubredditID         string      `json:"subreddit_id"`
				HideScore           bool        `json:"hide_score"`
				Edited              int         `json:"edited"`
				LinkFlairCSSClass   interface{} `json:"link_flair_css_class"`
				AuthorFlairCSSClass string      `json:"author_flair_css_class"`
				Downs               int         `json:"downs"`
				SecureMediaEmbed    struct {
				} `json:"secure_media_embed"`
				Saved           bool          `json:"saved"`
				RemovalReason   interface{}   `json:"removal_reason"`
				PostHint        string        `json:"post_hint"`
				Stickied        bool          `json:"stickied"`
				From            interface{}   `json:"from"`
				IsSelf          bool          `json:"is_self"`
				FromID          interface{}   `json:"from_id"`
				Permalink       string        `json:"permalink"`
				Locked          bool          `json:"locked"`
				Name            string        `json:"name"`
				Created         int           `json:"created"`
				URL             string        `json:"url"`
				AuthorFlairText string        `json:"author_flair_text"`
				Quarantine      bool          `json:"quarantine"`
				Title           string        `json:"title"`
				CreatedUtc      int           `json:"created_utc"`
				Distinguished   string        `json:"distinguished"`
				ModReports      []interface{} `json:"mod_reports"`
				Visited         bool          `json:"visited"`
				NumReports      interface{}   `json:"num_reports"`
				Ups             int           `json:"ups"`
			} `json:"data"`
		} `json:"children"`
		After  string      `json:"after"`
		Before interface{} `json:"before"`
	} `json:"data"`
}
