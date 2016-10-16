package main

import (
	zip2 "github.com/yeka/zip"
	"bytes"
	"os"
	"io"
	"fmt"
	"archive/zip"
	"log"
	"hash/crc32"
	"encoding/binary"
	"io/ioutil"
)

func main() {
	//fmt.Printf("%X\n", crc32.Update(0, crc32.IEEETable, []byte{0}))
	//return

	password := "zipcrypto"
	new(ZipReader).FromBytes(createzip(password), password)
	new(ZipReader).Open("java2.zip", password) // FA176C7F

	os.Remove("myzip/test1.zip")
	ioutil.WriteFile("myzip/test1.zip", createzip(password), os.ModePerm)
	return
	//z.Open("java2.zip", "zipcrypto") // FA176C7F
	//z.Open("java.zip", "123")

	//z.Open("a.zip")
	//z.Open("test.zip")
	fmt.Println("=======================================")
	return
	//unknown1()
	//show("java.zip", "zipcrypto")
	zipunzip()
	return
	fmt.Printf("%X\n", 878082192)
	show("banyak.zip", "123")
	//show("pass.zip", "zipcrypto")
	//show("java.zip", "zipcrypto")
	//show("a.zip", "aa") // Incorrect password
	//show("test.zip", "")

	fmt.Printf("%b\n", zip2.Crc32update(0, 0))
	h := crc32.NewIEEE()
	h.Write([]byte{0})
	fmt.Printf("%b\n", h.Sum32())
	fmt.Printf("%b\n", crc32.Update(0, crc32.IEEETable, []byte{0}))

}

func play1() {
	de := []byte("hello")
	//e := zip.NewZipCrypto([]byte("Passwordnya panjang banget kk"))
	//en := e.EncryptMessage([]byte("Hello Worlds Muahaha"))
	//d := zip.NewZipCrypto([]byte("Passwordnya panjang banget kk"))
	//de = d.DecryptMessage(en)
	fmt.Println(string(de))

	os.Exit(0)
	buf := new(bytes.Buffer)
	//a := zip.NewWriter(buf)
	//w, _ := a.Create("a.txt")
	//w.Write([]byte(`Hello World`))
	//a.Close()

	f, _ := os.Create("temp/test.zip")
	io.Copy(f, buf)
	f.Close()

	_ = zip.Deflate
}

func show(filename string, password string) {
	r, err := zip2.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Printf("------ %s ------\n", f.Name)

		f.SetDecryptionPassword(password)
		rc, err := f.Open()
		if err != nil {
			fmt.Println("Opening error")
			log.Print(err)
			continue
		}
		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, rc); err != nil {
			fmt.Println("buffering error")
			log.Print(err)
			continue
		}

		fmt.Printf("Size: %d\nContent:\n", buf.Len())
		fmt.Println(string(buf.Bytes()))
		//fmt.Println(buf.Bytes())
		rc.Close()
		fmt.Println()
		fmt.Println()
	}
}

func createzip(password string) []byte {
	fmt.Printf("\n\n========================== Writing ========================== \n\n")
	mydata, _ := ioutil.ReadFile("a.java")
	b := new(bytes.Buffer)
	z := zip2.NewWriter(b)
	//w, _ := z.Create("a.java")
	w, _ := z.Encrypt("a.java", password)
	fmt.Println(w.Write(mydata))
	z.Close()
	return b.Bytes()
}

func zipunzip() {
	fmt.Printf("\n\n========================== Writing ========================== \n\n")
	mydata := []byte(`Hello World`)
	password := "golang"
	b := new(bytes.Buffer)
	z := zip2.NewWriter(b)
	//w, _ := z.Create("a.txt")
	w, _ := z.Encrypt("a.txt", password)
	fmt.Println(w.Write(mydata))
	z.Close()

	fmt.Printf("\n\n========================== Reading ========================== \n\n")
	rz, _ := zip2.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
	rz.File[0].SetDecryptionPassword(password)
	r, _ := rz.File[0].Open()
	myresult := new(bytes.Buffer)
	io.Copy(myresult, r)
	r.Close()

	fmt.Println(mydata)
	fmt.Println(myresult.Bytes())
	fmt.Println(0 == bytes.Compare(mydata, myresult.Bytes()))
}

func unknown1() {
	hello := []byte("Hello world\n")
	b := new(bytes.Buffer)

	l := &zipwriterboongan{b}
	c := &zipwriterboongan{l}
	//m := &minues{l}
	//c := &minues2{m}
	//
	////i := io.MultiWriter(b, m, l, c)
	//c.Write(hello)
	c.Write(hello)
	fmt.Println(string(b.Bytes()))
}

