package rtcp

import (
	"fmt"
	"encoding/binary"
	"github.com/gdexlab/go-render/render" // for printing embedded structs
)

const (
	XRHeaderSize = 12 // bytes
	
	// The rtp version numbers in the extended report header are represented in two bits.
	// I am not certain exactly how, but this is my best guess:
	INT_0_IN_TWO_BITS = "00"
	INT_1_IN_TWO_BITS = "01"
	INT_2_IN_TWO_BITS = "10"
	INT_3_IN_TWO_BITS = "11"

	ADDITIONAL_PADDING_EXISTS = "1"
	NO_ADDITIONAL_PADDING = "0"

	NUMBER_OF_BLOCKS = 35 // not including header as a block and type 9 is skipped

	CHUNKTYPE_RUNLENGTH = iota + 1
	CHUNKTYPE_BITVECTOR
	CHUNKTYPE_TERMINATINGNULL
)

const (

	// https://www.iana.org/assignments/rtcp-xr-block-types/rtcp-xr-block-types.xhtml
	BLOCKTYPE_XRBT1 = iota + 1
	BLOCKTYPE_XRBT2
	BLOCKTYPE_XRBT3
	BLOCKTYPE_XRBT4
	BLOCKTYPE_XRBT5
	BLOCKTYPE_XRBT6
	BLOCKTYPE_XRBT7
	BLOCKTYPE_XRBT8
	BLOCKTYPE_XRBT9 // NOTE: skipping this one for lack of documentation
	BLOCKTYPE_XRBT10
	BLOCKTYPE_XRBT11
	BLOCKTYPE_XRBT12
	BLOCKTYPE_XRBT13
	BLOCKTYPE_XRBT14
	BLOCKTYPE_XRBT15
	BLOCKTYPE_XRBT16
	BLOCKTYPE_XRBT17
	BLOCKTYPE_XRBT18
	BLOCKTYPE_XRBT19
	BLOCKTYPE_XRBT20
	BLOCKTYPE_XRBT21
	BLOCKTYPE_XRBT22
	BLOCKTYPE_XRBT23
	BLOCKTYPE_XRBT24
	BLOCKTYPE_XRBT25
	BLOCKTYPE_XRBT26
	BLOCKTYPE_XRBT27
	BLOCKTYPE_XRBT28
	BLOCKTYPE_XRBT29
	BLOCKTYPE_XRBT30
	BLOCKTYPE_XRBT31
	BLOCKTYPE_XRBT32
	BLOCKTYPE_XRBT33
	BLOCKTYPE_XRBT34
	BLOCKTYPE_XRBT35
)

var (
	XRHeaderLength int 
	XRBT1Length int 
	XRBT2Length int 
	XRBT3Length int
	XRBT4Length int 
	XRBT5Length int 
	XRBT6Length int 
	XRBT7Length int 
	XRBT8Length int
	XRBT9Length int // NOTE: skipping this one for lack of documentation
	XRBT10Length int 
	XRBT11Length int 
	XRBT12Length int 
	XRBT13Length int 
	XRBT14Length int
	XRBT15Length int 
	XRBT16Length int 
	XRBT17Length int 
	XRBT18Length int 
	XRBT19Length int 
	XRBT20Length int 
	XRBT21Length int 
	XRBT22Length int
	XRBT23Length int
	XRBT24Length int 
	XRBT25Length int 
	XRBT26Length int 
	XRBT27Length int 
	XRBT28Length int 
	XRBT29Length int 
	XRBT30Length int 
	XRBT31Length int 
	XRBT32Length int 
	XRBT33Length int 
	XRBT34Length int 
	XRBT35Length int
)

type ExtendedReports struct {
	*XRHeader
	*XRBT1 
	*XRBT2 
	*XRBT3
	*XRBT4 
	*XRBT5 
	*XRBT6 
	*XRBT7 
	*XRBT8
	// type 9 is omitted due to lack of documentation
	*XRBT10 
	*XRBT11 
	*XRBT12 
	*XRBT13 
	*XRBT14
	*XRBT15 
	*XRBT16 
	*XRBT17 
	*XRBT18 
	*XRBT19 
	*XRBT20 
	*XRBT21 
	*XRBT22
	*XRBT23
	*XRBT24 
	*XRBT25 
	*XRBT26 
	*XRBT27 
	*XRBT28 
	*XRBT29 
	*XRBT30 
	*XRBT31 
	*XRBT32 
	*XRBT33 
	*XRBT34 
	*XRBT35 
}


// var Blocks map = new(map{int}struct)



// ---------------------------------------------------------------------------
// Code below is for Extended Report
// ---------------------------------------------------------------------------

