package main

import (
	"fmt"
	"crypto/md5"
	"parser_go/config"
	"testing"
)

func (p *parserOnePage) Testrun(t *testing.T) {

}

func (p *parserOnePage) Testinit(t *testing.T) {

}

func TestSplit(t *testing.T) {

}

func (p *parserOnePage) TestparsePage(t *testing.T) {

}
func (p *parserOnePage) TestmultiFiles(t *testing.T) {

}

func (p *parserOnePage) TestappEndHeader(t *testing.T) {

}

func (p *parserOnePage) TestrunRemoveAllTag(t *testing.T) {

}

func (p *parserOnePage) TestthenBaseHref(t *testing.T) {

}

func (p *parserOnePage) TestsaveIco(t *testing.T) {

}

func (p *parserOnePage) TestsaveModifyElemHref(t *testing.T) {

}

func (p *parserOnePage) TestsaveModifyIframe(t *testing.T) {

}
func (p *parserOnePage) TestmodifyForm(t *testing.T) {

}

//TODO test
func (p *parserOnePage) TestreplaceCssInHtml(t *testing.T) {

}

func (p *parserOnePage) TestsaveModifyJs(t *testing.T) {

}

func (p *parserOnePage) Testrequest(t *testing.T) {

}

func (p *parserOnePage) TestsaveModifayCss(t *testing.T) {

}

func (p *parserOnePage) TestsaveModifyImg(t *testing.T) {

}

func (p *parserOnePage) TestgetfileSrc(t *testing.T) {

}

func (p *parserOnePage) TestsrcFilter(t *testing.T) {

}

//saveModifayCss

func (p *parserOnePage) TestsaveFile(t *testing.T) {

}

func (p *parserOnePage) TestsaveFileGo(t *testing.T) {

}

func (p *parserOnePage) Testsave(t *testing.T) {

}

func (p *parserOnePage) TestsetTempName(t *testing.T) {

}

func (p *parserOnePage) TestsaveCookies(t *testing.T) {

}

func (p *parserOnePage) TestsetCustomDir(t *testing.T) {

}

func (p *parserOnePage) TestsetBaseLink(t *testing.T) {

}

func (p *parserOnePage) TestsetNotIframe(t *testing.T) {

}

func (p *parserOnePage) TestsetAjax(t *testing.T) {

}

func (p *parserOnePage) TestsetHrefAllLinks(t *testing.T) {

}

func (p *parserOnePage) TestsetOptions(t *testing.T) {

}

func (p *parserOnePage) TestCreateDirIfNotExist(t *testing.T) {

}

func TestGetInstance(t *testing.T) {
	conf := config.GetConfig()
	p := getInstance(conf) 
	if fmt.Sprintf("%T",p) != "*main.parserOnePage"{
		t.Errorf("GetInstance not parserOnePage %s",fmt.Sprintf("%T",p))
	}
}

func TestUrlAbsolute(t *testing.T) {

	var res string
	conf := config.GetConfig()
	p := getInstance(conf)

	res = p.urlAbsolute("https://example.com/dfdfdf.img", "https://example2.com/")
	if res != "https://example.com/dfdfdf.img" {
		t.Errorf("urlAbsolute 1 not correct -%s", res)
	}
	res = p.urlAbsolute("//example.com/dfdfdf.img", "https://example2.com/")
	if res != "https://example.com/dfdfdf.img" {
		t.Errorf("urlAbsolute 2 not correct -%s", res)
	}
	res = p.urlAbsolute("/dfdfdf.img", "https://example2.com/page1")
	if res != "https://example2.com/dfdfdf.img" {
		t.Errorf("urlAbsolute 3 not correct -%s", res)
	}
	res = p.urlAbsolute("../../dfdfdf.img", "https://example2.com/dir1/dir2/")
	if res != "https://example2.com/dfdfdf.img" {
		t.Errorf("urlAbsolute 4 not correct -%s", res)
	}
	res = p.urlAbsolute("dfdfdf.img", "http://example2.com/")
	if res != "http://example2.com/dfdfdf.img" {
		t.Errorf("urlAbsolute 5 not correct -%s", res)
	}
	res = p.urlAbsolute("dir1/dir2/dir3/dfdfdf.img", "http://example2.com/")
	if res != "http://example2.com/dir1/dir2/dir3/dfdfdf.img" {
		t.Errorf("urlAbsolute 6 not correct -%s", res)
	}

}

func TestFilterFileName(t *testing.T) {

	var res string
	p := getInstance(config.GetConfig())

	res = p.filterFileName("http://example2.com/dir1/dir2/dir3/dfdfdf.img")
	if res != "dfdfdf.img" {
		t.Errorf("FilterFileName 1 not correct -%s", res)
	}
	res = p.filterFileName("http://example2.com/dir1/dir2/dir3/dfdfdf.img?ddd=sds&dsd=43545")
	if res != "dfdfdf.img" {
		t.Errorf("FilterFileName 2 not correct -%s", res)
	}
	res = p.filterFileName("//example2.com/dir1/dir2/dir3/dfdfdf.img?ddd=sds&dsd=43545")
	if res != "dfdfdf.img" {
		t.Errorf("FilterFileName 3 not correct -%s", res)
	}
	res = p.filterFileName("//example2.com/dir1/dir2/dir3/?ddd=sds&dsd=43545")
	if res == "" {
		t.Errorf("FilterFileName 4 empty")
	}

	res = p.filterFileName("//example2.com/dir1/dir2/dir3/img1.JPG?ddd=sds&dsd=43545")
	p.temp_files[res] = "//example2.com/dir1/dir2/dir3/img1.JPG?ddd=sds&dsd=43545"
	if res != "img1.JPG" {
		t.Errorf("FilterFileName 5  %s", res)
	}

	res = p.filterFileName("//example2.com/dir1/dir2/dir3/img1.JPG?ddd=sds&dsd=43545")
	if res[0:4] != "img-" {
		t.Errorf("FilterFileName 6  %s", res)
	}

}

func TestSetAgentMob(t *testing.T) {
	p := getInstance(config.GetConfig())
	p.setAgentMob(true)
	if p.dirAgent != "mob" || p.userAgent != "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30" {
		t.Errorf("AgentMob 1  %s %s", p.dirAgent, p.userAgent)
	}
	p.setAgentMob(false)
	if p.dirAgent != "distr" || p.userAgent != "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36" {
		t.Errorf("AgentMob 2  %s %s", p.dirAgent, p.userAgent)
	}

}

func TestSubstrFind(t *testing.T) {

	var res string
	res = substrFind("find_text|dfdfdfdfdfdfgfgfg", '|')
	if res != "find_text" {
		t.Errorf("SubstrFind 1  %s", res)
	}
	res = substrFind("find_text|dfdfdfdf|dfdfgfgfg", '|')
	if res != "find_text" {
		t.Errorf("SubstrFind 1  %s", res)
	}

}

func TestMD5(t *testing.T) {
	res := MD5("text")
	b := md5.Sum([]byte("text"))

		if fmt.Sprintf("%x",b)!= res{
			t.Errorf("TestMD5 1  %s ", res)
		} 
	
	
}

// func TestWriteStringToFile(t *testing.T) {

// }
