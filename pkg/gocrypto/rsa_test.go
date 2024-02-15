package gocrypto

import (
	"crypto"
	"testing"
)

var (
	rsaRawData = []byte("rsa_abcdefghijklmnopqrstuvwxyz0123456789")
)

// PKCS#1
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNzg5i/VN3w5dDu1W+U4yCgRaL
kubJbCwi/RitEgRoV8OHhNiZUmpVZfqBIxIZMPrFnx1zTC2mto7BxtesbS9F3vW3
xggpuNIMjXeLD63mK0LSJ2VhNZ0YihpJ/eVCO439mDM7vtP1JQ4KveRMmAEIql1l
Im5/SiBYqiA5JP0XMwIDAQAB
-----END PUBLIC KEY-----
`)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCNzg5i/VN3w5dDu1W+U4yCgRaLkubJbCwi/RitEgRoV8OHhNiZ
UmpVZfqBIxIZMPrFnx1zTC2mto7BxtesbS9F3vW3xggpuNIMjXeLD63mK0LSJ2Vh
NZ0YihpJ/eVCO439mDM7vtP1JQ4KveRMmAEIql1lIm5/SiBYqiA5JP0XMwIDAQAB
AoGAK47nBmswT3KKLWkG/o6lc5T5eugl8itDJ4A9KzSEnBSRYDhjXD1folnP6AkA
zzInZbrpjfgRcctT8JwGtdVYFpJFJOO5/LoWS3SHHLiHtwBXmEBQowvkIky9iGB5
VGUnaCMFB8ddi4Y9CAu5wahxEA6rGUb0mHqsPQ3tBwFhkDECQQD3W+lNQp0K2/TZ
Tkl713IbzJ6+6JLGzxPlGln080wlyZ/HEJKWqF3ro/J85P59A5I3c4ZDWKQGp1ZG
eNVhYgN7AkEAksIxWIYP3Tdfji6OTUrn/DN3/ZEfggEzUQIPUWVd9i5oSkKICZ7h
u/UCJ8UVSOAhsgmMcOjSNLMQhzVvqWbxqQJBAKbfBoDsk20j/gYrXj+BlKVUYTOB
SqN8R3ujT1SEXbaQUo3EjF++rb2uGIRRJ63Gnvlxof4E6oLimL1p/ul3ackCQFyl
xXsqHwe7dlKPJ3y6Bhvb7isgm7B5y4ifcUYkZR4OC/6dY74XFFCRCwxKSfaYsAzy
JDv/bvyf8pY48MYT3AkCQQDG8ca9DtckMcP3wXk62LZrGZdCerkU7KgSo/ksObzx
W4majkDXHE/rXWrzIJkp7aSo1OBpEZU2K6C6htpA0a/3
-----END RSA PRIVATE KEY-----
`)

