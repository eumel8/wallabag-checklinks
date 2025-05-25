package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/go-resty/resty/v2"
)

var (
    wallabagURL    = getEnvOrFail("WALLABAG_URL")
    clientID       = getEnvOrFail("WALLABAG_CLIENT_ID")
    clientSecret   = getEnvOrFail("WALLABAG_CLIENT_SECRET")
    password       = getEnvOrFail("WALLABAG_PASSWORD")
    username       = getEnvOrFail("WALLABAG_USER")
)

type AuthResponse struct {
    AccessToken string `json:"access_token"`
}

type EntryList struct {
    Embedded struct {
        Items []struct {
            ID   int      `json:"id"`
            URL  string   `json:"url"`
            Tags []TagObj `json:"tags"`
        } `json:"items"`
    } `json:"_embedded"`
}

type TagObj struct {
    Label string `json:"label"`
}

type Entry struct {
    ID   int
    URL  string
    Tags []string
}

func getEnvOrFail(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("❌ Environment variable %s is not set", key)
    }
    return value
}

func getAccessToken(client *resty.Client) (string, error) {
    resp, err := client.R().
        SetFormData(map[string]string{
            "grant_type":    "password",
            "client_id":     clientID,
            "client_secret": clientSecret,
            "username":      username,
            "password":      password,
        }).
        SetResult(&AuthResponse{}).
        Post(wallabagURL + "/oauth/v2/token")

    if err != nil {
        return "", err
    }

    result := resp.Result().(*AuthResponse)
    if result.AccessToken == "" {
        log.Fatalf("❌ No access token: %v", err)
	os.Exit(1)
    }

    return result.AccessToken, nil
}

func getEntries(client *resty.Client, token string) ([]Entry, error) {
    var entries EntryList

    _, err := client.R().
        SetHeader("Authorization", "Bearer "+token).
        SetQueryParams(map[string]string{
            "perPage": "10000", // Adjust for more
        }).
        SetResult(&entries).
        Get(wallabagURL + "/api/entries.json")

    if err != nil {
        return nil, err
    }

    var result []Entry
    for _, item := range entries.Embedded.Items {
        tags := []string{}
        for _, tag := range item.Tags {
            tags = append(tags, tag.Label)
        }
        result = append(result, Entry{
            ID:   item.ID,
            URL:  item.URL,
            Tags: tags,
        })
    }

    return result, nil
}

func checkURL(url string) int {
    client := &http.Client{Timeout: 10 * time.Second}

    req, _ := http.NewRequest("HEAD", url, nil)
    resp, err := client.Do(req)
    if err != nil {
        return 0
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return 0
    }
    return resp.StatusCode
}

func tagEntry(client *resty.Client, token string, entryID int, currentTags []string, newTag string) error {
    // Skip if already tagged
    for _, tag := range currentTags {
        if tag == newTag {
            return nil
        }
    }

    // Append and send full tag list
    updatedTags := append(currentTags,newTag)
    data := map[string]interface{}{
        "tags": updatedTags,
    }
    body, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal tags: %w", err)
    }

    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+token).
        SetHeader("Content-Type", "application/json").
        SetBody(body).
        Post(fmt.Sprintf("%s/api/entries/%d/tags.json", wallabagURL, entryID))

    if err != nil {
        return fmt.Errorf("tagging request failed: %w", err)
    }

    if resp.IsError() {
        return fmt.Errorf("tagging failed with status %d: %s", resp.StatusCode(), resp.String())
    }

    return nil
}

func getTagIDByLabel(client *resty.Client, token, label string) (int, error) {
    var tags []struct {
        ID    int    `json:"id"`
        Label string `json:"label"`
    }

    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+token).
        SetResult(&tags).
        Get(fmt.Sprintf("%s/api/tags.json", wallabagURL))

    if err != nil || resp.IsError() {
        return 0, fmt.Errorf("failed to fetch tags: %v (%s)", err, resp.Status())
    }

    for _, tag := range tags {
        if tag.Label == label {
            return tag.ID, nil
        }
    }

    return 0, fmt.Errorf("tag '%s' not found", label)
}

func removeTag(client *resty.Client, token string, entryID int, tagID int) error {
    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+token).
        SetHeader("Content-Type", "application/json").
        Delete(fmt.Sprintf("%s/api/entries/%d/tags/%d.json", wallabagURL, entryID, tagID))

    if err != nil {
        return fmt.Errorf("tag deleting failed: %w", err)
    }

    if resp.IsError() {
        return fmt.Errorf("tag deleting failed with status %d: %s", resp.StatusCode(), resp.String())
    }
 
    return nil
}

func containsTag(tags []string, target string) bool {
    for _, tag := range tags {
        if tag == target {
            return true
        }
    }
    return false
}

func main() {
    restyClient := resty.New()

    token, err := getAccessToken(restyClient)
    if err != nil {
        log.Fatalf("❌ Failed to get token: %v", err)
    }

    entries, err := getEntries(restyClient, token)
    if err != nil {
        log.Fatalf("❌ Failed to fetch entries: %v", err)
    }

    fmt.Printf("🔎 Checking %d URLs...\n\n", len(entries))

    for _, entry := range entries {
        status := checkURL(entry.URL)
        hasDeadTag := containsTag(entry.Tags, "dead")

        if status == 0 {
            fmt.Printf("❌ DEAD - %s\n", entry.URL)
            if !hasDeadTag {
		err := tagEntry(restyClient, token, entry.ID, entry.Tags, "dead")
                if err != nil {
                    fmt.Printf("⚠️ Failed to tag entry %d: %v\n", entry.ID, err)
                } else {
                    fmt.Println("📝 Tagged with 'dead'")
                }
            }
        } else {
            fmt.Printf("✅ %d - %s\n", status, entry.URL)
            if hasDeadTag {
		tagID, err := getTagIDByLabel(restyClient, token, "dead")
                if err != nil {
                    log.Fatalf("Could not find 'dead' tag: %v", err)
                }
                err = removeTag(restyClient, token, entry.ID, tagID)
                if err != nil {
                    fmt.Printf("⚠️ Failed to remove 'dead' tag for entry %d: %v\n", entry.ID, err)
                } else {
                    fmt.Println("🧼 Removed 'dead' tag")
                }
            }
        }
        fmt.Println()
    }
}

