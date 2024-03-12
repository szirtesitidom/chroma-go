package types

import (
	"context"
	"reflect"
	"testing"

	openapi "github.com/amikos-tech/chroma-go/swagger"
	"github.com/amikos-tech/chroma-go/where"
	wheredoc "github.com/amikos-tech/chroma-go/where_document"
)

func TestConsistentHashEmbeddingFunction_EmbedDocuments(t *testing.T) {
	type args struct {
		documents []string
		dim       int
	}
	tests := []struct {
		name    string
		args    args
		want    func(got []*Embedding) bool
		wantErr bool
	}{
		{
			name:    "empty document list, expect empty embeddings list",
			args:    args{documents: []string{}, dim: 10},
			want:    func(got []*Embedding) bool { return len(got) == 0 },
			wantErr: false,
		},
		{
			name:    "with single document, expect single embedding",
			args:    args{documents: []string{"test document"}, dim: 10},
			want:    func(got []*Embedding) bool { return len(got) == 1 },
			wantErr: false,
		},
		{
			name:    "with single document and 100 dims",
			args:    args{documents: []string{"test document"}, dim: 100},
			want:    func(got []*Embedding) bool { return len(got) == 1 && len(*got[0].ArrayOfFloat32) == 100 },
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := &ConsistentHashEmbeddingFunction{dim: tt.args.dim}
			got, err := ef.EmbedDocuments(context.TODO(), tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbedDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want(got) {
				t.Errorf("generateServiceId() got = %v, want %v", got, tt.want(got))
			}
		})
	}
}

func TestConsistentHashEmbeddingFunction_EmbedQuery(t *testing.T) {
	type args struct {
		text string
		dim  int
	}
	tests := []struct {
		name    string
		args    args
		want    func(got *Embedding) bool
		wantErr bool
	}{
		{
			name:    "empty text, expect empty embedding",
			args:    args{text: "", dim: 10},
			want:    func(got *Embedding) bool { return got == nil },
			wantErr: true,
		},
		{
			name:    "with single document, expect single embedding",
			args:    args{text: "test document", dim: 10},
			want:    func(got *Embedding) bool { return len(*got.ArrayOfFloat32) == 10 },
			wantErr: false,
		},
		{
			name:    "with single document and 100 dims",
			args:    args{text: "test document", dim: 100},
			want:    func(got *Embedding) bool { return len(*got.ArrayOfFloat32) == 100 },
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := &ConsistentHashEmbeddingFunction{dim: tt.args.dim}
			got, err := ef.EmbedQuery(context.TODO(), tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want(got) {
				t.Errorf("EmbedQuery() got = %v, want %v", got, tt.want(got))
			}
		})
	}
}