// unmarshal gets information from a packet
// this function is called in packet.go on line 116
func (r ExtendedReports) Unmarshal(rawData []byte, XRBytesProcessed int) error  {
	XRBytesProcessed := XRHeaderSize

	*XRHeader.unmarshal()

	// bytesprocessed = int(h.Length+1) * 4
	// if bytesprocessed > len(rawData) {
	//, XRBytesProcessed int 	(XRBytesProcessedreturn nil, 0, errPacketTooShort
	// }
	// inPacket := rawData[:bytes, XRBytesProcessed intpr(XRBytesProcessedocessed]

		// if bytesTotal > len(rawData) {
		//
		// }

			
	


	for(XRBytesProcessed < *XRHeader.length) { // *XRHeader.length is the length of the whole ExtendedReports


		blockType, err := getBlockType(rawData[XRBytesProcessed:]) 


		switch blockType {
			case BLOCKTYPE_XRBT1:
				block := new(XRBT1)
			case BLOCKTYPE_XRBT2:
				block := new(XRBT2)
			case BLOCKTYPE_XRBT3:
				block := new(XRBT3)
			case BLOCKTYPE_XRBT4:
				block := new(XRBT4)
			case BLOCKTYPE_XRBT5:
				block := new(XRBT5)
			case BLOCKTYPE_XRBT6:
				block := new(XRBT6)
			case BLOCKTYPE_XRBT7:
				block := new(XRBT7)
			case BLOCKTYPE_XRBT8:
				block := new(XRBT8)
			case BLOCKTYPE_XRBT9:
				block := new(XRBT9)
			// type 9 is skipped over due to lack of documentation
			case BLOCKTYPE_XRBT10:
				block := new(XRBT10)
			case BLOCKTYPE_XRBT11:
				block := new(XRBT11)
			case BLOCKTYPE_XRBT12:
				block := new(XRBT12)
			case BLOCKTYPE_XRBT13:
				block := new(XRBT13)
			case BLOCKTYPE_XRBT14:
				block := new(XRBT14)
			case BLOCKTYPE_XRBT15:
				block := new(XRBT15)
			case BLOCKTYPE_XRBT16:
				block := new(XRBT16)
			case BLOCKTYPE_XRBT17:
				block := new(XRBT17)
			case BLOCKTYPE_XRBT18:
				block := new(XRBT18)
			case BLOCKTYPE_XRBT19:
				block := new(XRBT19)
			case BLOCKTYPE_XRBT20:
				block := new(XRBT20)
			case BLOCKTYPE_XRBT21:
				block := new(XRBT21)
			case BLOCKTYPE_XRBT22:
				block := new(XRBT22)
			case BLOCKTYPE_XRBT23:
				block := new(XRBT23)
			case BLOCKTYPE_XRBT24:
				block := new(XRBT24)
			case BLOCKTYPE_XRBT25:
				block := new(XRBT25)
			case BLOCKTYPE_XRBT26:
				block := new(XRBT26)
			case BLOCKTYPE_XRBT27:
				block := new(XRBT27)
			case BLOCKTYPE_XRBT28:
				block := new(XRBT28)
			case BLOCKTYPE_XRBT29:
				block := new(XRBT29)
			case BLOCKTYPE_XRBT30:
				block := new(XRBT30)
			case BLOCKTYPE_XRBT31:
				block := new(XRBT31)
			case BLOCKTYPE_XRBT32:
				block := new(XRBT32)
			case BLOCKTYPE_XRBT33:
				block := new(XRBT33)
			case BLOCKTYPE_XRBT34:
				block := new(XRBT34)
			case BLOCKTYPE_XRBT35:
				block := new(XRBT35)
			case BLOCKTYPE_XRBT36:
				block := new(XRBT36)
			default:
				err := fmt.Errorf("Extended reports block type is invalid. Cannot read Extended Reports.")
				return err
		}
	

		// NOTE: WHY IS PACKET RETURNED? just be sure you know

		// NOTE: block.unmarshal can return blocklength

		

		packet, XRBytesProcessed, err := block.unmarshal(rawPacket[XRBytesProcessed:], XRBytesProcessed)
		if err != nil {
			return err
		}
		// NOTE: might check the blow in every unmarshal function with checkIfPacketIsTooShort()
		if packet.len() < block.blockLength {
			err := errPacketTooShort
			return err
		}
		// output := render.AsCode(r)
		// fmt.Println(output)
   }

	return err
}



// marshal puts information into a packet
func (r ExtendedReports) Marshal() ([]byte, error) {

	return bytes, err
}

func (r *ExtendedReports) DestinationSSRC() []uint32 {
	out := make([]uint32, len(r.Reports)+1)
	for i, v := range r.Reports {
		out[i] = v.SSRC
	}
	out[len(r.Reports)] = r.SSRC
	return out
}

// ---------------------------------------------------------------------------
// Below are helper functions and types for the report blocks
// ---------------------------------------------------------------------------

var b binary.BigEndian

type chunk [2]byte
type word [4]byte


type runLengthChunk struct {
	
}

func (c runLengthChunk) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

func (c runLengthChunk) marshal() ([]byte, error) {

}

type bitVectorChunk struct {

}

func (c runLengthChunk) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

