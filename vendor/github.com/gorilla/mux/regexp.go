// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// newRouteRegexp parses a route template and returns a routeRegexp,
// used to match a host, a path or a query string.
//
// It will extract named variables, assemble a regexp to be matched, create
// a "reverse" template to build URLs and compile regexps to validate variable
// values used in URL building.
//
// Previously we accepted only Python-like identifiers for variable
// names ([a-zA-Z_][a-zA-Z0-9_]*), but currently the only restriction is that
// name and pattern can't be empty, and names can't contain a colon.
<<<<<<< HEAD
func newRouteRegexp(tpl string, matchHost, matchPrefix, matchQuery, strictSlash bool) (*routeRegexp, error) {
=======
func newRouteRegexp(tpl string, matchHost, matchPrefix, matchQuery, strictSlash, useEncodedPath bool) (*routeRegexp, error) {
>>>>>>> a24f8b0... add support for installing helm charts
	// Check if it is well-formed.
	idxs, errBraces := braceIndices(tpl)
	if errBraces != nil {
		return nil, errBraces
	}
	// Backup the original.
	template := tpl
	// Now let's parse it.
	defaultPattern := "[^/]+"
	if matchQuery {
		defaultPattern = "[^?&]*"
	} else if matchHost {
		defaultPattern = "[^.]+"
		matchPrefix = false
	}
	// Only match strict slash if not matching
	if matchPrefix || matchHost || matchQuery {
		strictSlash = false
	}
	// Set a flag for strictSlash.
	endSlash := false
	if strictSlash && strings.HasSuffix(tpl, "/") {
		tpl = tpl[:len(tpl)-1]
		endSlash = true
	}
	varsN := make([]string, len(idxs)/2)
	varsR := make([]*regexp.Regexp, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')
	reverse := bytes.NewBufferString("")
	var end int
	var err error
	for i := 0; i < len(idxs); i += 2 {
		// Set all values we are interested in.
		raw := tpl[end:idxs[i]]
		end = idxs[i+1]
		parts := strings.SplitN(tpl[idxs[i]+1:end-1], ":", 2)
		name := parts[0]
		patt := defaultPattern
		if len(parts) == 2 {
			patt = parts[1]
		}
		// Name or pattern can't be empty.
		if name == "" || patt == "" {
			return nil, fmt.Errorf("mux: missing name or pattern in %q",
				tpl[idxs[i]:end])
		}
		// Build the regexp pattern.
<<<<<<< HEAD
		varIdx := i / 2
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(varIdx), patt)
=======
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(i/2), patt)

>>>>>>> a24f8b0... add support for installing helm charts
		// Build the reverse template.
		fmt.Fprintf(reverse, "%s%%s", raw)

		// Append variable name and compiled pattern.
<<<<<<< HEAD
		varsN[varIdx] = name
		varsR[varIdx], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
=======
		varsN[i/2] = name
		varsR[i/2], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
>>>>>>> a24f8b0... add support for installing helm charts
		if err != nil {
			return nil, err
		}
	}
	// Add the remaining.
	raw := tpl[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	if strictSlash {
		pattern.WriteString("[/]?")
	}
	if matchQuery {
		// Add the default pattern if the query value is empty
		if queryVal := strings.SplitN(template, "=", 2)[1]; queryVal == "" {
			pattern.WriteString(defaultPattern)
		}
	}
	if !matchPrefix {
		pattern.WriteByte('$')
	}
	reverse.WriteString(raw)
	if endSlash {
		reverse.WriteByte('/')
	}
	// Compile full regexp.
	reg, errCompile := regexp.Compile(pattern.String())
	if errCompile != nil {
		return nil, errCompile
	}
<<<<<<< HEAD
	// Done!
	return &routeRegexp{
		template:    template,
		matchHost:   matchHost,
		matchQuery:  matchQuery,
		strictSlash: strictSlash,
		regexp:      reg,
		reverse:     reverse.String(),
		varsN:       varsN,
		varsR:       varsR,
=======

	// Check for capturing groups which used to work in older versions
	if reg.NumSubexp() != len(idxs)/2 {
		panic(fmt.Sprintf("route %s contains capture groups in its regexp. ", template) +
			"Only non-capturing groups are accepted: e.g. (?:pattern) instead of (pattern)")
	}

	// Done!
	return &routeRegexp{
		template:       template,
		matchHost:      matchHost,
		matchQuery:     matchQuery,
		strictSlash:    strictSlash,
		useEncodedPath: useEncodedPath,
		regexp:         reg,
		reverse:        reverse.String(),
		varsN:          varsN,
		varsR:          varsR,
>>>>>>> a24f8b0... add support for installing helm charts
	}, nil
}

// routeRegexp stores a regexp to match a host or path and information to
// collect and validate route variables.
type routeRegexp struct {
	// The unmodified template.
	template string
	// True for host match, false for path or query string match.
	matchHost bool
	// True for query string match, false for path and host match.
	matchQuery bool
	// The strictSlash value defined on the route, but disabled if PathPrefix was used.
	strictSlash bool
<<<<<<< HEAD
=======
	// Determines whether to use encoded path from getPath function or unencoded
	// req.URL.Path for path matching
	useEncodedPath bool
>>>>>>> a24f8b0... add support for installing helm charts
	// Expanded regexp.
	regexp *regexp.Regexp
	// Reverse template.
	reverse string
	// Variable names.
	varsN []string
	// Variable regexps (validators).
	varsR []*regexp.Regexp
}

// Match matches the regexp against the URL host or path.
func (r *routeRegexp) Match(req *http.Request, match *RouteMatch) bool {
	if !r.matchHost {
		if r.matchQuery {
			return r.matchQueryString(req)
<<<<<<< HEAD
		} else {
			return r.regexp.MatchString(req.URL.Path)
		}
	}
=======
		}
		path := req.URL.Path
		if r.useEncodedPath {
			path = getPath(req)
		}
		return r.regexp.MatchString(path)
	}

>>>>>>> a24f8b0... add support for installing helm charts
	return r.regexp.MatchString(getHost(req))
}

