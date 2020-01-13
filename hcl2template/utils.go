package hcl2template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func warningErrorsToDiags(block *hcl.Block, warnings []string, err error) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, warning := range warnings {
		diags = append(diags, &hcl.Diagnostic{
			Summary:  warning,
			Subject:  &block.DefRange,
			Severity: hcl.DiagWarning,
		})
	}
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Summary:  err.Error(),
			Subject:  &block.DefRange,
			Severity: hcl.DiagError,
		})
	}
	return diags
}

func isDir(name string) (bool, error) {
	s, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

func GetHCL2Files(filename, hclSuffix, jsonSuffix string) (hclFiles, jsonFiles []string, diags hcl.Diagnostics) {
	isDir, err := isDir(filename)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Cannot tell wether " + filename + " is a directory",
			Detail:   err.Error(),
		})
		return nil, nil, diags
	}
	if !isDir {
		if strings.HasSuffix(filename, jsonSuffix) {
			return nil, []string{filename}, diags
		}
		if strings.HasSuffix(filename, hclSuffix) {
			return []string{filename}, nil, diags
		}
		return nil, nil, diags
	}

	fileInfos, err := ioutil.ReadDir(filename)
	if err != nil {
		diag := &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Cannot read hcl directory",
			Detail:   err.Error(),
		}
		diags = append(diags, diag)
		return nil, nil, diags
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filename := filepath.Join(filename, fileInfo.Name())
		if strings.HasSuffix(filename, hclSuffix) {
			hclFiles = append(hclFiles, filename)
		} else if strings.HasSuffix(filename, jsonSuffix) {
			jsonFiles = append(jsonFiles, filename)
		}
	}

	return hclFiles, jsonFiles, diags
}