type zipwriterboongan struct {
	w io.Writer
}

func (m *zipwriterboongan) Write(p []byte) (i int, err error) {
	temp := make([]byte, len(p))
	for i, b := range p {
		temp[i] = b ^ 0xFF
	}
	m.w.Write(temp)
	return
}

type ftpwriterboongan struct {
	w io.Writer
}

func (m *ftpwriterboongan) Write(p []byte) (i int, err error) {
	temp := make([]byte, len(p))
	for i, b := range p {
		temp[i] = b ^ 0xFF
	}
	m.w.Write(temp)
	return
}

type ZipReader struct {
	b           *Bytes
	password    []byte
	files       []*ZipFile
	central_dir []*ZipCentralDirectory
	dir_end     *ZipDirectoryEnd
}

type ZipFile struct {
	zip        *ZipReader
	header     *ZipFileHeader
	data       *ZipFileData
	descriptor *ZipFileDescriptor
}

func (zf *ZipFile) init(offset, datasize uint64) {
	h := &ZipFileHeader{file: zf}
	h.init(offset)

	f := &ZipFileData{file: zf}
	f.content = zf.zip.b.IGet(datasize)

	e := &ZipFileDescriptor{file: zf}
	if u16x(h.general_purpose_bit_flag) & 8 > 0 {
		e.init(zf.zip.b.Pos())
	}

	zf.header = h
	zf.data = f
	zf.descriptor = e
}

type ZipFileHeader struct {
	file *ZipFile
	local_file_header_signature []byte
	version_needed_to_extract   []byte
	general_purpose_bit_flag    []byte
	compression_method          []byte
	last_mod_file_time          []byte
	last_mod_file_date          []byte
	crc_32                      []byte
	compressed_size             []byte
	uncompressed_size           []byte
	filename_length             []byte
	extra_field_length          []byte
	filename                    []byte
	extra_field                 []byte
}

func (fh *ZipFileHeader) init(offset uint64)  {
	b := fh.file.zip.b
	b.Offset(offset)
	fh.local_file_header_signature = b.IGet(4)
	fh.version_needed_to_extract   = b.IGet(2)
	fh.general_purpose_bit_flag    = b.IGet(2)
	fh.compression_method          = b.IGet(2)
	fh.last_mod_file_time          = b.IGet(2)
	fh.last_mod_file_date          = b.IGet(2)
	fh.crc_32                      = b.IGet(4)
	fh.compressed_size             = b.IGet(4)
	fh.uncompressed_size           = b.IGet(4)
	fh.filename_length             = b.IGet(2)
	fh.extra_field_length          = b.IGet(2)
	fh.filename                    = b.IGet(u16x(fh.filename_length))
	fh.extra_field                 = b.IGet(u16x(fh.extra_field_length))
}

type ZipFileData struct {
	file *ZipFile
	content []byte
}

type ZipFileDescriptor struct {
	file *ZipFile
	signature         []byte
	crc_32            []byte
	compressed_size   []byte
	uncompressed_size []byte
}

func (fd *ZipFileDescriptor) init(offset uint64)  {
	b := fd.file.zip.b
	b.Offset(offset)
	sig := b.IGet(4)
	if u32x(sig) == 0x08074b50 {
		fd.signature = sig
		fd.crc_32 = b.IGet(4)
	} else {
		fd.crc_32 = sig
	}
	fd.compressed_size = b.IGet(4)
	fd.uncompressed_size = b.IGet(4)
}

type ZipCentralDirectory struct {
	zip *ZipReader
	central_file_header_signature []byte
	version_made_by               []byte
	version_needed_to_extrac      []byte
	general_purpose_bit_flag      []byte
	compression_method            []byte
	last_mod_file_time            []byte
	last_mod_file_date            []byte
	crc_32                        []byte
	compressed_size               []byte
	uncompressed_size             []byte
	filename_length               []byte
	extra_field_length            []byte
	file_comment_length           []byte
	disk_num_start                []byte
	internal_file_attr            []byte
	external_file_attr            []byte
	local_header_relative_offset  []byte
	filename                      []byte
	extra_field                   []byte
	file_comment                  []byte
}

