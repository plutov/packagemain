package survey

func ExportSurvey(s *Survey, exporter Exporter) error {
	return exporter.Export(s)
}

type Exporter interface {
	Export(s *Survey) error
}

type S3Exporter struct{}

func (e *S3Exporter) Export(s *Survey) error {
	return nil
}

type GCSExporter struct{}

func (e *GCSExporter) Export(s *Survey) error {
	return nil
}
