package network

type UmbrellaRequest struct {
	Sn string
	Success bool
	Err string
}

func UmbrellaRequestTimeout() UmbrellaRequest {
	return UmbrellaRequest{
		Success: false,
		Err: "超时",
		Sn: "",
	}
}