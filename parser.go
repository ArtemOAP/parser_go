package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const DIR_SITE = "sites/"

type parserOnePage struct {
	timeSleep      int
	gol            int
	tempCookies    []*http.Cookie
	tempName       string
	temp_files_src map[string]string
	rootDir        string
	baseLink       string
	hrefAllLinks   string
	userAgent      string
	dirAgent       string
	mobAgent       string
	descAgent      string
	links          []string
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

var instance *parserOnePage

var ch chan int

func main() {

	parser := getInstance()

	//parser.temp_files_src["key1"] = "val1"
	//parser.links = append(parser.links, "fgfgfgfgfg", "fggfdgfdg", "fgfgfgfgfg", "fggfdgfdg", "fggfdgfdg")
	parser.setTempName("sdfsdfsdf")
	res := urlAbsolute("//fff.sdsds/sdsd/sdsdd", "https://diavita.com/page1/?q=1&w=765")
	fmt.Println(res)
	message("sdsdsd", "ffffff", "rrrrr")
	res = parser.filterFileName("//imeg.jpg?dfdf=1&f=434")
	fmt.Println(res)
	//parser.parsePage("https://toster.ru/q/431978")
	parser.baseLink = "https://toster.ru/q/431978"

	//parser.run()
	e := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch := make(chan int)
	st()
	go ex()
	for _, v := range e {
		ch <- v

	}
}
func ex() {
	d := <-ch
	fmt.Println(d)
}

func st() {

}

func (p *parserOnePage) run() *parserOnePage {

	baseLink := getInstance().baseLink
	if baseLink == "" {
		log.Fatal("Not init base links")
	}
	//TODO baseLink to one
	message("parse start", "keep", strconv.Itoa(p.countLink), "page(s)")
	p.setOptions()
	p.parsePage(baseLink)

	return p
}
func (p *parserOnePage) parsePage(link string) {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	p.thenBaseHref(doc)
	p.saveIco(doc)
	p.saveModifyElemHref(doc)
	p.saveModifyIframe(doc)
	p.saveModifyJs(doc)
	p.saveModifayCss(doc)

	if err != nil {
		log.Fatal(err)
	}

	p.save(doc, ".")
}

func (p *parserOnePage) thenBaseHref(doc *goquery.Document) {
	doc.Find("base").Remove()
}

func (p *parserOnePage) saveIco(doc *goquery.Document) {
	doc.Find("link[rel=icon],link[rel=apple-touch-icon]").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if link != "" {
			link = urlAbsolute(link, p.baseLink)
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

func (p *parserOnePage) saveModifyJs(doc *goquery.Document) {

	reg3 := regexp.MustCompile(`~^https:\/\/ajax.googleapis.com~`)

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src != "" {
			if !reg3.MatchString(src) {
				name := p.filterFileName(src)
				s.SetAttr("src", p.dirs["js_r"]+name)
				if _, err := os.Stat(p.dirs["js"] + name); os.IsNotExist(err) {

					link := urlAbsolute(src, p.baseLink)
					client := &http.Client{}
					req, err := http.NewRequest("GET", link, nil)
					if err == nil {
						req.Header.Set("Accept-language", "en")
						req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
						if len(p.tempCookies) > 0 {
							for _, c := range p.tempCookies {
								req.AddCookie(c)
							}
						}

						resp, err := client.Do(req)
						if err == nil && resp.StatusCode == 200 {
							fmt.Println(p.dir)
							defer resp.Body.Close()
							p.tempCookies = resp.Cookies()

							out, err := os.Create(p.dirs["js"] + name)
							defer out.Close()
							if err == nil {
								io.Copy(out, resp.Body)
							}
						}

					}

				}

			}
		}

	})

}

/**

 public function saveModifayCss($page)
    {
        foreach ($page->find('link[rel=stylesheet]') as $key => $css) {

            if (preg_match('/^https:\/\/fonts.googleapis.com/', $css->href) ||
                preg_match('/^http:\/\/ajax.googleapis.com/', $css->href)
            ) {
                // пропускаем <link href="https://fonts.googleapis.com/css?family=Prompt:400,400i,700&amp;subset=thai,vietnamese" rel="stylesheet"/>
                continue;
            }
            $name = $this->filterfileName($css->href);

            if ($css->href) {
                $linkCss    = $this->uri2absolute($css->href, $this->baseLink);
                $linkDirCss = substr($linkCss, 0, strrpos($linkCss, '/') + 1);

                $result_replace_url = preg_replace_callback('/url\((.*?)\)/',
                    function ($matches) use ($linkDirCss) {
                        return 'url(' . $this->getfile($matches[1], $linkDirCss) . ' )';
                    }
                    , $this->getContentFile($css->href));

                $result_replace_import = preg_replace_callback('/@import "(.*?)"/',
                    function ($matches) use ($linkDirCss) {
                        return '@import "' . $this->getfile($matches[1], $linkDirCss) . '"';
                    }, $result_replace_url);

                if (!file_exists($this->dirs['css'] . $name) && $this->getContentFile($css->href)) {
                    $this->message('➤ find css :', $name);
                    file_put_contents($this->dirs['css'] . $name, $result_replace_import);
                }

                if ($css->href) {
                    $css->href = $this->dirs['css_r'] . $name;

                }
            }
        };
    }
**/


func (p *parserOnePage) saveModifayCss(doc *goquery.Document) {

	reg1 := regexp.MustCompile(`/^https:\/\/fonts.googleapis.com/`)
	reg2 := regexp.MustCompile(`/^http:\/\/ajax.googleapis.com/`)
	reg3 := regexp.MustCompile(`/url\((.*?)\)/`)

	doc.Find("link[rel=stylesheet]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href != "" && !reg1.MatchString(href) && !reg2.MatchString(href) {

			name := p.filterFileName(href)
			linkDirCss := substrFind(href, '/')

			reg3.ReplaceAllStringFunc()
		}

	})
}

