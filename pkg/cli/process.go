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

package cli

import (
	"io"

	pdf "github.com/denisbetsi/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
)

// Command represents an execution context.
type Command struct {
	Mode          pdf.CommandMode    // VALIDATE  OPTIMIZE  SPLIT  MERGE  EXTRACT  TRIM  LISTATT ADDATT REMATT EXTATT  ENCRYPT  DECRYPT  CHANGEUPW  CHANGEOPW LISTP ADDP  WATERMARK  IMPORT  INSERTP REMOVEP ROTATE  NUP
	InFile        *string            //    *         *        *      -       *      *      *       *       *      *       *        *         *          *       *     *       *         -       *       *       *     -
	InFiles       []string           //    -         -        -      *       -      -      -       *       *      *       -        -         -          -       -     -       -         *       -       -       -     *
	InDir         *string            //    -         -        -      -       -      -      -       -       -      -       -        -         -          -       -     -       -         -       -       -       -     -
	OutFile       *string            //    -         *        -      *       -      *      -       -       -      -       *        *         *          *       -     -       *         *       *       *       -     *
	OutDir        *string            //    -         -        *      -       *      -      -       -       -      *       -        -         -          -       -     -       -         -       -       -       -     -
	PageSelection []string           //    -         -        -      -       *      *      -       -       -      -       -        -         -          -       -     -       *         -       *       *       -     *
	Conf          *pdf.Configuration //    *         *        *      *       *      *      *       *       *      *       *        *         *          *       *     *       *         *       *       *       *     *
	PWOld         *string            //    -         -        -      -       -      -      -       -       -      -       -        -         *          *       -     -       -         -       -       -       -     -
	PWNew         *string            //    -         -        -      -       -      -      -       -       -      -       -        -         *          *       -     -       -         -       -       -       -     -
	Watermark     *pdf.Watermark     //    -         -        -      -       -      -      -       -       -      -       -        -         -          -       -     -       -         -       -       -       -     -
	Span          int                //    -         -        *      -       -      -      -       -       -      -       -        -         -          -       -     -       -         -       -       -       -     -
	Import        *pdf.Import        //    -         -        -      -       -      -      -       -       -      -       -        -         -          -       -     -       -         *       -       -       -     -
	Rotation      int                //    -         -        -      -       -      -      -       -       -      -       -        -         -          -       -     -       -         -       -       -       *     -
	NUp           *pdf.NUp           //    -         -        -      -       -      -      -       -       -      -       -        -         -          -       -     -       -         -       -       -       -     *
	Input         io.ReadSeeker
	Inputs        []io.ReadSeeker
	Output        io.Writer
}

var cmdMap = map[pdf.CommandMode]func(cmd *Command) ([]string, error){
	pdf.VALIDATE:           Validate,
	pdf.OPTIMIZE:           Optimize,
	pdf.SPLIT:              Split,
	pdf.MERGE:              Merge,
	pdf.EXTRACTIMAGES:      ExtractImages,
	pdf.EXTRACTFONTS:       ExtractFonts,
	pdf.EXTRACTPAGES:       ExtractPages,
	pdf.EXTRACTCONTENT:     ExtractContent,
	pdf.EXTRACTMETADATA:    ExtractMetadata,
	pdf.TRIM:               Trim,
	pdf.ADDWATERMARKS:      AddWatermarks,
	pdf.LISTATTACHMENTS:    processAttachments,
	pdf.ADDATTACHMENTS:     processAttachments,
	pdf.REMOVEATTACHMENTS:  processAttachments,
	pdf.EXTRACTATTACHMENTS: processAttachments,
	pdf.ENCRYPT:            processEncryption,
	pdf.DECRYPT:            processEncryption,
	pdf.CHANGEUPW:          processEncryption,
	pdf.CHANGEOPW:          processEncryption,
	pdf.LISTPERMISSIONS:    processPermissions,
	pdf.SETPERMISSIONS:     processPermissions,
	pdf.IMPORTIMAGES:       ImportImages,
	pdf.INSERTPAGES:        InsertPages,
	pdf.REMOVEPAGES:        RemovePages,
	pdf.ROTATE:             Rotate,
	pdf.NUP:                NUp,
	pdf.INFO:               Info,
}

// Process executes a pdfcpu command.
func Process(cmd *Command) (out []string, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("unexpected panic attack: %v\n", r)
		}
	}()

	cmd.Conf.Cmd = cmd.Mode

	if f, ok := cmdMap[cmd.Mode]; ok {
		return f(cmd)
	}

	return nil, errors.Errorf("Process: Unknown command mode %d\n", cmd.Mode)
}

