package lightstep

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go"
)

// An Event is emitted by the LightStep tracer as a reporting mechanism. They are
// handled by registering an EventHandler callback via SetGlobalEventHandler. The
// emitted events may be cast to specific event types in order access additional
// information.
//
// NOTE: To ensure that events can be accurately identified, each event type contains
// a sentinel method matching the name of the type. This method is a no-op, it is
// only used for type coersion.
type Event interface {
	Event()
	String() string
}

// The ErrorEvent type can be used to filter events. The `Err` method
// returns the underlying error.
type ErrorEvent interface {
	Event
	error
	Err() error
}

// EventStartError occurs if the Options passed to NewTracer are invalid, and
// the Tracer has failed to start.
type EventStartError interface {
	ErrorEvent
	EventStartError()
}

type eventStartError struct {
	err error
}

func newEventStartError(err error) *eventStartError {
	return &eventStartError{err: err}
}

func (*eventStartError) Event()           {}
func (*eventStartError) EventStartError() {}

func (e *eventStartError) String() string {
	return e.err.Error()
}

func (e *eventStartError) Error() string {
	return e.err.Error()
}

func (e *eventStartError) Err() error {
	return e.err
}

// EventFlushErrorState lists the possible causes for a flush to fail.
type EventFlushErrorState string

// Constant strings corresponding to flush errors
const (
	FlushErrorTracerClosed   EventFlushErrorState = "flush failed, the tracer is closed."
	FlushErrorTracerDisabled EventFlushErrorState = "flush failed, the tracer is disabled."
	FlushErrorTransport      EventFlushErrorState = "flush failed, could not send report to Collector"
	FlushErrorReport         EventFlushErrorState = "flush failed, report contained errors"
	FlushErrorTranslate      EventFlushErrorState = "flush failed, could not translate report"
)

var (
	errFlushFailedTracerClosed = errors.New(string(FlushErrorTracerClosed))
)

// EventFlushError occurs when a flush fails to send. Call the `State` method to
// determine the type of error.
type EventFlushError interface {
	ErrorEvent
	EventFlushError()
	State() EventFlushErrorState
}

type eventFlushError struct {
	err   error
	state EventFlushErrorState
}

func newEventFlushError(err error, state EventFlushErrorState) *eventFlushError {
	return &eventFlushError{err: err, state: state}
}

func (*eventFlushError) Event()           {}
func (*eventFlushError) EventFlushError() {}

func (e *eventFlushError) State() EventFlushErrorState {
	return e.state
}

func (e *eventFlushError) String() string {
	return e.Error()
}

func (e *eventFlushError) Error() string {
	return e.Err().Error()
}

func (e *eventFlushError) Err() error {
	return e.err
}

// EventConnectionError occurs when the tracer fails to maintain it's connection
// with the Collector.
type EventConnectionError interface {
	ErrorEvent
	EventConnectionError()
}

type eventConnectionError struct {
	err error
}

func newEventConnectionError(err error) *eventConnectionError {
	return &eventConnectionError{err: err}
}

func (*eventConnectionError) Event()                {}
func (*eventConnectionError) EventConnectionError() {}

func (e *eventConnectionError) String() string {
	return e.Error()
}

func (e *eventConnectionError) Error() string {
	return e.Err().Error()
}

func (e *eventConnectionError) Err() error {
	return e.err
}

// EventStatusReport occurs on every successful flush. It contains all metrics
// collected since the previous successful flush.
type EventStatusReport interface {
	Event
	EventStatusReport()

	// StartTime is the earliest time a span was added to the report buffer.
	StartTime() time.Time

	// FinishTime is the latest time a span was added to the report buffer.
	FinishTime() time.Time

	// Duration is time between StartTime and FinishTime.
	Duration() time.Duration

	// SentSpans is the number of spans sent in the report buffer.
	SentSpans() int

	// DroppedSpans is the number of spans dropped that did not make it into
	// the report buffer.
	DroppedSpans() int

	// EncodingErrors is the number of encoding errors that occurred while
	// building the report buffer.
	EncodingErrors() int

	// FlushDuration is the time it took to send the report, including encoding,
	// buffer rotation, and network time.
	FlushDuration() time.Duration
}

// MetricEventStatusReport occurs every time metrics are sent successfully.. It
// contains all metrics collected since the previous successful flush.
type MetricEventStatusReport interface {
	Event
	MetricEventStatusReport()

	// StartTime is the earliest time a span was added to the report buffer.
	StartTime() time.Time

	// FinishTime is the latest time a span was added to the report buffer.
	FinishTime() time.Time

	// SentMetrics is the number of metrics sent in the report buffer.
	SentMetrics() int
}

type eventStatusReport struct {
	startTime      time.Time
	finishTime     time.Time
	sentSpans      int
	droppedSpans   int
	encodingErrors int
	flushDuration  time.Duration
}

func newEventStatusReport(
	startTime, finishTime time.Time,
	sentSpans, droppedSpans, encodingErrors int,
	flushDuration time.Duration,
) *eventStatusReport {
	return &eventStatusReport{
		startTime:      startTime,
		finishTime:     finishTime,
		sentSpans:      sentSpans,
		droppedSpans:   droppedSpans,
		encodingErrors: encodingErrors,
		flushDuration:  flushDuration,
	}
}

func (*eventStatusReport) Event() {}

func (*eventStatusReport) EventStatusReport() {}

func (s *eventStatusReport) SetSentSpans(sent int) {
	s.sentSpans = sent
}

func (s *eventStatusReport) StartTime() time.Time {
	return s.startTime
}

