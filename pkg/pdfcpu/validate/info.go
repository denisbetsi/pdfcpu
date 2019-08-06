/*
Copyright 2018 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validate

import (
	"github.com/denisbetsi/pdfcpu/pkg/log"
	pdf "github.com/denisbetsi/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
)

func handleDefault(xRefTable *pdf.XRefTable, o pdf.Object) (err error) {

	if xRefTable.ValidationMode == pdf.ValidationStrict {
		_, err = xRefTable.DereferenceStringOrHexLiteral(o, pdf.V10, nil)
	} else {
		_, err = xRefTable.Dereference(o)
	}

	return err
}

func validateInfoDictDate(xRefTable *pdf.XRefTable, o pdf.Object) (s string, err error) {
	if xRefTable.ValidationMode == pdf.ValidationRelaxed {
		return validateString(xRefTable, o, nil)
	}
	return validateDateObject(xRefTable, o, pdf.V10)
}

func validateInfoDictTrapped(o pdf.Object, xRefTable *pdf.XRefTable) error {

	sinceVersion := pdf.V13

	validate := func(s string) bool { return pdf.MemberOf(s, []string{"True", "False", "Unknown"}) }

	if xRefTable.ValidationMode == pdf.ValidationRelaxed {
		validate = func(s string) bool {
			return pdf.MemberOf(s, []string{"True", "False", "Unknown", "true", "false", "unknown"})
		}
	}

	_, err := xRefTable.DereferenceName(o, sinceVersion, validate)
	if err == nil {
		return nil
	}

	if xRefTable.ValidationMode == pdf.ValidationRelaxed {
		_, err = xRefTable.DereferenceBoolean(o, sinceVersion)
	}

	return err
}

func validateDocumentInfoDict(xRefTable *pdf.XRefTable, obj pdf.Object) (hasModDate bool, err error) {

	// Document info object is optional.

	d, err := xRefTable.DereferenceDict(obj)
	if err != nil || d == nil {
		return false, err
	}

	for k, v := range d {

		switch k {

		// text string, opt, since V1.1
		case "Title":
			xRefTable.Title, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V11, nil)

		// text string, optional
		case "Author":
			xRefTable.Author, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V10, nil)

		// text string, optional, since V1.1
		case "Subject":
			xRefTable.Subject, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V11, nil)

		// text string, optional, since V1.1
		case "Keywords":
			xRefTable.Keywords, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V11, nil)

		// text string, optional
		case "Creator":
			xRefTable.Creator, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V10, nil)

		// text string, optional
		case "Producer":
			xRefTable.Producer, err = xRefTable.DereferenceStringOrHexLiteral(v, pdf.V10, nil)

		// date, optional
		case "CreationDate":
			xRefTable.CreationDate, err = validateInfoDictDate(xRefTable, v)

		// date, required if PieceInfo is present in document catalog.
		case "ModDate":
			hasModDate = true
			xRefTable.ModDate, err = validateInfoDictDate(xRefTable, v)

		// name, optional, since V1.3
		case "Trapped":
			err = validateInfoDictTrapped(v, xRefTable)

		// text string, optional
		default:
			err = handleDefault(xRefTable, v)

		}

		if err != nil {
			return false, err
		}

	}

	return hasModDate, nil
}

func validateDocumentInfoObject(xRefTable *pdf.XRefTable) error {

	// Document info object is optional.
	if xRefTable.Info == nil {
		return nil
	}

	log.Validate.Println("*** validateDocumentInfoObject begin ***")

	hasModDate, err := validateDocumentInfoDict(xRefTable, *xRefTable.Info)
	if err != nil {
		return err
	}

	hasPieceInfo, err := xRefTable.CatalogHasPieceInfo()
	if err != nil {
		return err
	}

	if hasPieceInfo && !hasModDate {
		return errors.Errorf("validateDocumentInfoObject: missing required entry \"ModDate\"")
	}

	log.Validate.Println("*** validateDocumentInfoObject end ***")

	return nil
}
