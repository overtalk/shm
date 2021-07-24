package ishm

import (
	"log"
)

type CreateSHMParam struct {
	Key int64
	Size int64
}
type UpdateContent struct {
	EventType int16
	Topic string
	Content string
}
func UpdateCtx(shmparam CreateSHMParam, updatectx UpdateContent) (index int, err error){

	sm,err:= CreateWithKey(shmparam.Key,shmparam.Size)
	if err != nil {
		log.Fatal(err)
		return index,err
	}
	log.Print(sm)
	pos ,err:=sm.WriteObj(updatectx)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(pos)

	return pos,err
}
func GetCtx(shmparam CreateSHMParam) ( updatectx*UpdateContent, err error){

	sm,err:= CreateWithKey(shmparam.Key,0)
	if err != nil {
		log.Fatal(err)
		return updatectx,err
	}
	log.Print(sm)
	pos ,err:=sm.ReadObjCtx(updatectx)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(pos)
	return updatectx,err
}