func (cd *ZipCentralDirectory) init(offset uint64) (size uint64) {
	cd.zip.b.Offset(offset)
	cd.central_file_header_signature = cd.zip.b.IGet(4)
	cd.version_made_by               = cd.zip.b.IGet(2)
	cd.version_needed_to_extrac      = cd.zip.b.IGet(2)
	cd.general_purpose_bit_flag      = cd.zip.b.IGet(2)
	cd.compression_method            = cd.zip.b.IGet(2)
	cd.last_mod_file_time            = cd.zip.b.IGet(2)
	cd.last_mod_file_date            = cd.zip.b.IGet(2)
	cd.crc_32                        = cd.zip.b.IGet(4)
	cd.compressed_size               = cd.zip.b.IGet(4)
	cd.uncompressed_size             = cd.zip.b.IGet(4)
	cd.filename_length               = cd.zip.b.IGet(2)
	cd.extra_field_length            = cd.zip.b.IGet(2)
	cd.file_comment_length           = cd.zip.b.IGet(2)
	cd.disk_num_start                = cd.zip.b.IGet(2)
	cd.internal_file_attr            = cd.zip.b.IGet(2)
	cd.external_file_attr            = cd.zip.b.IGet(4)
	cd.local_header_relative_offset  = cd.zip.b.IGet(4)
	cd.filename                      = cd.zip.b.IGet(u16x(cd.filename_length))
	cd.extra_field                   = cd.zip.b.IGet(u16x(cd.extra_field_length))
	cd.file_comment                  = cd.zip.b.IGet(u16x(cd.file_comment_length))
	size = 46 + u16x(cd.filename_length) + u16x(cd.extra_field_length) + u16x(cd.file_comment_length)
	return
}

type ZipDirectoryEnd struct {
	zip *ZipReader
	signature                          []byte
	disknum                            []byte
	disknum_start_of_central_directory []byte
	entries_in_this_disk               []byte
	num_of_entries                     []byte
	size_of_central_dir                []byte
	offset_of_central_dir              []byte
	zipfile_comment_length             []byte
	zipfile_comment                    []byte
}

func (zde *ZipDirectoryEnd) init(offset uint64) {
	zde.zip.b.Offset(offset)
	zde.signature                          = zde.zip.b.IGet(4)
	zde.disknum                            = zde.zip.b.IGet(2)
	zde.disknum_start_of_central_directory = zde.zip.b.IGet(2)
	zde.entries_in_this_disk               = zde.zip.b.IGet(2)
	zde.num_of_entries                     = zde.zip.b.IGet(2)
	zde.size_of_central_dir                = zde.zip.b.IGet(4)
	zde.offset_of_central_dir              = zde.zip.b.IGet(4)
	zde.zipfile_comment_length             = zde.zip.b.IGet(2)
	zde.zipfile_comment                    = zde.zip.b.IGet(u16x(zde.zipfile_comment_length))
}

func (z *ZipReader) FromBytes(buf []byte, password string) {
	z.b = &Bytes{buf, 0}
	z.password = []byte(password)
	z.Init()
}

func (z *ZipReader) Open(filename string, password string) {
	buf, _ := ioutil.ReadFile(filename)
	z.FromBytes(buf, password)
}

