package proxy

func findSNI(data []byte) (offset int, length int, ok bool) {
	if len(data) < 9 {
		return 0, 0, false
	}
	if data[0] != 0x16 || data[5] != 0x01 {
		return 0, 0, false
	}
	// pos=9: [Version][Random]
	pos := 9
	pos += 34 // pos=43: [SessionIDLen]

	sessionIDLen := int(data[pos])
	pos++               // pos=44: [SessionID...]
	pos += sessionIDLen // pos=44+N: [CipherSuitesLen]

	CipherSuitesLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2               // pos=46+N: [CipherSuites...]
	pos += CipherSuitesLen // pos=46+N+M: [CompressionLen]

	compressionMethodsLen := int(data[pos])
	pos++                        // pos=47+N+M: [Compression...]
	pos += compressionMethodsLen // pos=47+N+M+K: [ExtensionsLen]

	extensionsLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2 // pos=49+N+M+K: [Ext1Type][Ext1Len]...
	end := pos + extensionsLen
	for pos < end {
		extType := int(data[pos])<<8 | int(data[pos+1])
		extLen := int(data[pos+2])<<8 | int(data[pos+3])
		pos += 4 // pos: [ExtData...]

		if extType == 0x0000 {
			pos += 2 // pos: [ServerNameType]
			pos++    // pos: [ServerNameLength]
			nameLen := int(data[pos])<<8 | int(data[pos+1])
			pos += 2 // pos: [ServerName] ← домен

			return pos, nameLen, true
		}

		pos += extLen // pos: [NextExtType]
	}
	return 0, 0, false
}
