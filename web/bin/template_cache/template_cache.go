package template_cache

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const (
	TEMPLATE_TTL     = 1
	TEMPLATE_NOCACHE = false
)

type TemplateData struct {
	ctime, mtime, size int64
	tmpl               *template.Template
}

type CacheTemplate struct {
	lock           sync.RWMutex
	template_ttl   int64
	template_path  string
	template_cache map[string]*TemplateData
	function_map   template.FuncMap
}

func CreateCache(template_path string) *CacheTemplate {

	c := new(CacheTemplate)

	c.template_ttl = TEMPLATE_TTL
	c.template_path = template_path
	c.template_cache = make(map[string]*TemplateData)

	c.function_map = template.FuncMap{
		"noescape": func(s string) template.HTML { return template.HTML(s) },
	}

	return c
}

func (c *CacheTemplate) _find(name string) (data *TemplateData, ok bool) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	data, ok = c.template_cache[name]

	if ok && data.ctime < time.Now().Unix() {
		stat, _ := os.Stat(c._file(name))
		if data.mtime != stat.ModTime().Unix() || data.size != stat.Size() {
			delete(c.template_cache, name)
		} else {
			data.ctime = time.Now().Unix() + c.template_ttl
			c.template_cache[name] = data
		}
	}

	return
}

func (c *CacheTemplate) _file(name string) string {

	return c.template_path + name
}

func (c *CacheTemplate) _load(name string) (*TemplateData, error) {

	c.lock.Lock()
	defer c.lock.Unlock()

	file_name := c._file(name)

	text, _ := ioutil.ReadFile(file_name)

	html := string(text)
	html = HtmlMinify(&html)

	tmpl, err := template.New(file_name).Funcs(c.function_map).Parse(html)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := TemplateData{
		tmpl: tmpl,
	}

	stat, err := os.Stat(file_name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data.size = stat.Size()
	data.mtime = stat.ModTime().Unix()
	data.ctime = time.Now().Unix() + c.template_ttl

	c.template_cache[name] = &data

	return &data, nil
}

func (c *CacheTemplate) Get(name string) (*template.Template, bool) {

	if TEMPLATE_NOCACHE {
		data, err := c._load(name)
		if err == nil {
			return data.tmpl, true
		}
		fmt.Println(err)
		return nil, false
	}

	data, ok := c._find(name)
	if ok {
		return data.tmpl, true
	} else {
		data, err := c._load(name)
		if err == nil {
			return data.tmpl, true
		}
	}
	return nil, false
}
