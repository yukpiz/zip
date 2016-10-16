package zip

import (
	"io"
	"bytes"
	"fmt"
	"hash/crc32"
)

type ZipCrypto struct {
	password []byte
	Keys [3]uint32
}

func NewZipCrypto(passphrase []byte) *ZipCrypto {
	z := &ZipCrypto{}
	z.password = passphrase
	z.Init()
	return z
}

func (z *ZipCrypto) Init() {
	z.Keys[0] = 0x12345678
	z.Keys[1] = 0x23456789
	z.Keys[2] = 0x34567890

	for i := 0; i < len(z.password); i++ {
		z.UpdateKeys(z.password[i])
	}
}

func (z *ZipCrypto) UpdateKeys(byteValue byte) {
	z.Keys[0] = Crc32update(z.Keys[0], byteValue);
	z.Keys[1] += z.Keys[0] & 0xff;
	z.Keys[1] = z.Keys[1] * 134775813 + 1;
	z.Keys[2] = Crc32update(z.Keys[2], (byte) (z.Keys[1] >> 24));
}

func (z *ZipCrypto) MagicByte() byte {
	var t uint32 = z.Keys[2] | 2
	return byte((t * (t ^ 1)) >> 8)
}

func (z *ZipCrypto) EncryptMessage(plaintext []byte) []byte {
	length := len(plaintext)
	CipherText := make([]byte, length)
	for i := 0; i < length; i++ {
		C := plaintext[i]
		CipherText[i] = plaintext[i] ^ z.MagicByte()
		z.UpdateKeys(C)
	}
	return CipherText
}

func (z *ZipCrypto) DecryptMessage(chipertext []byte) []byte {
	length := len(chipertext)
	PlainText := make([]byte, length)
	for i, c := range chipertext {
		v := c ^ z.MagicByte();
		z.UpdateKeys(v)
		PlainText[i] = v
	}
	return PlainText
}

func Crc32update(pCrc32 uint32, bval byte) uint32 {
	return crc32.IEEETable[(pCrc32 ^ uint32(bval)) & 0xff] ^ (pCrc32 >> 8)
}

func ZipCryptoDecryptor(r *io.SectionReader, password []byte) (*io.SectionReader, error) {
	//return r, nil
	z := NewZipCrypto(password)
	b := make([]byte, r.Size())

	r.Read(b)

	m := z.DecryptMessage(b)
	fmt.Printf("Header: %X %X %X\n%v\n", m[0:4], m[4:8], m[8:12], m[0:12])
	fmt.Printf("Header: %X\n", m[11])
	fmt.Printf("Content: %X\n", m[12:])
	fmt.Printf("Content: %v\n", string(m[12:]))
	//if m[11] > 0xFF {
	//	return nil, errors.New("incorrect password")
	//}
	//fmt.Printf("Size: %v\n", len(m))
	return io.NewSectionReader(bytes.NewReader(m), 12, int64(len(m))), nil
}

type zipCryptoWriter struct {
	w     io.Writer
	z     *ZipCrypto
	first bool
	fw    *fileWriter
}

func (z *zipCryptoWriter) Write(p []byte) (n int, err error) {
	if z.first {
		z.first = false
		header := []byte{0xF8, 0x53, 0xCF, 0x05, 0x2D, 0xDD, 0xAD, 0xC8, 0x66, 0x3F, 0x8C, 0xAC}
		header = z.z.EncryptMessage(header)

		crc := z.fw.ModifiedTime
		header[10] = byte(crc)
		header[11] = byte(crc >> 8)

		z.z.Init()
		z.w.Write(z.z.EncryptMessage(header))
		n += 12
		//z.z.Init()
	}
	z.w.Write(z.z.EncryptMessage(p))
	z.fw.FileHeader.CompressedSize += uint32(n)
	return
}

func ZipCryptoEncryptor(i io.Writer, pass passwordFn, fw *fileWriter) (io.Writer, error)  {
	fmt.Printf("Initializationg...\n")
	z := NewZipCrypto(pass())
	zc := &zipCryptoWriter{i, z, true, fw}
	return zc, nil
}