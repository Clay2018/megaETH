package bitcoin

type UTXO struct {
	Hash        string
	Index       uint32
	Value       int64
	PkScriptHex string
}
