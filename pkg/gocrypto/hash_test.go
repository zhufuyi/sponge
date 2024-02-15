package gocrypto

import (
	"crypto"
	"testing"
)

var hashRawData = []byte("hash_abcdefghijklmnopqrstuvwxyz0123456789")

func TestMd5(t *testing.T) {
	val := Md5(hashRawData)
	want := "98c0e2e94366eed32398f972e9742f4e"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha1(t *testing.T) {
	val := Sha1(hashRawData)
	want := "fec5700e21f47cb04127424cc09c99322925c15d"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha256(t *testing.T) {
	val := Sha256(hashRawData)
	want := "229c782bcccf23fb5e2a3f382b388df3d8edaa5502ace49ab6c80976023ad637"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func TestSha512(t *testing.T) {
	val := Sha512(hashRawData)
	want := "c1871959522cac1004ee87aaf0111d1b4569e07ff30673929e3691b119bc635960cbe63ab0ffba5acb6976a6110bb45f7cd56916662d595eac754c5f191cedfe"
	if val != want {
		t.Fatalf("got %v, want %v", val, want)
	}
	t.Log(val)
}

func BenchmarkMd5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Md5(hashRawData)
	}
}

func BenchmarkSha1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha1(hashRawData)
	}
}

func BenchmarkSha256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha256(hashRawData)
	}
}

func BenchmarkSha512(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sha512(hashRawData)
	}
}

func TestHash(t *testing.T) {
	type args struct {
		hashType crypto.Hash
		rawData  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "MD5",
			args: args{
				hashType: crypto.MD5,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA1",
			args: args{
				hashType: crypto.SHA1,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA224",
			args: args{
				hashType: crypto.SHA224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA256",
			args: args{
				hashType: crypto.SHA256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA384",
			args: args{
				hashType: crypto.SHA384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512",
			args: args{
				hashType: crypto.SHA512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "MD5SHA1",
			args: args{
				hashType: crypto.MD5SHA1,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_224",
			args: args{
				hashType: crypto.SHA3_224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_256",
			args: args{
				hashType: crypto.SHA3_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_384",
			args: args{
				hashType: crypto.SHA3_384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA3_512",
			args: args{
				hashType: crypto.SHA3_512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512_224",
			args: args{
				hashType: crypto.SHA512_224,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "SHA512_256",
			args: args{
				hashType: crypto.SHA512_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2s_256",
			args: args{
				hashType: crypto.BLAKE2s_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_256",
			args: args{
				hashType: crypto.BLAKE2b_256,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_384",
			args: args{
				hashType: crypto.BLAKE2b_384,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
		{
			name: "BLAKE2b_512",
			args: args{
				hashType: crypto.BLAKE2b_512,
				rawData:  hashRawData,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.hashType, tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func BenchmarkHash(b *testing.B) {
	b.Run("MD4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD4, hashRawData)
		}
	})

	b.Run("MD5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD5, hashRawData)
		}
	})

	b.Run("SHA1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA1, hashRawData)
		}
	})

	b.Run("SHA224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA224, hashRawData)
		}
	})

	b.Run("SHA256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA256, hashRawData)
		}
	})

	b.Run("SHA384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA384, hashRawData)
		}
	})

	b.Run("SHA512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512, hashRawData)
		}
	})

	b.Run("MD5SHA1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.MD5SHA1, hashRawData)
		}
	})

	b.Run("SHA3_224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_224, hashRawData)
		}
	})

	b.Run("SHA3_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_256, hashRawData)
		}
	})

	b.Run("SHA3_384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_384, hashRawData)
		}
	})

	b.Run("SHA3_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA3_512, hashRawData)
		}
	})

	b.Run("SHA512_224", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512_224, hashRawData)
		}
	})

	b.Run("SHA512_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.SHA512_256, hashRawData)
		}
	})

	b.Run("BLAKE2s_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2s_256, hashRawData)
		}
	})

	b.Run("BLAKE2b_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_256, hashRawData)
		}
	})

	b.Run("BLAKE2b_384", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_384, hashRawData)
		}
	})

	b.Run("BLAKE2b_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Hash(crypto.BLAKE2b_512, hashRawData)
		}
	})
}
