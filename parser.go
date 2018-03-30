package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"parser_go/config"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type parserOnePage struct {
	timeSleep      int
	gol            int
	tempCookies    []*http.Cookie
	tempName       string
	temp_files_src map[string]string
	temp_files     map[string]string
	rootDir        string
	baseLink       string
	hrefAllLinks   string
	userAgent      string
	dirAgent       string
	mobAgent       string
	descAgent      string
	links          map[string]string
	countLink      int
	dirs           map[string]string
	indFile        string
	index          string
	script         string
	notIframe      bool
	baseTeg        string
	ajax           bool
	dir            string
}

const (
	TypeImage = iota
	TypeSrc
)

var instance *parserOnePage
var ch chan int
var chName chan string
var chUrl chan string
var chTypeImages chan bool
var chTypeSrc chan bool

func main() {
	parser := getInstance(config.GetConfig())
	fmt.Println("start")
	parser.run()
}

func (p *parserOnePage) run() *parserOnePage {

	if len(p.links) < 1 {
		log.Fatal("Not init link(s)")
	}
	message("parse start", "keep", strconv.Itoa(len(p.links)), "page(s)")

	for name, link := range p.links {
		p.baseLink = link
		p.setTempName(name)
		p.setOptions()
		p.parsePage(link)
	}

	return p
}
func (p *parserOnePage) parsePage(link string) {
	err, res := p.request(link)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	p.thenBaseHref(doc)
	p.saveIco(doc)
	p.saveModifyElemHref(doc)
	p.saveModifyIframe(doc)
	p.saveModifyJs(doc)
	p.saveModifayCss(doc)
	p.saveModifyImg(doc)
	p.modifyForm(doc)


	name, url := "", ""
	for i := 1; i < 5; i++ {
		go p.multiFiles(name, url)
	}
	for name, url := range p.temp_files {
		chName <- name
		chUrl <- url
		chTypeImages <- true
	}
	for name, url := range p.temp_files_src {
		chName <- name
		chUrl <- url
		chTypeSrc <- true
	}
	chName <- ""
	chUrl <- ""
	fmt.Println(<-ch)


	/**
	        $this->replaceCssInHtml($patch);
	        $this->add($page,'<script type="text/javascript" src="/js/script.js"></script>');
		**/

	p.save(doc, p.rootDir+"/"+p.tempName+"/")
	message("pages url=", p.baseLink)
	message("OK All create")

}
func (p *parserOnePage) multiFiles(name string, url string) {

	for {

		name = <-chName
		url = <-chUrl
		if name == "" && url == "" {
			ch <- 10
			break
		}
		select {
		case <-chTypeImages:
			p.saveFileGo(url, p.dirs["img"], name)
		case <-chTypeSrc:
			p.saveFileGo(url, p.dirs["src"], name)
		}

	}

}

func (p *parserOnePage) thenBaseHref(doc *goquery.Document) {
	doc.Find("base").Remove()
}

func (p *parserOnePage) saveIco(doc *goquery.Document) {
	doc.Find("link[rel=icon],link[rel=apple-touch-icon]").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if link != "" {
			link = p.urlAbsolute(link, p.baseLink)
			name := p.filterFileName(link)
			p.saveFile(link, p.dirs["img"], name)
			s.SetAttr("href", p.dirs["img_r"]+name)
		}
	})
}

func (p *parserOnePage) saveModifyElemHref(doc *goquery.Document) {

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("href", p.hrefAllLinks)

		if t, _ := s.Attr("target"); t != "" {
			s.RemoveAttr("target")
		}
		if t, _ := s.Attr("onclick"); t != "" {
			s.RemoveAttr("onclick")
		}

	})
}

func (p *parserOnePage) saveModifyIframe(doc *goquery.Document) {
	if p.notIframe {
		doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
	}
}
func (p *parserOnePage) modifyForm(doc *goquery.Document) {
		doc.Find("form").Each(func(i int, s *goquery.Selection) {
			s.SetAttr("action","javascript:void(0)")
		})
}