func (s *eventStatusReport) FinishTime() time.Time {
	return s.finishTime
}

func (s *eventStatusReport) Duration() time.Duration {
	return s.finishTime.Sub(s.startTime)
}

func (s *eventStatusReport) FlushDuration() time.Duration {
	return s.flushDuration
}

func (s *eventStatusReport) SentSpans() int {
	return s.sentSpans
}

func (s *eventStatusReport) DroppedSpans() int {
	return s.droppedSpans
}

func (s *eventStatusReport) EncodingErrors() int {
	return s.encodingErrors
}

func (s *eventStatusReport) String() string {
	return fmt.Sprint(
		"STATUS REPORT start: ", s.startTime,
		", end: ", s.finishTime,
		", dropped spans: ", s.droppedSpans,
		", encoding errors: ", s.encodingErrors,
	)
}

// EventUnsupportedTracer occurs when a tracer being passed to a helper function
// fails to typecast as a LightStep tracer.
type EventUnsupportedTracer interface {
	ErrorEvent
	EventUnsupportedTracer()
	Tracer() opentracing.Tracer
}

type eventUnsupportedTracer struct {
	tracer opentracing.Tracer
	err    error
}

func newEventUnsupportedTracer(tracer opentracing.Tracer) EventUnsupportedTracer {
	return &eventUnsupportedTracer{
		tracer: tracer,
		err:    fmt.Errorf("unsupported tracer type: %v", reflect.TypeOf(tracer)),
	}
}

func (e *eventUnsupportedTracer) Event()                  {}
func (e *eventUnsupportedTracer) EventUnsupportedTracer() {}

func (e *eventUnsupportedTracer) Tracer() opentracing.Tracer {
	return e.tracer
}

func (e *eventUnsupportedTracer) String() string {
	return e.Error()
}

func (e *eventUnsupportedTracer) Error() string {
	return e.Err().Error()
}

func (e *eventUnsupportedTracer) Err() error {
	return e.err
}

// EventUnsupportedValue occurs when a tracer encounters an unserializable tag
// or log field.
type EventUnsupportedValue interface {
	ErrorEvent
	EventUnsupportedValue()
	Key() string
	Value() interface{}
}

type eventUnsupportedValue struct {
	key   string
	value interface{}
	err   error
}

func newEventUnsupportedValue(key string, value interface{}, err error) EventUnsupportedValue {
	if err == nil {
		err = fmt.Errorf(
			"value `%v` of type `%T` for key `%s` is an unsupported type",
			value, value, key,
		)
	}
	return &eventUnsupportedValue{
		key:   key,
		value: value,
		err:   err,
	}
}

func (e *eventUnsupportedValue) Event()                 {}
func (e *eventUnsupportedValue) EventUnsupportedValue() {}

func (e *eventUnsupportedValue) Key() string {
	return e.key
}

func (e *eventUnsupportedValue) Value() interface{} {
	return e.value
}

func (e *eventUnsupportedValue) String() string {
	return e.Error()
}

func (e *eventUnsupportedValue) Error() string {
	return e.Err().Error()
}

func (e *eventUnsupportedValue) Err() error {
	return e.err
}

const tracerDisabled = "the tracer has been disabled"

// EventTracerDisabled occurs when a tracer is disabled by either the user or
// the collector.
type EventTracerDisabled interface {
	Event
	EventTracerDisabled()
}

type eventTracerDisabled struct{}

func newEventTracerDisabled() EventTracerDisabled {
	return eventTracerDisabled{}
}

func (eventTracerDisabled) Event()               {}
func (eventTracerDisabled) EventTracerDisabled() {}
func (eventTracerDisabled) String() string {
	return tracerDisabled
}

type EventSystemMetricsMeasurementFailed interface {
	ErrorEvent
}

type eventSystemMetricsMeasurementFailed struct {
	err error
}

func newEventSystemMetricsMeasurementFailed(err error) *eventSystemMetricsMeasurementFailed {
	return &eventSystemMetricsMeasurementFailed{
		err: err,
	}
}

func (e *eventSystemMetricsMeasurementFailed) Event() {}

func (e *eventSystemMetricsMeasurementFailed) String() string {
	return e.Error()
}

func (e *eventSystemMetricsMeasurementFailed) Error() string {
	return e.Err().Error()
}

func (e *eventSystemMetricsMeasurementFailed) Err() error {
	return e.err
}

type eventSystemMetricsStatusReport struct {
	startTime   time.Time
	finishTime  time.Time
	sentMetrics int
}

func newEventSystemMetricsStatusReport(
	startTime, finishTime time.Time,
	sentMetrics int,
) *eventSystemMetricsStatusReport {
	return &eventSystemMetricsStatusReport{
		startTime:   startTime,
		finishTime:  finishTime,
		sentMetrics: sentMetrics,
	}
}

func (e *eventSystemMetricsStatusReport) Event() {}

func (e *eventSystemMetricsStatusReport) MetricEventStatusReport() {}

func (e *eventSystemMetricsStatusReport) String() string {
	return fmt.Sprint(
		"METRICS STATUS REPORT start: ", e.startTime,
		", end: ", e.finishTime,
		", sent metrics: ", e.sentMetrics,
	)
}

func (e *eventSystemMetricsStatusReport) StartTime() time.Time {
	return e.startTime
}

func (e *eventSystemMetricsStatusReport) FinishTime() time.Time {
	return e.finishTime
}

func (e *eventSystemMetricsStatusReport) SentMetrics() int {
	return e.sentMetrics
}
