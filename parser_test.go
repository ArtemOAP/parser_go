package main

import (
	"crypto/md5"
	"fmt"
	"parser_go/config"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_parserOnePage_run(t *testing.T) {
	tests := []struct {
		name string
		p    *parserOnePage
		want *parserOnePage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.run(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserOnePage.run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserOnePage_init(t *testing.T) {
	tests := []struct {
		name string
		p    *parserOnePage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.init()
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		str string
		del string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Split(tt.args.str, tt.args.del)
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Split() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Split() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parserOnePage_parsePage(t *testing.T) {
	conf := config.GetConfig()
	p := getInstance(conf)
	if fmt.Sprintf("%T", p) != "*main.parserOnePage" {
		t.Errorf("GetInstance not parserOnePage %s", fmt.Sprintf("%T", p))
	}
}

func Test_parserOnePage_multiFiles(t *testing.T) {
	tests := []struct {
		name string
		p    *parserOnePage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.multiFiles()
		})
	}
}

func Test_parserOnePage_appEndHeader(t *testing.T) {
	type args struct {
		doc *goquery.Document
	}
	tests := []struct {
		name string
		p    *parserOnePage
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.appEndHeader(tt.args.doc)
		})
	}
}

func Test_parserOnePage_runRemoveAllTag(t *testing.T) {
	type args struct {
		doc *goquery.Document
	}
	tests := []struct {
		name string
		p    *parserOnePage
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.runRemoveAllTag(tt.args.doc)
		})
	}
}

func Test_parserOnePage_thenBaseHref(t *testing.T) {
	type args struct {
		doc *goquery.Document
	}
	tests := []struct {
		name string
		p    *parserOnePage
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.thenBaseHref(tt.args.doc)
		})
	}
}

func Test_getInstance(t *testing.T) {
	type args struct {
		conf *config.Config
	}
	tests := []struct {
		name string
		args args
		want *parserOnePage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInstance(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserOnePage_urlAbsolute(t *testing.T) {
	conf := config.GetConfig()
	p := getInstance(conf)

	type args struct {
		link     string
		baseLink string
	}
	tests := []struct {
		name string
		p    *parserOnePage
		args args
		want string
	}{
		{"1 test", p, args{"https://example.com/dfdfdf.img", "https://example2.com/"}, "https://example.com/dfdfdf.img"},
		{"2 test", p, args{"//example.com/dfdfdf.img", "https://example2.com/"}, "https://example.com/dfdfdf.img"},
		{"3 test", p, args{"/dfdfdf.img", "https://example2.com/page1"}, "https://example2.com/dfdfdf.img"},
		{"4 test", p, args{"../../dfdfdf.img", "https://example2.com/dir1/dir2/"}, "https://example2.com/dfdfdf.img"},
		{"5 test", p, args{"dfdfdf.img", "http://example2.com/"}, "http://example2.com/dfdfdf.img"},
		{"6 test", p, args{"dir1/dir2/dir3/dfdfdf.img", "http://example2.com/"}, "http://example2.com/dir1/dir2/dir3/dfdfdf.img"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.urlAbsolute(tt.args.link, tt.args.baseLink); got != tt.want {
				t.Errorf("parserOnePage.urlAbsolute() = %v, want %v", got, tt.want)
			}
		})
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

	if fmt.Sprintf("%x", b) != res {
		t.Errorf("TestMD5 1  %s ", res)
	}

}