func (z *ZipReader) Init() {
	// First find Directory End
	var i uint64
	for i = uint64(len(z.b.b)) - 22; i >= 0; i-- {
		if z.b.BGet(i, 4).I64() == 0x06054b50 {
			z.dir_end = &ZipDirectoryEnd{zip: z}
			z.dir_end.init(i)
			break
		}
	}

	// Read all Central Directory Header
	central_dir_offset := u16x(z.dir_end.offset_of_central_dir)
	for i = 0; i < u16x(z.dir_end.num_of_entries); i++ {
		cdfh := &ZipCentralDirectory{zip: z}
		central_dir_offset += cdfh.init(central_dir_offset)
		z.central_dir = append(z.central_dir, cdfh)
	}

	for i = 0; i < u16x(z.dir_end.num_of_entries); i++ {
		zf := &ZipFile{zip: z}
		zf.init(u32x(z.central_dir[i].local_header_relative_offset), u32x(z.central_dir[i].compressed_size))
		z.files = append(z.files, zf)
	}

	fmt.Printf("%0 X\n", z.files[0].header)
	fmt.Printf("%0 X\n", z.files[0].descriptor)
	// PlayGround
	//fmt.Printf("%v\n", z.files[0].data.content)
	//zc := zip2.NewZipCrypto(z.password)
	//buf := zc.DecryptMessage(z.files[0].data.content)
	////fmt.Printf("%v\n", buf)
	//rc := flate.NewReader(bytes.NewBuffer(buf[12:]))
	//defer rc.Close()
	//res, err := ioutil.ReadAll(rc)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Printf("%v\n", string(res))

	//fmt.Printf("---------------- %v ----------------\n", string(h.filename))
		//fmt.Printf("Signature: %02X\n", h.local_file_header_signature)
		//fmt.Printf("Version: %v\n", h.version_needed_to_extract)
		//fmt.Printf("General Purpose Bit Flag: %08b\n", h.general_purpose_bit_flag)
		//fmt.Printf("Compression Method: %v\n", h.compression_method)
		//fmt.Printf("CRC32: %0 X, Compressed: %v, Uncompressed: %v\n", h.crc_32, u16x(h.compressed_size), u16x(h.uncompressed_size))
		//fmt.Printf("FileName Length: %v\n", u16x(h.filename_length))
		//fmt.Printf("ExtraField Length: %v\n", u16x(h.extra_field_length))
		//fmt.Printf("ExtraField: %0 X\n", h.extra_field)
		//if u16x(h.general_purpose_bit_flag) & 8 > 0 {
		//	crc32 := z.b.IGet(4)
		//	if u32x(crc32) == 0x08074b50 {
		//		// Signature may or may not be there
		//		crc32 = z.b.IGet(4)
		//	}
		//	csize := z.b.IGet(4)
		//	usize := z.b.IGet(4)
		//	fmt.Printf("DataDescriptor:\nCRC32: %0 X, Compressed: %v, Uncompressed: %v, Next Signature: %X\n", crc32, u16x(csize), u16x(usize), z.b.IGet(4))
		//}
		//
		//crcbuf := make([]byte, 4)
		//crcbuf[3] = h.crc_32[3]
		//crcbuf[2] = (h.crc_32[3] >> 8)
		//crcbuf[1] = (h.crc_32[3] >> 16)
		//crcbuf[0] = (h.crc_32[3] >> 24)
		//
		////fmt.Printf("CRC BUF: %x %v\n", crcbuf, crcbuf)
		//
		////fmt.Printf("% X\n", file_data)
		//if u16x(h.general_purpose_bit_flag) & 1 > 0 {
		//	fmt.Println("Encrypted")
		//	zc := zip2.NewZipCrypto(z.password)
		//
		//	//result := file_data[0]
		//	//for i := 0; i < 12; i++ {
		//	//	if i == 11 {
		//	//		//fmt.Printf("HASH MAGIC: %X %v\n", result ^ zc.MagicByte(), result ^ zc.MagicByte())
		//	//	}
		//	//	zc.UpdateKeys(result ^ zc.MagicByte())
		//	//	if i < 12 {
		//	//		result = file_data[i + 1]
		//	//	}
		//	//}
		//	file_header := zc.DecryptMessage(file_data[0:12])
		//	//fmt.Printf("Magic Byte: %X\n", zc.MagicByte())
		//	file_data = zc.DecryptMessage(file_data[12:])
		//	fmt.Printf("%0 X : %0 X\n", file_header, file_data)
		//	//fmt.Printf("%v %v\n", string(file_data[0:12]), string(file_data[12:]))
		//	//file_data = file_data[12:]
		//}
		//if u16x(h.compression_method) == 8 {
		//	//fmt.Printf("Compressed: %v\n", len(file_data))
		//	rc := flate.NewReader(bytes.NewBuffer(file_data))
		//	defer rc.Close()
		//	res, err := ioutil.ReadAll(rc)
		//	file_data = res
		//	if err != nil {
		//		fmt.Println(err.Error())
		//	}
		//}
		//crc := crc32.NewIEEE()
		//crc.Write(file_data)
		//fmt.Printf("GO CRC: %X\n", crc.Sum32())
		//fmt.Printf("\n%v\n", string(file_data))
	//}
}

func (z *ZipReader) DirectoryStructure() {

}

func u16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func u32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func u16x(b []byte) uint64 {
	return uint64(u16(b))
}

func u32x(b []byte) uint64 {
	return uint64(u32(b))
}

func u64x(b []byte) uint64 {
	return uint64(u32(b))
}


type Bytes struct {
	b []byte
	i uint64
}

// Get slice of byte
func (b *Bytes) Get(offset, size uint64) []byte {
	return b.b[offset : offset + size]
}

func (b *Bytes) BGet(offset, size uint64) *Bytes {
	return &Bytes{b.Get(offset, size), 0}
}

// Incremental Get
func (b *Bytes) IGet(size uint64) []byte {
	i := b.i
	b.i = i + size
	return b.b[i : b.i]
}

// Incremental Get
func (b *Bytes) IBGet(size uint64) *Bytes {
	return &Bytes{b.IGet(size), 0}
}

func (b *Bytes) I64() uint64 {
	switch len(b.b) {
	case 1:
		return uint64(b.b[0])
	case 2:
		return u16x(b.b)
	case 4:
		return u32x(b.b)
	case 8:
		return u64x(b.b)
	}
	return 0
}
// Set offset
func (b *Bytes) Offset(i uint64) {
	b.i = i
}

func (b *Bytes) Pos() uint64 {
	return b.i
}