func (c runLengthChunk) marshal() ([]byte, error) {

}

type terminatingNullChunk struct {
	// NullZeros Chunk
}

func (c runLengthChunk) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

func (c runLengthChunk) marshal() ([]byte, error) {

}

// this function gets chunk type
func whichChunkTypeIsThis(chunk chunk) string, error {

	


	return chunkType, err
}

func getBlockType(rawData []byte, XRBytesProcessed int) (blockType uint8, error) {

	return blockType, err
}

func checkIfPacketIsTooShort(packet int, blockLength uint8) error {
	if packet.len() < block.blockLength {
		err := errPacketTooShort
		return err
	}
	return err
}

func resetXRBytesProcessed() {
	XRBytesProcessed = XRHeaderSize
}

//https://www.iana.org/assignments/rtcp-xr-block-types/rtcp-xr-block-types.xhtml


// XR Packet Format Header
// as found on https://tools.ietf.org/html/rfc3611
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |V=2|P|reserved |   PT=XR=207   |             length            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                              SSRC                             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                         report blocks                         :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+



// ---------------------------------------------------------------------------
// Code below is for XRHeader
// ---------------------------------------------------------------------------

type XRHeader {
	Version uint8
	Padding bool
	PacketType uint8
	Length uint8
	SRCC uint32 

// unmarshal gets information from a packet
func (r *XRHeader) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {
	rawPacket := make([]byte, headerLength)

	rawPacket[0] |= rtpVersion << versionShift

	if r.Padding {
		rawPacket[0] |= 1 << paddingShift
	}

	if r.Count > 31 {
		return nil, errInvalidHeader
	}
	rawPacket[0] |= r.Count << countShift

	rawPacket[1] = uint8(r.Type)

	binary.BigEndian.PutUint16(rawPacket[2:], r.Length)

	return out, 
}



/ marshal puts information into a packet
func (h *XRHeader) marshal() ([]byte, error) {

	return 
}

fuc (h *XRHeader) DestinationSSRC() []uint32 {
	out := make([]uint32, len(r.Reports)+1)
	for i, v := range r.Reports {
		out[i] = v.SSRC
	}
	out[len(r.Reports)] = r.SSRC
	return out
}

// ---------------------------------------------------------------------------
// Code below is for XRBT1
// ---------------------------------------------------------------------------


// Report Block Types' Formats:
// as found on https://tools.ietf.org/html/rfc3611

// BT=1  :  Loss RLE Report Block Type
// BT=2  :  Duplicate RLE Report Block Type
// BT=3  :  Packet Receipt Times Report Block Type
// BT=4  :  Receiver Reference Time Report Block Type
// BT=5  :  DLRR (Delay since Last Receiver Reference) Report Block Type
// BT=6  :  Statistics Summary Report Block
// BT=7  :  VoIP Metrics Report Block

// -----------------------------------------------------------------

// The Loss RLE Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.1
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=1      | rsvd. |   T   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk 1              |             chunk 2           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                              ...                              :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk n-1            |             chunk n           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+



type XRBT1 struct {
	BlockType uint8
	Thinning uint8 
	BlockLength uint16
	SSRC uint32 // not sure if unsigned or not
	BeginningSequence uint16
	EndingSequence uint16
	Chunks []chunk
}


// unmarshal gets information from a packet
func (r *XRBT1) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT1) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT2
// ---------------------------------------------------------------------------

// The Duplicate RLE Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.2

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=2      | rsvd. |   T   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk 1              |             chunk 2           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                              ...                              :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk n-1            |             chunk n           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT2 struct {
	BlockType uint8
	Thinning uint8 
	BlockLength uint16
	SSRC uint32 // not sure if unsigned or not
	BeginningSequence uint16
	EndingSequence uint16
	Chunks []chunk
}

// unmarshal gets information from a packet
func (r *XRBT2) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT2) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT2
// ---------------------------------------------------------------------------


// The Packet Receipt Times Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.3

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=3      | rsvd. |   T   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Receip packet begin_seq                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Receip packet (begin_seq + 1) mod 65536        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                              ...                              :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Receip packet (end_seq - 1) mod 65536          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT3 struct {
	BlockType uint8
	Thinning uint8 
	BlockLength uint16
	SSRC uint32
	Seq uint16
	ReceiptTime uint32
}

// unmarshal gets information from a packet
func (r *XRBT3) unmarshal() ( error) {

}

// marshal puts information into a packet
func (r *XRBT3) marshal(bytes []byte) error {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT4
// ---------------------------------------------------------------------------

// The Receiver Reference Time Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.4

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=4      |   reserved    |       block length = 2        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              NTP timestamp, most significant word             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             NTP timestamp, least significant word             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT4 struct {
	BlockType uint8
	Thinning uint8 
	BlockLength uint16
	NTPTimestamp uint32 // NOTE: unsure about type for this
}

// unmarshal gets information from a packet
func (r *XRBT4) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT4) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT5
// ---------------------------------------------------------------------------

