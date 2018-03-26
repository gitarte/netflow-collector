package main

import "strconv"

// NetFlowV5 is representing single NetFlow transmission
type NetFlowV5 struct {
	Header  Header   `json:"Header"`
	Records []Record `json:"Records"`
	Raw     string   `json:"Raw"`
}

// Header structure of header NetFlow v5
type Header struct {
	Version          uint16 `json:"Version"`          // 00-01	version				NetFlow export format version number
	Count            uint16 `json:"Count"`            // 02-03	count				Number of flows exported in this packet (1-30)
	SysUptime        int32  `json:"sysUptime"`        // 04-07	sys_uptime			Current time in milliseconds since the export device booted
	Timestamp        string `json:"Timestamp"`        // 08-11	unix_secs			Current count of seconds since 0000 UTC 1970 and 12-15	unix_nsecs	Residual nanoseconds since 0000 UTC 1970
	FlowSequence     int32  `json:"FlowSequence"`     // 16-19	flow_sequence		Sequence counter of total flows seen
	EngineType       uint8  `json:"EngineType"`       // 20		engine_type			Type of flow-switching engine
	EngineID         uint8  `json:"EngineID"`         // 21		engine_id			Slot number of the flow-switching engine
	SamplingInterval uint16 `json:"SamplingInterval"` // 22-23	sampling_interval	First two bits hold the sampling mode; remaining 14 bits hold value of sampling interval
}

// Record structure of single record in NetFlow v5 peyload
type Record struct {
	SrcAddr string `json:"SrcAddr"` // 00-03	srcaddr	Source IP address
	DstAddr string `json:"DstAddr"` // 04-07	dstaddr	Destination IP address
	NextHop string `json:"NextHop"` // 08-11	nexthop	IP address of next hop router
	Input   uint16 `json:"Input"`   // 12-13	input	SNMP index of input interface
	Output  uint16 `json:"Output"`  // 14-15	output	SNMP index of output interface
	DPkts   uint32 `json:"DPkts"`   // 16-19	dPkts	Packets in the flow
	DOctets uint32 `json:"DOctets"` // 20-23	dOctets	Total number of Layer 3 bytes in the packets of the flow
	First   uint32 `json:"First"`   // 24-27	first	SysUptime at start of flow
	Last    uint32 `json:"Last"`    // 28-31	last	SysUptime at the time the last packet of the flow was received
	SrcPort uint16 `json:"SrcPort"` // 32-33	srcport	TCP/UDP source port number or equivalent
	DstPort uint16 `json:"DstPort"` // 34-35	dstport	TCP/UDP destination port number or equivalent
	// 36	pad1	Unused (zero) bytes
	TCPFlags uint8  `json:"TCPFlags"` // 37		tcp_flags	Cumulative OR of TCP flags
	Prot     uint8  `json:"Prot"`     // 38		prot		IP protocol type (for example, TCP = 6; UDP = 17)
	Tos      uint8  `json:"Tos"`      // 39		tos	IP 		type of service (ToS)
	SrcAs    uint16 `json:"SrcAs"`    // 40-41	src_as		Autonomous system number of the source, either origin or peer
	DstAs    uint16 `json:"DstAs"`    // 42-43	dst_as		Autonomous system number of the destination, either origin or peer
	SrcMask  uint8  `json:"SrcMask"`  // 44		src_mask	Source address prefix mask bits
	DstMask  uint8  `json:"DstMask"`  // 45		dst_mask	Destination address prefix mask bits
	// 46-47	pad2	Unused (zero) bytes
}

// DecodeHeader extracts header fields from given slice of bytes
func DecodeHeader(buf []byte) Header {
	var h Header
	h.Version = uint16(buf[0])*256 + uint16(buf[1]) //	NetFlow export format version number
	h.Count = uint16(buf[2])*256 + uint16(buf[3])   //	Number of flows exported in this packet (1-30)
	h.SysUptime = int32(buf[4])*256*256*256 +       //	Current time in milliseconds since the export device booted
		int32(buf[5])*256*256 +
		int32(buf[6])*256 +
		int32(buf[7])
	h.Timestamp = GetTimestamp(
		buf[8:12],  //	Current count of seconds since 0000 UTC 1970
		buf[12:16]) //	Residual nanoseconds since 0000 UTC 1970
	h.FlowSequence = int32(buf[16])*256*256*256 + //	Sequence counter of total flows seen
		int32(buf[17])*256*256 +
		int32(buf[18])*256 +
		int32(buf[19])
	h.EngineType = uint8(buf[20])                          //	Type of flow-switching engine
	h.EngineID = uint8(buf[21])                            //	Slot number of the flow-switching engine
	h.SamplingInterval = uint16(buf[22]) + uint16(buf[24]) //	First two bits hold the sampling mode; remaining 14 bits hold value of sampling interval
	return h
}

