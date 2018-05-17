package main

import (

	"testing"
)


func (p *parserOnePage) Testrun(t *testing.T) *parserOnePage {

	
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

func (p *parserOnePage) Testrequest(t *testing.T)  {

}

func (p *parserOnePage) TestsaveModifayCss(t *testing.T) {

	
}

func (p *parserOnePage) TestsaveModifyImg(t *testing.T)  {
	
}

func (p *parserOnePage) TestgetfileSrc(t *testing.T)  {


}

func (p *parserOnePage) TestsrcFilter(t *testing.T)  {

}

//saveModifayCss

func (p *parserOnePage) TestsaveFile(t *testing.T) {
	
}

func (p *parserOnePage) TestsaveFileGo(t *testing.T) {
	
}

func (p *parserOnePage) save(doc *goquery.Document, filePatch string) {

	//dir, _ := os.Getwd()

	file, err := os.Create(filePatch)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	html, err := doc.Selection.Html()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	fmt.Fprint(file, html)

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

func (p *parserOnePage) TestsaveCookies(t *testing.T) {

	
}

func (p *parserOnePage) TestsetCustomDir(t *testing.T) *parserOnePage {


}

func (p *parserOnePage) TestsetBaseLink(t *testing.T) *parserOnePage {
	
}

func (p *parserOnePage) TestsetNotIframe(t *testing.T) *parserOnePage {
	
}

func (p *parserOnePage) TestsetAjax(t *testing.T) *parserOnePage {
	
}

func (p *parserOnePage) TestsetHrefAllLinks(t *testing.T) *parserOnePage {
	
}

func (p *parserOnePage) TestsetOptions(t *testing.T) {


}

func (p *parserOnePage) TestCreateDirIfNotExist(t *testing.T) {

}

func TestgetInstance(t *testing.T) {


}

func (p *parserOnePage) TesturlAbsolute(t *testing.T)  {

	
}

func (p *parserOnePage) TestfilterFileName(t *testing.T) {



}

func (p *parserOnePage) TestsetAgentMob(t *testing.T)  {
	
}



func TestsubstrFind(t *testing.T)  {

}

func TestMD5(t *testing.T)  {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func TestWriteStringToFile(t *testing.T)  {
	
}
