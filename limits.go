package cloudevents

// Limits selects decoder byte/count/depth admission, not constructor/encoder
// limits or a total heap/CPU/time budget. Supply DefaultLimits explicitly; zero
// fields are not filled automatically and fail at their applicable checks.
// JSON token/raw-message materialization can precede attribute/batch admission.
// Kafka-specific fields cover even unowned metadata copied by DecodeKafka.
// HTTP leaves unrelated headers to the caller and requires reader cancellation;
// MaxAttributes applies to binary HTTP, not the redundant structured-header count.
type Limits struct {
	// MaxEventBytes caps one structured event/whole JSON batch input; binary HTTP/Kafka payload uses MaxDataBytes.
	MaxEventBytes int64
	// MaxDataBytes caps received data bytes after the applicable wire admission; it is not an encode cap.
	MaxDataBytes int64
	// MaxAttributes caps JSON context/binary HTTP/Kafka context counts; structured HTTP redundant headers have no separate count check.
	MaxAttributes int
	// MaxAttributeNameBytes caps decoded binding context-name bytes at their applicable checks.
	MaxAttributeNameBytes int
	// MaxAttributeValueBytes caps selected encoded/decoded context values, not unrelated HTTP headers.
	MaxAttributeValueBytes int
	// MaxDepth caps decoded JSON container depth; batch adds its outer array level.
	MaxDepth int
	// MaxBatchEvents caps decoded batch count after raw event materialization.
	MaxBatchEvents int
	// MaxKafkaKeyBytes caps copied record-key bytes.
	MaxKafkaKeyBytes int
	// MaxKafkaHeaders caps all received record headers, including transport-owned metadata.
	MaxKafkaHeaders int
	// MaxKafkaHeaderNameBytes caps each received Kafka header name before retained copies.
	MaxKafkaHeaderNameBytes int
	// MaxKafkaHeaderValueBytes caps each received Kafka header value before retained copies.
	MaxKafkaHeaderValueBytes int
}

// DefaultLimits accepts the CloudEvents interoperability floor while bounding
// selected input/retention sizes for ordinary library use, not all temporary allocations.
func DefaultLimits() Limits {
	return Limits{
		MaxEventBytes:            1 << 20,
		MaxDataBytes:             1 << 20,
		MaxAttributes:            64,
		MaxAttributeNameBytes:    64,
		MaxAttributeValueBytes:   16 << 10,
		MaxDepth:                 64,
		MaxBatchEvents:           1000,
		MaxKafkaKeyBytes:         1 << 20,
		MaxKafkaHeaders:          64,
		MaxKafkaHeaderNameBytes:  256,
		MaxKafkaHeaderValueBytes: 16 << 10,
	}
}
