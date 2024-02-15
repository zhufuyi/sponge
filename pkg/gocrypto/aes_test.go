package gocrypto

import (
	"testing"
)

var (
	aesRawData = []byte("aes_abcdefghijklmnopqrstuvwxyz0123456789")
	aesKey     = []byte("aeskey0123456789aeskey0123456789")
)

func TestAes(t *testing.T) {
	want := aesRawData

	// ECB default mod and key
	t.Run("default aes ebc", func(t *testing.T) {
		cypherData, _ := AesEncrypt(aesRawData)
		got, _ := AesDecrypt(cypherData)
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})

	// ECB
	t.Run("aes ecb", func(t *testing.T) {
		cypherData, _ := AesEncrypt(aesRawData, WithAesKey(aesKey), WithAesModeECB())
		got, _ := AesDecrypt(cypherData, WithAesKey(aesKey), WithAesModeECB())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})

	// CBC
	t.Run("aes cbc", func(t *testing.T) {
		cypherData, _ := AesEncrypt(aesRawData, WithAesKey(aesKey), WithAesModeCBC())
		got, _ := AesDecrypt(cypherData, WithAesKey(aesKey), WithAesModeCBC())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})

	// CFB
	t.Run("aes cfb", func(t *testing.T) {
		cypherData, _ := AesEncrypt(aesRawData, WithAesKey(aesKey), WithAesModeCFB())
		got, _ := AesDecrypt(cypherData, WithAesKey(aesKey), WithAesModeCFB())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})

	// CTR
	t.Run("aes ctr", func(t *testing.T) {
		cypherData, _ := AesEncrypt(aesRawData, WithAesKey(aesKey), WithAesModeCTR())
		got, _ := AesDecrypt(cypherData, WithAesKey(aesKey), WithAesModeCTR())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})
}

func TestAesHex(t *testing.T) {
	want := string(aesRawData)

	// ecb default mod and key
	t.Run("default aes ecb", func(t *testing.T) {
		cypherData, _ := AesEncryptHex(string(aesRawData))
		got, _ := AesDecryptHex(cypherData)
		if got != want {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", aesRawData, cypherData)
	})
}

func BenchmarkAes(b *testing.B) {
	b.Run("aes ecb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AesEncrypt(aesRawData, WithAesModeECB())
		}
	})
	b.Run("aes ecb decrypt", func(b *testing.B) {
		cypherData, err := AesEncrypt(aesRawData, WithAesModeECB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			AesDecrypt(cypherData, WithAesModeECB())
		}
	})

	b.Run("aes cbc encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AesEncrypt(aesRawData, WithAesModeCBC())
		}
	})
	b.Run("aes cbc decrypt", func(b *testing.B) {
		cypherData, err := AesEncrypt(aesRawData, WithAesModeCBC())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			AesDecrypt(cypherData, WithAesModeCBC())
		}
	})

	b.Run("aes cfb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AesEncrypt(aesRawData, WithAesModeCFB())
		}
	})
	b.Run("aes cfb decrypt", func(b *testing.B) {
		cypherData, err := AesEncrypt(aesRawData, WithAesModeCFB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			AesDecrypt(cypherData, WithAesModeCFB())
		}
	})

	b.Run("aes ctr encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AesEncrypt(aesRawData, WithAesModeCTR())
		}
	})
	b.Run("aes ctr decrypt", func(b *testing.B) {
		cypherData, err := AesEncrypt(aesRawData, WithAesModeCTR())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			AesDecrypt(cypherData, WithAesModeCTR())
		}
	})
}