// The DLRR Report Block has the following format
// found on https://tools.ietf.org/html/rfc3611#section-4.5

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=5      |   reserved    |         block length          |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |                 SSRC_1 first receiver)               | sub-
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ block
// |                         last RR (LRR)                         |   1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                   delay since last RR (DLRR)                  |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |                 SSRC_2 second receiver)              | sub-
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ block
// :                               ...                             :   2
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+

type XRBT5 struct {
	BlockType uint8
	BlockLength uint16
	SSRCs []uint32
	LRRs  []uint32
	DLRRs  []uint32
}

// unmarshal gets information from a packet
func (r *XRBT5) unmarshal() ( error) {

}

// marshal puts information into a packet
func (r *XRBT5) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT6
// ---------------------------------------------------------------------------

// The Statistics Summary Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.6

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=6      |L|D|J|ToH|rsvd.|       block length = 9        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        lost_packets                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        dup_packets                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         min_jitter                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         max_jitter                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         mean_jitter                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         dev_jitter                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | min_ttl_or_hl | max_ttl_or_hl |mean_ttl_or_hl | dev_ttl_or_hl |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT6 struct {
	BlockType uint8
	LossReportFlag bool
	DuplicateSupportFlag bool
	JitterFlag bool
	HopLimitFlag uint8
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	LostPackets uint32
	DuplicatePackets uint32
	MinimumJitter uint32
	MaximumJitter uint32
	MeanJitter uint32
	DeviationOfJitter uint32
	MinimumHopLimit uint8
	MaximumHopLimit uint8
	MeanHopLimit uint8
	DeviationOfHopLimit uint8
}

// unmarshal gets information from a packet
func (r *XRBT6) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT6) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT7
// ---------------------------------------------------------------------------

// The Statistics Summary Report Block has the following format:
// found on https://tools.ietf.org/html/rfc3611#section-4.7

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=7      |   reserved    |       block length = 8        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   loss rate   | discard rate  | burst density |  gap density  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       burst duration          |         gap duration          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     round trip delay          |       end system delay        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | signal level  |  noise level  |     RERL      |     Gmin      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   R factor    | ext. R factor |    MOS-LQ     |    MOS-CQ     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   RX config   |   reserved    |          JB nominal           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          JB maximum           |          JB abs max           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT7 struct {
	BlockType uint8
	BlockLength uint16
	SSRC uint32
	LossRate uint8
	DiscardRate uint8
	BurstDensity uint8
	GapDensity uint8
	GapDuration uint16
	RoundTripDelay uint16
	EndSystemDelay uint16
	SignalLevel uint8
	NoiseLevel uint8
	EchoReturnLoss uint8 // RERL
	GapThreshold uint8 // Gmin
	RFactor uint8
	MeanOpinionScoreListening uint8
	MeanOpinionScoreConversation uint8
// The next 3 values are in RX config (Receiver Configuration)
// https://tools.ietf.org/html/rfc3611#section-4.7.6  
	RXPacketLossConcealment uint8
	RXJitterBufferAdaptive uint8
	RXJItterBufferRate uint8
// -------------------------------------------------
	JitterNominalDelay uint16
	JitterMaximumDelay uint16
	JitterAbsMaximumDelay uint16
}

// unmarshal gets information from a packet
func (r *XRBT7) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT7) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT8
// ---------------------------------------------------------------------------


// Extended Network Quality (XNQ) Report Block
// https://tools.ietf.org/html/rfc5093

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=8      |   reserved    |      block length = 8         |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           vmaxdiff            |             vrange            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                              vsum                             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |               c               |            jbevents           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   reserved    |                     tdegnet                   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   reserved    |                     tdegjit                   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   reserved    |                        es                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   reserved    |                       ses                     |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+

type XRBT8 struct {
	BlockType uint8
	BlockLength uint16
	BeginningSequence uint16
	EndingSequence uint16
	IPDVdiff uint16
	IPDVrange uint16
	IPDVsum uint32
	RTCPcycles uint16 // c
	JitterAdaptationEvents uint16
	LossTotalSamplePeriods uint32
	JitterTotalSamplePeriods uint32
	UnavailablePacketsTime uint32
	UnavailablePacketsTimeSevere uint32
}

// unmarshal gets information from a packet
func (r *XRBT8) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT8) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT10
// ---------------------------------------------------------------------------

type XRBT9 struct {
	BlockType 
	BlockLength
}

// unmarshal gets information from a packet
func (r *XRBT9) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT9) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT10
// ---------------------------------------------------------------------------

// Post-Repair Loss RLE Report Block Type
// https://tools.ietf.org/html/rfc5725#section-3

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=10     | rsvd. |   T   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk 1              |             chunk 2           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                              ...                              :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk n-1            |             chunk n           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT10 struct {
	BlockType uint8
	Thinning uint8
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	Chunks []chunk
}

// unmarshal gets information from a packet
func (r *XRBT10) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT10) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT11
// ---------------------------------------------------------------------------

