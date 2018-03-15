package main

import (
	"fmt"
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

type Parser_one_page struct {
	time_sleep     int
	gol            int
	temp_cookies   string
	temp_name      string
	temp_files_src map[string]string
	_root_dir      string
	baseLink       string
	hrefAllLinks   string
	user_agent     string
	dir_agent      string
	mob_agent      string
	desc_agent     string
	links          []string
	countLink      int
	dirs           map[string]string
	indFile        string
	index          string
	script         string
	not_iframe     bool
	baseTeg        string
	ajax           bool
	dir            string
}

var instance *Parser_one_page

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
	parser.parsePage("https://toster.ru/q/431978")

}

func (p *Parser_one_page) run() *Parser_one_page {

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
func (p *Parser_one_page) parsePage(link string) {

	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		fmt.Println(link)
	})

	p.save(doc, ".")
}

/**
  function save($filepath='')
    {
        $ret = $this->root->innertext();
        if ($filepath!=='') file_put_contents($filepath, $ret, LOCK_EX);
        return $ret;
    }
**/

func (p *Parser_one_page) save(doc *goquery.Document, filePatch string) {

	file, err := os.Create("result.txt")
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

func (p *Parser_one_page) setTempName(tempName string) *Parser_one_page {

	if tempName == "" {
		panic("Not init temp name")
	}
	reg, err := regexp.Compile("/")

	if err != nil {
		log.Fatal(err)
	}

	var filter string

	filter = reg.ReplaceAllString(tempName, "-")
	p.temp_name = p.dir + filter

	return p
}

func (p *Parser_one_page) setCustomDir(customDir string) *Parser_one_page {

	if p.temp_name == "" {
		log.Fatal("Custom dir init first after temp name!!")
	}
	p.dir += customDir + "/"

	return p
}

func (p *Parser_one_page) setBaseLink(baseLink string) *Parser_one_page {
	p.baseLink = baseLink
	return p
}

func (p *Parser_one_page) setNotIframe(is_iframe bool) *Parser_one_page {
	p.not_iframe = is_iframe
	return p
}

func (p *Parser_one_page) setAjax(is_ajax bool) *Parser_one_page {
	p.ajax = is_ajax
	return p
}

func (p *Parser_one_page) setHrefAllLinks(hrefAllLinks string) *Parser_one_page {
	p.hrefAllLinks = hrefAllLinks
	return p
}

func (p *Parser_one_page) setOptions() {

	p.dirs["css"] = p._root_dir + "/" + p.temp_name + "/" + p.dir_agent + "/css/"
	p.dirs["css_r"] = p.dir_agent + "/css/"
	p.dirs["js"] = p._root_dir + "/" + p.temp_name + "/" + p.dir_agent + "/js/"
	p.dirs["js_r"] = p.dir_agent + "/js/"
	p.dirs["img"] = p._root_dir + "/" + p.temp_name + "/" + p.dir_agent + "/images/"
	p.dirs["img_r"] = p.dir_agent + "/images/"
	p.dirs["src"] = p._root_dir + "/" + p.temp_name + "/" + p.dir_agent + "/src/"
	p.dirs["src_r"] = p.dir_agent + "/src/"

	p.CreateDirIfNotExist(p.dirs["css"])
	p.CreateDirIfNotExist(p.dirs["js"])
	p.CreateDirIfNotExist(p.dirs["img"])
	p.CreateDirIfNotExist(p.dirs["src"])

	if p.dir_agent == "distr" {
		p.indFile = "index.php"
	} else {
		p.indFile = "index2.php"
	}
}

func (p *Parser_one_page) CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func getInstance() *Parser_one_page {

	if instance == nil {
		instance = &Parser_one_page{
			time_sleep:     10,
			gol:            2,
			temp_files_src: make(map[string]string),
			_root_dir:      ".",
			hrefAllLinks:   "#",
			mob_agent:      "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
			desc_agent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
			links:          make([]string, 0),
			countLink:      1,
			dirs:           make(map[string]string),
			index:          "index.php",
			script:         "<?php  $useragent=$_SERVER['HTTP_USER_AGENT'];   if(preg_match('/(android|bb\\d+|meego).+mobile|avantgo|bada\\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\\.(browser|link)|vodafone|wap|windows ce|xda|xiino/i',$useragent)||preg_match('/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\\-(n|u)|c55\\/|capi|ccwa|cdm\\-|cell|chtm|cldc|cmd\\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\\-s|devi|dica|dmob|do(c|p)o|ds(12|\\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\\-|_)|g1 u|g560|gene|gf\\-5|g\\-mo|go(\\.w|od)|gr(ad|un)|haie|hcit|hd\\-(m|p|t)|hei\\-|hi(pt|ta)|hp( i|ip)|hs\\-c|ht(c(\\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\\-(20|go|ma)|i230|iac( |\\-|\\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\\/)|klon|kpt |kwc\\-|kyo(c|k)|le(no|xi)|lg( g|\\/(k|l|u)|50|54|\\-[a-w])|libw|lynx|m1\\-w|m3ga|m50\\/|ma(te|ui|xo)|mc(01|21|ca)|m\\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\\-2|po(ck|rt|se)|prox|psio|pt\\-g|qa\\-a|qc(07|12|21|32|60|\\-[2-7]|i\\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\\-|oo|p\\-)|sdk\\/|se(c(\\-|0|1)|47|mc|nd|ri)|sgh\\-|shar|sie(\\-|m)|sk\\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\\-|v\\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\\-|tdg\\-|tel(i|m)|tim\\-|t\\-mo|to(pl|sh)|ts(70|m\\-|m3|m5)|tx\\-9|up(\\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\\-|your|zeto|zte\\-/i',substr($useragent,0,4)))  {      require_once ('index2.php');      die();  }  ?>",
			not_iframe:     true,
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

func (p *Parser_one_page) filterFileName(href string) string {
	reg := regexp.MustCompile(`(^.+\/)?(\?.+$)?`)
	return reg.ReplaceAllString(href, "")
}

func (p *Parser_one_page) setAgentMob(isMob bool) *Parser_one_page {
	if isMob {
		p.mob_agent = "mob"
		p.user_agent = p.mob_agent
		return p
	}
	p.mob_agent = "distr"
	p.user_agent = p.desc_agent
	return p
}

func message(args ...string) {
	string := ""
	for _, str := range args {
		string += str + " "
	}
	fmt.Printf("%s\n", string)
}