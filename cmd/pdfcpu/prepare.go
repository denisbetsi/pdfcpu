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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/denisbetsi/pdfcpu/pkg/api"
	"github.com/denisbetsi/pdfcpu/pkg/cli"
	"github.com/denisbetsi/pdfcpu/pkg/pdfcpu"
)

func hasPdfExtension(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

func ensurePdfExtension(filename string) {
	if !hasPdfExtension(filename) {
		fmt.Fprintf(os.Stderr, "%s needs extension \".pdf\".\n", filename)
		os.Exit(1)
	}
}

func defaultFilenameOut(filename string) string {
	return filename[:len(filename)-4] + "_new.pdf"
}

func printHelp(conf *pdfcpu.Configuration) {
	switch len(flag.Args()) {

	case 0:
		fmt.Fprintln(os.Stderr, usage)

	case 1:
		fmt.Fprintln(os.Stderr, cmdMap.HelpString(flag.Arg(0)))

	default:
		fmt.Fprintln(os.Stderr, "usage: pdfcpu help command\n\nToo many arguments.")

	}
}

func printPaperSizes(conf *pdfcpu.Configuration) {
	fmt.Fprintln(os.Stderr, paperSizes)
}

func printVersion(conf *pdfcpu.Configuration) {
	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageVersion)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "pdfcpu: %v\n build: %v\ncommit: %v\n", pdfcpu.VersionStr, date, commit)
}

func handleValidateCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 1 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageValidate)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	if mode != "" && mode != "strict" && mode != "s" && mode != "relaxed" && mode != "r" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageValidate)
		os.Exit(1)
	}

	switch mode {
	case "strict", "s":
		conf.ValidationMode = pdfcpu.ValidationStrict
	case "relaxed", "r":
		conf.ValidationMode = pdfcpu.ValidationRelaxed
	}

	process(cli.ValidateCommand(inFile, conf))
}

func handleOptimizeCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageOptimize)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := inFile
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	conf.StatsFileName = fileStats
	if len(fileStats) > 0 {
		fmt.Fprintf(os.Stdout, "stats will be appended to %s\n", fileStats)
	}

	process(cli.OptimizeCommand(inFile, outFile, conf))
}

func handleSplitCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 2 || len(flag.Args()) > 3 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageSplit)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outDir := flag.Arg(1)

	span := 1
	var err error
	if len(flag.Args()) == 3 {
		span, err = strconv.Atoi(flag.Arg(2))
		if err != nil || span < 1 {
			fmt.Fprintln(os.Stderr, "split: span is a numeric value >= 1")
			os.Exit(1)
		}
	}

	process(cli.SplitCommand(inFile, outDir, span, conf))
}

func handleMergeCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 3 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageMerge)
		os.Exit(1)
	}

	filesIn := []string{}
	outFile := ""
	for i, arg := range flag.Args() {
		ensurePdfExtension(arg)
		if i == 0 {
			outFile = arg
			continue
		}
		filesIn = append(filesIn, arg)
	}

	process(cli.MergeCommand(filesIn, outFile, conf))
}

func extractModeCompletion(modePrefix string) string {
	var modeStr string
	for _, mode := range []string{"image", "font", "page", "content", "meta"} {
		if !strings.HasPrefix(mode, modePrefix) {
			continue
		}
		if len(modeStr) > 0 {
			return ""
		}
		modeStr = mode
	}
	return modeStr
}

func handleExtractCommand(conf *pdfcpu.Configuration) {
	mode = extractModeCompletion(mode)
	if len(flag.Args()) != 2 || mode == "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageExtract)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)
	outDir := flag.Arg(1)

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	var cmd *cli.Command

	switch mode {

	case "image":
		cmd = cli.ExtractImagesCommand(inFile, outDir, pages, conf)

	case "font":
		cmd = cli.ExtractFontsCommand(inFile, outDir, pages, conf)

	case "page":
		cmd = cli.ExtractPagesCommand(inFile, outDir, pages, conf)

	case "content":
		cmd = cli.ExtractContentCommand(inFile, outDir, pages, conf)

	case "meta":
		cmd = cli.ExtractMetadataCommand(inFile, outDir, conf)

	default:
		fmt.Fprintf(os.Stderr, "unknown extract mode: %s\n", mode)
		os.Exit(1)

	}

	process(cmd)
}

func handleTrimCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 || selectedPages == "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageTrim)
		os.Exit(1)
	}

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := ""
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	process(cli.TrimCommand(inFile, outFile, pages, conf))
}

func handleListAttachmentsCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) != 1 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n", usageAttachList)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)
	process(cli.ListAttachmentsCommand(inFile, conf))
}

func handleAddAttachmentsCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 2 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", usageAttachAdd)
		os.Exit(1)
	}

	var inFile string
	fileNames := []string{}

	for i, arg := range flag.Args() {
		if i == 0 {
			inFile = arg
			ensurePdfExtension(inFile)
			continue
		}
		fileNames = append(fileNames, arg)
	}

	process(cli.AddAttachmentsCommand(inFile, "", fileNames, conf))
}

func handleRemoveAttachmentsCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 1 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", usageAttachRemove)
		os.Exit(1)
	}

	var inFile string
	fileNames := []string{}

	for i, arg := range flag.Args() {
		if i == 0 {
			inFile = arg
			ensurePdfExtension(inFile)
			continue
		}
		fileNames = append(fileNames, arg)
	}

	process(cli.RemoveAttachmentsCommand(inFile, "", fileNames, conf))
}

func handleExtractAttachmentsCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 2 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", usageAttachExtract)
		os.Exit(1)
	}

	var inFile string
	fileNames := []string{}
	var outDir string

	for i, arg := range flag.Args() {
		if i == 0 {
			inFile = arg
			ensurePdfExtension(inFile)
			continue
		}
		if i == 1 {
			outDir = arg
			continue
		}
		fileNames = append(fileNames, arg)
	}

	process(cli.ExtractAttachmentsCommand(inFile, outDir, fileNames, conf))
}

func handleListPermissionsCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) != 1 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n", usagePermList)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	process(cli.ListPermissionsCommand(inFile, conf))
}

func permCompletion(permPrefix string) string {
	var permStr string

	for _, perm := range []string{"none", "all"} {
		if !strings.HasPrefix(perm, permPrefix) {
			continue
		}
		if len(permStr) > 0 {
			return ""
		}
		permStr = perm
	}

	return permStr
}

func handleSetPermissionsCommand(conf *pdfcpu.Configuration) {
	if perm != "" {
		perm = permCompletion(perm)
	}
	if len(flag.Args()) != 1 || selectedPages != "" ||
		!(perm == "none" || perm == "all") {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", usagePermSet)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	if perm == "all" {
		conf.Permissions = pdfcpu.PermissionsAll
	}

	process(cli.SetPermissionsCommand(inFile, "", conf))
}

func handleDecryptCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageDecrypt)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := inFile
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	process(cli.DecryptCommand(inFile, outFile, conf))
}

func validateEncryptModeFlag() {
	if !pdfcpu.MemberOf(mode, []string{"rc4", "aes", ""}) {
		fmt.Fprintf(os.Stderr, "%s\n\n", "valid modes: rc4,aes default:aes")
		os.Exit(1)
	}

	// Default to AES encryption.
	if mode == "" {
		mode = "aes"
	}

	if key == "256" && mode == "rc4" {
		key = "128"
	}

	if mode == "rc4" {
		if key != "40" && key != "128" && key != "" {
			fmt.Fprintf(os.Stderr, "%s\n\n", "supported RC4 key lengths: 40,128 default:128")
			os.Exit(1)
		}
	}

	if mode == "aes" {
		if key != "40" && key != "128" && key != "256" && key != "" {
			fmt.Fprintf(os.Stderr, "%s\n\n", "supported AES key lengths: 40,128,256 default:256")
			os.Exit(1)
		}
	}

}

func validateEncryptFlags() {
	validateEncryptModeFlag()
	if perm != "none" && perm != "all" && perm != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", "supported permissions: none,all default:none (viewing always allowed!)")
		os.Exit(1)
	}
}

func handleEncryptCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageEncrypt)
		os.Exit(1)
	}

	if conf.OwnerPW == "" {
		fmt.Fprintln(os.Stderr, "missing non-empty owner password!")
		fmt.Fprintf(os.Stderr, "%s\n\n", usageEncrypt)
		os.Exit(1)
	}

	validateEncryptFlags()

	conf.EncryptUsingAES = mode != "rc4"

	kl, _ := strconv.Atoi(key)
	conf.EncryptKeyLength = kl

	if perm == "all" {
		conf.Permissions = pdfcpu.PermissionsAll
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := inFile
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	process(cli.EncryptCommand(inFile, outFile, conf))
}

func handleChangeUserPasswordCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) != 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageChangeUserPW)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := inFile
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)

	}

	pwOld := flag.Arg(1)
	pwNew := flag.Arg(2)

	process(cli.ChangeUserPWCommand(inFile, outFile, &pwOld, &pwNew, conf))
}

func handleChangeOwnerPasswordCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) != 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageChangeOwnerPW)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	outFile := inFile
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	pwOld := flag.Arg(1)
	pwNew := flag.Arg(2)
	if pwNew == "" {
		fmt.Fprintf(os.Stderr, "owner password cannot be empty")
		fmt.Fprintf(os.Stderr, "%s\n\n", usageChangeOwnerPW)
		os.Exit(1)
	}

	process(cli.ChangeOwnerPWCommand(inFile, outFile, &pwOld, &pwNew, conf))
}

func handleWatermarksCommand(conf *pdfcpu.Configuration, onTop bool) {
	if len(flag.Args()) < 2 || len(flag.Args()) > 3 {
		s := usageWatermark
		if onTop {
			s = usageStamp
		}
		fmt.Fprintf(os.Stderr, "%s\n\n", s)
		os.Exit(1)
	}

	selectedPages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v", err)
		os.Exit(1)
	}

	wm, err := pdfcpu.ParseWatermarkDetails(flag.Arg(0), onTop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	inFile := flag.Arg(1)
	ensurePdfExtension(inFile)

	outFile := ""
	if len(flag.Args()) == 3 {
		outFile = flag.Arg(2)
		ensurePdfExtension(outFile)
	}

	process(cli.AddWatermarksCommand(inFile, outFile, selectedPages, wm, conf))
}

func handleAddStampsCommand(conf *pdfcpu.Configuration) {
	handleWatermarksCommand(conf, true)
}

func handleAddWatermarksCommand(conf *pdfcpu.Configuration) {
	handleWatermarksCommand(conf, false)
}

func hasImageExtension(filename string) bool {
	s := strings.ToLower(filepath.Ext(filename))
	return pdfcpu.MemberOf(s, []string{".jpg", ".jpeg", ".png", ".tif", ".tiff"})
}

func ensureImageExtension(filename string) {
	if !hasImageExtension(filename) {
		fmt.Fprintf(os.Stderr, "%s needs an image extension (.jpg, .jpeg, .png, .tif, .tiff)\n", filename)
		os.Exit(1)
	}
}

func handleImportImagesCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 2 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageImportImages)
		os.Exit(1)
	}

	var outFile string
	outFile = flag.Arg(0)
	if hasPdfExtension(outFile) {
		// pdfcpu import outFile imageFile...
		imp := pdfcpu.DefaultImportConfig()
		imageFileNames := []string{}
		for i := 1; i < len(flag.Args()); i++ {
			arg := flag.Arg(i)
			ensureImageExtension(arg)
			imageFileNames = append(imageFileNames, arg)
		}
		process(cli.ImportImagesCommand(imageFileNames, outFile, imp, conf))
	}

	// pdfcpu import description outFile imageFile...
	imp, err := pdfcpu.ParseImportDetails(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if imp == nil {
		fmt.Fprintf(os.Stderr, "missing import description\n")
		os.Exit(1)
	}

	outFile = flag.Arg(1)
	ensurePdfExtension(outFile)

	imageFileNames := []string{}
	for i := 2; i < len(flag.Args()); i++ {
		arg := flag.Args()[i]
		ensureImageExtension(arg)
		imageFileNames = append(imageFileNames, arg)
	}

	process(cli.ImportImagesCommand(imageFileNames, outFile, imp, conf))
}

func handleInsertPagesCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usagePagesInsert)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)
	outFile := ""
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	process(cli.InsertPagesCommand(inFile, outFile, pages, conf))
}

func handleRemovePagesCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 || selectedPages == "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usagePagesRemove)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)
	outFile := ""
	if len(flag.Args()) == 2 {
		outFile = flag.Arg(1)
		ensurePdfExtension(outFile)
	}

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}
	if pages == nil {
		fmt.Fprintf(os.Stderr, "missing page selection\n")
		os.Exit(1)
	}

	process(cli.RemovePagesCommand(inFile, outFile, pages, conf))
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func handleRotateCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageRotate)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	rotation, err := strconv.Atoi(flag.Arg(1))
	if err != nil || abs(rotation)%90 > 0 {
		fmt.Fprintf(os.Stderr, "rotation must be a multiple of 90: %s\n", flag.Arg(1))
		os.Exit(1)
	}

	selectedPages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	process(cli.RotateCommand(inFile, "", rotation, selectedPages, conf))
}

func parseAfterNUpDetails(nup *pdfcpu.NUp, argInd int, filenameOut string) []string {

	var err error

	if nup.PageGrid {
		cols, err := strconv.Atoi(flag.Arg(argInd))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		rows, err := strconv.Atoi(flag.Arg(argInd + 1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		err = pdfcpu.ParseNUpGridDefinition(cols, rows, nup)
		argInd += 2
	} else {
		n, err := strconv.Atoi(flag.Arg(argInd))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		err = pdfcpu.ParseNUpValue(n, nup)
		argInd++
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	filenameIn := flag.Arg(argInd)
	if !hasPdfExtension(filenameIn) && !hasImageExtension(filenameIn) {
		fmt.Fprintf(os.Stderr, "inFile has to be a PDF or one or a sequence of image files: %s\n", filenameIn)
		os.Exit(1)
	}

	filenamesIn := []string{filenameIn}

	if hasPdfExtension(filenameIn) {
		if len(flag.Args()) > argInd+1 {
			usage := usageNUp
			if nup.PageGrid {
				usage = usageGrid
			}
			fmt.Fprintf(os.Stderr, "%s\n\n", usage)
			os.Exit(1)
		}
		if filenameIn == filenameOut {
			fmt.Fprintln(os.Stderr, "inFile and outFile can't be the same.")
			os.Exit(1)
		}
	} else {
		nup.ImgInputFile = true
		for i := argInd + 1; i < len(flag.Args()); i++ {
			arg := flag.Args()[i]
			ensureImageExtension(arg)
			filenamesIn = append(filenamesIn, arg)
		}
	}

	return filenamesIn
}

func handleNUpCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageNUp)
		os.Exit(1)
	}

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	nup := pdfcpu.DefaultNUpConfig()

	outFile := flag.Arg(0)
	if hasPdfExtension(outFile) {
		// pdfcpu nup outFile n inFile|imageFiles...
		// No optional 'description' argument provided.
		// We use the default nup configuration.
		inFiles := parseAfterNUpDetails(nup, 1, outFile)
		process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
	}

	// pdfcpu nup description outFile n inFile|imageFiles...
	if err = pdfcpu.ParseNUpDetails(flag.Arg(0), nup); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	outFile = flag.Arg(1)
	ensurePdfExtension(outFile)
	inFiles := parseAfterNUpDetails(nup, 2, outFile)
	process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
}

func handleGridCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) < 4 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageGrid)
		os.Exit(1)
	}

	pages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem with flag selectedPages: %v\n", err)
		os.Exit(1)
	}

	nup := pdfcpu.DefaultNUpConfig()

	nup.PageGrid = true
	outFile := flag.Arg(0)
	if hasPdfExtension(outFile) {
		// pdfcpu grid outFile m n inFile|imageFiles...
		// No optional 'description' argument provided.
		// We use the default nup configuration.
		inFiles := parseAfterNUpDetails(nup, 1, outFile)
		process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
	}

	// pdfcpu grid description outFile m n inFile|imageFiles...
	err = pdfcpu.ParseNUpDetails(flag.Arg(0), nup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	outFile = flag.Arg(1)
	ensurePdfExtension(outFile)
	inFiles := parseAfterNUpDetails(nup, 2, outFile)
	process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
}

func handleInfoCommand(conf *pdfcpu.Configuration) {
	if len(flag.Args()) == 0 || len(flag.Args()) > 1 || selectedPages != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", usageInfo)
		os.Exit(1)
	}

	inFile := flag.Arg(0)
	ensurePdfExtension(inFile)

	process(cli.InfoCommand(inFile, conf))
}
