// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swaggerui

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// SwaggerUIOpts configures the SwaggerUI middleware
type SwaggerUIOpts struct {
	// Host for the api
	Host string

	// BasePath for the API, defaults to: /
	BasePath string

	// Path combines with BasePath to construct the path to the UI, defaults to: "docs".
	Path string

	// SpecURL is the URL of the spec document.
	//
	// Defaults to: /swagger.json
	SpecURL string

	// Title for the documentation site, default to: API documentation
	Title string

	// Template specifies a custom template to serve the UI
	Template string

	// OAuthCallbackURL the url called after OAuth2 login
	OAuthCallbackURL string

	// The three components needed to embed swagger-ui

	// SwaggerURL points to the js that generates the SwaggerUI site.
	//
	// Defaults to: https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js
	SwaggerURL string

	SwaggerPresetURL string
	SwaggerStylesURL string

	Favicon32 string
	Favicon16 string
}

// EnsureDefaults in case some options are missing
func (r *SwaggerUIOpts) EnsureDefaults() {
	r.ensureDefaults()

	if r.Template == "" {
		r.Template = swaggeruiTemplate
	}
}

func (r *SwaggerUIOpts) EnsureDefaultsOauth2() {
	r.ensureDefaults()

	if r.Template == "" {
		r.Template = swaggerOAuthTemplate
	}
}

func (r *SwaggerUIOpts) ensureDefaults() {
	common := toCommonUIOptions(r)
	common.EnsureDefaults()
	fromCommonToAnyOptions(common, r)

	// swaggerui-specifics
	if r.OAuthCallbackURL == "" {
		r.OAuthCallbackURL = path.Join(r.BasePath, r.Path, "oauth2-callback")
	}
	if r.SwaggerURL == "" {
		r.SwaggerURL = swaggerLatest
	}
	if r.SwaggerPresetURL == "" {
		r.SwaggerPresetURL = swaggerPresetLatest
	}
	if r.SwaggerStylesURL == "" {
		r.SwaggerStylesURL = swaggerStylesLatest
	}
	if r.Favicon16 == "" {
		r.Favicon16 = swaggerFavicon16Latest
	}
	if r.Favicon32 == "" {
		r.Favicon32 = swaggerFavicon32Latest
	}
}

// SwaggerUI creates a middleware to serve a documentation site for a swagger spec.
//
// This allows for altering the spec before starting the http listener.
func SwaggerUI(opts SwaggerUIOpts, next http.Handler) http.Handler {
	opts.EnsureDefaults()

	pth := path.Join(opts.BasePath, opts.Path)
	tmpl := template.Must(template.New("swaggerui").Parse(opts.Template))
	assets := bytes.NewBuffer(nil)
	if err := tmpl.Execute(assets, opts); err != nil {
		panic(fmt.Errorf("cannot execute template: %w", err))
	}

	return serveUI(pth, assets.Bytes(), next)
}

const (
	swaggerLatest          = "https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"
	swaggerPresetLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"
	swaggerStylesLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui.css"
	swaggerFavicon32Latest = "https://unpkg.com/swagger-ui-dist/favicon-32x32.png"
	swaggerFavicon16Latest = "https://unpkg.com/swagger-ui-dist/favicon-16x16.png"
	swaggeruiTemplate      = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
		<title>{{ .Title }}</title>

    <link rel="stylesheet" type="text/css" href="{{ .SwaggerStylesURL }}" >
    <link rel="icon" type="image/png" href="{{ .Favicon32 }}" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .Favicon16 }}" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="{{ .SwaggerURL }}"> </script>
    <script src="{{ .SwaggerPresetURL }}"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
		host: "{{ .Host }}",
        url: '{{ .SpecURL }}',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
		oauth2RedirectUrl: '{{ .OAuthCallbackURL }}'
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`
)

const (
	swaggerOAuthTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{ .Title }}</title>
</head>
<body>
<script>
    'use strict';
    function run () {
        var oauth2 = window.opener.swaggerUIRedirectOauth2;
        var sentState = oauth2.state;
        var redirectUrl = oauth2.redirectUrl;
        var isValid, qp, arr;

        if (/code|token|error/.test(window.location.hash)) {
            qp = window.location.hash.substring(1).replace('?', '&');
        } else {
            qp = location.search.substring(1);
        }

        arr = qp.split("&");
        arr.forEach(function (v,i,_arr) { _arr[i] = '"' + v.replace('=', '":"') + '"';});
        qp = qp ? JSON.parse('{' + arr.join() + '}',
                function (key, value) {
                    return key === "" ? value : decodeURIComponent(value);
                }
        ) : {};

        isValid = qp.state === sentState;

        if ((
          oauth2.auth.schema.get("flow") === "accessCode" ||
          oauth2.auth.schema.get("flow") === "authorizationCode" ||
          oauth2.auth.schema.get("flow") === "authorization_code"
        ) && !oauth2.auth.code) {
            if (!isValid) {
                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "warning",
                    message: "Authorization may be unsafe, passed state was changed in server. The passed state wasn't returned from auth server."
                });
            }

            if (qp.code) {
                delete oauth2.state;
                oauth2.auth.code = qp.code;
                oauth2.callback({auth: oauth2.auth, redirectUrl: redirectUrl});
            } else {
                let oauthErrorMsg;
                if (qp.error) {
                    oauthErrorMsg = "["+qp.error+"]: " +
                        (qp.error_description ? qp.error_description+ ". " : "no accessCode received from the server. ") +
                        (qp.error_uri ? "More info: "+qp.error_uri : "");
                }

                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "error",
                    message: oauthErrorMsg || "[Authorization failed]: no accessCode received from the server."
                });
            }
        } else {
            oauth2.callback({auth: oauth2.auth, token: qp, isValid: isValid, redirectUrl: redirectUrl});
        }
        window.close();
    }

    if (document.readyState !== 'loading') {
        run();
    } else {
        document.addEventListener('DOMContentLoaded', function () {
            run();
        });
    }
</script>
</body>
</html>
`
)