// PKCS#8
var publicKeyPKCS8 = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDwCzwSw/IvMI1PcVJo9TKvTh9
HAF97mNhQAUB4OxVC0NEkjtwsLmf3u/gzs+cFKi5o+u6mczdfaati05jvKDfsYNo
Q7C+4wO1uuvPHJUOyJU42yHiKjBuXBVgaNKo+/QfAvl1kmBFTt64A9sC3s5NBoyN
2IC3hnVa1AOpNqlW2QIDAQAB
-----END PUBLIC KEY-----
`)

var privateKeyPKCS8 = []byte(`
-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAMPALPBLD8i8wjU9
xUmj1Mq9OH0cAX3uY2FABQHg7FULQ0SSO3CwuZ/e7+DOz5wUqLmj67qZzN19pq2L
TmO8oN+xg2hDsL7jA7W6688clQ7IlTjbIeIqMG5cFWBo0qj79B8C+XWSYEVO3rgD
2wLezk0GjI3YgLeGdVrUA6k2qVbZAgMBAAECgYEAvv+iWYxECG/1ZxGwkJvkozVi
CuDqq7+RBHD88cpPjuOAbUXp7ZjiZhWXJVllxTt7Lje9aMNs26kgmzDT+gkxRZ34
egj12yIH8OAtdClDquAR2vRGPLeNWSpO7Im+5DXrKBqMPH0n01cdK4+uEukyDgCA
+vbsx7pDipIGp9AHC2kCQQDj9voUHjhbeMnV+2tsCU8hDikwAJpDPAZ1V+0RT07x
dOHbz09nX8BoSgqL7TbodRSVPrzP0AvmWsA7wAoX/AeLAkEA29MAftsO7tl66RXP
xyhcODTsfokd09eifPZiGcJfgaHP3KkRarn+YBz/eGMBWKTCaatq6ommoAC06isY
9sKHqwJAO/CksMWBbAvGhk0lYbLQ65AdpFGEPkl6KUCFRRflWfexq2pHJpc2sDVH
sKMe3OBsGRH1825wspEKGqvT+5p5IQJBAMOmh3hgzGe11XmDWk0uFPZJ1HvC2nNk
J1EFkcbPg2XDeVgyejf9lvRAmvixVc9pxUd7tEtPfKhIOL164lsuRMUCQQCalVyg
SXqt1Q+/VDq0pT6yGeyD0FlHaVizfqQXbzIuL1a53pFHCvwESIf+pDQtf52vTKov
Q+SJkhpIKBtq0qMi
-----END PRIVATE KEY-----
`)

func Test_rsa(t *testing.T) {
	t.Run("rsa encrypt decrypt", func(t *testing.T) {
		want := rsaRawData
		cypherData, err := RsaEncrypt(publicKey, rsaRawData)
		if err != nil {
			t.Fatal(err)
		}
		got, err := RsaDecrypt(privateKey, cypherData)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%x]", rsaRawData, cypherData)
	})

	t.Run("rsa sign verify", func(t *testing.T) {
		signData, err := RsaSign(privateKey, rsaRawData)
		if err != nil {
			t.Fatal(err)
		}
		err = RsaVerify(publicKey, rsaRawData, signData)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("[%s]  <=>  [%x]", rsaRawData, signData)
	})

	t.Run("rsa sign verify pkcs8", func(t *testing.T) {
		signData, err := RsaSign(privateKeyPKCS8, rsaRawData, WithRsaFormatPKCS8())
		if err != nil {
			t.Fatal(err)
		}
		err = RsaVerify(publicKeyPKCS8, rsaRawData, signData, WithRsaFormatPKCS8())
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("[%s]  <=>  [%x]", rsaRawData, signData)
	})
}

func TestRsaHex(t *testing.T) {
	t.Run("default rsa encrypt hex", func(t *testing.T) {
		want := string(rsaRawData)
		cypherStr, err := RsaEncryptHex(publicKey, rsaRawData)
		if err != nil {
			t.Fatal(err)
		}
		got, err := RsaDecryptHex(privateKey, cypherStr)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got [%s], want [%s]", got, want)
		}
		t.Logf("[%s]  <=>  [%s]", rsaRawData, cypherStr)
	})

}

func TestRsaBase64(t *testing.T) {
	t.Run("default rsa encrypt hex", func(t *testing.T) {
		cypherStr, err := RsaSignBase64(privateKey, rsaRawData)
		if err != nil {
			t.Fatal(err)
		}
		err = RsaVerifyBase64(publicKey, rsaRawData, cypherStr)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("[%s]  <=>  [%s]", rsaRawData, cypherStr)
	})
}

func BenchmarkRsa(b *testing.B) {
	b.Run("rsa encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RsaEncrypt(publicKey, rsaRawData)
		}
	})
	b.Run("rsa decrypt", func(b *testing.B) {
		cypherData, err := RsaEncrypt(publicKey, rsaRawData)
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, cypherData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(cypherData, tmp)
			RsaDecrypt(privateKey, cypherData)
		}
	})

	b.Run("rsa sign", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RsaSign(privateKey, rsaRawData)
		}
	})
	b.Run("rsa verify", func(b *testing.B) {
		signData, err := RsaSign(privateKey, rsaRawData)
		if err != nil {
			b.Fatal(err)
		}
		var tmp []byte
		copy(tmp, signData)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			copy(signData, tmp)
			RsaVerify(publicKey, rsaRawData, signData)
		}
	})
}

func TestRsaOption(t *testing.T) {
	o := defaultRsaOptions()
	var opts []RsaOption
	opts = append(opts,
		WithRsaFormatPKCS1(),
		WithRsaFormatPKCS8(),
		WithRsaHashTypeMd5(),
		WithRsaHashTypeSha1(),
		WithRsaHashTypeSha256(),
		WithRsaHashTypeSha512(),
		WithRsaHashType(crypto.SHA256),
	)
	o.apply(opts...)
}
