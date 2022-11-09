package common

import "net/rpc"

func Call(srv, rpcname string, args, reply interface{}) error {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return errx
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	return err
}

func (e Error) Error() string {
	return e.Err
}
