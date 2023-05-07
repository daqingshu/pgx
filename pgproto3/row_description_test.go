package pgproto3_test

import (
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRowDescription_Decode1(t *testing.T) {
	hexString := "54000000fb000969640000006003000100000019ffffffffffff0000637265617465645f617400000060030002000004a00008ffffffff0000616374696f6e0000006003000300000019ffffffffffff00006f626a6563745f69640000006003000400000019ffffffffffff00006f626a6563745f747970650000006003000500000019ffffffffffff00007261775f6f626a6563740000006003000600000072ffffffffffff00007261775f6d6574610000006003000700000072ffffffffffff00007261775f646966660000006003000800000072ffffffffffff0000637265617465645f62790000006003000900000019ffffffffffff0000"
	decodedByteArray, err := hex.DecodeString(hexString)
	srcBytes := decodedByteArray
	rd := pgproto3.RowDescription{}
	err = rd.Decode(srcBytes[5:])
	u, _ := rd.MarshalJSON()
	fmt.Println(string(u))
	assert.NoError(t, err, "No errors on decode")
	dstBytes := []byte{}
	dstBytes = rd.Encode(dstBytes)
	fmt.Println(dstBytes)
	assert.EqualValues(t, []byte{}, dstBytes, "Expecting src & dest bytes to match")
}

func TestRowDescription_Decode(t *testing.T) {
	hexString := "54000000fb000969640000006003000100000019ffffffffffff0000637265617465645f617400000060030002000004a00008ffffffff0000616374696f6e0000006003000300000019ffffffffffff00006f626a6563745f69640000006003000400000019ffffffffffff00006f626a6563745f747970650000006003000500000019ffffffffffff00007261775f6f626a6563740000006003000600000072ffffffffffff00007261775f6d6574610000006003000700000072ffffffffffff00007261775f646966660000006003000800000072ffffffffffff0000637265617465645f62790000006003000900000019ffffffffffff0000"
	decodedByteArray, _ := hex.DecodeString(hexString)

	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		fields  pgproto3.RowDescription
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{
				{Name: []byte("id"), TableOID: 24579, TableAttributeNumber: 1, DataTypeOID: 25, DataTypeSize: -1, TypeModifier: -1, Format: 0},
				{Name: []byte("created_at"), TableOID: 24579, TableAttributeNumber: 2, DataTypeOID: 1184, DataTypeSize: 8, TypeModifier: -1, Format: 0},
				{Name: []byte("action"),TableOID:24579,TableAttributeNumber:3,DataTypeOID:25,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("object_id"),TableOID:24579,TableAttributeNumber:4,DataTypeOID:25,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("object_type"),TableOID:24579,TableAttributeNumber:5,DataTypeOID:25,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("raw_object"),TableOID:24579,TableAttributeNumber:6,DataTypeOID:114,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("raw_meta"),TableOID:24579,TableAttributeNumber:7,DataTypeOID:114,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("raw_diff"),TableOID:24579,TableAttributeNumber:8,DataTypeOID:114,DataTypeSize:-1,TypeModifier:-1,Format:0},
				{Name: []byte("created_by"),TableOID:24579,TableAttributeNumber:9,DataTypeOID:25,DataTypeSize:-1,TypeModifier:-1,Format:0},
				},
			},
			args: args{src: decodedByteArray},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &pgproto3.RowDescription{
				Fields: tt.fields.Fields,
			}
			tt.wantErr(t, dst.Decode(tt.args.src), fmt.Sprintf("Decode(%v)", tt.args.src))
		})
	}
}