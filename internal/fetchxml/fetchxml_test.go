package fetchxml

import (
	"testing"
)

func TestFetchXmlFeed(t *testing.T) {
	// Define the sample XML feed URL
	sampleURL := "https://blog.boot.dev/index.xml"

	// Call the FetchXmlFeed function
	rss, err := FetchXmlFeed(sampleURL)

	// Check if there was an error
	if err != nil {
		t.Fatalf("Error fetching XML feed: %s", err)
	}

	// Perform assertions on the returned RSS struct
	expectedTitle := "The Boot.dev Beat. February 2024"
	if rss.Channel.Items[0].Title != expectedTitle {
		t.Errorf("Unexpected title. Got: %s, Expected: %s", rss.Channel.Items[0].Title, expectedTitle)
	}

	expectedTitle2 := "Give Up Sooner"
	if rss.Channel.Items[1].Title != expectedTitle2 {
		t.Errorf("Unexpected title. Got: %s, Expected: %s", rss.Channel.Items[1].Title, expectedTitle2)
	}

}
