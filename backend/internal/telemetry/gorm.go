package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GormOTelPlugin is a custom, lightweight GORM v2 database tracing plugin.
// It hooks Query, Create, Update, and Delete transactions to trace execution and preloads inside SQLite.
type GormOTelPlugin struct {
	tracer trace.Tracer
}

// NewGormOTelPlugin instantiates a new GormOTelPlugin.
func NewGormOTelPlugin() *GormOTelPlugin {
	return &GormOTelPlugin{
		tracer: otel.GetTracerProvider().Tracer("northstar-gorm"),
	}
}

// Name returns the identifier of the tracing plugin.
func (p *GormOTelPlugin) Name() string {
	return "gorm-otel"
}

// Initialize registers standard GORM hook lifecycle callbacks.
func (p *GormOTelPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for Create
	_ = db.Callback().Create().Before("gorm:create").Register("otel:before_create", p.before("gorm.create"))
	_ = db.Callback().Create().After("gorm:create").Register("otel:after_create", p.after)

	// Register callbacks for Query
	_ = db.Callback().Query().Before("gorm:query").Register("otel:before_query", p.before("gorm.query"))
	_ = db.Callback().Query().After("gorm:query").Register("otel:after_query", p.after)

	// Register callbacks for Update
	_ = db.Callback().Update().Before("gorm:update").Register("otel:before_update", p.before("gorm.update"))
	_ = db.Callback().Update().After("gorm:update").Register("otel:after_update", p.after)

	// Register callbacks for Delete
	_ = db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", p.before("gorm.delete"))
	_ = db.Callback().Delete().After("gorm:delete").Register("otel:after_delete", p.after)

	return nil
}

func (p *GormOTelPlugin) before(op string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement == nil || db.Statement.Context == nil {
			return
		}
		// Start a new SQLite trace span linked to parent span context
		ctx, span := p.tracer.Start(db.Statement.Context, op, trace.WithSpanKind(trace.SpanKindClient))
		db.Statement.Context = ctx
		span.SetAttributes(attribute.String("db.system", "sqlite"))
	}
}

func (p *GormOTelPlugin) after(db *gorm.DB) {
	if db.Statement == nil || db.Statement.Context == nil {
		return
	}
	span := trace.SpanFromContext(db.Statement.Context)
	defer span.End()

	// Capture errors (excluding standard GORM ErrRecordNotFound)
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		span.RecordError(db.Error)
		span.SetStatus(2, db.Error.Error()) // 2 corresponds to codes.Error
	} else {
		span.SetStatus(1, "") // 1 corresponds to codes.Ok
	}

	// Capture database statement details if available
	if db.Statement.SQL.String() != "" {
		span.SetAttributes(attribute.String("db.statement", db.Statement.SQL.String()))
	}
}
