package solr

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// BulkIndex take json string slice and optional params which can use to specify addition parameters such as commit=true
func (si *SolrInterface) BulkIndex(docs []string, params *url.Values) (*SolrUpdateResponse, error) {
	if si.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	if params == nil {
		params = &url.Values{}
	}
	params.Set("wt", "json")

	var b bytes.Buffer
	b.WriteString("[")
	b.WriteString(strings.Join(docs, ","))
	b.WriteString("]")

	b1 := b.Bytes()

	r, err := HTTPPost(fmt.Sprintf("%s/%s/update/?%s", si.conn.url.String(), si.conn.core, params.Encode()),
		&b1, [][]string{{"Content-Type", "application/json"}}, si.conn.username, si.conn.password, si.conn.timeout)

	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}
	// check error in resp
	if !successStatus(resp) || hasError(resp) {
		return &SolrUpdateResponse{Success: false, Result: resp}, nil
	}

	return &SolrUpdateResponse{Success: true, Result: resp}, nil
}
