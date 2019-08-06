/*
	Copyright 2019 The pdfcpu Authors.

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

package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/denisbetsi/pdfcpu/pkg/log"
	pdf "github.com/denisbetsi/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
)

func imageObjNrs(ctx *pdf.Context, page int) []int {

	// TODO Exclude SMask image objects.

	o := []int{}

	for k, v := range ctx.Optimize.PageImages[page-1] {
		if v {
			o = append(o, k)
		}
	}

	return o
}

func imageFilenameWithoutExtension(dir, resID string, pageNr, objNr int) string {
	return filepath.Join(dir, fmt.Sprintf("%s_%d_%d", resID, pageNr, objNr))
}

func doExtractImages(ctx *pdf.Context, selectedPages pdf.IntSet) error {

	visited := pdf.IntSet{}

	for pageNr, v := range selectedPages {

		if v {

			log.Info.Printf("writing images for page %d\n", pageNr)

			for _, objNr := range imageObjNrs(ctx, pageNr) {

				if visited[objNr] {
					continue
				}

				visited[objNr] = true

				io, err := pdf.ExtractImageData(ctx, objNr)
				if err != nil {
					return err
				}

				if io == nil {
					continue
				}

				filename := imageFilenameWithoutExtension(ctx.Write.DirName, io.ResourceNames[0], pageNr, objNr)

				_, err = pdf.WriteImage(ctx.XRefTable, filename, io.ImageDict, objNr)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

// ExtractImages dumps embedded image resources from rs into outDir for selected pages.
func ExtractImages(rs io.ReadSeeker, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}

	fromStart := time.Now()
	ctx, durRead, durVal, durOpt, err := readValidateAndOptimize(rs, conf, fromStart)
	if err != nil {
		return err
	}

	fromWrite := time.Now()
	pages, err := pagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	ctx.Write.DirName = outDir
	if err = doExtractImages(ctx, pages); err != nil {
		return err
	}

	durWrite := time.Since(fromWrite).Seconds()
	durTotal := time.Since(fromStart).Seconds()
	log.Stats.Printf("XRefTable:\n%s\n", ctx)
	pdf.TimingStats("write images", durRead, durVal, durOpt, durWrite, durTotal)

	return nil
}

// ExtractImagesFile dumps embedded image resources from inFile into outDir for selected pages.
func ExtractImagesFile(inFile, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExtractImages(f, outDir, selectedPages, conf)
}

func fontObjNrs(ctx *pdf.Context, page int) []int {

	o := []int{}

	for k, v := range ctx.Optimize.PageFonts[page-1] {
		if v {
			o = append(o, k)
		}
	}

	return o
}

func doExtractFonts(ctx *pdf.Context, selectedPages pdf.IntSet) error {

	visited := pdf.IntSet{}

	for p, v := range selectedPages {

		if v {

			log.Info.Printf("writing fonts for page %d\n", p)

			for _, objNr := range fontObjNrs(ctx, p) {

				if visited[objNr] {
					continue
				}

				visited[objNr] = true

				fo, err := pdf.ExtractFontData(ctx, objNr)
				if err != nil {
					return err
				}

				if fo == nil {
					continue
				}

				fileName := fmt.Sprintf("%s/%s_%d_%d.%s", ctx.Write.DirName, fo.ResourceNames[0], p, objNr, fo.Extension)

				err = ioutil.WriteFile(fileName, fo.Data, os.ModePerm)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

// ExtractFonts dumps embedded fontfiles from rs into outDir for selected pages.
func ExtractFonts(rs io.ReadSeeker, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}

	fromStart := time.Now()
	ctx, durRead, durVal, durOpt, err := readValidateAndOptimize(rs, conf, fromStart)
	if err != nil {
		return err
	}

	fromWrite := time.Now()
	pages, err := pagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	ctx.Write.DirName = outDir
	if err = doExtractFonts(ctx, pages); err != nil {
		return err
	}

	durWrite := time.Since(fromWrite).Seconds()
	durTotal := time.Since(fromStart).Seconds()
	log.Stats.Printf("XRefTable:\n%s\n", ctx)
	pdf.TimingStats("write fonts", durRead, durVal, durOpt, durWrite, durTotal)

	return nil
}

// ExtractFontsFile dumps embedded fontfiles from inFile into outDir for selected pages.
func ExtractFontsFile(inFile, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExtractFonts(f, outDir, selectedPages, conf)
}

// singlePageFileName generates a filename for a Context and a specific page number.
func singlePageFileName(fileName string, pageNr int) string {
	baseFileName := filepath.Base(fileName)
	fn := strings.TrimSuffix(baseFileName, ".pdf")
	return fn + "_" + strconv.Itoa(pageNr) + ".pdf"
}

func writeSinglePagePDF(ctx *pdf.Context, pageNr int, outDir string) error {
	ctx.ResetWriteContext()
	w := ctx.Write
	w.SelectedPages[pageNr] = true
	w.DirName = outDir + "/"
	w.FileName = singlePageFileName("fn", pageNr)
	log.CLI.Printf("writing %s ...\n", w.DirName+w.FileName)
	return pdf.Write(ctx)
}

func writeSinglePagePDFs(ctx *pdf.Context, selectedPages pdf.IntSet, outDir string) error {
	for i, v := range selectedPages {
		if v {
			err := writeSinglePagePDF(ctx, i, outDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ExtractPages generates single page PDF files from rs in outDir for selected pages.
func ExtractPages(rs io.ReadSeeker, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
		conf.Cmd = pdf.EXTRACTPAGES
	}

	fromStart := time.Now()
	ctx, durRead, durVal, durOpt, err := readValidateAndOptimize(rs, conf, fromStart)
	if err != nil {
		return err
	}

	fromWrite := time.Now()
	pages, err := pagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	if err = writeSinglePagePDFs(ctx, pages, outDir); err != nil {
		return err
	}

	durWrite := time.Since(fromWrite).Seconds()
	durTotal := time.Since(fromStart).Seconds()
	log.Stats.Printf("XRefTable:\n%s\n", ctx)
	pdf.TimingStats("write PDFs", durRead, durVal, durOpt, durWrite, durTotal)

	return nil
}

// ExtractPagesFile generates single page PDF files from inFile in outDir for selected pages.
func ExtractPagesFile(inFile, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExtractPages(f, outDir, selectedPages, conf)
}

func contentObjNrs(ctx *pdf.Context, page int) ([]int, error) {

	objNrs := []int{}

	d, _, err := ctx.PageDict(page)
	if err != nil {
		return nil, err
	}

	o, found := d.Find("Contents")
	if !found || o == nil {
		return nil, nil
	}

	var objNr int

	ir, ok := o.(pdf.IndirectRef)
	if ok {
		objNr = ir.ObjectNumber.Value()
	}

	o, err = ctx.Dereference(o)
	if err != nil {
		return nil, err
	}

	if o == nil {
		return nil, nil
	}

	switch o := o.(type) {

	case pdf.StreamDict:

		objNrs = append(objNrs, objNr)

	case pdf.Array:

		for _, o := range o {

			ir, ok := o.(pdf.IndirectRef)
			if !ok {
				return nil, errors.Errorf("missing indref for page tree dict content no page %d", page)
			}

			sd, err := ctx.DereferenceStreamDict(ir)
			if err != nil {
				return nil, err
			}

			if sd == nil {
				continue
			}

			objNrs = append(objNrs, ir.ObjectNumber.Value())

		}

	}

	return objNrs, nil
}

func doExtractContent(ctx *pdf.Context, selectedPages pdf.IntSet) error {

	visited := pdf.IntSet{}

	for p, v := range selectedPages {

		if v {

			log.Info.Printf("writing content for page %d\n", p)

			objNrs, err := contentObjNrs(ctx, p)
			if err != nil {
				return err
			}

			if objNrs == nil {
				continue
			}

			for _, objNr := range objNrs {

				if visited[objNr] {
					continue
				}

				visited[objNr] = true

				b, err := pdf.ExtractStreamData(ctx, objNr)
				if err != nil {
					return err
				}

				if b == nil {
					continue
				}

				fileName := fmt.Sprintf("%s/%d_%d.txt", ctx.Write.DirName, p, objNr)

				err = ioutil.WriteFile(fileName, b, os.ModePerm)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

// ExtractContent dumps "PDF source" files from rs into outDir for selected pages.
func ExtractContent(rs io.ReadSeeker, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}

	fromStart := time.Now()
	ctx, durRead, durVal, durOpt, err := readValidateAndOptimize(rs, conf, fromStart)
	if err != nil {
		return err
	}

	fromWrite := time.Now()
	pages, err := pagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	ctx.Write.DirName = outDir
	if err = doExtractContent(ctx, pages); err != nil {
		return err
	}

	durWrite := time.Since(fromWrite).Seconds()
	durTotal := time.Since(fromStart).Seconds()
	log.Stats.Printf("XRefTable:\n%s\n", ctx)
	pdf.TimingStats("write content", durRead, durVal, durOpt, durWrite, durTotal)

	return nil
}

// ExtractContentFile dumps "PDF source" files from inFile into outDir for selected pages.
func ExtractContentFile(inFile, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExtractContent(f, outDir, selectedPages, conf)
}

func extractMetadataStream(ctx *pdf.Context, obj pdf.Object, objNr int, dt string) error {

	ir, _ := obj.(pdf.IndirectRef)
	sObjNr := ir.ObjectNumber.Value()
	b, err := pdf.ExtractStreamData(ctx, sObjNr)
	if err != nil {
		return err
	}

	if b == nil {
		return nil
	}

	fileName := fmt.Sprintf("%s/%d_%s.txt", ctx.Write.DirName, objNr, dt)

	return ioutil.WriteFile(fileName, b, os.ModePerm)
}

func doExtractMetadata(ctx *pdf.Context, selectedPages pdf.IntSet) error {

	for k, v := range ctx.XRefTable.Table {
		if v.Free || v.Compressed {
			continue
		}
		switch d := v.Object.(type) {

		case pdf.Dict:

			o, found := d.Find("Metadata")
			if !found || o == nil {
				continue
			}

			dt := "unknown"
			if d.Type() != nil {
				dt = *d.Type()
			}

			err := extractMetadataStream(ctx, o, k, dt)
			if err != nil {
				return err
			}

		case pdf.StreamDict:

			o, found := d.Find("Metadata")
			if !found || o == nil {
				continue
			}

			dt := "unknown"
			if d.Type() != nil {
				dt = *d.Type()
			}

			err := extractMetadataStream(ctx, o, k, dt)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// ExtractMetadata dumps all metadata dict entries for rs into outDir.
func ExtractMetadata(rs io.ReadSeeker, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}

	fromStart := time.Now()
	ctx, durRead, durVal, durOpt, err := readValidateAndOptimize(rs, conf, fromStart)
	if err != nil {
		return err
	}

	fromWrite := time.Now()
	pages, err := pagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	ctx.Write.DirName = outDir
	if err = doExtractMetadata(ctx, pages); err != nil {
		return err
	}

	durWrite := time.Since(fromWrite).Seconds()
	durTotal := time.Since(fromStart).Seconds()
	log.Stats.Printf("XRefTable:\n%s\n", ctx)
	pdf.TimingStats("write metadata", durRead, durVal, durOpt, durWrite, durTotal)

	return nil
}

// ExtractMetadataFile dumps all metadata dict entries for inFile into outDir.
func ExtractMetadataFile(inFile, outDir string, selectedPages []string, conf *pdf.Configuration) error {
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ExtractMetadata(f, outDir, selectedPages, conf)
}