// url builds a URL part using the given values.
func (r *routeRegexp) url(values map[string]string) (string, error) {
	urlValues := make([]interface{}, len(r.varsN))
	for k, v := range r.varsN {
		value, ok := values[v]
		if !ok {
			return "", fmt.Errorf("mux: missing route variable %q", v)
		}
		urlValues[k] = value
	}
	rv := fmt.Sprintf(r.reverse, urlValues...)
	if !r.regexp.MatchString(rv) {
		// The URL is checked against the full regexp, instead of checking
		// individual variables. This is faster but to provide a good error
		// message, we check individual regexps if the URL doesn't match.
		for k, v := range r.varsN {
			if !r.varsR[k].MatchString(values[v]) {
				return "", fmt.Errorf(
					"mux: variable %q doesn't match, expected %q", values[v],
					r.varsR[k].String())
			}
		}
	}
	return rv, nil
}

<<<<<<< HEAD
// getUrlQuery returns a single query parameter from a request URL.
// For a URL with foo=bar&baz=ding, we return only the relevant key
// value pair for the routeRegexp.
func (r *routeRegexp) getUrlQuery(req *http.Request) string {
=======
// getURLQuery returns a single query parameter from a request URL.
// For a URL with foo=bar&baz=ding, we return only the relevant key
// value pair for the routeRegexp.
func (r *routeRegexp) getURLQuery(req *http.Request) string {
>>>>>>> a24f8b0... add support for installing helm charts
	if !r.matchQuery {
		return ""
	}
	templateKey := strings.SplitN(r.template, "=", 2)[0]
	for key, vals := range req.URL.Query() {
		if key == templateKey && len(vals) > 0 {
			return key + "=" + vals[0]
		}
	}
	return ""
}

func (r *routeRegexp) matchQueryString(req *http.Request) bool {
<<<<<<< HEAD
	return r.regexp.MatchString(r.getUrlQuery(req))
=======
	return r.regexp.MatchString(r.getURLQuery(req))
>>>>>>> a24f8b0... add support for installing helm charts
}

// braceIndices returns the first level curly brace indices from a string.
// It returns an error in case of unbalanced braces.
func braceIndices(s string) ([]int, error) {
	var level, idx int
<<<<<<< HEAD
	idxs := make([]int, 0)
=======
	var idxs []int
>>>>>>> a24f8b0... add support for installing helm charts
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
	}
	return idxs, nil
}