// Multicast Acquisition Report Block Type
// https://tools.ietf.org/html/rfc6332

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=11     |   MA Method   |         Block Length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             the Primary Multicast Stream             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             Status            |             Rsvd.             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT11 struct {
	BlockType uint8
	MAmethod uint8
	BlockLength uint16
	SRCC uint32
	Status uint16
// 	0                   1                   2                   3
// 	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Type      |   Reserved    |            Length             |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   :                             Value                             :
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+	
	TLVElements []XRBT11TLVElement // TLV means type-length-value

}

type XRBT11TLVELement struct { 
	Type uint8
	Length uint16
	Value []byte
}

// unmarshal gets information from a packet
func (r *XRBT11) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT11) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT12
// ---------------------------------------------------------------------------

// Inter-Destination Media Synchronization (IDMS) Block
// https://tools.ietf.org/html/rfc7272#section-6

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=12     | SPST  |Resrv|P|         block length=7        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     PT      |               Resrv                             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              Media Stream Correlation Identifier              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    media source                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Packet Received NTP timestamp, most significant word     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Packet Received NTP timestamp, least significant word    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              Packet Received RTP timestamp                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              Packet Presented NTP timestamp                   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT12 struct {
	BlockType uint8
	SPST uint8 // Synchronization Packet Sender Type
	PacketNPTTimeStampFlag uint8 // NOTE: might change to bool
	BlockLength uint16

//	https://tools.ietf.org/html/rfc3551#section-6
	PayloadType string // NOTE: not sure about type, link above has more information

	MediaStreamCorrelation uint32
	SSRC uint32
	NTPtimestamps []uint32 // NTP means Network Time Stamp
	RTPtimestamps []uint32 
}

// unmarshal gets information from a packet
func (r *XRBT12) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT12) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT13
// ---------------------------------------------------------------------------

// RTCP XR Report Block for Explicit Congestion Notification (ECN) Summary
// https://tools.ietf.org/html/rfc6679#section-5.2

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=13     | Reserved      |         Block Length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  Media Sender                                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | ECT (0) Counter                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | ECT (1) Counter                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | ECN-CE Counter                | not-ECT Counter               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | Lost Packets Counter          | Duplication Counter           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT13 struct {
	BlockType uint8
	BlockLength uint16
	SSRC uint32 // NOTE: 28 bits instead of 32 bits as seen in above diagram???
	ECT0Counter uint32 // ECT means ECN Capable Transport.
	ECT1Counter uint32
	ECNCECounter uint16 // CE means Congestion Experienced.
	NonECTCounter uint16
	LostPacketsCounter uint16
	DuplicationCounter uint16

// https://tools.ietf.org/html/rfc6679#section-6
// SDP Signalling Extensions for ECN
// NOTE: I have no idea what the SDP signaling extensions are, but after looking at the des-
//     cription for block length, I am assuming they are nothing, but don't know.
}

// unmarshal gets information from a packet
func (r *XRBT13) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT13) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT14
// ---------------------------------------------------------------------------

// Measurement Identity and Information Reporting Using a
//  Source Description (SDES)
// https://tools.ietf.org/html/rfc6776 

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=14     |    Reserved   |      block length = 7         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                   packet source                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |            Reserved           |    first sequence number      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           extended first sequence  interval          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                 extended last sequence number                 |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              Measurement Duration (Interval)                  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Measurement Duration (Cumulative) - Seconds (bit 0-31)    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Measurement Duration (Cumulative) - Fraction (bit 0-31)   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT14 struct {
	BlockType uint8
	BlockLength uint16
	SSRC uint32
	FirstSequence uint16
	ExtendedFirstInterval uint32
	ExtendedLastInterval uint32
	MeasurementDurationInterval uint32
	MeasurementDurationCumulative uint64
}

// unmarshal gets information from a packet
func (r *XRBT14) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT14) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT15
// ---------------------------------------------------------------------------

// Packet Delay Variation Metric Reporting
// https://tools.ietf.org/html/rfc6798#section-3.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=15     | I |pdvtyp |Rsv|       block length=4          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Pos PDV Threshold/Peak     |     Pos PDV Percentile        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Neg PDV Threshold/Peak     |     Neg PDV Percentile        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Mean PDV             |           Reserved            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT15 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	DelayVariationMetricType uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	PositivePDVThresholdPeak uint16
	PositivePDVPercentile uint16
	NegativePDVThresholdPeak uint16
	MeanPDV uint16
}

// unmarshal gets information from a packet
func (r *XRBT15) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT15) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT16
// ---------------------------------------------------------------------------

// Block for Delay Metric Reporting
// https://tools.ietf.org/html/rfc6843

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    BT=16      | I |   resv.   |      block length = 6         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Source                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                  Mean Network Round-Trip Delay                |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                   Min Network Round-Trip Delay                |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                   Max Network Round-Trip Delay                |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |               End System Delay - Seconds (bit 0-31)           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              End System Delay - Fraction (bit 0-31)           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT16 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	MeanNetworkRoundTripDelay uint32
	MinNetworkRoundTripDelay uint32
	MaxNetworkRoundTripDelay uint32
	EndSystemDelay uint64
}