/**
protected function replaceCssInHtml($patch)
    {
        $html_all_page=file_get_contents($patch);
        $result = preg_replace_callback('/url\((.*?)\)/',
            function ($matches)  {
                return 'url(' . $this->getfile($matches[1], $this->uri2absolute($matches[1],$this->baseLink),true) . ' )';
            }
            , $html_all_page);
        file_put_contents($patch,$result);
        $this->saveCurlFile($this->multiCurl($this->temp_files_src),'src');
        unset($this->temp_files_src);
        
    }
**/
func (p *parserOnePage) replaceCssInHtml(patch string) {

	source, err := ioutil.ReadFile(patch)
	if err != nil {
		//TODO 
		//return;
	}
	//source

//	err = ioutil.WriteFile(fi.Name(), []byte(Value), 0644)


}



func (p *parserOnePage) saveModifyJs(doc *goquery.Document) {

	reg3 := regexp.MustCompile(`~^https:\/\/ajax.googleapis.com~`)

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src != "" {
			if !reg3.MatchString(src) {
				name := p.filterFileName(src)
				s.SetAttr("src", p.dirs["js_r"]+name)
				if _, err := os.Stat(p.dirs["js"] + name); os.IsNotExist(err) {

					link := p.urlAbsolute(src, p.baseLink)
					_, response := p.request(link)
					if response != nil {
						defer response.Body.Close()
						out, err := os.Create(p.dirs["js"] + name)
						defer out.Close()
						if err == nil {
							io.Copy(out, response.Body)

						}
					}

				}

			}

		}
	})

}

func (p *parserOnePage) request(url string) (error, *http.Response) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//TODO log
		return err, nil
	}
	req.Header.Set("Accept-language", "en")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	//add
	req.Header.Set("Accept-Encoding", "*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", p.userAgent)
	//TODO cookies
	//fmt.Println(p.tempCookies)

	if len(p.tempCookies) > 0 {
		for _, c := range p.tempCookies {

			if c != nil {
				req.AddCookie(c)
			}

		}
	}
	resp, err := client.Do(req)

	if err != nil {
		//TODO log
		return err, nil
	}
	if resp.StatusCode != 200 {
		return errors.New("Resp code = " + strconv.Itoa(resp.StatusCode)), nil
	}
	if resp != nil {
		p.tempCookies = resp.Cookies()
	}

	return nil, resp
}

func (p *parserOnePage) saveModifayCss(doc *goquery.Document) {

	reg1 := regexp.MustCompile(`/^https:\/\/fonts.googleapis.com/`)
	reg2 := regexp.MustCompile(`/^http:\/\/ajax.googleapis.com/`)
	reg3 := regexp.MustCompile(`/url\((.*?)\)/`)
	reg4 := regexp.MustCompile(`/@import "(.*?)"/`)

	doc.Find("link[rel=stylesheet]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href != "" && !reg1.MatchString(href) && !reg2.MatchString(href) {

			_, response := p.request(href)
			if response != nil {

				bodyb, _ := ioutil.ReadAll(response.Body)
				response.Body.Close()
				name := p.filterFileName(href)

				linkCss := p.urlAbsolute(href, p.baseLink)
				linkDirCss := substrFind(linkCss, '/')
				result_replace_url := reg3.ReplaceAllStringFunc(string(bodyb), func(str string) string {
					return "url(" + p.getfileSrc(str, linkDirCss, false) + ")"
				})
				result_replace_import := reg4.ReplaceAllStringFunc(result_replace_url, func(str string) string {
					return "url(" + p.getfileSrc(str, linkDirCss, false) + ")"
				})

				if _, err := os.Stat(p.dirs["css"] + name); os.IsNotExist(err) {
					message("➤ find css :", name)
					WriteStringToFile(p.dirs["css"]+name, result_replace_import)
				}

				s.SetAttr("href", p.dirs["css_r"]+name)
			} else {
				s.SetAttr("href", " ")
			}
		}

	})

}

func (p *parserOnePage) saveModifyImg(doc *goquery.Document) map[string]string {
	name := ""
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("src")
		name = p.filterFileName(link)
		if _, err := os.Stat(p.dirs["img"] + name); os.IsNotExist(err) {
			message("find images :", name)
			if p.temp_files[name] == "" {
				name = "img-" + name
			}
			p.temp_files[name] = p.urlAbsolute(link, p.baseLink)

		}
		s.SetAttr("src", p.dirs["img_r"]+name)

	})
	return p.temp_files
}

