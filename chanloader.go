package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"
)

const (
	version string = "0.0.2"
)

var (
	downloaded []int64
	boards     string = "a|b|c|d|e|f|g|gif|h|hr|k|m|o|p|r|s|t|u|v|vg|w|wg|i|ic|r9k|cm|hm|y|3|adv|an|cgl|ck|co|diy|fa|fit|hc|int|jp|lit|mlp|mu|n|po|pol|sci|soc|sp|tg|toy|trv|tv|vp|wsg|x|q"

	// Options
	minWidth     *int64         = flag.Int64P("min-width", "w", 0, "Minimum width of images")
	minHeight    *int64         = flag.Int64P("min-height", "h", 0, "Minimum height of images")
	refresh      *time.Duration = flag.DurationP("refresh", "r", time.Second*30, "Refresh rate (min 30s)")
	orignalNames *bool          = flag.BoolP("original-names", "o", false, "Save images under original filenames")
	showVersion  *bool          = flag.BoolP("version", "v", false, "Show version")
)

type Post struct {
	No             int64
	Resto          int64
	Sticky         bool
	Closed         bool
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
	Filedeleted    bool
	Spoiler        bool
	Custom_spoiler bool
	Omitted_posts  int64
	Omitted_images int64
}

func (p *Post) FullFileName(original bool) (s string) {
	if original {
		s = fmt.Sprintf("%s%s", p.Filename, p.Ext)
	} else {
		s = fmt.Sprintf("%d%s", p.Tim, p.Ext)
	}
	return
}

type Thread struct {
	Posts []Post
}

func checkError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func parseThreadId(s string) (b, id string, err error) {
	r, err := regexp.Compile(fmt.Sprintf("^(?:https?://boards.4chan.org)?/?(%s)/(?:thread/)?(\\d+).+$", boards))
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

func downloadImage(board string, p Post) {
	if p.W < *minWidth || p.H < *minHeight || p.Filedeleted {
		downloaded = append(downloaded, p.Tim)
		return
	}

	url := fmt.Sprintf("https://images.4chan.org/%s/src/%d%s", board, p.Tim, p.Ext)
	img, _, err := getUrl(url)
	checkError(err)

	err = ioutil.WriteFile(p.FullFileName(*orignalNames), img, 0644)
	checkError(err)

	downloaded = append(downloaded, p.Tim)
}

func getUrl(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, 0, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, 0, nil
	}
	return body, resp.StatusCode, nil

}

func wasDownloaded(p Post) bool {
	for _, e := range downloaded {
		if e == p.Tim {
			return true
		}
	}
	return false
}

func loadThread(board string, id string) {
	url := fmt.Sprintf("https://api.4chan.org/%s/res/%s.json", board, id)
	resp, status, err := getUrl(url)
	checkError(err)
	if status != 200 {
		if status == 404 {
			log.Println("Thread 404'd")
			os.Exit(1)
		}
		fmt.Printf("Something went wrong\nGot return code %d at '%s'.\n", status, url)
		return
	}

	thread, err := parseThreadFromJson(bytes.NewReader(resp))
	if err != nil {
		log.Println("Could not parse thread")
	}

	for _, e := range thread.Posts {
		if !wasDownloaded(e) {
			go downloadImage(board, e)
		}
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: chanloader [options] [http(s)://boards.4chan.org/]b/thread/123456[/thread-name]\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *refresh < time.Second*30 {
		*refresh = time.Second * 30
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	board, id, err := parseThreadId(flag.Arg(0))
	if err != nil {
		fmt.Printf("Invalid input: %s\n", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	}

	ticker := time.NewTicker(*refresh)
	for {
		log.Println("Loading thread")
		go loadThread(board, id)
		<-ticker.C
	}
}