func TestConsistentHashEmbeddingFunction_EmbedDocuments1(t *testing.T) {
	type fields struct {
		dim int
	}
	type args struct {
		ctx       context.Context
		documents []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Embedding
		wantErr bool
	}{
		{name: "empty document list, expect empty embeddings list",
			fields:  fields{dim: 10},
			args:    args{documents: []string{}, ctx: context.TODO()},
			want:    []*Embedding{},
			wantErr: false,
		},
		{name: "with single document, expect single embedding",
			fields:  fields{dim: 10},
			args:    args{documents: []string{"test document"}, ctx: context.TODO()},
			want:    []*Embedding{{ArrayOfFloat32: &[]float32{0.26666668, 0.53333336, .2, 0.46666667, 0.26666668, 0.46666667, 0.6, 0.06666667, 0.13333334, 0.33333334}, ArrayOfInt32: nil}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ConsistentHashEmbeddingFunction{
				dim: tt.fields.dim,
			}
			got, err := e.EmbedDocuments(tt.args.ctx, tt.args.documents)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbedDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmbedDocuments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsistentHashEmbeddingFunction_EmbedQuery1(t *testing.T) {
	type fields struct {
		dim int
	}
	type args struct {
		ctx      context.Context
		document string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Embedding
		wantErr bool
	}{
		{name: "empty document list, expect empty embeddings list",
			fields:  fields{dim: 10},
			args:    args{document: "", ctx: context.TODO()},
			want:    nil,
			wantErr: true,
		},
		{name: "with single document, expect single embedding",
			fields:  fields{dim: 10},
			args:    args{document: "test document", ctx: context.TODO()},
			want:    NewEmbeddingFromFloat32([]float32{0.26666668, 0.53333336, .2, 0.46666667, 0.26666668, 0.46666667, 0.6, 0.06666667, 0.13333334, 0.33333334}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ConsistentHashEmbeddingFunction{
				dim: tt.fields.dim,
			}
			got, err := e.EmbedQuery(tt.args.ctx, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmbedQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsistentHashEmbeddingFunction_EmbedRecords(t *testing.T) {
	type fields struct {
		dim int
	}
	type args struct {
		ctx     context.Context
		records []*Record
		force   bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "empty document list, expect empty embeddings list",
			fields: fields{dim: 10},
			args: args{
				ctx:     context.TODO(),
				records: []*Record{},
			},
			wantErr: false,
		},
		{
			name:   "with single document, expect single embedding",
			fields: fields{dim: 10},
			args: args{
				ctx: context.TODO(),
				records: []*Record{
					{
						Document: "test document",
						ID:       "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "record without doc content",
			fields: fields{dim: 10},
			args: args{
				ctx: context.TODO(),
				records: []*Record{
					{
						ID: "1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ConsistentHashEmbeddingFunction{
				dim: tt.fields.dim,
			}
			if err := e.EmbedRecords(tt.args.ctx, tt.args.records, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("EmbedRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmbedRecordsDefaultImpl(t *testing.T) {
	type args struct {
		e       EmbeddingFunction
		ctx     context.Context
		records []*Record
		force   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EmbedRecordsDefaultImpl(tt.args.e, tt.args.ctx, tt.args.records, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("EmbedRecordsDefaultImpl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmbeddableContext_Apply(t *testing.T) {
	type fields struct {
		Documents []string
	}
	type args struct {
		ctx               context.Context
		embeddingFunction EmbeddingFunction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Embedding
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EmbeddableContext{
				Documents: tt.fields.Documents,
			}
			got, err := e.Apply(tt.args.ctx, tt.args.embeddingFunction)
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmbedding_GetFloat32(t *testing.T) {
	type fields struct {
		ArrayOfFloat32 *[]float32
		ArrayOfInt32   *[]int32
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Embedding{
				ArrayOfFloat32: tt.fields.ArrayOfFloat32,
				ArrayOfInt32:   tt.fields.ArrayOfInt32,
			}
			if got := e.GetFloat32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmbedding_GetInt32(t *testing.T) {
	type fields struct {
		ArrayOfFloat32 *[]float32
		ArrayOfInt32   *[]int32
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Embedding{
				ArrayOfFloat32: tt.fields.ArrayOfFloat32,
				ArrayOfInt32:   tt.fields.ArrayOfInt32,
			}
			if got := e.GetInt32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmbedding_IsDefined(t *testing.T) {
	type fields struct {
		ArrayOfFloat32 *[]float32
		ArrayOfInt32   *[]int32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Embedding{
				ArrayOfFloat32: tt.fields.ArrayOfFloat32,
				ArrayOfInt32:   tt.fields.ArrayOfInt32,
			}
			if got := e.IsDefined(); got != tt.want {
				t.Errorf("IsDefined() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmbedding_ToAPI(t *testing.T) {
	type fields struct {
		ArrayOfFloat32 *[]float32
		ArrayOfInt32   *[]int32
	}
	tests := []struct {
		name   string
		fields fields
		want   openapi.EmbeddingsInner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Embedding{
				ArrayOfFloat32: tt.fields.ArrayOfFloat32,
				ArrayOfInt32:   tt.fields.ArrayOfInt32,
			}
			if got := e.ToAPI(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestF32ToInterface(t *testing.T) {
	type args struct {
		f []float32
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := F32ToInterface(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("F32ToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidEmbeddingValueError_Error(t *testing.T) {
	type fields struct {
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &InvalidEmbeddingValueError{
				Value: tt.fields.Value,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidMetadataValueError_Error(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &InvalidMetadataValueError{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConsistentHashEmbeddingFunction(t *testing.T) {
	tests := []struct {
		name string
		want EmbeddingFunction
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsistentHashEmbeddingFunction(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsistentHashEmbeddingFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmbeddingFromAPI(t *testing.T) {
	type args struct {
		apiEmbedding openapi.EmbeddingsInner
	}
	tests := []struct {
		name string
		args args
		want *Embedding
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmbeddingFromAPI(tt.args.apiEmbedding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmbeddingFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmbeddingFromFloat32(t *testing.T) {
	type args struct {
		embedding []float32
	}
	tests := []struct {
		name string
		args args
		want *Embedding
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmbeddingFromFloat32(tt.args.embedding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmbeddingFromFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmbeddingFromInt32(t *testing.T) {
	type args struct {
		embedding []int32
	}
	tests := []struct {
		name string
		args args
		want *Embedding
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmbeddingFromInt32(tt.args.embedding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmbeddingFromInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmbeddings(t *testing.T) {
	type args struct {
		embeddings []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Embedding
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmbeddings(tt.args.embeddings)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmbeddings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmbeddings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEmbeddingsFromFloat32(t *testing.T) {
	type args struct {
		embeddings [][]float32
	}
	tests := []struct {
		name string
		args args
		want []*Embedding
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmbeddingsFromFloat32(tt.args.embeddings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmbeddingsFromFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSHA256Generator(t *testing.T) {
	tests := []struct {
		name string
		want *SHA256Generator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSHA256Generator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSHA256Generator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewULIDGenerator(t *testing.T) {
	tests := []struct {
		name string
		want *ULIDGenerator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewULIDGenerator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewULIDGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUUIDGenerator(t *testing.T) {
	tests := []struct {
		name string
		want *UUIDGenerator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUUIDGenerator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUUIDGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256Generator_Generate(t *testing.T) {
	type args struct {
		document string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SHA256Generator{}
			if got := s.Generate(tt.args.document); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToAPIEmbeddings(t *testing.T) {
	type args struct {
		embeddings []*Embedding
	}
	tests := []struct {
		name string
		args args
		want []openapi.EmbeddingsInner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAPIEmbeddings(tt.args.embeddings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAPIEmbeddings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDistanceFunction(t *testing.T) {
	type args struct {
		str any
	}
	tests := []struct {
		name    string
		args    args
		want    DistanceFunction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDistanceFunction(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDistanceFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToDistanceFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestULIDGenerator_Generate(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &ULIDGenerator{}
			if got := u.Generate(tt.args.in0); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDGenerator_Generate(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UUIDGenerator{}
			if got := u.Generate(tt.args.in0); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIds(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithIds(tt.args.ids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithInclude(t *testing.T) {
	type args struct {
		include []QueryEnum
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithInclude(tt.args.include...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithLimit(t *testing.T) {
	type args struct {
		limit int32
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLimit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithNResults(t *testing.T) {
	type args struct {
		nResults int32
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithNResults(tt.args.nResults); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOffset(t *testing.T) {
	type args struct {
		offset int32
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithOffset(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQueryEmbedding(t *testing.T) {
	type args struct {
		queryEmbedding *Embedding
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithQueryEmbedding(tt.args.queryEmbedding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithQueryEmbedding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQueryEmbeddings(t *testing.T) {
	type args struct {
		queryEmbeddings []*Embedding
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithQueryEmbeddings(tt.args.queryEmbeddings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithQueryEmbeddings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQueryText(t *testing.T) {
	type args struct {
		queryText string
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithQueryText(tt.args.queryText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithQueryText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQueryTexts(t *testing.T) {
	type args struct {
		queryTexts []string
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithQueryTexts(tt.args.queryTexts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithQueryTexts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWhere(t *testing.T) {
	type args struct {
		operation where.WhereOperation
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithWhere(tt.args.operation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWhereDocument(t *testing.T) {
	type args struct {
		operation wheredoc.WhereDocumentOperation
	}
	tests := []struct {
		name string
		args args
		want CollectionQueryOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithWhereDocument(tt.args.operation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithWhereDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWhereDocumentMap(t *testing.T) {
	type args struct {
		where        map[string]interface{}
		queryBuilder *CollectionQueryBuilder
	}
	tests := []struct {
		name string
		args args
		want *CollectionQueryBuilder
	}{
		{
			name: "with where document map",
			args: args{where: map[string]interface{}{"test": "test"}, queryBuilder: &CollectionQueryBuilder{}},
			want: &CollectionQueryBuilder{WhereDocument: map[string]interface{}{"test": "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _ = WithWhereDocumentMap(tt.args.where)(tt.args.queryBuilder); !reflect.DeepEqual(tt.args.queryBuilder, tt.want) {
				t.Errorf("WithWhereDocumentMap() = %v, want %v", tt.args.queryBuilder, tt.want)
			}
		})
	}
}

func TestWithWhereMap(t *testing.T) {
	type args struct {
		where        map[string]interface{}
		queryBuilder *CollectionQueryBuilder
	}
	tests := []struct {
		name string
		args args
		want *CollectionQueryBuilder
	}{
		{
			name: "with where map",
			args: args{where: map[string]interface{}{"test": "test"}, queryBuilder: &CollectionQueryBuilder{}},
			want: &CollectionQueryBuilder{Where: map[string]interface{}{"test": "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _ = WithWhereMap(tt.args.where)(tt.args.queryBuilder); !reflect.DeepEqual(tt.args.queryBuilder, tt.want) {
				t.Errorf("WithWhereMap() = %v, want %v", tt.args.queryBuilder, tt.want)
			}
		})
	}
}
