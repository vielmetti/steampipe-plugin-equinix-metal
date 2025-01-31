package metal

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	metal "github.com/packethost/packngo"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func connect(_ context.Context, d *plugin.QueryData) (*metal.Client, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "metal"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*metal.Client), nil
	}

	var token string

	// First, use the token config
	metalConfig := GetConfig(d.Connection)
	if &metalConfig != nil {
		if metalConfig.Token != nil {
			token = *metalConfig.Token
		}
	}

	// Otherwise, default to using PACKET_AUTH_TOKEN env var (it appears the
	// product was later renamed to Equinix Metal)
	if token == "" {
		token = os.Getenv("PACKET_AUTH_TOKEN")
	}

	conn := metal.NewClientWithAuth("steampipe", token, http.DefaultClient)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn, nil
}

func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."
	}
	return ""
}

func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		errMsg := err.Error()
		for _, msg := range notFoundErrors {
			if strings.Contains(errMsg, msg) {
				return true
			}
		}
		return false
	}
}

func getParentProjectID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	parentProject := h.HydrateResults["listProject"].(metal.Project)
	return parentProject.ID, nil
}

func hrefToID(_ context.Context, d *transform.TransformData) (interface{}, error) {
	href := d.Value.(string)
	parts := strings.Split(href, "/")
	return parts[len(parts)-1], nil
}

func timestampToIsoTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.Value.(*metal.Timestamp)
	return i.Format(time.RFC3339), nil
}

func projectsToIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]metal.Project)
	ids := []string{}
	for _, i := range items {
		if i.ID != "" {
			ids = append(ids, i.ID)
		} else {
			parts := strings.Split(i.URL, "/")
			ids = append(ids, parts[len(parts)-1])
		}
	}
	return ids, nil
}

func facilitiesToIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]metal.Facility)
	ids := []string{}
	for _, i := range items {
		if i.ID != "" {
			ids = append(ids, i.ID)
		} else {
			parts := strings.Split(i.URL, "/")
			ids = append(ids, parts[len(parts)-1])
		}
	}
	return ids, nil
}

func sshKeysToIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]metal.SSHKey)
	ids := []string{}
	for _, i := range items {
		if i.ID != "" {
			ids = append(ids, i.ID)
		} else {
			parts := strings.Split(i.URL, "/")
			ids = append(ids, parts[len(parts)-1])
		}
	}
	return ids, nil
}

func volumesToIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]metal.Volume)
	ids := []string{}
	for _, i := range items {
		if i.ID != "" {
			ids = append(ids, i.ID)
		} else {
			parts := strings.Split(i.Href, "/")
			ids = append(ids, parts[len(parts)-1])
		}
	}
	return ids, nil
}

func usersToIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]metal.User)
	ids := []string{}
	for _, i := range items {
		if i.ID != "" {
			ids = append(ids, i.ID)
		} else {
			parts := strings.Split(i.URL, "/")
			ids = append(ids, parts[len(parts)-1])
		}
	}
	return ids, nil
}
