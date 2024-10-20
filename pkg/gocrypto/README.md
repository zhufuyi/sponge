## gocrypto

Commonly used `one-way encryption`, `symmetric encryption and decryption`, `asymmetric encryption and decryption` libraries, including hash, aes, des, rsa.

<br>

## Example of use

### Hash one-way encryption

```go
    import "github.com/zhufuyi/sponge/pkg/gocrypto"

    var hashRawData = []byte("hash_abcdefghijklmnopqrstuvwxyz0123456789")

    // independent hash functions
    gocrypto.Md5(hashRawData)
    gocrypto.Sha1(hashRawData)
    gocrypto.Sha256(hashRawData)
    gocrypto.Sha512(hashRawData)

    // hash collection, specify the execution of the corresponding hash function
    // according to the hash type
    gocrypto.Hash(crypto.MD5, hashRawData)
    gocrypto.Hash(crypto.SHA3_224, hashRawData)
    gocrypto.Hash(crypto.SHA256, hashRawData)
    gocrypto.Hash(crypto.SHA3_224, hashRawData)
    gocrypto.Hash(crypto.BLAKE2s_256, hashRawData)
```

<br>

### Password hash and checksum with salt

The password registered by the user is stored in the database through hash, and the password registered is compared with the hash value to judge whether the password is correct all the time, so as to ensure that only the user knows the plaintext of the password.

```go
    import "github.com/zhufuyi/sponge/pkg/gocrypto"

    pwd := "123"

    // hash
    hashStr, err := gocrypto.HashAndSaltPassword(pwd)
    if err != nil {
        return err
    }

    // check password
    ok := gocrypto.VerifyPassword(pwd, hashStr)
    if !ok {
        return errors.New("passwords mismatch")
    }
```

<br>

### AES encrypt and decrypt

AES (`Advanced Encryption Standard`) Advanced Encryption Standard, designed to replace `DES`, has four packet encryption modes: ECB CBC CFB CTR.

There are four functions `AesEncrypt`, `AesDecrypt`, `AesEncryptHex`, `AesDecryptHex`.

```go
    import "github.com/zhufuyi/sponge/pkg/gocrypto"

    var (
        aesRawData = []byte("aes_abcdefghijklmnopqrstuvwxyz0123456789")
        aesKey     = []byte("aesKey0123456789aesKey0123456789")
    )

    // AesEncrypt and AesDecrypt have default values for their arguments:
    // default mode is ECB, can be modified to CBC CTR CFB 
    // default key length is 16, which can be modified to 24 32

    // default mode is ECB, default key length is 16
    cypherData, _ := gocrypto.AesEncrypt(aesRawData) // encrypt
    raw, _ := gocrypto.AesDecrypt(cypherData) // decrypt, return to original

    // mode is ECB, key length is 32
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, gocrypto.WithAesKey(aesKey))  // encrypt
    raw, _ := gocrypto.AesDecrypt(cypherData, gocrypto.WithAesKey(aesKey)) // decrypt

    // mode is CTR, default key length is 16
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, gocrypto.WithAesModeCTR())  // encrypt
    raw, _ := gocrypto.AesDecrypt(cypherData, gocrypto.WithAesModeCTR())  // decrypt

    // mode is CBC, key length is 32
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, gocrypto.WithAesModeECB(), gocrypto.WithAesKey(aesKey)) // encrypt
    raw, _ := gocrypto.AesDecrypt(cypherData, gocrypto.WithAesModeECB(), gocrypto.WithAesKey(aesKey))   // decrypt


    // AesEncryptHex and AesDecryptHex functions, the ciphertext of these two functions is transcoded by hex,
    // and used in exactly the same way as AesEncrypt and AesDecrypt.
```
<br>

### DES encrypt and decrypt

DES (`Data Encryption Standard`) data encryption standard, is currently one of the most popular encryption algorithms, there are four packet encryption mode: ECB CBC CFB CTR.

There are four functions `DesEncrypt`, `DesDecrypt`, `DesEncryptHex`, `DesDecryptHex`.

