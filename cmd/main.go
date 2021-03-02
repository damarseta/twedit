package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"go.damarseta.id/twedit"
)

var (
	ctx        = context.Background()
	subreddits = []string{
		"curledfeetsies",
		"Catloaf",
		"catsareliquid",
		"Catswithjobs",
		"CatsStandingUp",
		"TuxedoCats",
	}
)

func main() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	if err := run(); err != nil {
		log.Fatal(err)
	}

}

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return
}

func run() (err error) {
	twc, err := LoadCredentials()
	if err != nil {
		return err
	}

	subs := subreddits[rand.Intn(len(subreddits))]
	client, err := reddit.NewReadonlyClient()
	if err != nil {
		return
	}

	client.OnRequestCompleted(logResponse)

	posts, _, err := client.Subreddit.TopPosts(ctx, subs, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
		},
		Time: "all",
	})
	if err != nil {
		return
	}

	pick := posts[rand.Intn(len(posts))]

	pc, _, err := client.Post.Get(ctx, pick.ID)
	if err != nil {
		fmt.Printf("error loading post detail: %s\n", err)
		return err
	}

	filename, err := getFilename(pc.Post.URL)
	if err != nil {
		return err
	}

	fpath, err := twedit.Download(pc.Post.URL, filename)
	if err != nil {
		return err
	}

	body, header, err := twedit.GetBody(fpath)
	if err != nil {
		return err
	}

	endpoint := "/1.1/statuses/update_with_media.json"
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", header)

	resp, err := twc.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resp.Status)

	tweet := &twittergo.Tweet{}
	err = resp.Parse(tweet)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ID:                         %v\n", tweet.Id())
	fmt.Printf("Tweet:                      %v\n", tweet.Text())
	fmt.Printf("User:                       %v\n", tweet.User().Name())
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:                 %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining:       %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:           %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
	if resp.HasMediaRateLimit() {
		fmt.Printf("Media Rate limit:           %v\n", resp.MediaRateLimit())
		fmt.Printf("Media Rate limit remaining: %v\n", resp.MediaRateLimitRemaining())
		fmt.Printf("Media Rate limit reset:     %v\n", resp.MediaRateLimitReset())
	} else {
		fmt.Printf("Could not parse media rate limit from response.\n")
	}
	return
}

func logResponse(req *http.Request, res *http.Response) {
	fmt.Printf("%s %s %s\n", req.Method, req.URL, res.Status)
}

func getFilename(rawurl string) (string, error) {
	newUrl := rawurl

	// Skip blank lines
	if newUrl == "" {
		return "", fmt.Errorf("empty url")
	}

	sep := `\`
	if strings.Index(newUrl, "/") > -1 {
		sep = "/"
	}

	if newUrl[len(newUrl)-1:] == sep {
		return "", fmt.Errorf("invalid url length")
	}

	newUrl = strings.Replace(newUrl, "http://", "", -1)
	newUrl = strings.Replace(newUrl, "https://", "", -1)

	urlParts := strings.Split(newUrl, sep)
	filename := urlParts[len(urlParts)-1]

	return filename, nil
}
