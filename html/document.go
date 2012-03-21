package html

/*
#cgo pkg-config: libxml-2.0

#include <libxml/HTMLtree.h>
#include <libxml/HTMLparser.h>
#include "helper.h"
*/
import "C"

import (
	"unsafe"
	"os"
	"gokogiri/xml"
	. "gokogiri/util"
)

//xml parse option
const (
	HTML_PARSE_RECOVER   = 1 << 0  /* Relaxed parsing */
	HTML_PARSE_NODEFDTD  = 1 << 2  /* do not default a doctype if not found */
	HTML_PARSE_NOERROR   = 1 << 5  /* suppress error reports */
	HTML_PARSE_NOWARNING = 1 << 6  /* suppress warning reports */
	HTML_PARSE_PEDANTIC  = 1 << 7  /* pedantic error reporting */
	HTML_PARSE_NOBLANKS  = 1 << 8  /* remove blank nodes */
	HTML_PARSE_NONET     = 1 << 11 /* Forbid network access */
	HTML_PARSE_NOIMPLIED = 1 << 13 /* Do not add implied html/body... elements */
	HTML_PARSE_COMPACT   = 1 << 16 /* compact small text nodes */
)

const EmptyHtmlDoc = ""

//default parsing option: relax parsing
var DefaultParseOption = HTML_PARSE_RECOVER |
	HTML_PARSE_NONET |
	HTML_PARSE_NOERROR |
	HTML_PARSE_NOWARNING

type HtmlDocument struct {
	*xml.XmlDocument
}

//default encoding in byte slice
var DefaultEncodingBytes = []byte(xml.DefaultEncoding)
var emptyHtmlDocBytes = []byte(EmptyHtmlDoc)

var ErrSetMetaEncoding = os.NewError("Set Meta Encoding failed")
var ERR_FAILED_TO_PARSE_HTML = os.NewError("failed to parse html input")
var emptyStringBytes = []byte{0}

//create a document
func NewDocument(p unsafe.Pointer, contentLen int, inEncoding, outEncoding []byte) (doc *HtmlDocument) {
	doc = &HtmlDocument{}
	doc.XmlDocument = xml.NewDocument(p, contentLen, inEncoding, outEncoding)
	node := doc.Node.(*xml.XmlNode)
	node.Document = doc
	return
}

//parse a string to document
func Parse(content, inEncoding, url []byte, options int, outEncoding []byte) (doc *HtmlDocument, err os.Error) {
	inEncoding  = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)

	var docPtr *C.xmlDoc
	contentLen := len(content)

	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer

		contentPtr = unsafe.Pointer(&content[0])
		if len(url) > 0 {
			url = AppendCStringTerminator(url)
			urlPtr = unsafe.Pointer(&url[0])
		}
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}

		docPtr = C.htmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		if docPtr == nil {
			err = ERR_FAILED_TO_PARSE_HTML
		} else {
			doc = NewDocument(unsafe.Pointer(docPtr), contentLen, inEncoding, outEncoding)
		}
	}
	if docPtr == nil {
		doc = CreateEmptyDocument(inEncoding, outEncoding)
	}
	return
}

func CreateEmptyDocument(inEncoding, outEncoding []byte) (doc *HtmlDocument) {
	C.xmlInitParser()
	docPtr := C.htmlNewDoc(nil, nil)
	doc = NewDocument(unsafe.Pointer(docPtr), 0, inEncoding, outEncoding)
	return
}

func (document *HtmlDocument) ParseFragment(input, url []byte, options int) (fragment *xml.DocumentFragment, err os.Error) {
	fragment, err = parsefragmentInDocument(document, input, url, options)
	return
}

func (doc *HtmlDocument) MetaEncoding() string {
	metaEncodingXmlCharPtr := C.htmlGetMetaEncoding((*C.xmlDoc)(doc.DocPtr()))
	return C.GoString((*C.char)(unsafe.Pointer(metaEncodingXmlCharPtr)))
}

func (doc *HtmlDocument) SetMetaEncoding(encoding string) (err os.Error) {
	var encodingPtr unsafe.Pointer = nil
	if len(encoding) > 0 {
		encodingBytes := AppendCStringTerminator([]byte(encoding))
		encodingPtr = unsafe.Pointer(&encodingBytes[0])
	}
	ret := int(C.htmlSetMetaEncoding((*C.xmlDoc)(doc.DocPtr()), (*C.xmlChar)(encodingPtr)))
	if ret == -1 {
		err = ErrSetMetaEncoding
	}
	return
}

func (document *HtmlDocument) Root() (element *xml.ElementNode) {
	p := unsafe.Pointer(document.Ptr)
	nodePtr := C.xmlDocGetRootElement((*C.xmlDoc)(p))
	element = xml.NewNode(unsafe.Pointer(nodePtr), document).(*xml.ElementNode)
	return
}

func (document *HtmlDocument) CreateElementNode(tag string) (element *xml.ElementNode) {
	tagBytes := GetCString([]byte(tag))
	tagPtr := unsafe.Pointer(&tagBytes[0])
	newNodePtr := C.xmlNewNode(nil, (*C.xmlChar)(tagPtr))
	newNode := xml.NewNode(unsafe.Pointer(newNodePtr), document)
	element = newNode.(*xml.ElementNode)
	return
}

func (document *HtmlDocument) CreateCData(data string) (cdata *xml.CDataNode) {
	var dataPtr unsafe.Pointer
	dataLen := len(data)
	if dataLen > 0 {
		dataBytes := []byte(data)
		dataPtr = unsafe.Pointer(&dataBytes[0])
	} else {
		dataPtr = unsafe.Pointer(&emptyStringBytes[0])
	}
	p := unsafe.Pointer(document.Ptr)
	nodePtr := C.xmlNewCDataBlock((*C.xmlDoc)(p), (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = xml.NewNode(unsafe.Pointer(nodePtr), document).(*xml.CDataNode)
	}
	return
}