// DecodeRecord extracts single record of number c fields from given slice of bytes
func DecodeRecord(buf []byte, c uint16) Record {
	var r Record

	recLen := uint16(48)
	recBuf := buf[c*recLen+24 : c*recLen+24+recLen]

	// 00-03	srcaddr	Source IP address
	r.SrcAddr = strconv.Itoa(int(recBuf[0])) + "." +
		strconv.Itoa(int(recBuf[1])) + "." +
		strconv.Itoa(int(recBuf[2])) + "." +
		strconv.Itoa(int(recBuf[3]))

	// 04-07	dstaddr	Destination IP address
	r.DstAddr = strconv.Itoa(int(recBuf[4])) + "." +
		strconv.Itoa(int(recBuf[5])) + "." +
		strconv.Itoa(int(recBuf[6])) + "." +
		strconv.Itoa(int(recBuf[7]))

	// 08-11	nexthop	IP address of next hop router
	r.NextHop = strconv.Itoa(int(recBuf[8])) + "." +
		strconv.Itoa(int(recBuf[9])) + "." +
		strconv.Itoa(int(recBuf[10])) + "." +
		strconv.Itoa(int(recBuf[11]))

	// 12-13	input	SNMP index of input interface
	r.Input = uint16(recBuf[12])*256 + uint16(recBuf[13])

	// 14-15	output	SNMP index of output interface
	r.Output = uint16(recBuf[14])*256 + uint16(recBuf[15])

	// 16-19	dPkts	Packets in the flow
	r.DPkts = uint32(recBuf[16])*256*256*256 +
		uint32(recBuf[17])*256*256 +
		uint32(recBuf[18])*256 +
		uint32(recBuf[19])

	// 20-23	dOctets	Total number of Layer 3 bytes in the packets of the flow
	r.DOctets = uint32(recBuf[20])*256*256*256 +
		uint32(recBuf[21])*256*256 +
		uint32(recBuf[22])*256 +
		uint32(recBuf[23])

	// 24-27	first	SysUptime at start of flow
	r.First = uint32(recBuf[24])*256*256*256 +
		uint32(recBuf[25])*256*256 +
		uint32(recBuf[26])*256 +
		uint32(recBuf[27])

	// 28-31	last	SysUptime at the time the last packet of the flow was received
	r.Last = uint32(recBuf[28])*256*256*256 +
		uint32(recBuf[29])*256*256 +
		uint32(recBuf[30])*256 +
		uint32(recBuf[31])

	// 32-33	srcport		TCP/UDP source port number or equivalent
	r.SrcPort = uint16(recBuf[32])*256 + uint16(recBuf[33])

	// 34-35	dstport		TCP/UDP destination port number or equivalent
	r.DstPort = uint16(recBuf[34])*256 + uint16(recBuf[35])

	// 37		tcp_flags	Cumulative OR of TCP flags
	r.TCPFlags = uint8(recBuf[37])

	// 38		prot		IP protocol type (for example, TCP = 6; UDP = 17)
	r.Prot = uint8(recBuf[38])

	// 39		tos			IP type of service (ToS)
	r.Tos = uint8(recBuf[39])

	// 40-41	src_as		Autonomous system number of the source, either origin or peer
	r.SrcAs = uint16(recBuf[40])*256 + uint16(recBuf[41])

	// 42-43	dst_as		Autonomous system number of the destination, either origin or peer
	r.DstAs = uint16(recBuf[42])*256 + uint16(recBuf[43])

	// 44		src_mask	Source address prefix mask bits
	r.SrcMask = uint8(recBuf[44])

	// 45		dst_mask	Destination address prefix mask bits
	r.DstMask = uint8(recBuf[45])

	return r
}