package meta

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// VorbisComment contains a list of name-value pairs.
//
// ref: https://www.xiph.org/flac/format.html#metadata_block_vorbis_comment
type VorbisComment struct {
	// Vendor name.
	Vendor string
	// A list of tags, each represented by a name-value pair.
	Tags [][2]string
}

// parseVorbisComment reads and parses the body of an VorbisComment metadata
// block.
func (block *Block) parseVorbisComment() error {
	// 32 bits: vendor length.
	var x uint32
	err := binary.Read(block.lr, binary.LittleEndian, &x)
	if err != nil {
		return unexpected(err)
	}

	// (vendor length) bits: Vendor.
	buf, err := readBytes(block.lr, int(x))
	if err != nil {
		return unexpected(err)
	}
	comment := new(VorbisComment)
	block.Body = comment
	comment.Vendor = string(buf)

	// Parse tags.
	// 32 bits: number of tags.
	err = binary.Read(block.lr, binary.LittleEndian, &x)
	if err != nil {
		return unexpected(err)
	}
	if x < 1 {
		return nil
	}
	comment.Tags = make([][2]string, x)
	for i := range comment.Tags {
		// 32 bits: vector length
		err = binary.Read(block.lr, binary.LittleEndian, &x)
		if err != nil {
			return unexpected(err)
		}

		// (vector length): vector.
		buf, err = readBytes(block.lr, int(x))
		if err != nil {
			return unexpected(err)
		}
		vector := string(buf)

		// Parse tag, which has the following format:
		//    NAME=VALUE
		pos := strings.Index(vector, "=")
		if pos == -1 {
			return fmt.Errorf("meta.Block.parseVorbisComment: unable to locate '=' in vector %q", vector)
		}
		comment.Tags[i][0] = vector[:pos]
		comment.Tags[i][1] = vector[pos+1:]
	}

	return nil
}
