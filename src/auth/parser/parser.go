/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"configcenter/src/auth/meta"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"github.com/emicklei/go-restful"
)

func ParseAttribute(req *restful.Request) (*meta.AuthAttribute, error) {

	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}

	meta := new(metadata.Metadata)
	if err := json.Unmarshal(body, meta); err != nil {
		return nil, err
	}

	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return nil, err
	}

	requestContext := &RequestContext{
		Header:   req.Request.Header,
		Method:   req.Request.Method,
		URI:      req.Request.URL.Path,
		Elements: elements,
		Body:     body,
		Metadata: *meta,
	}

	stream, err := newParseStream(requestContext)
	if err != nil {
		return nil, err
	}

	return stream.Parse()
}

// url example: /api/v3/create/model
var urlRegex = regexp.MustCompile(`^/api/([^/]+)/([^/]+)/([^/]+)/(.*)$`)

func urlParse(url string) (elements []string, err error) {
	if !urlRegex.MatchString(url) {
		return nil, errors.New("invalid url format")
	}

	return strings.Split(url, "/")[1:], nil
}

func filterAction(action string) (meta.Action, error) {
	switch action {
	case "find":
		return meta.Find, nil
	case "findMany":
		return meta.FindMany, nil

	case "create":
		return meta.Create, nil
	case "createMany":
		return meta.CreateMany, nil

	case "update":
		return meta.Update, nil
	case "updateMany":
		return meta.UpdateMany, nil

	case "delete":
		return meta.Delete, nil
	case "deleteMany":
		return meta.DeleteMany, nil

	default:
		return meta.Unknown, fmt.Errorf("unsupported action %s", action)
	}
}
