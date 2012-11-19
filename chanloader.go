package main

import (
	"encoding/json"
	"errors"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	version string = "0.1"
)

var (
	// threadRegex = regexp.Compile("/?([a-z]+)(?:res)?/(\\d+)")
	lastDownload int64
	board        string
	threadId     string
	boards       string         = "a|b|c|d|e|f|g|gif|h|hr|k|m|o|p|r|s|t|u|v|vg|w|wg|i|ic|r9k|cm|hm|y|3|adv|an|cgl|ck|co|diy|fa|fit|hc|int|jp|lit|mlp|mu|n|po|pol|sci|soc|sp|tg|toy|trv|tv|vp|wsg|x|q"
	minWidth     *int64         = flag.Int64P("min-width", "w", 0, "Minimum width of images")
	minHeight    *int64         = flag.Int64P("min-height", "h", 0, "Minimum height of images")
	refresh      *time.Duration = flag.DurationP("refresh", "r", time.Second*30, "Refresh rate (min 30s)")
	showVersion  *bool          = flag.BoolP("version", "v", false, "Show version")
)

type Post struct {
	No             int64
	Resto          int64
	Sticky         int64
	Closed         int64
	Now            string
	Time           int64
	Name           string
	Trip           string
	Id             string
	Capcode        string
	Country        string
	Country_name   string
	Email          string
	Sub            string
	Com            string
	Tim            int64
	Filename       string
	Ext            string
	Fsize          int64
	Md5            string
	W              int64
	H              int64
	Tn_w           int64
	Tn_h           int64
	Filedeleted    int64
	Spoiler        int64
	Custom_spoiler int64
	Omitted_posts  int64
	Omitted_images int64
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type Thread struct {
	Posts []Post
}

func parseThreadId(s string) (b, id string, err error) {
	r, err := regexp.Compile(fmt.Sprintf("^/?(%s)/(?:res/)?(\\d+)", boards))
	if err != nil {
		return b, id, err
	}

	groups := r.FindSubmatch([]byte(s))
	if len(groups) < 1 {
		return b, id, errors.New("Input does not match regex")
	}

	b = string(groups[1])
	id = string(groups[2])

	return b, id, nil
}

func parseThreadFromJson(r io.Reader) (Thread, error) {
	dec := json.NewDecoder(r)
	var t Thread
	for {
		if err := dec.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			return t, err
		}
	}
	return t, nil
}

func downloadImage(tim int64, ext string) {
	lastDownload = tim
	fileName := fmt.Sprintf("%s%s", strconv.FormatInt(tim, 10), ext)
	url := fmt.Sprintf("https://images.4chan.org/%s/src/%s", board, fileName)

	img, err := getUrl(url)
	checkError(err)

	err = ioutil.WriteFile(fileName, img, 0644)
	checkError(err)
}

func getUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, nil
	}
	return body, nil

}

func loadThread() {
	resp, err := getUrl(fmt.Sprintf("https://api.4chan.org/%s/res/%s.json", board, threadId))
	checkError(err)

	thread, err := parseThread(strings)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: chanloader [options] /b/res/123456\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	b, id, err := parseThreadId(flag.Arg(0))
	if err != nil {
		fmt.Printf("Invalid input: %s\n", flag.Arg(0))
		os.Exit(1)
	}
	board, threadId = b, id
	fmt.Println(b, id)

	ticker := time.NewTicker(*refresh)
	for {
		fmt.Println("Loading thread")
		go loadThread()
		<-ticker.C
	}

	var x string
	fmt.Scanln(&x)
}