```go
    import "github.com/zhufuyi/sponge/pkg/gocrypto"

    var (
        desRawData = []byte("des_abcdefghijklmnopqrstuvwxyz0123456789")
        desKey     = []byte("desKey0123456789desKey0123456789")
    )
// PKCS#1
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
xxxxxx
-----END PUBLIC KEY-----
`)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
xxxxxx
-----END RSA PRIVATE KEY-----
`)

   // DesEncrypt and DesDecrypt have default values for their arguments:
   // default mode is ECB, can be modified to CBC CTR CFB 
   // default key length is 16, which can be modified to 24 32

    // default mode is ECB, default key length is 16
    cypherData, _ := gocrypto.DesEncrypt(desRawData) // encrypt
    raw, _ := gocrypto.DesDecrypt(cypherData) // decrypt

    // mode is ECB, key length is 32
    cypherData, _ := gocrypto.DesEncrypt(desRawData, gocrypto.WithDesKey(desKey)) // encrypt
    raw, _ := gocrypto.DesDecrypt(cypherData, gocrypto.WithDesKey(desKey)) // decrypt

    // mode is CTR, default key length is 16
    cypherData, _ := gocrypto.DesEncrypt(desRawData, gocrypto.WithDesModeCTR()) // encrypt
    raw, _ := gocrypto.DesDecrypt(cypherData, gocrypto.WithDesModeCTR()) // decrypt

    // mode is CBC, key length is 32
    cypherData, _ := gocrypto.DesEncrypt(desRawData, gocrypto.WithDesModeECB(), gocrypto.WithDesKey(desKey)) // encrypt
    raw, _ := gocrypto.DesDecrypt(cypherData, gocrypto.WithDesModeECB(), gocrypto.WithDesKey(desKey))        // decrypt


    // DesEncryptHex and DesDecryptHex functions, the ciphertext of these two functions is transcoded by hex,
    // and used in exactly the same way as DesEncrypt and DesDecrypt.
```

<br>

### RSA asymmetric encryption and decryption

#### RSA encryption and decryption

The public key is used for encryption, and the private key is used for decryption. For example, if someone uses the public key to encrypt information and send it to you, you have the private key to decrypt the information content.

There are four functions: `RsaEncrypt`, `RsaDecrypt`, `RsaEncryptHex`, `RsaDecryptHex`.

```go
    import "github.com/zhufuyi/sponge/pkg/gocrypto"

    var rsaRawData = []byte("rsa_abcdefghijklmnopqrstuvwxyz0123456789")
    // PKCS#1
    var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
xxxxxx
-----END PUBLIC KEY-----
`)

    var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
xxxxxx
-----END RSA PRIVATE KEY-----
`)
	
    // RsaEncrypt and RsaDecrypt have default values for their arguments:
    // default key pair format: PKCS#1, can be modified to PKCS#8

    // default key pair format is PKCS#1
    cypherData, _ := gocrypto.RsaEncrypt(publicKey, rsaRawData) // encrypt
    raw, _ := gocrypto.RsaDecrypt(privateKey, cypherData) // decrypt

    // key pair format is PKCS#8
    cypherData, _ := gocrypto.RsaEncrypt(publicKey, rsaRawData, gocrypto.WithRsaFormatPKCS8()) // encrypt
    raw, _ := gocrypto.RsaDecrypt(privateKey, cypherData, gocrypto.WithRsaFormatPKCS8()) // decrypt


    // RsaEncryptHex and RsaDecryptHex functions, the ciphertext of these two functions is transcoded by hex,
    // and used in exactly the same way as RsaEncrypt and RsaDecrypt.
```

<br>

#### RSA signature and signature verification

The private key is used to sign, and the public key is used to verify the signature. For example, you sign your identity with the private key, and others verify whether your identity can be trusted through the public key.

There are four functions: `RsaSign`, `RsaVerify`, `RsaSignBase64`, `RsaVerifyBase64`.

```go
   import "github.com/zhufuyi/sponge/pkg/gocrypto"

    var rsaRawData = []byte("rsa_abcdefghijklmnopqrstuvwxyz0123456789")

    // RsaEncrypt and RsaDecrypt have default values for their arguments:
    // default key pair format is PKCS#1, can be modified to PKCS#8
    // default hash is sha1, can be modified to sha256, sha512

    // default key pair format is PKCS#1, default hash is sha1
    signData, _ := gocrypto.RsaSign(privateKey, rsaRawData) // signature
    err := gocrypto.RsaVerify(publicKey, rsaRawData, signData) // signature verification

    // default key pair format is PKCS#1, hash is sha256
    signData, _ := gocrypto.RsaSign(privateKey, rsaRawData, gocrypto.WithRsaHashTypeSha256()) // signature
    err := gocrypto.RsaVerify(publicKey, rsaRawData, signData, gocrypto.WithRsaHashTypeSha256()) // signature verification

    // key pair format is PKCS#8, default hash is sha1
    signData, _ := gocrypto.RsaSign(privateKey, rsaRawData, gocrypto.WithRsaFormatPKCS8()) // signature
    err := gocrypto.RsaVerify(publicKey, rsaRawData, signData, gocrypto.WithRsaFormatPKCS8()) // signature verification

    // key pair format is PKCS#8, hash is sha512
    signData, _ := gocrypto.RsaSign(privateKey, rsaRawData, gocrypto.WithRsaFormatPKCS8(), gocrypto.WithRsaHashTypeSha512()) // signature
    err := gocrypto.RsaVerify(publicKey, rsaRawData, signData, gocrypto.WithRsaFormatPKCS8(), gocrypto.WithRsaHashTypeSha512()) // signature verification


    // The ciphertext of RsaSignBase64 and RsaVerifyBase64 is base64 transcoded
    // and used exactly the same as RsaSign and RsaVerify.
```
