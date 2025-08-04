package zlog

import "go.uber.org/zap"

type Field = zap.Field

var (
	SkipField        = zap.Skip
	BinaryField      = zap.Binary
	BoolField        = zap.Bool
	BoolpField       = zap.Boolp
	ByteStringField  = zap.ByteString
	Complex128Field  = zap.Complex128
	Complex128pField = zap.Complex128p
	Complex64Field   = zap.Complex64
	Complex64pField  = zap.Complex64p
	Float64Field     = zap.Float64
	Float64pField    = zap.Float64p
	Float32Field     = zap.Float32
	Float32pField    = zap.Float32p
	IntField         = zap.Int
	IntpField        = zap.Intp
	Int64Field       = zap.Int64
	Int64pField      = zap.Int64p
	Int32Field       = zap.Int32
	Int32pField      = zap.Int32p
	Int16Field       = zap.Int16
	Int16pField      = zap.Int16p
	Int8Field        = zap.Int8
	Int8pField       = zap.Int8p
	StringField      = zap.String
	StringpField     = zap.Stringp
	UintField        = zap.Uint
	UintpField       = zap.Uintp
	Uint64Field      = zap.Uint64
	Uint64pField     = zap.Uint64p
	Uint32Field      = zap.Uint32
	Uint32pField     = zap.Uint32p
	Uint16Field      = zap.Uint16
	Uint16pField     = zap.Uint16p
	Uint8Field       = zap.Uint8
	Uint8pField      = zap.Uint8p
	UintptrField     = zap.Uintptr
	UintptrpField    = zap.Uintptrp
	ReflectField     = zap.Reflect
	NamespaceField   = zap.Namespace
	StringerField    = zap.Stringer
	TimeField        = zap.Time
	TimepField       = zap.Timep
	StackField       = zap.Stack
	StackSkipField   = zap.StackSkip
	DurationField    = zap.Duration
	DurationpField   = zap.Durationp
	ObjectField      = zap.Object
	InlineField      = zap.Inline
	ErrorField       = zap.Error
	AnyField         = zap.Any
)
