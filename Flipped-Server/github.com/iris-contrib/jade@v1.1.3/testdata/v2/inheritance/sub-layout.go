// Code generated by "jade.go"; DO NOT EDIT.

package jade

import (
	pool "github.com/valyala/bytebufferpool"
)

func Jade_sublayout(title string, buffer *pool.ByteBuffer) {

	buffer.WriteString(`<html><head><title>My Site - `)
	WriteEscString(title, buffer)
	buffer.WriteString(`</title><script src="/jquery.js"></script></head><body><div class="sidebar"><p>nothing</p></div><div class="primary"><p>nothing</p></div><div id="footer"><p>some footer content</p></div></body></html>`)

}