// unmarshal gets information from a packet
func (r *XRBT16) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT16) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT17
// ---------------------------------------------------------------------------

// Summary Statistics Metrics Reporting - Burst/Gap Loss Summary Statistics Block
// https://tools.ietf.org/html/rfc7004#section-3.1.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      BT=17    | I | Reserved  |        Block Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Source                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Burst Loss Rate        |         Gap Loss Rate         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Burst Duration Mean     |    Burst Duration Variance    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT17 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	BurstLossRate uint16
	GapLossRate uint16
	BurstDurationMean uint16
	BurstDurationVariance uint16
}

// unmarshal gets information from a packet
func (r *XRBT17) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT17) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT18
// ---------------------------------------------------------------------------

// Summary Statistics Metrics Reporting - Burst/Gap Discard Summary Statistics Block
// https://tools.ietf.org/html/rfc7004#section-3.2.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      BT=18    | I |  Reserved |        Block Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                         Source                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Burst Discard Rate   |        Gap Discard Rate       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT18 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	BurstDiscardRate uint16
	GapDiscardRate uint16
}

// unmarshal gets information from a packet
func (r *XRBT18) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT18) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT19
// ---------------------------------------------------------------------------

// Summary Statistics Metrics Reporting - Frame Impairment Statistics Summary Block
// https://tools.ietf.org/html/rfc7004#section-4.1.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      BT=19    |T|   Reserved  |        Block Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      discarded_frames                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          dup_frames                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      full_lost_frames                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      partial_lost_frames                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT19 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	DiscardedFrames uint32
	DuplicatedFrames uint32
	FullLostFrames uint32
	PartialLostFrames uint32
}

// unmarshal gets information from a packet
func (r *XRBT19) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT19) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT20
// ---------------------------------------------------------------------------

// Burst/Gap Loss Metric Reporting
// https://tools.ietf.org/html/rfc6958#section-3.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=20     | I |C|  resv.  |      Block Length = 5         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ -+-+-+-+-+-+-+-+
// | Threshold     |       Burst Durations (ms)           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |            Packets Lost in Bursts             |    Total...   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | ...Packets Expected in Bursts |     Bursts  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                ...S Burst Durations (ms-squared)     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT20 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	LossAndDiscardComboFLag uint8 // NOTE: might change to bool
	BlockLength uint8
	SSRC uint32
	Threshold uint8
	BurstDurationsSum uint32
	PacketLostBursts uint32
	TotalPacketsExpectedBursts uint32
	Bursts uint16
	SumOfSquaresBurstDurations uint32
}

// unmarshal gets information from a packet
func (r *XRBT20) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT20) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT21
// ---------------------------------------------------------------------------

// Burst/Gap Discard Metric Reporting

// block type (BT) in link below is corrected by the subsequent link (should be 21)

// https://tools.ietf.org/html/rfc7003#section-3.1
// https://www.rfc-editor.org/errata_search.php?eid=3735

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=21     | I |   resv    |      Block Length = 3         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Threshold   |         Packets Discarded in Bursts           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Total Packets Expected in Bursts        |   Reserved    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT21 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	Threshold uint8
	PacketsDiscardedBursts uint32
	TotalPacketsExpectedBursts uint32
}

// unmarshal gets information from a packet
func (r *XRBT21) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT21) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT22
// ---------------------------------------------------------------------------

// Block for MPEG-2
//   Transport Stream (TS) Program Specific Information (PSI) Independent
//                Decodability Statistics Metrics Reporting


// https://tools.ietf.org/html/rfc6990#section-1.4

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=22     |    Reserved   |         Block Length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Source                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      TS_sync_loss_count                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Sync_byte_error_count                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                  Continuity_count_error_count                 |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Transport_error_count                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        PCR_error_count                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                   PCR_repetition_error_count                  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |            PCR_discontinuity_indicator_error_count            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     PCR_accuracy_error_count                  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       PTS_error_count                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT22 struct {
	BlockType uint8
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	TSSyncLossCount uint32
	SyncByteErrorCount uint32
	ContinuityCountErrorCount uint32
	TransportErrorCount uint32
	PCRErrorCount uint32 // PCR means Primary Clock Reference
	PCRRepetitionErrorCount uint32
	PCRDiscontinuityIndicatorErrorCount uint32
	PCRAccuracyErrorCount uint32
	PTSErrorCount uint32 // PTS means Presentation Time Stamp
}

// NOTE: there is SDP signaling in this report

// unmarshal gets information from a packet
func (r *XRBT22) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT22) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT23
// ---------------------------------------------------------------------------

