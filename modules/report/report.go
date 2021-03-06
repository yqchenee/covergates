package report

import (
	"bytes"
	"fmt"
	"io"

	"github.com/covergates/covergates/core"
)

const (
	upArrow   = ":arrow_up_small:"
	downArrow = ":arrow_down_small:"
)

type filesMap map[string]*core.File

// Service of Report
type Service struct{}

// DiffReports coverage differences
func (service *Service) DiffReports(source, target *core.Report) (*core.CoverageReportDiff, error) {
	var m filesMap
	if target != nil && target.Coverage != nil {
		m = toFilesMap(target.Coverage.Files)
	}
	diffFiles := make([]*core.FileDiff, 0, len(source.Coverage.Files))
	for _, file := range source.Coverage.Files {
		diff := &core.FileDiff{
			File:                  file,
			StatementCoverageDiff: file.StatementCoverage,
		}
		if f, ok := m[file.Name]; ok {
			diff.StatementCoverageDiff -= f.StatementCoverage
			delete(m, file.Name)
		}
		diffFiles = append(diffFiles, diff)
	}
	for name := range m {
		diff := &core.FileDiff{
			File:                  m[name],
			StatementCoverageDiff: -m[name].StatementCoverage,
			Removed:               true,
		}
		diffFiles = append(diffFiles, diff)
	}
	coverageDiff := source.Coverage.StatementCoverage
	if target != nil && target.Coverage != nil {
		coverageDiff -= target.Coverage.StatementCoverage
	}

	return &core.CoverageReportDiff{
		StatementCoverageDiff: coverageDiff,
		Files:                 diffFiles,
	}, nil
}

// MarkdownReport generates coverage summary report in markdown format
func (service *Service) MarkdownReport(source, target *core.Report) (io.Reader, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("### Coverage: %.1f%%\n\n", source.Coverage.StatementCoverage*100))
	buf.WriteString("||File|Coverage|\n")
	buf.WriteString("|--|--|--------|\n")
	diff, err := service.DiffReports(source, target)
	if err != nil {
		return nil, err
	}
	for _, file := range diff.Files {
		if file.Removed {
			continue
		}
		mark := ""
		if file.StatementCoverageDiff > 0 {
			mark = upArrow
		} else if file.StatementCoverageDiff < 0 {
			mark = downArrow
		}

		buf.WriteString(fmt.Sprintf("|%s|%s|%.2f|\n", mark, file.File.Name, file.File.StatementCoverage))
	}
	return buf, nil
}

// MergeReport of two coverage reports, the report types are not necessary to be equal
func (service *Service) MergeReport(from, to *core.Report, changes []*core.FileChange) (*core.Report, error) {
	deleted := make(map[string]bool)
	for _, change := range changes {
		if change.Deleted {
			deleted[change.Path] = true
		}
	}
	target := toFilesMap(to.Coverage.Files)
	for _, file := range from.Coverage.Files {
		if _, ok := deleted[file.Name]; ok {
			continue
		}
		if _, ok := target[file.Name]; ok {
			continue
		}
		to.Coverage.Files = append(to.Coverage.Files, file)
	}
	to.Coverage.StatementCoverage = to.Coverage.AvgStatementCoverage()
	return to, nil
}

func toFilesMap(files []*core.File) filesMap {
	m := make(filesMap)
	for _, file := range files {
		m[file.Name] = file
	}
	return m
}