func (p *parserOnePage) getfileSrc(link string, url string, at_html bool) string {

	newlink, name := p.srcFilter(link)
	if _, err := os.Stat(p.dirs["src"] + name); !os.IsNotExist(err) {
		//TODO logs
		return ""
	}
	message("find src for css :", newlink)
	connectLink := func() string {
		reg := regexp.MustCompile(`/^htt(p|ps):/`)
		if !reg.MatchString(newlink) {
			return p.urlAbsolute(newlink, url)
		}
		return newlink
	}
	p.temp_files_src[name] = connectLink()
	message("connecting url:", p.temp_files_src[name])

	path := ""
	if !at_html {
		path = "../../"
	}
	path += p.dirs["src_r"]
	return path + name
}

func (p *parserOnePage) srcFilter(link string) (string, string) {
	newlink := strings.Trim(link, "\"'")
	reg1 := regexp.MustCompile(`(\S+\.(png|jpg|gif|jpeg)$)`)
	reg2 := regexp.MustCompile(`/[a-zA-Z,0-9,-]+\.(ttf|svg|woff|woff2)/`)
	reg5 := regexp.MustCompile(`(^(.+)\/)|(\?.+)`)
	name := reg5.ReplaceAllStringFunc(newlink, func(st string) string {
		return ""
	})
	if !reg1.MatchString(newlink) && !reg2.MatchString(newlink) {
		name = MD5(name)
	}
	return newlink, name
}

//saveModifayCss

func (p *parserOnePage) saveFile(url string, patch string, name string) {
	res, err := http.Get(url)

	if err == nil && res.StatusCode == 200 {
		defer res.Body.Close()
		out, err := os.Create(patch + name)
		defer out.Close()
		if err == nil {
			io.Copy(out, res.Body)
		}
	}
}

func (p *parserOnePage) saveFileGo(url string, patch string, name string) {
	err, res := p.request(url)
	if err == nil {
		defer res.Body.Close()
	}
	if err == nil && res.StatusCode == 200 {
		out, err := os.Create(patch + name)
		defer out.Close()
		if err == nil {
			io.Copy(out, res.Body)
		}
	}
}

func (p *parserOnePage) save(doc *goquery.Document, filePatch string) {

	dir, _ := os.Getwd()

	file, err := os.Create(dir + "/" + p.tempName + "/" + p.indFile)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	html, err := doc.Selection.Html()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	fmt.Fprintf(file, html)

}

func (p *parserOnePage) setTempName(tempName string) *parserOnePage {

	if tempName == "" {
		panic("Not init temp name")
	}
	reg, err := regexp.Compile("/")

	if err != nil {
		log.Fatal(err)
	}

	var filter string

	filter = reg.ReplaceAllString(tempName, "-")
	p.tempName = p.dir + filter

	return p
}

func (p *parserOnePage) setCustomDir(customDir string) *parserOnePage {

	if p.tempName == "" {
		log.Fatal("Custom dir init first after temp name!!")
	}
	p.dir += customDir + "/"

	return p
}

func (p *parserOnePage) setBaseLink(baseLink string) *parserOnePage {
	p.baseLink = baseLink
	return p
}

func (p *parserOnePage) setNotIframe(is_iframe bool) *parserOnePage {
	p.notIframe = is_iframe
	return p
}

func (p *parserOnePage) setAjax(is_ajax bool) *parserOnePage {
	p.ajax = is_ajax
	return p
}

func (p *parserOnePage) setHrefAllLinks(hrefAllLinks string) *parserOnePage {
	p.hrefAllLinks = hrefAllLinks
	return p
}

func (p *parserOnePage) setOptions() {

	if p.dirAgent == "" {
		p.setAgentMob(false)
	}
	p.dirs["css"] = p.rootDir + "/" + p.tempName + "/" + p.dirAgent + "/css/"
	p.dirs["css_r"] = p.dirAgent + "/css/"
	p.dirs["js"] = p.rootDir + "/" + p.tempName + "/" + p.dirAgent + "/js/"
	p.dirs["js_r"] = p.dirAgent + "/js/"
	p.dirs["img"] = p.rootDir + "/" + p.tempName + "/" + p.dirAgent + "/images/"
	p.dirs["img_r"] = p.dirAgent + "/images/"
	p.dirs["src"] = p.rootDir + "/" + p.tempName + "/" + p.dirAgent + "/src/"
	p.dirs["src_r"] = p.dirAgent + "/src/"

	p.CreateDirIfNotExist(p.dirs["css"])
	p.CreateDirIfNotExist(p.dirs["js"])
	p.CreateDirIfNotExist(p.dirs["img"])
	p.CreateDirIfNotExist(p.dirs["src"])

	if p.dirAgent == "distr" {
		p.indFile = "index.php"
	} else {
		p.indFile = "index2.php"
	}
}

