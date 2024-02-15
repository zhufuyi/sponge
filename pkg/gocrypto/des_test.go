package gocrypto

import "testing"

var (
	desRawData = []byte("des_abcdefghijklmnopqrstuvwxyz0123456789")
	deskey     = []byte("ABCDEFGH")
)

func TestDes(t *testing.T) {
	want := desRawData

	// ECB default mod and key
	t.Run("default des ebc", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData)
		got, _ := DesDecrypt(cypherData)
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// ECB
	t.Run("des ecb", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeECB())
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeECB())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CBC
	t.Run("des cbc", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCBC())
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCBC())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CFB
	t.Run("des cfb", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCFB())
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCFB())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})

	// CTR
	t.Run("des ctr", func(t *testing.T) {
		cypherData, _ := DesEncrypt(desRawData, WithDesKey(deskey), WithDesModeCTR())
		got, _ := DesDecrypt(cypherData, WithDesKey(deskey), WithDesModeCTR())
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", desRawData, cypherData)
	})
}

func BenchmarkDes(b *testing.B) {
	b.Run("des ecb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeECB())
		}
	})
	b.Run("des ecb decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeECB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeECB())
		}
	})

	b.Run("des cbc encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCBC())
		}
	})
	b.Run("des cbc decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCBC())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCBC())
		}
	})

	b.Run("des cfb encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCFB())
		}
	})
	b.Run("des cfb decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCFB())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCFB())
		}
	})

	b.Run("des ctr encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DesEncrypt(desRawData, WithDesModeCTR())
		}
	})
	b.Run("des ctr decrypt", func(b *testing.B) {
		cypherData, err := DesEncrypt(desRawData, WithDesModeCTR())
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			DesDecrypt(cypherData, WithDesModeCTR())
		}
	})
}

func TestDesHex(t *testing.T) {
	want := string(desRawData)

	t.Run("default des ecb", func(t *testing.T) {
		cypherStr, _ := DesEncryptHex(string(desRawData))
		got, _ := DesDecryptHex(cypherStr)
		if got != want {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%s]", desRawData, cypherStr)
	})
}
