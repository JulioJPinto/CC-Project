package client_tracer_protocol

type CTPHeader struct {
	flags struct {
		who_has bool
		i_have  bool
	}
}

type whoHasRequestBody struct {
	file_name   string
	byte_offset uint64
	byte_length uint64
}

func whoHasRequest()

func Header() {

}