// varGroupName builds a capturing group name for the indexed variable.
func varGroupName(idx int) string {
	return "v" + strconv.Itoa(idx)
}

// ----------------------------------------------------------------------------
// routeRegexpGroup
// ----------------------------------------------------------------------------

// routeRegexpGroup groups the route matchers that carry variables.
type routeRegexpGroup struct {
	host    *routeRegexp
	path    *routeRegexp
	queries []*routeRegexp
}

// setMatch extracts the variables from the URL once a route matches.
func (v *routeRegexpGroup) setMatch(req *http.Request, m *RouteMatch, r *Route) {
	// Store host variables.
	if v.host != nil {
<<<<<<< HEAD
		hostVars := v.host.regexp.FindStringSubmatch(getHost(req))
		if hostVars != nil {
			subexpNames := v.host.regexp.SubexpNames()
			varName := 0
			for i, name := range subexpNames[1:] {
				if name != "" && name == varGroupName(varName) {
					m.Vars[v.host.varsN[varName]] = hostVars[i+1]
					varName++
				}
			}
		}
	}
	// Store path variables.
	if v.path != nil {
		pathVars := v.path.regexp.FindStringSubmatch(req.URL.Path)
		if pathVars != nil {
			subexpNames := v.path.regexp.SubexpNames()
			varName := 0
			for i, name := range subexpNames[1:] {
				if name != "" && name == varGroupName(varName) {
					m.Vars[v.path.varsN[varName]] = pathVars[i+1]
					varName++
				}
			}
			// Check if we should redirect.
			if v.path.strictSlash {
				p1 := strings.HasSuffix(req.URL.Path, "/")
=======
		host := getHost(req)
		matches := v.host.regexp.FindStringSubmatchIndex(host)
		if len(matches) > 0 {
			extractVars(host, matches, v.host.varsN, m.Vars)
		}
	}
	path := req.URL.Path
	if r.useEncodedPath {
		path = getPath(req)
	}
	// Store path variables.
	if v.path != nil {
		matches := v.path.regexp.FindStringSubmatchIndex(path)
		if len(matches) > 0 {
			extractVars(path, matches, v.path.varsN, m.Vars)
			// Check if we should redirect.
			if v.path.strictSlash {
				p1 := strings.HasSuffix(path, "/")
>>>>>>> a24f8b0... add support for installing helm charts
				p2 := strings.HasSuffix(v.path.template, "/")
				if p1 != p2 {
					u, _ := url.Parse(req.URL.String())
					if p1 {
						u.Path = u.Path[:len(u.Path)-1]
					} else {
						u.Path += "/"
					}
					m.Handler = http.RedirectHandler(u.String(), 301)
				}
			}
		}
	}
	// Store query string variables.
	for _, q := range v.queries {
<<<<<<< HEAD
		queryVars := q.regexp.FindStringSubmatch(q.getUrlQuery(req))
		if queryVars != nil {
			subexpNames := q.regexp.SubexpNames()
			varName := 0
			for i, name := range subexpNames[1:] {
				if name != "" && name == varGroupName(varName) {
					m.Vars[q.varsN[varName]] = queryVars[i+1]
					varName++
				}
			}
=======
		queryURL := q.getURLQuery(req)
		matches := q.regexp.FindStringSubmatchIndex(queryURL)
		if len(matches) > 0 {
			extractVars(queryURL, matches, q.varsN, m.Vars)
>>>>>>> a24f8b0... add support for installing helm charts
		}
	}
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		return r.URL.Host
	}
	host := r.Host
	// Slice off any port information.
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host

}
<<<<<<< HEAD
=======

func extractVars(input string, matches []int, names []string, output map[string]string) {
	for i, name := range names {
		output[name] = input[matches[2*i+2]:matches[2*i+3]]
	}
}
>>>>>>> a24f8b0... add support for installing helm charts