func (p *parserOnePage) CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func getInstance(conf *config.Config) *parserOnePage {

	if instance == nil {
		instance = &parserOnePage{
			timeSleep:      conf.Parser.TimeSleep,
			gol:            conf.Parser.Gol,
			temp_files_src: make(map[string]string),
			temp_files:     make(map[string]string),
			rootDir:        conf.Parser.RootDir,
			hrefAllLinks:   conf.Parser.HrefAllLinks,
			mobAgent:       conf.Parser.MobAgent,
			descAgent:      conf.Parser.DescAgent,
			links:          conf.Parser.Links,
			countLink:      conf.Parser.CountLink,
			dirs:           make(map[string]string),
			index:          conf.Parser.Index,
			script:         conf.Parser.Script,
			notIframe:      conf.Parser.NotIframe,
			ajax:           conf.Parser.Ajax,
			dir:            conf.Parser.Dir,
		}
		ch = make(chan int)
		chName = make(chan string)
		chUrl = make(chan string)
		chTypeImages = make(chan bool)
		chTypeSrc = make(chan bool)
	}

	return instance
}

func (p *parserOnePage) urlAbsolute(link string, baseLink string) string {

	reg1, _ := regexp.Compile("^htt(p|ps)://")
	reg2, _ := regexp.Compile("^//")
	reg3 := regexp.MustCompile(`(^http://[^/?#]+)?(^https://[^/?#]+)?([^?#]*)?(\?[^#]*)?(#.*)?$`)
	reg4, _ := regexp.Compile("^/")

	if reg1.MatchString(link) && reg1.MatchString(baseLink) {
		return link
	}
	if p.baseTeg != "" {
		link = p.baseTeg + link
	}
	u, err := url.Parse(baseLink)
	if err != nil {
		log.Fatal(err)
	}
	if reg2.MatchString(link) {
		return u.Scheme + ":" + link
	}

	if !reg3.MatchString(link) {
		return ""
	}
	matchesLink := reg3.FindStringSubmatch(link)
	matchesBaseLink := reg3.FindStringSubmatch(baseLink)

	if matchesLink[1] != "" {
		return link
	}
	if matchesLink[2] != "" {
		return link
	}

	if !reg3.MatchString(baseLink) {
		return ""
	}
	if matchesLink[3] == "" {
		if matchesBaseLink[1] == "" {
			return matchesBaseLink[2] + matchesBaseLink[3] + matchesBaseLink[4]
		}
		return matchesBaseLink[1] + matchesBaseLink[3] + matchesBaseLink[4]
	}

	patch := reg2.ReplaceAllString(matchesLink[3], "")
	patches := strings.Split(patch, "/")

	if patches[0] == "" {
		return u.Scheme + "://" + u.Host + matchesLink[3] + matchesLink[4]
	}

	patchBase := reg4.ReplaceAllString(matchesBaseLink[3], "")

	patchesBase := strings.Split(patchBase, "/")
	if count := len(patchesBase); count > 0 {
		patchesBase = patchesBase[:count-1]
	}

	for _, p := range patches {

		if p == "." {
			continue
		} else if p == ".." {
			if count := len(patchesBase); count > 0 {
				patchesBase = patchesBase[:count-1]
			}
		} else {
			patchesBase = append(patchesBase, p)

		}
	}
	return u.Scheme + "://" + u.Host + "/" + strings.Join(patchesBase, "/") + matchesLink[4]
}

func (p *parserOnePage) filterFileName(href string) string {
	reg := regexp.MustCompile(`(^.+\/)?(\?.+$)?`)
	return reg.ReplaceAllString(href, "")
}

func (p *parserOnePage) setAgentMob(isMob bool) *parserOnePage {
	if isMob {
		p.dirAgent = "mob"
		p.userAgent = p.mobAgent
		return p
	}
	p.dirAgent = "distr"
	p.userAgent = p.descAgent
	return p
}

func message(args ...string) {
	string := ""
	for _, str := range args {
		string += str + " "
	}
	fmt.Printf("%s\n", string)
}

func substrFind(str string, s rune) string {
	i := 0
	k := '0'
	newString := []rune(str)
	for i, k = range newString {
		if s == k {
			break
		}
	}
	return string(newString[:i])
}

func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}