// De-Jitter Buffer Metric Reporting
// https://tools.ietf.org/html/rfc7005#section-4.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=23    | I |C|  resv    |       Block Length=3          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Source                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          DJB nominal          |        DJB maximum            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     DJB high-water mark       |      DJB low-water mark       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+



type XRBT23 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	LossAndDiscardComboFLag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	DejitterMaximumDelay uint16
	DejitterHighWaterMark uint16
	DejitterBufferLowWaterMark uint16
}

// NOTE: there is SDP signaling in this report

// unmarshal gets information from a packet
func (r *XRBT23) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT23) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT24
// ---------------------------------------------------------------------------

// Discard Count Metric Reporting

// https://tools.ietf.org/html/rfc7002#section-3.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=24     | I |DT |  resv |      Block Length = 2         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Discard Count                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT24 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	DiscardType uint8
	BlockLength uint16
	SSRC uint16
	DiscardCount uint32
}

// NOTE: there is SDP signaling in this report

// unmarshal gets information from a packet
func (r *XRBT24) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT24) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT25
// ---------------------------------------------------------------------------

// Discarded Packets
// https://tools.ietf.org/html/rfc7097#section-3

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=25     |rsvd |E|   T   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk 1              |             chunk 2           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// :                              ...                              :
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          chunk n-1            |             chunk n           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT25 struct {
	BlockType uint8
	EarlyPacketsDiscardedFlag uint8 // NOTE: might change to bool
	Thinning uint8 // NOTE: There is no documentation on this and the types below other than above diagram (so this is a guess)
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	Chunks []chunk
}

// NOTE: SDP signaling is optional

// unmarshal gets information from a packet
func (r *XRBT25) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT25) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT26
// ---------------------------------------------------------------------------

// Bytes Discarded Metric Report Block
// https://tools.ietf.org/html/rfc7243#section-3

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=26     | I |E|Reserved |       Block length=2          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              RTP payload bytes discarded             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+



type XRBT26 struct {
	BlockType uint8
	IntervalMetricFlag uint8
	EarlyPacketsDiscardedFlag uint8 // NOTE: might change to bool
	SSRC uint32
	DiscardedRTPPayloadBytes uint32
}

// SDP signaling is optional

// unmarshal gets information from a packet
func (r *XRBT26) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT26) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT27
// ---------------------------------------------------------------------------

// Synchronization Delay and Offset Metrics Reporting
// https://tools.ietf.org/html/rfc7244#section-3.1



// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=27     |   Reserved    |         Block length=2        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     Source                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |               Initial Synchronization Delay                   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT27 struct {
	BlockType uint8
	BlockLength uint8
	SSRC uint32
	InitialSyncDelay uint32
}

// unmarshal gets information from a packet
func (r *XRBT27) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT27) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT28
// ---------------------------------------------------------------------------

// Flow General Synchronization Offset Metrics Block
// https://tools.ietf.org/html/rfc7244#section-4.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    BT=28      | I | Reserved  |         Block length=3        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Synchronization Offset, most significant word         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Synchronization Offset, least significant word        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT28 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint8
	SSRC uint32
	SyncOffset uint64
}

// The is SDP signaling here

// unmarshal gets information from a packet
func (r *XRBT28) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT28) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT29
// ---------------------------------------------------------------------------

// Mean Opinion Score (MOS) Metric Reporting Block
// https://tools.ietf.org/html/rfc7266#section-3.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=29     | I |  Reserved |       Block Length            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Segment  1                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Segment 2                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// ..................
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Segment n                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+


// 3.2.1.  Single-Channel Audio/Video per SSRC Segment

//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |S|     CAID      |    PT       |           MOS Value           |
// 	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	

type segmentSingleChannel struct {
	SegmentType uint8
	CalculationAlgorithm uint8
	PayloadType uint8
	MOSValue uint16 //  MOS means mean opinion score
}

// 3.2.2.  Multi-Channel Audio per SSRC Segment

//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |S|     CAID      |    PT       |CHID |        MOS Value        |
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type segmentMultiChannel struct {
	SegmentType uint8
	CalculationAlgorithm uint8
	PayloadType uint8
	ChannelIdentifier uint8
	MOSValue uint16 //  MOS means mean opinion score
}

type XRBT29 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint16
	SSRC uint32
	Segments []struct
}



// unmarshal gets information from a packet
func (r *XRBT29) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {
	// if channelType == "multi" {
	// 	newSegment = segmentMultiChannel{ add struct fields for values above } 
	// }
	// if channelType == "single" {
	// 	newSegment := segmentSingleChannel{ add stuct fields for values above } 
	// }	
	// 
   // r.Segments = r.Segments.append(newSegment)
	// 
}

// marshal puts information into a packet
func (r *XRBT29) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT30
// ---------------------------------------------------------------------------