// ValidateCommand creates a new command to validate a file.
func ValidateCommand(inFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.VALIDATE
	return &Command{
		Mode:   pdf.VALIDATE,
		InFile: &inFile,
		Conf:   conf}
}

// OptimizeCommand creates a new command to optimize a file.
func OptimizeCommand(inFile, outFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.OPTIMIZE
	return &Command{
		Mode:    pdf.OPTIMIZE,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

// SplitCommand creates a new command to split a file into single page files.
func SplitCommand(inFile, dirNameOut string, span int, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.SPLIT
	return &Command{
		Mode:   pdf.SPLIT,
		InFile: &inFile,
		OutDir: &dirNameOut,
		Span:   span,
		Conf:   conf}
}

// MergeCommand creates a new command to merge files.
func MergeCommand(inFiles []string, outFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.MERGE
	return &Command{
		Mode:    pdf.MERGE,
		InFiles: inFiles,
		OutFile: &outFile,
		Conf:    conf}
}

// ExtractImagesCommand creates a new command to extract embedded images.
// (experimental
func ExtractImagesCommand(inFile string, outDir string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTIMAGES
	return &Command{
		Mode:          pdf.EXTRACTIMAGES,
		InFile:        &inFile,
		OutDir:        &outDir,
		PageSelection: pageSelection,
		Conf:          conf}
}

// ExtractFontsCommand creates a new command to extract embedded fonts.
// (experimental)
func ExtractFontsCommand(inFile string, outDir string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTFONTS
	return &Command{
		Mode:          pdf.EXTRACTFONTS,
		InFile:        &inFile,
		OutDir:        &outDir,
		PageSelection: pageSelection,
		Conf:          conf}
}

// ExtractPagesCommand creates a new command to extract specific pages of a file.
func ExtractPagesCommand(inFile string, outDir string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTPAGES
	return &Command{
		Mode:          pdf.EXTRACTPAGES,
		InFile:        &inFile,
		OutDir:        &outDir,
		PageSelection: pageSelection,
		Conf:          conf}
}

// ExtractContentCommand creates a new command to extract page content streams.
func ExtractContentCommand(inFile string, outDir string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTCONTENT
	return &Command{
		Mode:          pdf.EXTRACTCONTENT,
		InFile:        &inFile,
		OutDir:        &outDir,
		PageSelection: pageSelection,
		Conf:          conf}
}

// ExtractMetadataCommand creates a new command to extract metadata streams.
func ExtractMetadataCommand(inFile string, outDir string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTMETADATA
	return &Command{
		Mode:   pdf.EXTRACTMETADATA,
		InFile: &inFile,
		OutDir: &outDir,
		Conf:   conf}
}

// TrimCommand creates a new command to trim the pages of a file.
func TrimCommand(inFile, outFile string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.TRIM
	return &Command{
		Mode:          pdf.TRIM,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// ListAttachmentsCommand create a new command to list attachments.
func ListAttachmentsCommand(inFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.LISTATTACHMENTS
	return &Command{
		Mode:   pdf.LISTATTACHMENTS,
		InFile: &inFile,
		Conf:   conf}
}

// AddAttachmentsCommand creates a new command to add attachments.
func AddAttachmentsCommand(inFile, outFile string, fileNames []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.ADDATTACHMENTS
	return &Command{
		Mode:    pdf.ADDATTACHMENTS,
		InFile:  &inFile,
		OutFile: &outFile,
		InFiles: fileNames,
		Conf:    conf}
}

// RemoveAttachmentsCommand creates a new command to remove attachments.
func RemoveAttachmentsCommand(inFile, outFile string, fileNames []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.REMOVEATTACHMENTS
	return &Command{
		Mode:    pdf.REMOVEATTACHMENTS,
		InFile:  &inFile,
		OutFile: &outFile,
		InFiles: fileNames,
		Conf:    conf}
}

// ExtractAttachmentsCommand creates a new command to extract attachments.
func ExtractAttachmentsCommand(inFile string, outDir string, fileNames []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.EXTRACTATTACHMENTS
	return &Command{
		Mode:    pdf.EXTRACTATTACHMENTS,
		InFile:  &inFile,
		OutDir:  &outDir,
		InFiles: fileNames,
		Conf:    conf}
}

// EncryptCommand creates a new command to encrypt a file.
func EncryptCommand(inFile, outFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.ENCRYPT
	return &Command{
		Mode:    pdf.ENCRYPT,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

// DecryptCommand creates a new command to decrypt a file.
func DecryptCommand(inFile, outFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.DECRYPT
	return &Command{
		Mode:    pdf.DECRYPT,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

// ChangeUserPWCommand creates a new command to change the user password.
func ChangeUserPWCommand(inFile, outFile string, pwOld, pwNew *string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.CHANGEUPW
	return &Command{
		Mode:    pdf.CHANGEUPW,
		InFile:  &inFile,
		OutFile: &outFile,
		PWOld:   pwOld,
		PWNew:   pwNew,
		Conf:    conf}
}

// ChangeOwnerPWCommand creates a new command to change the owner password.
func ChangeOwnerPWCommand(inFile, outFile string, pwOld, pwNew *string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.CHANGEOPW
	return &Command{
		Mode:    pdf.CHANGEOPW,
		InFile:  &inFile,
		OutFile: &outFile,
		PWOld:   pwOld,
		PWNew:   pwNew,
		Conf:    conf}
}

// ListPermissionsCommand create a new command to list permissions.
func ListPermissionsCommand(inFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.LISTPERMISSIONS
	return &Command{
		Mode:   pdf.LISTPERMISSIONS,
		InFile: &inFile,
		Conf:   conf}
}

// SetPermissionsCommand creates a new command to add permissions.
func SetPermissionsCommand(inFile, outFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.SETPERMISSIONS
	return &Command{
		Mode:    pdf.SETPERMISSIONS,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

func processAttachments(cmd *Command) (out []string, err error) {
	switch cmd.Mode {

	case pdf.LISTATTACHMENTS:
		out, err = ListAttachments(cmd)

	case pdf.ADDATTACHMENTS:
		out, err = AddAttachments(cmd)

	case pdf.REMOVEATTACHMENTS:
		out, err = RemoveAttachments(cmd)

	case pdf.EXTRACTATTACHMENTS:
		out, err = ExtractAttachments(cmd)
	}

	return out, err
}

func processEncryption(cmd *Command) (out []string, err error) {
	switch cmd.Mode {

	case pdf.ENCRYPT:
		return Encrypt(cmd)

	case pdf.DECRYPT:
		return Decrypt(cmd)

	case pdf.CHANGEUPW:
		return ChangeUserPassword(cmd)

	case pdf.CHANGEOPW:
		return ChangeOwnerPassword(cmd)
	}

	return nil, nil
}

func processPermissions(cmd *Command) (out []string, err error) {
	switch cmd.Mode {

	case pdf.LISTPERMISSIONS:
		return ListPermissions(cmd)

	case pdf.SETPERMISSIONS:
		return SetPermissions(cmd)
	}

	return nil, nil
}

// AddWatermarksCommand creates a new command to add Watermarks to a file.
func AddWatermarksCommand(inFile, outFile string, pageSelection []string, wm *pdf.Watermark, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.ADDWATERMARKS
	return &Command{
		Mode:          pdf.ADDWATERMARKS,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Watermark:     wm,
		Conf:          conf}
}

// ImportImagesCommand creates a new command to import images.
func ImportImagesCommand(imageFiles []string, outFile string, imp *pdf.Import, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.IMPORTIMAGES
	return &Command{
		Mode:    pdf.IMPORTIMAGES,
		InFiles: imageFiles,
		OutFile: &outFile,
		Import:  imp,
		Conf:    conf}
}

// InsertPagesCommand creates a new command to insert a blank page before selected pages.
func InsertPagesCommand(inFile, outFile string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.INSERTPAGES
	return &Command{
		Mode:          pdf.INSERTPAGES,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// RemovePagesCommand creates a new command to remove selected pages.
func RemovePagesCommand(inFile, outFile string, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.REMOVEPAGES
	return &Command{
		Mode:          pdf.REMOVEPAGES,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// RotateCommand creates a new command to rotate pages.
func RotateCommand(inFile, outFile string, rotation int, pageSelection []string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.ROTATE
	return &Command{
		Mode:          pdf.ROTATE,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Rotation:      rotation,
		Conf:          conf}
}

// NUpCommand creates a new command to render PDFs or image files in n-up fashion.
func NUpCommand(inFiles []string, outFile string, pageSelection []string, nUp *pdf.NUp, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.NUP
	return &Command{
		Mode:          pdf.NUP,
		InFiles:       inFiles,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		NUp:           nUp,
		Conf:          conf}
}

// InfoCommand creates a new command to output information about inFile.
func InfoCommand(inFile string, conf *pdf.Configuration) *Command {
	if conf == nil {
		conf = pdf.NewDefaultConfiguration()
	}
	conf.Cmd = pdf.INFO
	return &Command{
		Mode:   pdf.INFO,
		InFile: &inFile,
		Conf:   conf}
}
