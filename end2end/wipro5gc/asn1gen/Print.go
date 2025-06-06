// Package asn1gen - ASN1C generated code package
/**
 * This file was generated by the Objective Systems ASN1C Compiler
 * (https://obj-sys.com).  Version: 7.7.2, Date: 05-Feb-2024.
 *
 * Copyright (c) 2020-2023 Objective Systems, Inc.
 *
 * This software is furnished under a license and may be used and copied
 * only in accordance with the terms of such license and with the
 * inclusion of the above copyright notice. This software or any other
 * copies thereof may not be provided or otherwise made available to any
 * other person. No title to and ownership of the software is hereby
 * transferred.
 *
 * The information in this software is subject to change without notice
 * and should not be construed as a commitment by Objective Systems, Inc.
 *
 * PROPRIETARY NOTICE
 *
 * This software is an unpublished work subject to a confidentiality agreement
 * and is protected by copyright and trade secret law.  Unauthorized copying,
 * redistribution or other use of this work is prohibited.
 *
 * The above notice of copyright on this source code product does not
 * indicate any actual or intended publication of such source code.
 *
 * Command:  asn1c /home/imgadmin/asn1c-v772/golang/sample_per/ngap/ngap.asn -i /home/imgadmin/asn1c-v772/golang/sample_per/ngap -o src -oh src -genprint -genprttostr -gentest -aper -go -genmake src/makefile -prjdir ../ngap
 *
 **************************************************************************/
package asn1gen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func PrintToString(val interface{}) string {
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	return StripSpew(spew.Sdump(val))
}

func Print(val interface{}) {
	fmt.Print(PrintToString(val))
}

// StripSpew accepts a string of output from spew and returns a string that
// has <nil> pointers and most type information stripped out.
func StripSpew(spewOutput string) string {
	var sb strings.Builder
	var matchFound bool = false
	var newpiece string

	// Compile regular expression to detect nil pointers.
	renil, _ := regexp.Compile(`<nil>`)

	// Compile a regular expression to detect type information.
	retype, _ := regexp.Compile(`(^.*?\:\s+)\(.*?\)`)

	pieces := strings.Split(spewOutput, "\n")
	for _, piece := range pieces {
		if len(piece) == 0 {
			sb.WriteByte('\n')
			continue
		}

		// See if a nil pointer is indicated.  If so, we remove this item from
		// the output.
		if renil.MatchString(piece) {
			matchFound = true
			continue
		} else {
			newpiece = piece
		}

		// See if type information is indicated.  If so, we remove it.
		if retype.MatchString(piece) {
			matchFound = true
			newpiece = retype.ReplaceAllString(piece, "$1")
		}
		sb.WriteString(newpiece)
		sb.WriteByte('\n')
	}
	if matchFound {
		return sb.String()
	}

	return spewOutput
}