// Blocks for Concealment Metrics Reporting on Audio Applications - Loss Concealment Metrics Block
// https://tools.ietf.org/html/rfc7294#section-3.1

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    BT=30      | I |plc|  resv |       block length=6          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Source                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                 On-Time Playout Duration                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                  Loss Concealment Duration                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              Buffer Adjustment Concealment Duration           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Playout Interrupt Count    |           Reserved            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                 Mean Playout Interrupt Size                   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT30 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	PacketLossConcealmentMethod uint8
	BlockLength uint8
	SSRC uint32
}

// unmarshal gets information from a packet
func (r *XRBT30) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT30) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT31
// ---------------------------------------------------------------------------

// Blocks for Concealment Metrics Reporting on Audio Applications - Concealed Seconds Metrics
// https://tools.ietf.org/html/rfc7294#section-4.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    BT=31      | I |plc|  resv |       block length=4          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Source                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Unimpaired Seconds                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Concealed Seconds                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// | Severely Concealed Seconds    | Reserved      | SCS Threshold |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT31 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	PacketLossConcealmentMethod uint8
	BlockLength uint8
	SSRC uint32
	UnimpairedSeconds uint32
	ConcealedSeconds uint32
	SeverelyConcealedSeconds uint16
	SCSThreshold uint8
}

// SDP siginaling is optional

// unmarshal gets information from a packet
func (r *XRBT31) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT31) marshal() ([]byte, error) {

}

// ---------------------------------------------------------------------------
// Code below is for XRBT32
// ---------------------------------------------------------------------------

// Block for MPEG2
//  Transport Stream (TS) Program Specific Information (PSI) Decodability
//                       Statistics Metrics Reporting
// https://tools.ietf.org/html/rfc7380#section-3

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      BT=32    |    Reserved   |         block length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    source                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          begin_seq            |             end_seq           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        PAT_error_count        |      PAT_error_2_count        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        PMT_error_count        |      PMT_error_2_count        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       PID_error_count         |      CRC_error_count          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        CAT_error_count        |        Reserved               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT32 struct {
	BlockType uint8
	BlockLength uint16
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	PATErrorCount uint16 // PAT means Program Association Table
	PATError2Count uint16
	PMTErrorCount uint16
	PMTError2Count uint16
	PIDErrorCount uint16
	CRCErrorCount uint16
	CATErrorCount uint16
}

// SDP signaling is optional

// unmarshal gets information from a packet
func (r *XRBT32) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT32) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT33
// ---------------------------------------------------------------------------

// Post-Repair Loss Count Metrics
// https://tools.ietf.org/html/rfc7509#section-3.1

// 0               1               2               3               4
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=33     |   Reserved    |      Block length = 4         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Source                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       begin_seq               |          end_seq              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  Post-repair loss count       |     Repaired loss count       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT33 struct {
	BlockType uint8
	BlockLength uint8
	SSRC uint32
	BeginningSequence uint16
	EndingSequence uint16
	PostRepairLossCount uint16
	RepairedLossCount uint16
}

// This report uses SDP signaling.

// unmarshal gets information from a packet
func (r *XRBT33) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT33) marshal() ([]byte, error) {

}



// ---------------------------------------------------------------------------
// Code below is for XRBT34
// ---------------------------------------------------------------------------

// Loss Concealment Metrics for Video Applications Block
// https://tools.ietf.org/html/rfc7867#section-4

// 0               1               2               3
// 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    BT=34      | I | V |  RSV  |       Block Length            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Source                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Impaired Duration                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Concealed Duration                       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                  Mean Frame Freeze Duration (optional)        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    MIFP       |    MCFP       |     FFSC      |     Reserved  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT34 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	VideoLossConcealmentMethodType uint8
	BlockLength uint8
	SSRC uint32
	ImpairedDuration uint32
	ConcealmentDuration uint32
	MainFrameFreezeDuration uint32
	MainImpairedFrameProportion uint8
	FractionOfFramesSubjectToConcealment uint8
}

// There is SDP signaling with this report

// unmarshal gets information from a packet
func (r *XRBT34) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT34) marshal() ([]byte, error) {

}


// ---------------------------------------------------------------------------
// Code below is for XRBT35
// ---------------------------------------------------------------------------

// Independent Rep Burst/Gap Discard Metrics
// https://tools.ietf.org/html/rfc8015#section-3.1

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     BT=35     | I |   resv    |      Block Length = 5         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Threshold   |       Burst Durations (ms)           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Packets Discarded in Bursts          |      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Bursts     |           Total Packets Expected in Bursts    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        Discard Count                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type XRBT35 struct {
	BlockType uint8
	IntervalMetricFlag uint8 // NOTE: might change to bool
	BlockLength uint8
	Source uint32
	Threshold uint8
	BurstDurations uint
	PacketsDiscardedBursts uint32
	Bursts uint8
	TotalPacketsExpectedBursts uint32
	DiscardCount uint32


}

// unmarshal gets information from a packet
func (r *XRBT35) unmarshal(rawData []byte, XRBytesProcessed int) (XRBytesProcessed  error {

}

// marshal puts information into a packet
func (r *XRBT35) marshal() ([]byte, error) {

}