/**

 public function getfile($link, $url,$html=false)
    {

        list($newLink, $name) = $this->srcFilter($link);

        if (!file_exists($this->dirs['src'] . $name)) {
            $this->message('find src for css :',$newLink);
            $connectLink = function ($newLink) use ($url) {
                if (!preg_match('/^https:/', $newLink) && !preg_match('/^http:/', $newLink)) {
                    return $this->uri2absolute($newLink, $url);
                }
                return $newLink;
            };
            $this->temp_files_src[$name]=$connectLink($newLink);
            $this->message('connecting url:',$this->temp_files_src[$name]);
        }
        if($html){
            return '' . $this->dirs['src_r'] . $name;
        }
        return '../../' . $this->dirs['src_r'] . $name;
    }

**/

func (p *parserOnePage) getfileSrc(link string, url string, at_html bool) {

	newlink, name := p.srcFilter(link)
	if _, err := os.Stat(p.dirs["src"] + name); os.IsNotExist(err) {
		message("find src for css :", newlink)
		connectLink := func() string {
			reg := regexp.MustCompile(`/^htt(p|ps):/`)
			if !reg.MatchString(newlink) {
				return urlAbsolute(newlink, url)
			}
			return newlink
		}

	}

}

func (p *parserOnePage) srcFilter(link string) (string, string) {
	newlink := strings.Trim(link, "\"'")
	reg1 := regexp.MustCompile(`(\S+\.(png|jpg|gif|jpeg)$)`)
	reg2 := regexp.MustCompile(`/[a-zA-Z,0-9,-]+\.(ttf|svg|woff|woff2)/`)
	reg4 := regexp.MustCompile(`/(^(.+)\//)(/\?.+)/`)

	reg5 := regexp.MustCompile(`(^(.+)\/)|(\?.+)`)
	name := reg5.ReplaceAllStringFunc("sdsdsd/dadasd/qqqqq?dfdfdf", func(st string) string {
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
		fmt.Println(p.dir)
		defer res.Body.Close()
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

func getInstance() *parserOnePage {

	if instance == nil {
		instance = &parserOnePage{
			timeSleep:      10,
			gol:            2,
			temp_files_src: make(map[string]string),
			rootDir:        ".",
			hrefAllLinks:   "#",
			mobAgent:       "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
			descAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
			links:          make([]string, 0),
			countLink:      1,
			dirs:           make(map[string]string),
			index:          "index.php",
			script:         "<?php  $useragent=$_SERVER['HTTP_USER_AGENT'];   if(preg_match('/(android|bb\\d+|meego).+mobile|avantgo|bada\\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\\.(browser|link)|vodafone|wap|windows ce|xda|xiino/i',$useragent)||preg_match('/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\\-(n|u)|c55\\/|capi|ccwa|cdm\\-|cell|chtm|cldc|cmd\\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\\-s|devi|dica|dmob|do(c|p)o|ds(12|\\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\\-|_)|g1 u|g560|gene|gf\\-5|g\\-mo|go(\\.w|od)|gr(ad|un)|haie|hcit|hd\\-(m|p|t)|hei\\-|hi(pt|ta)|hp( i|ip)|hs\\-c|ht(c(\\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\\-(20|go|ma)|i230|iac( |\\-|\\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\\/)|klon|kpt |kwc\\-|kyo(c|k)|le(no|xi)|lg( g|\\/(k|l|u)|50|54|\\-[a-w])|libw|lynx|m1\\-w|m3ga|m50\\/|ma(te|ui|xo)|mc(01|21|ca)|m\\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\\-2|po(ck|rt|se)|prox|psio|pt\\-g|qa\\-a|qc(07|12|21|32|60|\\-[2-7]|i\\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\\-|oo|p\\-)|sdk\\/|se(c(\\-|0|1)|47|mc|nd|ri)|sgh\\-|shar|sie(\\-|m)|sk\\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\\-|v\\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\\-|tdg\\-|tel(i|m)|tim\\-|t\\-mo|to(pl|sh)|ts(70|m\\-|m3|m5)|tx\\-9|up(\\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\\-|your|zeto|zte\\-/i',substr($useragent,0,4)))  {      require_once ('index2.php');      die();  }  ?>",
			notIframe:      true,
			ajax:           false,
			dir:            DIR_SITE,
		}
	}

	return instance
}

func urlAbsolute(link string, baseLink string) string {
	message("links = ", link)
	message("base = ", baseLink)

	reg1, _ := regexp.Compile("^htt(p|ps)://")
	reg2, _ := regexp.Compile("^//")
	reg3 := regexp.MustCompile(`(^http://[^/?#]+)?(^https://[^/?#]+)?([^?#]*)?(\?[^#]*)?(#.*)?$`)
	reg4, _ := regexp.Compile("^/")

	if reg1.MatchString(link) && reg1.MatchString(baseLink) {
		return link
	}
	if getInstance().baseTeg != "" {
		link = getInstance().baseTeg + link
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
