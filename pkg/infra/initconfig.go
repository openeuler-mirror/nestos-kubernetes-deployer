package infra

type InitConfig struct{
	OsType string
}


type BootConfigAssembler interface {
	Assemble(assets Assets) InitConfig
}
